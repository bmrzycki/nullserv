#!/bin/bash

base=$(dirname "$0")
cd "$base"
f='-ldflags "'

# Names copied from https://github.com/ahmetb/govvv

date=$(date --rfc-3339=seconds 2>/dev/null | \
    tr ' ' 'T')
[[ -z $date ]] && date="unknown"
f+="-X main.BuildDate=$date"

sha=$(git rev-parse --short=14 HEAD 2>/dev/null)
[[ -z $sha ]] && sha="unknown"
f+=" -X main.GitCommit=$sha"

state="unknown"
if [[ $sha != "unknown" ]]; then
    tmp=$(git ls-files -mud 2>/dev/null)
    state="clean"
    [[ -n $tmp ]] && state="dirty"
fi
f+=" -X main.GitState=$state"

if [[ -f $base/VERSION ]]; then
    version=$(head -n 1 "$base/VERSION" 2>/dev/null)
fi
[[ -z $version ]] && version="unknown"
f+=" -X main.Version=$version"

echo "$f\""
exit 0
