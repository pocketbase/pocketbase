PocketBase Superuser dashboard UI
======================================================================

This is the PocketBase Superuser dashboard UI (built with Shablon and Vite).

Although it could be used independently, it is intended to be embedded and extended
as part of the PocketBase app executable (hence the `dist` directory and `embed.go` file).

> [!WARNING]
> The UI kit and extension APIs remains deliberately undocumented for the time being until a stable PocketBase release is published ([#7612](https://github.com/pocketbase/pocketbase/discussions/7612)).

## Development

Download the repository and run the appropriate console commands:

```sh
# install dependencies
npm install

# start a dev server with hot reload at localhost:5173
npm run dev

# or generate production ready bundle in dist/ directory
npm run build
```
