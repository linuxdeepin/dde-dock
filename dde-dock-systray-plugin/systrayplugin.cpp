#include <QtDBus/QDBusConnection>

#include "systrayplugin.h"
#include "abstractdockitem.h"

SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

QList<AbstractDockItem*> SystrayPlugin::items()
{
    //clear m_items.
    this->clearItems();

    // get xids of trayicons.
    if (!m_dbusTrayManager) {
        m_dbusTrayManager = new com::deepin::dde::TrayManager("com.deepin.dde.TrayManager",
                                                              "/com/deepin/dde/TrayManager",
                                                              QDBusConnection::sessionBus(),
                                                              this);
    }

    QList<uint> trayIcons = m_dbusTrayManager->trayIcons();
    qDebug() << "Found trayicons: " << trayIcons;

    QList<WId> winIds;
    foreach (QVariant trayIcon, trayIcons) {
        winIds << trayIcon.toUInt();
    }

    // generate items.
    foreach (WId winId, winIds) {
        m_items << DockTrayItem::fromWinId(winId);
    }

    return m_items;
}

void SystrayPlugin::clearItems()
{
    foreach (AbstractDockItem * item, m_items) {
        item->deleteLater();
    }
    m_items.clear();
}
