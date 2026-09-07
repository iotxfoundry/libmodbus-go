package modbus

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// TCPServer is a pure-Go multi-client Modbus TCP server.
// It accepts connections via net.Listener, dup's each connection's fd into
// blocking mode, and hands the fd to the C libmodbus context for recv/send.
type TCPServer struct {
	addr       string
	mapping    *ModbusMapping
	maxClients int
	debug      bool
	slaveID    int

	onConnect    func(addr net.Addr)
	onDisconnect func(addr net.Addr, err error)

	mu       sync.Mutex
	closed   bool
	listener net.Listener
	files    map[int]*os.File // dupFd -> *os.File (prevent GC of fd)
	wg       sync.WaitGroup
	done     chan struct{}
	sem      chan struct{} // semaphore for max clients
}

// NewTCPServer creates a new Modbus TCP server that listens on addr and
// serves requests against the shared data mapping.
func NewTCPServer(addr string, mapping *ModbusMapping) *TCPServer {
	return &TCPServer{
		addr:       addr,
		mapping:    mapping,
		maxClients: 10,
		slaveID:    -1, // -1 means use default (no SetSlave call)
		files:      make(map[int]*os.File),
		done:       make(chan struct{}),
	}
}

// SetMaxClients sets the maximum number of simultaneous client connections.
func (s *TCPServer) SetMaxClients(n int) {
	s.maxClients = n
}

// SetDebug enables or disables verbose debug output for all client contexts.
func (s *TCPServer) SetDebug(debug bool) {
	s.debug = debug
}

// SetOnConnect registers a callback invoked when a new client connects.
func (s *TCPServer) SetOnConnect(fn func(net.Addr)) {
	s.onConnect = fn
}

// SetOnDisconnect registers a callback invoked when a client disconnects.
// The error argument describes the reason for disconnection.
func (s *TCPServer) SetOnDisconnect(fn func(net.Addr, error)) {
	s.onDisconnect = fn
}

// SetSlaveID sets the Modbus slave/unit ID for all client connections.
// By default, the server uses the libmodbus default (MODBUS_TCP_SLAVE = 0xFF),
// which responds to all unit IDs. Set a specific ID to filter requests.
func (s *TCPServer) SetSlaveID(id int) {
	s.slaveID = id
}

// Serve starts the TCP server. It blocks until Close is called or a fatal
// accept error occurs.
func (s *TCPServer) Serve() error {
	return s.ServeContext(context.Background())
}

// ServeContext starts the TCP server. It blocks until ctx is canceled, Close
// is called, or a fatal accept error occurs. When ctx is canceled, in-flight
// client connections are closed and ServeContext returns ctx.Err().
func (s *TCPServer) ServeContext(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("modbus: tcp server listen %s: %w", s.addr, err)
	}
	s.mu.Lock()
	s.listener = ln
	s.mu.Unlock()

	s.sem = make(chan struct{}, s.maxClients)

	// Unblock Accept when ctx is canceled, the server is closed, or
	// ServeContext returns (e.g. fatal accept error).
	serveDone := make(chan struct{})
	defer close(serveDone)
	go func() {
		select {
		case <-ctx.Done():
			ln.Close()
		case <-s.done:
		case <-serveDone:
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			default:
				return fmt.Errorf("modbus: tcp server accept: %w", err)
			}
		}

		// Try to acquire a semaphore slot (non-blocking).
		select {
		case s.sem <- struct{}{}:
			s.wg.Add(1)
			go s.handleConn(ctx, conn)
		default:
			conn.Close()
		}
	}
}

// Close gracefully shuts down the server. It stops accepting new connections,
// closes all dup'd file descriptors, and waits for all handler goroutines
// to finish. It is safe to call Close multiple times.
func (s *TCPServer) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	close(s.done)
	var err error
	if s.listener != nil {
		err = s.listener.Close()
	}
	for fd, file := range s.files {
		// shutdown 唤醒阻塞在同一 fd 上的 C 层 recv (close 不会唤醒阻塞读),
		// 使 handleConn 立即退出, Close 的 wg.Wait 不必等对端下一拍请求.
		_ = shutdownSocket(fd)
		file.Close()
		delete(s.files, fd)
	}
	s.mu.Unlock()

	s.wg.Wait()
	return err
}

// handleConn manages a single client connection. It dup's the TCP fd into
// blocking mode, creates a per-connection C libmodbus context, and runs a
// receive/reply loop until the client disconnects, ctx is canceled, or an
// error occurs.
func (s *TCPServer) handleConn(ctx context.Context, conn net.Conn) (err error) {
	defer s.wg.Done()
	defer func() { <-s.sem }()

	remoteAddr := conn.RemoteAddr()

	// Extract the raw fd via *net.TCPConn.File (duplicates the fd in blocking mode).
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		conn.Close()
		return fmt.Errorf("modbus: expected *net.TCPConn, got %T", conn)
	}

	file, err := tcpConn.File()
	if err != nil {
		conn.Close()
		return fmt.Errorf("modbus: get tcp conn file: %w", err)
	}
	dupFd := int(file.Fd())

	// Close the original conn; the dup'd fd keeps the socket alive.
	conn.Close()

	// Register the dup'd file to prevent GC from closing the fd.
	s.mu.Lock()
	s.files[dupFd] = file
	s.mu.Unlock()

	// Shutdown the connection's fd when ctx is canceled to unblock the C receive.
	// close(2) 不会唤醒阻塞在 recv 上的读线程, 必须 shutdown(2) 才能立即解除阻塞.
	connDone := make(chan struct{})
	defer close(connDone)
	go func() {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			if f, ok := s.files[dupFd]; ok {
				_ = shutdownSocket(dupFd)
				f.Close()
			}
			s.mu.Unlock()
		case <-connDone:
		}
	}()

	defer func() {
		s.mu.Lock()
		delete(s.files, dupFd)
		s.mu.Unlock()
		_ = shutdownSocket(dupFd)
		file.Close()
		if s.onDisconnect != nil {
			s.onDisconnect(remoteAddr, err)
		}
	}()

	// Create a per-connection C libmodbus context.
	mb, err := NewTCP("0.0.0.0", 0)
	if err != nil {
		return fmt.Errorf("modbus: new tcp context: %w", err)
	}
	defer mb.Free()

	if s.debug {
		mb.SetDebug(true)
	}

	if s.slaveID >= 0 {
		if err := mb.SetSlave(s.slaveID); err != nil {
			return fmt.Errorf("modbus: set slave: %w", err)
		}
	}

	if err := mb.SetSocket(dupFd); err != nil {
		return fmt.Errorf("modbus: set socket: %w", err)
	}

	if s.onConnect != nil {
		s.onConnect(remoteAddr)
	}

	// Receive/Reply loop.
	for {
		req, recvErr := mb.Receive()
		if recvErr != nil {
			err = recvErr
			return
		}
		replyErr := mb.Reply(req, s.mapping)
		if replyErr != nil {
			if s.debug {
				log.Printf("modbus: reply error: %s", replyErr)
			}
			err = replyErr
			return
		}
	}
}
