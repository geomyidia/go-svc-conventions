# go-svc-conventions

## Build and Run

```shell
$ make
$ make run
```

Or, to run with a different HTTP port:

```shell
$ APP_HTTPD_PORT=5151 make run 
```

Checking with `curl`:

```shell
$ curl http://localhost:1515/rest/ping
```
```
pong
```

## Development

If you'd like to pull the deps into `vendor` dir:

```shell
$ make deps
```