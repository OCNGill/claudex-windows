@echo off
setlocal
if "%CLAUDEX_WINDOWS_NOTIFICATIONS_ENABLED%"=="" set "CLAUDEX_WINDOWS_NOTIFICATIONS_ENABLED=%CLAUDEX_NOTIFICATIONS_ENABLED%"
if /I "%CLAUDEX_WINDOWS_NOTIFICATIONS_ENABLED%"=="false" exit /b 0
powershell -NoProfile -ExecutionPolicy Bypass -Command "$shell = New-Object -ComObject WScript.Shell; [void]$shell.Popup('Agent complete', 5, 'Claudex Windows', 0x40)"
