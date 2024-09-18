// libmodbusgo golang binding for libmodbus
//
// # libmodbus functions convert
//
//	| C                                  | GO                              | comment |
//	| ---------------------------------- | ------------------------------- | ------- |
//	| modbus_close()                     | Modbus.Close()                  |         |
//	| modbus_connect()                   | Modbus.Connect()                |         |
//	| modbus_disable_quirks()            | Modbus.DisableQuirks()          |         |
//	| modbus_enable_quirks()             | Modbus.EnableQuirks()           |         |
//	| modbus_flush()                     | Modbus.Flush()                  |         |
//	| modbus_free()                      | Modbus.Free()                   |         |
//	| modbus_get_byte_from_bits()        | ModbusGetByteFromBits()         |         |
//	| modbus_get_byte_timeout()          | Modbus.GetByteTimeout()         |         |
//	| modbus_get_float()                 | ModbusGetFloat()                |         |
//	| modbus_get_float_abcd()            | ModbusGetFloatAbcd()            |         |
//	| modbus_get_float_badc()            | ModbusGetFloatBadc()            |         |
//	| modbus_get_float_cdab()            | ModbusGetFloatCdab()            |         |
//	| modbus_get_float_dcba()            | ModbusGetFloatDcba()            |         |
//	| modbus_get_header_length()         | Modbus.GetHeaderLength()        |         |
//	| modbus_get_indication_timeout()    | Modbus.GetIndicationTimeout()   |         |
//	| modbus_get_response_timeout()      | Modbus.GetResponseTimeout()     |         |
//	| modbus_get_slave()                 | Modbus.GetSlave()               |         |
//	| modbus_get_socket()                | Modbus.GetSocket()              |         |
//	| modbus_mapping_free()              | ModbusMapping.Free()            |         |
//	| modbus_mapping_new()               | ModbusMapping.New()             |         |
//	| modbus_mapping_new_start_address() | ModbusMapping.NewStartAddress() |         |
//	| modbus_mask_write_register()       | Modbus.MaskWriteRegister()      |         |
//	| modbus_new_rtu()                   | ModbusNewRtu()                  |         |
//	| modbus_new_tcp()                   | ModbusNewTcp()                  |         |
//	| modbus_new_tcp_pi()                | ModbusNewTcpPi()                |         |
//	| modbus_read_bits()                 | Modbus.ReadBits()               |         |
//	| modbus_read_input_bits()           | Modbus.ReadInputBits()          |         |
//	| modbus_read_input_registers()      | Modbus.ReadInputRegisters()     |         |
//	| modbus_read_registers()            | Modbus.ReadRegisters()          |         |
//	| modbus_receive()                   | Modbus.Receive()                |         |
//	| modbus_receive_confirmation()      | Modbus.ReceiveConfirmation()    |         |
//	| modbus_reply()                     | Modbus.Reply()                  |         |
//	| modbus_reply_exception()           | Modbus.ReplyException()         |         |
//	| modbus_report_slave_id()           | Modbus.ReportSlaveId()          |         |
//	| modbus_rtu_get_rts()               | Modbus.RtuGetRts()              |         |
//	| modbus_rtu_get_rts_delay()         | Modbus.RtuGetRtsDelay()         |         |
//	| modbus_rtu_get_serial_mode()       | Modbus.RtuGetSerialMode()       |         |
//	| modbus_rtu_set_custom_rts()        | Modbus.RtuSetCustomRts()        |         |
//	| modbus_rtu_set_rts()               | Modbus.RtuSetRts()              |         |
//	| modbus_rtu_set_rts_delay()         | Modbus.RtuSetRtsDelay()         |         |
//	| modbus_rtu_set_serial_mode()       | Modbus.RtuSetSerialMode()       |         |
//	| modbus_send_raw_request()          | Modbus.SendRawRequest()         |         |
//	| modbus_set_bits_from_byte()        | ModbusSetBitsFromByte()         |         |
//	| modbus_set_bits_from_bytes()       | ModbusSetBitsFromBytes()        |         |
//	| modbus_set_byte_timeout()          | Modbus.SetByteTimeout()         |         |
//	| modbus_set_debug()                 | Modbus.SetDebug()               |         |
//	| modbus_set_error_recovery()        | Modbus.SetErrorRecovery()       |         |
//	| modbus_set_float()                 | ModbusSetFloat()                |         |
//	| modbus_set_float_abcd()            | ModbusSetFloatAbcd()            |         |
//	| modbus_set_float_badc()            | ModbusSetFloatBadc()            |         |
//	| modbus_set_float_cdab()            | ModbusSetFloatCdab()            |         |
//	| modbus_set_float_dcba()            | ModbusSetFloatDcba()            |         |
//	| modbus_set_indication_timeout()    | Modbus.SetIndicationTimeout()   |         |
//	| modbus_set_response_timeout()      | Modbus.SetResponseTimeout()     |         |
//	| modbus_set_slave()                 | Modbus.SetSlave()               |         |
//	| modbus_set_socket()                | Modbus.SetSocket()              |         |
//	| modbus_strerror()                  | ModbusStrError()                |         |
//	| modbus_tcp_accept()                | Modbus.TcpAccept()              |         |
//	| modbus_tcp_listen()                | Modbus.TcpListen()              |         |
//	| modbus_tcp_pi_accept()             | Modbus.TcpPiAccept()            |         |
//	| modbus_tcp_pi_listen()             | Modbus.TcpPiListen()            |         |
//	| modbus_write_and_read_registers()  | Modbus.WriteAndReadRegisters()  |         |
//	| modbus_write_bit()                 | Modbus.WriteBit()               |         |
//	| modbus_write_bits()                | Modbus.WriteBits()              |         |
//	| modbus_write_register()            | Modbus.WriteRegister()          |         |
//	| modbus_write_registers()           | Modbus.WriteRegisters()         |         |
//
// # cgo, C, go type convert
//
//	|C Language Type        | CGO Type      |Go Language Type|     SDK        |
//	|---------------------- | ------------- | -------------- | -------------- |
//	|char                   | C.char        | byte           |                |
//	|singed char            | C.schar       | int8           |                |
//	|unsigned char          | C.uchar       | uint8          | BYTE           |
//	|short                  | C.short       | int16          |                |
//	|unsigned short         | C.ushort      | uint16         | WORD,USHORT    |
//	|int                    | C.int         | int32          | BOOL,LONG      |
//	|unsigned int           | C.uint        | uint32         | DWORD,UINT     |
//	|long                   | C.long        | int32          |                |
//	|unsigned long          | C.ulong       | uint32         |                |
//	|long long int          | C.longlong    | int64          |                |
//	|unsigned long long int | C.ulonglong   | uint64         |                |
//	|float                  | C.float       | float32        |                |
//	|double                 | C.double      | float64        |                |
//	|size_t                 | C.size_t      | uint           |                |
//	|void*                  | unsafe.Pointer| unsafe.Pointer | LPVOID.HANDLE  |
//
// # cgo, C, go string convert
//
//   - func C.CString(string) *C.char
//
//     Go string to C string.
//     The C string is allocated in the C heap using malloc.
//     It is the caller's responsibility to arrange for it to be
//     freed, such as by calling C.free (be sure to include stdlib.h
//     if C.free is needed).
//
//   - func C.CBytes([]byte) unsafe.Pointer
//
//     Go []byte slice to C array.
//     The C array is allocated in the C heap using malloc.
//     It is the caller's responsibility to arrange for it to be
//     freed, such as by calling C.free (be sure to include stdlib.h
//     if C.free is needed).
//
//   - func C.GoString(*C.char) string
//
//     C string to Go string
//
//   - func C.GoStringN(*C.char, C.int) string
//
//     C data with explicit length to Go string
//
//   - func C.GoBytes(unsafe.Pointer, C.int) []byte
//
//     C data with explicit length to Go []byte
//
// # C header -> CGO header
//
//  1. delete __stdcall
//
//  2. delete CALLBACK
//
//  3. enum XXX{}; -> type enum _XXX {}XXX;
//
//  4. delete function with init parameters, pUser = NULL
//
//  5. struct must have tag name
package libmodbusgo
