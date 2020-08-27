#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"
URL="https://github.com/openservicemesh/osm/releases/download/$OSM_VERSION/osm-$OSM_VERSION-$OS-$OSM_ARCH.tar.gz"

printf "INFO\tDownloading osmctl from: %s" "$URL"
printf "\n\n"

if curl -L "$URL" | tar xz; then
  printf "\n"
  printf "INFO\tosmctl %s has been downloaded!\n" "$OSM_VERSION"
  printf "\n"
else
  printf "\n"
  printf "ERROR\tUnable to download osmctl\n"
  exit 1
fi

printf "INFO\tDeleting deployment......\n"
if ! ./$OS-$OSM_ARCH/osm delete; then
	printf "ERROR\tUnable to delete\n"
	exit 1
fi
printf "INFO\tDeleted successfully\n"