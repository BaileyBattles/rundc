.PHONY: all
all:
	go build -o bin/rundc main.go 

run:
	.bin/rundc run /bin/bash