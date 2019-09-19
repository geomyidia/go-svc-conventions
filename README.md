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

Checking the HTTP daemon with `curl`:

```shell
$ curl http://localhost:1515/rest/ping
```
```
pong
```

Checking the gRPC daemon with our client:

```shell
$ ./bin/client echo hej, check this stuff out
```
```
2019-09-19T14:27:27-05:00 INFO [github.com/geomyidia/zylog/logger.SetupLogging:109] ▶ Logging initialized.
2019-09-19T14:27:27-05:00 INFO [main.main:46] ▶ Echo: [hej, check this stuff out]
```
```shell
$ ./bin/client health
```
```
2019-09-19T14:27:43-05:00 INFO [github.com/geomyidia/zylog/logger.SetupLogging:109] ▶ Logging initialized.
2019-09-19T14:27:43-05:00 INFO [main.main:52] ▶ Services: OK
2019-09-19T14:27:43-05:00 INFO [main.main:53] ▶ Errors: NULL
```
```shell
$ ./bin/client
```
```
2019-09-19T14:27:58-05:00 INFO [github.com/geomyidia/zylog/logger.SetupLogging:109] ▶ Logging initialized.
2019-09-19T14:27:58-05:00 INFO [main.main:59] ▶ Reply: PONG
```


## Development

If you'd like to pull the deps into `vendor` dir:

```shell
$ make deps
```