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
#include "xembedtraywidget.h"
#include "indicatortray.h"
#include "indicatortraywidget.h"
#include "snitraywidget.h"
#include "system-trays/systemtrayscontroller.h"
#include "dbus/sni/statusnotifierwatcher_interface.h"

#include <QSettings>
#include <QLabel>
#include <mutex>

class FashionTrayItem;
namespace Dock {
class TipsWidget;
}
class TrayPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "tray.json")

public:
    explicit TrayPlugin(QObject *parent = 0);

    const QString pluginName() const override;
    void init(PluginProxyInterface *proxyInter) override;
    bool pluginIsDisable() override;
    void displayModeChanged(const Dock::DisplayMode mode) override;
    void positionChanged(const Dock::Position position) override;
    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void refreshIcon(const QString &itemKey) override;
    void pluginSettingsChanged() override;

    Dock::Position dockPosition() const;
    bool traysSortedInFashionMode();
    void saveValue(const QString &itemKey, const QString &key, const QVariant &value);
    const QVariant getValue(const QString &itemKey, const QString &key, const QVariant& fallback = QVariant());

private:
    void loadIndicator();
    bool isSystemTrayItem(const QString &itemKey);
    QString itemKeyOfTrayWidget(AbstractTrayWidget *trayWidget);
    Dock::DisplayMode displayMode();

private slots:
    void initXEmbed();
    void initSNI();
    void addTrayWidget(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void sniItemsChanged();
    void xembedItemsChanged();
    void trayXEmbedAdded(const QString &itemKey, quint32 winId);
    void traySNIAdded(const QString &itemKey, const QString &sniServicePath);
    void trayIndicatorAdded(const QString &itemKey, const QString &indicatorName);
    void trayRemoved(const QString &itemKey, const bool deleteObject = true);
    void xembedItemChanged(quint32 winId);
    void switchToMode(const Dock::DisplayMode mode);
    void onRequestWindowAutoHide(const bool autoHide);
    void onRequestRefershWindowVisible();
    void onSNIItemStatusChanged(SNITrayWidget::ItemStatus status);

private:
    DBusTrayManager *m_trayInter;
    org::kde::StatusNotifierWatcher *m_sniWatcher;
    FashionTrayItem *m_fashionItem;
    SystemTraysController *m_systemTraysController;
    QTimer *m_refreshXEmbedItemsTimer;
    QTimer *m_refreshSNIItemsTimer;

    QMap<QString, AbstractTrayWidget *> m_trayMap;
    QMap<QString, SNITrayWidget *> m_passiveSNITrayMap;     //这个目前好像无用了
    QMap<QString, IndicatorTray*> m_indicatorMap;           //这个有键盘跟license

    Dock::TipsWidget *m_tipsLabel;
    bool m_pluginLoaded;
    std::mutex        m_sniMutex;
};

#endif // TRAYPLUGIN_H
