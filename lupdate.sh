#!/bin/bash
cp ".transifexrc" ${HOME}/

lupdate ./ -ts -no-obsolete translations/dde-dock.ts
#lupdate ./ -ts -no-obsolete plugins/dcc-dock-plugin/translations/dcc-dock-plugin.ts

tx push -s --branch m23
