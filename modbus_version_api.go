package modbus

/*
#include "modbus.h"

extern int modbus_version_check(unsigned int major, unsigned int minor, unsigned int micro);
*/
import "C"

// VersionCheck evaluates to True if the version is greater than or equal to @major, @minor and @micro
func VersionCheck(major uint, minor uint, micro uint) bool {
	ret := C.modbus_version_check(C.uint(major), C.uint(minor), C.uint(micro))
	return ret == 1
}
