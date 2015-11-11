PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.deepin.io/dde/daemon

ifndef USE_GCCGO
	GOLDFLAGS = -ldflags '-s -w'
else
	GOLDFLAGS = -s -w  -Os -O2
endif

ifdef GODEBUG
	GOLDFLAGS =
endif

ifndef USE_GCCGO
	GOBUILD = go build ${GOLDFLAGS}
else
	GOLDFLAGS += $(shell pkg-config --libs gio-2.0  x11 xi xtst xcursor xfixes xkbfile libpulse libudev gdk-pixbuf-xlib-2.0 gtk+-3.0 sqlite3 fontconfig)
	GOBUILD = go build -compiler gccgo -gccgoflags "${GOLDFLAGS}"
endif

BINARIES =  \
	    dde-preload \
	    dde-session-daemon \
	    dde-system-daemon \
	    desktop-toggle \
	    grub2 \
	    grub2ext \
	    search \
	    theme-thumb-tool \
	    backlight_helper \
	    langselectora

LANGUAGES = $(basename $(notdir $(wildcard misc/po/*.po)))

all: build

prepare:
	@mkdir -p out/bin
	@if [ ! -d ${GOPATH_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX}; \
		fi

out/bin/%:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/bin/${@F}

out/bin/default-terminal: bin/default-terminal/default-terminal.c
	gcc -o $@ $(shell pkg-config --cflags --libs gio-unix-2.0) $^

ifdef USE_GCCGO
out/bin/theme-thumb-tool:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" \
		go build -compiler gccgo -gccgoflags \
		"$(shell pkg-config --libs glib-2.0 gdk-3.0 cairo-ft poppler-glib libmetacity-private )" \
		-o $@  ${GOPKG_PREFIX}/bin/${@F}
endif

out/locale/%/LC_MESSAGES/dde-daemon.mo:misc/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

pot:
	deepin-update-pot misc/po/locale_config.ini

build: prepare $(addprefix out/bin/, ${BINARIES}) out/bin/deepin-default-terminal

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

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/dde-daemon/*   ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}/var/cache/appearance
	cp -r misc/thumbnail ${DESTDIR}/var/cache/appearance/

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out/bin
	rm -rf out/locale

rebuild: clean build
