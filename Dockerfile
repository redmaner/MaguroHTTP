FROM golang:alpine AS builder

RUN apk add git \
 	&& go get -u github.com/cespare/xxhash \
	&& go get -u golang.org/x/crypto/acme \
  && go get -u golang.org/x/time/rate \
  && go get -u golang.org/x/net/idna \
  && go get -u github.com/hashicorp/hcl \
	&& go get -u github.com/redmaner/MicroHTTP

WORKDIR /go/src/github.com/redmaner/MicroHTTP

RUN go build -o microhttp

FROM alpine:latest
RUN apk add --no-cache ca-certificates \
	&& mkdir -p /usr/lib/microhttp/www \
	&& mkdir -p /usr/bin \
	&& echo "<html><head></head><body><h1>Welcome to MicroHTTP</h1></body></html>" > /usr/lib/microhttp/www/index.html

COPY --from=builder /go/src/github.com/redmaner/MicroHTTP/microhttp /usr/bin

COPY ./opt/config/docker.json /usr/lib/microhttp/main.json

EXPOSE 80 443

CMD ["/usr/bin/microhttp", "/usr/lib/microhttp/main.json"]
