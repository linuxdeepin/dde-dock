PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.deepin.io/dde/daemon
GOBUILD = go build

BINARIES =  \
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

out/bin/%: prepare
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/bin/${@F}

out/bin/default-file-manager: bin/default-file-manager/main.c
	gcc $^ $(shell pkg-config --cflags --libs gio-unix-2.0) -o $@

out/bin/desktop-toggle: bin/desktop-toggle/main.c
	gcc $^ $(shell pkg-config --cflags --libs x11) -o $@

out/pam_deepin_auth.so: misc/pam-module/deepin_auth.c
	gcc -fPIC -shared -Wall $(shell pkg-config --libs libsystemd) -o $@ $^
	chmod -x $@

out/locale/%/LC_MESSAGES/dde-daemon.mo: misc/po/%.po
	mkdir -p $(@D)
	msgfmt -o $@ $<

translate: $(addsuffix /LC_MESSAGES/dde-daemon.mo, $(addprefix out/locale/, ${LANGUAGES}))

pot:
	deepin-update-pot misc/po/locale_config.ini

POLICIES=accounts Grub2 Fprintd
ts:
	for i in $(POLICIES); do \
		deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.daemon.$$i.policy.in misc/ts/com.deepin.daemon.$$i.policy; \
	done

ts_to_policy:
	for i in $(POLICIES); do \
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.daemon.$$i.policy.in misc/ts/com.deepin.daemon.$$i.policy misc/polkit-action/com.deepin.daemon.$$i.policy; \
	done

build: prepare out/bin/default-terminal out/bin/default-file-manager out/bin/desktop-toggle out/pam_deepin_auth.so $(addprefix out/bin/, ${BINARIES}) ts_to_policy

test: prepare
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" go test -v ./...

install: build translate install-dde-data install-icons install-pam-module
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

	mkdir -pv ${DESTDIR}/var/lib/polkit-1/localauthority/10-vendor.d
	cp misc/polkit-localauthority/*.pkla ${DESTDIR}/var/lib/polkit-1/localauthority/10-vendor.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-daemon
	cp -r misc/dde-daemon/*   ${DESTDIR}${PREFIX}/share/dde-daemon/

	mkdir -pv ${DESTDIR}${PREFIX}/share/pam-configs
	cp -r misc/pam-configs/* ${DESTDIR}${PREFIX}/share/pam-configs

	mkdir -pv ${DESTDIR}/var/cache/appearance
	cp -r misc/thumbnail ${DESTDIR}/var/cache/appearance/

	mkdir -pv ${DESTDIR}/lib/systemd/system/
	cp -f misc/systemd/services/* ${DESTDIR}/lib/systemd/system/

	mkdir -pv ${DESTDIR}/etc/pam.d/
	cp -f misc/etc/pam.d/* ${DESTDIR}/etc/pam.d/

	mkdir -pv ${DESTDIR}/etc/default/grub.d
	cp -f misc/etc/default/grub.d/* ${DESTDIR}/etc/default/grub.d

	mkdir -pv ${DESTDIR}/etc/grub.d
	cp -f misc/etc/grub.d/* ${DESTDIR}/etc/grub.d

	mkdir -pv ${DESTDIR}/etc/acpi/events
	cp -f misc/etc/acpi/events/* ${DESTDIR}/etc/acpi/events/

	mkdir -pv ${DESTDIR}/etc/acpi/actions
	cp -f misc/etc/acpi/actions/* ${DESTDIR}/etc/acpi/actions/

	mkdir -pv ${DESTDIR}/etc/pulse/daemon.conf.d
	cp -f misc/etc/pulse/daemon.conf.d/*.conf ${DESTDIR}/etc/pulse/daemon.conf.d/

	mkdir -pv ${DESTDIR}/lib/udev/rules.d
	cp -f misc/udev-rules/*.rules ${DESTDIR}/lib/udev/rules.d/

install-pam-module:
	mkdir -pv ${DESTDIR}/${PAM_MODULE_DIR}
	cp -f out/pam_deepin_auth.so ${DESTDIR}/${PAM_MODULE_DIR}

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

check_code_quality: prepare
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" go vet ./...
