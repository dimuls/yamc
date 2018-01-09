# yamc
Yet another memory cache.

Features:
* supports string values, lists and dictionaries
* each key has time to live (TTL)
* HTTP restful API
* has go client
* authorization support
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

## Parameters

### Accounts
To run server `accounts` file required. It must contain YAML encoded map of login and password. File path can be set by `-accounts` flag. Default is `./accounts`

### Cleaning period
It is store cleaning period. Cleaning removes expired keys, lists and dicts. Can be set by `-cleaning-period` flag. Default is `60s`. Must be `time.Duration` string. 

## Running

```bash
$ yamc -accounts ./accounts -cleaning-period 60s
```

Runs on port 8080. No root privileges required.