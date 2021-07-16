# go-svc-conventions

## Key Concepts

This repo demonstrates the following:

* Areas of responsibility separated into application "components" (config, logging, db connections, various servers, etc.)
* Wait groups and contexts for graceful shutdowns

In particular, what are shown in this code:

* Encapsulation of functionality you want to assemble code using "components"
* Embedding of this in a struct as a "server" (the term Mat Ryer uses) or as an "application" (the term I use in this repo)
* Providing both reuse and symmetry between server and client code
* Putting in place mechanisms that facilitate lower-effort, lower-impact future refactorings
* Providing symmetry between HTTP and gRPC handlers (note that business logic should be done elsewhere! then called via imported functions inside your handlers)
* Pulling of config into memory

Much of this has been taken from my long experiences in the world of non-Go fault-tolerant and highly available services, but is validated by long-time Go hackers with very similar and pragmatic views.

Here are some links of interest:

* [https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)
* [https://www.youtube.com/watch?v=rWBSMsLG8po](https://www.youtube.com/watch?v=rWBSMsLG8po)

## Build and Run

```shell
$ make
$ make run
```

```
2021-07-16T01:32:12-05:00 INFO ▶ Logging initialized.
2021-07-16T01:32:12-05:00 DEBUG ▶ Setting up HTTP daemon ...
2021-07-16T01:32:12-05:00 DEBUG ▶ Setting up HTTPD routes ...
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /echo                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Echo-fm (1 handlers)
[GIN-debug] GET    /health                   --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Health-fm (1 handlers)
[GIN-debug] GET    /ping                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Ping-fm (1 handlers)
[GIN-debug] GET    /version                  --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Version-fm (1 handlers)
2021-07-16T01:32:12-05:00 DEBUG ▶ HTTP daemon set up.
2021-07-16T01:32:12-05:00 DEBUG ▶ Setting up gRPC daemon ...
2021-07-16T01:32:12-05:00 DEBUG ▶ gRPC implementation set up.
2021-07-16T01:32:12-05:00 INFO ▶ gRPC daemon listening on localhost:2525 ...
2021-07-16T01:32:12-05:00 INFO ▶ HTTP daemon listening on localhost:5099 ...
```

Or, to run with a different HTTP port:

```shell
$ APP_HTTPD_PORT=5151 make run 
```

Using the default port, check the HTTP daemon with `curl`:

```shell
$ curl http://localhost:5099/ping
pong
$ curl http://localhost:5099/health
Services: OK
Errors: NULL
$ curl -XPOST http://localhost:5099/echo -d "Stuff"
Stuff
$ curl http://localhost:5099/version
Version: 0.2.0
Build Date: 2021-07-16T06:32:10Z
Git Commit: ba7858a
Git Branch: main
Git Summary: 0.1.0-11-gba7858a-dirty
```

You should see server debug notices for these requests:

```
2021-07-16T01:32:48-05:00 DEBUG ▶ Got ping request
2021-07-16T01:32:55-05:00 DEBUG ▶ Got health request
2021-07-16T01:33:05-05:00 DEBUG ▶ Got echo request: Stuff
2021-07-16T01:33:16-05:00 DEBUG ▶ Got version request
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
$ ./bin/client ping
```
```
2019-09-19T14:27:58-05:00 INFO [github.com/geomyidia/zylog/logger.SetupLogging:109] ▶ Logging initialized.
2019-09-19T14:27:58-05:00 INFO [main.main:59] ▶ Reply: PONG
```
