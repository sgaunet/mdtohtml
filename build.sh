#!/usr/bin/env bash

IMG="sgaunet/mdtohtml:latest"

docker build --build-arg="VERSION=development" . -t "$IMG"
rc=$?

if [ "$rc" != "0" ]
then
  echo "Build FAILED"
  exit 1
fi

docker push "$IMG"