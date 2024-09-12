package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib -lmodbus -Wl,-rpath=/usr/local/lib
#include "modbus.h"
*/
import "C"
import "time"

// ModbusStrError modbus_strerror - return the error message
//
// The modbus_strerror() function shall return a pointer to an error message string corresponding to
// the error number specified by the errnum argument. As libmodbus defines additional error
// numbers over and above those defined by the operating system, applications should use
// modbus_strerror() in preference to the standard strerror() function.
func ModbusStrError(errnum int) string {
	return C.GoString(C.modbus_strerror(C.int(errnum)))
}

// ModbusSetSlave modbus_set_slave - set slave number in the context
//
// The modbus_set_slave() function shall set the slave number in the libmodbus context.
//
// It is usually only required to set the slave ID in RTU. The meaning of this ID will be different
// if your program acts as client (master) or server (slave).
//
// As RTU client, modbus_set_slave() sets the ID of the remote device you want to communicate.
// Be sure to set the slave ID before issuing any Modbus requests on the serial bus. If you
// communicate with several servers (slaves), you can set the slave ID of the remote device before
// each request.
//
// As RTU server, the slave ID allows the various clients to reach your service. You should use a
// free ID, once set, this ID should be known by the clients of the network. According to the
// protocol, a Modbus device must only accept message holding its slave number or the special
// broadcast number.
//
// In TCP, the slave number is only required if the message must reach a device on a serial
// network. Some not compliant devices or software (such as modpoll) uses the slave ID as unit
// identifier, that's incorrect (cf page 23 of Modbus Messaging Implementation Guide v1.0b) but
// without the slave value, the faulty remote device or software drops the requests! The special
// value MODBUS_TCP_SLAVE (0xFF) can be used in TCP mode to restore the default value.
//
// The broadcast address is MODBUS_BROADCAST_ADDRESS. This special value must be use when
// you want all Modbus devices of the network receive the request.
func (x *Modbus) ModbusSetSlave(slave int) (err error) {
	code := C.modbus_set_slave(x.ctx, C.int(slave))
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusGetSlave modbus_get_slave - get slave number in the context
//
// The modbus_get_slave() function shall get the slave number in the libmodbus context.
func (x *Modbus) ModbusGetSlave() (slave int, err error) {
	code := C.modbus_get_slave(x.ctx)
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	slave = int(code)
	return
}

// ModbusSetErrorRecovery modbus_set_error_recovery - set the error recovery mode
//
// The modbus_set_error_recovery() function shall set the error recovery mode to apply when the connection
// fails or the byte received is not expected. The argument error_recovery may be bitwise-or'ed with zero or
// more of the following constants.
//
// By default there is no error recovery (MODBUS_ERROR_RECOVERY_NONE) so the application is responsible for
// controlling the error values returned by libmodbus functions and for handling them if necessary.
//
// When MODBUS_ERROR_RECOVERY_LINK is set, the library will attempt an reconnection after a delay defined by
// response timeout of the libmodbus context. This mode will try an infinite close/connect loop until success
// on send call and will just try one time to re-establish the connection on select/read calls (if the connection
// was down, the values to read are certainly not available any more after reconnection, except for slave/server).
// This mode will also run flush requests after a delay based on the current response timeout in some situations
// (eg. timeout of select call). The reconnection attempt can hang for several seconds if the network to the remote
// target unit is down.
//
// When MODBUS_ERROR_RECOVERY_PROTOCOL is set, a sleep and flush sequence will be used to clean up the ongoing
// communication, this can occurs when the message length is invalid, the TID is wrong or the received function
// code is not the expected one. The response timeout delay will be used to sleep.
//
// The modes are mask values and so they are complementary.
//
// It's not recommended to enable error recovery for a Modbus slave/server.
func (x *Modbus) ModbusSetErrorRecovery(errorRecovery ModbusErrorRecoveryMode) (err error) {
	code := C.modbus_set_error_recovery(x.ctx, C.modbus_error_recovery_mode(errorRecovery))
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusConnect modbus_connect - establish a Modbus connection
//
// The modbus_connect() function shall establish a connection to a Modbus server, a network or a bus
// using the context information of libmodbus context given in argument.
func (x *Modbus) ModbusConnect() (err error) {
	code := C.modbus_connect(x.ctx)
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusSetSocket modbus_set_socket - set socket of the context
//
// The modbus_set_socket() function shall set the socket or file descriptor in the libmodbus context.
// This function is useful for managing multiple client connections to the same server.
func (x *Modbus) ModbusSetSocket(s int) (err error) {
	code := C.modbus_set_socket(x.ctx, C.int(s))
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusGetSocket modbus_get_socket - get the current socket of the context
//
// The modbus_get_socket() function shall return the current socket or file descriptor of the libmodbus context.
func (x *Modbus) ModbusGetSocket() (s int, err error) {
	code := C.modbus_get_socket(x.ctx)
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	s = int(code)
	return
}

// ModbusSetResponseTimeout modbus_set_response_timeout - set timeout for response
//
// The modbus_set_response_timeout() function shall set the timeout interval used to wait for a response.
// When a byte timeout is set, if elapsed time for the first byte of response is longer than the given timeout,
// an ETIMEDOUT error will be raised by the function waiting for a response. When byte timeout is disabled,
// the full confirmation response must be received before expiration of the response timeout.
//
// The value of to_usec argument must be in the range 0 to 999999.
func (x *Modbus) ModbusSetResponseTimeout(timeout time.Duration) (err error) {
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_response_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusGetResponseTimeout modbus_get_response_timeout - get timeout for response
//
// The modbus_get_response_timeout() function shall return the timeout interval used to wait for a response
// in the to_sec and to_usec arguments.
func (x *Modbus) ModbusGetResponseTimeout() (timeout time.Duration, err error) {
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_response_timeout(x.ctx, &to_sec, &to_usec)
	if code <= 0 {
		err = ErrorCode(code).Error()
		return
	}
	timeout = time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond
	return
}

// ModbusFree modbus_free - free a libmodbus context
//
// The modbus_free() function shall free an allocated modbus_t structure.
func (x *Modbus) ModbusFree() {
	C.modbus_free(x.ctx)
}
