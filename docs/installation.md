## Installation

There are two ways to install the tool:

- use docker image
- building from sources

### Using docker

To run the tool from a docker image, run the command:

```bash
docker run --rm -v "$PWD:/tmp" -p 8888:8888 dclg/redis-inventory inventory <redis-url>
```

Mounting temp dir (`-v $PWD:/tmp`) makes it possible to run sequence of `index` and `display` while keeping the state
between runs.

Exposing a port (`-p 8888:8888`) is needed for `chart` output type.

If you are using the tool regularly, create an alias:

```bash
alias ri="docker run --rm -v \"$PWD:/tmp\" -p 8888:8888  dclg/redis-inventory"
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

2. Build the binary:

```bash
cd redis-inventory
go build -o redis-inventory main.go
```

3. Run it or move it to one of your PATH directories:

```bash
mv redis-inventory /usr/local/bin/
```
