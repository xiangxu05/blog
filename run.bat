@echo off
echo ========================================
echo Personal Blog System Startup Script
echo ========================================
echo.

echo Starting blog system...
echo Target: cmd/blogsys/main.go
echo.

:: Change to project directory
cd /d "%~dp0"

:: Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/
    pause
    exit /b 1
)

:: Check if main.go exists
if not exist "cmd\blogsys\main.go" (
    echo Error: cmd/blogsys/main.go not found
    echo Please make sure the file exists
    pause
    exit /b 1
)

:: Run the blog system
echo Running blog system...
go run ./cmd/blogsys/main.go

:: Check if the program exited with error
if errorlevel 1 (
    echo.
    echo Error: Blog system exited with error code %errorlevel%
    echo Press any key to exit...
    pause >nul
    exit /b %errorlevel%
) else (
    echo.
    echo Blog system stopped normally
    echo Press any key to exit...
    pause >nul
)