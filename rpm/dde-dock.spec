%global sname deepin-dock

Name:           dde-dock
Version:        5.4.9
Release:        1
Summary:        Deepin desktop-environment - Dock module
License:        GPLv3
URL:            http://shuttle.corp.deepin.com/cache/repos/eagle/release-candidate/RERFNS4wLjAuNjU3NQ/pool/main/d/dde-dock/
Source0:        %{name}-%{version}.orig.tar.xz	

BuildRequires:  cmake
BuildRequires:  gcc-c++
BuildRequires:  pkgconfig(dbusmenu-qt5)
BuildRequires:  pkgconfig(dde-network-utils)
BuildRequires:  dtkwidget-devel >= 5.1
BuildRequires:  dtkcore-devel >= 5.1
BuildRequires:  pkgconfig(dframeworkdbus) >= 2.0
BuildRequires:  pkgconfig(gsettings-qt)
BuildRequires:  pkgconfig(gtk+-2.0)
BuildRequires:  pkgconfig(Qt5Core)
BuildRequires:  pkgconfig(Qt5Gui)
BuildRequires:  pkgconfig(Qt5DBus)
BuildRequires:  pkgconfig(Qt5X11Extras)
BuildRequires:  pkgconfig(Qt5Svg)
BuildRequires:  pkgconfig(x11)
BuildRequires:  pkgconfig(xtst)
BuildRequires:  pkgconfig(xext)
BuildRequires:  pkgconfig(xcb-composite)
BuildRequires:  pkgconfig(xcb-ewmh)
BuildRequires:  pkgconfig(xcb-icccm)
BuildRequires:  pkgconfig(xcb-image)
BuildRequires:  qt5-linguist
Requires:       dbusmenu-qt5
Requires:       dde-network-utils
Requires:       dde-qt-dbus-factory
Requires:       xcb-util-wm
Requires:       xcb-util-image
Requires:       libxcb

%description
Deepin desktop-environment - Dock module.

%package devel
Summary:        Development package for %{sname}
Requires:       %{name}%{?_isa} = %{version}-%{release}

%description devel
Header files and libraries for %{sname}.

%package onboard-plugin
Summary:        deepin desktop-environment - dock plugin
Requires:       %{name}%{?_isa} = %{version}-%{release}

%description onboard-plugin
deepin desktop-environment - dock plugin.

%prep
%setup -q -n %{name}-%{version}
sed -i '/TARGETS/s|lib|%{_lib}|' plugins/*/CMakeLists.txt \
                                 plugins/plugin-guide/plugins-developer-guide.md

sed -i -E '30,39d' CMakeLists.txt

sed -i 's|/lib|/%{_lib}|' frame/controller/dockpluginscontroller.cpp \
                          frame/panel/mainpanelcontrol.cpp \
                          plugins/tray/system-trays/systemtrayscontroller.cpp


sed -i 's|/lib|/libexec|g' plugins/show-desktop/showdesktopplugin.cpp

sed -i 's|/usr/lib/dde-dock/plugins|%{_libdir}/dde-dock/plugins|' plugins/plugin-guide/plugins-developer-guide.md
sed -i 's|local/lib/dde-dock/plugins|local/%{_lib}/dde-dock/plugins|' plugins/plugin-guide/plugins-developer-guide.md

%build
export PATH=%{_qt5_bindir}:$PATH
%cmake -DCMAKE_INSTALL_PREFIX=%{_prefix} -DARCHITECTURE=%{_arch} .
%make_build

%install
%make_install INSTALL_ROOT=%{buildroot}

%ldconfig_scriptlets

%files
%license LICENSE
%{_sysconfdir}/%{name}/indicator/keybord_layout.json
%{_bindir}/%{name}
%{_libdir}/%{name}/
%{_datadir}/%{name}/
%{_datadir}/dbus-1/services/*.service
%{_datarootdir}/glib-2.0/schemas/com.deepin.dde.dock.module.gschema.xml
%{_datarootdir}/polkit-1/actions/com.deepin.dde.dock.overlay.policy

%files devel
%{_includedir}/%{name}/
%{_libdir}/pkgconfig/%{name}.pc
%{_libdir}/cmake/DdeDock/DdeDockConfig.cmake

%files onboard-plugin
%{_libdir}/dde-dock/plugins/libonboard.so


%changelog
* Thu Mar 23 2021 uoser <uoser@uniontech.com> - 5.4.9-1
- Update to 5.4.9

