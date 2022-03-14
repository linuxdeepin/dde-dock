#!/bin/bash
cp ".transifexrc" ${HOME}/

lupdate ./ -ts -no-obsolete translations/org.deepin.dde.dock.ts
lupdate ./ -ts -no-obsolete plugins/dcc-dock-plugin/translations/dcc-dock-plugin.ts

tx push -s -b m20
