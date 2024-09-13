package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"
*/
import "C"

const (
	// MODBUS_RTU_MAX_ADU_LENGTH Modbus_Application_Protocol_V1_1b.pdf Chapter 4 Section 1 Page 5
	// RS232 / RS485 ADU = 253 bytes + slave (1 byte) + CRC (2 bytes) = 256 bytes
	MODBUS_RTU_MAX_ADU_LENGTH = C.MODBUS_RTU_MAX_ADU_LENGTH
)

const (
	MODBUS_RTU_RS232 = C.MODBUS_RTU_RS232
	MODBUS_RTU_RS485 = C.MODBUS_RTU_RS485
)

const (
	MODBUS_RTU_RTS_NONE = C.MODBUS_RTU_RTS_NONE
	MODBUS_RTU_RTS_UP   = C.MODBUS_RTU_RTS_UP
	MODBUS_RTU_RTS_DOWN = C.MODBUS_RTU_RTS_DOWN
)

type SetRtsCallback func(ctx *Modbus, on int)
