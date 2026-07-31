//go:build !windows

package modbus

import "syscall"

// shutdownSocket wakes up any receive/send call blocked on the socket fd.
// It is best-effort: errors are ignored by callers.
func shutdownSocket(fd int) error {
	return syscall.Shutdown(fd, syscall.SHUT_RDWR)
}
