package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib -lmodbus -Wl,-rpath=/usr/local/lib
#include "modbus.h"
*/
import "C"
import (
	"iter"
	"time"
	"unsafe"
)

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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
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
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	timeout = time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond
	return
}

// ModbusSetByteTimeout modbus_set_byte_timeout - set timeout between bytes
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
func (x *Modbus) ModbusSetByteTimeout(timeout time.Duration) (err error) {
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_byte_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusGetByteTimeout modbus_get_byte_timeout - get timeout between bytes
//
// The modbus_get_byte_timeout() function shall store the timeout interval between two
// consecutive bytes of the same message in the to_sec and to_usec arguments.
func (x *Modbus) ModbusGetByteTimeout() (timeout time.Duration, err error) {
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_byte_timeout(x.ctx, &to_sec, &to_usec)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	timeout = time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond
	return
}

// ModbusSetIndicationTimeout modbus_set_indication_timeout - set timeout between indications
//
// The modbus_set_indication_timeout() function shall set the timeout interval used by a server
// to wait for a request from a client.
//
// The value of to_usec argument must be in the range 0 to 999999.
//
// If both to_sec and to_usec are zero, this timeout will not be used at all. In this case,
// the server will wait forever.
func (x *Modbus) ModbusSetIndicationTimeout(timeout time.Duration) (err error) {
	usec := timeout - time.Duration(timeout.Seconds())*time.Second
	code := C.modbus_set_indication_timeout(x.ctx, C.uint32_t(timeout.Seconds()), C.uint32_t(usec.Microseconds()))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusGetIndicationTimeout modbus_get_indication_timeout - get timeout used to wait for an indication
// (request received by a server).
//
// The modbus_get_indication_timeout() function shall store the timeout interval used to wait for an
// indication in the to_sec and to_usec arguments. Indication is the term used by the Modbus protocol
// to designate a request received by the server.
//
// The default value is zero, it means the server will wait forever.
func (x *Modbus) ModbusGetIndicationTimeout() (timeout time.Duration, err error) {
	to_sec := C.uint32_t(0)
	to_usec := C.uint32_t(0)
	code := C.modbus_get_indication_timeout(x.ctx, &to_sec, &to_usec)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	timeout = time.Duration(to_sec)*time.Second + time.Duration(to_usec)*time.Microsecond
	return
}

// ModbusGetHeaderLength modbus_get_header_length - retrieve the current header length
//
// The modbus_get_header_length() function shall retrieve the current header length from
// the backend. This function is convenient to manipulate a message and so it's limited to
// low-level operations.
func (x *Modbus) ModbusGetHeaderLength() (length int) {
	code := C.modbus_get_header_length(x.ctx)
	length = int(code)
	return
}

// ModbusFree modbus_free - free a libmodbus context
//
// The modbus_free() function shall free an allocated modbus_t structure.
func (x *Modbus) ModbusFree() {
	C.modbus_free(x.ctx)
}

// ModbusClose modbus_close - close a Modbus connection
//
// The modbus_close() function shall close the connection established with the backend set in the context.
func (x *Modbus) ModbusClose() {
	C.modbus_close(x.ctx)
}

// ModbusFlush modbus_flush - flush non-transmitted data
//
// The modbus_flush() function shall discard data received but not read to the socket or file descriptor associated
// to the context 'ctx'.
func (x *Modbus) ModbusFlush() (err error) {
	code := C.modbus_flush(x.ctx)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusSetDebug modbus_set_debug - set debug flag of the context
//
// The modbus_set_debug() function shall set the debug flag of the modbus_t context by using the argument flag.
// By default, the boolean flag is set to FALSE. When the flag value is set to TRUE, many verbose messages are
// displayed on stdout and stderr. For example, this flag is useful to display the bytes of the Modbus messages.
func (x *Modbus) ModbusSetDebug(flag bool) (err error) {
	f := C.FALSE
	if flag {
		f = C.TRUE
	}
	code := C.modbus_set_debug(x.ctx, C.int(f))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusReadBits modbus_read_bits - read many bits
//
// The modbus_read_bits() function shall read the status of the nb bits (coils) to the address addr of the remote
// device. The result of reading is stored in dest array as unsigned bytes (8 bits) set to TRUE or FALSE.
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint8_t)).
//
// The function uses the Modbus function code 0x01 (read coil status).
func (x *Modbus) ModbusReadBits(addr int, nb int) (out []byte, err error) {
	dest := make([]C.uint8_t, nb)
	code := C.modbus_read_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	for _, v := range dest {
		out = append(out, byte(v))
	}
	return
}

// ModbusReadInputBits modbus_read_input_bits - read many input bits
//
// The modbus_read_input_bits() function shall read the content of the nb input bits to the address addr of the remote
// device. The result of reading is stored in dest array as unsigned bytes (8 bits) set to TRUE or FALSE.
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint8_t)).
//
// The function uses the Modbus function code 0x02 (read input status).
func (x *Modbus) ModbusReadInputBits(addr int, nb int) (out []byte, err error) {
	dest := make([]C.uint8_t, nb)
	code := C.modbus_read_input_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	for _, v := range dest {
		out = append(out, byte(v))
	}
	return
}

