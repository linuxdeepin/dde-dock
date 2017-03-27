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
	GOLDFLAGS += $(shell pkg-config --libs gio-2.0  x11 xi xtst xcursor xfixes xkbfile libpulse libudev gdk-pixbuf-xlib-2.0 gtk+-3.0 fontconfig librsvg-2.0 libcanberra libbamf3 gudev-1.0 libinput xcb xcb-record)
	GOBUILD = go build -compiler gccgo -gccgoflags "${GOLDFLAGS}"
endif

BINARIES =  \
	    dde-session-initializer\
	    dde-session-daemon \
	    dde-system-daemon \
	    grub2 \
	    grub2ext \
	    search \
	    theme-thumb-tool \
	    backlight_helper \
	    langselector \
	    soundeffect

LANGUAGES = $(basename $(notdir $(wildcard misc/po/*.po)))

all: build

prepare:
	@mkdir -p out/bin
	@if [ ! -d ${GOPATH_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX}; \
		fi

out/bin/%:
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/bin/${@F}

out/bin/default-terminal: bin/default-terminal/default-terminal.c
	gcc -o $@ $(shell pkg-config --cflags --libs gio-unix-2.0) $^

out/bin/desktop-toggle: bin/desktop-toggle/main.c
	gcc -o $@ $(shell pkg-config --cflags --libs x11) $^

out/locale/%/LC_MESSAGES/dde-daemon.mo:misc/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

pot:
	deepin-update-pot misc/po/locale_config.ini

build: prepare out/bin/default-terminal out/bin/desktop-toggle $(addprefix out/bin/, ${BINARIES})

test: prepare
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" go test -v ./...

install: build translate install-dde-data
	mkdir -pv ${DESTDIR}${PREFIX}/lib/deepin-daemon
	cp out/bin/* ${DESTDIR}${PREFIX}/lib/deepin-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/share/locale
	cp -r out/locale/* ${DESTDIR}${PREFIX}/share/locale

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}${PREFIX}/share/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1
	cp -r misc/services ${DESTDIR}${PREFIX}/share/dbus-1/
	cp -r misc/system-services ${DESTDIR}${PREFIX}/share/dbus-1/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp misc/polkit-action/* ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/dde-daemon/*   ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}/var/cache/appearance
	cp -r misc/thumbnail ${DESTDIR}/var/cache/appearance/

#	mkdir -pv ${DESTDIR}${PREFIX}/share/applications
#	cp -r misc/applications/* ${DESTDIR}${PREFIX}/share/applications/

	mkdir -pv ${DESTDIR}${PREFIX}/share/icons/hicolor
	cp -r misc/icons/* ${DESTDIR}${PREFIX}/share/icons/hicolor/

install-dde-data:
	mkdir -pv ${DESTDIR}${PREFIX}/share/dde/
	cp -r misc/data ${DESTDIR}${PREFIX}/share/dde/

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out

rebuild: clean build
