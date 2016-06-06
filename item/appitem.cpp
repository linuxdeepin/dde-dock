#include "appitem.h"

#include <QPainter>

#define APP_STATUS_KEY          "app-status"
#define APP_ICON_KEY            "icon"
#define APP_MENU_KEY            "menu"
#define APP_XIDS_KEY            "app-xids"

#define APP_ACTIVE_STATUS       "active"
#define APP_NORMAL_STATUS       "normal"

DBusClientManager *AppItem::ClientInter = nullptr;
uint AppItem::ActiveWindowId = 0;

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(App, parent),
      m_itemEntry(new DBusDockEntry(entry.path(), this))
{
    initClientManager();

    m_data = m_itemEntry->data();

    connect(m_itemEntry, static_cast<void (DBusDockEntry::*)(const QString&, const QString&)>(&DBusDockEntry::DataChanged), this, &AppItem::entryDataChanged);
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    const QRect itemRect = rect();
    const int iconSize = std::min(itemRect.width(), itemRect.height());

    QRect iconRect;
    iconRect.setWidth(iconSize);
    iconRect.setHeight(iconSize);
    iconRect.moveTopLeft(itemRect.center() - iconRect.center());

    QPainter painter(this);

    // draw current active background
    if (m_windows.contains(ActiveWindowId))
    {
        painter.fillRect(rect(), Qt::blue);
    } else if (m_data[APP_STATUS_KEY] == APP_ACTIVE_STATUS)
    {
        // draw active background
        painter.fillRect(rect(), Qt::cyan);
    } else {
        // draw normal background
        painter.fillRect(rect(), Qt::gray);
    }

    // draw icon
    painter.fillRect(iconRect, Qt::yellow);

    // draw text
    painter.drawText(rect(), m_itemEntry->id());
}

void AppItem::mouseReleaseEvent(QMouseEvent *e)
{
    Q_UNUSED(e);

    // TODO: dbus signature changed
    m_itemEntry->Activate();
}

void AppItem::initClientManager()
{
    if (ClientInter)
        return;

    ClientInter = new DBusClientManager(this);
    connect(ClientInter, &DBusClientManager::ActiveWindowChanged, [&] (const uint wid) {
        ActiveWindowId = wid;
    });
}

void AppItem::entryDataChanged(const QString &key, const QString &value)
{
    // update data
    m_data[key] = value;

    qDebug() << m_data;

    if (key == APP_STATUS_KEY)
        return update();
}
