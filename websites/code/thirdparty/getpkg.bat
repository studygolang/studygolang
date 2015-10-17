@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go get -u -v github.com/go-sql-driver/mysql
go get -u -v github.com/studygolang/mux
go get -u -v github.com/gorilla/sessions
go get -u -v github.com/robfig/cron
go get -u -v github.com/dchest/captcha
go get -u -v github.com/andybalholm/cascadia
go get -u -v github.com/qiniu/api.v6

if not exist "src/golang.org/x/text" (
	git clone https://github.com/golang/text src/golang.org/x/text
)
go install golang.org/x/text/...

if not exist "src/golang.org/x/crypto" (
	git clone https://github.com/golang/crypto src/golang.org/x/crypto
)
go install golang.org/x/crypto/...

if not exist "src/golang.org/x/crypto" (
	git clone https://github.com/golang/crypto src/golang.org/x/crypto
)
go install golang.org/x/crypto/...

go get -u -v github.com/PuerkitoBio/goquery

set GOPATH=%OLDGOPATH%

:end
echo finished
