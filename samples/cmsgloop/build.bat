@echo off
go build -buildmode=c-shared -o gw.dll ./dll
gcc -municode -o main.exe main.c
