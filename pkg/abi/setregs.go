package abi

//#cgo CFLAGS: -D_POSIX_C_SOURCE=199309L
//#include "setregs.h"
import "C"

func SetReturnValue(pid, rax int) {
	C.setReturnValue(C.int(pid), C.int(rax))
}

func SetSyscall(pid, orig_rax int) {
	C.setSyscall(C.int(pid), C.int(orig_rax))
}
