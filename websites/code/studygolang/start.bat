@echo off

setlocal

if exist start.bat goto ok
echo start.bat must be run from its folder
goto end

:ok

start /b bin\studygolang.exe >> log\panic.log 2>&1 &

echo start successfully

:end