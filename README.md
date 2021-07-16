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

```
2021-07-16T01:15:58-05:00 INFO [github.com/geomyidia/zylog/logger.SetupLogging:109] ▶ Logging initialized.
2021-07-16T01:15:58-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.SetupServer:11] ▶ Setting up HTTP daemon ...
2021-07-16T01:15:58-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).SetupRoutes:35] ▶ Setting up HTTPD routes ...
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /echo                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Echo-fm (1 handlers)
[GIN-debug] GET    /health                   --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Health-fm (1 handlers)
[GIN-debug] GET    /ping                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Ping-fm (1 handlers)
2021-07-16T01:15:58-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.SetupServer:13] ▶ HTTP daemon set up.
2021-07-16T01:15:58-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/grpcd.SetupServer:11] ▶ Setting up gRPC daemon ...
2021-07-16T01:15:58-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/grpcd.SetupServer:13] ▶ gRPC implementation set up.
2021-07-16T01:15:58-05:00 INFO [github.com/geomyidia/go-svc-conventions/pkg/components/grpcd.(*GRPCHandlerServer).Serve:72] ▶ gRPC daemon listening on localhost:2525 ...
2021-07-16T01:15:58-05:00 INFO [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Serve:70] ▶ HTTP daemon listening on localhost:5099 ...
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
```

You should see server debug notices for these requests:

```
2021-07-16T01:18:23-05:00 INFO [github.com/geomyidia/go-svc-conventions/pkg/components/grpcd.(*GRPCHandlerServer).Serve:72] ▶ gRPC daemon listening on localhost:2525 ...
2021-07-16T01:18:23-05:00 INFO [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Serve:70] ▶ HTTP daemon listening on localhost:5099 ...
2021-07-16T01:18:25-05:00 DEBUG [github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPHandlerServer).Echo:52] ▶ Got echo request: Stuff
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

## Development

If you'd like to pull the deps into `vendor` dir:

```shell
$ make deps
```
