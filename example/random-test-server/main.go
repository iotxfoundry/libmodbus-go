package main

import (
	"log"

	"github.com/iotxfoundry/libmodbus-go"
)

func main() {
	ctx, err := modbus.NewTCP("127.0.0.1", 1502)
	if err != nil {
		log.Println("NewTCP error:", err)
		return
	}
	defer ctx.Free()
	defer func() { _ = ctx.Close() }()

	ctx.SetDebug(true)

	mbMapping, err := modbus.NewMapping(500, 500, 500, 500)
	if err != nil {
		log.Println("NewMapping error:", err)
		return
	}
	defer mbMapping.Free()

	_, err = ctx.TCPListen(1)
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = ctx.TCPAccept()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		req, err := ctx.Receive()
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
