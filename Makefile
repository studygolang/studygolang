.PHONY: build reload start stop

v=""

BUILD = $(shell git symbolic-ref HEAD | cut -b 12-)-$(shell git rev-parse HEAD)

build:
	if [ ! -d log ]; then mkdir log; fi

	go build -ldflags "-X global.Build=$(BUILD)" -o bin/studygolang github.com/studygolang/studygolang/cmd/studygolang

reload:
	./reload.sh	

start:
	./start.sh

stop:
	./stop.sh