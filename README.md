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

### flags:
```
-batch=false to turn batch downloading off (don't forget about the maximum open files limit!)
-size=INT to change the size of the value* (0 and below set batch=false; otherwise ignored if batch=false)
-v to display successful downloads.
-verbose to display successful downloads and any errors downloading.
*every one of these spawns a goroutine, increasing the size may actually be worse. size is set to 75 by default.
```

The Download path is the $PWD/instance.name.social