// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef COMMON_H
#define COMMON_H

#include <QString>
#include <QMap>
#include <QDir>

const QString configDock = "com.deepin.dde.dock";
const QString keyHideMode             = "Hide_Mode";
const QString keyDisplayMode          = "Display_Mode";
const QString keyPosition             = "Position";
const QString keyIconSize             = "Icon_Size";
const QString keyDockedApps           = "Docked_Apps";
const QString keyShowTimeout          = "Show_Timeout";
const QString keyHideTimeout          = "Hide_Timeout";
const QString keyWindowSizeFashion    = "Window_Size_Fashion";
const QString keyWindowSizeEfficient  = "Window_Size_Efficient";
const QString keyWinIconPreferredApps = "Win_Icon_Preferred_Apps";
const QString keyOpacity              = "Opacity";
const QString keyPluginSettings       = "Plugin_Settings";
const QString keyForceQuitApp         = "Force_Quit_App";
const QString keyRecentApp            = "Recent_App";
const QString keyShowRecent           = "Show_Recent";
const QString keyShowMultiWindow      = "Show_MultiWindow";
const QString keyQuickTrayName       = "Dock_Quick_Tray_Name";
const QString keyShowWindowName      = "Dock_Show_Window_Name";
const QString keyQuickPlugins        = "Dock_Quick_Plugins";

const QString scratchDir = QDir::homePath() + "/.local/dock/scratch/";

const QString windowPatternsFile = "/usr/share/dde/data/window_patterns.json";
const QString desktopHashPrefix = "d:";
const QString windowHashPrefix = "w:";

// 驻留应用desktop file模板
const QString dockedItemTemplate = R"([Desktop Entry]
Name=%1
Exec=%2
Icon=%3
Type=Application
Terminal=false
StartupNotify=false
)";

const QString frontendWindowWmClass = "dde-dock";
const int configureNotifyDelay = 100;
const int smartHideTimerDelay = 400;

const int bestIconSize = 48;
const int menuItemHintShowAllWindows = 1;

const int MotifHintFunctions = 1;
const int MotifHintDecorations = 2;
const int MotifHintInputMode = 4;
const int MotifHintStatus = 8;

const int MotifFunctionNone = 0;
const int MotifFunctionAll = 1;
const int MotifFunctionResize = 2;
const int MotifFunctionMove = 4;
const int MotifFunctionMinimize = 8;
const int MotifFunctionMaximize = 16;
const int MotifFunctionClose = 32;

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
