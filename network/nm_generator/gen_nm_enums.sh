#!/bin/bash

# Copyright (C) 2015 Deepin Technology Co., Ltd.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 3 of the License, or
# (at your option) any later version.

GEN_NM_ENUMS_TPL=nm_enums_gen.go.tpl
GEN_NM_ENUMS_FILE=../nm_enums_gen.go

if [ ! -d /usr/include/libnm ]; then
    echo "please install libnm-dev and retry"
    exit 1
fi

headers=()
for f in /usr/include/libnm/*.h; do
  if grep -q 'NM_\w*_LAST *=' $f; then
	  sed '/NM_\w*_LAST *=/d' $f > /tmp/`basename $f`;
	  headers+=(/tmp/`basename $f`);
  else
	headers+=($f);
  fi;
done
glib-mkenums --template ${GEN_NM_ENUMS_TPL} ${headers[@]} > ${GEN_NM_ENUMS_FILE}
gofmt -w ${GEN_NM_ENUMS_FILE}

echo "GEN ${GEN_NM_ENUMS_FILE}"
