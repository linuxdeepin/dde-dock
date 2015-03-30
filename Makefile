PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.linuxdeepin.com/dde-daemon

ifndef USE_GCCGO
    GOBUILD = go build
else
   LDFLAGS = $(shell pkg-config --libs gio-2.0 x11 xi xtst xcursor xfixes libpulse libudev gdk-3.0 gdk-pixbuf-xlib-2.0 gtk+-3.0 sqlite3 fontconfig)
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

pot:
	deepin-update-pot misc/po/locale_config.ini

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

	mkdir -pv ${DESTDIR}/etc/systemd/system
	cp -rf misc/etc/systemd/system/* ${DESTDIR}/etc/systemd/system/
	mkdir -pv ${DESTDIR}/etc/systemd/system/graphical.target.wants
	ln -sv ${DESTDIR}/etc/systemd/system/dbus-com.deepin.daemon.Accounts.service ${DESTDIR}/etc/systemd/system/graphical.target.wants/dbus-com.deepin.daemon.Accounts.service

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1
	cp -r misc/services ${DESTDIR}${PREFIX}/share/dbus-1/
	cp -r misc/system-services ${DESTDIR}${PREFIX}/share/dbus-1/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp misc/polkit-action/* ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/glib-2.0/schemas
	cp misc/schemas/*.xml ${DESTDIR}${PREFIX}/share/glib-2.0/schemas/
	cp misc/schemas/wrap/*.xml ${DESTDIR}${PREFIX}/share/glib-2.0/schemas/
	cp misc/schemas/wrap/*.convert ${DESTDIR}${PREFIX}/share/glib-2.0/schemas/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/usr/share/dde-daemon/*   ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/share/personalization/thumbnail
	cp -r misc/thumbnail/* ${DESTDIR}${PREFIX}/share/personalization/thumbnail/

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out/bin
	rm -rf out/locale

rebuild: clean build
