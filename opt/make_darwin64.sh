#!/usr/bin/env bash

source ./opt/build.cfg

rm -f ./opt/microhttp_darwin64
rm -rf ./.temp

export GOOS="darwin"
export GOHOSTARCH="amd64"

go get -u golang.org/x/crypto/acme
go get -u golang.org/x/time/rate
go get -u github.com/nu7hatch/gouuid
go get -u github.com/cespare/xxhash
go get -u golang.org/x/crypto/bcrypt
go get -u golang.org/x/net/idna
go get -u github.com/hashicorp/hcl

go build -o ./opt/microhttp_darwin64

mkdir -p ./.temp
mkdir -p ./out
cp ./opt/microhttp_darwin64 ./.temp/microhttp
cp ./opt/config/example.json ./.temp/main.json

cd ./.temp
zip ../out/microhttp_"$VERSION"_darwin64.zip *
cd ..

rm -rf ./.temp
