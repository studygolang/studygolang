@echo off

setlocal

if exist clean.bat goto ok
echo clean.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;

go clean -i -r studygolang

set GOPATH=%OLDGOPATH%

:end
echo finished