<p align="center">
    <a href="https://pocketbase.io" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/5qimnm5.png" alt="PocketBase - open source backend in 1 file" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml" target="_blank" rel="noopener"><img src="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml/badge.svg" alt="build" /></a>
    <a href="https://github.com/pocketbase/pocketbase/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/release/pocketbase/pocketbase.svg" alt="Latest releases" /></a>
    <a href="https://pkg.go.dev/github.com/pocketbase/pocketbase" target="_blank" rel="noopener"><img src="https://godoc.org/github.com/ganigeorgiev/fexpr?status.svg" alt="Go package documentation" /></a>
</p>

[PocketBase](https://pocketbase.io) is an open source Go backend, consisting of:

- embedded database (_SQLite_) with **realtime subscriptions**
- built-in **files and users management**
- convenient **Admin dashboard UI**
- and simple **REST-ish API**

**For documentation and examples, please visit https://pocketbase.io/docs.**

> ⚠️ Please keep in mind that PocketBase is still under active development
> and therefore full backward compatibility is not guaranteed before reaching v1.0.0.


## API SDK clients

The easiest way to interact with the API is to use one of the official SDK clients:

- **JavaScript - [pocketbase/js-sdk](https://github.com/pocketbase/js-sdk)** (_browser and node_)
- **Dart - [pocketbase/dart-sdk](https://github.com/pocketbase/dart-sdk)** (_web, mobile, desktop_)


## Overview

PocketBase could be [downloaded directly as a standalone app](https://github.com/pocketbase/pocketbase/releases) or it could be used as a Go framework/toolkit which allows you to build
your own custom app specific business logic and still have a single portable executable at the end.

### Installation

```sh
# go 1.18+
go get github.com/pocketbase/pocketbase
```
> For Windows, you may have to use go 1.19+ due to an incorrect js mime type in the Windows Registry (see [issue#6](https://github.com/pocketbase/pocketbase/issues/6)).

### Example

```go
package main

import (
    "log"
    "net/http"

    "github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        // add new "GET /hello" route to the app router (echo)
        e.Router.AddRoute(echo.Route{
            Method: http.MethodGet,
            Path:   "/hello",
            Handler: func(c echo.Context) error {
                return c.String(200, "Hello world!")
            },
            Middlewares: []echo.MiddlewareFunc{
                apis.ActivityLogger(app),
            },
        })

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Running and building

Running/building the application is the same as for any other Go program, aka. just `go run` and `go build`.

**PocketBase embeds SQLite, but doesn't require CGO.**

If CGO is enabled (aka. `CGO_ENABLED=1`), it will use [mattn/go-sqlite3](https://pkg.go.dev/github.com/mattn/go-sqlite3) driver, otherwise - [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite).
Enable CGO only if you really need to squeeze the read/write query performance at the expense of complicating cross compilation.

To build the minimal standalone executable, like the prebuilt ones in the releases page, you can simply run `go build` inside the `examples/base` directory:

0. [Install Go 1.18+](https://go.dev/doc/install) (_if you haven't already_)
1. Clone/download the repo
2. Navigate to `examples/base`
3. Run `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build`
    (_https://go.dev/doc/install/source#environment_)
4. Start the created executable by running `./base serve`.

The supported build targets by the non-cgo driver at the moment are:
```
darwin  amd64
darwin  arm64
freebsd amd64
freebsd arm64
linux   386
linux   amd64
linux   arm
linux   arm64
linux   ppc64le
linux   riscv64
windows amd64
windows arm64
```

### Testing

PocketBase comes with mixed bag of unit and integration tests.
To run them, use the default `go test` command:
```sh
go test ./...
```

Check also the [Testing guide](http://pocketbase.io/docs/testing) to learn how to write your own custom application tests.

## Security

If you discover a security vulnerability within PocketBase, please send an e-mail to **support at pocketbase.io**.

All reports will be promptly addressed, and you'll be credited accordingly.


## Contributing

PocketBase is free and open source project licensed under the [MIT License](LICENSE.md).
You are free to do whatever you want with it, even offering it as a paid service.

You could help continuing its development by:

- [Contribute to the source code](CONTRIBUTING.md)
- [Suggest new features and report issues](https://github.com/pocketbase/pocketbase/issues)
- [Donate a small amount](https://pocketbase.io/support-us)

PRs for _small features_ (eg. adding new OAuth2 providers), bug and documentation fixes, etc. are more than welcome.

But please refrain creating PRs for _big features_ without previously discussing the implementation details. Reviewing big PRs often requires a lot of time and tedious back-and-forth communication.
PocketBase has a [roadmap](https://github.com/orgs/pocketbase/projects/2)
and I try to work on issues in a specific order and such PRs often come in out of nowhere and skew all initial planning.

Don't get upset if I close your PR, even if it is well executed and tested. This doesn't mean that it will never be merged.
Later we can always refer to it and/or take pieces of your implementation when the time comes to work on the issue (don't worry you'll be credited in the release notes).

_Please also note that PocketBase was initially created to serve as a new backend for my other open source project - [Presentator](https://presentator.io) (see [#183](https://github.com/presentator/presentator/issues/183)),
so all feature requests will be first aligned with what we need for Presentator v3._
