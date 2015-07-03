#include <QtDBus/QDBusConnection>

#include "systrayplugin.h"

SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

void SystrayPlugin::init(DockPluginProxyInterface * proxier)
{
    m_proxier = proxier;

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
}

QStringList SystrayPlugin::uuids()
{
    return m_items.keys();
}

QWidget * SystrayPlugin::getItem(QString uuid)
{
    return m_items.value(uuid);
}

QString SystrayPlugin::name()
{
    return QString("systray");
}

// private methods
void SystrayPlugin::clearItems()
{
    foreach (QWidget * item, m_items) {
        item->deleteLater();
    }
    m_items.clear();
}

// private slots
void SystrayPlugin::onAdded(WId winId)
{
    QString key = QString::number(winId);

    DockTrayItem *item = DockTrayItem::fromWinId(winId);
    m_items[key] = item;

    m_proxier->itemAddedEvent(key);
}

void SystrayPlugin::onRemoved(WId winId)
{
    QString key = QString::number(winId);

    m_proxier->itemRemovedEvent(key);
}
