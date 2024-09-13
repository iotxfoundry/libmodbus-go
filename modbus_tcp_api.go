package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib/libmodbus.a
#include <stdlib.h>
#include "modbus.h"
*/
import "C"
import (
	"unsafe"
)

// ModbusNewTcp modbus_new_tcp - create a libmodbus context for TCP/IPv4
//
// The modbus_new_tcp() function shall allocate and initialize a modbus_t structure to communicate with a Modbus TCP
// IPv4 server.
//
// The ip argument specifies the IP address of the server to which the client wants to establish a connection. A NULL
// value can be used to listen any addresses in server mode.
//
// The port argument is the TCP port to use. Set the port to MODBUS_TCP_DEFAULT_PORT to use the default one (502). It's
// convenient to use a port number greater than or equal to 1024 because it's not necessary to have administrator
// privileges.
func ModbusNewTcp(addr string, port int) *Modbus {
	saddr := C.CString(addr)
	defer C.free(unsafe.Pointer(saddr))
	ctx := C.modbus_new_tcp(saddr, C.int(port))
	if ctx == nil {
		return nil
	}
	return &Modbus{ctx: ctx}
}

// TcpListen modbus_tcp_listen - create and listen a TCP Modbus socket (IPv4)
//
// The modbus_tcp_listen() function shall create a socket and listen to maximum nb_connection incoming connections on
// the specified IP address. The context ctx must be allocated and initialized with modbus_new_tcp before to set the IP
// address to listen, if IP address is set to NULL or '0.0.0.0', any addresses will be listen.
func (x *Modbus) TcpListen(nb int) (socket int, err error) {
	code := C.modbus_tcp_listen(x.ctx, C.int(nb))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	socket = int(code)
	x.socket = socket
	return
}

// TcpAccept modbus_tcp_accept - accept a new connection on a TCP Modbus socket (IPv4)
//
// The modbus_tcp_accept() function shall extract the first connection on the queue of pending connections, create a
// new socket and store it in libmodbus context given in argument. If available, accept4() with SOCK_CLOEXEC will be
// called instead of accept().
func (x *Modbus) TcpAccept() (err error) {
	code := C.modbus_tcp_accept(x.ctx, (*C.int)(unsafe.Pointer(&x.socket)))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	return
}

// ModbusNewTcpPi modbus_new_tcp_pi - create a libmodbus context for TCP Protocol Independent
//
// The modbus_new_tcp_pi() function shall allocate and initialize a modbus_t structure to communicate with a Modbus TCP
// IPv4 or IPv6 server.
//
// The node argument specifies the host name or IP address of the host to connect to, eg. "192.168.0.5" , "::1" or
// "server.com". A NULL value can be used to listen any addresses in server mode.
//
// The service argument is the service name/port number to connect to. To use the default Modbus port, you can provide
// an NULL value or the string "502". On many Unix systems, it's convenient to use a port number greater than or equal
// to 1024 because it's not necessary to have administrator privileges.
//
// v3.1.8 handles NULL value for service (no EINVAL error).
func ModbusNewTcpPi(node string, service string) *Modbus {
	cnode := C.CString(node)
	defer C.free(unsafe.Pointer(cnode))
	cservice := C.CString(service)
	defer C.free(unsafe.Pointer(cservice))
	ctx := C.modbus_new_tcp_pi(cnode, cservice)
	if ctx == nil {
		return nil
	}
	return &Modbus{ctx: ctx}
}

// TcpPiListen modbus_tcp_pi_listen - create and listen a TCP PI Modbus socket (IPv6)
//
// The modbus_tcp_pi_listen() function shall create a socket and listen to maximum nb_connection incoming connections
// on the specified nodes. The context ctx must be allocated and initialized with modbus_new_tcp_pi before to set the
// node to listen, if node is set to NULL or '0.0.0.0', any addresses will be listen.
func (x *Modbus) TcpPiListen(nb int) (socket int, err error) {
	code := C.modbus_tcp_pi_listen(x.ctx, C.int(nb))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	socket = int(code)
	x.socket = socket
	return
}

// TcpPiAccept modbus_tcp_pi_accept - accept a new connection on a TCP PI Modbus socket (IPv6)
//
// The modbus_tcp_pi_accept() function shall extract the first connection on the queue of pending connections, create a
// new socket and store it in libmodbus context given in argument. If available, accept4() with SOCK_CLOEXEC will be
// called instead of accept().
func (x *Modbus) TcpPiAccept() (err error) {
	code := C.modbus_tcp_pi_accept(x.ctx, (*C.int)(unsafe.Pointer(&x.socket)))
	if code < 0 {
		err = ModbusStrError()
		return
	}
	return
}
