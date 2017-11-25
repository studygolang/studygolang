#!/usr/bin/env bash

set -e

if [ ! -f install.sh ]; then
	echo 'install must be run within its container folder' 1>&2
	exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN=

if [ ! -d log ]; then
	mkdir log
fi

gofmt -w -s src

BUILD="`git symbolic-ref HEAD | cut -b 12-`-`git rev-parse HEAD`"

go install -ldflags "-X global.Build="$BUILD server/studygolang
go install server/indexer
go install server/crawler
go install server/migrator

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'finished'

