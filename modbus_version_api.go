package libmodbusgo

/*
#include "modbus.h"

extern int modbus_version_check(unsigned int major, unsigned int minor, unsigned int micro);
*/
import "C"

// ModbusVersionCheck Evaluates to True if the version is greater than @major, @minor and @micro
func ModbusVersionCheck(major uint, minor uint, micro uint) bool {
	ret := C.modbus_version_check(C.uint(major), C.uint(minor), C.uint(micro))
	return ret == 1
}
