// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "docksettings.h"
#include "settings.h"
#include "common.h"

#include <DConfig>

#include <QDebug>
#include <QJsonObject>
#include <QJsonDocument>
#include <qobjectdefs.h>

DCORE_USE_NAMESPACE

DockSettings::DockSettings(QObject *parent)
 : QObject (parent)
 , m_dockSettings(Settings::ConfigPtr(configDock))
{
    init();
}

void DockSettings::init()
{
    // 绑定属性
    if (m_dockSettings) {
            connect(m_dockSettings, &DConfig::valueChanged, this, [&] (const QString &key) {
                if (key == keyHideMode) {
                    Q_EMIT hideModeChanged(HideModeHandler(m_dockSettings->value(keyHideMode).toString()).toEnum());
                } else if(key == keyDisplayMode) {
                    Q_EMIT displayModeChanged(DisplayMode(DisplayModeHandler(m_dockSettings->value(key).toString()).toEnum()));
                } else if (key == keyPosition) {
                    Q_EMIT positionModeChanged(Position(PositionModeHandler(m_dockSettings->value(key).toString()).toEnum()));
                } else if (key == keyForceQuitApp){
                    Q_EMIT forceQuitAppChanged(ForceQuitAppModeHandler(m_dockSettings->value(key).toString()).toEnum());
                } else if (key == keyShowRecent) {
                    Q_EMIT showRecentChanged(m_dockSettings->value(key).toBool());
                } else if (key == keyShowMultiWindow) {
                    Q_EMIT showMultiWindowChanged(m_dockSettings->value(key).toBool());
                } else if ( key == keyQuickTrayName) {
                    Q_EMIT quickTrayNameChanged(m_dockSettings->value(keyQuickTrayName).toStringList());
                } else if ( key == keyShowWindowName) {
                    Q_EMIT windowNameShowModeChanged(m_dockSettings->value(keyShowWindowName).toInt());
                } else if ( key == keyQuickPlugins) {
                    Q_EMIT quickPluginsChanged(m_dockSettings->value(keyQuickPlugins).toStringList());
                } else if ( key == keyWindowSizeFashion) {
                    Q_EMIT windowSizeFashionChanged(m_dockSettings->value(keyWindowSizeFashion).toUInt());
                } else if ( key == keyWindowSizeEfficient) {
                    Q_EMIT windowSizeEfficientChanged(m_dockSettings->value(keyWindowSizeEfficient).toUInt());
                }
            });
    }
}

HideMode DockSettings::getHideMode()
{
    if (m_dockSettings) {
        QString mode = m_dockSettings->value(keyHideMode).toString();
        HideModeHandler handler(mode);
        return handler.toEnum();
    }
    return HideMode::KeepShowing;
}

void DockSettings::setHideMode(HideMode mode)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyHideMode, HideModeHandler(mode).toString());
    }
}

DisplayMode DockSettings::getDisplayMode()
{
    if (m_dockSettings) {
        QString mode = m_dockSettings->value(keyDisplayMode).toString();
        DisplayModeHandler handler(mode);
        return handler.toEnum();
    }
    return DisplayMode::Efficient;
}

void DockSettings::setDisplayMode(DisplayMode mode)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyDisplayMode, DisplayModeHandler(mode).toString());
    }
}

Position DockSettings::getPositionMode()
{
    Position ret = Position::Bottom;
    if (m_dockSettings) {
        QString mode = m_dockSettings->value(keyPosition).toString();
        PositionModeHandler handler(mode);
        ret = handler.toEnum();
    }
    return ret;
}

void DockSettings::setPositionMode(Position mode)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyPosition, PositionModeHandler(mode).toString());
    }
}

ForceQuitAppMode DockSettings::getForceQuitAppMode()
{
    ForceQuitAppMode ret = ForceQuitAppMode::Enabled;
    if (m_dockSettings) {
        QString mode = m_dockSettings->value(keyForceQuitApp).toString();
        ForceQuitAppModeHandler handler(mode);
        ret = handler.toEnum();
    }
    return ret;
}

void DockSettings::setForceQuitAppMode(ForceQuitAppMode mode)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyForceQuitApp, ForceQuitAppModeHandler(mode).toString());
    }
}

uint DockSettings::getIconSize()
{
    uint size = 36;
    if (m_dockSettings) {
        size = m_dockSettings->value(keyIconSize).toUInt();
    }
    return size;
}

void DockSettings::setIconSize(uint size)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyIconSize, size);
    }
}

