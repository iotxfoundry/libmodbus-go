package modbus

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestGetHighLowByte(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		wantHigh byte
		wantLow  byte
	}{
		{"zero", 0x0000, 0x00, 0x00},
		{"byte swap", 0x1234, 0x12, 0x34},
		{"high bits", 0xABCD, 0xAB, 0xCD},
		{"max uint16", 0xFFFF, 0xFF, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHighByte(uint16(tt.value)); got != tt.wantHigh {
				t.Errorf("GetHighByte(%#04x) = %#02x, want %#02x", tt.value, got, tt.wantHigh)
			}
			if got := GetLowByte(uint16(tt.value)); got != tt.wantLow {
				t.Errorf("GetLowByte(%#04x) = %#02x, want %#02x", tt.value, got, tt.wantLow)
			}
			// Wider types must agree on the low 16 bits.
			if got := GetHighByte(uint64(tt.value)); got != tt.wantHigh {
				t.Errorf("GetHighByte[uint64](%#04x) = %#02x, want %#02x", tt.value, got, tt.wantHigh)
			}
			if got := GetLowByte(int64(tt.value)); got != tt.wantLow {
				t.Errorf("GetLowByte[int64](%#04x) = %#02x, want %#02x", tt.value, got, tt.wantLow)
			}
		})
	}
}

func TestIntConversions(t *testing.T) {
	t.Run("int16 round trip", func(t *testing.T) {
		for _, v := range []int16{0, 1, -1, 0x1234, math.MinInt16, math.MaxInt16} {
			tab := SetInt16ToInt8(v)
			if got := GetInt16FromInt8(tab); got != v {
				t.Errorf("int16 round trip: got %d, want %d", got, v)
			}
		}
	})
	t.Run("int32 round trip", func(t *testing.T) {
		for _, v := range []int32{0, 1, -1, 0x12345678, math.MinInt32, math.MaxInt32} {
			tab := SetInt32ToInt16(v)
			if got := GetInt32FromInt16(tab); got != v {
				t.Errorf("int32 round trip: got %d, want %d", got, v)
			}
		}
	})
	t.Run("int64 round trip", func(t *testing.T) {
		for _, v := range []int64{0, 1, -1, 0x1122334455667788, math.MinInt64, math.MaxInt64} {
			tab := SetInt64ToInt16(v)
			if got := GetInt64FromInt16(tab); got != v {
				t.Errorf("int64 round trip: got %d, want %d", got, v)
			}
		}
	})
	t.Run("big endian layout", func(t *testing.T) {
		if got := SetInt32ToInt16(0x00010002); !reflect.DeepEqual(got, []int16{1, 2}) {
			t.Errorf("SetInt32ToInt16(0x00010002) = %v, want [1 2]", got)
		}
		if got := SetInt64ToInt16(0x0001000200030004); !reflect.DeepEqual(got, []int16{1, 2, 3, 4}) {
			t.Errorf("SetInt64ToInt16 = %v, want [1 2 3 4]", got)
		}
	})
}

