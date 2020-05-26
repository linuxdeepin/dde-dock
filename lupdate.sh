#!/bin/bash
rm ./translations/dde-dock.ts
lupdate ./ -ts ./translations/dde-dock.ts
tx push -s -b maintain/20_0102
