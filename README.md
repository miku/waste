waste
=====

This is just an weekend experiment playing around with the [Docker
SDK](https://docs.docker.com/engine/api/sdks/). This service does nothing
useful, but it does it in a containerized fashion.

Basically: An *HTTP request* comes in, a *tar archive* is created from the
*request body*, which is *copied* to a freshly created *container* (alpine by
default). A single *command* is *run* inside the container to display the content
of the file. The *stdout* of the command is *streamed* back to the HTTP
*client*. That's all. And yes, it even has a *timeout*.

To start the waste server, simply:

```shell
$ waste

██╗    ██╗ █████╗ ███████╗████████╗███████╗
██║    ██║██╔══██╗██╔════╝╚══██╔══╝██╔════╝
██║ █╗ ██║███████║███████╗   ██║   █████╗
██║███╗██║██╔══██║╚════██║   ██║   ██╔══╝
╚███╔███╔╝██║  ██║███████║   ██║   ███████╗
 ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝


Welcome to waste - your friendly "cat as a service" provider.

This server accepts HTTP requests and will copy the request body into a
container, run the "cat" command on the input and stream the output back to
stdout.

Example, inspect a local file:

    $ curl http://localhost:3000 --data-binary @README.md

Or run the docker homepage through a docker container first:

    $ curl http://localhost:3000 --data-binary @<(curl -sL http://www.docker.com)

Version: 0.1.0
Startup: 2017-09-03 18:25:10.953849056 +0200 CEST m=+0.004078683

DEBU[0000] docker is up: 1.30

```

In a different terminal curl something to the server:

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
DEBU[0001] request body contains newlines: true
DEBU[0001] archived 1070 bytes from request body
DEBU[0001] running with a timeout of 10s
DEBU[0001] creating new docker client
DEBU[0001] pulling image from docker.io/library/alpine
DEBU[0005] creating container from alpine
DEBU[0005] copying data into container
DEBU[0005] 1582 bytes written into container
DEBU[0005] stat: {body 1070 -rw-r--r-- 2017-09-03 16:34:07 +0000 UTC }
DEBU[0005] starting container ce1e84c442ab22579110b40537e6485d0a872bc81cbb2e165de2c6fabc254b4d
DEBU[0006] waiting for container ce1e84c442ab22579110b40537e6485d0a872bc81cbb2e165de2c6fabc254b4d
DEBU[0006] 1091 bytes read from application
DEBU[0006] removing container ce1e84c442ab22579110b40537e6485d0a872bc81cbb2e165de2c6fabc254b4d
DEBU[0006] operation finished successfully
```

Try
---

```shell
$ git clone https://github.com/miku/waste.git
$ cd waste
$ make
$ ./waste
```

Don't be too new
----------------

```shell
$ DOCKER_API_VERSION=1.35 ./waste
```

Resources
---------

> The Engine API is the API served by Docker Engine. It allows you to control
every aspect of Docker from within your own applications, build tools to manage
and monitor applications running on Docker, and even use it to build apps on
Docker itself.

* [Engine API](https://docs.docker.com/engine/api/)
* [godoc.org/moby/moby](https://godoc.org/github.com/moby/moby)
