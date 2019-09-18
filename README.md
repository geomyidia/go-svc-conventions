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

If you'd like to pull the deps into `vendor` dir:

```shell
$ make deps
```