package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iotxfoundry/libmodbus-go"
)

func main() {
	mapping, err := modbus.NewMapping(500, 500, 500, 500)
	if err != nil {
		log.Fatalln("NewMapping error:", err)
	}
	defer mapping.Free()

	server := modbus.NewTCPServer("127.0.0.1:1502", mapping)
	server.SetMaxClients(10)
	server.SetDebug(true)
	server.SetOnConnect(func(addr net.Addr) {
		log.Printf("Client connected: %s", addr)
	})
	server.SetOnDisconnect(func(addr net.Addr, err error) {
		log.Printf("Client disconnected: %s (%v)", addr, err)
	})

	// Run server in background
	go func() {
		log.Println("Starting Modbus TCP server on 127.0.0.1:1502")
		if err := server.Serve(); err != nil {
			log.Fatalln("Server error:", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Println("Close error:", err)
	}
	log.Println("Server stopped")
}
