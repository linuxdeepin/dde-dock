#include "systemtrayplugin.h"

#include <QWindow>
#include <QWidget>

#define FASHION_MODE_ITEM   "fashion-mode-item"

SystemTrayPlugin::SystemTrayPlugin(QObject *parent)
    : QObject(parent),
      m_trayInter(new DBusTrayManager(this))
{

}

const QString SystemTrayPlugin::pluginName() const
{
    return "system-tray";
}

void SystemTrayPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    connect(m_trayInter, &DBusTrayManager::Added, this, &SystemTrayPlugin::pluginAdded);
    connect(m_trayInter, &DBusTrayManager::Removed, this, &SystemTrayPlugin::pluginRemoved);

    m_trayInter->RetryManager();
}

PluginsItemInterface::ItemType SystemTrayPlugin::pluginType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Complex;
}

QWidget *SystemTrayPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    const quint32 trayWinId = itemKey.toUInt();

    return m_trayList[trayWinId];
}

void SystemTrayPlugin::pluginAdded(const quint32 winId)
{
    if (m_trayList.contains(winId))
        return;

    TrayWidget *trayWidget = new TrayWidget(winId);

    m_trayList[winId] = trayWidget;

    m_proxyInter->itemAdded(this, QString::number(winId));
}

void SystemTrayPlugin::pluginRemoved(const quint32 winId)
{
    TrayWidget *widget = m_trayList[winId];
    if (!widget)
        return;

    m_trayList.remove(winId);

    m_proxyInter->itemRemoved(this, QString::number(winId));
    widget->deleteLater();
}
