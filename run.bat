@echo off
if "%1"=="" (
    echo Usage: run.bat ^<port^>
    exit /b 1
)

set APP_EXTERNAL_PORT=%1
docker-compose build --no-cache
docker-compose up
