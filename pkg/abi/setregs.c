#include <sys/ptrace.h>
#include <sys/user.h>
#include <stdio.h>
#include "setregs.h"

void setReturnValue(int pid, int rax) {
    struct user_regs_struct regs;
    ptrace( PTRACE_GETREGS, pid, 0, &regs);
    regs.rax = rax;
    ptrace( PTRACE_SETREGS, pid, 0, &regs);
}

void setSyscall(int pid, int orig_rax) {
    struct user_regs_struct regs;
    ptrace( PTRACE_GETREGS, pid, 0, &regs);
    regs.orig_rax = orig_rax;
    ptrace( PTRACE_SETREGS, pid, 0, &regs);
}

void setRegs(int pid, int rax, int rdi, int rsi,
             int rdx, int r10, int r8, int r9) {
    struct user_regs_struct regs;
    ptrace( PTRACE_GETREGS, pid, 0, &regs);
    regs.rax = rax;
    regs.rdi = rdi;
    regs.rsi = rsi;
    regs.rdx = rdx;
    regs.r10 = r10;
    regs.r8 = r8;
    regs.r9 = r9;
    ptrace( PTRACE_SETREGS, pid, 0, &regs);
}