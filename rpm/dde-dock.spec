%global sname deepin-dock
%global repo dde-dock
%global __provides_exclude_from ^%{_libdir}/%{repo}/.*\\.so$

%if 0%{?fedora}
%global start_logo start-here
Name:           %{sname}
%else
Name:           %{repo}
%endif
Version:        5.4.4
Release:        1%{?fedora:%dist}
Summary:        Deepin desktop-environment - Dock module
License:        GPLv3
%if 0%{?fedora}
URL:            https://github.com/linuxdeepin/dde-dock
Source0:        %{url}/archive/%{version}/%{repo}-%{version}.tar.gz
%else
URL:            http://shuttle.corp.deepin.com/cache/repos/eagle/release-candidate/RERFNS4wLjAuNjU3NQ/pool/main/d/dde-dock/
Source0:        %{name}_%{version}.orig.tar.xz
%endif

BuildRequires:  cmake
BuildRequires:  gcc-c++
BuildRequires:  pkgconfig(dbusmenu-qt5)
BuildRequires:  pkgconfig(dde-network-utils)
BuildRequires:  dtkwidget-devel >= 5.1
BuildRequires:  dtkgui-devel >= 5.2.2.16
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
BuildRequires:  gtest-devel
Requires:       dbusmenu-qt5
%if 0%{?fedora}
BuildRequires:  qt5-qtbase-private-devel
BuildRequires:  make
Requires:       deepin-network-utils
Requires:       deepin-qt-dbus-factory
%else
Requires:       dde-network-utils
Requires:       dde-qt-dbus-factory
%endif
Requires:       xcb-util-wm
Requires:       xcb-util-image

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
%autosetup -p1 -n %{repo}-%{version}
sed -i '/TARGETS/s|lib|%{_lib}|' plugins/*/CMakeLists.txt \
                                 plugins/plugin-guide/plugins-developer-guide.md

sed -i 's|/lib|/%{_lib}|' frame/controller/dockpluginscontroller.cpp \
                          frame/panel/mainpanelcontrol.cpp \
                          plugins/tray/system-trays/systemtrayscontroller.cpp


sed -i 's|/lib|/libexec|g' plugins/show-desktop/showdesktopplugin.cpp \
                           frame/panel/mainpanelcontrol.cpp

sed -i 's:libdir.*:libdir=%{_libdir}:' dde-dock.pc.in

sed -i 's|/usr/lib/dde-dock/plugins|%{_libdir}/dde-dock/plugins|' plugins/plugin-guide/plugins-developer-guide.md
sed -i 's|local/lib/dde-dock/plugins|local/%{_lib}/dde-dock/plugins|' plugins/plugin-guide/plugins-developer-guide.md

%if 0%{?fedora}
# set icon to Fedora logo
sed -i 's|deepin-launcher|%{start_logo}|' frame/item/launcheritem.cpp
%endif

%build
export PATH=%{_qt5_bindir}:$PATH
%if 0%{?fedora}
%cmake -DCMAKE_INSTALL_PREFIX=%{_prefix} -DARCHITECTURE=%{_arch}
%cmake_build
%else
%cmake -DCMAKE_INSTALL_PREFIX=%{_prefix} -DARCHITECTURE=%{_arch} .
%make_build
%endif

%install
%if 0%{?fedora}
%cmake_install
%else
%make_install INSTALL_ROOT=%{buildroot}
%endif

%files
%license LICENSE
%{_sysconfdir}/%{repo}/
%{_bindir}/%{repo}
%{_libdir}/%{repo}/
%{_datadir}/%{repo}/
%{_datarootdir}/glib-2.0/schemas/com.deepin.dde.dock.module.gschema.xml
%{_datarootdir}/polkit-1/actions/com.deepin.dde.dock.overlay.policy

%files devel
%doc plugins/plugin-guide
%{_includedir}/%{repo}/
%{_libdir}/pkgconfig/%{repo}.pc
%{_libdir}/cmake/DdeDock/DdeDockConfig.cmake

%files onboard-plugin
%{_libdir}/dde-dock/plugins/libonboard.so


%changelog
* Wed Jun 10 2020 uoser <uoser@uniontech.com> - 5.1.0.13
- Update to 5.1.0.13

