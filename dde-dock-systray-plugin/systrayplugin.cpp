#include <QtDBus/QDBusConnection>

#include "systrayplugin.h"
#include "abstractdockitem.h"


SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

void SystrayPlugin::init()
{
    if (!m_dbusTrayManager) {
        m_dbusTrayManager = new com::deepin::dde::TrayManager("com.deepin.dde.TrayManager",
                                                              "/com/deepin/dde/TrayManager",
                                                              QDBusConnection::sessionBus(),
                                                              this);
    }

    QList<uint> trayIcons = m_dbusTrayManager->trayIcons();
    qDebug() << "Found trayicons: " << trayIcons;

    foreach (uint trayIcon, trayIcons) {
        m_items[QString::number(trayIcon)] = DockTrayItem::fromWinId(trayIcon);
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

void SystrayPlugin::clearItems()
{
    foreach (QWidget * item, m_items) {
        item->deleteLater();
    }
    m_items.clear();
}
