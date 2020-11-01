#!/bin/sh

set -e

: "${KUMA_VERSION:=}"
: "${KUMA_ARCH:=amd64}"
: "${DISTRO:=$(grep -oP '(?<=^ID=).+' /etc/os-release | tr -d '"')}"

OS=`uname -s`
if [ "$OS" = "Linux" ]; then
  DISTRO=$(grep -oP '(?<=^ID=).+' /etc/os-release | tr -d '"')
  if [ "$DISTRO" = "amzn" ]; then
    DISTRO="centos"
  fi
elif [ "$OS" = "Darwin" ]; then
  DISTRO="darwin"
else
  printf "ERROR\tOperating system %s not supported by Kuma\n" "$OS"
  exit 1
fi

if [ -z "$DISTRO" ]; then
  printf "ERROR\tUnable to detect the operating system\n"
  exit 1
fi

URL="https://kong.bintray.com/kuma/kuma-$KUMA_VERSION-$DISTRO-$KUMA_ARCH.tar.gz"

if ! type "grep" > /dev/null 2>&1; then
  printf "ERROR\tgrep cannot be found\n"
  exit 1;
fi
if ! type "curl" > /dev/null 2>&1; then
  printf "ERROR\tcurl cannot be found\n"
  exit 1;
fi
if ! type "tar" > /dev/null 2>&1; then
  printf "ERROR\ttar cannot be found\n"
  exit 1;
fi
if ! type "gzip" > /dev/null 2>&1; then
  printf "ERROR\tgzip cannot be found\n"
  exit 1;
fi

if [ -z "$KUMA_VERSION" ]; then
  # Fetching latest Kuma version
  KUMA_VERSION=`curl -s https://kuma.io/latest_version`
  if [ $? -ne 0 ]; then
    printf "ERROR\tUnable to fetch latest Kuma version.\n"
    exit 1
  fi
  if [ -z "$KUMA_VERSION" ]; then
    printf "ERROR\tUnable to fetch latest Kuma version because of a problem with Kuma.\n"
    exit 1
  fi
fi

if ! curl -s --head $URL | head -n 1 | grep "HTTP/1.[01] [23].." > /dev/null; then
  printf "ERROR\tUnable to download Kuma at the following URL: %s\n" "$URL"
  exit 1
fi

if curl -L "$URL" | tar xz; then
  printf "\n"
  printf "INFO\tKuma %s has been downloaded!\n" "$KUMA_VERSION"
  printf "\n"
else
  printf "\n"
  printf "ERROR\tUnable to download Kuma\n"
  exit 1
fi

if ! VERSION=$KUMA_VERSION MODE=$KUMA_MODE PLATFORM=$KUMA_PLATFORM ZONE=$KUMA_ZONE ./scripts/deploy.sh; then
	printf "ERROR\tUnable to deploy\n"
	exit 1
fi

if ! rm -rf kuma-$KUMA_VERSION; then
	printf "ERROR\tUnable to clear temperory files!"
fi

printf "INFO\tKuma service mesh has been installed successfully!!\n"
# printf "Visit: https://meshery.io/adaptors/kuma/dashboard for more information\n"

