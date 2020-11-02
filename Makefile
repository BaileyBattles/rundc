.PHONY: all
all:
	go build -gcflags="all=-N -l" -o bin/rundc main.go

run:
	.bin/rundc run /bin/bash

debug:
	dlv --listen=:2345 --headless=true --api-version=2 exec ./bin/rundc pull alpine

