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

## Using cached index

You can scan redis instance once and then use different visualisations. For scanning use `index` command. To display
cached index use `display` command.

### See more

- [index command](cobra/redis-inventory_index.md)
- [display command](cobra/redis-inventory_display.md)
- [inventoory command](cobra/redis-inventory_inventory.md)
