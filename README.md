waste
=====

Waste is just an weekend experiment playing around with the Docker SDK. This
service does nothing useful, but it does it containerized.

Basically it works like this: An HTTP request comes in, a tar archive is
created from the request body, which is copied to a freshly created container
(alpine by default). A single command is run inside the container to display
the content of the file. The stdout of the command is streamed back to the HTTP
client. Yes, it even has a timeout.

To start the waste server, simply:

```shell
$ waste

██╗    ██╗ █████╗ ███████╗████████╗███████╗
██║    ██║██╔══██╗██╔════╝╚══██╔══╝██╔════╝
██║ █╗ ██║███████║███████╗   ██║   █████╗
██║███╗██║██╔══██║╚════██║   ██║   ██╔══╝
╚███╔███╔╝██║  ██║███████║   ██║   ███████╗
 ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝


Welcome to waste - your premium "cat as a service" provider.

This server accepts HTTP requests and will copy the request body into a
container, run the "cat" command on the input and stream the output back to
stdout.

Example, inspect a local file:

    $ curl http://localhost:3000 --data-binary @README.md

Or run the docker webpage to a docker container first:

    $ curl http://localhost:3000 --data-binary @<(curl -s http://www.docker.io)

Version: 0.1.0
Startup: 2017-09-03 18:25:10.953849056 +0200 CEST m=+0.004078683

DEBU[0000] docker is up: 1.30

```

On a different shell curl something to the server:

```shell
$ curl http://localhost:3000 --data-binary @LICENSE

MIT License

Copyright (c) 2017 Martin Czygan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

The log of the server will tell you what is going on:

```shell
DEBU[0056] request body contains newlines: false
DEBU[0056] archived 0 bytes from request body
DEBU[0056] running with a timeout of 10s
DEBU[0056] creating new docker client
DEBU[0056] pulling image from docker.io/library/alpine
DEBU[0059] creating container from alpine
DEBU[0059] copying data into container
DEBU[0059] 512 bytes written into container
DEBU[0059] stat: {body 0 -rw-r--r-- 2017-09-03 16:26:07 +0000 UTC }
DEBU[0059] starting container f5ee0d60f89cb2542284c6393b51dd355aa492aa8126b80cc0e36648773fe91f
DEBU[0059] request body contains newlines: true
DEBU[0059] archived 1070 bytes from request body
DEBU[0059] running with a timeout of 10s
DEBU[0059] creating new docker client
DEBU[0059] pulling image from docker.io/library/alpine
DEBU[0059] waiting for container f5ee0d60f89cb2542284c6393b51dd355aa492aa8126b80cc0e36648773fe91f
DEBU[0060] 0 bytes read from application
DEBU[0060] removing container f5ee0d60f89cb2542284c6393b51dd355aa492aa8126b80cc0e36648773fe91f
DEBU[0060] operation finished successfully
DEBU[0061] creating container from alpine
DEBU[0061] copying data into container
DEBU[0061] 1582 bytes written into container
DEBU[0061] stat: {body 1070 -rw-r--r-- 2017-09-03 16:26:10 +0000 UTC }
DEBU[0061] starting container f9dc98df5ed0b401d21127940bd1007aec2997de77fcdb485dcfa1528a0f4c22
DEBU[0062] waiting for container f9dc98df5ed0b401d21127940bd1007aec2997de77fcdb485dcfa1528a0f4c22
DEBU[0062] 1091 bytes read from application
DEBU[0062] removing container f9dc98df5ed0b401d21127940bd1007aec2997de77fcdb485dcfa1528a0f4c22
DEBU[0062] operation finished successfully
```

Resources
---------

> The Engine API is the API served by Docker Engine. It allows you to control
every aspect of Docker from within your own applications, build tools to manage
and monitor applications running on Docker, and even use it to build apps on
Docker itself.

* [Engine API](https://docs.docker.com/engine/api/)
* [godoc.org/moby/moby](https://godoc.org/github.com/moby/moby)
