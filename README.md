# PocketBase with PostgresSQL

**Features**

- ✅ PostgresSQL support
- ✅ Support horizontal scaling
- ✅ Realtime events works normally with horizontal scaling
- ✅ 100% test case pass rate across total 4701 unit tests
- ✅ Fully compatible with latest PocketBase SDKs/Docs

**Demo App**

To demonstrate the horizontal scaling and realtime capabilities, I have deployed a **realtime chat demo app** on two different PocketBase instances.

- Instance 1: [pocketbase-chat-01.mimimiao.com](https://pocketbase-chat-01.mimimiao.com)
- Instance 2: [pocketbase-chat-02.mimimiao.com](https://pocketbase-chat-02.mimimiao.com)

![Realtime Chat Demo](./.github/chat-page-demo.png)

> [!TIP]
> The above demo app is free hosted on ClawCloud Run and Neon DB, see [Free Hosting Guide](https://github.com/fondoger/pocketbase/tree/realtime-demo) for more details.

**Get Started**

Everything is the same as the original PocketBase, except that an additional environment variable `POSTGRES_URL` is required.

```sh
export POSTGRES_URL=postgres://user:pass@127.0.0.1:5432/postgres?sslmode=disable
./pocketbase serve
```

See: [pocketbase/pocketbase](https://github.com/pocketbase/pocketbase)

**Deploy with Docker**

1. Start PostgresSQL:

   ```sh
   docker run -d \
     --name postgres \
     -p 5432:5432 \
     -e POSTGRES_USER=user \
     -e POSTGRES_PASSWORD=pass \
     -v ./pg_data:/var/lib/postgresql/data \
     postgres:alpine
   ```

2. Start PocketBase (Don't use @latest tag, use a specific version in production)

   ```sh
   docker run -d \
     --network=host \
     --name pocketbase \
     -v ./pb_data:/data \
     -e PB_HTTP_ADDR=127.0.0.1:8090 \
     -e PB_DATA_DIR="/data" \
     -e POSTGRES_URL="postgres://user:pass@127.0.0.1:5432/postgres?sslmode=disable" \
     ghcr.io/fondoger/pocketbase:latest
   ```

3. Get admin password reset link
   ```sh
   docker logs -f pocketbase
   ```

**Available Environment Variables**
| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `POSTGRES_URL` | PostgresSQL connection URL (required) | required |
| `POSTGRES_DATA_DB` | Database name for the PocketBase data | `pb-data` |
| `POSTGRES_AUX_DB` | Database name for the PocketBase logs data | `pb-auxiliary` |
| `PB_DATA_DIR` | Directory to store the PocketBase data | `./pb_data` |
| `PB_PUBLIC_DIR` | Directory to store the public files | `./pb_public` |
| `PB_HOOKS_DIR` | Directory to store the custom hooks | `./pb_hooks` |
| `PB_REALTIME_BRIDGE` | Enable/Disable the realtime bridge. Disable it if you don't need horizontal scaling or don't need realtime feature. | `true` |
| `PB_HTTP_ADDR` | TCP address to listen for the HTTP server | `127.0.0.1:8090` if no domain specified |
| `PB_HTTPS_ADDR` | TCP address to listen for the HTTPS server | - |
| `PB_PATH_PREFIX` | URL path prefix for the HTTP server (Useful when reuse same domain for diffrent sites behind nginx) | - |
| `PB_ALLOWED_ORIGINS` | Comma separated list of allowed CORS origins | `*` (all origins) |

**Limitations**

- Local file system is not synced across multiple instances.
  > You need to add a S3 storage account if you are deploying multiple instances and need the file upload feature.
- The built-in SQLite Backup feature is not supported
  > PostgresSQL have many mature and stable backup solutions. Eg: `pg_dump`, `docker-pg-backup`, `postgres-backup-s3`.

**Disclaimer**

I don't want to create another fork of PocketBase, but I really need PostgresSQL support for my project. The eventual goal is to merge this code back into the main PocketBase repository.

If you want to contribute, please first go to [pocketbase/pocketbase](https://github.com/pocketbase/pocketbase). This repository only handles the bugs related to PostgresSQL.

---

## Original README.md file

<p align="center">
    <a href="https://pocketbase.io" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/5qimnm5.png" alt="PocketBase - open source backend in 1 file" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml" target="_blank" rel="noopener"><img src="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml/badge.svg" alt="build" /></a>
    <a href="https://github.com/pocketbase/pocketbase/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/release/pocketbase/pocketbase.svg" alt="Latest releases" /></a>
    <a href="https://pkg.go.dev/github.com/pocketbase/pocketbase" target="_blank" rel="noopener"><img src="https://godoc.org/github.com/pocketbase/pocketbase?status.svg" alt="Go package documentation" /></a>
</p>

[PocketBase](https://pocketbase.io) is an open source Go backend that includes:

- embedded database (_SQLite_) with **realtime subscriptions**
- built-in **files and users management**
- convenient **Admin dashboard UI**
- and simple **REST-ish API**

**For documentation and examples, please visit https://pocketbase.io/docs.**

> [!WARNING]
> Please keep in mind that PocketBase is still under active development
> and therefore full backward compatibility is not guaranteed before reaching v1.0.0.

## API SDK clients

The easiest way to interact with the PocketBase Web APIs is to use one of the official SDK clients:

- **JavaScript - [pocketbase/js-sdk](https://github.com/pocketbase/js-sdk)** (_Browser, Node.js, React Native_)
- **Dart - [pocketbase/dart-sdk](https://github.com/pocketbase/dart-sdk)** (_Web, Mobile, Desktop, CLI_)

You could also check the recommendations in https://pocketbase.io/docs/how-to-use/.


## Overview

### Use as standalone app

You could download the prebuilt executable for your platform from the [Releases page](https://github.com/pocketbase/pocketbase/releases).
Once downloaded, extract the archive and run `./pocketbase serve` in the extracted directory.

The prebuilt executables are based on the [`examples/base/main.go` file](https://github.com/pocketbase/pocketbase/blob/master/examples/base/main.go) and comes with the JS VM plugin enabled by default which allows to extend PocketBase with JavaScript (_for more details please refer to [Extend with JavaScript](https://pocketbase.io/docs/js-overview/)_).

### Use as a Go framework/toolkit

PocketBase is distributed as a regular Go library package which allows you to build
your own custom app specific business logic and still have a single portable executable at the end.

Here is a minimal example:

0. [Install Go 1.23+](https://go.dev/doc/install) (_if you haven't already_)

1. Create a new project directory with the following `main.go` file inside it:
    ```go
    package main

    import (
        "log"

        "github.com/pocketbase/pocketbase"
        "github.com/pocketbase/pocketbase/core"
    )

    func main() {
        app := pocketbase.New()

        app.OnServe().BindFunc(func(se *core.ServeEvent) error {
            // registers new "GET /hello" route
            se.Router.GET("/hello", func(re *core.RequestEvent) error {
                return re.String(200, "Hello world!")
            })

            return se.Next()
        })

        if err := app.Start(); err != nil {
            log.Fatal(err)
        }
    }
    ```

2. To init the dependencies, run `go mod init myapp && go mod tidy`.

3. To start the application, run `go run main.go serve`.

4. To build a statically linked executable, you can run `CGO_ENABLED=0 go build` and then start the created executable with `./myapp serve`.

_For more details please refer to [Extend with Go](https://pocketbase.io/docs/go-overview/)._

### Building and running the repo main.go example

To build the minimal standalone executable, like the prebuilt ones in the releases page, you can simply run `go build` inside the `examples/base` directory:

0. [Install Go 1.23+](https://go.dev/doc/install) (_if you haven't already_)
1. Clone/download the repo
2. Navigate to `examples/base`
3. Run `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build`
   (_https://go.dev/doc/install/source#environment_)
4. Start the created executable by running `./base serve`.

Note that the supported build targets by the pure Go SQLite driver at the moment are:

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
linux   s390x
windows amd64
windows arm64
```

### Testing

PocketBase comes with mixed bag of unit and integration tests.
To run them, use the standard `go test` command:

```sh
go test ./...
```

Check also the [Testing guide](http://pocketbase.io/docs/testing) to learn how to write your own custom application tests.

## Security

If you discover a security vulnerability within PocketBase, please send an e-mail to **support at pocketbase.io**.

All reports will be promptly addressed and you'll be credited in the fix release notes.

## Contributing

PocketBase is free and open source project licensed under the [MIT License](LICENSE.md).
You are free to do whatever you want with it, even offering it as a paid service.

You could help continuing its development by:

- [Contribute to the source code](CONTRIBUTING.md)
- [Suggest new features and report issues](https://github.com/pocketbase/pocketbase/issues)

PRs for new OAuth2 providers, bug fixes, code optimizations and documentation improvements are more than welcome.

But please refrain creating PRs for _new features_ without previously discussing the implementation details.
PocketBase has a [roadmap](https://github.com/orgs/pocketbase/projects/2) and I try to work on issues in specific order and such PRs often come in out of nowhere and skew all initial planning with tedious back-and-forth communication.

Don't get upset if I close your PR, even if it is well executed and tested. This doesn't mean that it will never be merged.
Later we can always refer to it and/or take pieces of your implementation when the time comes to work on the issue (don't worry you'll be credited in the release notes).
