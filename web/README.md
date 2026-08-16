# Smolpho Web

Smolpho Web is the read-only React interface for the Smolpho indexer API. It
shows market totals and the published, paginated list of indexed positions.

## Development

Start the Smolpho API from the repository root:

```sh
task demo:api
```

Then start the web application:

```sh
cd web
npm install
npm run dev
```

Vite serves the application at `http://localhost:5173` and proxies `/api` and
`/healthz` to `http://127.0.0.1:8080`. Override the proxy destination with
`VITE_API_PROXY_TARGET`.

For a separately hosted API, set `VITE_API_URL` to its full origin. That origin
must allow the web application's origin through the API's `-cors-origin` flag.

## Build

```sh
npm run build
```

The production bundle is written to `web/dist`.

Accounting values are parsed as `bigint` and displayed as coefficients of
`1e18`. The API does not currently expose token symbols.
