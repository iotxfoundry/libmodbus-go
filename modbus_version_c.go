// Package libmodbus Wrappers for Go callback functions to be passed into C.
package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"

int modbus_version_check(unsigned int major, unsigned int minor, unsigned int micro)
{
	return LIBMODBUS_VERSION_CHECK(major, minor, micro);
}
*/
import "C"
