#ifndef SYSTEMTRAYPLUGIN_H
#define SYSTEMTRAYPLUGIN_H

#include "pluginsiteminterface.h"
#include "traywidget.h"
#include "tipswidget.h"
#include "dbus/dbustraymanager.h"

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

private:
    void updateTipsContent();

private slots:
    void trayAdded(const quint32 winId);
    void trayRemoved(const quint32 winId);
    void trayChanged(const quint32 winId);
    void switchToMode(const Dock::DisplayMode mode);

private:
    DBusTrayManager *m_trayInter;
    FashionTrayItem *m_fashionItem;
    QMap<quint32, TrayWidget *> m_trayList;

    TipsWidget *m_tipsWidget;
};

#endif // SYSTEMTRAYPLUGIN_H
