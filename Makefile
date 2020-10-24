.PHONY: all
all:
	go build -o bin/main pkg/main.go 

run:
	./bin/main