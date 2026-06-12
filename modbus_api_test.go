package modbus

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand/v2"
	"testing"
)

const (
	LOOP          = 1
	SERVER_ID     = 17
	ADDRESS_START = 0
	ADDRESS_END   = 99
)

func inittbl(nb int) (tab_rq_bits []byte, tab_rq_registers []uint16, tab_rw_rq_registers []uint16) {
	tab_rq_bits = make([]byte, nb)
	tab_rq_registers = make([]uint16, nb)
	tab_rw_rq_registers = make([]uint16, nb)
	for i := range nb {
		tab_rq_registers[i] = uint16(rand.UintN(65535))
		tab_rw_rq_registers[i] = ^tab_rq_registers[i]
		tab_rq_bits[i] = byte(tab_rq_registers[i] % 2)
	}
	return
}

func setup(outChan chan struct{}, errChan chan error, port int) {
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		errChan <- err
		return
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	mbMapping, err := NewMapping(500, 500, 500, 500)
	if err != nil {
		errChan <- err
		return
	}
	defer mbMapping.Free()
	_, err = ctx.TcpListen(1)
	if err != nil {
		errChan <- err
		return
	}
	outChan <- struct{}{}
	_, err = ctx.TcpAccept()
	if err != nil {
		errChan <- err
		return
	}

	for {
		req, err := ctx.Receive()
		if err != nil {
			log.Printf("receive error: %s", err)
			break
		}
		err = ctx.Reply(req, mbMapping)
		if err != nil {
			log.Printf("reply error: %s", err)
			break
		}
	}
}

func TestModbus_WriteBit(t *testing.T) {
	nb := ADDRESS_END - ADDRESS_START
	outChan := make(chan struct{})
	errChan := make(chan error, 1)
	const port = 15020
	go setup(outChan, errChan, port)
	select {
	case <-outChan:
	case err := <-errChan:
		t.Error("setup error:", err)
		t.FailNow()
	}
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Error("NewTCP error:", err)
		t.FailNow()
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	err = ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tab_rq_bits, _, _ := inittbl(nb)

	addr := ADDRESS_START

	err = ctx.WriteBit(addr, tab_rq_bits[0] != 0)
	if err != nil {
		t.Errorf("ERROR modbus_write_bit (%s)\n", err)
		t.Errorf("Address = %d, value = %d\n", addr, tab_rq_bits[0])
		t.FailNow()
	} else {
		out, err := ctx.ReadBits(addr, 1)
		if err != nil || out[0] != tab_rq_bits[0] {
			t.Errorf("ERROR modbus_read_bits single (%s)\n", err)
			t.Errorf("Address = %d\n", addr)
			t.FailNow()
		}
	}
}

func TestModbus_WriteBits(t *testing.T) {
	nb := ADDRESS_END - ADDRESS_START
	outChan := make(chan struct{})
	errChan := make(chan error, 1)
	const port = 15021
	go setup(outChan, errChan, port)
	select {
	case <-outChan:
	case err := <-errChan:
		t.Error("setup error:", err)
		t.FailNow()
	}
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Error("NewTCP error:", err)
		t.FailNow()
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	err = ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tab_rq_bits, _, _ := inittbl(nb)

	addr := ADDRESS_START

	err = ctx.WriteBits(addr, tab_rq_bits[:nb])
	if err != nil {
		t.Errorf("ERROR modbus_write_bits (%s)\n", err)
		t.Errorf("Address = %d, nb = %d\n", addr, nb)
		t.FailNow()
	} else {
		out, err := ctx.ReadBits(addr, nb)
		if err != nil {
			t.Errorf("ERROR modbus_read_bits\n")
			t.Errorf("Address = %d, nb = %d\n", addr, nb)
			t.FailNow()
		} else {
			for i := range nb {
				if out[i] != tab_rq_bits[i] {
					t.Errorf("ERROR modbus_read_bits\n")
					t.Errorf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
						addr,
						tab_rq_bits[i],
						tab_rq_bits[i],
						out[i],
						out[i])
					t.FailNow()
				}
			}
		}
	}
}

func TestModbus_WriteRegister(t *testing.T) {
	nb := ADDRESS_END - ADDRESS_START
	outChan := make(chan struct{})
	errChan := make(chan error, 1)
	const port = 15022
	go setup(outChan, errChan, port)
	select {
	case <-outChan:
	case err := <-errChan:
		t.Error("setup error:", err)
		t.FailNow()
	}
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Error("NewTCP error:", err)
		t.FailNow()
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	err = ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, tab_rq_registers, _ := inittbl(nb)

	addr := ADDRESS_START

	err = ctx.WriteRegister(addr, tab_rq_registers[0])
	if err != nil {
		t.Errorf("ERROR modbus_write_register (%s)\n", err)
		t.Errorf("Address = %d, value = %d\n", addr, tab_rq_registers[0])
		t.FailNow()
	} else {
		out, err := ctx.ReadRegisters(addr, 1)
		if err != nil || out[0] != tab_rq_registers[0] {
			t.Errorf("ERROR modbus_read_registers single (%s)\n", err)
			t.Errorf("Address = %d\n", addr)
			t.FailNow()
		}
	}
}

