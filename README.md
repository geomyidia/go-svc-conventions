# go-svc-conventions

## Key Concetps

This repo demonstrates the following:

* Areas of responsibility separated into application "components" (config, logging, db connections, various servers, etc.)
* Wait groups and contexts for graceful shutdowns

In particular, what are shown in this code:

* Encapsulate the different ways in which you want to assemble your code using "components"
* Embed this in a struct as a "server" (the term Mat Ryer uses) or as an "application" (the term I use in this repo)
* Provide both reuse and symmetry between server and client code
* Put in place mechanisms that facilitate lower-effort, lower-impact future refactorings
* Provide symmetry between HTTP and gRPC handlers (note that business logic should be done elsewhere! then called via imported functions inside your handlers)
* Pull config into memory

Much of this has been taken from my long experiences in the world of non-Go fault-tolerant and highly available services, but is validated by long-time Go hackers with very similar and pragmatic views.

Here are some links of interest:

* [https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)
* [https://www.youtube.com/watch?v=rWBSMsLG8po](https://www.youtube.com/watch?v=rWBSMsLG8po)

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
$ curl http://localhost:5099/ping
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
