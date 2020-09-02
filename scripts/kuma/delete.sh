#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"
URL="https://github.com/openservicemesh/osm/releases/download/$OSM_VERSION/osm-$OSM_VERSION-$OS-$OSM_ARCH.tar.gz"

if curl -L "$URL" | tar xz; then
  printf "INFO\tosmctl %s has been downloaded!\n" "$OSM_VERSION"
else
  printf "ERROR\tUnable to download osmctl\n"
  exit 1
fi

if ! ./$OS-$OSM_ARCH/osm delete; then
	printf "ERROR\tUnable to delete\n"
	exit 1
fi
