package libmodbusgo

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func ExampleModbus_Close() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()
}

func ExampleModbus_Free() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()
}

func ExampleModbusNewTcp() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()
}

func ExampleModbusNewTcpPi() {
	ctx := ModbusNewTcpPi("::1", "1502")
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()
}

func ExampleModbus_Connect() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()
}

func ExampleModbus_DisableQuirks() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.DisableQuirks(MODBUS_QUIRK_ALL)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbus_EnableQuirks() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.EnableQuirks(MODBUS_QUIRK_MAX_SLAVE | MODBUS_QUIRK_REPLY_TO_BROADCAST)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbus_GetByteTimeout() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	timeout, err := ctx.GetByteTimeout()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(timeout)
}

func ExampleModbus_GetIndicationTimeout() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	timeout, err := ctx.GetIndicationTimeout()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(timeout)
}

func ExampleModbus_GetResponseTimeout() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	timeout, err := ctx.GetResponseTimeout()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(timeout)
}

func ExampleModbus_SetResponseTimeout() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.SetResponseTimeout(0)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbusMappingNew() {
	BITS_ADDRESS := 0
	BITS_NB := 10
	INPUT_BITS_ADDRESS := 0
	INPUT_BITS_NB := 10
	REGISTERS_ADDRESS := 0
	REGISTERS_NB := 10
	INPUT_REGISTERS_ADDRESS := 0
	INPUT_REGISTERS_NB := 10
	mm := ModbusMappingNew(
		BITS_ADDRESS+BITS_NB,
		INPUT_BITS_ADDRESS+INPUT_BITS_NB,
		REGISTERS_ADDRESS+REGISTERS_NB,
		INPUT_REGISTERS_ADDRESS+INPUT_REGISTERS_NB,
	)
	defer mm.Free()
	fmt.Println(mm.NbBits())
	fmt.Println(mm.NbInputBits())
	fmt.Println(mm.NbInputRegisters())
	fmt.Println(mm.NbRegisters())

	// Output:
	// 10
	// 10
	// 10
	// 10
}

func ExampleModbusMappingNewStartAddress() {
	BITS_ADDRESS := 0
	BITS_NB := 10
	INPUT_BITS_ADDRESS := 0
	INPUT_BITS_NB := 10
	REGISTERS_ADDRESS := 0
	REGISTERS_NB := 10
	INPUT_REGISTERS_ADDRESS := 0
	INPUT_REGISTERS_NB := 10
	mm := ModbusMappingNewStartAddress(
		uint(BITS_ADDRESS),
		uint(BITS_NB),
		uint(INPUT_BITS_ADDRESS),
		uint(INPUT_BITS_NB),
		uint(REGISTERS_ADDRESS),
		uint(REGISTERS_NB),
		uint(INPUT_REGISTERS_ADDRESS),
		uint(INPUT_REGISTERS_NB),
	)
	defer mm.Free()
	fmt.Println(mm.NbBits())
	fmt.Println(mm.NbInputBits())
	fmt.Println(mm.NbInputRegisters())
	fmt.Println(mm.NbRegisters())

	// Output:
	// 10
	// 10
	// 10
	// 10
}

func ExampleModbusMapping_Free() {
	BITS_ADDRESS := 0
	BITS_NB := 10
	INPUT_BITS_ADDRESS := 0
	INPUT_BITS_NB := 10
	REGISTERS_ADDRESS := 0
	REGISTERS_NB := 10
	INPUT_REGISTERS_ADDRESS := 0
	INPUT_REGISTERS_NB := 10
	mm := ModbusMappingNewStartAddress(
		uint(BITS_ADDRESS),
		uint(BITS_NB),
		uint(INPUT_BITS_ADDRESS),
		uint(INPUT_BITS_NB),
		uint(REGISTERS_ADDRESS),
		uint(REGISTERS_NB),
		uint(INPUT_REGISTERS_ADDRESS),
		uint(INPUT_REGISTERS_NB),
	)
	defer mm.Free()
	fmt.Println(mm.NbBits())
	fmt.Println(mm.NbInputBits())
	fmt.Println(mm.NbInputRegisters())
	fmt.Println(mm.NbRegisters())

	// Output:
	// 10
	// 10
	// 10
	// 10
}