func TestModbus_WriteRegisters(t *testing.T) {
	nb := ADDRESS_END - ADDRESS_START
	outChan := make(chan struct{})
	errChan := make(chan error, 1)
	const port = 15023
	go setup(outChan, errChan, port)
	select {
	case <-outChan:
	case err := <-errChan:
		t.Error("setup error:", err)
		t.FailNow()
	}
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Error("NewTCP error:", err)
		t.FailNow()
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	err = ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, tab_rq_registers, _ := inittbl(nb)

	addr := ADDRESS_START

	err = ctx.WriteRegisters(addr, tab_rq_registers[:nb])
	if err != nil {
		t.Errorf("ERROR modbus_write_registers (%s)\n", err)
		t.Errorf("Address = %d, nb = %d\n", addr, nb)
		t.FailNow()
	} else {
		out, err := ctx.ReadRegisters(addr, nb)
		if err != nil {
			t.Errorf("ERROR modbus_read_registers\n")
			t.Errorf("Address = %d, nb = %d\n", addr, nb)
			t.FailNow()
		} else {
			for i := range nb {
				if out[i] != tab_rq_registers[i] {
					t.Errorf("ERROR modbus_read_registers\n")
					t.Errorf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
						addr,
						tab_rq_registers[i],
						tab_rq_registers[i],
						out[i],
						out[i])
					t.FailNow()
				}
			}
		}
	}
}

func TestModbus_WriteAndReadRegisters(t *testing.T) {
	nb := ADDRESS_END - ADDRESS_START
	outChan := make(chan struct{})
	errChan := make(chan error, 1)
	const port = 15024
	go setup(outChan, errChan, port)
	select {
	case <-outChan:
	case err := <-errChan:
		t.Error("setup error:", err)
		t.FailNow()
	}
	ctx, err := NewTCP("127.0.0.1", port)
	if err != nil {
		t.Error("NewTCP error:", err)
		t.FailNow()
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	err = ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, _, tab_rw_rq_registers := inittbl(nb)

	addr := ADDRESS_START

	out, err := ctx.WriteAndReadRegisters(addr, tab_rw_rq_registers[:nb], addr, nb)
	if err != nil {
		t.Errorf("ERROR modbus_read_and_write_registers (%s)\n", err)
		t.Errorf("Address = %d, nb = %d\n", addr, nb)
		t.FailNow()
	} else {
		for i := range nb {
			if out[i] != tab_rw_rq_registers[i] {
				t.Errorf("ERROR modbus_read_and_write_registers READ\n")
				t.Errorf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
					addr,
					out[i],
					tab_rw_rq_registers[i],
					out[i],
					tab_rw_rq_registers[i])
				t.FailNow()
			}
		}

		out, err := ctx.ReadRegisters(addr, nb)
		if err != nil {
			t.Errorf("ERROR modbus_read_registers (%s)\n", err)
			t.Errorf("Address = %d, nb = %d\n", addr, nb)
			t.FailNow()
		} else {
			for i := range nb {
				if tab_rw_rq_registers[i] != out[i] {
					t.Errorf("ERROR modbus_read_and_write_registers WRITE\n")
					t.Errorf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
						addr,
						tab_rw_rq_registers[i],
						tab_rw_rq_registers[i],
						out[i],
						out[i])
					t.FailNow()
				}
			}
		}
	}
}

func TestModbus_MarshalJSON(t *testing.T) {
	nb := 10
	mbMapping, err := NewMapping(nb, nb, nb, nb)
	if err != nil {
		t.Error("NewMapping error:", err)
		t.FailNow()
	}
	for i := range nb {
		if err := mbMapping.SetTabBits(i, byte(i%2)); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	for i := range nb {
		if err := mbMapping.SetTabInputBits(i, byte((i+1)%2)); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	for i := range nb {
		if err := mbMapping.SetTabInputRegisters(i, uint16(i+1)); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	for i := range nb {
		if err := mbMapping.SetTabRegisters(i, uint16(i)); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	err = enc.Encode(mbMapping)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(len(buffer.Bytes()))
	dec := json.NewDecoder(buffer)
	mbMapping = &ModbusMapping{}
	err = dec.Decode(mbMapping)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if mbMapping.NbBits() != nb {
		t.FailNow()
	}
	if mbMapping.NbInputBits() != nb {
		t.FailNow()
	}
	if mbMapping.NbRegisters() != nb {
		t.FailNow()
	}
	if mbMapping.NbInputRegisters() != nb {
		t.FailNow()
	}
	for k, v := range mbMapping.TabBits() {
		t.Log(k, v)
	}
	for k, v := range mbMapping.TabInputBits() {
		t.Log(k, v)
	}
	for k, v := range mbMapping.TabInputRegisters() {
		t.Log(k, v)
	}
	for k, v := range mbMapping.TabRegisters() {
		t.Log(k, v)
	}
}
