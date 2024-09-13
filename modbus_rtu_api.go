package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib -lmodbus -Wl,-rpath=/usr/local/lib
#include <stdlib.h>
#include "modbus.h"

extern void set_rts_cgo(modbus_t *ctx, int on);
*/
import "C"
import (
	"time"
	"unsafe"
)

// modbus_new_rtu modbus_new_rtu - create a libmodbus context for RTU
//
// The modbus_new_rtu() function shall allocate and initialize a modbus_t structure to communicate in RTU mode on a
// serial line.
//
// The device argument specifies the name of the serial port handled by the OS, eg. "/dev/ttyS0" or "/dev/ttyUSB0". On
// Windows, it's necessary to prepend COM name with "\.\" for COM number greater than 9, eg. "\\.\COM10". See http://
// msdn.microsoft.com/en-us/library/aa365247(v=vs.85).aspx for details
//
// The baud argument specifies the baud rate of the communication, eg. 9600, 19200, 57600, 115200, etc.
//
// The parity argument can have one of the following values:
//
//   - N for none
//   - E for even
//   - O for odd
//
// The data_bits argument specifies the number of bits of data, the allowed values are 5, 6, 7 and 8.
//
// The stop_bits argument specifies the bits of stop, the allowed values are 1 and 2.
//
// Once the modbus_t structure is initialized, you can connect to the serial bus with modbus_connect.
//
// In RTU, your program can act as server or client:
//
// server is called slave in Modbus terminology, your program will expose data to the network by processing and
// answering the requests of one of several clients. It up to you to define the slave ID of your service with
// modbus_set_slave, this ID should be used by the client to communicate with your program.
//
// client is called master in Modbus terminology, your program will send requests to servers to read or write data from
// them. Before issuing the requests, you should define the slave ID of the remote device with modbus_set_slave. The
// slave ID is not an argument of the read/write functions because it's very frequent to talk with only one server so
// you can set it once and for all. The slave ID it not used in TCP communications so this way the API is common to
// both.
func ModbusNewRtu(device string, baud int, parity byte, dataBit int, stopBit int) *Modbus {
	cdevice := C.CString(device)
	defer C.free(unsafe.Pointer(cdevice))
	ctx := C.modbus_new_rtu(cdevice, C.int(baud), C.char(parity), C.int(dataBit), C.int(stopBit))
	if ctx == nil {
		return nil
	}
	return &Modbus{ctx: ctx}
}

// ModbusRtuSetSerialMode modbus_rtu_set_serial_mode - set the serial mode
//
// The modbus_rtu_set_serial_mode() function shall set the selected serial mode:
//
//   - MODBUS_RTU_RS232, the serial line is set for RS-232 communication. RS-232 (Recommended Standard 232) is the
//     traditional name for a series of standards for serial binary single-ended data and control signals connecting
//     between a DTE (Data Terminal Equipment) and a DCE (Data Circuit-terminating Equipment). It is commonly used in
//     computer serial ports.
//
//   - MODBUS_RTU_RS485, the serial line is set for RS-485 communication. EIA-485, also known as TIA/EIA-485 or RS-485,
//     is a standard defining the electrical characteristics of drivers and receivers for use in balanced digital
//     multipoint systems. This standard is widely used for communications in industrial automation because it can be
//     used effectively over long distances and in electrically noisy environments.
//
// This function is only supported on Linux kernels 2.6.28 onwards.
func (x *Modbus) ModbusRtuSetSerialMode(mode int) (err error) {
	code := C.modbus_rtu_set_serial_mode(x.ctx, C.int(mode))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusRtuGetSerialMode modbus_rtu_get_serial_mode - get the current serial mode
//
// The modbus_rtu_get_serial_mode() function shall return the serial mode currently used by the libmodbus context:
//
//   - MODBUS_RTU_RS232, the serial line is set for RS-232 communication. RS-232 (Recommended Standard 232) is the
//     traditional name for a series of standards for serial binary single-ended data and control signals connecting
//     between a DTE (Data Terminal Equipment) and a DCE (Data Circuit-terminating Equipment). It is commonly used in
//     computer serial ports
//
//   - MODBUS_RTU_RS485, the serial line is set for RS-485 communication. EIA-485, also known as TIA/EIA-485 or RS-485,
//     is a standard defining the electrical characteristics of drivers and receivers for use in balanced digital
//     multipoint systems. This standard is widely used for communications in industrial automation because it can be
//     used effectively over long distances and in electrically noisy environments. This function is only available on
//     Linux kernels 2.6.28 onwards and can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuGetSerialMode() (mode int, err error) {
	code := C.modbus_rtu_get_serial_mode(x.ctx)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	mode = int(code)
	return
}

