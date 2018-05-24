PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.deepin.io/dde/daemon
GOBUILD = go build

ARCH=$(shell uname -m)

ifdef USE_GCCGO
	extra_gccgo_flags = -Os -O2
	ifeq ($(ARCH), sw_64)
		extra_gccgo_flags += -mieee
	endif
	GOBUILD=gccgo_build.pl -p "gio-2.0  x11 xi xtst xcursor xfixes xkbfile libpulse libpulse-simple alsa libudev gdk-pixbuf-xlib-2.0 gtk+-3.0 fontconfig librsvg-2.0 gudev-1.0 libinput xcb xcb-record libnl-genl-3.0" -f "$(extra_gccgo_flags)"
endif

BINARIES =  \
	    dde-session-initializer\
	    dde-session-daemon \
	    dde-system-daemon \
	    grub2 \
	    search \
	    theme-thumb-tool \
	    backlight_helper \
	    langselector \
	    soundeffect \
	    dde-lockservice \
	    dde-authority \
	    dde-authority-pam \
	    default-terminal \
	    dde-greeter-setter

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

out/bin/default-file-manager: bin/default-file-manager/main.c
	gcc $^ $(shell pkg-config --cflags --libs gio-unix-2.0) -o $@

out/bin/desktop-toggle: bin/desktop-toggle/main.c
	gcc $^ $(shell pkg-config --cflags --libs x11) -o $@

out/locale/%/LC_MESSAGES/dde-daemon.mo:misc/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

pot:
	deepin-update-pot misc/po/locale_config.ini

ts:
	deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.daemon.Grub2.policy.in misc/ts/com.deepin.daemon.Grub2.policy
	deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.daemon.accounts.policy.in misc/ts/com.deepin.daemon.accounts.policy

ts_to_policy:
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.daemon.Grub2.policy.in misc/ts/com.deepin.daemon.Grub2.policy misc/polkit-action/com.deepin.daemon.Grub2.policy
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.daemon.accounts.policy.in misc/ts/com.deepin.daemon.accounts.policy misc/polkit-action/com.deepin.daemon.accounts.policy

build: prepare out/bin/default-terminal out/bin/default-file-manager out/bin/desktop-toggle $(addprefix out/bin/, ${BINARIES}) ts_to_policy

test: prepare
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" go test -v ./...

install: build translate install-dde-data install-icons
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
	cp misc/polkit-action/*.policy ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/dde-daemon/*   ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}/var/cache/appearance
	cp -r misc/thumbnail ${DESTDIR}/var/cache/appearance/

	mkdir -pv ${DESTDIR}/lib/systemd/system/
	cp -f misc/systemd/services/* ${DESTDIR}/lib/systemd/system/

	mkdir -pv ${DESTDIR}/etc/pam.d/
	cp -f misc/etc/pam.d/* ${DESTDIR}/etc/pam.d/

install-dde-data:
	mkdir -pv ${DESTDIR}${PREFIX}/share/dde/
	cp -r misc/data ${DESTDIR}${PREFIX}/share/dde/

install-icons:
	python3 misc/icons/install_to_hicolor.py -d status -o out/icons misc/icons/status
	mkdir -pv ${DESTDIR}${PREFIX}/share/icons/
	cp -r out/icons/hicolor ${DESTDIR}${PREFIX}/share/icons/

clean:
	rm -rf ${GOPATH_DIR}
	rm -rf out

rebuild: clean build
