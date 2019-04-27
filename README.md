# emoji-dl
Simple CLI to grab emojis from Fediverse Instances implimenting the Mastodon API
. (e.g. Mastodon and Pleroma)

## Build

to compile it yourself, clone the repo and run

```
go build
```

this project *should* be module aware, if you have below go 1.11, it must be cloned into your $GOPATH.

an alternative build script is included, but is only intended to make creating releases significantly simpler.

## usage:
```
./emoji-dl instance.name
or
./emoji-dl https://instance.name
```

### usage:
the simplest form is `emoji-dl instance-name-here`, which will silently run and then dump into a folder named after the instance (e.g. `emoji-dl mastodon-social` will dump it into a folder called mastodon.social). I intend to only use golang's build in support for flags, so they are a bit tempermental and work slightly differently than the expected UNIX norms.

Here is a list of all flags.
```
-h displays a help message.
-batch=false turns off batch downloading.
-size=INT changes how many files download all at once. (0 and below set batch=false; otherwise ignored if batch=false)
-v to display successful downloads.
-verbose displays successful downloads as well as any errors recieved.
```

**verbose is seperated from v because error messages aren't actually *bad***: Displaying errors might create unnecessary confusion. They're more useful if you experience actual crashes, or are trying to hack on to emoji-dl. I would strongly recommend just using the -v flag instead.

**batch is affected by the global file's open limit**: you can check how many files you can have open at once with `ulimit -Sn`. I will not be fixing issues if they relate to you setting the Size higher than your ulimit. The ability to turn it off is honestly for the morbidly curious or those with ridicious amounts of CPU cores and RAM.

**A larger Size value isn't (inheritly) better**: Size actually determines the maximum number of goroutines since every single file is 1 goroutine. Having way too many open can hog all your CPU/RAM usage. I set the default to 75 because I was unable to reach a GOAWAY response from a server with that value. Files are simply requeued when it fails to download, but it can still take way more time than was neccessary if you find yourself hitting goaway frequently.

The Download path is $PWD/instance.name.social

Windows is not supported. I will take pull requests to fix Windows related issues, but supported platform comes first.
