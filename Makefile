.PHONY: all
all:
	go build -o bin/main rundc/main.go 

run:
	./bin/main