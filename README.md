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

## The `components` Subpackage

This is obviosuly not a standard practice in Go. To be clear, I am not advocating for such a convention. Rather, the use of the `components` subpackage is meant to underscore the packages that:

1. are central to the app, while also
1. being code that either needs to talk to other components or vice versa

To the last point, (and regardless of what they are actually called or where they live) most "components" should offer an API and a state struct (e.g., connection info) that may be utilised by other "components" in the system.

## Build and Run

```shell
$ export PATH=$PATH:~/go/bin
$ make
$ make run
```

```
2021-07-23T14:04:27-05:00 INFO ▶ Logging initialized.
[watermill] 2021/07/23 14:04:27.291402 router.go:183: 	level=DEBUG msg="Adding plugins" count=1
[watermill] 2021/07/23 14:04:27.291430 router.go:151: 	level=DEBUG msg="Adding middleware" count=1
2021-07-23T14:04:27-05:00 DEBUG ▶ Setting up database connection ...
2021-07-23T14:04:27-05:00 DEBUG ▶ Setting up HTTP daemon ...
2021-07-23T14:04:27-05:00 DEBUG ▶ Setting up HTTPD routes ...
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /echo                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPServer).Echo-fm (1 handlers)
[GIN-debug] GET    /health                   --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPServer).Health-fm (1 handlers)
[GIN-debug] GET    /ping                     --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPServer).Ping-fm (1 handlers)
[GIN-debug] GET    /version                  --> github.com/geomyidia/go-svc-conventions/pkg/components/httpd.(*HTTPServer).Version-fm (1 handlers)
2021-07-23T14:04:27-05:00 DEBUG ▶ HTTP daemon set up.
2021-07-23T14:04:27-05:00 DEBUG ▶ Setting up gRPC daemon ...
2021-07-23T14:04:27-05:00 DEBUG ▶ gRPC implementation set up.
2021-07-23T14:04:27-05:00 INFO ▶ gRPC daemon listening on localhost:2525 ...
2021-07-23T14:04:27-05:00 INFO ▶ HTTP daemon listening on localhost:5099 ...
2021-07-23T14:04:27-05:00 DEBUG ▶ Setting up event bus auditor ...
2021-07-23T14:04:27-05:00 INFO ▶ Auditor is listening for new events ...
badger 2021/07/23 14:04:27 INFO: All 0 tables opened in 0s
badger 2021/07/23 14:04:27 INFO: Discard stats nextEmptySlot: 0
badger 2021/07/23 14:04:27 INFO: Set nextTxnTs to 0
badger 2021/07/23 14:04:27 INFO: Deleting empty file: ./data/badger/000031.vlog
2021-07-23T14:04:27-05:00 INFO ▶ Connected to database: ./data/badger
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
Version: 0.3.0-dev
Build Date: 2021-07-23T19:19:33Z
Git Commit: 1e9402b
Git Branch: release/0.3.x
Git Summary: 0.2.0-5-g1e9402b-dirty
```

You should see server debug notices for these requests:

```
2021-07-23T14:19:40-05:00 DEBUG ▶ Received HTTP ping request
2021-07-23T14:19:40-05:00 DEBUG ▶ Publishing to topic '*' ...
2021-07-23T14:19:40-05:00 WARNING ▶ Auditor got event: &{ID:e441545a-2888-46de-975b-b518577c99e2 Name:ping Data:DATA}
2021-07-23T14:19:44-05:00 DEBUG ▶ Received HTTP health request
2021-07-23T14:19:47-05:00 DEBUG ▶ Received HTTP echo request: Stuff
2021-07-23T14:19:49-05:00 DEBUG ▶ Received HTTP version request
```

Checking the gRPC daemon with our client:

```shell
$ ./bin/client ping
PONG
$ ./bin/client health
Services: OK
Errors: NULL
$ ./bin/client echo hej, check this stuff out
Echo: [hej, check this stuff out]
$ ./bin/client version
Version: 0.3.0-dev
BuildDate: 2021-07-23T19:19:33Z
GitCommit: 1e9402b
GitBranch: release/0.3.x
GitSummary: 0.2.0-5-g1e9402b-dirty
```

And resulting output:

```
2021-07-23T14:20:22-05:00 DEBUG ▶ Received gRPC ping request
2021-07-23T14:20:22-05:00 DEBUG ▶ Publishing to topic '*' ...
2021-07-23T14:20:22-05:00 WARNING ▶ Auditor got event: &{ID:0f89bfeb-0788-4fac-bce6-8ff9fb0cdd80 Name:ping Data:DATA}
2021-07-23T14:20:29-05:00 DEBUG ▶ Received gRPC health request
2021-07-23T14:20:36-05:00 DEBUG ▶ Received gRPC echo request: data:"[hej, check this stuff out]"
2021-07-23T14:20:43-05:00 DEBUG ▶ Received gRPC version request
2021-07-23T14:20:43-05:00 DEBUG ▶ Publishing to topic '*' ...
2021-07-23T14:20:43-05:00 WARNING ▶ Auditor got event: &{ID:c037eda9-696d-45c0-8507-fd570becc9ee Name:version Data:DATA}
```
