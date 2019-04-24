#!/usr/bin/fish
# this build script is very basic. This is just intended to build all
# "supported"(big airquote) platforms executables.

mkdir -p build
set -x GOOS linux
go build -o ./build/emoji-dl-linux
set -x GOOS darwin
go build -o ./build/emoji-dl-osx
set -x GOOS freebsd
go build -o ./build/emoji-dl-free
set -x GOOS netbsd
go build -o ./build/emoji-dl-net
set -x GOOS linux; set -x GOARCH 386
go build -o ./build/emoji-dl-linux-386