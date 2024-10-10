package main

import (
	"log"
	"math/rand/v2"

	libmodbusgo "github.com/iotxfoundry/libmodbus-go"
)

const (
	LOOP          = 1
	SERVER_ID     = 17
	ADDRESS_START = 0
	ADDRESS_END   = 99
)

func main() {
	ctx := libmodbusgo.ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		log.Println("ModbusNewTcp error")
		return
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	err := ctx.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}

	nb := ADDRESS_END - ADDRESS_START

	tab_rq_bits := make([]uint8, nb)
	tab_rq_registers := make([]uint16, nb)
	tab_rw_rq_registers := make([]uint16, nb)

	nbFail := 0
	for range LOOP {
		for addr := ADDRESS_START; addr < ADDRESS_END; addr++ {

			// Random numbers (short)
			for i := 0; i < nb; i++ {
				tab_rq_registers[i] = uint16(rand.UintN(65535))
				tab_rw_rq_registers[i] = ^tab_rq_registers[i]
				tab_rq_bits[i] = uint8(tab_rq_registers[i] % 2)
			}

			nb = ADDRESS_END - addr

			// WRITE BIT
			err = ctx.WriteBit(addr, tab_rq_bits[0])
			if err != nil {
				log.Printf("ERROR modbus_write_bit (%s)\n", err)
				log.Printf("Address = %d, value = %d\n", addr, tab_rq_bits[0])
				nbFail++
				return
			} else {
				out, err := ctx.ReadBits(addr, 1)
				if err != nil || out[0] != tab_rq_bits[0] {
					log.Printf("ERROR modbus_read_bits single (%s)\n", err)
					log.Printf("Address = %d\n", addr)
					nbFail++
					return
				}
			}

			// MULTIPLE BITS
			err = ctx.WriteBits(addr, tab_rq_bits[:nb])
			if err != nil {
				log.Printf("ERROR modbus_write_bits (%s)\n", err)
				log.Printf("Address = %d, nb = %d\n", addr, nb)
				nbFail++
			} else {
				out, err := ctx.ReadBits(addr, nb)
				if err != nil {
					log.Printf("ERROR modbus_read_bits\n")
					log.Printf("Address = %d, nb = %d\n", addr, nb)
					nbFail++
				} else {
					for i := 0; i < nb; i++ {
						if out[i] != tab_rq_bits[i] {
							log.Printf("ERROR modbus_read_bits\n")
							log.Printf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
								addr,
								tab_rq_bits[i],
								tab_rq_bits[i],
								out[i],
								out[i])
							nbFail++
						}
					}
				}
			}

			// SINGLE REGISTER
			err = ctx.WriteRegister(addr, tab_rq_registers[0])
			if err != nil {
				log.Printf("ERROR modbus_write_register (%s)\n", err)
				log.Printf("Address = %d, value = %d (0x%X)\n",
					addr,
					tab_rq_registers[0],
					tab_rq_registers[0])
				nbFail++
			} else {
				out, err := ctx.ReadRegisters(addr, 1)
				if err != nil {
					log.Printf("ERROR modbus_read_registers single (%s)\n", err)
					log.Printf("Address = %d\n", addr)
					nbFail++
				} else {
					if tab_rq_registers[0] != out[0] {
						log.Printf("ERROR modbus_read_registers single\n")
						log.Printf("Address = %d, value = %d (0x%X) != %d (0x%X)\n",
							addr,
							tab_rq_registers[0],
							tab_rq_registers[0],
							out[0],
							out[0])
						nbFail++
					}
				}
			}

			// MULTIPLE REGISTERS
			err = ctx.WriteRegisters(addr, tab_rq_registers[:nb])
			if err != nil {
				log.Printf("ERROR modbus_write_registers (%s)\n", err)
				log.Printf("Address = %d, nb = %d\n", addr, nb)
				nbFail++
			} else {
				out, err := ctx.ReadRegisters(addr, nb)
				if err != nil {
					log.Printf("ERROR modbus_read_registers (%s)\n", err)
					log.Printf("Address = %d, nb = %d\n", addr, nb)
					nbFail++
				} else {
					for i := 0; i < nb; i++ {
						if tab_rq_registers[i] != out[i] {
							log.Printf("ERROR modbus_read_registers\n")
							log.Printf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
								addr,
								tab_rq_registers[i],
								tab_rq_registers[i],
								out[i],
								out[i])
							nbFail++
						}
					}
				}
			}
			/* R/W MULTIPLE REGISTERS */
			out, err := ctx.WriteAndReadRegisters(addr, tab_rw_rq_registers[:nb], addr, nb)
			if err != nil {
				log.Printf("ERROR modbus_read_and_write_registers (%s)\n", err)
				log.Printf("Address = %d, nb = %d\n", addr, nb)
				nbFail++
			} else {
				for i := 0; i < nb; i++ {
					if out[i] != tab_rw_rq_registers[i] {
						log.Printf("ERROR modbus_read_and_write_registers READ\n")
						log.Printf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
							addr,
							out[i],
							tab_rw_rq_registers[i],
							out[i],
							tab_rw_rq_registers[i])
						nbFail++
					}
				}

				out, err := ctx.ReadRegisters(addr, nb)
				if err != nil {
					log.Printf("ERROR modbus_read_registers (%s)\n", err)
					log.Printf("Address = %d, nb = %d\n", addr, nb)
					nbFail++
				} else {
					for i := 0; i < nb; i++ {
						if tab_rw_rq_registers[i] != out[i] {
							log.Printf("ERROR modbus_read_and_write_registers WRITE\n")
							log.Printf("Address = %d, value %d (0x%X) != %d (0x%X)\n",
								addr,
								tab_rw_rq_registers[i],
								tab_rw_rq_registers[i],
								out[i],
								out[i])
							nbFail++
						}
					}
				}
			}
		}

		log.Printf("Test: ")
		if nbFail > 0 {
			log.Printf("%d FAILS\n", nbFail)
		} else {
			log.Printf("SUCCESS\n")
		}
	}
}
