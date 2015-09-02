#!/bin/sh
set -ex
gox -os darwin -output "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"

VERSION="$(git tag -l --points-at HEAD | grep ^v)"
git push --tags origin master
ghr -u ryotarai $VERSION pkg/
