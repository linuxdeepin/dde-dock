
#include "dockitem.h"
#include "dbus/dbusmenu.h"
#include "dbus/dbusmenumanager.h"

#include <QMouseEvent>
#include <QJsonObject>

DockItem::DockItem(const ItemType type, QWidget *parent)
    : QWidget(parent),
//      m_side(DockSettings::Top),
      m_type(type),

      m_menuManagerInter(new DBusMenuManager(this))
{
}

void DockItem::setDockSide(const DockSide side)
{
    m_side = side;

    update();
}

DockItem::ItemType DockItem::itemType() const
{
    return m_type;
}

void DockItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
}

void DockItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::RightButton)
        return showContextMenu();
}

void DockItem::showContextMenu()
{
    const QString menuJson = contextMenu();
    if (menuJson.isEmpty())
        return;

    QDBusPendingReply<QDBusObjectPath> result = m_menuManagerInter->RegisterMenu();

    result.waitForFinished();
    if (result.isError())
    {
        qWarning() << result.error();
        return;
    }

    const QPoint p = mapToGlobal(pos());
    QJsonObject menuObject;
    menuObject.insert("x", QJsonValue(p.x() + width() / 2));
    menuObject.insert("y", QJsonValue(p.y()));
    menuObject.insert("isDockMenu", QJsonValue(true));
    menuObject.insert("menuJsonContent", QJsonValue(menuJson));

    const QDBusObjectPath path = result.argumentAt(0).value<QDBusObjectPath>();
    DBusMenu *menuInter = new DBusMenu(path.path(), this);

    connect(menuInter, &DBusMenu::ItemInvoked, this, &DockItem::invokedMenuItem);
    connect(menuInter, &DBusMenu::MenuUnregistered, menuInter, &DBusMenu::deleteLater, Qt::QueuedConnection);

    menuInter->ShowMenu(QString(QJsonDocument(menuObject).toJson()));
}

void DockItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    Q_UNUSED(itemId)
    Q_UNUSED(checked)
}

const QString DockItem::contextMenu() const
{
    return QString();
}
