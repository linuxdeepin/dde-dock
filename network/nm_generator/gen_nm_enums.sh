#!/bin/bash

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
