package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include <errno.h>

int get_errno_cgo()
{
	return errno;
}
*/
import "C"
