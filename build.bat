@echo off
setlocal enabledelayedexpansion

set "SRC_DIR=src"
set "BIN_DIR=%LOCALAPPDATA%\claudex-windows\bin"
set "CONFIG_DIR=%LOCALAPPDATA%\claudex-windows"

if "%1"=="" goto :help
if /I "%1"=="build" goto :build
if /I "%1"=="build-hooks" goto :build_hooks
if /I "%1"=="deps" goto :deps
if /I "%1"=="test" goto :test
if /I "%1"=="fmt" goto :fmt
if /I "%1"=="vet" goto :vet
if /I "%1"=="check" goto :check
if /I "%1"=="install" goto :install
if /I "%1"=="install-hooks" goto :install_hooks
if /I "%1"=="install-mcp" goto :install_mcp
if /I "%1"=="install-project" goto :install_project
if /I "%1"=="clean" goto :clean
if /I "%1"=="run" goto :run
if /I "%1"=="help" goto :help

:help
echo Available commands:
echo   build            - Build claudex-windows.exe
echo   build-hooks      - Build claudex-windows-hooks.exe
echo   deps             - Install/update dependencies
echo   test             - Run tests
echo   fmt              - Format code with go fmt
echo   vet              - Run go vet
echo   check            - Run fmt, vet, and test
echo   install          - Install binaries and hooks under %%LOCALAPPDATA%%\claudex-windows
echo   install-hooks    - Install hooks only
echo   install-mcp       - Configure recommended MCP servers
echo   install-project  - Install profiles/hooks to current project .claude
echo   clean            - Remove build artifacts
echo   run              - Build and run claudex-windows.exe
echo   help             - Show this help
exit /b 0

:version
for /f "delims=" %%V in ('git describe --tags --always --dirty 2^>nul') do set "VERSION=%%V"
if "%VERSION%"=="" set "VERSION=dev"
exit /b 0

:build
call :version
echo Building claudex-windows %%VERSION%%...
pushd %SRC_DIR%
go build -ldflags "-X main.Version=%VERSION%" -o ..\claudex-windows.exe .\cmd\claudex
popd
echo Built: claudex-windows %%VERSION%%
exit /b %ERRORLEVEL%

:build_hooks
echo Building claudex-windows-hooks...
pushd %SRC_DIR%
go build -o ..\bin\claudex-windows-hooks.exe .\cmd\claudex-hooks
popd
echo Built: bin\claudex-windows-hooks.exe
exit /b %ERRORLEVEL%

:deps
echo Installing dependencies...
pushd %SRC_DIR%
go mod tidy
popd
echo Dependencies installed
exit /b %ERRORLEVEL%

:test
echo Running tests...
pushd %SRC_DIR%
go test -v ./...
popd
echo Tests complete
exit /b %ERRORLEVEL%

:fmt
echo Formatting code...
pushd %SRC_DIR%
go fmt ./...
popd
echo Formatting complete
exit /b %ERRORLEVEL%

:vet
echo Vetting code...
pushd %SRC_DIR%
go vet ./...
popd
echo Vet complete
exit /b %ERRORLEVEL%

:check
call :fmt
if errorlevel 1 exit /b 1
call :vet
if errorlevel 1 exit /b 1
call :test
exit /b %ERRORLEVEL%

:install_hooks
call :build_hooks
if errorlevel 1 exit /b 1
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%CONFIG_DIR%\hooks" mkdir "%CONFIG_DIR%\hooks"
copy /y bin\claudex-windows-hooks.exe "%BIN_DIR%\claudex-windows-hooks.exe" >nul
copy /y %SRC_DIR%\scripts\proxies\*.bat "%CONFIG_DIR%\hooks" >nul
copy /y %SRC_DIR%\scripts\proxies\*.ps1 "%CONFIG_DIR%\hooks" >nul
echo Installed hooks to %CONFIG_DIR%\hooks
exit /b 0

:install_mcp
call :build
if errorlevel 1 exit /b 1
claudex-windows.exe --setup-mcp
exit /b %ERRORLEVEL%

:install
call :build
if errorlevel 1 exit /b 1
call :build_hooks
if errorlevel 1 exit /b 1
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
xcopy /e /i /y %SRC_DIR%\profiles "%CONFIG_DIR%\profiles" >nul
copy /y claudex-windows.exe "%BIN_DIR%\claudex-windows.exe" >nul
copy /y bin\claudex-windows-hooks.exe "%BIN_DIR%\claudex-windows-hooks.exe" >nul
if not exist "%CONFIG_DIR%\hooks" mkdir "%CONFIG_DIR%\hooks"
copy /y %SRC_DIR%\scripts\proxies\*.bat "%CONFIG_DIR%\hooks" >nul
copy /y %SRC_DIR%\scripts\proxies\*.ps1 "%CONFIG_DIR%\hooks" >nul
echo Installed to %CONFIG_DIR%
echo Add %BIN_DIR% to your PATH if needed
exit /b 0

:install_project
echo Installing claudex-windows to current project...
if not exist .claude mkdir .claude
if exist "%CONFIG_DIR%\profiles" (
    xcopy /e /i /y "%CONFIG_DIR%\profiles" .claude\profiles >nul
    echo Copied profiles from %CONFIG_DIR%
) else if exist profiles (
    xcopy /e /i /y profiles .claude\profiles >nul
    echo Copied profiles from local directory
) else (
    echo No profiles directory found
    exit /b 1
)
if exist "%CONFIG_DIR%\hooks" (
    xcopy /e /i /y "%CONFIG_DIR%\hooks" .claude\hooks >nul
    echo Copied hooks from %CONFIG_DIR%
) else if exist .claude\hooks (
    echo Hooks already exist in .claude\
) else (
    echo No hooks directory found
)
echo Project installation complete
exit /b 0

:clean
echo Cleaning build artifacts...
del /q claudex-windows.exe 2>nul
rmdir /s /q bin 2>nul
echo Cleaned
exit /b 0

:run
call :build
if errorlevel 1 exit /b 1
claudex-windows.exe
exit /b %ERRORLEVEL%
