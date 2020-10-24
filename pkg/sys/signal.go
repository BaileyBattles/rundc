package sys

// +marshal
// +stateify savable
type SignalInfo struct {
	Signo    int32 // Signal number
	Errno    int32 // Errno value
	Code     int32 // Signal code
	Trapno   int
	Pid      int64
	uid      int64
	status   int
	utime    uint64
	stime    uint64
	overrun  int
	timerid  int
	Addr     uintptr
	band     int64
	fd       int
	addr_lsb int16
}
