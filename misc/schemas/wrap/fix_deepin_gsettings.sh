#!/bin/bash

echo "==> rename files"
mv -vf wm-schemas.convert deepin-wm-schemas.convert &>/dev/null
mv -vf gsettings-desktop-schemas.convert deepin-gsettings-desktop-schemas.convert &>/dev/null
mv -vf gnome-settings-daemon.convert deepin-gnome-settings-daemon.convert &>/dev/null
for f in *.xml; do
  mv -vf "${f}" "${f/#org.gnome/com.deepin.wrap.gnome}" &>/dev/null
done

echo "==> rename gsettings path"
for f in *.xml *.convert; do
  sed -e 's=org\.gnome=com.deepin.wrap.gnome=g' \
      -e 's=/org/gnome=/com/deepin/wrap/gnome=g' -i "${f}"
done

echo "==> show unmet files"
ls -1 | grep -v deepin
