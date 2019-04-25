#!/usr/bin/fish
# this build script is very basic. It is intended to build all
# "supported"(big airquote) platforms executables.
# if you just want to compile the application,
# I recommend just using `go build` in the project root folder.

# build compiles the application for an ARCH and OS.
function build 
    set -x GOOS $argv[1]
    set -x GOARCH $argv[2]
    set file emoji-dl-{$argv[1]}-{$argv[2]}
    set bath {$argv[3]}/{$file}
    go build -o $bath
    # read write execute for user, only read for group and other
    chmod 744 $bath
    echo $file
    sha256sum $bath
    zip {$file}.zip $bath
end

mkdir -p binaries
build linux amd64 binaries
build linux 386 binaries
build darwin amd64 binaries
build freebsd amd64 binaries
build netbsd amd64 binaries
