#!/bin/bash

if [ -n "$(gofmt -l service delivery cmd example)" ]; then
  echo "Go code is not formatted:"
  gofmt -d .
  exit 1
fi
