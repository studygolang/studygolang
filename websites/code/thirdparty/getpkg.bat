@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go get -u github.com/go-sql-driver/mysql
go get -u github.com/studygolang/mux
go get -u github.com/gorilla/sessions

set GOPATH=%OLDGOPATH%

:end
echo finished