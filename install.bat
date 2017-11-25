@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set GOBIN=
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

if not exist log mkdir log

gofmt -w -s src

go install server/studygolang
go install server/indexer
go install server/crawler
go install server/migrator

set GOPATH=%OLDGOPATH%

:end
echo finished