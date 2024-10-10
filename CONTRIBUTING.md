# Contributing to PocketBase

Thanks for taking the time to improve PocketBase!

This document describes how to prepare a PR for a change in the main repository.

- [Prerequisites](#prerequisites)
- [Making changes in the Go code](#making-changes-in-the-go-code)
- [Making changes in the Admin UI](#making-changes-in-the-admin-ui)

## Prerequisites

- Go 1.23+ (for making changes in the Go code)
- Node 18+ (for making changes in the Admin UI)

If you haven't already, you can fork the main repository and clone your fork so that you can work locally:

```
git clone https://github.com/your_username/pocketbase.git
```

> [!IMPORTANT]
> It is recommended to create a new branch from master for each of your bugfixes and features.
> This is required if you are planning to submit multiple PRs in order to keep the changes separate for review until they eventually get merged.

## Making changes in the Go code

PocketBase is distributed as a Go package, which means that in order to run the project you'll have to create a Go `main` program that imports the package.

The repository already includes such program, located in `examples/base`, that is also used for the prebuilt executables.

So, let's assume that you already done some changes in the PocketBase Go code and you want now to run them:

1. Navigate to `examples/base`
2. Run `go run main.go serve`

This will start a web server on `http://localhost:8090` with the embedded prebuilt Admin UI from `ui/dist`. And that's it!

**Before making a PR to the main repository, it is a good idea to:**

- Add unit/integration tests for your changes (we are using the standard `testing` go package).
  To run the tests, you could execute (while in the root project directory):

  ```sh
  go test ./...

  # or using the Makefile
  make test
  ```

- Run the linter - **golangci** ([see how to install](https://golangci-lint.run/usage/install/#local-installation)):

  ```sh
  golangci-lint run -c ./golangci.yml ./...

  # or using the Makefile
  make lint
  ```

## Making changes in the Admin UI

PocketBase Admin UI is a single-page application (SPA) built with Svelte and Vite.

To start the Admin UI:

1. Navigate to the `ui` project directory
2. Run `npm install` to install the node dependencies
3. Start vite's dev server
   ```sh
   npm run dev
   ```

You could open the browser and access the running Admin UI at `http://localhost:3000`.

Since the Admin UI is just a client-side application, you need to have the PocketBase backend server also running in the background (either manually running the `examples/base/main.go` or download a prebuilt executable).

> [!NOTE]
> By default, the Admin UI is expecting the backend server to be started at `http://localhost:8090`, but you could change that by creating a new `ui/.env.development.local` file with `PB_BACKEND_URL = YOUR_ADDRESS` variable inside it.

Every change you make in the Admin UI should be automatically reflected in the browser at `http://localhost:3000` without reloading the page.

Once you are done with your changes, you have to build the Admin UI with `npm run build`, so that it can be embedded in the go package. And that's it - you can make your PR to the main PocketBase repository.
