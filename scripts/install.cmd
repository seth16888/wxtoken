@echo off
for /f %%i in ('git describe --tags --abbrev^=0') do set VERSION=%%i
echo Install version %VERSION%
go install -ldflags="-X github.com/seth16888/wxtoken/internal/cmd.Version=%VERSION%" ./...
