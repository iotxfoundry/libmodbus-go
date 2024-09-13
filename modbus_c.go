package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include <errno.h>

int get_errno()
{
	return errno;
}
*/
import "C"
