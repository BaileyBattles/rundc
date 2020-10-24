#ifndef __ABI_PUTREGS__
#define __ABI_PUTREGS__

void setReturnValue(int pid, int rax);

void setSyscall(int pid, int rax);

void setRegs(int pid, int rax, int rdi, int rsi,
             int rdx, int r10, int r8, int r9);

#endif