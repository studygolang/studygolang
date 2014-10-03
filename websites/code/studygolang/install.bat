@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;%~dp0..\thirdparty

if not exist log mkdir log

gofmt -w src

:: -tags "debug" 表示测试
go install -tags "debug" server/studygolang
go install -tags "debug" server/crawlarticle

set GOPATH=%OLDGOPATH%

:end
echo finished