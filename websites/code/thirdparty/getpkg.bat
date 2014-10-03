@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

:go get -u github.com/go-sql-driver/mysql
:go get -u github.com/studygolang/mux
:go get -u github.com/gorilla/sessions
:go get -u github.com/robfig/cron
:go get -u github.com/PuerkitoBio/goquery
go get -u github.com/qiniu/api
:hg clone https://code.google.com/p/go.net src/go.net
:go install go.net/websocket

set GOPATH=%OLDGOPATH%

:end
echo finished