uint DockSettings::getShowTimeout()
{
    uint time = 100;
    if (m_dockSettings) {
        time = m_dockSettings->value(keyShowTimeout).toUInt();
    }
    return time;
}

void DockSettings::setShowTimeout(uint time)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyShowTimeout, time);
    }
}

uint DockSettings::getHideTimeout()
{
    uint time = 0;
    if (m_dockSettings) {
        time = m_dockSettings->value(keyHideTimeout).toUInt();
    }
    return time;
}

void DockSettings::setHideTimeout(uint time)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyHideTimeout, time);
    }
}

uint DockSettings::getWindowSizeEfficient()
{
    uint size = 40;
    if (m_dockSettings) {
        size = m_dockSettings->value(keyWindowSizeEfficient).toUInt();
    }
    return size;
}

void DockSettings::setWindowSizeEfficient(uint size)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyWindowSizeEfficient, size);
    }
}

uint DockSettings::getWindowSizeFashion()
{
    uint size = 48;
    if (m_dockSettings) {
        size = m_dockSettings->value(keyWindowSizeFashion).toUInt();
    }
    return size;
}

void DockSettings::setWindowSizeFashion(uint size)
{
    if (m_dockSettings) {
        m_dockSettings->setValue(keyWindowSizeFashion, size);
    }
}

void DockSettings::saveStringList(const QString &key, const QStringList &values)
{
    if (!m_dockSettings)
        return;

    m_dockSettings->setValue(key, values);
}

QStringList DockSettings::loadStringList(const QString &key) const
{
    QStringList ret;
    if (!m_dockSettings)
        return ret;

    for(const auto &var : m_dockSettings->value(key).toList()) {
        if (var.isValid())
            ret.push_back(var.toString());
    }

    return ret;
}

QStringList DockSettings::getDockedApps()
{
    return loadStringList(keyDockedApps);
}

void DockSettings::setDockedApps(const QStringList &apps)
{
    saveStringList(keyDockedApps, apps);
}

QStringList DockSettings::getRecentApps() const
{
    return loadStringList(keyRecentApp);
}

void DockSettings::setRecentApps(const QStringList &apps)
{
    saveStringList(keyRecentApp, apps);
}

QVector<QString> DockSettings::getWinIconPreferredApps()
{
    QVector<QString> ret;
    if (m_dockSettings) {
        for(const auto &var : m_dockSettings->value(keyWinIconPreferredApps).toList()) {
            if (var.isValid())
                ret.push_back(var.toString());
        }
    }

    return ret;
}

void DockSettings::setShowRecent(bool visible)
{
    if (!m_dockSettings)
        return;

    m_dockSettings->setValue(keyShowRecent, visible);
}

bool DockSettings::showRecent() const
{
    if (!m_dockSettings)
        return false;

    return m_dockSettings->value(keyShowRecent).toBool();
}

void DockSettings::setShowMultiWindow(bool showMultiWindow)
{
    if (!m_dockSettings)
        return;

    m_dockSettings->setValue(keyShowMultiWindow, showMultiWindow);
}

bool DockSettings::showMultiWindow() const
{
    if (!m_dockSettings)
        return false;

    return m_dockSettings->value(keyShowMultiWindow).toBool();
}

QString DockSettings::getPluginSettings()
{
    QString ret;
    if (m_dockSettings) {
        ret = m_dockSettings->value(keyPluginSettings).toString();
    }

    qInfo() << "getpluginsettings:" << ret;
    return ret;
}

void DockSettings::setPluginSettings(QString jsonStr)
{
    if (jsonStr.isEmpty())
        return;

    if (m_dockSettings) {
        m_dockSettings->setValue(keyPluginSettings, jsonStr);
    }
}

void DockSettings::mergePluginSettings(QString jsonStr)
{
    QString origin = getPluginSettings();
    QJsonObject originSettings = plguinSettingsStrToObj(origin);
    QJsonObject needMergeSettings = plguinSettingsStrToObj(jsonStr);
    for (auto pluginsIt = needMergeSettings.begin(); pluginsIt != needMergeSettings.end(); pluginsIt++) {
        const QString &pluginName = pluginsIt.key();
        const QJsonObject &needMergeSettingsObj = pluginsIt.value().toObject();
        QJsonObject originSettingsObj = originSettings.value(pluginName).toObject();
        for (auto settingsIt = needMergeSettingsObj.begin(); settingsIt != needMergeSettingsObj.end(); settingsIt++) {
            originSettingsObj.insert(settingsIt.key(), settingsIt.value());
        }

        // 重写plugin对应的设置
        originSettings.remove(pluginName);
        originSettings.insert(pluginName, originSettingsObj);
    }

    setPluginSettings(QJsonDocument(originSettings).toJson(QJsonDocument::JsonFormat::Compact));
}

