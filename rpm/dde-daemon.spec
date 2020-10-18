%global _smp_mflags -j1

%if 0%{?fedora} == 0
%global debug_package   %{nil}
%global _unpackaged_files_terminate_build 0
%global _missing_build_ids_terminate_build 0
%define __debug_install_post   \
   %{_rpmconfigdir}/find-debuginfo.sh %{?_find_debuginfo_opts} "%{_builddir}/%{?buildsubdir}"\
%{nil}
%endif

%global sname deepin-daemon
%global repo dde-daemon
%global release_name server-industry

%if 0%{?fedora}
Name:           %{sname}
%else
Name:           %{repo}
%endif
Version:        5.12.0.19
Release:        1%{?fedora:%dist}
Summary:        Daemon handling the DDE session settings
License:        GPLv3
%if 0%{?fedora}
URL:            https://github.com/linuxdeepin/dde-daemon
Source0:        %{url}/archive/%{version}/%{repo}-%{version}.tar.gz
# upstream default mono font set to 'Noto Mono', which is not yet available in
# Fedora. We change to 'Noto Sans Mono'
Source1:        fontconfig.json
%else
URL:            http://shuttle.corp.deepin.com/cache/tasks/18802/unstable-amd64/
Source0:        %{repo}-%{version}.orig.tar.xz
%endif
Source2:        deepin-auth

BuildRequires:  python3
%if 0%{?fedora}
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 aarch64 %{arm}}
BuildRequires:  golang(pkg.deepin.io/dde/api/dxinput) >= 3.1.26
BuildRequires:  golang(github.com/linuxdeepin/go-dbus-factory/org.bluez)
BuildRequires:  golang(github.com/linuxdeepin/go-x11-client)
BuildRequires:  golang(github.com/BurntSushi/xgb)
BuildRequires:  golang(github.com/BurntSushi/xgbutil)
BuildRequires:  golang(github.com/axgle/mahonia)
BuildRequires:  golang(github.com/msteinert/pam)
BuildRequires:  golang(github.com/nfnt/resize)
BuildRequires:  golang(github.com/cryptix/wav)
BuildRequires:  golang(gopkg.in/alecthomas/kingpin.v2)
BuildRequires:  golang(gopkg.in/yaml.v2)
BuildRequires:  golang(github.com/gosexy/gettext)
BuildRequires:  golang(github.com/jinzhu/gorm)
BuildRequires:  golang(github.com/jinzhu/gorm/dialects/sqlite)
BuildRequires:  golang(github.com/kelvins/sunrisesunset)
BuildRequires:  golang(github.com/rickb777/date)
BuildRequires:  golang(github.com/teambition/rrule-go)
BuildRequires:  golang(github.com/davecgh/go-spew/spew)
%else
BuildRequires:  gocode
BuildRequires:  ddcutil-devel
BuildRequires:  resize-devel
BuildRequires:  gorm-devel
BuildRequires:  inflection-devel
%endif
BuildRequires:  compiler(go-compiler)
BuildRequires:  deepin-gettext-tools
BuildRequires:  fontpackages-devel
BuildRequires:  librsvg2-tools
BuildRequires:  pam-devel >= 1.3.1
BuildRequires:  glib2-devel
BuildRequires:  gtk3-devel
BuildRequires:  systemd-devel
BuildRequires:  alsa-lib-devel
BuildRequires:  pulseaudio-libs-devel
BuildRequires:  gdk-pixbuf2-xlib-devel
BuildRequires:  libnl3-devel
BuildRequires:  libgudev-devel
BuildRequires:  libinput-devel
BuildRequires:  librsvg2-devel
BuildRequires:  libXcursor-devel
BuildRequires:  pkgconfig(sqlite3)

Requires:       bluez-libs%{?_isa}
Requires:       deepin-desktop-base
Requires:       deepin-desktop-schemas
%if 0%{?fedora}
Requires:       deepin-session-ui
Requires:       deepin-polkit-agent
%else
Requires:       dde-session-ui
Requires:       dde-polkit-agent
%endif
Requires:       rfkill
Requires:       gvfs
Requires:       iw

Recommends:     iso-codes
Recommends:     imwheel
Recommends:     mobile-broadband-provider-info
Recommends:     google-noto-mono-fonts
Recommends:     google-noto-sans-fonts
%if 0%{?fedora}
Recommends:     google-noto-sans-mono-fonts
%endif

%description
Daemon handling the DDE session settings

%prep
%autosetup -n %{repo}-%{version}
patch langselector/locale.go < rpm/locale.go.patch
install -m 644 %{SOURCE2} misc/etc/pam.d/deepin-auth

# Fix library exec path
sed -i '/deepin/s|lib|libexec|' Makefile
sed -i '/systemd/s|lib|usr/lib|' Makefile
sed -i 's:/lib/udev/rules.d:%{_udevrulesdir}:' Makefile
sed -i '/${DESTDIR}\/usr\/lib\/deepin-daemon\/service-trigger/s|${DESTDIR}/usr/lib/deepin-daemon/service-trigger|${DESTDIR}/usr/libexec/deepin-daemon/service-trigger|g' Makefile
sed -i '/${DESTDIR}${PREFIX}\/lib\/deepin-daemon/s|${DESTDIR}${PREFIX}/lib/deepin-daemon|${DESTDIR}${PREFIX}/usr/libexec/deepin-daemon|g' Makefile
sed -i 's|lib/NetworkManager|libexec|' network/utils_test.go

