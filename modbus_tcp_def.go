package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"
*/
import "C"

const (
	MODBUS_TCP_DEFAULT_PORT = C.MODBUS_TCP_DEFAULT_PORT
	MODBUS_TCP_SLAVE        = C.MODBUS_TCP_SLAVE
)

const (
	// MODBUS_TCP_MAX_ADU_LENGTH Modbus_Application_Protocol_V1_1b.pdf Chapter 4 Section 1 Page 5
	// TCP MODBUS ADU = 253 bytes + MBAP (7 bytes) = 260 bytes
	MODBUS_TCP_MAX_ADU_LENGTH = C.MODBUS_TCP_MAX_ADU_LENGTH
)
