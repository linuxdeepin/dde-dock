#include "appitem.h"

#include "util/themeappicon.h"
#include "util/imagefactory.h"
#include "xcb/xcb_misc.h"

#include <X11/X.h>
#include <X11/Xlib.h>

#include <QPainter>
#include <QDrag>
#include <QMouseEvent>
#include <QApplication>
#include <QHBoxLayout>
#include <QGraphicsScene>
#include <QGraphicsItemAnimation>
#include <QTimeLine>

#define APP_DRAG_THRESHOLD      20

const static qreal Frames[] = { 0, 0.327013, 0.987033, 1.77584, 2.61157, 3.45043, 4.26461, 5.03411, 5.74306, 6.37782, 6.92583, 7.37484, 7.71245, 7.92557, 8, 7.86164, 7.43184, 6.69344, 5.64142, 4.2916, 2.68986, 0.91694, -0.91694, -2.68986, -4.2916, -5.64142, -6.69344, -7.43184, -7.86164, -8, -7.86164, -7.43184, -6.69344, -5.64142, -4.2916, -2.68986, -0.91694, 0.91694, 2.68986, 4.2916, 5.64142, 6.69344, 7.43184, 7.86164, 8, 7.93082, 7.71592, 7.34672, 6.82071, 6.1458, 5.34493, 4.45847, 3.54153, 2.65507, 1.8542, 1.17929, 0.653279, 0.28408, 0.0691776, 0 };

int AppItem::IconBaseSize;
QPoint AppItem::MousePressPos;

AppItem::AppItem(const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(parent),
      m_appNameTips(new QLabel(this)),
      m_appPreviewTips(new PreviewContainer(this)),
      m_itemEntry(new DBusDockEntry(entry.path(), this)),

      m_itemView(new QGraphicsView(this)),
      m_itemScene(new QGraphicsScene(this)),

      m_draging(false),

      m_smallIcon(QPixmap()),
      m_largeIcon(QPixmap()),

      m_horizontalIndicator(QPixmap(":/indicator/resources/indicator.png")),
      m_verticalIndicator(QPixmap(":/indicator/resources/indicator_ver.png")),
      m_activeHorizontalIndicator(QPixmap(":/indicator/resources/indicator_active.png")),
      m_activeVerticalIndicator(QPixmap(":/indicator/resources/indicator_active_ver.png")),
      m_updateIconGeometryTimer(new QTimer(this)),

      m_smallWatcher(new QFutureWatcher<QPixmap>(this)),
      m_largeWatcher(new QFutureWatcher<QPixmap>(this))
{
    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->addWidget(m_itemView);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setAccessibleName(m_itemEntry->name());
    setAcceptDrops(true);
//    setMouseTracking(true);
    setLayout(centralLayout);

    m_itemView->setScene(m_itemScene);
    m_itemView->setAlignment(Qt::AlignCenter);
    m_itemView->setVisible(false);
    m_itemView->setFrameStyle(QFrame::NoFrame);
    m_itemView->setContentsMargins(0, 0, 0, 0);
    m_itemView->setRenderHints(QPainter::SmoothPixmapTransform);
    m_itemView->setViewportUpdateMode(QGraphicsView::SmartViewportUpdate);
    m_itemView->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_itemView->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);

    m_id = m_itemEntry->id();
    m_active = m_itemEntry->active();

    m_appNameTips->setObjectName(m_itemEntry->name());
    m_appNameTips->setAccessibleName(m_itemEntry->name() + "-tips");
    m_appNameTips->setVisible(false);
    m_appNameTips->setStyleSheet("color:white;"
                                 "padding:0px 3px;");

    m_updateIconGeometryTimer->setInterval(500);
    m_updateIconGeometryTimer->setSingleShot(true);

    m_appPreviewTips->setVisible(false);

    connect(m_itemEntry, &DBusDockEntry::ActiveChanged, this, &AppItem::activeChanged);
    connect(m_itemEntry, &DBusDockEntry::TitlesChanged, this, &AppItem::updateTitle, Qt::QueuedConnection);
    connect(m_itemEntry, &DBusDockEntry::IconChanged, this, &AppItem::refershIcon);
    connect(m_itemEntry, &DBusDockEntry::ActiveChanged, this, static_cast<void (AppItem::*)()>(&AppItem::update));

    connect(m_updateIconGeometryTimer, &QTimer::timeout, this, &AppItem::updateWindowIconGeometries, Qt::QueuedConnection);

    connect(m_appPreviewTips, &PreviewContainer::requestActivateWindow, this, &AppItem::requestActivateWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestPreviewWindow, this, &AppItem::requestPreviewWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestCancelPreview, this, &AppItem::requestCancelPreview, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestHidePreview, this, &AppItem::hidePopup, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestCheckWindows, m_itemEntry, &DBusDockEntry::Check);

    updateTitle();
    refershIcon();
}

