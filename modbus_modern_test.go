package modbus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"
)

// waitListening polls addr until a TCP connection succeeds or the deadline
// passes. It replaces time.Sleep-based synchronization in tests.
func waitListening(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server at %s did not start listening within 2s", addr)
}

func TestError_ErrorsIs(t *testing.T) {
	err := &Error{code: EMBXILFUN, message: "test"}

	if !errors.Is(err, EMBXILFUN) {
		t.Error("errors.Is(err, EMBXILFUN) = false, want true")
	}
	if errors.Is(err, EMBXILADD) {
		t.Error("errors.Is(err, EMBXILADD) = true, want false")
	}

	// The check must also work through wrapping.
	wrapped := fmt.Errorf("modbus operation failed: %w", err)
	if !errors.Is(wrapped, EMBXILFUN) {
		t.Error("errors.Is(wrapped, EMBXILFUN) = false, want true")
	}

	var ec ErrorCode
	if !errors.As(err, &ec) || ec != EMBXILFUN {
		t.Errorf("errors.As(err, &ec) got %v, want %v", ec, EMBXILFUN)
	}
}

func TestModbus_Free_Idempotent(t *testing.T) {
	ctx, err := NewTCP("127.0.0.1", 0)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}

	ctx.Free()
	ctx.Free() // must not panic or double-free

	if _, err := ctx.GetSlave(); err == nil {
		t.Error("GetSlave after Free succeeded, want error")
	}
	if err := ctx.Close(); err != nil {
		t.Errorf("Close after Free error: %v", err)
	}
}

func TestModbusMapping_Free_Idempotent(t *testing.T) {
	mm, err := NewMapping(8, 8, 8, 8)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	mm.Free()
	mm.Free() // must not panic or double-free
}

func TestTCPServer_ServeContext_Cancel(t *testing.T) {
	mapping, err := NewMapping(100, 100, 100, 100)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mapping.Free()

	srv := NewTCPServer("127.0.0.1:15060", mapping)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ServeContext(ctx)
	}()

	waitListening(t, "127.0.0.1:15060")
	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("ServeContext returned %v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ServeContext did not return after cancel")
	}
}

func TestModbus_ReceiveContext_Cancel(t *testing.T) {
	const port = 15061

	srv, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer srv.Free()
	if _, err := srv.TCPListen(1); err != nil {
		t.Fatalf("TCPListen error: %v", err)
	}

	// Connect a client that never sends anything, so Receive blocks. The Dial
	// itself confirms the listener is ready (no probe connection, which would
	// pollute the accept backlog).
	client, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer client.Close()

	if _, err := srv.TCPAccept(); err != nil {
		t.Fatalf("TCPAccept error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	reqCh := make(chan error, 1)
	go func() {
		_, err := srv.ReceiveContext(ctx)
		reqCh <- err
	}()

	// Let the receive block, then cancel it.
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-reqCh:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("ReceiveContext returned %v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ReceiveContext did not return after cancel")
	}
}

func TestModbus_ReceiveConfirmationContext_Cancel(t *testing.T) {
	const port = 15062

	srv, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer srv.Free()
	if _, err := srv.TCPListen(1); err != nil {
		t.Fatalf("TCPListen error: %v", err)
	}

	// A silent peer keeps ReceiveConfirmation blocked.
	client, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer client.Close()

	if _, err := srv.TCPAccept(); err != nil {
		t.Fatalf("TCPAccept error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	rspCh := make(chan error, 1)
	go func() {
		_, err := srv.ReceiveConfirmationContext(ctx)
		rspCh <- err
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-rspCh:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("ReceiveConfirmationContext returned %v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ReceiveConfirmationContext did not return after cancel")
	}
}
