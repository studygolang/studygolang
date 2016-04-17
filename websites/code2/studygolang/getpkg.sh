#!/usr/bin/env bash

set -e

if [ ! -f getpkg.sh ]; then
    echo 'getpkg.sh must be run within its container folder' 1>&2
    exit 1
fi

OLDGOPATH="$GOPATH"
export GOPATH=`pwd`

cd src

if [ -d "vendor/github.com" ]; then
	if [ "$1" = "update" ]; then
		gvt update -all
	fi
elif [ -f "vendor/manifest" ]; then
	gvt restore -connections 8
else
	pkgs=("github.com/polaris1119/middleware" "github.com/fatih/structs"
	"github.com/go-xorm/xorm" "github.com/fatih/set" "github.com/dchest/captcha"
	"github.com/robfig/cron" "github.com/gorilla/sessions"
	"golang.org/x/net/websocket" "github.com/polaris1119/slices"
	"github.com/gorilla/schema" "github.com/facebookgo/grace/gracehttp")

	for pkg in "${pkgs[@]}"; do
		gvt fetch "$pkg"
	done
fi

cd ..

export GOPATH="$OLDGOPATH"
export PATH="$OLDPATH"

echo 'finished'
