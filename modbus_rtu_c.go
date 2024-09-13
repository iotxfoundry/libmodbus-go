package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"

extern void set_rts_go(modbus_t *ctx, int on);

void set_rts_cgo(modbus_t *ctx, int on)
{
	set_rts_go(ctx, on);
}
*/
import "C"