func TestBitsAndBytes(t *testing.T) {
	t.Run("SetBitsFromByte/GetByteFromBits round trip", func(t *testing.T) {
		for _, v := range []byte{0x00, 0xA5, 0xFF, 0x01, 0x80} {
			dest := make([]byte, 16)
			if err := SetBitsFromByte(dest, 2, v); err != nil {
				t.Fatalf("SetBitsFromByte(%#x) error: %v", v, err)
			}
			got, err := GetByteFromBits(dest, 2, 8)
			if err != nil {
				t.Fatalf("GetByteFromBits error: %v", err)
			}
			if got != v {
				t.Errorf("round trip %#x: got %#x", v, got)
			}
		}
	})
	t.Run("GetByteFromBits single bit", func(t *testing.T) {
		// libmodbus stores one bit per byte; reading nb=1 returns that bit.
		src := []byte{0, 0, 1, 0, 1}
		for i, want := range src {
			got, err := GetByteFromBits(src, i, 1)
			if err != nil {
				t.Fatalf("GetByteFromBits error: %v", err)
			}
			if got != want {
				t.Errorf("bit %d = %d, want %d", i, got, want)
			}
		}
	})
	t.Run("SetBitsFromBytes succeeds", func(t *testing.T) {
		dest := make([]byte, 16)
		tab := []byte{1, 0, 1, 0, 1, 0, 1, 0}
		if err := SetBitsFromBytes(dest, 0, uint(len(tab)), tab); err != nil {
			t.Errorf("SetBitsFromBytes error: %v", err)
		}
	})
	t.Run("boundary validation", func(t *testing.T) {
		dest := make([]byte, 4)
		if err := SetBitsFromByte(dest, 0, 0xFF); err == nil {
			t.Error("SetBitsFromByte with short dest: want error, got nil")
		}
		if err := SetBitsFromByte(dest, -1, 0xFF); err == nil {
			t.Error("SetBitsFromByte with negative index: want error, got nil")
		}
		if err := SetBitsFromBytes(dest, 0, 5, []byte{1}); err == nil {
			t.Error("SetBitsFromBytes with nb > len(tab): want error, got nil")
		}
		if _, err := GetByteFromBits(dest, 0, 9); err == nil {
			t.Error("GetByteFromBits with nb > 8: want error, got nil")
		}
		if _, err := GetByteFromBits(dest, 3, 4); err == nil {
			t.Error("GetByteFromBits past end: want error, got nil")
		}
	})
}

func TestFloatCodecs(t *testing.T) {
	codecs := []struct {
		name   string
		encode func(float32, []uint16) error
		decode func([]uint16) (float32, error)
	}{
		{"default", EncodeFloat, DecodeFloat},
		{"ABCD", EncodeFloatABCD, DecodeFloatABCD},
		{"DCBA", EncodeFloatDCBA, DecodeFloatDCBA},
		{"BADC", EncodeFloatBADC, DecodeFloatBADC},
		{"CDAB", EncodeFloatCDAB, DecodeFloatCDAB},
	}
	values := []float32{0, 1.5, -1.5, math.Pi, math.MaxFloat32, math.SmallestNonzeroFloat32}
	for _, c := range codecs {
		t.Run(c.name, func(t *testing.T) {
			for _, want := range values {
				reg := make([]uint16, 2)
				if err := c.encode(want, reg); err != nil {
					t.Fatalf("encode(%v) error: %v", want, err)
				}
				got, err := c.decode(reg)
				if err != nil {
					t.Fatalf("decode error: %v", err)
				}
				if got != want {
					t.Errorf("round trip %v: got %v", want, got)
				}
			}
			if err := c.encode(1.0, make([]uint16, 1)); err == nil {
				t.Error("encode with short dest: want error, got nil")
			}
			if _, err := c.decode(make([]uint16, 1)); err == nil {
				t.Error("decode with short src: want error, got nil")
			}
		})
	}
}

func TestMappingTabAccessors(t *testing.T) {
	mm, err := NewMapping(8, 8, 8, 8)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mm.Free()

	if n := mm.NbBits(); n != 8 {
		t.Errorf("NbBits = %d, want 8", n)
	}
	if err := mm.SetTabRegisters(0, 0xBEEF); err != nil {
		t.Fatalf("SetTabRegisters error: %v", err)
	}
	if err := mm.SetTabInputRegisters(1, 0xF00D); err != nil {
		t.Fatalf("SetTabInputRegisters error: %v", err)
	}
	if err := mm.SetTabBits(2, 1); err != nil {
		t.Fatalf("SetTabBits error: %v", err)
	}
	if err := mm.SetTabInputBits(3, 1); err != nil {
		t.Fatalf("SetTabInputBits error: %v", err)
	}

	if v, err := mm.GetTabRegisters(0); err != nil || v != 0xBEEF {
		t.Errorf("GetTabRegisters = (%#x, %v), want (0xBEEF, nil)", v, err)
	}
	if v, err := mm.GetTabInputRegisters(1); err != nil || v != 0xF00D {
		t.Errorf("GetTabInputRegisters = (%#x, %v), want (0xF00D, nil)", v, err)
	}
	if v, err := mm.GetTabBits(2); err != nil || v != 1 {
		t.Errorf("GetTabBits = (%v, %v), want (1, nil)", v, err)
	}
	if v, err := mm.GetTabInputBits(3); err != nil || v != 1 {
		t.Errorf("GetTabInputBits = (%v, %v), want (1, nil)", v, err)
	}

	// Out-of-range access must fail, not panic.
	if _, err := mm.GetTabRegisters(99); err == nil {
		t.Error("GetTabRegisters out of range: want error, got nil")
	}
	if _, err := mm.GetTabInputRegisters(99); err == nil {
		t.Error("GetTabInputRegisters out of range: want error, got nil")
	}
	if err := mm.SetTabRegisters(99, 1); err == nil {
		t.Error("SetTabRegisters out of range: want error, got nil")
	}
}

