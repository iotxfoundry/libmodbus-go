package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"
*/
import "C"

const (
	LIBMODBUS_VERSION_MAJOR  = C.LIBMODBUS_VERSION_MAJOR  // The major version, (1, if %LIBMODBUS_VERSION is 1.2.3)
	LIBMODBUS_VERSION_MINOR  = C.LIBMODBUS_VERSION_MINOR  // The minor version (2, if %LIBMODBUS_VERSION is 1.2.3)
	LIBMODBUS_VERSION_MICRO  = C.LIBMODBUS_VERSION_MICRO  // The micro version (3, if %LIBMODBUS_VERSION is 1.2.3)
	LIBMODBUS_VERSION_STRING = C.LIBMODBUS_VERSION_STRING // The full version, in string form (suited for string concatenation)
	LIBMODBUS_VERSION_HEX    = C.LIBMODBUS_VERSION_HEX    // Numerically encoded version, eg. v1.2.3 is 0x010203
)
