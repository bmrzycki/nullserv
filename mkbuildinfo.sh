#!/bin/sh

base=`dirname "$0"`
base=`cd "$base" && pwd`

date=`date -R 2>/dev/null`
test -z "$date" && date="unknown"

version=`head -n 1 "$base/VERSION" 2>/dev/null`
test -z "$version" && version="unknown"

sha=`git -C "$base" describe --always --long --dirty --tags 2>/dev/null`
test -z "$sha" && sha="unknown"

cdate=`git -C "$base" show  -s --format=%cd \
           --date=format:'%a, %d %b %Y %H:%M:%S %z' HEAD 2>/dev/null`
test -z "$cdate" && version="unknown"

echo "Generating $base/buildinfo.go"
sed \
    -e "s/%DATE%/$date/g" \
    -e "s/%VERSION%/$version/g" \
    -e "s/%SHA%/$sha/g" \
    -e "s/%CDATE%/$cdate/g" \
    "$base/buildinfo.go.in" > "$base/buildinfo.go"
