@echo off

setlocal

if exist clean.bat goto ok
echo clean.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;%~dp0..\thirdparty

go clean -i -r studygolang

rd /s /q pid
rd /s /q bin
rd /s /q pkg
del /q /f /a log

set GOPATH=%OLDGOPATH%

:end
echo finished