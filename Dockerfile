FROM golang:alpine AS builder

RUN mkdir /gobuild

COPY ./*.go /gobuild/

WORKDIR /gobuild

RUN go get -v  github.com/gbrlsnchs/jwt \
	&& go build -o microhttp *.go

FROM alpine:latest
COPY --from=builder /gobuild/microhttp /usr/local/bin
COPY ./config.json /usr/local/bin/microhttp.json
CMD ["/usr/local/bin/microhttp", "/usr/local/bin/microhttp.json"]
