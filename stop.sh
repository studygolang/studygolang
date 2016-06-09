#!/usr/bin/env bash

if [ ! -f stop.sh ]; then
	echo 'stop.sh must be run within its container folder' 1>&2
	exit 1
fi

kill `cat pid/*.pid`

sleep 1
rm -rf pid/*.pid