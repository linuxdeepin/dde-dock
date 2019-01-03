/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
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

#ifndef TRAYPLUGIN_H
#define TRAYPLUGIN_H

#include "pluginsiteminterface.h"
#include "dbus/dbustraymanager.h"
#include "xwindowtraywidget.h"
#include "indicatortray.h"
#include "indicatortraywidget.h"
#include "snitraywidget.h"
#include "system-trays/systemtrayscontroller.h"
#include "dbus/sni/statusnotifierwatcher_interface.h"

#include <QSettings>
#include <QLabel>

class FashionTrayItem;
class TipsWidget;
class TrayPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "tray.json")

public:
    explicit TrayPlugin(QObject *parent = 0);

    const QString pluginName() const Q_DECL_OVERRIDE;
    void init(PluginProxyInterface *proxyInter) Q_DECL_OVERRIDE;
    void displayModeChanged(const Dock::DisplayMode mode) Q_DECL_OVERRIDE;
    void positionChanged(const Dock::Position position) Q_DECL_OVERRIDE;
    QWidget *itemWidget(const QString &itemKey) Q_DECL_OVERRIDE;
    QWidget *itemTipsWidget(const QString &itemKey) Q_DECL_OVERRIDE;
    QWidget *itemPopupApplet(const QString &itemKey) Q_DECL_OVERRIDE;
    bool itemAllowContainer(const QString &itemKey) Q_DECL_OVERRIDE;
    bool itemIsInContainer(const QString &itemKey) Q_DECL_OVERRIDE;
    int itemSortKey(const QString &itemKey) Q_DECL_OVERRIDE;
    void setSortKey(const QString &itemKey, const int order) Q_DECL_OVERRIDE;
    void setItemIsInContainer(const QString &itemKey, const bool container) Q_DECL_OVERRIDE;
    void refreshIcon(const QString &itemKey) Q_DECL_OVERRIDE;

    Dock::Position dockPosition() const;
    bool traysSortedInFashionMode();
    void saveValue(const QString &key, const QVariant &value);
    const QVariant getValue(const QString &key, const QVariant& fallback = QVariant());

private:
    void loadIndicator();
    bool isSystemTrayItem(const QString &itemKey);
    QString itemKeyOfTrayWidget(AbstractTrayWidget *trayWidget);

private slots:
    void addTrayWidget(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void sniItemsChanged();
    void trayListChanged();
    void trayXWindowAdded(const QString &itemKey, quint32 winId);
    void traySNIAdded(const QString &itemKey, const QString &sniServicePath);
    void trayIndicatorAdded(const QString &itemKey);
    void trayRemoved(const QString &itemKey, const bool deleteObject = true);
    void trayChanged(quint32 winId);
    void switchToMode(const Dock::DisplayMode mode);
    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);
    void onRequestWindowAutoHide(const bool autoHide);
    void onRequestRefershWindowVisible();
    void onSNIItemStatusChanged(SNITrayWidget::ItemStatus status);

private:
    DBusTrayManager *m_trayInter;
    org::kde::StatusNotifierWatcher *m_sniWatcher;
    FashionTrayItem *m_fashionItem;
    SystemTraysController *m_systemTraysController;
    QDBusConnectionInterface *m_dbusDaemonInterface;

    QMap<QString, AbstractTrayWidget *> m_trayMap;
    QMap<QString, IndicatorTray*> m_indicatorMap;
    QString m_sniHostService;

    TipsWidget *m_tipsLabel;
};

#endif // TRAYPLUGIN_H
