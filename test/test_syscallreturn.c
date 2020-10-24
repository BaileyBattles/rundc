#include <unistd.h> 
#include <stdio.h>

int main() {
    write(1, "hello", 5);
}