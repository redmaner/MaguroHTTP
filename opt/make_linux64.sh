#!/usr/bin/env bash

rm -f ./opt/microhttp_linux64
rm -f ./microhttp_linux64.tar.gz
rm -rf ./.temp

export GOOS="linux"
export GOHOSTARCH="amd64"

go build -o ./opt/microhttp_linux64 *.go

mkdir -p ./.temp
cp ./opt/microhttp_linux64 ./.temp/microhttp
cp ./opt/systemd/microhttp.service ./.temp/microhttp.service
cp ./opt/config/example.json ./.temp/main.json

cd ./.temp
zip ../opt/microhttp_linux64.zip *
cd ..

rm -rf ./.temp