AppItem::~AppItem()
{
    m_appNameTips->deleteLater();
    m_appPreviewTips->deleteLater();
}

const QString AppItem::appId() const
{
    return m_id;
}

// Update _NET_WM_ICON_GEOMETRY property for windows that every item
// that manages, so that WM can do proper animations for specific
// window behaviors like minimization.
void AppItem::updateWindowIconGeometries()
{
    const QRect r(mapToGlobal(QPoint(0, 0)),
                  mapToGlobal(QPoint(width(),height())));
    auto *xcb_misc = XcbMisc::instance();

    for (auto it(m_titles.cbegin()); it != m_titles.cend(); ++it)
        xcb_misc->set_window_icon_geometry(it.key(), r);
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

void AppItem::moveEvent(QMoveEvent *e)
{
    DockItem::moveEvent(e);

    m_updateIconGeometryTimer->start();
}

int AppItem::itemBaseHeight()
{
    if (DockDisplayMode == Efficient)
        return IconBaseSize * 1.2;
    else
        return IconBaseSize * 1.5;
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (m_draging || m_itemView->isVisible())
        return;

    QPainter painter(this);
    if (!painter.isActive())
        return;
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setRenderHint(QPainter::SmoothPixmapTransform, true);

    const QRectF itemRect = rect();

    // draw background
    QRectF backgroundRect = itemRect;
    if (DockDisplayMode == Efficient)
    {
        switch (DockPosition)
        {
        case Top:
        case Bottom:
//            backgroundRect = itemRect;//.marginsRemoved(QMargins(2, 0, 2, 0));
//            backgroundRect = itemRect.marginsRemoved(QMargins(0, 1, 0, 1));
        case Left:
        case Right:
            backgroundRect = itemRect.marginsRemoved(QMargins(1, 1, 1, 1));
//            backgroundRect = itemRect.marginsRemoved(QMargins(1, 0, 1, 0));
        }
    }

    if (DockDisplayMode == Efficient)
    {
        if (m_active)
        {
            painter.fillRect(backgroundRect, QColor(44, 167, 248, 255 * 0.3));

            const int activeLineWidth = itemRect.height() > 50 ? 4 : 2;
            QRectF activeRect = backgroundRect;
            switch (DockPosition)
            {
            case Top:       activeRect.setBottom(activeRect.top() + activeLineWidth);   break;
            case Bottom:    activeRect.setTop(activeRect.bottom() - activeLineWidth);   break;
            case Left:      activeRect.setRight(activeRect.left() + activeLineWidth);   break;
            case Right:     activeRect.setLeft(activeRect.right() - activeLineWidth);   break;
            }

            painter.fillRect(activeRect, QColor(44, 167, 248, 255));
        }
        else if (!m_titles.isEmpty())
            painter.fillRect(backgroundRect, QColor(255, 255, 255, 255 * 0.2));
    //    else
    //        painter.fillRect(backgroundRect, Qt::gray);
    }
    else
    {
        if (!m_titles.isEmpty())
        {
            QPoint p;
            QPixmap pixmap;
            QPixmap activePixmap;
            switch (DockPosition)
            {
            case Top:
                pixmap = m_horizontalIndicator;
                activePixmap = m_activeHorizontalIndicator;
                p.setX((itemRect.width() - pixmap.width()) / 2);
                p.setY(1);
                break;
            case Bottom:
                pixmap = m_horizontalIndicator;
                activePixmap = m_activeHorizontalIndicator;
                p.setX((itemRect.width() - pixmap.width()) / 2);
                p.setY(itemRect.height() - pixmap.height() - 1);
                break;
            case Left:
                pixmap = m_verticalIndicator;
                activePixmap = m_activeVerticalIndicator;
                p.setX(1);
                p.setY((itemRect.height() - pixmap.height()) / 2);
                break;
            case Right:
                pixmap = m_verticalIndicator;
                activePixmap = m_activeVerticalIndicator;
                p.setX(itemRect.width() - pixmap.width() - 1);
                p.setY((itemRect.height() - pixmap.height()) / 2);
                break;
            }

            if (m_active)
                painter.drawPixmap(p, activePixmap);
            else
                painter.drawPixmap(p, pixmap);
        }
    }

    // icon
    const QPixmap pixmap = DockDisplayMode == Efficient ? m_smallIcon : m_largeIcon;
    if (pixmap.isNull())
        return;

    // icon pos
    const auto ratio = qApp->devicePixelRatio();
    const int iconX = itemRect.center().x() - pixmap.rect().center().x() / ratio;
    const int iconY = itemRect.center().y() - pixmap.rect().center().y() / ratio;

    // draw ligher/normal icon
    if (!m_hover)
        painter.drawPixmap(iconX, iconY, pixmap);
    else
        painter.drawPixmap(iconX, iconY, ImageFactory::lighterEffect(pixmap));
}

void AppItem::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() == Qt::MiddleButton) {
        m_itemEntry->NewInstance();
    } else if (e->button() == Qt::LeftButton) {

        m_itemEntry->Activate();

        if (!m_titles.isEmpty())
            return;

        // start launching effects
        m_itemScene->clear();
        const auto ratio = qApp->devicePixelRatio();
        const QRect r = rect();
        const QPixmap icon = DockDisplayMode == Efficient ? m_smallIcon : m_largeIcon;
        QGraphicsPixmapItem *item = m_itemScene->addPixmap(icon);
        item->setTransformationMode(Qt::SmoothTransformation);
        item->setPos(QPointF(r.center()) - QPointF(icon.rect().center()) / ratio);
        m_itemScene->setSceneRect(r);
        m_itemView->setSceneRect(r);
        m_itemView->setFixedSize(r.size());
//        m_itemView->setSceneRect((r.width() - icon.width()) / 2, (r.height() - icon.height()) / 2, r.width(), r.height());

        QTimeLine *tl = new QTimeLine;
        tl->setDuration(1200);
        tl->setFrameRange(0, 60);
        tl->setLoopCount(1);
        tl->setEasingCurve(QEasingCurve::Linear);
        tl->setStartFrame(0);
        tl->start();

        QGraphicsItemAnimation *ani = new QGraphicsItemAnimation;
        ani->setItem(item);
        ani->setTimeLine(tl);

        const int px = qreal(-icon.rect().center().x()) / ratio;
        const int py = qreal(-icon.rect().center().y()) / ratio - 18.;
        const QPoint pos = r.center() + QPoint(0, 18);
        for (int i(0); i != 60; ++i)
        {
            ani->setPosAt(i / 60.0, pos);
            ani->setTranslationAt(i / 60.0, px, py);
            ani->setRotationAt(i / 60.0, Frames[i]);
        }

        connect(tl, &QTimeLine::finished, tl, &QTimeLine::deleteLater);
        connect(tl, &QTimeLine::finished, ani, &QGraphicsItemAnimation::deleteLater);
        connect(tl, &QTimeLine::finished, m_itemScene, &QGraphicsScene::clear);
        connect(tl, &QTimeLine::finished, m_itemView, &QGraphicsView::hide);

        m_itemView->setVisible(true);
    }
}

