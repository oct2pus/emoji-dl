# gomoji
Simple CLI to grab emojis from Fediverse Instances implimenting the Mastodon API. (e.g. Mastodon and Pleroma)

## Build

Included is a binary file and the source code, if you'd prefer to compile it yourself do:

```
go build
```

## usage:
```
./gomoji instance.name
or
./gomoji https://instance.name
```

The download path is set to this folder, it'll create a sub directory with `instance.name` as the title, and place all images inside of that.
Currently has an assumption on the filename, modify line 35 to whatever you want to call the file to trim it properly.

## Todo

1. Remove hard coded path
2. allow user to specify file location
