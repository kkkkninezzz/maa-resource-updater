@echo off
go build -ldflags="-s -w" -o maa-resource-updater.exe ./cmd/main/main.go