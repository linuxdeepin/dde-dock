TEMPLATE = subdirs
SUBDIRS = dde-dock \
          dde-dock-systray-plugin \
          dde-dock-shutdown-plugin \
          dde-dock-trash-plugin

TRANSLATIONS += translations/dde-dock.ts

qm_files.files += translations/*.qm
qm_files.path   = /usr/share/dde-dock/translations/

INSTALLS += qm_files
