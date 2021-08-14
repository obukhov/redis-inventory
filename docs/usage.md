Usage
=====

## General interface

```bash
redis-inventory inventory <host>:<port> [--output=<output type>] [--output-params=<querstring serialized params>]
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

### Chart

| Option name  | Description                                     | Default   |
|--------------|-------------------------------------------------|-----------|
| depth        | Maximum nesting level for keys before grouping  | `10`      |
| port         | Use custom character to pad nested level        | `8888`    |

## Using cached index

You can scan redis instance once and then use different visualisations. For scanning use `index` command. To display
cached index use `display` command.

## Customising level separators

It is quite common to use ":" as a level separator in redis. Sometimes, though, you want to use other punctuation (for
example "_" or "-") to also be a level separator. Redis-inventory can use custom characters for splitting keys:

```bash
redis-inventory inventory <host>:<port> --separators=":_-"
```

It will work even if those characters sometimes are used not for level separation, as a node with just one child will
automatically be merged with it.

## See more

- [index command](cobra/redis-inventory_index.md)
- [display command](cobra/redis-inventory_display.md)
- [inventory command](cobra/redis-inventory_inventory.md)
