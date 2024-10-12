package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo linux,amd64 LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_amd64/modbus/lib/libmodbus.a
#cgo linux,arm LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_armv7/modbus/lib/libmodbus.a
#cgo linux,arm64 LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_armv8/modbus/lib/libmodbus.a
#cgo linux,386 LDFLAGS: -static -L${SRCDIR}/3rdParty/linux_386/modbus/lib/libmodbus.a
#cgo windows LDFLAGS: -static -L${SRCDIR}/3rdParty/windows_amd64/modbus/lib/modbus.lib

#include "modbus.h"

extern int modbus_version_check(unsigned int major, unsigned int minor, unsigned int micro);
*/
import "C"

// ModbusVersionCheck Evaluates to True if the version is greater than @major, @minor and @micro
func ModbusVersionCheck(major uint, minor uint, micro uint) bool {
	ret := C.modbus_version_check(C.uint(major), C.uint(minor), C.uint(micro))
	return ret == 1
}
