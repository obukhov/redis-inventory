Redis Inventory
===============

[![Build Status](https://travis-ci.com/obukhov/redis-inventory.svg?branch=master)](https://travis-ci.com/obukhov/redis-inventory)
[![Coverage Status](https://coveralls.io/repos/github/obukhov/redis-inventory/badge.svg?branch=master)](https://coveralls.io/github/obukhov/redis-inventory?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/obukhov/redis-inventory)](https://goreportcard.com/report/github.com/obukhov/redis-inventory)
[![Docker Pulls](https://img.shields.io/docker/pulls/dclg/redis-inventory)](https://hub.docker.com/repository/docker/dclg/redis-inventory)

Redis inventory is a tool to analyse redis memory usage by keys patterns and displaying it in hierarchically. The name
is inspired by "Disk Inventory X" tool doing similar analysis for disk usage.

Example:

```bash
$ redis-inventory inventory localhost:63795 --output=table --output-params="padSpaces=2&depth=2&human=1"                                                                                                                                                                                       643ms  Do 22 Jul 2021 22:01:41 UTC
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

Read more about [usage](docs/usage.md)

## Installation

There are several ways to install this tools:

- using docker
- building from sources

### Using docker

To run the tool from a docker image, run the command:

```bash
docker run --rm dclg/redis-inventory inventory <HOST>:<PORT>
```

If you plan to run sequence of `index` and `display` so want to utilize file cache, add volume mount to local dir.

```bash
docker run --rm -v "$PWD:/tmp" dclg/redis-inventory index <HOST>:<PORT>
```

If you use the tool regularly, create an alias:

```bash
alias ri="docker run --rm -v "$PWD:/tmp" dclg/redis-inventory"
```

So you can run it just like so:

```bash
$ ri index <host>:<port>
11:52PM INF Start indexing
Scanning keys ... done! [3.18K in 262ms]
11:52PM INF Finish scanning and saved index as a file /tmp/redis-inventory.json

$ ri display --output-params="depth=2&padSpaces=2"
11:53PM INF Loading index
+---------------------+----------+-----------+
| KEY                 | BYTESIZE | KEYSCOUNT |
+---------------------+----------+-----------+
|   dev:              |  1053120 |      1587 |
|     article:        |   142824 |       210 |
|     blogpost:       |   153416 |       240 |
|     collections:    |   145912 |       217 |
|     events:         |   144360 |       217 |
|     friends:foobar: |   169040 |       249 |
|     news:           |   124944 |       191 |
|     user:           |   172624 |       263 |
|   prod:             |  1078048 |      1589 |
|     article:        |   157264 |       231 |
|     blogpost:       |   145632 |       214 |
|     collections:    |   165992 |       240 |
|     events:         |   143536 |       219 |
|     friends:foobar: |   168864 |       251 |
|     news:           |   132096 |       197 |
|     user:           |   164664 |       237 |
+---------------------+----------+-----------+
11:53PM INF Done

$ ri display --output-params="depth=1&human=1&padding=🔥"
11:57PM INF Loading index
+---------+----------+-----------+
| KEY     | BYTESIZE | KEYSCOUNT |
+---------+----------+-----------+
| 🔥dev:  |       1M |     1,587 |
| 🔥prod: |       1M |     1,589 |
+---------+----------+-----------+
```

### Building from sources

You have to have [golang installed](https://golang.org/doc/install) on your computer.

1. Checkout the repository:

```bash
git clone git@github.com:obukhov/redis-inventory.git
```

2. Build a binary:

```bash
cd redis-inventory
go build -o redis-inventory main.go
```

3. Run it or move it one of your PATH directories:

```bash
mv redis-inventory /usr/local/bin/
```
