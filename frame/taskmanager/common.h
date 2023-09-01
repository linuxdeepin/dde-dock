// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef COMMON_H
#define COMMON_H

#include <QMap>
#include <QDir>
#include <QString>
#include <QStandardPaths>

const QString configDock              = "com.deepin.dde.dock";
const QString configAppearance        = "com.deepin.dde.appearance";

const QString keyOpacity              = "Opacity";
const QString keyPosition             = "Position";
const QString keyIconSize             = "Icon_Size";
const QString keyHideMode             = "Hide_Mode";
const QString keyRecentApp            = "Recent_App";
const QString keyShowRecent           = "Show_Recent";
const QString keyDockedApps           = "Docked_Apps";
const QString keyDisplayMode          = "Display_Mode";
const QString keyShowTimeout          = "Show_Timeout";
const QString keyHideTimeout          = "Hide_Timeout";
const QString keyForceQuitApp         = "Force_Quit_App";
const QString keyPluginSettings       = "Plugin_Settings";
const QString keyShowMultiWindow      = "Show_MultiWindow";
const QString keyWindowSizeFashion    = "Window_Size_Fashion";
const QString keyWindowSizeEfficient  = "Window_Size_Efficient";
const QString keyWinIconPreferredApps = "Win_Icon_Preferred_Apps";

constexpr auto DesktopFileActionKey = u8"Desktop Action ";
constexpr auto DDEApplicationManager1ObjectPath = u8"/org/desktopspec/ApplicationManager1";
constexpr auto ApplicationManager1DBusName= u8"org.desktopspec.ApplicationManager1";

static const QString scratchDir = QStandardPaths::writableLocation(QStandardPaths::GenericDataLocation).append("/deepin/dde-dock/scratch/");

const QString desktopHashPrefix = "d:";
const QString windowHashPrefix = "w:";

// 驻留应用desktop file模板
// 由于Icon存储的直接是icon base64压缩后的数据需要“”，防止被desktopfile当成stringlist，从而导致获取icon失败
const QString dockedItemTemplate = R"([Desktop Entry]
Name=%1
Exec=%2
Icon="%3"
Type=Application
Terminal=false
StartupNotify=false
)";

const QString frontendWindowWmClass     = "dde-dock";
const QString ddeLauncherWMClass        = "dde-launcher";

const int smartHideTimerDelay           = 400;
const int configureNotifyDelay          = 100;

const int bestIconSize                  = 48;
const int menuItemHintShowAllWindows    = 1;

const int MotifHintStatus               = 8;
const int MotifHintFunctions            = 1;
const int MotifHintInputMode            = 4;
const int MotifHintDecorations          = 2;

const int MotifFunctionNone             = 0;
const int MotifFunctionAll              = 1;
const int MotifFunctionMove             = 4;
const int MotifFunctionClose            = 32;
const int MotifFunctionResize           = 2;
const int MotifFunctionMinimize         = 8;
const int MotifFunctionMaximize         = 16;


static inline QByteArray sessionType() {
    static QByteArray type = qgetenv("XDG_SESSION_TYPE");
    return type;
}

static inline bool isWaylandSession() {
    return sessionType().compare("wayland") == 0;
}

static inline bool isX11Session() {
    return sessionType().compare("x11") == 0;
}

#endif // COMMON_H
