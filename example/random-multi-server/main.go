package main

import (
	"errors"
	"log"
	"syscall"

	libmodbusgo "github.com/iotxfoundry/libmodbus-go"
)

func FD_SET(fd int, p *syscall.FdSet) {
	if fd < 0 || fd/64 >= len(p.Bits) {
		return
	}
	p.Bits[fd/64] |= 1 << (uint(fd) % 64)
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}

func FD_CLR(fd int, p *syscall.FdSet) {
	if fd < 0 || fd/64 >= len(p.Bits) {
		return
	}
	p.Bits[fd/64] &^= 1 << (uint(fd) % 64)
}

func FD_ISSET(fd int, p *syscall.FdSet) bool {
	if fd < 0 || fd/64 >= len(p.Bits) {
		return false
	}
	return p.Bits[fd/64]&(1<<(uint(fd)%64)) != 0
}

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

	socket, err := ctx.TcpListen(1)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("TcpListen")

	clients := [10]int{}
	readfds := syscall.FdSet{}
	for {
		FD_ZERO(&readfds)
		FD_SET(socket, &readfds)
		maxFd := socket

		for _, client := range clients {
			if client > 0 {
				FD_SET(client, &readfds)
				if client > maxFd {
					maxFd = client
				}
			}
		}

		activity, err := syscall.Select(maxFd+1, &readfds, nil, nil, nil)
		if activity < 0 {
			if err != nil && !errors.Is(err, syscall.EINTR) {
				log.Printf("select error: %s\n", err)
				break
			}
			continue
		}

		if FD_ISSET(socket, &readfds) {
			cs, err := ctx.TcpAccept()
			if err != nil {
				log.Printf("TcpAccept error: %s\n", err)
				continue
			}

			added := false
			for k, client := range clients {
				if client == 0 {
					clients[k] = cs
					added = true
					break
				}
			}
			if !added {
				log.Printf("Too many clients, closing new socket: %d\n", cs)
				syscall.Close(cs)
			}
		}

		for k, client := range clients {
			if client > 0 && FD_ISSET(client, &readfds) {
				ctx.SetSocket(client)
				req, err := ctx.Receive()
				if err != nil {
					syscall.Close(client)
					clients[k] = 0
					log.Printf("Client disconnected, socket fd: %d\n", client)
					continue
				}
				err = ctx.Reply(req, mbMapping)
				if err != nil {
					log.Printf("reply error: %s", err)
					break
				}
			}
		}
	}
}