// ModbusRtuSetRts modbus_rtu_set_rts - set the RTS mode in RTU
//
// The modbus_rtu_set_rts() function shall set the Request To Send mode to communicate on a RS-485 serial bus. By
// default, the mode is set to MODBUS_RTU_RTS_NONE and no signal is issued before writing data on the wire.
//
// To enable the RTS mode, the values MODBUS_RTU_RTS_UP or MODBUS_RTU_RTS_DOWN must be used, these modes enable the RTS
// mode and set the polarity at the same time. When MODBUS_RTU_RTS_UP is used, an ioctl call is made with RTS flag
// enabled then data is written on the bus after a delay of 1 ms, then another ioctl call is made with the RTS flag
// disabled and again a delay of 1 ms occurs. The MODBUS_RTU_RTS_DOWN mode applies the same procedure but with an
// inverted RTS flag.
//
// This function can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuSetRts(mode int) (err error) {
	code := C.modbus_rtu_set_rts(x.ctx, C.int(mode))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusRtuGetRts modbus_rtu_get_rts - get the current RTS mode in RTU
//
// The modbus_rtu_get_rts() function shall get the current Request To Send mode of the libmodbus context ctx. The
// possible returned values are:
//
//   - MODBUS_RTU_RTS_NONE
//   - MODBUS_RTU_RTS_UP
//   - MODBUS_RTU_RTS_DOWN
//
// This function can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuGetRts() (mode int, err error) {
	code := C.modbus_rtu_get_rts(x.ctx)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	mode = int(code)
	return
}

// ModbusRtuSetCustomRts modbus_rtu_set_custom_rts - set a function to be used for custom RTS implementation
//
// The modbus_rtu_set_custom_rts() function shall set a custom function to be called when the RTS pin is to be set
// before and after a transmission. By default this is set to an internal function that toggles the RTS pin using an
// ioctl call.
//
// Note that this function adheres to the RTS mode, the values MODBUS_RTU_RTS_UP or MODBUS_RTU_RTS_DOWN must be used
// for the function to be called.
//
// This function can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuSetCustomRts(cb SetRtsCallback) (err error) {
	mapSetRtsCallback.Store(x.ctx, cb)
	code := C.modbus_rtu_set_custom_rts(x.ctx, (C.set_rts)(C.set_rts_cgo))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusRtuSetRtsDelay modbus_rtu_set_rts_delay - set the RTS delay in RTU
//
// The modbus_rtu_set_rts_delay() function shall set the Request To Send delay period of the libmodbus context 'ctx'.
//
// This function can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuSetRtsDelay(us time.Duration) (err error) {
	code := C.modbus_rtu_set_rts_delay(x.ctx, C.int(us))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusRtuGetRtsDelay modbus_rtu_get_rts_delay - get the current RTS delay in RTU
//
// The modbus_rtu_get_rts_delay() function shall get the current Request To Send delay period of the libmodbus context
// 'ctx'.
//
// This function can only be used with a context using a RTU backend.
func (x *Modbus) ModbusRtuGetRtsDelay() (us time.Duration, err error) {
	code := C.modbus_rtu_get_rts_delay(x.ctx)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	us = time.Duration(code) * time.Microsecond
	return
}
