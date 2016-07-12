#include "mainpanel.h"
#include "item/appitem.h"

#include <QBoxLayout>
#include <QDragEnterEvent>

DockItem *MainPanel::DragingItem = nullptr;

MainPanel::MainPanel(QWidget *parent)
    : QFrame(parent),
      m_position(Dock::Top),
      m_displayMode(Dock::Fashion),
      m_itemLayout(new QBoxLayout(QBoxLayout::LeftToRight)),

      m_itemAdjustTimer(new QTimer(this)),
      m_itemController(DockItemController::instance(this))
{
    m_itemLayout->setSpacing(0);
    m_itemLayout->setContentsMargins(0, 0, 0, 0);

    setAcceptDrops(true);
    setObjectName("MainPanel");
    setStyleSheet("QWidget #MainPanel {"
                  "border:" xstr(PANEL_BORDER) "px solid rgba(162, 162, 162, .2);"
                  "background-color:rgba(10, 10, 10, .6);"
                  "}"
                  // Top
                  "QWidget #MainPanel[displayMode='0'][position='0'] {"
                  "border-bottom-left-radius:5px;"
                  "border-bottom-right-radius:5px;"
                  "}"
                  // Right
                  "QWidget #MainPanel[displayMode='0'][position='1'] {"
                  "border-top-left-radius:5px;"
                  "border-bottom-left-radius:5px;"
                  "}"
                  // Bottom
                  "QWidget #MainPanel[displayMode='0'][position='2'] {"
                  "border-top-left-radius:6px;"
                  "border-top-right-radius:6px;"
                  "}"
                  // Left
                  "QWidget #MainPanel[displayMode='0'][position='3'] {"
                  "border-top-right-radius:5px;"
                  "border-bottom-right-radius:5px;"
                  "}"
                  "QWidget #MainPanel[position='0'] {"
                  "padding:0 " xstr(PANEL_PADDING) "px;"
                  "border-top:none;"
                  "}"
                  "QWidget #MainPanel[position='1'] {"
                  "padding:" xstr(PANEL_PADDING) "px 0;"
                  "border-right:none;"
                  "}"
                  "QWidget #MainPanel[position='2'] {"
                  "padding:0 " xstr(PANEL_PADDING) "px;"
                  "border-bottom:none;"
                  "}"
                  "QWidget #MainPanel[position='3'] {"
                  "padding:" xstr(PANEL_PADDING) "px 0;"
                  "border-left:none;"
                  "}");

    connect(m_itemController, &DockItemController::itemInserted, this, &MainPanel::itemInserted);
    connect(m_itemController, &DockItemController::itemRemoved, this, &MainPanel::itemRemoved, Qt::DirectConnection);
    connect(m_itemController, &DockItemController::itemMoved, this, &MainPanel::itemMoved);
    connect(m_itemAdjustTimer, &QTimer::timeout, this, &MainPanel::adjustItemSize);

    m_itemAdjustTimer->setSingleShot(true);
    m_itemAdjustTimer->setInterval(100);

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        initItemConnection(item);
        m_itemLayout->addWidget(item);
    }

    setLayout(m_itemLayout);
}

void MainPanel::updateDockPosition(const Position dockPosition)
{
    m_position = dockPosition;

    switch (m_position)
    {
    case Position::Top:
    case Position::Bottom:          m_itemLayout->setDirection(QBoxLayout::LeftToRight);    break;
    case Position::Left:
    case Position::Right:           m_itemLayout->setDirection(QBoxLayout::TopToBottom);    break;
    }

    m_itemAdjustTimer->start();
}

void MainPanel::updateDockDisplayMode(const DisplayMode displayMode)
{
    m_displayMode = displayMode;

//    const QList<DockItem *> itemList = m_itemController->itemList();
//    for (auto item : itemList)
//    {
//        if (item->itemType() == DockItem::Placeholder)
//            item->setVisible(displayMode == Dock::Efficient);
//    }

    // reload qss
    setStyleSheet(styleSheet());
}

int MainPanel::displayMode()
{
    return int(m_displayMode);
}

int MainPanel::position()
{
    return int(m_position);
}

void MainPanel::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    m_itemAdjustTimer->start();
}

void MainPanel::dragEnterEvent(QDragEnterEvent *e)
{
    DockItem *dragSourceItem = qobject_cast<DockItem *>(e->source());
    if (!dragSourceItem)
        return;

    e->accept();

    if (dragSourceItem)
        DragingItem->show();
}

void MainPanel::dragMoveEvent(QDragMoveEvent *e)
{
    DockItem *item = itemAt(e->pos());
    if (item == DragingItem)
        return;

    m_itemController->itemMove(DragingItem, item);
}

void MainPanel::dragLeaveEvent(QDragLeaveEvent *e)
{
    Q_UNUSED(e)

    DragingItem->hide();
}

void MainPanel::dropEvent(QDropEvent *e)
{
    Q_UNUSED(e)
}

