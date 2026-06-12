package modbus

import (
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