for file in $(grep "/usr/lib/deepin-daemon" * -nR |awk -F: '{print $1}')
do
    sed -i 's|/usr/lib/deepin-daemon|/usr/libexec/deepin-daemon|g' $file
done

# Fix grub.cfg path
sed -i 's|boot/grub|boot/grub2|' grub2/{grub2,grub_params,theme}.go

# Fix activate services failed (Permission denied)
# dbus service
pushd misc/system-services/
sed -i '$aSystemdService=deepin-accounts-daemon.service' com.deepin.system.Power.service \
    com.deepin.daemon.{Accounts,Apps,Daemon}.service \
    com.deepin.daemon.{Gesture,SwapSchedHelper,Timedated}.service
sed -i '$aSystemdService=dbus-com.deepin.dde.lockservice.service' com.deepin.dde.LockService.service
popd
# systemd service
cat > misc/systemd/services/dbus-com.deepin.dde.lockservice.service <<EOF
[Unit]
Description=Deepin Lock Service
Wants=user.slice dbus.socket
After=user.slice dbus.socket

[Service]
Type=dbus
BusName=com.deepin.dde.LockService
ExecStart=%{_libexecdir}/%{sname}/dde-lockservice

[Install]
WantedBy=graphical.target
EOF

# Replace reference of google-chrome to chromium-browser
sed -i 's/google-chrome/chromium-browser/g' misc/dde-daemon/mime/data.json

%build
BUILDID="0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')"
%if 0%{?fedora}
export GOPATH="$(pwd)/build:%{gopath}"
%else
export GOPATH=/usr/share/gocode
%endif
%make_build GO_BUILD_FLAGS=-trimpath GOBUILD="go build -compiler gc -ldflags \"-B $BUILDID\""

%install
BUILDID="0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')"
%if 0%{?fedora}
export GOPATH="$(pwd)/build:%{gopath}"
%else
export GOPATH=/usr/share/gocode
%endif
%make_install PAM_MODULE_DIR=%{_libdir}/security GOBUILD="go build -compiler gc -ldflags \"-B $BUILDID\""

# fix systemd/logind config
install -d %{buildroot}/usr/lib/systemd/logind.conf.d/
cat > %{buildroot}/usr/lib/systemd/logind.conf.d/10-%{sname}.conf <<EOF
[Login]
HandlePowerKey=ignore
HandleSuspendKey=ignore
EOF

%if 0%{?fedora}
# install default settings
install -Dm644 %{SOURCE1} \
    %{buildroot}%{_datadir}/deepin-default-settings/fontconfig.json
%endif

%find_lang %{repo}

%post
if [ $1 -ge 1 ]; then
  systemd-sysusers %{sname}.conf
  %{_sbindir}/alternatives --install %{_bindir}/x-terminal-emulator \
    x-terminal-emulator %{_libexecdir}/%{sname}/default-terminal 30
fi

%preun
if [ $1 -eq 0 ]; then
  %{_sbindir}/alternatives --remove x-terminal-emulator \
    %{_libexecdir}/%{sname}/default-terminal
fi

%postun
if [ $1 -eq 0 ]; then
  rm -f /var/cache/deepin/mark-setup-network-services
  rm -f /var/log/deepin.log 
fi

%files -f %{repo}.lang
%doc README.md
%license LICENSE
%{_sysconfdir}/default/grub.d/10_deepin.cfg
%{_sysconfdir}/grub.d/35_deepin_gfxmode
%{_sysconfdir}/pam.d/deepin-auth
%{_sysconfdir}/pam.d/deepin-auth-keyboard
%{_sysconfdir}/NetworkManager/conf.d/deepin.dde.daemon.conf
%{_sysconfdir}/modules-load.d/i2c_dev.conf
%{_libexecdir}/%{sname}/
%{_prefix}/lib/systemd/logind.conf.d/10-%{sname}.conf
%{_datadir}/dbus-1/services/*.service
%{_datadir}/dbus-1/system-services/*.service
%{_datadir}/dbus-1/system.d/*.conf
%{_datadir}/icons/hicolor/*/status/*
%{_datadir}/%{repo}/
%{_datadir}/dde/
%{_datadir}/polkit-1/actions/*.policy
%{_var}/cache/appearance/
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Accounts.pkla
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Grub2.pkla
%{_sysconfdir}/acpi/actions/deepin_lid.sh
%{_sysconfdir}/acpi/events/deepin_lid
# This directory is not provided by any other package.
%dir %{_sysconfdir}/pulse/daemon.conf.d
%{_sysconfdir}/pulse/daemon.conf.d/10-deepin.conf
%{_udevrulesdir}/80-deepin-fprintd.rules
%{_datadir}/pam-configs/deepin-auth
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Fprintd.pkla
%{_libdir}/security/pam_deepin_auth.so
%{_unitdir}/dbus-com.deepin.dde.lockservice.service
%{_unitdir}/deepin-accounts-daemon.service
%{_unitdir}/hwclock_stop.service
%if 0%{?fedora}
%{_datadir}/deepin-default-settings/
%endif

%changelog
* Thu May 28 2020 uoser <uoser@uniontech.com> - 5.9.4-2
- Project init.

