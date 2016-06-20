#include "appitem.h"

#include "util/themeappicon.h"

#include <QPainter>
#include <QDrag>
#include <QMouseEvent>

#define APP_DRAG_THRESHOLD      20

QPoint AppItem::MousePressPos;
DBusClientManager *AppItem::ClientInter = nullptr;
//uint AppItem::ActiveWindowId = 0;

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(App, parent),
      m_itemEntry(new DBusDockEntry(entry.path(), this)),
      m_draging(false)
{
    initClientManager();

    m_id = m_itemEntry->id();

    connect(m_itemEntry, &DBusDockEntry::TitlesChanged, this, &AppItem::updateTitle);
    connect(m_itemEntry, &DBusDockEntry::ActiveChanged, this, static_cast<void (AppItem::*)()>(&AppItem::update));

    updateTitle();
    updateIcon();
}

const QString AppItem::appId() const
{
    return m_id;
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (m_draging || !m_itemEntry->isValid())
        return;

    const QRect itemRect = rect();
    const int iconSize = std::min(itemRect.width(), itemRect.height());

    QRect iconRect;
    iconRect.setWidth(iconSize);
    iconRect.setHeight(iconSize);
    iconRect.moveTopLeft(itemRect.center() - iconRect.center());

    QPainter painter(this);

    // draw background
    const QRect backgroundRect = rect().marginsRemoved(QMargins(3, 3, 3, 3));
    if (m_itemEntry->active())
        painter.fillRect(backgroundRect, Qt::blue);
    else if (!m_titles.isEmpty())
        painter.fillRect(backgroundRect, Qt::cyan);
    else
        painter.fillRect(backgroundRect, Qt::gray);

    // draw icon
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);

    // draw text
    painter.setPen(Qt::red);
    painter.drawText(rect(), m_itemEntry->title());
}

void AppItem::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() != Qt::LeftButton)
        return;

    const QPoint distance = MousePressPos - e->pos();
    if (distance.manhattanLength() < APP_DRAG_THRESHOLD)
        m_itemEntry->Activate();
}

void AppItem::mousePressEvent(QMouseEvent *e)
{
    DockItem::mousePressEvent(e);

    if (e->button() == Qt::LeftButton)
        MousePressPos = e->pos();
}

void AppItem::mouseMoveEvent(QMouseEvent *e)
{
    e->accept();

    // handle drag
    if (e->buttons() != Qt::LeftButton)
        return;

    const QPoint distance = e->pos() - MousePressPos;
    if (distance.manhattanLength() < APP_DRAG_THRESHOLD)
        return;

    startDrag();
}

void AppItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    updateIcon();
}

void AppItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    Q_UNUSED(itemId)
    Q_UNUSED(checked)
}

const QString AppItem::contextMenu() const
{
    return m_itemEntry->menu();
}

void AppItem::startDrag()
{
    m_draging = true;
    update();

    QPixmap pixmap(25, 25);
    pixmap.fill(Qt::red);

    QDrag *drag = new QDrag(this);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center());
    drag->setMimeData(new QMimeData);

    emit dragStarted();
    const Qt::DropAction result = drag->exec(Qt::MoveAction);

    qDebug() << "dnd result: " << result;

    m_draging = false;
    update();
    setVisible(true);
}

void AppItem::initClientManager()
{
    if (ClientInter)
        return;

    ClientInter = new DBusClientManager(this);
//    connect(ClientInter, &DBusClientManager::ActiveWindowChanged, [&] (const uint wid) {
//        ActiveWindowId = wid;
//    });
}

void AppItem::updateTitle()
{
    m_titles = m_itemEntry->titles();

    update();
}

void AppItem::updateIcon()
{
    const QString icon = m_itemEntry->icon();
    const int iconSize = qMin(width(), height()) * 0.6;

    m_icon = ThemeAppIcon::getIcon(icon, iconSize);
}
