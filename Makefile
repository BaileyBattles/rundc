.PHONY: all
all:
	go build -gcflags="all=-N -l" -o bin/rundc rundc/entrypoint/main.go
	go build -gcflags="all=-N -l" -o bin/main main.go

run:
	.bin/rundc run /bin/bash

debug:
	dlv --listen=:2345 --headless=true --api-version=2 exec ./bin/main run echo hello

