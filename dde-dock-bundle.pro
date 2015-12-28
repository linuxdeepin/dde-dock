TEMPLATE = subdirs
SUBDIRS = dde-dock \
          dde-dock-systray-plugin \
          dde-dock-shutdown-plugin \
          dde-dock-trash-plugin

# Automating generation .qm files from .ts files
system($$PWD/translate_generation.sh)

qm_files.files += translations/*.qm
qm_files.path   = /usr/share/dde-dock/translations/

INSTALLS += qm_files
