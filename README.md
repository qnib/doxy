# doxy
Docker unix-socket proxy to provide unharmful, read-only API calls

## Usage

```bash
$ ./doxy_darwin --help
*snip*
GLOBAL OPTIONS:
   --docker-socket value  Docker host to connect to. (default: "/var/run/docker.sock") [$DOXY_DOCKER_SOCKET]
   --proxy-socket value   Proxy socket to be created (default: "/tmp/doxy.sock") [$DOXY_PROXY_SOCKET]
   --debug                Print proxy requests [$DOXY_DEBUG]
   --pattern-file value   File holding line-separated regex-patterns to be allowed (comments allowed, use #) (default: "/etc/doxy.pattern") [$DOXY_PATTERN_FILE]
   --help, -h             show help
   --version, -v          print the version
$ ./doxy_darwin
2017/08/18 11:37:43 [II] Start Version: 0.1.0
2017/08/18 11:37:43 Error reading patterns file (open /etc/doxy.pattern: no such file or directory), using default patterns
2017/08/18 11:37:43 [gk-soxy] Listening on /tmp/doxy.sock
```

## Filter mechanism

### Request Method

For starters the proxy only allows `GET` requests.

```bash
$ docker -H unix:///tmp/doxy.sock run ubuntu bash
docker: Error response from daemon: Only GET requests are allowed, req.Method: POST.
See 'docker run --help'.
```

### Regex

Once the method is checked, a list of regular expressions are checked. In version 0.1.0 the list reads:

```bash
# List and inspect containers
^/(v\d\.\d+/)?containers(/\w+)?/json$
# List and inspect services
^/(v\d\.\d+/)?services(/[0-9a-f]+)?$
# List and inspect tasks
^/(v\d\.\d+/)?tasks(/\w+)?$
# List and inspect networks
^/(v\d\.\d+/)?networks(/\w+)?$
# List and inspect nodes
^/(v\d\.\d+/)?nodes(/\w+)?$
# Show engine info
^/(v\d\.\d+/)?info$
# Healthcheck
^/_ping$
```

Thus, an export of a container filesystem is not allowed.

```bash
$ docker -H unix:///tmp/doxy.sock export -o test.tar $(docker ps -lq)
Error response from daemon: '/v1.31/containers/a62250e0890a/export' is not allowed.
```

## Debug output

The tool uses [negroni](https://github.com/urfave/negroni), a nice web middleware in golang.
When providing the `-debug` flag, the `Logger()` middleware will be added.

```bash
$ ./doxy_darwin -debug
2017/08/18 11:44:50 [II] Start Version: 0.1.0
2017/08/18 11:44:50 Error reading patterns file (open /etc/doxy.pattern: no such file or directory), using default patterns
2017/08/18 11:44:50 0  : ^/(v\d\.\d+/)?containers(/\w+)?/json$
2017/08/18 11:44:50 1  : ^/(v\d\.\d+/)?services(/[0-9a-f]+)?$
2017/08/18 11:44:50 2  : ^/(v\d\.\d+/)?tasks(/\w+)?$
2017/08/18 11:44:50 3  : ^/(v\d\.\d+/)?networks(/\w+)?$
2017/08/18 11:44:50 4  : ^/(v\d\.\d+/)?nodes(/\w+)?$
2017/08/18 11:44:50 5  : ^/(v\d\.\d+/)?info$
2017/08/18 11:44:50 6  : ^/_ping$
2017/08/18 11:44:50 [gk-soxy] Listening on /tmp/doxy.sock
[negroni] 2017-08-18T11:45:00+02:00 | 200 | 	 3.800713ms | docker | GET /_ping
[negroni] 2017-08-18T11:45:00+02:00 | 403 | 	 34.067µs | docker | GET /v1.31/containers/a62250e0890a/export
[negroni] 2017-08-18T11:45:04+02:00 | 200 | 	 1.800044ms | docker | GET /_ping
[negroni] 2017-08-18T11:45:04+02:00 | 200 | 	 2.055015ms | docker | GET /v1.31/containers/json
```