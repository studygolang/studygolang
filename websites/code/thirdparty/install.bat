@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go install github.com/go-sql-driver/mysql
go install github.com/studygolang/mux
go install github.com/gorilla/sessions
go install go.net/websocket

set GOPATH=%OLDGOPATH%

:end
echo finished