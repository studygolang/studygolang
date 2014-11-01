@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;%~dp0..\thirdparty

if not exist log mkdir log

set VERSION=git symbolic-ref HEAD | cut -b 12-
set VERSION=%VERSION%-git rev-parse HEAD

gofmt -w src

:: -tags "debug" 表示测试
go install -tags "debug" -ldflags "-X util.version 1.0.0 -X util.date %date%" ./...

set GOPATH=%OLDGOPATH%

:end
echo finished