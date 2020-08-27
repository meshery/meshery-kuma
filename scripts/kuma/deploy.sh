#!/bin/sh

set -e

# For flat deployment mode
if [ "$MODE" = "flat" ]; then
  if ! kuma-$VERSION/bin/kumactl install control-plane | kubectl apply -f -; then
  	printf "ERROR\tUnable to apply manifests\n"
  	exit 1
  fi
fi

# For distributed deployment mode
if [ "$MODE" = "distributed" ]; then
  if [ "$PLATFORM" = "universal" ]; then
    if ! kuma-$VERSION/bin/KUMA_MODE_MODE=remote KUMA_MODE_REMOTE_ZONE=$ZONE kuma-cp run ; then
    	printf "ERROR\tUnable to create remote control-plane\n"
    	exit 1
    fi
    if ! kumactl appy -f dataplane-universal.yml; then
    	printf "ERROR\tUnable to install ingress data plane\n"
    	exit 1
    fi
    if ! kuma-$VERSION/bin/kumactl install dns | kubectl apply -f -; then
    	printf "ERROR\tUnable to install dns resolvers\n"
    	exit 1
    fi
  else
    if ! kuma-$VERSION/bin/kumactl install control-plane --mode=remote --zone=$ZONE | kubectl apply -f -; then
    	printf "ERROR\tUnable to create global control-plane\n"
    	exit 1
    fi
    if ! kuma-$VERSION/bin/kumactl install ingress | kubectl apply -f -; then
    	printf "ERROR\tUnable to install ingress data plane\n"
    	exit 1
    fi
    if ! kuma-$VERSION/bin/kumactl install dns | kubectl apply -f -; then
    	printf "ERROR\tUnable to install dns resolvers\n"
    	exit 1
    fi
  fi
fi