void MainPanel::initItemConnection(DockItem *item)
{
    connect(item, &DockItem::dragStarted, this, &MainPanel::itemDragStarted);
    connect(item, &DockItem::menuUnregistered, this, &MainPanel::requestRefershWindowVisible);
    connect(item, &DockItem::requestWindowAutoHide, this, &MainPanel::requestWindowAutoHide);
}

DockItem *MainPanel::itemAt(const QPoint &point)
{
    const QList<DockItem *> itemList = m_itemController->itemList();

    for (auto item : itemList)
    {
        QRect rect;
        rect.setTopLeft(item->pos());
        rect.setSize(item->size());

        if (rect.contains(point))
            return item;
    }

    return nullptr;
}

void MainPanel::adjustItemSize()
{
    Q_ASSERT(sender() == m_itemAdjustTimer);

    QSize itemSize;
    switch (m_position)
    {
    case Top:
    case Bottom:
        itemSize.setHeight(height() - PANEL_BORDER);
        itemSize.setWidth(AppItem::itemBaseWidth());
        break;

    case Left:
    case Right:
        itemSize.setHeight(AppItem::itemBaseHeight());
        itemSize.setWidth(width() - PANEL_BORDER);
        break;

    default:
        Q_ASSERT(false);
    }

    if (itemSize.height() < 0 || itemSize.width() < 0)
        return;

    int totalAppItemCount = 0;
    int totalWidth = 0;
    int totalHeight = 0;
    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        QMetaObject::invokeMethod(item, "setVisible", Qt::QueuedConnection, Q_ARG(bool, true));

        switch (item->itemType())
        {
        case DockItem::App:
        case DockItem::Launcher:
            item->setFixedSize(itemSize);
            ++totalAppItemCount;
            totalWidth += itemSize.width();
            totalHeight += itemSize.height();
            break;
        case DockItem::Plugins:
            if (m_displayMode == Fashion)
            {
                item->setFixedSize(itemSize);
                ++totalAppItemCount;
                totalWidth += itemSize.width();
                totalHeight += itemSize.height();
            }
            else
            {
                const QSize size = item->sizeHint();
                item->setFixedSize(size);
                if (m_position == Dock::Top || m_position == Dock::Bottom)
                    item->setFixedHeight(itemSize.height());
                else
                    item->setFixedWidth(itemSize.width());
                totalWidth += size.width();
                totalHeight += size.height();
            }
            break;
        default:;
        }
    }

    const int w = width() - PANEL_BORDER * 2 - PANEL_PADDING * 2;
    const int h = height() - PANEL_BORDER * 2 - PANEL_PADDING * 2;

    // test if panel can display all items completely
    bool containsCompletely = false;
    switch (m_position)
    {
    case Dock::Top:
    case Dock::Bottom:
        containsCompletely = totalWidth <= w;     break;

    case Dock::Left:
    case Dock::Right:
        containsCompletely = totalHeight <= h;   break;

    default:
        Q_ASSERT(false);
    }

    // abort adjust.
    if (containsCompletely)
        return;

    // now, we need to decrease item size to fit panel size
    int overflow;
    int base;
    if (m_position == Dock::Top || m_position == Dock::Bottom)
    {
//        qDebug() << "width: " << totalWidth << width();
        overflow = totalWidth;
        base = w;
    }
    else
    {
//        qDebug() << "height: " << totalHeight << height();
        overflow = totalHeight;
        base = h;
    }

    const int decrease = double(overflow - base) / totalAppItemCount;
    int extraDecrease = overflow - base - decrease * totalAppItemCount;

    for (auto item : itemList)
    {
        if (item->itemType() == DockItem::Placeholder)
            continue;
        if (item->itemType() == DockItem::Plugins)
            if (m_displayMode != Dock::Fashion)
                continue;

        switch (m_position)
        {
        case Dock::Top:
        case Dock::Bottom:
            item->setFixedWidth(item->width() - decrease - bool(extraDecrease));
            break;

        case Dock::Left:
        case Dock::Right:
            item->setFixedHeight(item->height() - decrease - bool(extraDecrease));
            break;
        }

        if (extraDecrease)
            --extraDecrease;
    }

    // ensure all extra space assigned
    Q_ASSERT(extraDecrease == 0);

    update();
}

void MainPanel::itemInserted(const int index, DockItem *item)
{
    // hide new item, display it after size adjust finished
    item->hide();

    initItemConnection(item);
    m_itemLayout->insertWidget(index, item);

    m_itemAdjustTimer->start();
}

void MainPanel::itemRemoved(DockItem *item)
{
    m_itemLayout->removeWidget(item);

    m_itemAdjustTimer->start();
}

void MainPanel::itemMoved(DockItem *item, const int index)
{
    // remove old item
    m_itemLayout->removeWidget(item);
    // insert new position
    m_itemLayout->insertWidget(index, item);
}

void MainPanel::itemDragStarted()
{
    DragingItem = qobject_cast<DockItem *>(sender());

    QRect rect;
    rect.setTopLeft(mapToGlobal(pos()));
    rect.setSize(size());

    DragingItem->setVisible(rect.contains(QCursor::pos()));
}
