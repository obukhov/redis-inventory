Redis Inventory
===============

[![Build Status](https://travis-ci.com/obukhov/redis-inventory.svg?branch=master)](https://travis-ci.com/obukhov/redis-inventory)
[![Coverage Status](https://coveralls.io/repos/github/obukhov/redis-inventory/badge.svg?branch=master)](https://coveralls.io/github/obukhov/redis-inventory?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/obukhov/redis-inventory)](https://goreportcard.com/report/github.com/obukhov/redis-inventory)
[![Docker Pulls](https://img.shields.io/docker/pulls/dclg/redis-inventory)](https://hub.docker.com/repository/docker/dclg/redis-inventory)

Redis inventory is a tool to analyse Redis memory usage by key patterns and displaying it hierarchically. The name is
inspired by "Disk Inventory X" tool doing similar analysis for disk usage.

Blog post on how it works: [Medium](https://medium.com/geekculture/analyzing-memory-usage-in-redis-by-key-patterns-e36420a2c225) | [My blog](https://dclg.net/posts/2021-08-16---Redis-inventory/redis-inventory)

Example:

```bash
$ redis-inventory inventory <redis-url> --output=table --output-params="padSpaces=2&depth=2&human=1"
```

`<redis-url>` can be provided in [one of two formats](docs/redis-url.md)

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

It also can render sunburst diagrams to visualize it:
```bash
$ redis-inventory inventory <redis-url> --output=chart --maxChildren=100
```


[![Full sunburst diagram](docs/images/diagram1000.png)](docs/images/diagram.png)

Read more about [usage](docs/usage.md)

## Installation

There are two ways to install the tool:

- use docker image
- build binary from sources

### Using docker

To run the tool from a docker image, run the command:

```bash
docker run --rm dclg/redis-inventory inventory <HOST>:<PORT>
```

Read more about [installation](docs/installation.md)

### Kudos

- [Radix](https://github.com/mediocregopher/radix) — Redis client for Go,
- [Cobra](https://github.com/spf13/cobra) — CLI applications framework,
- [Go-Pretty](https://github.com/jedib0t/go-pretty) — tools to render tables and progress bars,
- [AnyChart](https://www.anychart.com/) — beautiful sunburst diagrams.