#!/usr/bin/env bash

source ./opt/build.cfg

rm -f ./opt/microhttp_windows64.exe
rm -rf ./.temp

export GOOS="windows"
export GOHOSTARCH="amd64"

go get -v -u golang.org/x/crypto/acme
go get -v -u golang.org/x/time/rate
go get -v -u github.com/cespare/xxhash

go build -o ./opt/microhttp_windows64.exe *.go

mkdir -p ./.temp
mkdir -p ./out
cp ./opt/microhttp_windows64.exe ./.temp/microhttp.exe
cp ./opt/config/example.json ./.temp/main.json

cd ./.temp
zip ../out/microhttp_"$VERSION"_windows64.zip *
cd ..

rm -rf ./.temp
