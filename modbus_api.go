package modbus

/*
#include "modbus.h"
#include <unistd.h>

extern int get_errno_cgo();
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"runtime"
	"time"
	"unsafe"
)

// newCError creates an Error from the current C errno.
// MUST be called immediately after the failing C call.
func newCError() error {
	errno := C.get_errno_cgo()
	return &Error{
		code:    ErrorCode(errno),
		message: C.GoString(C.modbus_strerror(errno)),
	}
}

// SetSlave modbus_set_slave - set slave number in the context
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
func (x *Modbus) SetSlave(slave int) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_set_slave(x.ctx, C.int(slave))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetSlave modbus_get_slave - get slave number in the context
//
// The modbus_get_slave() function shall get the slave number in the libmodbus context.
func (x *Modbus) GetSlave() (slave int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0, err
	}
	code := C.modbus_get_slave(x.ctx)
	if code < 0 {
		return 0, newCError()
	}
	return int(code), nil
}

// SetErrorRecovery modbus_set_error_recovery - set the error recovery mode
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
func (x *Modbus) SetErrorRecovery(errorRecovery ErrorRecoveryMode) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_set_error_recovery(x.ctx, C.modbus_error_recovery_mode(errorRecovery))
	if code < 0 {
		return newCError()
	}
	return nil
}

// Connect modbus_connect - establish a Modbus connection
//
// The modbus_connect() function shall establish a connection to a Modbus server, a network or a bus
// using the context information of libmodbus context given in argument.
func (x *Modbus) Connect() (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_connect(x.ctx)
	if code < 0 {
		return newCError()
	}
	x.closed = false
	return nil
}

// SetSocket modbus_set_socket - set socket of the context
//
// The modbus_set_socket() function shall set the socket or file descriptor in the libmodbus context.
// This function is useful for managing multiple client connections to the same server.
func (x *Modbus) SetSocket(s int) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_set_socket(x.ctx, C.int(s))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetSocket modbus_get_socket - get the current socket of the context
//
// The modbus_get_socket() function shall return the current socket or file descriptor of the libmodbus context.
func (x *Modbus) GetSocket() (s int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0, err
	}
	code := C.modbus_get_socket(x.ctx)
	if code < 0 {
		return 0, newCError()
	}
	return int(code), nil
}

// SetResponseTimeout modbus_set_response_timeout - set timeout for response
//
// The modbus_set_response_timeout() function shall set the timeout interval used to wait for a response.
// When a byte timeout is set, if elapsed time for the first byte of response is longer than the given timeout,
// an ETIMEDOUT error will be raised by the function waiting for a response. When byte timeout is disabled,
// the full confirmation response must be received before expiration of the response timeout.
//
// The value of to_usec argument must be in the range 0 to 999999.
func (x *Modbus) SetResponseTimeout(timeout time.Duration) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	if timeout < 0 {
		return fmt.Errorf("modbus: response timeout must be non-negative, got %v", timeout)
	}
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_response_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetResponseTimeout modbus_get_response_timeout - get timeout for response
//
// The modbus_get_response_timeout() function shall return the timeout interval used to wait for a response
// in the to_sec and to_usec arguments.
func (x *Modbus) GetResponseTimeout() (timeout time.Duration, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0, err
	}
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_response_timeout(x.ctx, &to_sec, &to_usec)
	if code < 0 {
		return 0, newCError()
	}
	return time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond, nil
}

// SetByteTimeout modbus_set_byte_timeout - set timeout between bytes
//
// The modbus_set_byte_timeout() function shall set the timeout interval between two consecutive
// bytes of the same message. The timeout is an upper bound on the amount of time elapsed before select()
// returns, if the time elapsed is longer than the defined timeout, an ETIMEDOUT error will be raised by
// the function waiting for a response.
//
// The value of to_usec argument must be in the range 0 to 999999.
//
// If both to_sec and to_usec are zero, this timeout will not be used at all. In this case,
// modbus_set_response_timeout() governs the entire handling of the response, the full confirmation
// response must be received before expiration of the response timeout. When a byte timeout is set,
// the response timeout is only used to wait for until the first byte of the response.
func (x *Modbus) SetByteTimeout(timeout time.Duration) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	if timeout < 0 {
		return fmt.Errorf("modbus: byte timeout must be non-negative, got %v", timeout)
	}
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_byte_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetByteTimeout modbus_get_byte_timeout - get timeout between bytes
//
// The modbus_get_byte_timeout() function shall store the timeout interval between two
// consecutive bytes of the same message in the to_sec and to_usec arguments.
func (x *Modbus) GetByteTimeout() (timeout time.Duration, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0, err
	}
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_byte_timeout(x.ctx, &to_sec, &to_usec)
	if code < 0 {
		return 0, newCError()
	}
	return time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond, nil
}

// SetIndicationTimeout modbus_set_indication_timeout - set timeout between indications
//
// The modbus_set_indication_timeout() function shall set the timeout interval used by a server
// to wait for a request from a client.
//
// The value of to_usec argument must be in the range 0 to 999999.
//
// If both to_sec and to_usec are zero, this timeout will not be used at all. In this case,
// the server will wait forever.
func (x *Modbus) SetIndicationTimeout(timeout time.Duration) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	if timeout < 0 {
		return fmt.Errorf("modbus: indication timeout must be non-negative, got %v", timeout)
	}
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_indication_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetIndicationTimeout modbus_get_indication_timeout - get timeout used to wait for an indication
// (request received by a server).
//
// The modbus_get_indication_timeout() function shall store the timeout interval used to wait for an
// indication in the to_sec and to_usec arguments. Indication is the term used by the Modbus protocol
// to designate a request received by the server.
//
// The default value is zero, it means the server will wait forever.
func (x *Modbus) GetIndicationTimeout() (timeout time.Duration, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0, err
	}
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_indication_timeout(x.ctx, &to_sec, &to_usec)
	if code < 0 {
		return 0, newCError()
	}
	return time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond, nil
}

// GetHeaderLength modbus_get_header_length - retrieve the current header length
//
// The modbus_get_header_length() function shall retrieve the current header length from
// the backend. This function is convenient to manipulate a message and so it's limited to
// low-level operations.
func (x *Modbus) GetHeaderLength() (length int) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return 0
	}
	return int(C.modbus_get_header_length(x.ctx))
}

// Free modbus_free - free a libmodbus context
//
// The modbus_free() function shall free an allocated modbus_t structure.
// It is safe to call Free multiple times.
func (x *Modbus) Free() {
	x.mu.Lock()
	defer x.mu.Unlock()
	if x.ctx == nil {
		return
	}
	mapSetRtsCallback.Delete(unsafe.Pointer(x.ctx))
	C.modbus_free(x.ctx)
	x.ctx = nil
	runtime.SetFinalizer(x, nil)
}

// Close modbus_close - close a Modbus connection
//
// The modbus_close() function shall close the connection established with the backend set in the context.
// It is safe to call Close multiple times.
func (x *Modbus) Close() error {
	x.mu.Lock()
	defer x.mu.Unlock()
	if x.closed || x.ctx == nil {
		return nil
	}
	C.modbus_close(x.ctx)
	x.closed = true
	return nil
}

// Destroy releases all resources: closes the connection and frees the context.
// It is safe to call Destroy multiple times.
// This is the recommended way to clean up a Modbus instance.
func (x *Modbus) Destroy() {
	_ = x.Close()
	x.Free()
}

// Flush modbus_flush - flush non-transmitted data
//
// The modbus_flush() function shall discard data received but not read to the socket or file descriptor associated
// to the context 'ctx'.
func (x *Modbus) Flush() (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_flush(x.ctx)
	if code < 0 {
		return newCError()
	}
	return nil
}

// SetDebug modbus_set_debug - set debug flag of the context
//
// The modbus_set_debug() function shall set the debug flag of the modbus_t context by using the argument flag.
// By default, the boolean flag is set to FALSE. When the flag value is set to TRUE, many verbose messages are
// displayed on stdout and stderr. For example, this flag is useful to display the bytes of the Modbus messages.
func (x *Modbus) SetDebug(flag bool) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	f := C.FALSE
	if flag {
		f = C.TRUE
	}
	code := C.modbus_set_debug(x.ctx, C.int(f))
	if code < 0 {
		return newCError()
	}
	return nil
}

// ReadBits modbus_read_bits - read many bits (coils)
//
// The modbus_read_bits() function shall read the status of the nb bits (coils) to the address addr of the remote
// device. The result of reading is stored in dest array as unsigned bytes (8 bits) set to TRUE or FALSE.
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint8_t)).
//
// The function uses the Modbus function code 0x01 (read coil status).
func (x *Modbus) ReadBits(addr int, nb int) (out []byte, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	dest := make([]C.uint8_t, nb)
	code := C.modbus_read_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return nil, newCError()
	}
	out = cUint8ToBytes(dest)
	return out, nil
}

// ReadInputBits modbus_read_input_bits - read many input bits
//
// The modbus_read_input_bits() function shall read the content of the nb input bits to the address addr of the remote
// device. The result of reading is stored in dest array as unsigned bytes (8 bits) set to TRUE or FALSE.
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint8_t)).
//
// The function uses the Modbus function code 0x02 (read input status).
func (x *Modbus) ReadInputBits(addr int, nb int) (out []byte, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	dest := make([]C.uint8_t, nb)
	code := C.modbus_read_input_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return nil, newCError()
	}
	out = cUint8ToBytes(dest)
	return out, nil
}

// ReadRegisters modbus_read_registers - read many registers
//
// The modbus_read_registers() function shall read the content of the nb holding registers to the address addr of the
// remote device. The result of reading is stored in dest array as word values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x03 (read holding registers).
func (x *Modbus) ReadRegisters(addr int, nb int) (out []uint16, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	dest := make([]C.uint16_t, nb)
	code := C.modbus_read_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return nil, newCError()
	}
	out = cUint16ToSlice(dest)
	return out, nil
}

// ReadInputRegisters modbus_read_input_registers - read many input registers
//
// The modbus_read_input_registers() function shall read the content of the nb input registers to address addr of the
// remote device. The result of the reading is stored in dest array as word values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x04 (read input registers). The holding registers and input registers
// have different historical meaning, but nowadays it's more common to use holding registers only.
func (x *Modbus) ReadInputRegisters(addr int, nb int) (out []uint16, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	dest := make([]C.uint16_t, nb)
	code := C.modbus_read_input_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return nil, newCError()
	}
	out = cUint16ToSlice(dest)
	return out, nil
}

// WriteBit modbus_write_bit - write a single bit
//
// The modbus_write_bit() function shall write the status of status at the address addr of the remote device. The value
// must be set to TRUE or FALSE.
//
// The function uses the Modbus function code 0x05 (force single coil).
func (x *Modbus) WriteBit(addr int, status bool) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	v := C.FALSE
	if status {
		v = C.TRUE
	}
	code := C.modbus_write_bit(x.ctx, C.int(addr), C.int(v))
	if code < 0 {
		return newCError()
	}
	return nil
}

// WriteRegister modbus_write_register - write a single register
//
// The modbus_write_register() function shall write the value of value holding registers at the address addr of the
// remote device.
//
// The function uses the Modbus function code 0x06 (preset single register).
func (x *Modbus) WriteRegister(addr int, value uint16) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_write_register(x.ctx, C.int(addr), C.uint16_t(value))
	if code < 0 {
		return newCError()
	}
	return nil
}

// WriteBits modbus_write_bits - write many bits
//
// The modbus_write_bits() function shall write the status of the nb bits (coils) from src at the address addr of the
// remote device. The src array must contains bytes set to TRUE or FALSE.
//
// The function uses the Modbus function code 0x0F (force multiple coils).
func (x *Modbus) WriteBits(addr int, data []byte) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	nb := len(data)
	dest := bytesToCUint8(data)
	code := C.modbus_write_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return newCError()
	}
	return nil
}

// WriteRegisters modbus_write_registers - write many registers
//
// The modbus_write_registers() function shall write the content of the nb holding registers from the array src at
// address addr of the remote device.
//
// The function uses the Modbus function code 0x10 (preset multiple registers).
func (x *Modbus) WriteRegisters(addr int, data []uint16) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	nb := len(data)
	dest := uint16ToCUint16(data)
	code := C.modbus_write_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		return newCError()
	}
	return nil
}

// MaskWriteRegister modbus_mask_write_register - mask a single register
//
// The modbus_mask_write_register() function shall modify the value of the holding register at the address 'addr' of
// the remote device using the algorithm:
//
// new value = (current value AND 'and') OR ('or' AND (NOT 'and'))
//
// The function uses the Modbus function code 0x16 (mask single register).
func (x *Modbus) MaskWriteRegister(addr int, andMask uint16, orMask uint16) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_mask_write_register(x.ctx, C.int(addr), C.uint16_t(andMask), C.uint16_t(orMask))
	if code < 0 {
		return newCError()
	}
	return nil
}

// WriteAndReadRegisters modbus_write_and_read_registers - write and read many registers in a single transaction
//
// The modbus_write_and_read_registers() function shall write the content of the write_nb holding registers from the
// array 'src' to the address write_addr of the remote device then shall read the content of the read_nb holding
// registers to the address read_addr of the remote device. The result of reading is stored in dest array as word
// values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x17 (write/read registers).
func (x *Modbus) WriteAndReadRegisters(writeAddr int, src []uint16, readAddr int, readNb int) (dest []uint16, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	writeNb := len(src)
	if writeNb == 0 || readNb == 0 {
		return nil, fmt.Errorf("modbus: write and read registers requires non-empty source and positive read count")
	}
	csrc := uint16ToCUint16(src)
	cdest := make([]C.uint16_t, readNb)
	code := C.modbus_write_and_read_registers(x.ctx, C.int(writeAddr), C.int(writeNb), (*C.uint16_t)(unsafe.SliceData(csrc)), C.int(readAddr), C.int(readNb), (*C.uint16_t)(unsafe.SliceData(cdest)))
	if code < 0 {
		return nil, newCError()
	}
	dest = cUint16ToSlice(cdest)
	return dest, nil
}

// ReportSlaveID modbus_report_slave_id - returns a description of the controller
//
// The modbus_report_slave_id() function shall send a request to the controller to obtain a description of the
// controller.
//
// The response stored in dest contains:
//
//   - the slave ID, this unique ID is in reality not unique at all so it's not possible to depend on it to know how the
//     information are packed in the response.
//   - the run indicator status (0x00 = OFF, 0xFF = ON)
//   - additional data specific to each controller. For example, libmodbus returns the version of the library as a
//     string.
//
// The function writes at most max_dest bytes from the response to dest so you must ensure that dest is large enough.
func (x *Modbus) ReportSlaveID() (dest *SlaveIDReport, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	cdest := make([]C.uint8_t, MODBUS_MAX_PDU_LENGTH)
	code := C.modbus_report_slave_id(x.ctx, C.int(MODBUS_MAX_PDU_LENGTH), unsafe.SliceData(cdest))
	if code < 0 {
		return nil, newCError()
	}
	if code < 2 {
		return nil, fmt.Errorf("modbus: report slave id returned insufficient data (%d bytes)", code)
	}
	buff := cUint8ToBytes(cdest[:code:code])
	dest = &SlaveIDReport{
		SlaveId:            buff[0],
		RunIndicatorStatus: buff[1],
		AdditionalData:     buff[2:],
	}
	return dest, nil
}

// NewMappingWithStart modbus_mapping_new_start_address - allocate four arrays of bits and
// registers accessible from their starting addresses
//
// The modbus_mapping_new_start_address() function shall allocate four arrays to store bits, input bits, registers and
// inputs registers. The pointers are stored in modbus_mapping_t structure. All values of the arrays are initialized to
// zero.
//
// The different starting addresses make it possible to place the mapping at any address in each address space. This
// way, you can give access to clients to values stored at high addresses without allocating memory from the
// address zero, for example to make available registers from 340 to 349, you can use:
//
//	mb_mapping = modbus_mapping_new_start_address(0, 0, 0, 0, 340, 10, 0, 0);
//
// If it isn't necessary to allocate an array for a specific type of data, you can pass the zero value in argument,
// the associated pointer will be NULL.
//
// This function is convenient to handle requests in a Modbus server/slave.
func NewMappingWithStart(
	startBits uint,
	nbBits uint,
	startInputBits uint,
	nbInputBits uint,
	startRegisters uint,
	nbRegisters uint,
	startInputRegisters uint,
	nbInputRegisters uint,
) (*ModbusMapping, error) {
	mn := C.modbus_mapping_new_start_address(
		C.uint(startBits),
		C.uint(nbBits),
		C.uint(startInputBits),
		C.uint(nbInputBits),
		C.uint(startRegisters),
		C.uint(nbRegisters),
		C.uint(startInputRegisters),
		C.uint(nbInputRegisters),
	)
	if mn == nil {
		return nil, newCError()
	}
	mm := &ModbusMapping{mb: mn}
	runtime.SetFinalizer(mm, (*ModbusMapping).Free)
	return mm, nil
}

// NewMapping modbus_mapping_new - allocate four arrays of bits and registers
//
// The modbus_mapping_new() function shall allocate four arrays to store bits, input bits, registers and inputs
// registers. The pointers are stored in modbus_mapping_t structure. All values of the arrays are initialized to zero.
//
// This function is equivalent to a call of the modbus_mapping_new_start_address function with all start addresses to 0.
//
// If it isn't necessary to allocate an array for a specific type of data, you can pass the zero value in argument,
// the associated pointer will be NULL.
//
// This function is convenient to handle requests in a Modbus server/slave.
func NewMapping(nbBits int, nbInputBits int, nbRegisters int, nbInputRegisters int) (*ModbusMapping, error) {
	mn := C.modbus_mapping_new(C.int(nbBits), C.int(nbInputBits), C.int(nbRegisters), C.int(nbInputRegisters))
	if mn == nil {
		return nil, newCError()
	}
	mm := &ModbusMapping{mb: mn}
	runtime.SetFinalizer(mm, (*ModbusMapping).Free)
	return mm, nil
}

type mapping struct {
	StartBits           int      `json:"start_bits"`
	StartInputBits      int      `json:"start_input_bits"`
	StartInputRegisters int      `json:"start_input_registers"`
	StartRegisters      int      `json:"start_registers"`
	Bits                []byte   `json:"bits"`
	InputBits           []byte   `json:"input_bits"`
	InputRegisters      []uint16 `json:"input_registers"`
	Registers           []uint16 `json:"registers"`
}

// MarshalJSON implements json.Marshaler for ModbusMapping.
func (mm *ModbusMapping) MarshalJSON() ([]byte, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return nil, fmt.Errorf("modbus: mapping is nil")
	}
	mp := &mapping{
		StartBits:           mm.startBits(),
		StartInputBits:      mm.startInputBits(),
		StartInputRegisters: mm.startInputRegisters(),
		StartRegisters:      mm.startRegisters(),
	}
	mp.Bits = mm.tabBitsSlice()
	mp.InputBits = mm.tabInputBitsSlice()
	mp.InputRegisters = mm.tabInputRegistersSlice()
	mp.Registers = mm.tabRegistersSlice()
	return json.Marshal(mp)
}

// UnmarshalJSON implements json.Unmarshaler for ModbusMapping.
func (mm *ModbusMapping) UnmarshalJSON(data []byte) error {
	mp := &mapping{}
	buffer := bytes.NewBuffer(data)
	dec := json.NewDecoder(buffer)
	err := dec.Decode(mp)
	if err != nil {
		return err
	}
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb != nil {
		C.modbus_mapping_free(mm.mb)
		mm.mb = nil
	}
	mn := C.modbus_mapping_new_start_address(
		C.uint(mp.StartBits),
		C.uint(len(mp.Bits)),
		C.uint(mp.StartInputBits),
		C.uint(len(mp.InputBits)),
		C.uint(mp.StartRegisters),
		C.uint(len(mp.Registers)),
		C.uint(mp.StartInputRegisters),
		C.uint(len(mp.InputRegisters)),
	)
	if mn == nil {
		return fmt.Errorf("modbus: new mapping error")
	}
	mm.mb = mn
	runtime.SetFinalizer(mm, (*ModbusMapping).Free)
	for k, v := range mp.Bits {
		tab := unsafe.Slice(mm.mb.tab_bits, int(mm.mb.nb_bits))
		tab[k] = C.uint8_t(v)
	}
	for k, v := range mp.InputBits {
		tab := unsafe.Slice(mm.mb.tab_input_bits, int(mm.mb.nb_input_bits))
		tab[k] = C.uint8_t(v)
	}
	for k, v := range mp.Registers {
		tab := unsafe.Slice(mm.mb.tab_registers, int(mm.mb.nb_registers))
		tab[k] = C.uint16_t(v)
	}
	for k, v := range mp.InputRegisters {
		tab := unsafe.Slice(mm.mb.tab_input_registers, int(mm.mb.nb_input_registers))
		tab[k] = C.uint16_t(v)
	}
	return nil
}

// helper methods for internal use (must be called with lock held)

func (mm *ModbusMapping) startBits() int           { return int(mm.mb.start_bits) }
func (mm *ModbusMapping) startInputBits() int      { return int(mm.mb.start_input_bits) }
func (mm *ModbusMapping) startInputRegisters() int { return int(mm.mb.start_input_registers) }
func (mm *ModbusMapping) startRegisters() int      { return int(mm.mb.start_registers) }

func (mm *ModbusMapping) tabBitsSlice() []byte {
	n := int(mm.mb.nb_bits)
	tab := unsafe.Slice(mm.mb.tab_bits, n)
	return cUint8ToBytes(tab)
}

func (mm *ModbusMapping) tabInputBitsSlice() []byte {
	n := int(mm.mb.nb_input_bits)
	tab := unsafe.Slice(mm.mb.tab_input_bits, n)
	return cUint8ToBytes(tab)
}

func (mm *ModbusMapping) tabRegistersSlice() []uint16 {
	n := int(mm.mb.nb_registers)
	tab := unsafe.Slice(mm.mb.tab_registers, n)
	return cUint16ToSlice(tab)
}

func (mm *ModbusMapping) tabInputRegistersSlice() []uint16 {
	n := int(mm.mb.nb_input_registers)
	tab := unsafe.Slice(mm.mb.tab_input_registers, n)
	return cUint16ToSlice(tab)
}

// Public accessors

func (mm *ModbusMapping) NbBits() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.nb_bits)
}

func (mm *ModbusMapping) StartBits() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.start_bits)
}

func (mm *ModbusMapping) NbInputBits() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.nb_input_bits)
}

func (mm *ModbusMapping) StartInputBits() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.start_input_bits)
}

func (mm *ModbusMapping) NbInputRegisters() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.nb_input_registers)
}

func (mm *ModbusMapping) StartInputRegisters() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.start_input_registers)
}

func (mm *ModbusMapping) NbRegisters() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.nb_registers)
}

func (mm *ModbusMapping) StartRegisters() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0
	}
	return int(mm.mb.start_registers)
}

// TabBits returns an iterator over all bits with their absolute addresses.
// The iteration uses a snapshot so it is safe against concurrent modifications.
func (mm *ModbusMapping) TabBits() iter.Seq2[int, byte] {
	mm.mu.Lock()
	if mm.mb == nil {
		mm.mu.Unlock()
		return func(yield func(int, byte) bool) {}
	}
	start := int(mm.mb.start_bits)
	snap := mm.tabBitsSlice()
	mm.mu.Unlock()
	return func(yield func(int, byte) bool) {
		for k, v := range snap {
			if !yield(start+k, v) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabBits(addr int) (byte, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0, fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_bits)
	if offset < 0 || offset >= int(mm.mb.nb_bits) {
		return 0, fmt.Errorf("modbus: bits address %d out of range [%d, %d)", addr, int(mm.mb.start_bits), int(mm.mb.start_bits)+int(mm.mb.nb_bits))
	}
	tab := unsafe.Slice(mm.mb.tab_bits, int(mm.mb.nb_bits))
	return byte(tab[offset]), nil
}

func (mm *ModbusMapping) SetTabBits(addr int, v byte) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_bits)
	if offset < 0 || offset >= int(mm.mb.nb_bits) {
		return fmt.Errorf("modbus: bits address %d out of range [%d, %d)", addr, int(mm.mb.start_bits), int(mm.mb.start_bits)+int(mm.mb.nb_bits))
	}
	tab := unsafe.Slice(mm.mb.tab_bits, int(mm.mb.nb_bits))
	tab[offset] = C.uint8_t(v)
	return nil
}

// TabInputBits returns an iterator over all input bits with their absolute addresses.
// The iteration uses a snapshot so it is safe against concurrent modifications.
func (mm *ModbusMapping) TabInputBits() iter.Seq2[int, byte] {
	mm.mu.Lock()
	if mm.mb == nil {
		mm.mu.Unlock()
		return func(yield func(int, byte) bool) {}
	}
	start := int(mm.mb.start_input_bits)
	snap := mm.tabInputBitsSlice()
	mm.mu.Unlock()
	return func(yield func(int, byte) bool) {
		for k, v := range snap {
			if !yield(start+k, v) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabInputBits(addr int) (byte, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0, fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_input_bits)
	if offset < 0 || offset >= int(mm.mb.nb_input_bits) {
		return 0, fmt.Errorf("modbus: input bits address %d out of range [%d, %d)", addr, int(mm.mb.start_input_bits), int(mm.mb.start_input_bits)+int(mm.mb.nb_input_bits))
	}
	tab := unsafe.Slice(mm.mb.tab_input_bits, int(mm.mb.nb_input_bits))
	return byte(tab[offset]), nil
}

func (mm *ModbusMapping) SetTabInputBits(addr int, v byte) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_input_bits)
	if offset < 0 || offset >= int(mm.mb.nb_input_bits) {
		return fmt.Errorf("modbus: input bits address %d out of range [%d, %d)", addr, int(mm.mb.start_input_bits), int(mm.mb.start_input_bits)+int(mm.mb.nb_input_bits))
	}
	tab := unsafe.Slice(mm.mb.tab_input_bits, int(mm.mb.nb_input_bits))
	tab[offset] = C.uint8_t(v)
	return nil
}

// TabInputRegisters returns an iterator over all input registers with their absolute addresses.
// The iteration uses a snapshot so it is safe against concurrent modifications.
func (mm *ModbusMapping) TabInputRegisters() iter.Seq2[int, uint16] {
	mm.mu.Lock()
	if mm.mb == nil {
		mm.mu.Unlock()
		return func(yield func(int, uint16) bool) {}
	}
	start := int(mm.mb.start_input_registers)
	snap := mm.tabInputRegistersSlice()
	mm.mu.Unlock()
	return func(yield func(int, uint16) bool) {
		for k, v := range snap {
			if !yield(start+k, v) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabInputRegisters(addr int) (uint16, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0, fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_input_registers)
	if offset < 0 || offset >= int(mm.mb.nb_input_registers) {
		return 0, fmt.Errorf("modbus: input registers address %d out of range [%d, %d)", addr, int(mm.mb.start_input_registers), int(mm.mb.start_input_registers)+int(mm.mb.nb_input_registers))
	}
	tab := unsafe.Slice(mm.mb.tab_input_registers, int(mm.mb.nb_input_registers))
	return uint16(tab[offset]), nil
}

func (mm *ModbusMapping) SetTabInputRegisters(addr int, v uint16) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_input_registers)
	if offset < 0 || offset >= int(mm.mb.nb_input_registers) {
		return fmt.Errorf("modbus: input registers address %d out of range [%d, %d)", addr, int(mm.mb.start_input_registers), int(mm.mb.start_input_registers)+int(mm.mb.nb_input_registers))
	}
	tab := unsafe.Slice(mm.mb.tab_input_registers, int(mm.mb.nb_input_registers))
	tab[offset] = C.uint16_t(v)
	return nil
}

// TabRegisters returns an iterator over all registers with their absolute addresses.
// The iteration uses a snapshot so it is safe against concurrent modifications.
func (mm *ModbusMapping) TabRegisters() iter.Seq2[int, uint16] {
	mm.mu.Lock()
	if mm.mb == nil {
		mm.mu.Unlock()
		return func(yield func(int, uint16) bool) {}
	}
	start := int(mm.mb.start_registers)
	snap := mm.tabRegistersSlice()
	mm.mu.Unlock()
	return func(yield func(int, uint16) bool) {
		for k, v := range snap {
			if !yield(start+k, v) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabRegisters(addr int) (uint16, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return 0, fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_registers)
	if offset < 0 || offset >= int(mm.mb.nb_registers) {
		return 0, fmt.Errorf("modbus: registers address %d out of range [%d, %d)", addr, int(mm.mb.start_registers), int(mm.mb.start_registers)+int(mm.mb.nb_registers))
	}
	tab := unsafe.Slice(mm.mb.tab_registers, int(mm.mb.nb_registers))
	return uint16(tab[offset]), nil
}

func (mm *ModbusMapping) SetTabRegisters(addr int, v uint16) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return fmt.Errorf("modbus: mapping is nil")
	}
	offset := addr - int(mm.mb.start_registers)
	if offset < 0 || offset >= int(mm.mb.nb_registers) {
		return fmt.Errorf("modbus: registers address %d out of range [%d, %d)", addr, int(mm.mb.start_registers), int(mm.mb.start_registers)+int(mm.mb.nb_registers))
	}
	tab := unsafe.Slice(mm.mb.tab_registers, int(mm.mb.nb_registers))
	tab[offset] = C.uint16_t(v)
	return nil
}

// Free modbus_mapping_free - free a modbus_mapping_t structure
//
// The function shall free the four arrays of modbus_mapping_t structure and finally the modbus_mapping_t itself
// referenced by mb_mapping.
// It is safe to call Free multiple times.
func (mm *ModbusMapping) Free() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.mb == nil {
		return
	}
	C.modbus_mapping_free(mm.mb)
	mm.mb = nil
	runtime.SetFinalizer(mm, nil)
}

// SendRawRequest modbus_send_raw_request - send a raw request
//
// The modbus_send_raw_request() function shall send a request via the socket of the context ctx. This function must be
// used for debugging purposes because you have to take care to make a valid request by hand. The function only adds to
// the message, the header or CRC of the selected backend, so raw_req must start and contain at least a slave/unit
// identifier and a function code. This function can be used to send request not handled by the library.
//
// The public header of libmodbus provides a list of supported Modbus functions codes, prefixed by MODBUS_FC_ (eg.
// MODBUS_FC_READ_HOLDING_REGISTERS), to help build of raw requests.
func (x *Modbus) SendRawRequest(raw []byte) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	req := bytesToCUint8(raw)
	code := C.modbus_send_raw_request(x.ctx, unsafe.SliceData(req), C.int(len(raw)))
	if code < 0 {
		return newCError()
	}
	return nil
}

// SendRawRequestTid sends a raw request with a specific transaction ID.
func (x *Modbus) SendRawRequestTid(raw []byte, tid int) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	req := bytesToCUint8(raw)
	code := C.modbus_send_raw_request_tid(x.ctx, unsafe.SliceData(req), C.int(len(raw)), C.int(tid))
	if code < 0 {
		return newCError()
	}
	return nil
}

// Receive modbus_receive - receive an indication request
//
// The modbus_receive() function shall receive an indication request from the socket of the context ctx. This function
// is used by a Modbus slave/server to receive and analyze indication request sent by the masters/clients.
//
// If you need to use another socket or file descriptor than the one defined in the context ctx, see the function
// modbus_set_socket.
func (x *Modbus) Receive() (req []byte, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	recv := make([]C.uint8_t, MODBUS_MAX_ADU_LENGTH)
	code := C.modbus_receive(x.ctx, unsafe.SliceData(recv))
	if code < 0 {
		return nil, newCError()
	}
	req = cUint8ToBytes(recv[:code:code])
	return req, nil
}

// ReceiveConfirmation modbus_receive_confirmation - receive a confirmation request
//
// The modbus_receive_confirmation() function shall receive a request via the socket of the context ctx. This function
// must be used for debugging purposes because the received response isn't checked against the initial request. This
// function can be used to receive request not handled by the library.
//
// The maximum size of the response depends on the used backend, in RTU the rsp array must be MODBUS_RTU_MAX_ADU_LENGTH
// bytes and in TCP it must be MODBUS_TCP_MAX_ADU_LENGTH bytes. If you want to write code compatible with both, you can
// use the constant MODBUS_MAX_ADU_LENGTH (maximum value of all libmodbus backends). Take care to allocate enough
// memory to store responses to avoid crashes of your server.
func (x *Modbus) ReceiveConfirmation() (rsp []byte, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return nil, err
	}
	recv := make([]C.uint8_t, MODBUS_MAX_ADU_LENGTH)
	code := C.modbus_receive_confirmation(x.ctx, unsafe.SliceData(recv))
	if code < 0 {
		return nil, newCError()
	}
	rsp = cUint8ToBytes(recv[:code:code])
	return rsp, nil
}

// Lock ordering: Modbus.mu must be acquired before ModbusMapping.mu.
// Reply modbus_reply - send a response to the received request
//
// The modbus_reply() function shall send a response to received request. The request req given in argument is
// analyzed, a response is then built and sent by using the information of the modbus context ctx.
//
// If the request indicates to read or write a value the operation will done in the modbus mapping mb_mapping according
// to the type of the manipulated data.
//
// If an error occurs, an exception response will be sent.
//
// This function is designed for Modbus servers.
func (x *Modbus) Reply(req []byte, mm *ModbusMapping) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	raw := bytesToCUint8(req)
	mm.mu.Lock()
	defer mm.mu.Unlock()
	code := C.modbus_reply(x.ctx, unsafe.SliceData(raw), C.int(len(req)), mm.mb)
	if code < 0 {
		return newCError()
	}
	return nil
}

// ReplyException modbus_reply_exception - send an exception response
//
// The modbus_reply_exception() function shall send an exception response based on the 'exception_code' in argument.
func (x *Modbus) ReplyException(req []byte, ecode uint) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	raw := bytesToCUint8(req)
	code := C.modbus_reply_exception(x.ctx, unsafe.SliceData(raw), C.uint(ecode))
	if code < 0 {
		return newCError()
	}
	return nil
}

// EnableQuirks modbus_enable_quirks - enable a list of quirks according to a mask
func (x *Modbus) EnableQuirks(quirksMask Quirks) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_enable_quirks(x.ctx, C.uint(quirksMask))
	if code < 0 {
		return newCError()
	}
	return nil
}

// DisableQuirks modbus_disable_quirks - disable a list of quirks according to a mask
func (x *Modbus) DisableQuirks(quirksMask Quirks) (err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.ensureCtx(); err != nil {
		return err
	}
	code := C.modbus_disable_quirks(x.ctx, C.uint(quirksMask))
	if code < 0 {
		return newCError()
	}
	return nil
}

// GetHighByte returns the high byte of a 16-bit or wider integer.
func GetHighByte[T int16 | uint16 | int32 | uint32 | int64 | uint64](data T) byte {
	return byte((uint64(data) >> 8) & 0xFF)
}

// GetLowByte returns the low byte of a 16-bit or wider integer.
func GetLowByte[T int16 | uint16 | int32 | uint32 | int64 | uint64](data T) byte {
	return byte(uint64(data) & 0xFF)
}

// GetInt64FromInt16 converts 4 int16 values to int64 (big-endian).
func GetInt64FromInt16(tab []int16) int64 {
	return int64(tab[0])<<48 | int64(tab[1])<<32 | int64(tab[2])<<16 | int64(tab[3])
}

// GetInt32FromInt16 converts 2 int16 values to int32 (big-endian).
func GetInt32FromInt16(tab []int16) int32 {
	return int32(tab[0])<<16 | int32(tab[1])
}

// GetInt16FromInt8 converts 2 int8 values to int16 (big-endian).
func GetInt16FromInt8(tab []int8) int16 {
	return int16(tab[0])<<8 | int16(tab[1])
}

// SetInt16ToInt8 converts an int16 to 2 int8 values (big-endian).
func SetInt16ToInt8(value int16) []int8 {
	return []int8{int8(value >> 8), int8(value)}
}

// SetInt32ToInt16 converts an int32 to 2 int16 values (big-endian).
func SetInt32ToInt16(value int32) []int16 {
	return []int16{int16(value >> 16), int16(value)}
}

// SetInt64ToInt16 converts an int64 to 4 int16 values (big-endian).
func SetInt64ToInt16(value int64) []int16 {
	return []int16{int16(value >> 48), int16(value >> 32), int16(value >> 16), int16(value)}
}

// SetBitsFromByte modbus_set_bits_from_byte - set many bits from a single byte value
//
// The modbus_set_bits_from_byte() function shall set many bits from a single byte. All 8 bits from the byte value will
// be written to dest array starting at index position.
func SetBitsFromByte(dest []byte, index int, value byte) error {
	if index < 0 || index > len(dest)-8 {
		return fmt.Errorf("modbus: set bits from byte requires dest length >= index+8, got %d", len(dest))
	}
	C.modbus_set_bits_from_byte((*C.uint8_t)(unsafe.SliceData(dest)), C.int(index), C.uint8_t(value))
	return nil
}

// SetBitsFromBytes modbus_set_bits_from_bytes - set many bits from an array of bytes
//
// The modbus_set_bits_from_bytes function shall set bits by reading an array of bytes. All the bits of the bytes read
// from the first position of the array tab_byte are written as bits in the dest array starting at position index.
func SetBitsFromBytes(dest []byte, index int, nb uint, tab []byte) error {
	if int(nb) > len(tab) {
		return fmt.Errorf("modbus: set bits from bytes requires tab length >= nb, got %d", len(tab))
	}
	if index < 0 || (int(nb) > 0 && index > len(dest)-int(nb)) {
		return fmt.Errorf("modbus: set bits from bytes requires dest length >= index+nb, got %d", len(dest))
	}
	C.modbus_set_bits_from_bytes((*C.uint8_t)(unsafe.SliceData(dest)), C.int(index), C.uint(nb), (*C.uint8_t)(unsafe.SliceData(tab)))
	return nil
}

// GetByteFromBits modbus_get_byte_from_bits - get the value from many bits
//
// The modbus_get_byte_from_bits() function shall extract a value from many bits. All nb_bits bits from src at position
// index will be read as a single value. To obtain a full byte, set nb_bits to 8.
func GetByteFromBits(src []byte, index int, nb uint) (byte, error) {
	if int(nb) > 8 {
		return 0, fmt.Errorf("modbus: get byte from bits requires nb <= 8, got %d", nb)
	}
	if len(src) < index+int(nb) {
		return 0, fmt.Errorf("modbus: get byte from bits requires src length >= index+nb, got %d", len(src))
	}
	return byte(C.modbus_get_byte_from_bits((*C.uint8_t)(unsafe.SliceData(src)), C.int(index), C.uint(nb))), nil
}

// DecodeFloat modbus_get_float - get a float value from 2 registers (DCBA byte order)
func DecodeFloat(src []uint16) (float32, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("modbus: decode float requires at least 2 registers, got %d", len(src))
	}
	return float32(C.modbus_get_float((*C.uint16_t)(unsafe.SliceData(src)))), nil
}

// DecodeFloatABCD modbus_get_float_abcd - get a float value from 2 registers in ABCD byte order
func DecodeFloatABCD(src []uint16) (float32, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("modbus: decode float requires at least 2 registers, got %d", len(src))
	}
	return float32(C.modbus_get_float_abcd((*C.uint16_t)(unsafe.SliceData(src)))), nil
}

// DecodeFloatDCBA modbus_get_float_dcba - get a float value from 2 registers in DCBA byte order
func DecodeFloatDCBA(src []uint16) (float32, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("modbus: decode float requires at least 2 registers, got %d", len(src))
	}
	return float32(C.modbus_get_float_dcba((*C.uint16_t)(unsafe.SliceData(src)))), nil
}

// DecodeFloatBADC modbus_get_float_badc - get a float value from 2 registers in BADC byte order
func DecodeFloatBADC(src []uint16) (float32, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("modbus: decode float requires at least 2 registers, got %d", len(src))
	}
	return float32(C.modbus_get_float_badc((*C.uint16_t)(unsafe.SliceData(src)))), nil
}

// DecodeFloatCDAB modbus_get_float_cdab - get a float value from 2 registers in CDAB byte order
func DecodeFloatCDAB(src []uint16) (float32, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("modbus: decode float requires at least 2 registers, got %d", len(src))
	}
	return float32(C.modbus_get_float_cdab((*C.uint16_t)(unsafe.SliceData(src)))), nil
}

// EncodeFloat modbus_set_float - set a float value to 2 registers (ABCD byte order)
func EncodeFloat(f float32, dest []uint16) error {
	if len(dest) < 2 {
		return fmt.Errorf("modbus: encode float requires at least 2 registers, got %d", len(dest))
	}
	C.modbus_set_float(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
	return nil
}

// EncodeFloatABCD modbus_set_float_abcd - set a float value in 2 registers using ABCD byte order
func EncodeFloatABCD(f float32, dest []uint16) error {
	if len(dest) < 2 {
		return fmt.Errorf("modbus: encode float requires at least 2 registers, got %d", len(dest))
	}
	C.modbus_set_float_abcd(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
	return nil
}

// EncodeFloatDCBA modbus_set_float_dcba - set a float value in 2 registers using DCBA byte order
func EncodeFloatDCBA(f float32, dest []uint16) error {
	if len(dest) < 2 {
		return fmt.Errorf("modbus: encode float requires at least 2 registers, got %d", len(dest))
	}
	C.modbus_set_float_dcba(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
	return nil
}

// EncodeFloatBADC modbus_set_float_badc - set a float value in 2 registers using BADC byte order
func EncodeFloatBADC(f float32, dest []uint16) error {
	if len(dest) < 2 {
		return fmt.Errorf("modbus: encode float requires at least 2 registers, got %d", len(dest))
	}
	C.modbus_set_float_badc(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
	return nil
}

// EncodeFloatCDAB modbus_set_float_cdab - set a float value in 2 registers using CDAB byte order
func EncodeFloatCDAB(f float32, dest []uint16) error {
	if len(dest) < 2 {
		return fmt.Errorf("modbus: encode float requires at least 2 registers, got %d", len(dest))
	}
	C.modbus_set_float_cdab(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
	return nil
}

// cUint8ToBytes converts a []C.uint8_t to a []byte.
func cUint8ToBytes(src []C.uint8_t) []byte {
	out := make([]byte, len(src))
	for i, v := range src {
		out[i] = byte(v)
	}
	return out
}

// cUint16ToSlice converts a []C.uint16_t to a []uint16.
func cUint16ToSlice(src []C.uint16_t) []uint16 {
	out := make([]uint16, len(src))
	for i, v := range src {
		out[i] = uint16(v)
	}
	return out
}

// bytesToCUint8 converts a []byte to a []C.uint8_t.
func bytesToCUint8(src []byte) []C.uint8_t {
	dest := make([]C.uint8_t, len(src))
	for i, v := range src {
		dest[i] = C.uint8_t(v)
	}
	return dest
}

// uint16ToCUint16 converts a []uint16 to a []C.uint16_t.
func uint16ToCUint16(src []uint16) []C.uint16_t {
	dest := make([]C.uint16_t, len(src))
	for i, v := range src {
		dest[i] = C.uint16_t(v)
	}
	return dest
}