// rawReadReq builds a Modbus PDU (unit id + function code + data) suitable for
// SendRawRequest/SendRawRequestTID, which add the backend header themselves.
func rawReadReq(unit byte, addr, nb uint16) []byte {
	req := make([]byte, 6)
	req[0] = unit
	req[1] = 0x03 // read holding registers
	binary.BigEndian.PutUint16(req[2:], addr)
	binary.BigEndian.PutUint16(req[4:], nb)
	return req
}

// tcpReadReq builds a full Modbus TCP frame (MBAP header + PDU) as a peer would
// write to a socket.
func tcpReadReq(tid uint16, unit byte, addr, nb uint16) []byte {
	req := make([]byte, 12)
	binary.BigEndian.PutUint16(req[0:], tid) // transaction id
	binary.BigEndian.PutUint16(req[2:], 0)   // protocol id
	binary.BigEndian.PutUint16(req[4:], 6)   // length: unit + fc + 4 bytes
	req[6] = unit
	req[7] = 0x03
	binary.BigEndian.PutUint16(req[8:], addr)
	binary.BigEndian.PutUint16(req[10:], nb)
	return req
}

func TestQuirksAndErrorRecovery(t *testing.T) {
	const port = 15070
	srvCtx, srvReady, srvErr := startReplyServer(t, port)
	defer srvCtx.Free()
	<-srvReady
	if err := peekErr(srvErr); err != nil {
		t.Fatalf("server setup error: %v", err)
	}

	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer ctx.Free()
	if err := ctx.Connect(); err != nil {
		t.Fatalf("Connect error: %v", err)
	}

	if err := ctx.EnableQuirks(MODBUS_QUIRK_ALL); err != nil {
		t.Errorf("EnableQuirks error: %v", err)
	}
	if err := ctx.DisableQuirks(MODBUS_QUIRK_ALL); err != nil {
		t.Errorf("DisableQuirks error: %v", err)
	}
	if err := ctx.SetErrorRecovery(MODBUS_ERROR_RECOVERY_LINK | MODBUS_ERROR_RECOVERY_PROTOCOL); err != nil {
		t.Errorf("SetErrorRecovery error: %v", err)
	}
}

func TestSendRawRequestAndConfirmation(t *testing.T) {
	const port = 15071
	srvCtx, srvReady, srvErr := startReplyServer(t, port)
	defer srvCtx.Free()
	<-srvReady
	if err := peekErr(srvErr); err != nil {
		t.Fatalf("server setup error: %v", err)
	}

	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer ctx.Free()
	if err := ctx.Connect(); err != nil {
		t.Fatalf("Connect error: %v", err)
	}

	// SendRawRequest: server replies through its mapping; read the confirmation back.
	if err := ctx.SendRawRequest(rawReadReq(SERVER_ID, 0, 2)); err != nil {
		t.Fatalf("SendRawRequest error: %v", err)
	}
	rsp, err := ctx.ReceiveConfirmation()
	if err != nil {
		t.Fatalf("ReceiveConfirmation error: %v", err)
	}
	if len(rsp) < 9 || rsp[7] != 0x03 {
		t.Errorf("unexpected confirmation: % x", rsp)
	}

	// SendRawRequestTID: the transaction id must be echoed in the response.
	const tid = 0x2A
	if err := ctx.SendRawRequestTID(rawReadReq(SERVER_ID, 0, 1), tid); err != nil {
		t.Fatalf("SendRawRequestTID error: %v", err)
	}
	rsp, err = ctx.ReceiveConfirmation()
	if err != nil {
		t.Fatalf("ReceiveConfirmation error: %v", err)
	}
	if got := binary.BigEndian.Uint16(rsp[0:2]); got != tid {
		t.Errorf("transaction id = %#x, want %#x", got, tid)
	}
}

