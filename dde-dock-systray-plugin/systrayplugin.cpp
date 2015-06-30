#include "systrayplugin.h"

SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

QList<AbstractDockItem*> SystrayPlugin::items()
{
    //clear m_items.
    this->clearItems();

    // get xids of trayicons.
    QList<WId> winIds;
    winIds << 79691780 << 65011722;

    // generate items.
    WId winId;
    foreach (winId, winIds) {
        m_items << DockTrayItem::fromWinId(winId);
    }

    return m_items;
}

void SystrayPlugin::clearItems()
{
    AbstractDockItem *item;
    foreach (item, m_items) {
        item->deleteLater();
    }
    m_items.clear();
}
