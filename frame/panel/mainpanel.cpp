/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "mainpanel.h"
#include "item/appitem.h"

#include <QBoxLayout>
#include <QDragEnterEvent>

DockItem *MainPanel::DragingItem = nullptr;
PlaceholderItem *MainPanel::RequestDockItem = nullptr;

const char *RequestDockKey = "RequestDock";

MainPanel::MainPanel(QWidget *parent)
    : DBlurEffectWidget(parent),
      m_position(Dock::Top),
      m_displayMode(Dock::Fashion),
      m_itemLayout(new QBoxLayout(QBoxLayout::LeftToRight)),

      m_itemAdjustTimer(new QTimer(this)),
      m_itemController(DockItemController::instance(this))
{
    m_itemLayout->setSpacing(0);
    m_itemLayout->setContentsMargins(0, 0, 0, 0);

    setBlurRectXRadius(0);
    setBlurRectYRadius(0);
    setBlendMode(BehindWindowBlend);

    setAcceptDrops(true);
    setAccessibleName("dock-mainpanel");
    setObjectName("MainPanel");
    setStyleSheet("QWidget #MainPanel {"
//                  "background-color:rgba(10, 10, 10, .6);"
                  "}"
//                  "QWidget #MainPanel[displayMode='1'] {"
//                  "border:none;"
//                  "}"
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

    connect(m_itemController, &DockItemController::itemInserted, this, &MainPanel::itemInserted, Qt::DirectConnection);
    connect(m_itemController, &DockItemController::itemRemoved, this, &MainPanel::itemRemoved, Qt::DirectConnection);
    connect(m_itemController, &DockItemController::itemMoved, this, &MainPanel::itemMoved);
    connect(m_itemController, &DockItemController::itemManaged, this, &MainPanel::manageItem);
    connect(m_itemController, &DockItemController::itemUpdated, m_itemAdjustTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_itemAdjustTimer, &QTimer::timeout, this, &MainPanel::adjustItemSize, Qt::QueuedConnection);

    m_itemAdjustTimer->setSingleShot(true);
    m_itemAdjustTimer->setInterval(100);

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        manageItem(item);
        m_itemLayout->addWidget(item);
    }

    setLayout(m_itemLayout);
}

///
/// \brief MainPanel::updateDockPosition change panel layout with spec position.
/// \param dockPosition
///
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

///
/// \brief MainPanel::updateDockDisplayMode change panel style with spec mode.
/// \param displayMode
///
void MainPanel::updateDockDisplayMode(const DisplayMode displayMode)
{
    m_displayMode = displayMode;

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        // we need to hide container item at fashion mode.
        switch (item->itemType())
        {
        case DockItem::Container:
            item->setVisible(displayMode == Dock::Efficient);
            break;
        default:;
        }
    }

    // reload qss
    setStyleSheet(styleSheet());
}

///
/// \brief MainPanel::displayMode interface for Q_PROPERTY, never use this func.
/// \return
///
int MainPanel::displayMode() const
{
    return int(m_displayMode);
}

///
/// \brief MainPanel::position interface for Q_PROPERTY, never use this func.
/// \return
///
int MainPanel::position() const
{
    return int(m_position);
}

void MainPanel::setEffectEnabled(const bool enabled)
{
    if (enabled)
        setMaskColor(QColor(0, 0, 0, 255 * 0.4));
    else
        setMaskColor(QColor(55, 63, 71));

    m_itemAdjustTimer->start();
}

void MainPanel::moveEvent(QMoveEvent* e)
{
    DBlurEffectWidget::moveEvent(e);

    emit geometryChanged();
}

void MainPanel::resizeEvent(QResizeEvent *e)
{
    DBlurEffectWidget::resizeEvent(e);

    m_itemAdjustTimer->start();
//    m_effectWidget->resize(e->size());

    emit geometryChanged();
}

void MainPanel::dragEnterEvent(QDragEnterEvent *e)
{
    DockItem *item = itemAt(e->pos());
    if (item && item->itemType() == DockItem::Container)
        return;

    DockItem *dragSourceItem = qobject_cast<DockItem *>(e->source());
    if (dragSourceItem)
    {
        e->accept();
        if (DragingItem)
            DragingItem->show();
        return;
    } else {
        DragingItem = nullptr;
    }

    if (!e->mimeData()->formats().contains(RequestDockKey))
        return;
    if (m_itemController->appIsOnDock(e->mimeData()->data(RequestDockKey)))
        return;

    e->accept();
}

