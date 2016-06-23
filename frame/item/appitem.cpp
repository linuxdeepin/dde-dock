#include "appitem.h"

#include "util/themeappicon.h"

#include <QPainter>
#include <QDrag>
#include <QMouseEvent>

#define APP_DRAG_THRESHOLD      20

int AppItem::IconBaseSize;
QPoint AppItem::MousePressPos;

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(App, parent),
      m_itemEntry(new DBusDockEntry(entry.path(), this)),
      m_draging(false)
{
    setAcceptDrops(true);

    m_id = m_itemEntry->id();
    m_active = m_itemEntry->active();

    connect(m_itemEntry, &DBusDockEntry::ActiveChanged, this, &AppItem::activeChanged);
    connect(m_itemEntry, &DBusDockEntry::TitlesChanged, this, &AppItem::updateTitle);
    connect(m_itemEntry, &DBusDockEntry::ActiveChanged, this, static_cast<void (AppItem::*)()>(&AppItem::update));

    updateTitle();
    updateIcon();
}

const QString AppItem::appId() const
{
    return m_id;
}

void AppItem::setIconBaseSize(const int size)
{
    IconBaseSize = size;
}

int AppItem::iconBaseSize()
{
    return IconBaseSize;
}

int AppItem::itemBaseWidth()
{
    if (DockDisplayMode == Dock::Fashion)
        return itemBaseHeight() * 1.1;
    else
        return itemBaseHeight() * 1.4;
}

int AppItem::itemBaseHeight()
{
    return IconBaseSize * 1.5;
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (m_draging)
        return;

    QPainter painter(this);
    if (!painter.isActive())
        return;

    const QRect itemRect = rect();

    // draw background
    const QRect backgroundRect = itemRect.marginsRemoved(QMargins(1, 2, 1, 2));
    if (DockDisplayMode == Efficient)
    {
        if (m_active)
        {
            painter.fillRect(backgroundRect, QColor(70, 100, 200, 120));

            const int activeLineWidth = 3;
            QRect activeRect = backgroundRect;
            switch (DockPosition)
            {
            case Top:       activeRect.setBottom(activeRect.top() + activeLineWidth);   break;
            case Bottom:    activeRect.setTop(activeRect.bottom() - activeLineWidth);   break;
            case Left:      activeRect.setRight(activeRect.left() + activeLineWidth);   break;
            case Right:     activeRect.setLeft(activeRect.right() - activeLineWidth);   break;
            }

            painter.fillRect(activeRect, QColor(47, 168, 247));
        }
        else if (!m_titles.isEmpty())
            painter.fillRect(backgroundRect, QColor(255, 255, 255, 50));
    //    else
    //        painter.fillRect(backgroundRect, Qt::gray);
    }
    else
    {
        if (!m_titles.isEmpty())
        {
            const int activeLineWidth = 1;
            const int activeLineLength = 20;
            QRect activeRect = itemRect;
            switch (DockPosition)
            {
            case Top:
                activeRect.setBottom(activeRect.top() + activeLineWidth);
                activeRect.setWidth(activeLineLength);
                activeRect.moveLeft((itemRect.width() - activeRect.width()) / 2);
                break;
            case Bottom:
                activeRect.setTop(activeRect.bottom() - activeLineWidth);
                activeRect.setWidth(activeLineLength);
                activeRect.moveLeft((itemRect.width() - activeRect.width()) / 2);
                break;
            case Left:
                activeRect.setRight(activeRect.left() + activeLineWidth);
                activeRect.setHeight(activeLineLength);
                activeRect.moveTop((itemRect.height() - activeRect.height()) / 2);
                break;
            case Right:
                activeRect.setLeft(activeRect.right() - activeLineWidth);
                activeRect.setHeight(activeLineLength);
                activeRect.moveTop((itemRect.height() - activeRect.height()) / 2);
                break;
            }

            painter.fillRect(activeRect, QColor(163, 167, 166));
        }
    }

    // draw icon
    if (DockDisplayMode == Efficient)
        painter.drawPixmap(itemRect.center() - m_smallIcon.rect().center(), m_smallIcon);
    else
        painter.drawPixmap(itemRect.center() - m_largeIcon.rect().center(), m_largeIcon);
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
    if (e->button() == Qt::RightButton)
    {
        if (perfectIconRect().contains(e->pos()))
            return showContextMenu();
        else
            return QWidget::mousePressEvent(e);
    }

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

void AppItem::dragEnterEvent(QDragEnterEvent *e)
{
    // ignore drag from panel
    if (e->source())
        return;

    e->accept();
}

void AppItem::dropEvent(QDropEvent *e)
{
    QStringList uriList;
    for (auto uri : e->mimeData()->urls())
        uriList << uri.toString();

//    qDebug() << uriList;
    m_itemEntry->HandleDragDrop(uriList);
}

void AppItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    Q_UNUSED(checked);

    m_itemEntry->HandleMenuItem(itemId);
}

const QString AppItem::contextMenu() const
{
    return m_itemEntry->menu();
}

void AppItem::startDrag()
{
    m_draging = true;
    update();

    QDrag *drag = new QDrag(this);
    drag->setPixmap(m_smallIcon);
    drag->setHotSpot(m_smallIcon.rect().center());
    drag->setMimeData(new QMimeData);

    emit dragStarted();
    const Qt::DropAction result = drag->exec(Qt::MoveAction);
    Q_UNUSED(result);

    // drag out of dock panel
    if (!drag->target())
        m_itemEntry->RequestUndock();

    m_draging = false;
    setVisible(true);
    update();
}

void AppItem::updateTitle()
{
    m_titles = m_itemEntry->titles();

    update();
}

void AppItem::updateIcon()
{
    const QString icon = m_itemEntry->icon();
    const int iconSize = qMin(width(), height());

    m_smallIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.6);
    m_largeIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.8);
}

void AppItem::activeChanged()
{
    m_active = !m_active;
}
