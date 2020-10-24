#include <unistd.h> 
#include <stdio.h>

int main() {
    int fd = unlink("apath");
    if (fd != 77) {
        printf("File descriptor = %d\n", fd);
    }
}