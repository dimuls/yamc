# yamc
Yet another memory cache.

Features:
* supports string values, lists and dictionaries
* each key has time to live (TTL)
* HTTP restful API
* has go client
* fully tested


## Installation

```bash
$ go get github.com/someanon/yamc
```

## Testing

```bash
$ go get -t github.com/someanon/yamc
$ go test ./...
```

## Running

```bash
$ yamc
```

Runs on port 8080. No root privileges required.