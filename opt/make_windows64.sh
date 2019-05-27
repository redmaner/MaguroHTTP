#!/usr/bin/env bash

source ./opt/build.cfg

rm -f ./opt/magurohttp_windows64.exe
rm -rf ./.temp

export GOOS="windows"
export GOHOSTARCH="amd64"

go get -u golang.org/x/crypto/acme
go get -u golang.org/x/time/rate
go get -u github.com/cespare/xxhash
go get -u golang.org/x/net/idna
go get -u github.com/hashicorp/hcl

go build -o ./opt/magurohttp_windows64.exe

mkdir -p ./.temp
mkdir -p ./out
cp ./opt/magurohttp_windows64.exe ./.temp/magurohttp.exe
cp ./opt/config/docker.hcl ./.temp/main.config

cd ./.temp
zip ../out/magurohttp_"$VERSION"_windows64.zip *
cd ..

rm -rf ./.temp
