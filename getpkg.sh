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

if [ -d "vendor/github.com" ]; then
	if [ "$1" = "update" ]; then
		gvt update -all
	fi
elif [ -f "vendor/manifest" ]; then
	gvt restore -connections 8
else
	pkgs=("github.com/polaris1119/middleware" "github.com/fatih/structs"
	"github.com/go-xorm/xorm" "github.com/fatih/set" "github.com/dchest/captcha"
	"github.com/robfig/cron" "github.com/gorilla/sessions" "github.com/polaris1119/echoutils"
	"golang.org/x/net/websocket" "github.com/polaris1119/slices" "github.com/qiniu/api.v6"
	"github.com/polaris1119/times" "github.com/PuerkitoBio/goquery" "github.com/go-validator/validator"
	"github.com/gorilla/schema" "github.com/facebookgo/grace/gracehttp")

	for pkg in "${pkgs[@]}"; do
		gvt fetch "$pkg"
	done
fi

cd ..

export GOPATH="$OLDGOPATH"

echo 'finished'
