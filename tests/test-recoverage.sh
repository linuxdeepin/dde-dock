#!/bin/bash

BUILD_DIR=build
REPORT_DIR=report

cd ../
rm -rf $BUILD_DIR
mkdir $BUILD_DIR
cd $BUILD_DIR
cmake ../
make -j 16

cd tests/

./dde_dock_unit_test
lcov -c -d ./ -o cover.info
lcov -e cover.info '*/frame/*' '*/dde-dock/widgets/*' -o code.info
lcov -r code.info '*/dbus/*' '*/xcb/*' -o final.info

rm -rf ../../tests/$REPORT_DIR
mkdir -p ../../tests/$REPORT_DIR
genhtml -o ../../tests/$REPORT_DIR final.info
