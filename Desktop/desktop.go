package main

import "dlib/dbus"

type DesktopManager struct {
        DockShowMode string

        DesktopIconShowChanged func(iconName string, show bool)
        DockShowModeChanged    func(modeName string)
        ScreenHotAreaChanged   func(areaName, areaAction string)
}

func (desk *DesktopManager) GetDBusInfo () dbus.DBusInfo {
        return dbus.DBusInfo {
                "com.deepin.daemon.desktop",
                "/com/deepin/daemon/desktop",
                "com.deepin.daemon.desktop",
        }
}

func (desk *DesktopManager) GetDesktopShowList() []string {
        return nil
}

func (desk *DesktopManager) SetDesktopShowIcon(iconName string, show bool) bool {
        return true
}

func (desk *DesktopManager) SetDockShowMode(modeName string) bool {
        return true
}

func (desk *DesktopManager) SetScreenHotArea(areaName, areaAction string) bool {
        return true
}

func (desk *DesktopManager) GetScreenAreaAction(areaName string) string {
        return ""
}

func main () {
        desk := DesktopManager {}
        dbus.InstallOnSession (&desk)
        select {}
}
