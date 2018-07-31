#!/usr/bin/env bash

set -o errexit
set -o xtrace

if [ -z $@]
  then
    OSES="linux"
  else
    OSES=$@
fi

TOOLS_PKG='github.com/buglloc/rip'
SOURCE_GOPATH=`pwd`/.gopath
VENDOR_GOPATH=`pwd`/vendor

# set up the $GOPATH to use the vendored dependencies as
# well as the source for the yadi
rm -rf .gopath/
mkdir -p .gopath/src/"$(dirname "${TOOLS_PKG}")"
ln -sf `pwd` .gopath/src/$TOOLS_PKG
export GOPATH="${SOURCE_GOPATH}:${VENDOR_GOPATH}"

# remove previous builds
rm -rf .bin

for os in $OSES; do
  if [ "$os" == "windows" ]
    then
      OUTPUT=".bin/${os}/rip.exe"
    else
      OUTPUT=".bin/${os}/rip"
  fi
  mkdir -p ".bin/${os}"
  CGO_ENABLED=0 GOOS="${os}" GOARCH=amd64 go build -o "${OUTPUT}" main.go
done
