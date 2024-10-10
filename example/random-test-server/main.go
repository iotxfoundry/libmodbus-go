package main

import (
	"log"

	libmodbusgo "github.com/iotxfoundry/libmodbus-go"
)

func main() {
	ctx := libmodbusgo.ModbusNewTcp("127.0.0.1", 1502)
	if ctx == nil {
		log.Println("ModbusNewTcp error")
		return
	}
	defer ctx.Free()
	defer ctx.Close()

	ctx.SetDebug(true)

	mbMapping := libmodbusgo.ModbusMappingNew(500, 500, 500, 500)
	if mbMapping == nil {
		log.Println("ModbusMappingNew error")
		return
	}
	defer mbMapping.Free()

	_, err := ctx.TcpListen(1)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = ctx.TcpAccept()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		req, err := ctx.TcpReceive()
		if err != nil {
			log.Printf("receive error: %s", err)
			break
		}
		err = ctx.Reply(req, mbMapping)
		if err != nil {
			log.Printf("reply error: %s", err)
			break
		}
		// for k, v := range mbMapping.TabBits() {
		// 	log.Println(k, v)
		// }
	}
}
