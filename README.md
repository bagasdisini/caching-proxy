# Caching Proxy

A CLI-based caching proxy server built in GoLang that forwards requests to an origin server, caches the responses, and returns cached responses for repeated requests. It also allows for cache clearing via a CLI command. **Without using any external libraries**.

This project is served as a solution to the [Caching Proxy by roadmap.sh](https://roadmap.sh/projects/caching-server).

## Features

* Caches responses from the origin server.
* Serves cached responses for repeated requests.
* Adds headers to indicate whether the response is from cache or the origin server.
* Provides an option to clear the cache via CLI.

## Requirements

* [Go 1.23](https://go.dev/dl/) or higher

## Usage

Clone the repository and build the project using the following commands:

```bash
git clone https://github.com/bagasdisini/caching-proxy.git
cd caching-proxy
go build -o caching-proxy
```

### Starting the Caching Proxy Server
To start the server, run the following command:

```bash
./caching-proxy --port <number> --origin <url>
```

* `--port` is the port on which the caching proxy server will run.
* `--origin` is the URL of the server to which the requests will be forwarded.

#### Example

```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

This will start the proxy server on port `3000` and forward requests to `http://dummyjson.com`.

#### Example Request

If you run the server with the example above, and then make a request to `http://localhost:3000/products`, the proxy will forward the request to `http://dummyjson.com/products`, return the response, and cache it. On subsequent requests to the same URL, the cached response will be served.

### Cache Headers

* `X-Cache: HIT` — Indicates the response was served from cache.
* `X-Cache: MISS` — Indicates the response was forwarded to the origin server.

### Clearing the Cache

To clear the cache, use the following command:

```bash
./caching-proxy --clear-cache
```

This will remove all cached responses, allowing fresh requests to be forwarded to the origin server.