void AppItem::mousePressEvent(QMouseEvent *e)
{
    m_updateIconGeometryTimer->stop();
    hidePopup();

    if (e->button() == Qt::RightButton)
    {
        if (perfectIconRect().contains(e->pos()))
        {
            QMetaObject::invokeMethod(this, "showContextMenu", Qt::QueuedConnection);
            return;
        } else {
            return QWidget::mousePressEvent(e);
        }
    }

    if (e->button() == Qt::LeftButton)
        MousePressPos = e->pos();
}

void AppItem::mouseMoveEvent(QMouseEvent *e)
{
    e->accept();

    // handle preview
//    if (e->buttons() == Qt::NoButton)
//        return showPreview();

    // handle drag
    if (e->buttons() != Qt::LeftButton)
        return;

    const QPoint pos = e->pos();
    if (!rect().contains(pos))
        return;

    const QPoint distance = pos - MousePressPos;
    if (distance.manhattanLength() > APP_DRAG_THRESHOLD)
        return startDrag();
}

void AppItem::wheelEvent(QWheelEvent *e)
{
    QWidget::wheelEvent(e);

    m_itemEntry->PresentWindows();
}

void AppItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    refershIcon();
}

void AppItem::dragEnterEvent(QDragEnterEvent *e)
{
    // ignore drag from panel
    if (e->source())
        return;

    // ignore request dock event
    if (e->mimeData()->formats().contains("RequestDock"))
        return e->ignore();

    e->accept();
}

