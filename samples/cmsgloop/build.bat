@echo off
go build -buildmode=c-shared -o gw.dll ./dll
gcc -municode -mwindows -o main.exe main.c
