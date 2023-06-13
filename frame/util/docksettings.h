// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"

#include <QObject>

using namespace Dock;

class HideModeHandler {
    HideMode modeEnum;
    QString modeStr;

public:
    HideModeHandler(HideMode mode) : modeEnum(mode), modeStr("") {}
    HideModeHandler(QString mode) : modeEnum(HideMode::KeepShowing), modeStr(mode) {}

    bool equal(HideModeHandler hideMode) {
        return toString() == hideMode.toString() || toEnum() == hideMode.toEnum();
    }

    QString toString() {
        switch (modeEnum) {
        case HideMode::KeepShowing:
            return "keep-showing";
        case HideMode::KeepHidden:
            return "keep-hidden";
        case HideMode::SmartHide:
            return "smart-hide";
        default:
            return "keep-showing";
        }
        // 默认保持始终显示
    }

    HideMode toEnum() {
        if (modeStr == "keep-hidden")
            return HideMode::KeepHidden;
        if (modeStr == "smart-hide")
            return HideMode::SmartHide;

        return HideMode::KeepShowing;
    }
};

class DisplayModeHandler {
    DisplayMode modeEnum;
    QString modeStr;

public:
    DisplayModeHandler(DisplayMode mode) : modeEnum(mode), modeStr("") {}
    DisplayModeHandler(QString mode) : modeEnum(DisplayMode::Efficient), modeStr(mode) {}

    bool equal(DisplayModeHandler displayMode) {
        return toString() == displayMode.toString() || toEnum() == displayMode.toEnum();
    }

    QString toString() {
        switch (modeEnum) {
        case DisplayMode::Fashion:
            return "fashion";
        case DisplayMode::Efficient:
            return "efficient";
        default:
            return "fashion";
        }
    }

    DisplayMode toEnum() {
        if (modeStr == "fashion")
            return DisplayMode::Fashion;

        return DisplayMode::Efficient;
    }
};


class PositionModeHandler {
    Position modeEnum;
    QString modeStr;

public:
    PositionModeHandler(Position mode) : modeEnum(mode), modeStr("") {}
    PositionModeHandler(QString mode) : modeEnum(Position::Bottom), modeStr(mode) {}

    bool equal(PositionModeHandler displayMode) {
        return toString() == displayMode.toString() || toEnum() == displayMode.toEnum();
    }

    QString toString() {
        switch (modeEnum) {
        case Position::Top:
            return "top";
        case Position::Right:
            return "right";
        case Position::Left:
            return "left";
        case Position::Bottom:
            return "bottom";
        default:
            return "bottom";
        }
    }

    Position toEnum() {
        if (modeStr == "top")
            return Position::Top;
        if (modeStr == "right")
            return Position::Right;
        if (modeStr == "bottom")
            return Position::Bottom;
        if (modeStr == "left")
            return Position::Left;

        return Position::Bottom;
    }
};

// 强制退出应用菜单状态
enum class ForceQuitAppMode {
    Enabled,        // 开启
    Disabled,       // 关闭
    Deactivated,    // 置灰
};

class ForceQuitAppModeHandler {
    ForceQuitAppMode modeEnum;
    QString modeStr;

public:
    ForceQuitAppModeHandler(ForceQuitAppMode mode) : modeEnum(mode), modeStr("") {}
    ForceQuitAppModeHandler(QString mode) : modeEnum(ForceQuitAppMode::Enabled), modeStr(mode) {}

    bool equal(ForceQuitAppModeHandler displayMode) {
        return toString() == displayMode.toString() || toEnum() == displayMode.toEnum();
    }

    QString toString() {
        switch (modeEnum) {
        case ForceQuitAppMode::Enabled:
            return "enabled";
        case ForceQuitAppMode::Disabled:
            return "disabled";
        case ForceQuitAppMode::Deactivated:
            return "deactivated";
        default:
            return "enabled";
        }
    }

    ForceQuitAppMode toEnum() {
        if (modeStr == "disabled")
            return ForceQuitAppMode::Disabled;
        if (modeStr == "deactivated")
            return ForceQuitAppMode::Deactivated;

        return ForceQuitAppMode::Enabled;
    }
};

class Settings;
namespace Dtk {
namespace Core {
class DConfig;
}
}

using namespace Dtk::Core;

// 任务栏组策略配置类
class DockSettings: public QObject
{
    Q_OBJECT

public:
    static inline DockSettings *instance() {
        static DockSettings instance;
        return &instance;
    }
    void init();

    HideMode getHideMode();
    void setHideMode(HideMode mode);
    DisplayMode getDisplayMode();
    void setDisplayMode(DisplayMode mode);
    Position getPositionMode();
    void setPositionMode(Position mode);
    ForceQuitAppMode getForceQuitAppMode();
    void setForceQuitAppMode(ForceQuitAppMode mode);
    uint getIconSize();
    void setIconSize(uint size);
    uint getShowTimeout();
    void setShowTimeout(uint time);
    uint getHideTimeout();
    void setHideTimeout(uint time);
    uint getWindowSizeEfficient();
    void setWindowSizeEfficient(uint size);
    uint getWindowSizeFashion();
    void setWindowSizeFashion(uint size);
    QStringList getDockedApps();
    void setDockedApps(const QStringList &apps);
    QStringList getRecentApps() const;
    void setRecentApps(const QStringList &apps);
    QVector<QString> getWinIconPreferredApps();
    void setShowRecent(bool visible);
    bool showRecent() const;

    QStringList getTrayItemsOnDock();
    void setTrayItemOnDock(QString item);
    void removeTrayItemOnDock(QString item);
    void updateTrayItemsOnDock(QStringList items);

    int  getWindowNameShowMode();
    void setWindowNameShowMode(int value);

    QStringList getQuickPlugins();
    void setQuickPlugin(QString plugin);
    void removeQuickPlugin(QString plugin);
    void updateQuickPlugins(QStringList plugins);


    void setShowMultiWindow(bool showMultiWindow);
    bool showMultiWindow() const;

    // plugin settings
    QString getPluginSettings();
    void setPluginSettings(QString jsonStr);
    void mergePluginSettings(QString jsonStr);
    void removePluginSettings(QString pluginName, QStringList settingkeys);
    QJsonObject plguinSettingsStrToObj(QString jsonStr);

Q_SIGNALS:
    // 隐藏模式改变
    void hideModeChanged(HideMode mode);
    // 显示样式改变
    void displayModeChanged(DisplayMode mode);
    // 显示位置改变
    void positionModeChanged(Position mode);
    // 强制退出应用开关改变
    void forceQuitAppChanged(ForceQuitAppMode mode);
    // 是否显示最近打开应用改变
    void showRecentChanged(bool);
    // 是否显示多开应用改变
    void showMultiWindowChanged(bool);
    // 显示在任务栏上的托盘图标改变
    void quickTrayNameChanged(QStringList items);
    // 小窗口显示窗口名称模式改变
    void windowNameShowModeChanged(int mode);
    // 显示在任务栏上的快捷插件改变
    void quickPluginsChanged(QStringList plugins);
    // 时尚模式下，dock尺寸信息改变
    void windowSizeFashionChanged(uint size);
    // 高效模式dock尺寸改变
    void windowSizeEfficientChanged(uint size);

private:
    DockSettings(QObject *paret = nullptr);
    DockSettings(const DockSettings &);
    DockSettings& operator= (const DockSettings &);
    void saveStringList(const QString &key, const QStringList &values);
    QStringList loadStringList(const QString &key) const;

private:
    DConfig *m_dockSettings;
};

#endif // DOCKSETTINGS_H
