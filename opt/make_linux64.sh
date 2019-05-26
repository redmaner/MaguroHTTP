#!/usr/bin/env bash

source ./opt/build.cfg

rm -f ./opt/microhttp_linux64
rm -f ./microhttp_linux64.tar.gz
rm -rf ./.temp

export GOOS="linux"
export GOHOSTARCH="amd64"

go get -u golang.org/x/crypto/acme
go get -u golang.org/x/time/rate
go get -u github.com/cespare/xxhash
go get -u golang.org/x/net/idna
go get -u github.com/hashicorp/hcl

go build -o ./opt/microhttp_linux64

mkdir -p ./.temp
mkdir -p ./out
cp ./opt/microhttp_linux64 ./.temp/microhttp
cp ./opt/systemd/microhttp.service ./.temp/microhttp.service
cp ./opt/config/docker.hcl ./.temp/main.config

cd ./.temp
zip ../out/microhttp_"$VERSION"_linux64.zip *
cd ..

rm -rf ./.temp
