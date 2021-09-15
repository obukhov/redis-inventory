Redis URL
========

Redis URL can be provided in one of two formats:

- `<host>:<port>` - simplified format,
- `redis://[:<password>@]<host>:<port>[/<dbIndex>]` - in case if you want to specify password or DB Index for the connection.

Examples (can be used with provided docker-compose):
- `redis-inventory localhost:63795`
- `redis-inventory redis://localhost:63795/1`
- `redis-inventory redis://:12345@localhost:63795`
- `redis-inventory redis://:12345@localhost:63795/3`