func TestReplyException(t *testing.T) {
	const port = 15072

	// Server side: accept one client, receive a request, answer with an exception.
	srv, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer srv.Free()
	if _, err := srv.TCPListen(1); err != nil {
		t.Fatalf("TCPListen error: %v", err)
	}

	client, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP client error: %v", err)
	}
	defer client.Free()

	accepted := make(chan error, 1)
	go func() {
		if _, err := srv.TCPAccept(); err != nil {
			accepted <- err
			return
		}
		req, err := srv.Receive()
		if err != nil {
			accepted <- err
			return
		}
		accepted <- srv.ReplyException(req, uint(MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS))
	}()

	if err := client.Connect(); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	if err := client.SendRawRequest(rawReadReq(SERVER_ID, 0, 1)); err != nil {
		t.Fatalf("SendRawRequest error: %v", err)
	}
	rsp, err := client.ReceiveConfirmation()
	if err != nil {
		t.Fatalf("ReceiveConfirmation error: %v", err)
	}
	// Exception response: function code 0x83 + exception code.
	if len(rsp) < 9 || rsp[7] != 0x83 || rsp[8] != byte(MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS) {
		t.Errorf("unexpected exception response: % x", rsp)
	}
	if err := <-accepted; err != nil {
		t.Errorf("server side error: %v", err)
	}
}

func TestProxy(t *testing.T) {
	const backendPort = 15073
	const frontendPort = 15074

	// Backend: a normal reply server holding the real mapping.
	backendSrv, backendReady, backendErr := startReplyServer(t, backendPort)
	defer backendSrv.Free()
	<-backendReady
	if err := peekErr(backendErr); err != nil {
		t.Fatalf("backend setup error: %v", err)
	}

	// Backend client context: what the proxy forwards to.
	backend, err := NewTCP("127.0.0.1", backendPort)
	if err != nil {
		t.Fatalf("NewTCP backend error: %v", err)
	}
	defer backend.Free()
	if err := backend.Connect(); err != nil {
		t.Fatalf("backend Connect error: %v", err)
	}

	// Frontend: listener + accepted connection, driven by a dummy network peer.
	frontend, err := NewTCP("127.0.0.1", frontendPort)
	if err != nil {
		t.Fatalf("NewTCP frontend error: %v", err)
	}
	defer frontend.Free()
	if _, err := frontend.TCPListen(1); err != nil {
		t.Fatalf("TCPListen error: %v", err)
	}
	peer, err := net.Dial("tcp", "127.0.0.1:15074")
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer peer.Close()
	if _, err := frontend.TCPAccept(); err != nil {
		t.Fatalf("TCPAccept error: %v", err)
	}

	// Forward a read request: backend answers, response lands on the peer socket.
	req := tcpReadReq(9, SERVER_ID, 0, 1)
	if err := frontend.Proxy(backend, req); err != nil {
		t.Fatalf("Proxy error: %v", err)
	}

	buf := make([]byte, 64)
	peer.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := peer.Read(buf)
	if err != nil {
		t.Fatalf("peer read error: %v", err)
	}
	rsp := buf[:n]
	if len(rsp) < 9 || rsp[7] != 0x03 {
		t.Errorf("proxied response: % x", rsp)
	}
}

// startReplyServer starts a libmodbus reply-loop server on port and returns
// its context plus a channel closed when the server is listening.
func startReplyServer(t *testing.T, port int) (*Modbus, chan struct{}, chan error) {
	t.Helper()
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	mapping, err := NewMapping(100, 100, 100, 100)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	t.Cleanup(mapping.Free)

	ready := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		if _, err := ctx.TCPListen(1); err != nil {
			errCh <- err
			return
		}
		close(ready)
		if _, err := ctx.TCPAccept(); err != nil {
			errCh <- err
			return
		}
		for {
			req, err := ctx.Receive()
			if err != nil {
				return
			}
			if err := ctx.Reply(req, mapping); err != nil {
				return
			}
		}
	}()
	return ctx, ready, errCh
}