void MainPanel::dragMoveEvent(QDragMoveEvent *e)
{
    e->accept();

    DockItem *dst = itemAt(e->pos());
    if (!dst)
        return;

    // internal drag swap
    if (e->source())
    {
        if (dst == DragingItem)
            return;
        if (!DragingItem)
            return;
        if (m_itemController->itemIsInContainer(DragingItem))
            return;

        m_itemController->itemMove(DragingItem, dst);
    } else {
        DragingItem = nullptr;

        if (!RequestDockItem)
        {
            DockItem *insertPositionItem = itemAt(e->pos());
            if (!insertPositionItem)
                return;
            const auto type = insertPositionItem->itemType();
            if (type != DockItem::App && type != DockItem::Stretch)
                return;
            RequestDockItem = new PlaceholderItem;
            m_itemController->placeholderItemAdded(RequestDockItem, insertPositionItem);
        } else {
            if (dst == RequestDockItem)
                return;

            m_itemController->itemMove(RequestDockItem, dst);
        }
    }
}

void MainPanel::dragLeaveEvent(QDragLeaveEvent *e)
{
    Q_UNUSED(e)

    if (RequestDockItem)
    {
        const QRect r(static_cast<QWidget *>(parent())->pos(), size());
        const QPoint p(QCursor::pos());

        if (r.contains(p))
            return;

        m_itemController->placeholderItemRemoved(RequestDockItem);
        RequestDockItem->deleteLater();
        RequestDockItem = nullptr;
    }

    if (DragingItem && DragingItem->itemType() != DockItem::Plugins)
        DragingItem->hide();
}

void MainPanel::dropEvent(QDropEvent *e)
{
    Q_UNUSED(e)

    DragingItem = nullptr;

    if (RequestDockItem)
    {
        m_itemController->placeholderItemDocked(e->mimeData()->data(RequestDockKey), RequestDockItem);
        m_itemController->placeholderItemRemoved(RequestDockItem);
        RequestDockItem->deleteLater();
        RequestDockItem = nullptr;
    }
}

///
/// \brief MainPanel::manageItem manage a dock item, all dock item should be managed after construct.
/// \param item
///
void MainPanel::manageItem(DockItem *item)
{
    connect(item, &DockItem::dragStarted, this, &MainPanel::itemDragStarted, Qt::UniqueConnection);
    connect(item, &DockItem::itemDropped, this, &MainPanel::itemDropped, Qt::UniqueConnection);
    connect(item, &DockItem::requestRefershWindowVisible, this, &MainPanel::requestRefershWindowVisible, Qt::UniqueConnection);
    connect(item, &DockItem::requestWindowAutoHide, this, &MainPanel::requestWindowAutoHide, Qt::UniqueConnection);
}

///
/// \brief MainPanel::itemAt get a dock item which placed at spec point,
/// \param point
/// \return if no item at spec point, return nullptr
///
DockItem *MainPanel::itemAt(const QPoint &point)
{
    const QList<DockItem *> itemList = m_itemController->itemList();

    for (auto item : itemList)
    {
        if (!item->isVisible())
            continue;

        QRect rect;
        rect.setTopLeft(item->pos());
        rect.setSize(item->size());

        if (rect.contains(point))
            return item;
    }

    return nullptr;
}

