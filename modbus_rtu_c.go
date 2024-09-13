package libmodbusgo

/*
#cgo CFLAGS: -I${SRCDIR}
#include "modbus.h"

void set_rts_cgo(modbus_t *ctx, int on)
{
	set_rts_go(ctx, on);
}
*/
import "C"
import "sync"

var mapSetRtsCallback = sync.Map{}

//export set_rts_go
func set_rts_go(ctx *C.modbus_t, on C.int) {
	mapSetRtsCallback.Range(func(key, value any) bool {
		f, ok := value.(SetRtsCallback)
		if ok {
			f(&Modbus{ctx: ctx}, int(on))
		} else {
			mapSetRtsCallback.Delete(key)
		}
		return true
	})
}