void DockSettings::removePluginSettings(QString pluginName, QStringList settingkeys)
{
    if (pluginName.isEmpty())
        return;

    QString origin = getPluginSettings();
    QJsonObject originSettings = plguinSettingsStrToObj(origin);
    if (settingkeys.size() == 0) {
        originSettings.remove(pluginName);
    } else {
        for (auto pluginsIt = originSettings.begin(); pluginsIt != originSettings.end(); pluginsIt++) {
            const QString &pluginName = pluginsIt.key();
            if (pluginName != pluginName)
                continue;

            QJsonObject originSettingsObj = originSettings.value(pluginName).toObject();
            for (const auto &key : settingkeys) {
                originSettingsObj.remove(key);
            }

            // 重写plugin对应的设置
            originSettings.remove(pluginName);
            originSettings.insert(pluginName, originSettingsObj);
        }
    }

    setPluginSettings(QJsonDocument(originSettings).toJson(QJsonDocument::JsonFormat::Compact));
}

// 借鉴自dde-dock
QJsonObject DockSettings::plguinSettingsStrToObj(QString jsonStr)
{
    QJsonObject ret;
    const QJsonObject &pluginSettingsObject = QJsonDocument::fromJson(jsonStr.toLocal8Bit()).object();
    if (pluginSettingsObject.isEmpty()) {
        return ret;
    }

    for (auto pluginsIt = pluginSettingsObject.constBegin(); pluginsIt != pluginSettingsObject.constEnd(); ++pluginsIt) {
        const QString &pluginName = pluginsIt.key();
        const QJsonObject &settingsObject = pluginsIt.value().toObject();
        QJsonObject newSettingsObject = ret.value(pluginName).toObject();
        for (auto settingsIt = settingsObject.constBegin(); settingsIt != settingsObject.constEnd(); ++settingsIt) {
            newSettingsObject.insert(settingsIt.key(), settingsIt.value());
        }
        // TODO: remove not exists key-values
        ret.insert(pluginName, newSettingsObject);
    }

    return ret;
}

QStringList DockSettings::getQuickPlugins()
{
    if (!m_dockSettings)
        return QStringList();

    return m_dockSettings->value(keyQuickPlugins).toStringList();
}

void DockSettings::setQuickPlugin(QString plugin)
{
    if (!m_dockSettings)
        return;
    QStringList plugins = m_dockSettings->value(keyQuickPlugins).toStringList();
    if (plugins.contains(plugin)) return;
    m_dockSettings->setValue(keyQuickPlugins, plugins << plugin);
}

void DockSettings::removeQuickPlugin(QString plugin)
{
    if (!m_dockSettings)
        return;
    QStringList plugins = m_dockSettings->value(keyQuickPlugins).toStringList();
    plugins.removeOne(plugin);
    m_dockSettings->setValue(keyQuickPlugins, plugins);
}

void DockSettings::updateQuickPlugins(QStringList plugins)
{
    if (!m_dockSettings)
        return;
    m_dockSettings->setValue(keyQuickPlugins, plugins);
}

QStringList DockSettings::getTrayItemsOnDock()
{
    if (!m_dockSettings)
        return QStringList();
    return m_dockSettings->value(keyQuickTrayName).toStringList();
}

void DockSettings::setTrayItemOnDock(QString item)
{
    if (!m_dockSettings)
        return;
    QStringList items = m_dockSettings->value(keyQuickTrayName).toStringList();
    if (items.contains(item)) return;
    m_dockSettings->setValue(keyQuickTrayName, items << item);
}

void DockSettings::removeTrayItemOnDock(QString item)
{
    if (!m_dockSettings)
        return;
    QStringList items = m_dockSettings->value(keyQuickTrayName).toStringList();
    items.removeOne(item);
    m_dockSettings->setValue(keyQuickTrayName, items);
}

void DockSettings::updateTrayItemsOnDock(QStringList items)
{
    if (!m_dockSettings)
        return;
    m_dockSettings->setValue(keyQuickTrayName, items);
}

int DockSettings::getWindowNameShowMode()
{
    if (!m_dockSettings)
        return 0;
    return m_dockSettings->value(keyShowWindowName).toInt();
}

void DockSettings::setWindowNameShowMode(int value)
{
     if (!m_dockSettings)
        return;
    m_dockSettings->setValue(keyShowWindowName, value);
}
