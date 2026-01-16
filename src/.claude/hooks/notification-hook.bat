@echo off
setlocal
set "HOOKS_BIN=%CLAUDEX_WINDOWS_HOOKS_BIN%"
if "%HOOKS_BIN%"=="" set "HOOKS_BIN=%CLAUDEX_HOOKS_BIN%"
if "%HOOKS_BIN%"=="" set "HOOKS_BIN=claudex-windows-hooks"
"%HOOKS_BIN%" notification