void AppItem::dragMoveEvent(QDragMoveEvent *e)
{
    DockItem::dragMoveEvent(e);

    if (m_titles.isEmpty())
        return;

    if (!PopupWindow->isVisible() || PopupWindow->getContent() != m_appPreviewTips)
        showPreview();
}

void AppItem::dropEvent(QDropEvent *e)
{
    QStringList uriList;
    for (auto uri : e->mimeData()->urls()) {
        uriList << uri.toEncoded();
    }

    qDebug() << "accept drop event with URIs: " << uriList;
    m_itemEntry->HandleDragDrop(uriList);
}

void AppItem::leaveEvent(QEvent *e)
{
    DockItem::leaveEvent(e);

    if (m_appPreviewTips->isVisible())
        m_appPreviewTips->prepareHide();
}

void AppItem::showHoverTips()
{
    if (!m_titles.isEmpty())
        return showPreview();

    DockItem::showHoverTips();
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

QWidget *AppItem::popupTips()
{
    if (m_draging)
        return nullptr;

    if (!m_titles.isEmpty())
    {
        const quint32 currentWindow = m_itemEntry->currentWindow();
        Q_ASSERT(m_titles.contains(currentWindow));
        m_appNameTips->setText(m_titles[currentWindow]);
    } else {
        m_appNameTips->setText(m_itemEntry->name());
    }

    return m_appNameTips;
}

void AppItem::startDrag()
{
    m_draging = true;
    update();

    const QPixmap dragPix = DockDisplayMode == Dock::Fashion ? m_largeIcon : m_smallIcon;

    QDrag *drag = new QDrag(this);
    drag->setPixmap(dragPix);
    drag->setHotSpot(dragPix.rect().center());
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
    m_appPreviewTips->setWindowInfos(m_titles);
    m_updateIconGeometryTimer->start();

    update();
}

void AppItem::refershIcon()
{
    const QString icon = m_itemEntry->icon();
    const int iconSize = qMin(width(), height());

    if (DockDisplayMode == Efficient)
    {
        m_smallIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.7);
        m_largeIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.9);
    } else {
        m_smallIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.6);
        m_largeIcon = ThemeAppIcon::getIcon(icon, iconSize * 0.8);
    }

    update();

    m_updateIconGeometryTimer->start();
}

void AppItem::activeChanged()
{
    m_active = !m_active;
}

void AppItem::showPreview()
{
    if (m_titles.isEmpty())
        return;

    // test cursor position
//    const QRect r = rect();
//    const QPoint p = mapFromGlobal(QCursor::pos());

//    switch (DockPosition)
//    {
//    case Top:       if (p.y() != r.top())       return;     break;
//    case Left:      if (p.x() != r.left())      return;     break;
//    case Right:     if (p.x() != r.right())     return;     break;
//    case Bottom:    if (p.y() != r.bottom())    return;     break;
//    default:        return;
//    }

    m_appPreviewTips->setWindowInfos(m_titles);
    m_appPreviewTips->updateSnapshots();
    m_appPreviewTips->updateLayoutDirection(DockPosition);

    showPopupWindow(m_appPreviewTips, true);
}

