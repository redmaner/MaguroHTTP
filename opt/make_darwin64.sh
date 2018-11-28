#!/usr/bin/env bash

source ./opt/build.cfg

rm -f ./opt/microhttp_darwin64
rm -rf ./.temp

export GOOS="windows"
export GOHOSTARCH="amd64"

go get -v -u github.com/gbrlsnchs/jwt
go get -v -u github.com/redmaner/smux

go build -o ./opt/microhttp_darwin64 *.go

mkdir -p ./.temp
mkdir -p ./out
cp ./opt/microhttp_darwin64 ./.temp/microhttp
cp ./opt/config/example.json ./.temp/main.json

cd ./.temp
zip ../out/microhttp_"$VERSION"_darwin64.zip *
cd ..

rm -rf ./.temp