// ModbusReadRegisters modbus_read_registers - read many registers
//
// The modbus_read_registers() function shall read the content of the nb holding registers to the address addr of the
// remote device. The result of reading is stored in dest array as word values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x03 (read holding registers).
func (x *Modbus) ModbusReadRegisters(addr int, nb int) (out []uint16, err error) {
	dest := make([]C.uint16_t, nb)
	code := C.modbus_read_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	for _, v := range dest {
		out = append(out, uint16(v))
	}
	return
}

// ModbusReadInputRegisters modbus_read_input_registers - read many input registers
//
// The modbus_read_input_registers() function shall read the content of the nb input registers to address addr of the
// remote device. The result of the reading is stored in dest array as word values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x04 (read input registers). The holding registers and input registers
// have different historical meaning, but nowadays it's more common to use holding registers only.
func (x *Modbus) ModbusReadInputRegisters(addr int, nb int) (out []uint16, err error) {
	dest := make([]C.uint16_t, nb)
	code := C.modbus_read_input_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	for _, v := range dest {
		out = append(out, uint16(v))
	}
	return
}

// ModbusWriteBit modbus_write_bit - write a single bit
//
// The modbus_write_bit() function shall write the status of status at the address addr of the remote device. The value
// must be set to TRUE or FALSE.
//
// The function uses the Modbus function code 0x05 (force single coil).
func (x *Modbus) ModbusWriteBit(addr int, status bool) (err error) {
	f := C.FALSE
	if status {
		f = C.TRUE
	}
	code := C.modbus_write_bit(x.ctx, C.int(addr), C.int(f))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusWriteRegister modbus_write_register - write a single register
//
// he modbus_write_register() function shall write the value of value holding registers at the address addr of the
// remote evice.
//
// he function uses the Modbus function code 0x06 (preset single register).
func (x *Modbus) ModbusWriteRegister(addr int, value uint16) (err error) {
	code := C.modbus_write_register(x.ctx, C.int(addr), C.uint16_t(value))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusWriteBits modbus_write_bits - write many bits
//
// The modbus_write_bits() function shall write the status of the nb bits (coils) from src at the address addr of the
// remote device. The src array must contains bytes set to TRUE or FALSE.
//
// The function uses the Modbus function code 0x0F (force multiple coils).
func (x *Modbus) ModbusWriteBits(addr int, data []bool) (err error) {
	nb := len(data)
	dest := make([]C.uint8_t, nb)
	for _, v := range data {
		f := C.FALSE
		if v {
			f = C.TRUE
		}
		dest = append(dest, C.uint8_t(f))
	}
	code := C.modbus_write_bits(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusWriteRegisters modbus_write_registers - write many registers
//
// The modbus_write_registers() function shall write the content of the nb holding registers from the array src at
// address addr of the remote device.
//
// The function uses the Modbus function code 0x10 (preset multiple registers).
func (x *Modbus) ModbusWriteRegisters(addr int, data []uint16) (err error) {
	nb := len(data)
	dest := make([]C.uint16_t, nb)
	for _, v := range data {
		dest = append(dest, C.uint16_t(v))
	}
	code := C.modbus_write_registers(x.ctx, C.int(addr), C.int(nb), unsafe.SliceData(dest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusMaskWriteRegister modbus_mask_write_register - mask a single register
//
// The modbus_mask_write_register() function shall modify the value of the holding register at the address 'addr' of
// the remote device using the algorithm:
//
// new value = (current value AND 'and') OR ('or' AND (NOT 'and'))
//
// The function uses the Modbus function code 0x16 (mask single register).
func (x *Modbus) ModbusMaskWriteRegister(addr int, andMask uint16, orMask uint16) (err error) {
	code := C.modbus_mask_write_register(x.ctx, C.int(addr), C.uint16_t(andMask), C.uint16_t(orMask))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusWriteAndReadRegisters modbus_write_and_read_registers - write and read many registers in a single transaction
//
// The modbus_write_and_read_registers() function shall write the content of the write_nb holding registers from the
// array 'src' to the address write_addr of the remote device then shall read the content of the read_nb holding
// registers to the address read_addr of the remote device. The result of reading is stored in dest array as word
// values (16 bits).
//
// You must take care to allocate enough memory to store the results in dest (at least nb * sizeof(uint16_t)).
//
// The function uses the Modbus function code 0x17 (write/read registers).
func (x *Modbus) ModbusWriteAndReadRegisters(addr int, writeAddr int, src []uint16, readAddr int, readNb int) (dest []uint16, err error) {
	writeNb := len(src)
	csrc := make([]C.uint16_t, writeNb)
	for _, v := range src {
		csrc = append(csrc, C.uint16_t(v))
	}
	cdest := make([]C.uint16_t, readNb)
	code := C.modbus_write_and_read_registers(x.ctx, C.int(writeAddr), C.int(writeNb), (*C.uint16_t)(unsafe.Pointer(&csrc[0])), C.int(readAddr), C.int(readNb), (*C.uint16_t)(unsafe.Pointer(&cdest[0])))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	for _, v := range cdest {
		dest = append(dest, uint16(v))
	}
	return
}

// ModbusReportSlaveId modbus_report_slave_id - returns a description of the controller
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
func (x *Modbus) ModbusReportSlaveId() (dest *ReportSlaveId, err error) {
	cdest := make([]C.uint8_t, MODBUS_MAX_PDU_LENGTH)
	code := C.modbus_report_slave_id(x.ctx, C.int(MODBUS_MAX_PDU_LENGTH), unsafe.SliceData(cdest))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	buff := []byte{}
	for i := range code {
		buff = append(buff, byte(cdest[i]))
	}
	dest = &ReportSlaveId{
		SlaveId:            buff[0],
		RunIndicatorStatus: buff[1],
		AdditionalData:     buff[2:],
	}
	return
}

// ModbusMappingNewStartAddress modbus_mapping_new_start_address - allocate four arrays of bits and
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
// The newly created mb_mapping will have the following arrays:
//
//   - tab_bits set to NULL
//   - tab_input_bits set to NULL
//   - tab_input_registers allocated to store 10 registers (uint16_t)
//   - tab_registers set to NULL.
//
// The clients can read the first register by using the address 340 in its request. On the server side, you should
// use the first index of the array to set the value at this client address:
//
//	mb_mapping->tab_registers[0] = 42;
//
// If it isn't necessary to allocate an array for a specific type of data, you can pass the zero value in argument,
// the associated pointer will be NULL.
//
// This function is convenient to handle requests in a Modbus server/slave.
func ModbusMappingNewStartAddress(
	startBits uint,
	nbBits uint,
	startInputBits uint,
	nbInputBits uint,
	startRegisters uint,
	nbRegisters uint,
	startInputRegisters uint,
	nbInputRegisters uint,
) *ModbusMapping {
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
		return nil
	}
	return &ModbusMapping{
		mb: mn,
	}
}

// ModbusMappingNew modbus_mapping_new - allocate four arrays of bits and registers
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
func ModbusMappingNew(nbBits int, nbInputBits int, nbRegisters int, nbInputRegisters int) *ModbusMapping {
	mn := C.modbus_mapping_new(C.int(nbBits), C.int(nbInputBits), C.int(nbRegisters), C.int(nbInputRegisters))
	if mn == nil {
		return nil
	}
	return &ModbusMapping{
		mb: mn,
	}
}

func (mm *ModbusMapping) NbBits() int {
	return int(mm.mb.nb_bits)
}

func (mm *ModbusMapping) StartBits() int {
	return int(mm.mb.start_bits)
}

func (mm *ModbusMapping) NbInputBits() int {
	return int(mm.mb.nb_input_bits)
}

func (mm *ModbusMapping) StartInputBits() int {
	return int(mm.mb.start_input_bits)
}

func (mm *ModbusMapping) NbInputRegisters() int {
	return int(mm.mb.nb_input_registers)
}

func (mm *ModbusMapping) StartInputRegisters() int {
	return int(mm.mb.start_input_registers)
}

func (mm *ModbusMapping) NbRegisters() int {
	return int(mm.mb.nb_registers)
}

func (mm *ModbusMapping) StartRegisters() int {
	return int(mm.mb.start_registers)
}

func (mm *ModbusMapping) TabBits() iter.Seq2[int, bool] {
	return func(yield func(int, bool) bool) {
		tab := unsafe.Slice(mm.mb.tab_bits, mm.NbBits())
		for k, v := range tab {
			if !yield(int(mm.StartBits()+k), v == C.TRUE) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabBits(addr int) bool {
	tab := unsafe.Slice(mm.mb.tab_bits, mm.NbBits())
	return tab[addr-mm.StartBits()] == C.TRUE
}

func (mm *ModbusMapping) SetTabBits(addr int, v bool) {
	tab := unsafe.Slice(mm.mb.tab_bits, mm.NbBits())
	f := C.FALSE
	if v {
		f = C.TRUE
	}
	tab[addr-mm.StartBits()] = C.uint8_t(f)
}

func (mm *ModbusMapping) TabInputBits() iter.Seq2[int, bool] {
	return func(yield func(int, bool) bool) {
		tab := unsafe.Slice(mm.mb.tab_input_bits, mm.NbInputBits())
		for k, v := range tab {
			if !yield(int(mm.StartInputBits()+k), v == C.TRUE) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabInputBits(addr int) bool {
	tab := unsafe.Slice(mm.mb.tab_input_bits, mm.NbInputBits())
	return tab[addr-mm.StartInputBits()] == C.TRUE
}

func (mm *ModbusMapping) SetTabInputBits(addr int, v bool) {
	tab := unsafe.Slice(mm.mb.tab_input_bits, mm.NbInputBits())
	f := C.FALSE
	if v {
		f = C.TRUE
	}
	tab[addr-mm.StartInputBits()] = C.uint8_t(f)
}

func (mm *ModbusMapping) TabInputRegisters() iter.Seq2[int, uint16] {
	return func(yield func(int, uint16) bool) {
		tab := unsafe.Slice(mm.mb.tab_input_registers, mm.NbInputRegisters())
		for k, v := range tab {
			if !yield(int(mm.StartInputRegisters()+k), uint16(v)) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabInputRegisters(addr int) uint16 {
	tab := unsafe.Slice(mm.mb.tab_input_registers, mm.NbInputRegisters())
	return uint16(tab[addr-mm.StartInputRegisters()])
}

func (mm *ModbusMapping) SetTabInputRegisters(addr int, v uint16) {
	tab := unsafe.Slice(mm.mb.tab_input_registers, mm.NbInputRegisters())
	tab[addr-mm.StartInputRegisters()] = C.uint16_t(v)
}

func (mm *ModbusMapping) TabRegisters() iter.Seq2[int, uint16] {
	return func(yield func(int, uint16) bool) {
		tab := unsafe.Slice(mm.mb.tab_registers, mm.NbRegisters())
		for k, v := range tab {
			if !yield(int(mm.StartRegisters()+k), uint16(v)) {
				return
			}
		}
	}
}

func (mm *ModbusMapping) GetTabRegisters(addr int) uint16 {
	tab := unsafe.Slice(mm.mb.tab_registers, mm.NbRegisters())
	return uint16(tab[addr-mm.StartRegisters()])
}

func (mm *ModbusMapping) SetTabRegisters(addr int, v uint16) {
	tab := unsafe.Slice(mm.mb.tab_registers, mm.NbRegisters())
	tab[addr-mm.StartRegisters()] = C.uint16_t(v)
}

// Free modbus_mapping_free - free a modbus_mapping_t structure
//
// The function shall free the four arrays of modbus_mapping_t structure and finally the modbus_mapping_t itself
// referenced by mb_mapping.
func (mm *ModbusMapping) Free() {
	C.modbus_mapping_free(mm.mb)
}

// ModbusSendRawRequest modbus_send_raw_request - send a raw request
//
// The modbus_send_raw_request() function shall send a request via the socket of the context ctx. This function must be
// used for debugging purposes because you have to take care to make a valid request by hand. The function only adds to
// the message, the header or CRC of the selected backend, so raw_req must start and contain at least a slave/unit
// identifier and a function code. This function can be used to send request not handled by the library.
//
// The public header of libmodbus provides a list of supported Modbus functions codes, prefixed by MODBUS_FC_ (eg.
// MODBUS_FC_READ_HOLDING_REGISTERS), to help build of raw requests.
func (x *Modbus) ModbusSendRawRequest(raw []byte) (err error) {
	req := []C.uint8_t{}
	for _, v := range raw {
		req = append(req, C.uint8_t(v))
	}
	code := C.modbus_send_raw_request(x.ctx, unsafe.SliceData(req), C.int(len(raw)))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

func (x *Modbus) ModbusSendRawRequestTid(raw []byte, tid int) (err error) {
	req := []C.uint8_t{}
	for _, v := range raw {
		req = append(req, C.uint8_t(v))
	}
	code := C.modbus_send_raw_request_tid(x.ctx, unsafe.SliceData(req), C.int(len(raw)), C.int(tid))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusReceive modbus_receive - receive an indication request
//
// The modbus_receive() function shall receive an indication request from the socket of the context ctx. This function
// is used by a Modbus slave/server to receive and analyze indication request sent by the masters/clients.
//
// If you need to use another socket or file descriptor than the one defined in the context ctx, see the function
// modbus_set_socket.
func (x *Modbus) ModbusReceive() (req []byte, err error) {
	var recv *C.uint8_t
	code := C.modbus_receive(x.ctx, recv)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	data := unsafe.Slice(recv, int(code))
	for _, v := range data {
		req = append(req, byte(v))
	}
	return
}

// ModbusReceiveConfirmation modbus_receive_confirmation - receive a confirmation request
//
// The modbus_receive_confirmation() function shall receive a request via the socket of the context ctx. This function
// must be used for debugging purposes because the received response isn't checked against the initial request. This
// function can be used to receive request not handled by the library.
//
// The maximum size of the response depends on the used backend, in RTU the rsp array must be MODBUS_RTU_MAX_ADU_LENGTH
// bytes and in TCP it must be MODBUS_TCP_MAX_ADU_LENGTH bytes. If you want to write code compatible with both, you can
// use the constant MODBUS_MAX_ADU_LENGTH (maximum value of all libmodbus backends). Take care to allocate enough
// memory to store responses to avoid crashes of your server.
func (x *Modbus) ModbusReceiveConfirmation() (rsp []byte, err error) {
	var recv *C.uint8_t
	code := C.modbus_receive_confirmation(x.ctx, recv)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	data := unsafe.Slice(recv, int(code))
	for _, v := range data {
		rsp = append(rsp, byte(v))
	}
	return
}

// ModbusReply modbus_reply - send a response to the received request
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
func (x *Modbus) ModbusReply(req []byte, mm *ModbusMapping) (err error) {
	raw := []C.uint8_t{}
	for _, v := range req {
		raw = append(raw, C.uint8_t(v))
	}
	code := C.modbus_reply(x.ctx, unsafe.SliceData(raw), C.int(len(req)), mm.mb)
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusReplyException modbus_reply_exception - send an exception response
//
// The modbus_reply_exception() function shall send an exception response based on the 'exception_code' in argument.
//
// The libmodbus provides the following exception codes:
//
//   - MODBUS_EXCEPTION_ILLEGAL_FUNCTION (1)
//   - MODBUS_EXCEPTION_ILLEGAL_DATA_ADDRESS (2)
//   - MODBUS_EXCEPTION_ILLEGAL_DATA_VALUE (3)
//   - MODBUS_EXCEPTION_SLAVE_OR_SERVER_FAILURE (4)
//   - MODBUS_EXCEPTION_ACKNOWLEDGE (5)
//   - MODBUS_EXCEPTION_SLAVE_OR_SERVER_BUSY (6)
//   - MODBUS_EXCEPTION_NEGATIVE_ACKNOWLEDGE (7)
//   - MODBUS_EXCEPTION_MEMORY_PARITY (8)
//   - MODBUS_EXCEPTION_NOT_DEFINED (9)
//   - MODBUS_EXCEPTION_GATEWAY_PATH (10)
//   - MODBUS_EXCEPTION_GATEWAY_TARGET (11)
//
// The initial request req is required to build a valid response.
func (x *Modbus) ModbusReplyException(req []byte, ecode uint) (err error) {
	raw := []C.uint8_t{}
	for _, v := range req {
		raw = append(raw, C.uint8_t(v))
	}
	code := C.modbus_reply_exception(x.ctx, unsafe.SliceData(raw), C.uint(ecode))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusEnableQuirks modbus_enable_quirks - enable a list of quirks according to a mask
//
// The function is only useful when you are confronted with equipment which does not respect the protocol, which
// behaves strangely or when you wish to move away from the standard.
//
// In that case, you can enable a specific quirk to workaround the issue, libmodbus offers the following flags:
//
//   - MODBUS_QUIRK_MAX_SLAVE allows slave adresses between 247 and 255.
//   - MODBUS_QUIRK_REPLY_TO_BROADCAST force a reply to a broacast request when the device is a slave in RTU mode
//     (should be enabled on the slave device).
//
// You can combine the flags by using the bitwise OR operator.
func (x *Modbus) ModbusEnableQuirks(quirksMask uint) (err error) {
	code := C.modbus_enable_quirks(x.ctx, C.uint(quirksMask))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

// ModbusDisableQuirks modbus_disable_quirks - disable a list of quirks according to a mask
//
// The function shall disable the quirks according to the provided mask. It's useful to revert changes applied by a
// previous call to modbus_enable_quirks
//
// To reset all quirks, you can use the specific value MODBUS_QUIRK_ALL.
//
//	modbus_enable_quirks(ctx, MODBUS_QUIRK_MAX_SLAVE | MODBUS_QUIRK_REPLY_TO_BROADCAST);
//
//	...
//
//	// Reset all quirks
//	modbus_disable_quirks(ctx, MODBUS_QUIRK_ALL);
func (x *Modbus) ModbusDisableQuirks(quirksMask uint) (err error) {
	code := C.modbus_disable_quirks(x.ctx, C.uint(quirksMask))
	if code < 0 {
		err = ErrorCode(code).Error()
		return
	}
	return
}

func ModbusGetHighByte[T int16 | uint16 | int32 | uint32 | int64 | uint64](data T) byte {
	return byte((uint64(data) >> 8) & 0xFF)
}

func ModbusGetLowByte[T int16 | uint16 | int32 | uint32 | int64 | uint64](data T) byte {
	return byte(uint64(data) & 0xFF)
}

func ModbusGetInt64FromInt16(tab []int16) int64 {
	return int64(tab[0])<<48 | int64(tab[1])<<32 | int64(tab[2])<<16 | int64(tab[3])
}

func ModbusGetInt32FromInt16(tab []int16) int32 {
	return int32(tab[0])<<16 | int32(tab[1])
}

func ModbusGetInt16FromInt8(tab []int8) int16 {
	return int16(tab[0])<<8 | int16(tab[1])
}

func ModbusSetInt16ToInt8(value int16) []int8 {
	return []int8{int8(value >> 8), int8(value)}
}

func ModbusSetInt32ToInt16(value int32) []int16 {
	return []int16{int16(value >> 16), int16(value)}
}

func ModbusSetInt64ToInt16(value int64) []int16 {
	return []int16{int16(value >> 48), int16(value >> 32), int16(value >> 16), int16(value)}
}

// ModbusSetBitsFromByte modbus_set_bits_from_byte - set many bits from a single byte value
//
// The modbus_set_bits_from_byte() function shall set many bits from a single byte. All 8 bits from the byte value will
// be written to dest array starting at index position.
func ModbusSetBitsFromByte(dest []byte, index int, value byte) {
	C.modbus_set_bits_from_byte((*C.uint8_t)(unsafe.SliceData(dest)), C.int(index), C.uint8_t(value))
}

// ModbusSetBitsFromBytes modbus_set_bits_from_bytes - set many bits from an array of bytes
//
// The modbus_set_bits_from_bytes function shall set bits by reading an array of bytes. All the bits of the bytes read
// from the first position of the array tab_byte are written as bits in the dest array starting at position index.
func ModbusSetBitsFromBytes(dest []byte, index int, nb uint, tab []byte) {
	C.modbus_set_bits_from_bytes((*C.uint8_t)(unsafe.SliceData(dest)), C.int(index), C.uint(nb), (*C.uint8_t)(unsafe.SliceData(tab)))
}

// ModbusGetByteFromBits modbus_get_byte_from_bits - get the value from many bits
//
// The modbus_get_byte_from_bits() function shall extract a value from many bits. All nb_bits bits from src at position
// index will be read as a single value. To obtain a full byte, set nb_bits to 8.
func ModbusGetByteFromBits(src []byte, index int, nb uint) byte {
	return byte(C.modbus_get_byte_from_bits((*C.uint8_t)(unsafe.SliceData(src)), C.int(index), C.uint(nb)))
}

// ModbusGetFloat modbus_get_float - get a float value from 2 registers
//
// The modbus_get_float() function shall get a float from 4 bytes in Modbus format (DCBA byte order). The src array
// must be a pointer on two 16 bits values, for example, if the first word is set to 0x4465 and the second to 0x229a,
// the float value will be 916.540649.
func ModbusGetFloat(src []uint16) float32 {
	return float32(C.modbus_get_float((*C.uint16_t)(unsafe.SliceData(src))))
}

// ModbusGetFloatAbcd modbus_get_float_abcd - get a float value from 2 registers in ABCD byte order
//
// The modbus_get_float_abcd() function shall get a float from 4 bytes in usual Modbus format. The src array must be a
// pointer on two 16 bits values, for example, if the first word is set to 0x0020 and the second to 0xF147, the float
// value will be read as 123456.0.
func ModbusGetFloatAbcd(src []uint16) float32 {
	return float32(C.modbus_get_float_abcd((*C.uint16_t)(unsafe.SliceData(src))))
}

// ModbusGetFloatDcba modbus_get_float_dcba - get a float value from 2 registers in DCBA byte order
//
// The modbus_get_float_dcba() function shall get a float from 4 bytes in inverted Modbus format (DCBA order instead of
// ABCD). The src array must be a pointer on two 16 bits values, for example, if the first word is set to 0x47F1 and
// the second to 0x2000, the float value will be read as 123456.0.
func ModbusGetFloatDcba(src []uint16) float32 {
	return float32(C.modbus_get_float_dcba((*C.uint16_t)(unsafe.SliceData(src))))
}

// ModbusGetFloatBadc modbus_get_float_badc - get a float value from 2 registers in BADC byte order
//
// The modbus_get_float_badc() function shall get a float from 4 bytes with swapped bytes (BADC instead of ABCD). The
// src array must be a pointer on two 16 bits values, for example, if the first word is set to 0x2000 and the second to
// 0x47F1, the float value will be read as 123456.0.
func ModbusGetFloatBadc(src []uint16) float32 {
	return float32(C.modbus_get_float_badc((*C.uint16_t)(unsafe.SliceData(src))))
}

// ModbusGetFloatCdab modbus_get_float_cdab - get a float value from 2 registers in CDAB byte order
//
// The modbus_get_float_cdab() function shall get a float from 4 bytes with swapped words (CDAB order instead of ABCD).
// The src array must be a pointer on two 16 bits values, for example, if the first word is set to F147 and the second
// to 0x0020, the float value will be read as 123456.0.
func ModbusGetFloatCdab(src []uint16) float32 {
	return float32(C.modbus_get_float_cdab((*C.uint16_t)(unsafe.SliceData(src))))
}

// ModbusSetFloat modbus_set_float - set a float value from 2 registers
//
// The modbus_set_float() function shall set a float to 4 bytes in Modbus format (ABCD). The dest array must be pointer
// on two 16 bits values to be able to store the full result of the conversion.
func ModbusSetFloat(f float32, dest []uint16) {
	C.modbus_set_float(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
}

// ModbusSetFloatAbcd modbus_set_float_abcd - set a float value in 2 registers using ABCD byte order
//
// The modbus_set_float_abcd() function shall set a float to 4 bytes in usual Modbus format. The dest array must be
// pointer on two 16 bits values to be able to store the full result of the conversion.
func ModbusSetFloatAbcd(f float32, dest []uint16) {
	C.modbus_set_float_abcd(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
}

// ModbusSetFloatDcba modbus_set_float_dcba - set a float value in 2 registers using DCBA byte order
//
// The modbus_set_float_dcba() function shall set a float to 4 bytes in inverted Modbus format (DCBA order). The dest
// array must be pointer on two 16 bits values to be able to store the full result of the conversion.
func ModbusSetFloatDcba(f float32, dest []uint16) {
	C.modbus_set_float_dcba(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
}

// ModbusSetFloatBadc modbus_set_float_badc - set a float value in 2 registers using BADC byte order
//
// The modbus_set_float_badc() function shall set a float to 4 bytes in swapped bytes Modbus format (BADC instead of
// ABCD). The dest array must be pointer on two 16 bits values to be able to store the full result of the conversion.
func ModbusSetFloatBadc(f float32, dest []uint16) {
	C.modbus_set_float_badc(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
}

// ModbusSetFloatCdab modbus_set_float_cdab - set a float value in 2 registers using CDAB byte order
//
// The modbus_set_float_cdab() function shall set a float to 4 bytes in swapped words Modbus format (CDAB order instead
// of ABCD). The dest array must be pointer on two 16 bits values to be able to store the full result of the conversion.
func ModbusSetFloatCdab(f float32, dest []uint16) {
	C.modbus_set_float_cdab(C.float(f), (*C.uint16_t)(unsafe.SliceData(dest)))
}
