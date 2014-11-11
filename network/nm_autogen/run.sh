#!/usr/bin/env bash

go run $(find -type f -iname "*.go" -not -iname "*test.go") \
    -f -b -w --front-end-dir "$HOME/workspace/code/deepin/dde-control-center/modules/network/edit_autogen/"
