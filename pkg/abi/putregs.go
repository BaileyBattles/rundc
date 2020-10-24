package abi

//#cgo CFLAGS: -D_POSIX_C_SOURCE=199309L
//#include "putregs.h"
import "C"

func PutRegs(pid int) {
	C.putregs(C.int(pid))
}
