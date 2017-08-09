@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

cd src

gvt restore -connections 8

cd ..

set GOPATH=%OLDGOPATH%

:end
echo finished

