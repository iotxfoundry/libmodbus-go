package libmodbusgo

import (
	"bytes"
	"encoding/gob"
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

var (
	nb                  = ADDRESS_END - ADDRESS_START
	tab_rq_bits         = make([]uint8, nb)
	tab_rq_registers    = make([]uint16, nb)
	tab_rw_rq_registers = make([]uint16, nb)
)

func inittbl() {
	for i := range nb {
		tab_rq_registers[i] = uint16(rand.UintN(65535))
		tab_rw_rq_registers[i] = ^tab_rq_registers[i]
		tab_rq_bits[i] = uint8(tab_rq_registers[i] % 2)
	}
}

func setup(outChan chan struct{}) {
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		log.Println("ModbusNewTcp error")
		return
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	mbMapping := ModbusMappingNew(500, 500, 500, 500)
	if mbMapping == nil {
		log.Println("ModbusMappingNew error")
		return
	}
	defer mbMapping.Free()
	_, err := ctx.TcpListen(1)
	if err != nil {
		log.Fatalln(err)
		return
	}
	outChan <- struct{}{}
	_, err = ctx.TcpAccept()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		req, err := ctx.TcpReceive()
		if err != nil {
			log.Printf("receive error: %s", err)
			break
		}
		err = ctx.Reply(req, mbMapping)
		if err != nil {
			log.Printf("reply error: %s", err)
			break
		}
		// for k, v := range mbMapping.TabBits() {
		// 	log.Println(k, v)
		// }
	}
}

func TestModbus_WriteBit(t *testing.T) {
	outChan := make(chan struct{})
	go setup(outChan)
	<-outChan
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		t.Error("ModbusNewTcp error")
		t.FailNow()
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inittbl()

	addr := ADDRESS_START

	err = ctx.WriteBit(addr, tab_rq_bits[0])
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
	outChan := make(chan struct{})
	go setup(outChan)
	<-outChan
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		t.Error("ModbusNewTcp error")
		t.FailNow()
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inittbl()

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
			for i := 0; i < nb; i++ {
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
	outChan := make(chan struct{})
	go setup(outChan)
	<-outChan
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		t.Error("ModbusNewTcp error")
		t.FailNow()
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inittbl()

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
	outChan := make(chan struct{})
	go setup(outChan)
	<-outChan
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		t.Error("ModbusNewTcp error")
		t.FailNow()
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inittbl()

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
			for i := 0; i < nb; i++ {
				if out[i] != tab_rq_registers[i] {
					t.Errorf("ERROR modbus_read_registers\n")
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

func TestModbus_WriteAndReadRegisters(t *testing.T) {
	outChan := make(chan struct{})
	go setup(outChan)
	<-outChan
	ctx := ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		t.Error("ModbusNewTcp error")
		t.FailNow()
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inittbl()

	addr := ADDRESS_START

	out, err := ctx.WriteAndReadRegisters(addr, tab_rw_rq_registers[:nb], addr, nb)
	if err != nil {
		t.Errorf("ERROR modbus_read_and_write_registers (%s)\n", err)
		t.Errorf("Address = %d, nb = %d\n", addr, nb)
		t.FailNow()
	} else {
		for i := 0; i < nb; i++ {
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
			for i := 0; i < nb; i++ {
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

func TestModbus_MarshalBinary(t *testing.T) {
	nb := 10
	mbMapping := ModbusMappingNew(nb, nb, nb, nb)
	for i := range nb {
		mbMapping.SetTabBits(i, byte(i%2))
	}
	for i := range nb {
		mbMapping.SetTabInputBits(i, byte((i+1)%2))
	}
	for i := range nb {
		mbMapping.SetTabInputRegisters(i, uint16(i+1))
	}
	for i := range nb {
		mbMapping.SetTabRegisters(i, uint16(i))
	}
	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(mbMapping)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(len(buffer.Bytes()))
	dec := gob.NewDecoder(buffer)
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
