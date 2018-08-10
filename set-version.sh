#!/bin/bash

if [ $# != 1 ]; then
  echo specify a version to release.
  exit 1
fi

VERSION=$1

sed -i '' -e "s/^const version =.*$/const version = \"$VERSION\"/" cli/cli.go
