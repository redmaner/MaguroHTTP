FROM golang:alpine AS builder

RUN mkdir /gobuild && apk add git

COPY ./*.go /gobuild/

WORKDIR /gobuild

RUN go get -v  github.com/gbrlsnchs/jwt  \
	&& go get -v github.com/redmaner/smux \
	&& go build -o microhttp *.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates \
	&& mkdir -p /usr/lib/microhttp/www \
	&& mkdir -p /usr/bin \
	&& echo "<html><head></head><body><h1>Welcome to MicroHTTP</h1></body></html>" > /usr/lib/microhttp/www/index.html

COPY --from=builder /gobuild/microhttp /usr/bin

COPY ./opt/config/docker.json /usr/lib/microhttp/main.json

EXPOSE 80 443

CMD ["/usr/bin/microhttp", "/usr/lib/microhttp/main.json"]
