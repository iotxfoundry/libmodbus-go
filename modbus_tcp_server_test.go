package modbus

import (
	"net"
	"testing"
	"time"
)

func TestTCPServer_StartStop(t *testing.T) {
	mapping, err := NewMapping(100, 100, 100, 100)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	srv := NewTCPServer("127.0.0.1:1503", mapping)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve()
	}()

	time.Sleep(100 * time.Millisecond)

	if err := srv.Close(); err != nil {
		t.Errorf("Close error: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Serve returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Serve did not return after Close")
	}
}

func TestTCPServer_SingleClient(t *testing.T) {
	mapping, err := NewMapping(500, 500, 500, 500)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	if err := mapping.SetTabRegisters(0, 42); err != nil {
		t.Fatalf("SetTabRegisters error: %v", err)
	}

	srv := NewTCPServer("127.0.0.1:1504", mapping)
	go func() {
		if err := srv.Serve(); err != nil {
			t.Logf("Serve error: %v", err)
		}
	}()
	defer srv.Close()

	time.Sleep(100 * time.Millisecond)

	ctx, err := NewTCP("127.0.0.1", 1504)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer ctx.Free()
	defer ctx.Close()

	if err := ctx.Connect(); err != nil {
		t.Fatalf("Connect error: %v", err)
	}

	regs, err := ctx.ReadRegisters(0, 1)
	if err != nil {
		t.Fatalf("ReadRegisters error: %v", err)
	}

	if len(regs) != 1 {
		t.Fatalf("ReadRegisters returned %d values, want 1", len(regs))
	}
	if regs[0] != 42 {
		t.Errorf("ReadRegisters got %d, want 42", regs[0])
	}
}

func TestTCPServer_MultipleClients(t *testing.T) {
	mapping, err := NewMapping(500, 500, 500, 500)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	if err := mapping.SetTabRegisters(0, 1234); err != nil {
		t.Fatalf("SetTabRegisters error: %v", err)
	}

	srv := NewTCPServer("127.0.0.1:1505", mapping)
	srv.SetMaxClients(3)
	go func() {
		if err := srv.Serve(); err != nil {
			t.Logf("Serve error: %v", err)
		}
	}()
	defer srv.Close()

	time.Sleep(100 * time.Millisecond)

	type clientResult struct {
		val uint16
		err error
	}

	const numClients = 3
	results := make([]clientResult, numClients)
	done := make(chan struct{})

	for i := range numClients {
		go func(idx int) {
			defer func() { done <- struct{}{} }()

			ctx, err := NewTCP("127.0.0.1", 1505)
			if err != nil {
				results[idx] = clientResult{err: err}
				return
			}
			defer ctx.Free()
			defer ctx.Close()

			if err := ctx.Connect(); err != nil {
				results[idx] = clientResult{err: err}
				return
			}

			regs, err := ctx.ReadRegisters(0, 1)
			if err != nil {
				results[idx] = clientResult{err: err}
				return
			}
			if len(regs) != 1 {
				results[idx] = clientResult{err: err}
				return
			}
			results[idx] = clientResult{val: regs[0]}
		}(i)
	}

	for range numClients {
		<-done
	}

	for i, r := range results {
		if r.err != nil {
			t.Errorf("client %d error: %v", i, r.err)
			continue
		}
		if r.val != 1234 {
			t.Errorf("client %d got %d, want 1234", i, r.val)
		}
	}
}

func TestTCPServer_GracefulClose(t *testing.T) {
	mapping, err := NewMapping(100, 100, 100, 100)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	srv := NewTCPServer("127.0.0.1:1506", mapping)

	go func() {
		if err := srv.Serve(); err != nil {
			t.Logf("Serve error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	if err := srv.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

// TestTCPServer_CloseUnblocksIdleConnection 回归: 挂着不发包的客户端连接会让
// handleConn 阻塞在 C 层 recv; close(2) 不会唤醒阻塞读, 若清理只 close 不
// shutdown, Close 的 wg.Wait 只能等对端下一拍请求才放行(表现为关停耗时数秒).
// 修复后 shutdown(2) 立即唤醒 recv, Close 应亚秒级返回.
func TestTCPServer_CloseUnblocksIdleConnection(t *testing.T) {
	mapping, err := NewMapping(100, 100, 100, 100)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	// 选一个空闲端口, 避免与其它测试冲突.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate port: %v", err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	srv := NewTCPServer(addr, mapping)
	go func() {
		if err := srv.Serve(); err != nil {
			t.Logf("Serve error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// 挂一条不发包的空闲连接(模拟外部主站轮询间隙).
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("dial idle client: %v", err)
	}
	defer conn.Close()
	time.Sleep(100 * time.Millisecond) // 等 handleConn 进入阻塞 recv

	done := make(chan struct{})
	go func() {
		_ = srv.Close()
		close(done)
	}()

	select {
	case <-done:
		// Close 及时返回
	case <-time.After(2 * time.Second):
		t.Error("Close 阻塞超过 2s: 空闲连接的阻塞 recv 未被 shutdown 唤醒")
	}
}
