package libmodbusgo

import "fmt"

func ExampleModbus_Close() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Free()
	defer ctx.Close()

	// Output:
}

func ExampleModbus_Free() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Free()
	defer ctx.Close()
}

func ExampleModbusNewTcp() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Free()
	defer ctx.Close()
}

func ExampleModbus_Connect() {
	ctx := ModbusNewTcp("127.0.0.1", 502)
	err := ctx.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ctx.Free()
	defer ctx.Close()
}
