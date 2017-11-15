#!/usr/bin/env bash

set -e

if [ ! -f getpkg.sh ]; then
    echo 'getpkg.sh must be run within its container folder' 1>&2
    exit 1
fi

if ! type gvt >/dev/null 2>&1; then
	echo >&2 "This script requires the gvt tool."
	echo >&2 "You may obtain it with the following command:"
	echo >&2 "go get github.com/polaris1119/gvt"
	exit 1
fi

OLDGOPATH="$GOPATH"
export GOPATH=`pwd`

cd src

if [ "$1" = "update" ]; then
    if [ -d "vendor/github.com" ]; then
        gvt update -all
    fi
elif [ -f "vendor/manifest" ]; then
    gvt restore -connections 8 -precaire
fi

cd ..

export GOPATH="$OLDGOPATH"

echo 'finished'
