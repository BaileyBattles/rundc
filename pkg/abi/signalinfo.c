#include <assert.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>
#include <sys/ptrace.h>
#include <unistd.h>
#include "signalinfo.h"

struct SignalInfo getSignalInfo(int pid) 
{
    siginfo_t si;
    struct SignalInfo signalInfo;
    memset(&si, 0, sizeof(siginfo_t));
    memset(&signalInfo, 0, sizeof(struct SignalInfo));
    if(ptrace(PTRACE_GETSIGINFO, pid, NULL, &si) != 0) {
        const int err = errno;
        fprintf(stderr, "error grabbing signal, errno: %d\n", err);
        errno = err;
        return signalInfo;
    }
    printf("Address = %p\n", si.si_addr);
    signalInfo.address = si.si_addr;
    signalInfo.pid = 234321;
    signalInfo.si_signo = 1;
    signalInfo.valid = 1;
    return signalInfo;
}