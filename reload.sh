#!/usr/bin/env bash

set -e

if [ ! -f reload.sh ]; then
	echo 'reload.sh must be run within its container folder' 1>&2
	exit 1
fi

kill -USR2 `cat pid/*.pid`

echo 'reload successfully'
