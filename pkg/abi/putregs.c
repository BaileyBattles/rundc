#include <sys/ptrace.h>
#include <sys/user.h>
#include <stdio.h>
#include "putregs.h"

void putregs(int pid) {
    struct user_regs_struct regs;
    printf("Placing regs\n\n\n");
    ptrace( PTRACE_GETREGS, pid, 0, &regs);
    regs.rax = 25;
    ptrace( PTRACE_SETREGS, pid, 0, &regs);
}