func TestClientReadWrite(t *testing.T) {
	const port = 15075
	srvCtx, srvReady, srvErr := startReplyServer(t, port)
	defer srvCtx.Free()
	<-srvReady
	if err := peekErr(srvErr); err != nil {
		t.Fatalf("server setup error: %v", err)
	}

	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer ctx.Free()
	if err := ctx.Connect(); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	if err := ctx.SetSlave(SERVER_ID); err != nil {
		t.Fatalf("SetSlave error: %v", err)
	}

	t.Run("registers round trip", func(t *testing.T) {
		want := []uint16{0x1111, 0x2222, 0x3333}
		if err := ctx.WriteRegisters(0, want); err != nil {
			t.Fatalf("WriteRegisters error: %v", err)
		}
		got, err := ctx.ReadRegisters(0, len(want))
		if err != nil {
			t.Fatalf("ReadRegisters error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("registers = %v, want %v", got, want)
		}
	})
	t.Run("single register", func(t *testing.T) {
		if err := ctx.WriteRegister(10, 0xABCD); err != nil {
			t.Fatalf("WriteRegister error: %v", err)
		}
		got, err := ctx.ReadRegisters(10, 1)
		if err != nil {
			t.Fatalf("ReadRegisters error: %v", err)
		}
		if len(got) != 1 || got[0] != 0xABCD {
			t.Errorf("register = %v, want [0xABCD]", got)
		}
	})
	t.Run("mask write register", func(t *testing.T) {
		if err := ctx.WriteRegister(20, 0x00FF); err != nil {
			t.Fatalf("WriteRegister error: %v", err)
		}
		// (value AND andMask) OR (orMask AND NOT andMask)
		if err := ctx.MaskWriteRegister(20, 0x0F0F, 0x0011); err != nil {
			t.Fatalf("MaskWriteRegister error: %v", err)
		}
		got, err := ctx.ReadRegisters(20, 1)
		if err != nil {
			t.Fatalf("ReadRegisters error: %v", err)
		}
		want := (0x00FF & 0x0F0F) | (0x0011 &^ 0x0F0F)
		if len(got) != 1 || got[0] != uint16(want) {
			t.Errorf("masked register = %#x, want %#x", got[0], want)
		}
	})
	t.Run("write and read registers", func(t *testing.T) {
		src := []uint16{0xAAAA, 0xBBBB}
		dest, err := ctx.WriteAndReadRegisters(30, src, 0, 3)
		if err != nil {
			t.Fatalf("WriteAndReadRegisters error: %v", err)
		}
		if len(dest) != 3 {
			t.Errorf("write-and-read returned %d registers, want 3", len(dest))
		}
		got, err := ctx.ReadRegisters(30, 2)
		if err != nil {
			t.Fatalf("ReadRegisters error: %v", err)
		}
		if !reflect.DeepEqual(got, src) {
			t.Errorf("written registers = %v, want %v", got, src)
		}
	})
	t.Run("bits round trip", func(t *testing.T) {
		want := []byte{1, 0, 1, 1, 0, 0, 1, 0, 1, 0}
		if err := ctx.WriteBits(0, want); err != nil {
			t.Fatalf("WriteBits error: %v", err)
		}
		got, err := ctx.ReadBits(0, len(want))
		if err != nil {
			t.Fatalf("ReadBits error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("bits = %v, want %v", got, want)
		}
	})
	t.Run("single bit", func(t *testing.T) {
		if err := ctx.WriteBit(50, true); err != nil {
			t.Fatalf("WriteBit error: %v", err)
		}
		got, err := ctx.ReadBits(50, 1)
		if err != nil {
			t.Fatalf("ReadBits error: %v", err)
		}
		if len(got) != 1 || got[0] != 1 {
			t.Errorf("bit = %v, want [1]", got)
		}
	})
	t.Run("read only operations", func(t *testing.T) {
		if _, err := ctx.ReadInputBits(0, 8); err != nil {
			t.Errorf("ReadInputBits error: %v", err)
		}
		if _, err := ctx.ReadInputRegisters(0, 4); err != nil {
			t.Errorf("ReadInputRegisters error: %v", err)
		}
		if _, err := ctx.ReportSlaveID(); err != nil {
			t.Errorf("ReportSlaveID error: %v", err)
		}
	})
	t.Run("socket and misc", func(t *testing.T) {
		if _, err := ctx.GetSocket(); err != nil {
			t.Errorf("GetSocket error: %v", err)
		}
		if n := ctx.GetHeaderLength(); n <= 0 {
			t.Errorf("GetHeaderLength = %d, want > 0", n)
		}
		if err := ctx.Flush(); err != nil {
			t.Errorf("Flush error: %v", err)
		}
	})
}

func TestTimeouts(t *testing.T) {
	ctx, err := NewTCP("127.0.0.1", 0)
	if err != nil {
		t.Fatalf("NewTCP error: %v", err)
	}
	defer ctx.Free()

	cases := []struct {
		name string
		set  func(time.Duration) error
		get  func() (time.Duration, error)
	}{
		{"response", ctx.SetResponseTimeout, ctx.GetResponseTimeout},
		{"byte", ctx.SetByteTimeout, ctx.GetByteTimeout},
		{"indication", ctx.SetIndicationTimeout, ctx.GetIndicationTimeout},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want := 250 * time.Millisecond
			if err := c.set(want); err != nil {
				t.Fatalf("set error: %v", err)
			}
			got, err := c.get()
			if err != nil {
				t.Fatalf("get error: %v", err)
			}
			if got != want {
				t.Errorf("timeout = %v, want %v", got, want)
			}
		})
	}
	if err := ctx.SetResponseTimeout(-1); err == nil {
		t.Error("negative timeout: want error, got nil")
	}
}

