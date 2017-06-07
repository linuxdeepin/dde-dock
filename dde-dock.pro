
TEMPLATE = subdirs
SUBDIRS = frame \
          plugins


# Automating generation .qm files from .ts files
CONFIG(release, debug|release) {
    !system($$PWD/translate_generation.sh): error("Failed to generate translation")
}
TRANSLATIONS    = translations/dde-dock.ts

qm_files.path = $${PREFIX}/share/dde-dock/translations/
qm_files.files = translations/*.qm

INSTALLS = qm_files
