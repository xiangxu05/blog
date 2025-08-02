
@echo off
echo ========================================
echo Personal Blog Database Regeneration Script
echo ========================================
echo.

echo [1/4] Checking database file status...
if exist "data\blogData.db" (
    echo Found existing database file, creating backup...
    copy "data\blogData.db" "data\blogData.db.backup" >nul 2>&1
    if errorlevel 1 (
        echo Warning: Backup failed, continuing...
    ) else (
        echo Backup completed: data\blogData.db.backup
    )
    
    echo Deleting existing database file...
    del "data\blogData.db" >nul 2>&1
    if errorlevel 1 (
        echo Error: Cannot delete existing database file
        pause
        exit /b 1
    )
)
echo.

echo [2/4] Initializing database structure...
go run ./cmd/init_db/main.go
if errorlevel 1 (
    echo Error: Database initialization failed
    echo Restoring backup...
    if exist "data\blogData.db.backup" (
        copy "data\blogData.db.backup" "data\blogData.db" >nul 2>&1
        echo Backup restored
    )
    pause
    exit /b 1
)
echo Database initialization completed!
echo.

echo [3/4] Generating DAO layer code...
go run ./cmd/gen/main.go
if errorlevel 1 (
    echo Error: DAO code generation failed
    pause
    exit /b 1
)
echo DAO code generation completed!
echo.

echo [4/4] Cleaning backup files...
if exist "data\blogData.db.backup" (
    del "data\blogData.db.backup" >nul 2>&1
    echo Backup files cleaned
)
echo.

echo ========================================
echo Database regeneration completed!
echo ========================================
echo Generated files:
echo   - Database: data\blogData.db
echo   - DAO code: dao\*.gen.go
echo.
echo Press any key to exit...
pause >nul
