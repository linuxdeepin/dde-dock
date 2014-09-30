PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.linuxdeepin.com/dde-daemon

ifndef USE_GCCGO
    GOBUILD = go build
else
   LDFLAGS = $(shell pkg-config --libs gio-2.0 x11 xi xtst libpulse libudev gdk-3.0 gdk-pixbuf-xlib-2.0 gtk+-3.0 sqlite3 fontconfig)
   GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
endif

BINARIES =  \
    backlight_helper \
    dde-session-daemon \
    dde-system-daemon \
    desktop-toggle \
    grub2 \
    grub2ext \
    gtk-thumb-tool \
    search \
    theme-thumb-tool \
    langselector

LANGUAGES = $(basename $(notdir $(wildcard misc/po/*.po)))

all: build

prepare:
	@if [ ! -d ${GOPATH_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX}; \
	fi

out/bin/%:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/bin/${@F}

ifdef USE_GCCGO
out/bin/gtk-thumb-tool:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" \
            go build -compiler gccgo -gccgoflags \
            "$(shell pkg-config --libs gtk+-2.0 libmetacity-private)" \
            -o $@  ${GOPKG_PREFIX}/bin/${@F}
endif

out/locale/%/LC_MESSAGES/dde-daemon.mo:misc/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

build: prepare $(addprefix out/bin/, ${BINARIES})

test: prepare
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" go test -v ./...

install: build translate
	mkdir -pv ${DESTDIR}${PREFIX}/lib/deepin-daemon
	cp out/bin/* ${DESTDIR}${PREFIX}/lib/deepin-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/share/locale
	cp -r out/locale/* ${DESTDIR}${PREFIX}/share/locale

	mkdir -pv ${DESTDIR}/etc/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}/etc/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1
	cp -r misc/services ${DESTDIR}${PREFIX}/share/dbus-1/
	cp -r misc/system-services ${DESTDIR}${PREFIX}/share/dbus-1/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp misc/polkit-action/* ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/glib-2.0/schemas
	cp misc/schemas/* ${DESTDIR}${PREFIX}/share/glib-2.0/schemas

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/lang     ${DESTDIR}${PREFIX}/share/dde-daemon/
	cp -r misc/template ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/bin
	cp misc/tool/wireless_script*.sh ${DESTDIR}${PREFIX}/bin/wireless-script

	mkdir -pv ${DESTDIR}${PREFIX}/share/personalization/thumbnail/autogen
	cp -r misc/thumb_bg ${DESTDIR}${PREFIX}/share/personalization/thumbnail/autogen

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out/bin
	rm -rf out/locale

rebuild: clean build
