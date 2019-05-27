FROM golang:alpine AS builder

RUN apk add git \
 	&& go get -u github.com/cespare/xxhash \
	&& go get -u golang.org/x/crypto/acme \
  && go get -u golang.org/x/time/rate \
  && go get -u golang.org/x/net/idna \
  && go get -u github.com/hashicorp/hcl \
	&& go get -u github.com/redmaner/MaguroHTTP

WORKDIR /go/src/github.com/redmaner/MaguroHTTP

RUN go build -o magurohttp

FROM alpine:latest
RUN apk add --no-cache ca-certificates \
	&& mkdir -p /usr/lib/magurohttp/www \
	&& mkdir -p /usr/bin \
	&& echo "<html><head></head><body><h1>Welcome to MaguroHTTP</h1></body></html>" > /usr/lib/magurohttp/www/index.html

COPY --from=builder /go/src/github.com/redmaner/MaguroHTTP/magurohttp /usr/bin

COPY ./opt/config/docker.json /usr/lib/magurohttp/main.json

EXPOSE 80 443

CMD ["/usr/bin/magurohttp", "/usr/lib/magurohttp/main.json"]