func ExampleModbusNewRtu() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.SetSlave(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := ctx.ReadRegisters(0, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

func ExampleModbus_SetSlave() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.SetSlave(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := ctx.ReadRegisters(0, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

func ExampleModbus_ReadRegisters() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	err = ctx.SetSlave(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := ctx.ReadRegisters(0, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

func ExampleModbus_ReceiveConfirmation() {
	var ctx *Modbus // just for test

	rsp, err := ctx.ReceiveConfirmation()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func ExampleModbus_ReportSlaveId() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	report, err := ctx.ReportSlaveId()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(report.SlaveId)
	fmt.Println(report.RunIndicatorStatus)
}

func ExampleModbus_RtuSetSerialMode() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	ctx.SetSlave(1)
	ctx.RtuSetSerialMode(MODBUS_RTU_RS485)
	ctx.RtuSetRts(MODBUS_RTU_RTS_UP)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	report, err := ctx.ReportSlaveId()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(report.SlaveId)
	fmt.Println(report.RunIndicatorStatus)
}

func ExampleModbus_RtuSetRts() {
	ctx := ModbusNewRtu("/dev/ttyUSB0", 115200, 'N', 8, 1)
	if ctx == nil {
		fmt.Println("Unable to create the libmodbus context")
		return
	}
	defer ctx.Free()
	ctx.SetSlave(1)
	ctx.RtuSetSerialMode(MODBUS_RTU_RS485)
	ctx.RtuSetRts(MODBUS_RTU_RTS_UP)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	report, err := ctx.ReportSlaveId()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(report.SlaveId)
	fmt.Println(report.RunIndicatorStatus)
}

func ExampleModbus_SendRawRequest() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

	req := []byte{0xFF, MODBUS_FC_READ_HOLDING_REGISTERS, 0x00, 0x01, 0x0, 0x05}

	err = ctx.SendRawRequest(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp, err := ctx.ReceiveConfirmation()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func ExampleModbus_SetErrorRecovery() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()
	err := ctx.SetErrorRecovery(MODBUS_ERROR_RECOVERY_LINK | MODBUS_ERROR_RECOVERY_PROTOCOL)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Close()

}

func ExampleModbus_SetSocket() {
	const NB_CONNECTION = 100
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()

	ss, err := ctx.TcpListen(NB_CONNECTION)
	if err != nil {
		fmt.Println(err)
		return
	}

	var mm *ModbusMapping // TODO:

	var ms int
	rs := &unix.FdSet{}
	rs.Zero()
	rs.Set(ss)

	// ... TODO:

	if rs.IsSet(ms) {
		ctx.SetSocket(ms)
		req, err := ctx.Receive()
		if err == nil {
			ctx.Reply(req, mm)
		}
	}
}

func ExampleModbus_TcpAccept() {
	const NB_CONNECTION = 100
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()

	_, err := ctx.TcpListen(NB_CONNECTION)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ctx.TcpAccept()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbus_TcpListen() {
	const NB_CONNECTION = 100
	ctx := ModbusNewTcp("127.0.0.1", 502)
	if ctx == nil {
		return
	}
	defer ctx.Free()

	_, err := ctx.TcpListen(NB_CONNECTION)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ctx.TcpAccept()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbus_TcpPiListen() {
	const NB_CONNECTION = 100
	ctx := ModbusNewTcpPi("::0", "502")
	if ctx == nil {
		return
	}
	defer ctx.Free()

	_, err := ctx.TcpPiListen(NB_CONNECTION)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ctx.TcpPiAccept()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExampleModbus_TcpPiAccept() {
	const NB_CONNECTION = 100
	ctx := ModbusNewTcpPi("::0", "502")
	if ctx == nil {
		return
	}
	defer ctx.Free()

	_, err := ctx.TcpPiListen(NB_CONNECTION)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ctx.TcpPiAccept()
	if err != nil {
		fmt.Println(err)
		return
	}
}
