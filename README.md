Redis Inventory
===============

[![Build Status](https://travis-ci.com/obukhov/redis-inventory.svg?branch=master)](https://travis-ci.com/obukhov/redis-inventory)
[![Coverage Status](https://coveralls.io/repos/github/obukhov/redis-inventory/badge.svg?branch=master)](https://coveralls.io/github/obukhov/redis-inventory?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/obukhov/redis-inventory)](https://goreportcard.com/report/github.com/obukhov/redis-inventory)
[![Docker Pulls](https://img.shields.io/docker/pulls/dclg/redis-inventory)](https://hub.docker.com/repository/docker/dclg/redis-inventory)

Tool to see redis memory usage by keys in hierarchical way.

Example:
```bash
$ go run main.go inventory localhost:63795 --output=table --output-params="padSpaces=2&depth=2&human=1"                                                                                                                                                                                       643ms  Do 22 Jul 2021 22:01:41 UTC
```

Outputs it as a nice table
```bash
12:39PM INF Start scanning
+---------------------+----------+-----------+
| KEY                 | BYTESIZE | KEYSCOUNT |
+---------------------+----------+-----------+
|   dev:              |     2.9M |     4,555 |
|     article:        |   413.7K |       616 |
|     blogpost:       |   408.5K |       630 |
|     collections:    |   426.7K |       627 |
|     events:         |   391.2K |       614 |
|     friends:foobar: |   501.1K |       745 |
|     news:           |   388.8K |       593 |
|     user:           |     481K |       730 |
|   prod:             |     2.9M |     4,531 |
|     article:        |   397.1K |       614 |
|     blogpost:       |   409.4K |       627 |
|     collections:    |   374.7K |       560 |
|     events:         |   384.2K |       588 |
|     friends:foobar: |     503K |       755 |
|     news:           |   407.9K |       618 |
|     user:           |   492.3K |       769 |
+---------------------+----------+-----------+
12:39PM INF Finish scanning
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
| human        | Display numbers in human-friendly way (0 or 1) | 0         |



If padding is not specified in either way, nested keys are displayed with full paths as following:
```bash
+----------------------+----------+-----------+
| KEY                  | BYTESIZE | KEYSCOUNT |
+----------------------+----------+-----------+
| dev:                 |     2.9M |     4,555 |
| dev:article:         |   413.7K |       616 |
| dev:blogpost:        |   408.5K |       630 |
| dev:collections:     |   426.7K |       627 |
| dev:events:          |   391.2K |       614 |
| dev:friends:foobar:  |   501.1K |       745 |
| dev:news:            |   388.8K |       593 |
| dev:user:            |     481K |       730 |
...
```

### Json

| Option name  | Description                                  | Default   |
|--------------|----------------------------------------------|-----------|
| padSpaces    | Number of spaces to indent the nested level  | `0`       |
| padding      | Use custom character to pad nested level     | `""`      |

If padding is not specified in either way json is not pretty-printed.
