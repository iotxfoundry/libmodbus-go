package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib -lmodbus -Wl,-rpath=/usr/local/lib
#include "modbus.h"

*/
import "C"
import (
	"fmt"
)

// Modbus function codes
const (
	MODBUS_FC_READ_COILS               = C.MODBUS_FC_READ_COILS
	MODBUS_FC_READ_DISCRETE_INPUTS     = C.MODBUS_FC_READ_DISCRETE_INPUTS
	MODBUS_FC_READ_HOLDING_REGISTERS   = C.MODBUS_FC_READ_HOLDING_REGISTERS
	MODBUS_FC_READ_INPUT_REGISTERS     = C.MODBUS_FC_READ_INPUT_REGISTERS
	MODBUS_FC_WRITE_SINGLE_COIL        = C.MODBUS_FC_WRITE_SINGLE_COIL
	MODBUS_FC_WRITE_SINGLE_REGISTER    = C.MODBUS_FC_WRITE_SINGLE_REGISTER
	MODBUS_FC_READ_EXCEPTION_STATUS    = C.MODBUS_FC_READ_EXCEPTION_STATUS
	MODBUS_FC_WRITE_MULTIPLE_COILS     = C.MODBUS_FC_WRITE_MULTIPLE_COILS
	MODBUS_FC_WRITE_MULTIPLE_REGISTERS = C.MODBUS_FC_WRITE_MULTIPLE_REGISTERS
	MODBUS_FC_REPORT_SLAVE_ID          = C.MODBUS_FC_REPORT_SLAVE_ID
	MODBUS_FC_MASK_WRITE_REGISTER      = C.MODBUS_FC_MASK_WRITE_REGISTER
	MODBUS_FC_WRITE_AND_READ_REGISTERS = C.MODBUS_FC_WRITE_AND_READ_REGISTERS
)

const (
	MODBUS_BROADCAST_ADDRESS = C.MODBUS_BROADCAST_ADDRESS
)

// Modbus_Application_Protocol_V1_1b.pdf (chapter 6 section 1 page 12)
// Quantity of Coils to read (2 bytes): 1 to 2000 (0x7D0)
// (chapter 6 section 11 page 29)
// Quantity of Coils to write (2 bytes): 1 to 1968 (0x7B0)
const (
	MODBUS_MAX_READ_BITS  = C.MODBUS_MAX_READ_BITS
	MODBUS_MAX_WRITE_BITS = C.MODBUS_MAX_WRITE_BITS
)

// Modbus_Application_Protocol_V1_1b.pdf (chapter 6 section 3 page 15)
// Quantity of Registers to read (2 bytes): 1 to 125 (0x7D)
// (chapter 6 section 12 page 31)
// Quantity of Registers to write (2 bytes) 1 to 123 (0x7B)
// (chapter 6 section 17 page 38)
// Quantity of Registers to write in R/W registers (2 bytes) 1 to 121 (0x79)
const (
	MODBUS_MAX_READ_REGISTERS     = C.MODBUS_MAX_READ_REGISTERS
	MODBUS_MAX_WRITE_REGISTERS    = C.MODBUS_MAX_WRITE_REGISTERS
	MODBUS_MAX_WR_WRITE_REGISTERS = C.MODBUS_MAX_WR_WRITE_REGISTERS
	MODBUS_MAX_WR_READ_REGISTERS  = C.MODBUS_MAX_WR_READ_REGISTERS
)

// MODBUS_MAX_PDU_LENGTH The size of the MODBUS PDU is limited by the size constraint inherited from
// the first MODBUS implementation on Serial Line network (max. RS485 ADU = 256
// bytes). Therefore, MODBUS PDU for serial line communication = 256 - Server
// address (1 byte) - CRC (2 bytes) = 253 bytes.
const MODBUS_MAX_PDU_LENGTH = C.MODBUS_MAX_PDU_LENGTH

// MODBUS_MAX_ADU_LENGTH
//
// Consequently:
//   - RTU MODBUS ADU = 253 bytes + Server address (1 byte) + CRC (2 bytes) = 256
//     bytes.
//   - TCP MODBUS ADU = 253 bytes + MBAP (7 bytes) = 260 bytes.
//
// so the maximum of both backend in 260 bytes. This size can used to allocate
// an array of bytes to store responses and it will be compatible with the two
// backends.
const MODBUS_MAX_ADU_LENGTH = C.MODBUS_MAX_ADU_LENGTH

// MODBUS_ENOBASE Random number to avoid errno conflicts
const MODBUS_ENOBASE = C.MODBUS_ENOBASE

// ModbusException Protocol exceptions
type ModbusException int

