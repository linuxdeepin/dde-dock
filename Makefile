PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.linuxdeepin.com/dde-daemon

BINARIES =  \
    backlight_helper \
    dde-session-daemon \
    dde-system-daemon \
    desktop-toggle \
    grub2 \
    gtk-thumb-tool \
    search \
    theme-thumb-tool

LANGUAGES = $(basename $(notdir $(wildcard data/po/*.po)))

all: build

prepare:
	@if [ ! -d ${GOPATH_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX}; \
	fi

out/bin/%:
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" go build -o $@  ${GOPKG_PREFIX}/bin/${@F}

out/locale/%/LC_MESSAGES/dde-daemon.mo:data/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

build: prepare $(addprefix out/bin/, ${BINARIES})

install: build translate
	mkdir -pv ${DESTDIR}${PREFIX}/lib/deepin-daemon
	cp out/bin/* ${DESTDIR}${PREFIX}/lib/deepin-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/share/locale
	cp -r out/locale/* ${DESTDIR}${PREFIX}/share/locale

	mkdir -pv ${DESTDIR}/etc/dbus-1/system.d
	cp data/conf/*.conf ${DESTDIR}/etc/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1
	cp -r data/services ${DESTDIR}${PREFIX}/share/dbus-1/
	cp -r data/system-services ${DESTDIR}${PREFIX}/share/dbus-1/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp data/polkit-action/* ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/glib-2.0/schemas
	cp data/schemas/* ${DESTDIR}${PREFIX}/share/glib-2.0/schemas

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r data/lang     ${DESTDIR}${PREFIX}/share/dde-daemon/
	cp -r data/template ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/bin
	cp data/tool/wireless_script*.sh ${DESTDIR}${PREFIX}/bin/wireless-script

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out/bin
	rm -rf out/locale

rebuild: clean build
