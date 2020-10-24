#ifndef __ABI_SIGNALINFO__
#define __ABI_SIGNALINFO__

//Struct to represent siginfo_t
//Because pointers are hard to deal with, set valid = 0 on failure of Ptrace
struct SignalInfo {
    int si_signo;
    int si_errno;
    int si_code;
    int si_trapno;
    int pid;
    int uid;
    void *address;
    int valid; //set to invalid on failure
};

void *get_addr(siginfo_t si);
struct SignalInfo getSignalInfo(int pid);


#endif