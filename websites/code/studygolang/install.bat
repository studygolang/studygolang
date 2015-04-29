@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;%~dp0..\thirdparty

if not exist log mkdir log

for /f "delims=" %%t in ('git symbolic-ref HEAD') do set VERNAME=%%t
set VERNAME=%VERNAME:~11%

for /f "delims=" %%t in ('git rev-parse HEAD') do set VERCODE=%%t

set VERSION=%VERNAME%-%VERCODE%

gofmt -w src

:: -tags "debug" 表示测试
go install -tags "debug" -ldflags "-X util.version.Version "%VERSION% ./...

set GOPATH=%OLDGOPATH%

:end
echo finished