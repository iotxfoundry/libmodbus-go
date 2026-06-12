package modbus

/*
#include "modbus.h"
*/
import "C"
import (
	"fmt"
	"sync"
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

// Exception Protocol exceptions
type Exception int

const (
	MODBUS_EXCEPTION_ILLEGAL_FUNCTION        Exception = C.MODBUS_EXCEPTION_ILLEGAL_FUNCTION
	MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS    Exception = C.MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS
	MODBUS_EXCEPTION_ILLEGAL_DATA_VALUE      Exception = C.MODBUS_EXCEPTION_ILLEGAL_DATA_VALUE
	MODBUS_EXCEPTION_SLAVE_OR_SERVER_FAILURE Exception = C.MODBUS_EXCEPTION_SLAVE_OR_SERVER_FAILURE
	MODBUS_EXCEPTION_ACKNOWLEDGE             Exception = C.MODBUS_EXCEPTION_ACKNOWLEDGE
	MODBUS_EXCEPTION_SLAVE_OR_SERVER_BUSY    Exception = C.MODBUS_EXCEPTION_SLAVE_OR_SERVER_BUSY
	MODBUS_EXCEPTION_NEGATIVE_ACKNOWLEDGE    Exception = C.MODBUS_EXCEPTION_NEGATIVE_ACKNOWLEDGE
	MODBUS_EXCEPTION_MEMORY_PARITY           Exception = C.MODBUS_EXCEPTION_MEMORY_PARITY
	MODBUS_EXCEPTION_NOT_DEFINED             Exception = C.MODBUS_EXCEPTION_NOT_DEFINED
	MODBUS_EXCEPTION_GATEWAY_PATH            Exception = C.MODBUS_EXCEPTION_GATEWAY_PATH
	MODBUS_EXCEPTION_GATEWAY_TARGET          Exception = C.MODBUS_EXCEPTION_GATEWAY_TARGET
	MODBUS_EXCEPTION_MAX                     Exception = C.MODBUS_EXCEPTION_MAX
)

type ErrorCode int

func (c ErrorCode) Error() string {
	return C.GoString(C.modbus_strerror(C.int(c)))
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
	VersionMajor = C.libmodbus_version_major
	VersionMinor = C.libmodbus_version_minor
	VersionMicro = C.libmodbus_version_micro
)

type Modbus struct {
	mu     sync.Mutex
	ctx    *C.modbus_t
	socket int
	closed bool
}

func (x *Modbus) ensureCtx() error {
	if x.ctx == nil {
		return fmt.Errorf("modbus: context is nil (already freed)")
	}
	return nil
}

type ModbusMapping struct {
	mu sync.Mutex
	mb *C.modbus_mapping_t
}

type ErrorRecoveryMode byte

const (
	MODBUS_ERROR_RECOVERY_NONE     ErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_NONE
	MODBUS_ERROR_RECOVERY_LINK     ErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_LINK
	MODBUS_ERROR_RECOVERY_PROTOCOL ErrorRecoveryMode = C.MODBUS_ERROR_RECOVERY_PROTOCOL
)

type Quirks byte

const (
	MODBUS_QUIRK_NONE               Quirks = C.MODBUS_QUIRK_NONE
	MODBUS_QUIRK_MAX_SLAVE          Quirks = C.MODBUS_QUIRK_MAX_SLAVE
	MODBUS_QUIRK_REPLY_TO_BROADCAST Quirks = C.MODBUS_QUIRK_REPLY_TO_BROADCAST
	MODBUS_QUIRK_ALL                Quirks = C.MODBUS_QUIRK_ALL
)

type Error struct {
	code    ErrorCode
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.code, e.message)
}

func (e *Error) Code() ErrorCode {
	return e.code
}

type SlaveIDReport struct {
	SlaveId            byte
	RunIndicatorStatus byte
	AdditionalData     []byte
}
