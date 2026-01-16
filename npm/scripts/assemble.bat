@echo off
setlocal
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0assemble.ps1"