func TestMappingAccessorsAndJSON(t *testing.T) {
	mm, err := NewMappingWithStart(1, 8, 2, 8, 3, 8, 4, 8)
	if err != nil {
		t.Fatalf("NewMappingWithStart error: %v", err)
	}
	defer mm.Free()

	if got := mm.StartBits(); got != 1 {
		t.Errorf("StartBits = %d, want 1", got)
	}
	if got := mm.NbBits(); got != 8 {
		t.Errorf("NbBits = %d, want 8", got)
	}
	if got := mm.NbInputBits(); got != 8 {
		t.Errorf("NbInputBits = %d, want 8", got)
	}
	if got := mm.NbRegisters(); got != 8 {
		t.Errorf("NbRegisters = %d, want 8", got)
	}
	if got := mm.NbInputRegisters(); got != 8 {
		t.Errorf("NbInputRegisters = %d, want 8", got)
	}
	if got := mm.StartRegisters(); got != 3 {
		t.Errorf("StartRegisters = %d, want 3", got)
	}

	// Populate then iterate with the range-over-func iterators.
	if err := mm.SetTabRegisters(3, 0x1234); err != nil {
		t.Fatalf("SetTabRegisters error: %v", err)
	}
	seen := map[int]uint16{}
	for i, v := range mm.TabRegisters() {
		seen[i] = v
	}
	if len(seen) != 8 {
		t.Errorf("TabRegisters yielded %d entries, want 8", len(seen))
	}
	if seen[3] != 0x1234 {
		t.Errorf("TabRegisters[3] = %#x, want 0x1234", seen[3])
	}
	for range mm.TabBits() {
	}
	for range mm.TabInputBits() {
	}
	for range mm.TabInputRegisters() {
	}

	// JSON round trip.
	data, err := json.Marshal(mm)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	mm2, err := NewMapping(1, 1, 1, 1)
	if err != nil {
		t.Fatalf("NewMapping error: %v", err)
	}
	defer mm2.Free()
	if err := json.Unmarshal(data, mm2); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if got := mm2.NbRegisters(); got != 8 {
		t.Errorf("unmarshaled NbRegisters = %d, want 8", got)
	}
}

func peekErr(ch chan error) error {
	select {
	case err := <-ch:
		return err
	default:
		return nil
	}
}
