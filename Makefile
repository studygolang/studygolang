.PHONY: build reload start stop

v=""

export GOPROXY=https://goproxy.cn
export GO111MODULE=on

BUILD = $(shell git symbolic-ref HEAD | cut -b 12-)-$(shell git rev-parse HEAD)

build:
	if [ ! -d log ]; then mkdir log; fi

	gofmt -w -s .

	go build -ldflags "-X global.Build=$(BUILD)" -o bin/studygolang github.com/studygolang/studygolang/cmd/studygolang

	@echo "build successfully!"

reload:
	kill -USR2 `cat pid/*.pid`

	echo 'reload successfully'

start:
	if [ ! -d pid ]; then mkdir pid; fi
	export GOTRACEBACK=crash
	ulimit -c unlimited

	bin/studygolang >> log/panic.log 2>&1 &

	@echo "start successfully"

stop:
	kill `cat pid/*.pid`
	sleep 1
	rm -rf pid/*.pid