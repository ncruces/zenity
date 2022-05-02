#!/bin/bash

TAG=$(git tag --points-at HEAD)
echo 'package main; const tag = "'$TAG'"' > tag.go

go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo

GOOS=windows GOARCH=386   CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath &&
zip -9 zenity_win32.zip zenity.exe

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath &&
zip -9 zenity_win64.zip zenity.exe

rm resource.syso

GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o zenity_macos_x64 &&
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o zenity_macos_arm &&
go run github.com/randall77/makefat zenity zenity_macos_x64 zenity_macos_arm &&
zip -9 zenity_macos.zip zenity

zip -9 zenity_brew.zip zenity zenity.exe
rm zenity zenity_macos_* zenity.exe

GOOS=linux go build -tags dev
go build -tags dev
git restore tag.go
