# yamc
Yet another memory cache.

Features:
* supports string values, lists and dictionaries
* each key has time to live (TTL)
* HTTP restful API
* has go client
* authorization support
* dumping/loading to/from file
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

### Accounts path / `--accounts-path`
To run server `accounts` file required. It must contain YAML encoded map of login and password. File path can be set by `--accounts-path` flag. Default is `./accounts`

### Cleaning period / `--cleaning-period`
It is store cleaning period. Cleaning removes expired keys, lists and dicts. Can be set by `--cleaning-period` flag. Default is `60s`. Must be `time.Duration` string and `>= 100ms`.

### Dumping period / `--dumping-period`
It is store dumping period. Dumping dumps store items to file. Can be set by `--dumping-period` flag. Default is `60s`. Must be `time.Duration` string and `>= 60s`. 

### Dump path / `--dump-path`
Path to dump file. If file exists on service start store items are loaded from it. Then each dumping period items are dumped to there. Can be set by `--dump-path` flag. Default is `./dump`.


## Running

```bash
$ yamc -accounts ./accounts -cleaning-period 60s
```

Runs on port 8080. No root privileges required.

## Server documentation

Located [here](https://github.com/someanon/yamc/tree/master/server).
