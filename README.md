Redis Inventory [beta]
=====================

[![Build Status](https://travis-ci.com/obukhov/redis-inventory.svg?branch=master)](https://travis-ci.com/obukhov/redis-inventory)

[![Coverage Status](https://coveralls.io/repos/github/obukhov/redis-inventory/badge.svg?branch=master)](https://coveralls.io/github/obukhov/redis-inventory?branch=master)

[![Go Report Card](https://goreportcard.com/badge/github.com/obukhov/redis-inventory)](https://goreportcard.com/report/github.com/obukhov/redis-inventory)

Tool to see redis memory usage by keys in hierarchical way.

Example:
```bash
$ go run main.go inventory localhost:63795 --output=table --output-params="padSpaces=2&depth=2"                                                                                                                                                                                       643ms  Do 22 Jul 2021 22:01:41 UTC
```

Outputs it as a nice table
```bash
12:04AM INF Start scanning
+---------------------+----------+-------+
| KEY                 | BYTESIZE | COUNT |
+---------------------+----------+-------+
|   prod:             | 1902856  |       |
|     friends:foobar: | 333552   |       |
|     events:         | 233768   |       |
|     user:           | 290544   |       |
|     collections:    | 253984   |       |
|     news:           | 269800   |       |
|     article:        | 267568   |       |
|     blogpost:       | 253640   |       |
|   dev:              | 2016720  |       |
|     news:           | 270688   |       |
|     article:        | 256824   |       |
|     friends:foobar: | 309440   |       |
|     blogpost:       | 275392   |       |
|     user:           | 344976   |       |
|     events:         | 290728   |       |
|     collections:    | 268672   |       |
+---------------------+----------+-------+
12:04AM INF Finish scanning
```

Not all the features are implemented, for details see the [project](https://github.com/obukhov/redis-inventory/projects/1)

## General interface

```bash
go run main.go inventory <host>:<port> [--output=<output type>] [--output-params=<querstring serialized params>]
```

## Output type

### Table

| Option name  | Description                                    | Default   | 
|--------------|------------------------------------------------|-----------|
| padSpaces    | Number of spaces to indent the nested level    | `0`       |
| padding      | Use custom character to pad nested level       | `""`      |
| depth        | Maximum nesting level for keys before grouping | 10        |


If padding is not specified in either way, nested keys are displayed with full paths as following:
```bash
+----------------------+----------+-------+
| KEY                  | BYTESIZE | COUNT |
+----------------------+----------+-------+
| prod:                | 1889136  |       |
| prod:user:           | 287776   |       |
| prod:collections:    | 251624   |       |
| prod:news:           | 268496   |       |
| prod:article:        | 264560   |       |
...
```

### Json

| Option name  | Description                                  | Default   | 
|--------------|----------------------------------------------|-----------|
| padSpaces    | Number of spaces to indent the nested level  | `0`       |
| padding      | Use custom character to pad nested level     | `""`      |

If padding is not specified in either way json is not pretty-printed.