# Letter

```
[8.bytes TYPE] [16.byts ENDPOINT]
[PARAMS]

[BODY]
```

`<#b7>`
`<35#b7>`
`<root@127.0.0.1:16400>`
`<@127.0.0.1:16400>`

# Channel

- `00`
    - `info.properties`
    - `messages.list`
    - `tiny.db`
    - `small.db`
    - `medium.db`
    - `large.db`
- `b7`
    - `info.properties`
    - `messages.list`
    - `tiny.db`
    - `small.db`
    - `medium.db`
    - `large.db`

## info.properties

```properties
After=b5
Type=text
Name=Example Channel!
LatestMessageIndex=15
```

## messages.list

```
[uint64 timestamp][byte chunk][uint64 line][16 bytes user]
```

## tiny.db / small.db / medium.db / large.db

Databases of messages, all contain data in chunks of bytes as follows:
- `16 bytes (~4-16 characters)` (tiny)
- `256 bytes (~64-256 characters)` (small)
- `1024 bytes (~256-1024 characters)` (medium)
- `8192 bytes (~2048-8192 characters)` (large)

A message is a UTF-8 encoded markdown text. There are several special sequences of characters:

- `<@user_id>` — Mention a user
- `<#b7>` — Mention a channel