///
/// \brief MainPanel::adjustItemSize adjust all dock item size to fit panel size,
/// for optimize cpu usage, DO NOT call this func immediately, you should use m_itemAdjustTimer
/// to delay do this operate.
///
void MainPanel::adjustItemSize()
{
    Q_ASSERT(sender() == m_itemAdjustTimer);

    const auto ratio = devicePixelRatioF();

    QSize itemSize;
    switch (m_position)
    {
    case Top:
    case Bottom:
        itemSize.setHeight(height() - PANEL_BORDER);
        itemSize.setWidth(std::round(qreal(AppItem::itemBaseWidth()) / ratio));
        break;

    case Left:
    case Right:
        itemSize.setHeight(std::round(qreal(AppItem::itemBaseHeight()) / ratio));
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
        const auto itemType = item->itemType();
        if (m_itemController->itemIsInContainer(item))
            continue;
        if (m_displayMode == Fashion &&
            itemType == DockItem::Container)
            continue;

        QMetaObject::invokeMethod(item, "setVisible", Qt::QueuedConnection, Q_ARG(bool, true));

        switch (item->itemType())
        {
        case DockItem::Placeholder:
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
        case DockItem::Container:
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
        case DockItem::Stretch:
            break;
        default:
            Q_UNREACHABLE();
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
        containsCompletely = totalWidth <= w;   break;

    case Dock::Left:
    case Dock::Right:
        containsCompletely = totalHeight <= h;  break;

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
        overflow = totalWidth;
        base = w;
    }
    else
    {
        overflow = totalHeight;
        base = h;
    }

    const int decrease = double(overflow - base) / totalAppItemCount;
    int extraDecrease = overflow - base - decrease * totalAppItemCount;

    for (auto item : itemList)
    {
        const DockItem::ItemType itemType = item->itemType();
        if (itemType == DockItem::Stretch || itemType == DockItem::Container)
            continue;
        if (itemType == DockItem::Plugins)
        {
            if (m_displayMode != Dock::Fashion)
                continue;
            if (m_itemController->itemIsInContainer(item))
                continue;
        }

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
}

///
/// \brief MainPanel::itemInserted insert dock item into index position.
/// the new inserted item will be hideen first, and then shown after size
/// adjust finished.
/// \param index
/// \param item
///
void MainPanel::itemInserted(const int index, DockItem *item)
{
    // hide new item, display it after size adjust finished
    item->setVisible(false);
    item->setParent(this);

    manageItem(item);
    m_itemLayout->insertWidget(index, item);

    m_itemAdjustTimer->start();
}

///
/// \brief MainPanel::itemRemoved take out spec item from panel, this function
/// will NOT delete item, and NOT disconnect any signals between item and panel.
/// \param item
///
void MainPanel::itemRemoved(DockItem *item)
{
    m_itemLayout->removeWidget(item);

    m_itemAdjustTimer->start();
}

///
/// \brief MainPanel::itemMoved move item to spec index.
/// the index is start from 0 and counted before remove spec item.
/// \param item
/// \param index
///
void MainPanel::itemMoved(DockItem *item, const int index)
{
    // remove old item
    m_itemLayout->removeWidget(item);
    // insert new position
    m_itemLayout->insertWidget(index, item);
}

///
/// \brief MainPanel::itemDragStarted handle managed item draging
///
void MainPanel::itemDragStarted()
{
    DragingItem = qobject_cast<DockItem *>(sender());

    if (DragingItem->itemType() == DockItem::Plugins)
    {
        if (static_cast<PluginsItem *>(DragingItem)->allowContainer())
        {
            qobject_cast<PluginsItem *>(DragingItem)->hidePopup();
            m_itemController->setDropping(true);
        }
    }

    QRect rect;
    rect.setTopLeft(mapToGlobal(pos()));
    rect.setSize(size());

    DragingItem->setVisible(rect.contains(QCursor::pos()));
}

///
/// \brief MainPanel::itemDropped handle managed item dropped.
/// \param destnation
///
void MainPanel::itemDropped(QObject *destnation)
{
    m_itemController->setDropping(false);

    if (m_displayMode == Dock::Fashion)
        return;

    DockItem *src = qobject_cast<DockItem *>(sender());
//    DockItem *dst = qobject_cast<DockItem *>(destnation);

    if (!src)
        return;

    const bool itemIsInContainer = m_itemController->itemIsInContainer(src);

    // drag from container
    if (itemIsInContainer && src->itemType() == DockItem::Plugins && destnation == this)
        m_itemController->itemDragOutFromContainer(src);

    // drop to container
    if (!itemIsInContainer && src->parent() == this && destnation != this)
        m_itemController->itemDroppedIntoContainer(src);

    m_itemAdjustTimer->start();
}
