#!/bin/bash

go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo

GOOS=windows GOARCH=386 go build -ldflags="-s -w" -trimpath &&
zip -9 zenity_win32.zip zenity.exe

GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath &&
zip -9 zenity_win64.zip zenity.exe

GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath &&
zip -9 zenity_macos.zip zenity

go build -tags dev