const (
	MODBUS_EXCEPTION_ILLEGAL_FUNCTION        ModbusException = C.MODBUS_EXCEPTION_ILLEGAL_FUNCTION
	MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS    ModbusException = C.MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS
	MODBUS_EXCEPTION_ILLEGAL_DATA_VALUE      ModbusException = C.MODBUS_EXCEPTION_ILLEGAL_DATA_VALUE
	MODBUS_EXCEPTION_SLAVE_OR_SERVER_FAILURE ModbusException = C.MODBUS_EXCEPTION_SLAVE_OR_SERVER_FAILURE
	MODBUS_EXCEPTION_ACKNOWLEDGE             ModbusException = C.MODBUS_EXCEPTION_ACKNOWLEDGE
	MODBUS_EXCEPTION_SLAVE_OR_SERVER_BUSY    ModbusException = C.MODBUS_EXCEPTION_SLAVE_OR_SERVER_BUSY
	MODBUS_EXCEPTION_NEGATIVE_ACKNOWLEDGE    ModbusException = C.MODBUS_EXCEPTION_NEGATIVE_ACKNOWLEDGE
	MODBUS_EXCEPTION_MEMORY_PARITY           ModbusException = C.MODBUS_EXCEPTION_MEMORY_PARITY
	MODBUS_EXCEPTION_NOT_DEFINED             ModbusException = C.MODBUS_EXCEPTION_NOT_DEFINED
	MODBUS_EXCEPTION_GATEWAY_PATH            ModbusException = C.MODBUS_EXCEPTION_GATEWAY_PATH
	MODBUS_EXCEPTION_GATEWAY_TARGET          ModbusException = C.MODBUS_EXCEPTION_GATEWAY_TARGET
	MODBUS_EXCEPTION_MAX                     ModbusException = C.MODBUS_EXCEPTION_MAX
)

type ErrorCode int

func (c ErrorCode) Error() (e *Error) {
	msg := C.modbus_strerror(C.int(c))
	return &Error{
		code:    c,
		message: C.GoString(msg),
	}
}

const (
	EMBXILFUN  ErrorCode = C.EMBXILFUN
	EMBXILADD  ErrorCode = C.EMBXILADD
	EMBXILVAL  ErrorCode = C.EMBXILVAL
	EMBXSFAIL  ErrorCode = C.EMBXSFAIL
	EMBXACK    ErrorCode = C.EMBXACK
	EMBXSBUSY  ErrorCode = C.EMBXSBUSY
	EMBXNACK   ErrorCode = C.EMBXNACK
	EMBXMEMPAR ErrorCode = C.EMBXMEMPAR
	EMBXGPATH  ErrorCode = C.EMBXGPATH
	EMBXGTAR   ErrorCode = C.EMBXGTAR
)

// Native libmodbus error codes
const (
	EMBBADCRC   ErrorCode = C.EMBBADCRC
	EMBBADDATA  ErrorCode = C.EMBBADDATA
	EMBBADEXC   ErrorCode = C.EMBBADEXC
	EMBUNKEXC   ErrorCode = C.EMBUNKEXC
	EMBMDATA    ErrorCode = C.EMBMDATA
	EMBBADSLAVE ErrorCode = C.EMBBADSLAVE
)

var (
	LibmodbusVersionMajor = C.libmodbus_version_major
	LibmodbusVersionMinor = C.libmodbus_version_minor
	LibmodbusVersionMicro = C.libmodbus_version_micro
)

type Modbus struct {
	ctx *C.modbus_t
}

type ModbusMapping struct {
	mb *C.modbus_mapping_t
}

type ModbusErrorRecoveryMode byte

const (
	MODBUS_ERROR_RECOVERY_NONE     ModbusErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_NONE
	MODBUS_ERROR_RECOVERY_LINK     ModbusErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_LINK
	MODBUS_ERROR_RECOVERY_PROTOCOL ModbusErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_PROTOCOL
)

type ModbusQuirks byte

const (
	MODBUS_QUIRK_NONE               ModbusQuirks = C.MODBUS_QUIRK_NONE
	MODBUS_QUIRK_MAX_SLAVE          ModbusQuirks = C.MODBUS_QUIRK_MAX_SLAVE
	MODBUS_QUIRK_REPLY_TO_BROADCAST ModbusQuirks = C.MODBUS_QUIRK_REPLY_TO_BROADCAST
	MODBUS_QUIRK_ALL                ModbusQuirks = C.MODBUS_QUIRK_ALL
)

type Error struct {
	code    ErrorCode
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d message:%s", e.code, e.message)
}

func (e *Error) Code() ErrorCode {
	return e.code
}

type ReportSlaveId struct {
	SlaveId            byte
	RunIndicatorStatus byte
	AdditionalData     []byte
}
