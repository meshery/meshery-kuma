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

if curl -L "$URL" | tar xz; then
  printf "\n"
  printf "INFO\tkumactl %s has been downloaded\n" "$KUMA_VERSION"
  printf "\n"
else
  printf "\n"
  printf "ERROR\tUnable to download kumactl\n"
  exit 1
fi

# For flat deployment mode
if [ "$KUMA_MODE" = "flat" ]; then
  if ! kuma-$KUMA_VERSION/bin/kumactl install control-plane | kubectl delete -f -; then
  	printf "ERROR\tUnable to delete manifests\n"
  	exit 1
  fi
fi