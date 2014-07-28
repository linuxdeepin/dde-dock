GOROOT = $(shell go env GOROOT)
PREFIX = /usr
TOPDIR = ${CURDIR}
GOPKG_PREFIX = "pkg.linuxdeepin.com/dde-daemon"
GOPATH_LOCAL = "${CURDIR}/gopath:${GOPATH}"

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

out/bin/%:
	env GOPATH=${GOPATH_LOCAL} go build -o $@  ${GOPKG_PREFIX}/bin/${@F}

out/locale/%/LC_MESSAGES/dde-daemon.mo:data/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

build: $(addprefix out/bin/, ${BINARIES})

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
	cp -r data/template ${DESTDIR}${PREFIX}/share/dde-daeman/

	mkdir -pv ${DESTDIR}${PREFIX}/bin
	cp data/tool/wireless_script*.sh ${DESTDIR}${PREFIX}/bin/wireless-script

clean:
	rm -rf out/bin
	rm -rf out/locale
