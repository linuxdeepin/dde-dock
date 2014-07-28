#!/usr/bin/env bash

go run $(find -type f -iname "*.go" -not -iname "*test.go") -f -b -w
