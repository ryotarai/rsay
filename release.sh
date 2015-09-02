#!/bin/sh
set -e
gox -os darwin -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"

VERSION="$(git tag -l --points-at HEAD | grep ^v)"
if [ "_$VERSION" = "_" ]; then
    echo "no tag"
    exit 1
fi

ghr -u ryotarai $VERSION pkg/
