package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib/libmodbus.a
#include <stdlib.h>
#include "modbus.h"
*/
import "C"
import (
	"unsafe"
)

// ModbusNewTcp modbus_new_tcp - create a libmodbus context for TCP/IPv4
//
// The modbus_new_tcp() function shall allocate and initialize a modbus_t structure to communicate with a Modbus TCP
// IPv4 server.
//
// The ip argument specifies the IP address of the server to which the client wants to establish a connection. A NULL
// value can be used to listen any addresses in server mode.
//
// The port argument is the TCP port to use. Set the port to MODBUS_TCP_DEFAULT_PORT to use the default one (502). It's
// convenient to use a port number greater than or equal to 1024 because it's not necessary to have administrator
// privileges.
func ModbusNewTcp(addr string, port int) *Modbus {
	saddr := C.CString(addr)
	defer C.free(unsafe.Pointer(saddr))
	ctx := C.modbus_new_tcp(saddr, C.int(port))
	if ctx == nil {
		return nil
	}
	return &Modbus{ctx: ctx}
}

// TcpListen modbus_tcp_listen - create and listen a TCP Modbus socket (IPv4)
//
// The modbus_tcp_listen() function shall create a socket and listen to maximum nb_connection incoming connections on
// the specified IP address. The context ctx must be allocated and initialized with modbus_new_tcp before to set the IP
// address to listen, if IP address is set to NULL or '0.0.0.0', any addresses will be listen.
func (x *Modbus) TcpListen(nb int) (socket int, err error) {
	code := C.modbus_tcp_listen(x.ctx, C.int(nb))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	socket = int(code)
	x.socket = socket
	return
}

// TcpAccept modbus_tcp_accept - accept a new connection on a TCP Modbus socket (IPv4)
//
// The modbus_tcp_accept() function shall extract the first connection on the queue of pending connections, create a
// new socket and store it in libmodbus context given in argument. If available, accept4() with SOCK_CLOEXEC will be
// called instead of accept().
func (x *Modbus) TcpAccept() (err error) {
	code := C.modbus_tcp_accept(x.ctx, (*C.int)(unsafe.Pointer(&x.socket)))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	return
}

// ModbusNewTcpPi modbus_new_tcp_pi - create a libmodbus context for TCP Protocol Independent
//
// The modbus_new_tcp_pi() function shall allocate and initialize a modbus_t structure to communicate with a Modbus TCP
// IPv4 or IPv6 server.
//
// The node argument specifies the host name or IP address of the host to connect to, eg. "192.168.0.5" , "::1" or
// "server.com". A NULL value can be used to listen any addresses in server mode.
//
// The service argument is the service name/port number to connect to. To use the default Modbus port, you can provide
// an NULL value or the string "502". On many Unix systems, it's convenient to use a port number greater than or equal
// to 1024 because it's not necessary to have administrator privileges.
//
// v3.1.8 handles NULL value for service (no EINVAL error).
func ModbusNewTcpPi(node string, service string) *Modbus {
	cnode := C.CString(node)
	defer C.free(unsafe.Pointer(cnode))
	cservice := C.CString(service)
	defer C.free(unsafe.Pointer(cservice))
	ctx := C.modbus_new_tcp_pi(cnode, cservice)
	if ctx == nil {
		return nil
	}
	return &Modbus{ctx: ctx}
}

// TcpPiListen modbus_tcp_pi_listen - create and listen a TCP PI Modbus socket (IPv6)
//
// The modbus_tcp_pi_listen() function shall create a socket and listen to maximum nb_connection incoming connections
// on the specified nodes. The context ctx must be allocated and initialized with modbus_new_tcp_pi before to set the
// node to listen, if node is set to NULL or '0.0.0.0', any addresses will be listen.
func (x *Modbus) TcpPiListen(nb int) (socket int, err error) {
	code := C.modbus_tcp_pi_listen(x.ctx, C.int(nb))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	socket = int(code)
	x.socket = socket
	return
}

// TcpPiAccept modbus_tcp_pi_accept - accept a new connection on a TCP PI Modbus socket (IPv6)
//
// The modbus_tcp_pi_accept() function shall extract the first connection on the queue of pending connections, create a
// new socket and store it in libmodbus context given in argument. If available, accept4() with SOCK_CLOEXEC will be
// called instead of accept().
func (x *Modbus) TcpPiAccept() (err error) {
	code := C.modbus_tcp_pi_accept(x.ctx, (*C.int)(unsafe.Pointer(&x.socket)))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	return
}

// TcpReceive modbus_receive - receive an indication request
//
// The modbus_receive() function shall receive an indication request from the socket of the context ctx. This function
// is used by a Modbus slave/server to receive and analyze indication request sent by the masters/clients.
//
// If you need to use another socket or file descriptor than the one defined in the context ctx, see the function
// modbus_set_socket.
func (x *Modbus) TcpReceive() (req []byte, err error) {
	recv := make([]C.uint8_t, MODBUS_TCP_MAX_ADU_LENGTH)
	code := C.modbus_receive(x.ctx, unsafe.SliceData(recv))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	for i := range code {
		req = append(req, byte(recv[i]))
	}
	return
}

// RtuReceiveConfirmation modbus_receive_confirmation - receive a confirmation request
//
// The modbus_receive_confirmation() function shall receive a request via the socket of the context ctx. This function
// must be used for debugging purposes because the received response isn't checked against the initial request. This
// function can be used to receive request not handled by the library.
//
// The maximum size of the response depends on the used backend, in RTU the rsp array must be MODBUS_RTU_MAX_ADU_LENGTH
// bytes and in TCP it must be MODBUS_TCP_MAX_ADU_LENGTH bytes. If you want to write code compatible with both, you can
// use the constant MODBUS_MAX_ADU_LENGTH (maximum value of all libmodbus backends). Take care to allocate enough
// memory to store responses to avoid crashes of your server.
func (x *Modbus) TcpReceiveConfirmation() (rsp []byte, err error) {
	recv := make([]C.uint8_t, MODBUS_TCP_MAX_ADU_LENGTH)
	code := C.modbus_receive_confirmation(x.ctx, unsafe.SliceData(recv))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	for i := range code {
		rsp = append(rsp, byte(recv[i]))
	}
	return
}
