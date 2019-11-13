# MaguroHTTP
MaguroHTTP is a simple and small HTTP webserver and reverse proxy, with an extensive focus on performance and security.

You can find more info @ https://magurohttp.io

# Repository overview
This repository contains the source code of MaguroHTTP, which is written in GO. This repository consists of:

* / - The root of the repository contains the code of the command line interface of MaguroHTTP
* /tuna - The tuna folder holds the MaguroHTTP server
* /router - The router folder holds the MaguroHTTP HTTP router
* /guard - The guard folder holds various security focused HTTP Middleware
* /html - The html folder provides serveral objects to help with HTML templating
* /debug - The debug folder holds the MaguroHTTP logging facility
* /cache - The cache folder provides a memory allocated, time associated, concurrency safe cache
* /opt - The opt folder has serveral tools which can be used to build and deploy MaguroHTTP
