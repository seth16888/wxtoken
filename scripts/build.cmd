@echo off
for /f %%i in ('git describe --tags --abbrev^=0') do set VERSION=%%i
echo Building version %VERSION%
go build -ldflags="-X github.com/seth16888/wxtoken/internal/cmd.Version=%VERSION%" -o bin/ ./
