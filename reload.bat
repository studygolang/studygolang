@echo off

setlocal

if exist reload.bat goto ok
echo reload.bat must be run from its folder
goto end

:ok

:: stop
taskkill /im studygolang.exe /f
del /q /f /a pid

:: start
start /b bin\studygolang >> log\panic.log 2>&1 &

echo reload successfully

:end