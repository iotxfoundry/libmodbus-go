# libmodbus-go

[![Go Version](http://img.shields.io/github/go-mod/go-version/iotxfoundry/libmodbus-go)][gomod]
[![GitHub release](http://img.shields.io/github/release/iotxfoundry/libmodbus-go.svg?style=flat-square)][release]
[![LGPL-2.1 license](https://img.shields.io/badge/license-LGPL2.1-blue?style=flat-square)][license]
[![libmodbus version](https://img.shields.io/badge/libmodbus-v3.2.0-blue)][libmodbus]

[gomod]: https://github.com/iotxfoundry/libmodbus-go/blob/main/go.md
[release]: https://github.com/iotxfoundry/libmodbus-go/releases
[license]: https://github.com/iotxfoundry/libmodbus-go/blob/main/LICENSE
[libmodbus]: https://github.com/stephane/libmodbus/releases/tag/v3.2.0

golang binding for libmodbus

## functions

| C                                  | Go                              | comment |
| ---------------------------------- | ------------------------------- | ------- |
| modbus_close()                     | Modbus.Close()                  |         |
| modbus_connect()                   | Modbus.Connect()                |         |
| modbus_disable_quirks()            | Modbus.DisableQuirks()          |         |
| modbus_enable_quirks()             | Modbus.EnableQuirks()           |         |
| modbus_flush()                     | Modbus.Flush()                  |         |
| modbus_free()                      | Modbus.Free()                   |         |
| modbus_get_byte_from_bits()        | GetByteFromBits()               |         |
| modbus_get_byte_timeout()          | Modbus.GetByteTimeout()         |         |
| modbus_get_float()                 | DecodeFloat()                   |         |
| modbus_get_float_abcd()            | DecodeFloatABCD()               |         |
| modbus_get_float_badc()            | DecodeFloatBADC()               |         |
| modbus_get_float_cdab()            | DecodeFloatCDAB()               |         |
| modbus_get_float_dcba()            | DecodeFloatDCBA()               |         |
| modbus_get_header_length()         | Modbus.GetHeaderLength()        |         |
| modbus_get_indication_timeout()    | Modbus.GetIndicationTimeout()   |         |
| modbus_get_response_timeout()      | Modbus.GetResponseTimeout()     |         |
| modbus_get_slave()                 | Modbus.GetSlave()               |         |
| modbus_get_socket()                | Modbus.GetSocket()              |         |
| modbus_mapping_free()              | ModbusMapping.Free()            |         |
| modbus_mapping_new()               | NewMapping()                    |         |
| modbus_mapping_new_start_address() | NewMappingWithStart()           |         |
| modbus_mask_write_register()       | Modbus.MaskWriteRegister()      |         |
| modbus_new_rtu()                   | NewRTU()                        |         |
| modbus_new_tcp()                   | NewTCP()                        |         |
| modbus_new_tcp_pi()                | NewTCPPI()                      |         |
| modbus_proxy()                     | Modbus.Proxy()                  |         |
| modbus_read_bits()                 | Modbus.ReadBits()               |         |
| modbus_read_input_bits()           | Modbus.ReadInputBits()          |         |
| modbus_read_input_registers()      | Modbus.ReadInputRegisters()     |         |
| modbus_read_registers()            | Modbus.ReadRegisters()          |         |
| modbus_receive()                   | Modbus.Receive()                |         |
| modbus_receive_confirmation()      | Modbus.ReceiveConfirmation()    |         |
| modbus_reply()                     | Modbus.Reply()                  |         |
| modbus_reply_exception()           | Modbus.ReplyException()         |         |
| modbus_report_slave_id()           | Modbus.ReportSlaveID()          |         |
| modbus_rtu_get_rts()               | Modbus.RTUGetRTS()              |         |
| modbus_rtu_get_rts_delay()         | Modbus.RTUGetRTSDelay()         |         |
| modbus_rtu_get_serial_mode()       | Modbus.RTUGetSerialMode()       |         |
| modbus_rtu_set_custom_rts()        | Modbus.RTUSetCustomRTS()        |         |
| modbus_rtu_set_rts()               | Modbus.RTUSetRTS()              |         |
| modbus_rtu_set_rts_delay()         | Modbus.RTUSetRTSDelay()         |         |
| modbus_rtu_set_serial_mode()       | Modbus.RTUSetSerialMode()       |         |
| modbus_send_raw_request()          | Modbus.SendRawRequest()         |         |
| modbus_send_raw_request_tid()      | Modbus.SendRawRequestTID()      |         |
| modbus_set_bits_from_byte()        | SetBitsFromByte()               |         |
| modbus_set_bits_from_bytes()       | SetBitsFromBytes()              |         |
| modbus_set_byte_timeout()          | Modbus.SetByteTimeout()         |         |
| modbus_set_debug()                 | Modbus.SetDebug()               |         |
| modbus_set_error_recovery()        | Modbus.SetErrorRecovery()       |         |
| modbus_set_float()                 | EncodeFloat()                   |         |
| modbus_set_float_abcd()            | EncodeFloatABCD()               |         |
| modbus_set_float_badc()            | EncodeFloatBADC()               |         |
| modbus_set_float_cdab()            | EncodeFloatCDAB()               |         |
| modbus_set_float_dcba()            | EncodeFloatDCBA()               |         |
| modbus_set_indication_timeout()    | Modbus.SetIndicationTimeout()   |         |
| modbus_set_response_timeout()      | Modbus.SetResponseTimeout()     |         |
| modbus_set_slave()                 | Modbus.SetSlave()               |         |
| modbus_set_socket()                | Modbus.SetSocket()              |         |
| modbus_tcp_accept()                | Modbus.TCPAccept()              |         |
| modbus_tcp_listen()                | Modbus.TCPListen()              |         |
| modbus_tcp_pi_accept()             | Modbus.TCPPIAccept()            |         |
| modbus_tcp_pi_listen()             | Modbus.TCPPIListen()            |         |
| modbus_write_and_read_registers()  | Modbus.WriteAndReadRegisters()  |         |
| modbus_write_bit()                 | Modbus.WriteBit()               |         |
| modbus_write_bits()                | Modbus.WriteBits()              |         |
| modbus_write_register()            | Modbus.WriteRegister()          |         |
| modbus_write_registers()           | Modbus.WriteRegisters()         |         |
