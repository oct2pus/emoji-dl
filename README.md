# emoji-dl
Simple CLI to grab emojis from Fediverse Instances implimenting the Mastodon API
. (e.g. Mastodon and Pleroma)

## Build

Included is a binary file and the source code, if you'd prefer to compile it yourself do:

```
go build
```

## usage:
```
./emoji-dl instance.name
or
./emoji-dl https://instance.name
```

The Download path is the $PWD/instance.name.social