package abi

//#cgo CFLAGS: -D_POSIX_C_SOURCE=199309L
//#include <signal.h>
//#include "signalinfo.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type SignalInfo struct {
	Pid  int
	Addr unsafe.Pointer
}

func GetSignalInfo(pid int) (*SignalInfo, error) {
	si := C.getSignalInfo(C.int(pid))
	if !ValidSignalInfo(si) {
		return nil, fmt.Errorf("signal info struct returned from signalinfo.c has valid field set to 0")
	}
	return &SignalInfo{
		Pid:  int(si.pid),
		Addr: unsafe.Pointer(si.address),
	}, nil
}

func ValidSignalInfo(si C.struct_SignalInfo) bool {
	if si.valid == 0 {
		return false
	}
	return true
}
