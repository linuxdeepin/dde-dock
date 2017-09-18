/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef SYSTEMTRAYPLUGIN_H
#define SYSTEMTRAYPLUGIN_H

#include "pluginsiteminterface.h"
#include "traywidget.h"
#include "tipswidget.h"
#include "dbus/dbustraymanager.h"

#include <QSettings>
#include <QLabel>

class FashionTrayItem;
class SystemTrayPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "system-tray.json")

public:
    explicit SystemTrayPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);
    void displayModeChanged(const Dock::DisplayMode mode);

    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);
    QWidget *itemPopupApplet(const QString &itemKey);

    bool itemAllowContainer(const QString &itemKey);
    bool itemIsInContainer(const QString &itemKey);
    int itemSortKey(const QString &itemKey);
    void setItemIsInContainer(const QString &itemKey, const bool container);

private:
    void updateTipsContent();
    const QString getWindowClass(quint32 winId);

private slots:
    void trayListChanged();
    void trayAdded(const quint32 winId);
    void trayRemoved(const quint32 winId);
    void trayChanged(const quint32 winId);
    void switchToMode(const Dock::DisplayMode mode);

private:
    DBusTrayManager *m_trayInter;
    FashionTrayItem *m_fashionItem;
    QMap<quint32, TrayWidget *> m_trayList;

    TrayApplet *m_trayApplet;
    QLabel *m_tipsLabel;

    QSettings *m_containerSettings;
};

#endif // SYSTEMTRAYPLUGIN_H
