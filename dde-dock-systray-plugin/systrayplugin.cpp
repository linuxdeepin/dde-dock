#include <QDBusConnection>
#include <QWidget>

#include "systrayplugin.h"
#include "compositetrayitem.h"
#include "trayicon.h"

static const QString CompositeItemKey = "composite_item_key";

SystrayPlugin::SystrayPlugin()
{
    m_compositeItem = new CompositeTrayItem;
}

SystrayPlugin::~SystrayPlugin()
{
    m_compositeItem->deleteLater();
}

void SystrayPlugin::init(DockPluginProxyInterface * proxy)
{
    m_proxy = proxy;
    m_compositeItem->setMode(proxy->dockMode());

    if (!m_dbusTrayManager) {
        m_dbusTrayManager = new com::deepin::dde::TrayManager("com.deepin.dde.TrayManager",
                                                              "/com/deepin/dde/TrayManager",
                                                              QDBusConnection::sessionBus(),
                                                              this);
        connect(m_dbusTrayManager, &TrayManager::Added, this, &SystrayPlugin::onAdded);
        connect(m_dbusTrayManager, &TrayManager::Removed, this, &SystrayPlugin::onRemoved);
    }

    QList<uint> trayIcons = m_dbusTrayManager->trayIcons();
    qDebug() << "Found trayicons: " << trayIcons;

    foreach (uint trayIcon, trayIcons) {
        onAdded(trayIcon);
    }

    m_proxy->itemAddedEvent(CompositeItemKey);
}

QString SystrayPlugin::getPluginName()
{
    return "System Tray";
}

QStringList SystrayPlugin::ids()
{
    return QStringList(CompositeItemKey);
}

QString SystrayPlugin::getTitle(QString)
{
    return "";
}

QString SystrayPlugin::getName(QString)
{
    return getPluginName();
}

QString SystrayPlugin::getCommand(QString)
{
    return "";
}

bool SystrayPlugin::canDisable(QString)
{
    return false;
}

bool SystrayPlugin::isDisabled(QString)
{
    return false;
}

void SystrayPlugin::setDisabled(QString, bool)
{

}

QWidget * SystrayPlugin::getItem(QString)
{
    return m_compositeItem;
}

QWidget * SystrayPlugin::getApplet(QString)
{
    return NULL;
}

void SystrayPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode)
{
    m_compositeItem->setMode(newMode);
    m_proxy->infoChanged(DockPluginInterface::ItemSize, CompositeItemKey);
}

QString SystrayPlugin::getMenuContent(QString)
{
    return "";
}

void SystrayPlugin::invokeMenuItem(QString, QString, bool)
{

}

// private slots
void SystrayPlugin::onAdded(WId winId)
{
    QString key = QString::number(winId);

    TrayIcon * icon = new TrayIcon(winId);

    m_compositeItem->addTrayIcon(key, icon);

    m_proxy->infoChanged(DockPluginInterface::ItemSize, CompositeItemKey);
}

void SystrayPlugin::onRemoved(WId winId)
{
    QString key = QString::number(winId);

    m_compositeItem->remove(key);

    m_proxy->infoChanged(DockPluginInterface::ItemSize, CompositeItemKey);
}
