# MicroHTTP
MicroHTTP is a small, fast, stable and secure HTTP server and reverse proxy writtin in Go.

## Features
MicroHTTP is currently being developed and supports the following features:<br>
* HTTP server
* Reverse proxy
* Support for virtual hosts to host multiple websites on one webserver
* TLSv1.2 support with automatic strong ciphers
* Support for RSA and Elliptic Curve certificates
* Configurable Strict-Transport-Security
* Security headers are set by default
* HTTP Headers are easily configurable
* HTTP Methods can be easily enabled and disabled. This applies to the whole website, and can even be configured for indivual pages if so desired
* Automatic response Content-Type
* Configurable request Content-Type
* Firewall support for both HTTP server and reverse proxy
* Flexible configuration with json<br>

## Backlog
The following features will be added shortly:<br>
* Support for HTTP2
* Caching support for both HTTP server and reverse proxy
* Network logging and application logging, both configurable<br>

## Not supported
The following features are not supported for now or ever:<br>
* TLSv1.3 (will be added in the future)
* Diffe-Hellman exchange (DHE) cipher suites
* FastCGI
* WebDAV
* PHP
* Trash and proprietary or prehistoric technologies
