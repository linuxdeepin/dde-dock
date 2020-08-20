%global _smp_mflags -j1

%global debug_package   %{nil}
%global _unpackaged_files_terminate_build 0
%global _missing_build_ids_terminate_build 0
%define __debug_install_post   \
   %{_rpmconfigdir}/find-debuginfo.sh %{?_find_debuginfo_opts} "%{_builddir}/%{?buildsubdir}"\
%{nil}

%global sname deepin-daemon
%global release_name server-industry

Name:           dde-daemon
Version:        5.10.0.23
Release:        5
Summary:        Daemon handling the DDE session settings
License:        GPLv3
URL:            http://shuttle.corp.deepin.com/cache/tasks/18802/unstable-amd64/
Source0:        %{name}-%{version}-%{release_name}.orig.tar.xz
Patch0:         0001-fix-building-error.patch

BuildRequires:  python37
BuildRequires:  compiler(go-compiler)
BuildRequires:  deepin-gettext-tools
BuildRequires:  fontpackages-devel
BuildRequires:  librsvg2-tools
BuildRequires:  pam-devel >= 1.3.1
BuildRequires:  pam >= 1.3.1
BuildRequires:  golang-github-linuxdeepin-go-x11-client-devel
BuildRequires:  golang-golang-org-net-devel
BuildRequires:  glib2-devel
BuildRequires:  gtk3-devel
BuildRequires:  systemd-devel
BuildRequires:  golang-github-axgle-mahonia-devel
BuildRequires:  golang-golang-x-xerrors-devel
BuildRequires:  golang-gopkg-alecthomas-kingpin-devel
BuildRequires:  dde-api-devel
BuildRequires:  golang.org-x-image-devel
BuildRequires:  text-devel
BuildRequires:  resize-devel
BuildRequires:  golang-github-teambition-rrule-go-devel
BuildRequires:  golang-github-rickb777-date-devel
BuildRequires:  golang-github-mozillazg-go-pinyin-devel
BuildRequires:  golang-github-kelvins-sunrisesunset-devel
BuildRequires:  gorm-devel
BuildRequires:  golang-github-cryptix-wav-devel
BuildRequires:  golang-github-mattn-go-sqlite3-devel
BuildRequires:  golang-github-alecthomas-template-devel
BuildRequires:  inflection-devel
BuildRequires:  golang-github-rickb777-plural-devel
BuildRequires:  golang-github-alecthomas-units-devel
BuildRequires:  alsa-lib-devel
BuildRequires:  alsa-lib
BuildRequires:  pulseaudio-libs-devel
BuildRequires:  gdk-pixbuf2-xlib-devel
BuildRequires:  gdk-pixbuf2-xlib
BuildRequires:  libnl3-devel
BuildRequires:  libnl3
BuildRequires:  libgudev-devel
BuildRequires:  libgudev
BuildRequires:  golang-github-davecgh-go-spew-devel
BuildRequires:  libinput-devel
BuildRequires:  libinput
BuildRequires:  golang-github-gosexy-gettext-devel
BuildRequires:  librsvg2-devel
BuildRequires:  librsvg2
BuildRequires:  golang-github-msteinert-pam-devel
BuildRequires:  go-lib-devel
BuildRequires:  golang-github-linuxdeepin-go-dbus-factory-devel
BuildRequires:  deepin-gir-generator
BuildRequires:  libXcursor-devel

Requires:       bluez-libs
Requires:       deepin-desktop-base
Requires:       deepin-desktop-schemas
Requires:       dde-session-ui
Requires:       dde-polkit-agent
Requires:       rfkill
Requires:       gvfs
Requires:       iw

Recommends:     iso-codes
Recommends:     imwheel
Recommends:     mobile-broadband-provider-info
Recommends:     google-noto-mono-fonts
Recommends:     google-noto-sans-fonts

%description
Daemon handling the DDE session settings

%prep
%setup -q -n %{name}-%{version}-%{release_name}
%patch0 -p1

# Fix library exec path
sed -i '/deepin/s|lib|libexec|' Makefile
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
export GOPATH=/usr/share/gocode
%make_build GOBUILD="go build -compiler gc -ldflags \"-B $BUILDID\""
#make GOPATH=/usr/share/gocode

%install
export GOPATH=/usr/share/gocode
%make_install

# fix systemd/logind config
install -d %{buildroot}/usr/lib/systemd/logind.conf.d/
cat > %{buildroot}/usr/lib/systemd/logind.conf.d/10-%{sname}.conf <<EOF
[Login]
HandlePowerKey=ignore
HandleSuspendKey=ignore
EOF

install -d %{buildroot}/%{_libdir}/security/
install -Dm755 %{buildroot}/pam_deepin_auth.so %{buildroot}/%{_libdir}/security/pam_deepin_auth.so

%find_lang %{name}

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

%files -f %{name}.lang
%doc README.md
%license LICENSE
%{_sysconfdir}/default/grub.d/10_deepin.cfg
%{_sysconfdir}/grub.d/35_deepin_gfxmode
%{_sysconfdir}/pam.d/deepin-auth-keyboard
%{_libexecdir}/%{sname}/
%{_prefix}/lib/systemd/logind.conf.d/10-%{sname}.conf
%{_datadir}/dbus-1/services/*.service
%{_datadir}/dbus-1/system-services/*.service
%{_datadir}/dbus-1/system.d/*.conf
%{_datadir}/icons/hicolor/*/status/*
%{_datadir}/%{name}/
%{_datadir}/dde/
%{_datadir}/polkit-1/actions/*.policy
%{_var}/cache/appearance/
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Accounts.pkla
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Grub2.pkla
%{_sysconfdir}/acpi/actions/deepin_lid.sh
%{_sysconfdir}/acpi/events/deepin_lid
%{_sysconfdir}/pulse/daemon.conf.d/10-deepin.conf
/lib/udev/rules.d/80-deepin-fprintd.rules
%{_datadir}/pam-configs/deepin-auth
/var/lib/polkit-1/localauthority/10-vendor.d/com.deepin.daemon.Fprintd.pkla
%{_libdir}/security/pam_deepin_auth.so
/lib/systemd/system/dbus-com.deepin.dde.lockservice.service
/lib/systemd/system/deepin-accounts-daemon.service

%changelog
* Thu May 28 2020 uoser <uoser@uniontech.com> - 5.9.4-2
- Project init.

