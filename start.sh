#!/usr/bin/env bash

set -e

if [ ! -f start.sh ]; then
	echo 'start.sh must be run within its container folder' 1>&2
	exit 1
fi

if [ ! -d log ]; then
	mkdir log
fi

if [ ! -d pid ]; then
	mkdir pid
fi

export GOTRACEBACK=crash
ulimit -c unlimited

bin/studygolang >> log/panic.log 2>&1 &

echo "start successfully"
