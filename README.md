# MicroHTTP v0.11
MicroHTTP is a small, fast, stable and secure HTTP server and reverse proxy writtin in Go.

## Features
MicroHTTP is currently being developed and supports the following features:<br>
* HTTP2 server
* HTTP2 proxy
* Support for virtual hosts, to serve multiple websites and / or proxies on one host
* TLSv1.2 support with automatic strong ciphers
* Support for RSA and Elliptic Curve certificates
* Configurable Strict-Transport-Security
* Security headers are set by default
* HTTP Headers are easily configurable
* HTTP Methods can be easily enabled and disabled. This applies to the whole website, and can even be configured for indivual pages if so desired
* Automatic response Content-Type
* Configurable request Content-Type
* Firewall support for server and proxy
* Flexible configuration with json<br>

## Backlog
The following features will be added shortly:<br>
* Automatic TLS certifcates with LetsEncrypt
* Caching support for both HTTP server and reverse proxy

## Not supported
The following features are not supported for now or ever:<br>
* TLSv1.3 (will be added in the future)
* Diffe-Hellman exchange (DHE) cipher suites
* FastCGI
* WebDAV
* PHP
* Websocket
* Trash and proprietary or prehistoric technologies
