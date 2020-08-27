#!/bin/sh

set -e

if ! type "kubectl" > /dev/null 2>&1; then
  printf "ERROR\tgrep cannot be found\n"
  exit 1;
fi

printf "INFO\Deleting SMI deployment......\n"
if ! kubectl delete -f ./scripts/smi/smi.yaml; then
	printf "ERROR\tUnable to delete\n"
	exit 1
fi
printf "INFO\tDeleted successfull!!\n"
