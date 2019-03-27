/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "overlay-warning-plugin.h"

#include <QIcon>
#include <QStorageInfo>

#include <ddialog.h>

#define PLUGIN_STATE_KEY    "enable"
#define OverlayFileSystemType    "overlay"

DWIDGET_USE_NAMESPACE

OverlayWarningPlugin::OverlayWarningPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
{
}

const QString OverlayWarningPlugin::pluginName() const
{
    return "overlay-warning";
}

const QString OverlayWarningPlugin::pluginDisplayName() const
{
    return QString();
}

QWidget *OverlayWarningPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_warningWidget;
}

QWidget *OverlayWarningPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

void OverlayWarningPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void OverlayWarningPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());

    if (pluginIsDisable())
    {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    }
}

bool OverlayWarningPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString OverlayWarningPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    showCloseOverlayDialog();

    return QString();
}

void OverlayWarningPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (!pluginIsDisable()) {
        m_warningWidget->update();
    }
}

int OverlayWarningPlugin::itemSortKey(const QString &itemKey)
{
    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(mode);

    if (mode == Dock::DisplayMode::Fashion) {
        return m_proxyInter->getValue(this, key, 2).toInt();
    } else {
        return m_proxyInter->getValue(this, key, 5).toInt();
    }
}

void OverlayWarningPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

void OverlayWarningPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "overlay-warning plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_warningWidget = new PluginWidget;

    if (!isOverlayRoot()) {
        return;
    }

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());

    QTimer::singleShot(3000, this, &OverlayWarningPlugin::showCloseOverlayDialog);
}

bool OverlayWarningPlugin::isOverlayRoot()
{
    // ignore live/recovery mode
    if (QString(QFile("/proc/cmdline").readAll()).contains("boot=live")) {
        return false;
    }
    return QString(QStorageInfo::root().fileSystemType()) == OverlayFileSystemType;
}

void OverlayWarningPlugin::showCloseOverlayDialog()
{
    const int result = QProcess::execute("/usr/bin/pkexec overlayroot-disable");

    if (result == 0) {
        QProcess::startDetached("reboot");
    } else {
        qDebug() << "close overlayroot failed, the return code is" << result;
    }
}
