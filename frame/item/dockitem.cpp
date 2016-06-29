
#include "dockitem.h"
#include "dbus/dbusmenu.h"
#include "dbus/dbusmenumanager.h"

#include <QMouseEvent>
#include <QJsonObject>

Position DockItem::DockPosition = Position::Top;
DisplayMode DockItem::DockDisplayMode = DisplayMode::Efficient;

DockItem::DockItem(const ItemType type, QWidget *parent)
    : QWidget(parent),
      m_type(type),
      m_hover(false),

      m_menuManagerInter(new DBusMenuManager(this))
{
}

void DockItem::setDockPosition(const Position side)
{
    DockPosition = side;
}

void DockItem::setDockDisplayMode(const DisplayMode mode)
{
    DockDisplayMode = mode;
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

void DockItem::enterEvent(QEvent *e)
{
    m_hover = true;

    update();

    return QWidget::enterEvent(e);
}

void DockItem::leaveEvent(QEvent *e)
{
    m_hover = false;

    update();

    return QWidget::leaveEvent(e);
}

const QRect DockItem::perfectIconRect() const
{
    const QRect itemRect = rect();
    const int iconSize = std::min(itemRect.width(), itemRect.height()) * 0.8;

    QRect iconRect;
    iconRect.setWidth(iconSize);
    iconRect.setHeight(iconSize);
    iconRect.moveTopLeft(itemRect.center() - iconRect.center());

    return iconRect;
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

    QPoint p;
    QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    const QRect r = rect();
    switch (DockPosition)
    {
    case Top:       p += QPoint(r.width() / 2, r.height());      break;
    case Bottom:    p += QPoint(r.width() / 2, 0);               break;
    case Left:      p += QPoint(r.width(), r.height() / 2);      break;
    case Right:     p += QPoint(0, r.height() / 2);              break;
    }

    QJsonObject menuObject;
    menuObject.insert("x", QJsonValue(p.x()));
    menuObject.insert("y", QJsonValue(p.y()));
    menuObject.insert("isDockMenu", QJsonValue(true));
    menuObject.insert("menuJsonContent", QJsonValue(menuJson));

    switch (DockPosition)
    {
    case Top:       menuObject.insert("direction", "top");      break;
    case Bottom:    menuObject.insert("direction", "bottom");   break;
    case Left:      menuObject.insert("direction", "left");     break;
    case Right:     menuObject.insert("direction", "right");    break;
    }

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
