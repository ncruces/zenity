#!/bin/bash

GOOS=windows GOARCH=386 go build -ldflags="-s -w" &&
zip -9 zenity_win32.zip zenity.exe

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" &&
zip -9 zenity_win64.zip zenity.exe

GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" &&
zip -9 zenity_macos.zip zenity

go build
