/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
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
#include <QApplication>
#include <QScreen>
#include <QGraphicsView>

#include <window/mainwindow.h>

#include <item/systemtraypluginitem.h>

static DockItem *DraggingItem = nullptr;
static PlaceholderItem *RequestDockItem = nullptr;

const char *RequestDockKey = "RequestDock";
const char *RequestDockKeyFallback = "text/plain";

const char *DesktopMimeType = "application/x-desktop";

MainPanel::MainPanel(QWidget *parent)
    : DBlurEffectWidget(parent),
      m_position(Dock::Top),
      m_displayMode(Dock::Fashion),
      m_itemLayout(new QBoxLayout(QBoxLayout::LeftToRight)),
      m_showDesktopItem(new ShowDesktopItem(this)),
      m_itemAdjustTimer(new QTimer(this)),
      m_checkMouseLeaveTimer(new QTimer(this)),
      m_itemController(DockItemController::instance(this)),
      m_appDragWidget(nullptr)
{
    m_itemLayout->setSpacing(0);
    m_itemLayout->setContentsMargins(0, 0, 0, 0);

    setBlurRectXRadius(0);
    setBlurRectYRadius(0);
    setBlendMode(BehindWindowBlend);

    setAcceptDrops(true);
    setAccessibleName("dock-mainpanel");
    setObjectName("MainPanel");
    setMouseTracking(true);

    QFile qssFile(":/qss/frame.qss");

    qssFile.open(QFile::ReadOnly);
    if(qssFile.isOpen()) {
        setStyleSheet(qssFile.readAll());
        qssFile.close();
    }

    connect(m_itemController, &DockItemController::itemInserted, this, &MainPanel::itemInserted, Qt::DirectConnection);
    connect(m_itemController, &DockItemController::itemRemoved, this, &MainPanel::itemRemoved, Qt::DirectConnection);
    connect(m_itemController, &DockItemController::itemMoved, this, &MainPanel::itemMoved);
    connect(m_itemController, &DockItemController::itemManaged, this, &MainPanel::manageItem);
    connect(m_itemController, &DockItemController::itemUpdated, m_itemAdjustTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_itemAdjustTimer, &QTimer::timeout, this, &MainPanel::adjustItemSize, Qt::QueuedConnection);
    connect(m_checkMouseLeaveTimer, &QTimer::timeout, this, &MainPanel::checkMouseReallyLeave, Qt::QueuedConnection);
    connect(&DockSettings::Instance(), &DockSettings::opacityChanged, this, &MainPanel::setMaskAlpha);

    m_itemAdjustTimer->setSingleShot(true);
    m_itemAdjustTimer->setInterval(100);

    m_checkMouseLeaveTimer->setSingleShot(true);
    m_checkMouseLeaveTimer->setInterval(300);

    const auto &itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        manageItem(item);
        m_itemLayout->addWidget(item);
    }

    m_showDesktopItem->setFixedSize(10, height());
    m_itemLayout->addWidget(m_showDesktopItem);

    setLayout(m_itemLayout);
}

MainPanel::~MainPanel() { }

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
    case Position::Bottom:
        m_itemLayout->setDirection(QBoxLayout::LeftToRight);
        m_showDesktopItem->setFixedSize(10, height());
        break;
    case Position::Left:
    case Position::Right:
        m_itemLayout->setDirection(QBoxLayout::TopToBottom);
        m_showDesktopItem->setFixedSize(width(), 10);
        break;
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

    m_showDesktopItem->setVisible(displayMode == Dock::Efficient);

    const auto &itemList = m_itemController->itemList();
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
        setMaskColor(DarkColor);
    else
        setMaskColor(QColor(55, 63, 71));

    setMaskAlpha(DockSettings::Instance().Opacity());

    m_itemAdjustTimer->start();
}

bool MainPanel::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == static_cast<QGraphicsView *>(m_appDragWidget)->viewport()) {
        QDropEvent *e = static_cast<QDropEvent *>(event);
        bool isContains = rect().contains(mapFromGlobal(m_appDragWidget->mapToGlobal(e->pos())));
        if (isContains) {
            if (event->type() == QEvent::DragMove) {
                handleDragMove(static_cast<QDragMoveEvent *>(event), true);
            } else if (event->type() == QEvent::Drop) {
                m_appDragWidget->hide();
                return true;
            }
        }
    }
    return false;
}

void MainPanel::moveEvent(QMoveEvent* e)
{
    DBlurEffectWidget::moveEvent(e);

    QTimer::singleShot(500, this, &MainPanel::geometryChanged);
}

void MainPanel::resizeEvent(QResizeEvent *e)
{
    DBlurEffectWidget::resizeEvent(e);

    m_itemAdjustTimer->start();
//    m_effectWidget->resize(e->size());

    QTimer::singleShot(500, this, &MainPanel::geometryChanged);
}

void MainPanel::dragEnterEvent(QDragEnterEvent *e)
{
    // 不知道为什么有可能会收不到dragLeaveEvent，因此使用timer来检测鼠标是否已经离开dock
    m_checkMouseLeaveTimer->start();

    // call dragEnterEvent of MainWindow to show dock when dock is hidden
    static_cast<MainWindow *>(window())->dragEnterEvent(e);

    DockItem *item = itemAt(e->pos());
    if (item && item->itemType() == DockItem::Container)
        return;

    DockItem *dragSourceItem = qobject_cast<DockItem *>(e->source());
    if (dragSourceItem)
    {
        e->accept();
        if (DraggingItem)
            DraggingItem->show();
        return;
    } else {
        DraggingItem = nullptr;
    }

    m_draggingMimeKey = e->mimeData()->formats().contains(RequestDockKey) ? RequestDockKey : RequestDockKeyFallback;

    // dragging item is NOT a desktop file
    if (QMimeDatabase().mimeTypeForFile(e->mimeData()->data(m_draggingMimeKey)).name() != DesktopMimeType) {
        m_draggingMimeKey.clear();
        return;
    }

    // dragging item has been docked
    if (m_itemController->appIsOnDock(e->mimeData()->data(m_draggingMimeKey))) {
        m_draggingMimeKey.clear();
        return;
    }

    e->accept();
}

void MainPanel::dragMoveEvent(QDragMoveEvent *e)
{
    handleDragMove(e, false);
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

    if (DraggingItem) {
        DockItem::ItemType type = DraggingItem->itemType();
        if (type != DockItem::Plugins && type != DockItem::SystemTrayPlugin)
            DraggingItem->hide();
    }
}

void MainPanel::dropEvent(QDropEvent *e)
{
    DraggingItem = nullptr;

    if (RequestDockItem)
    {
        m_itemController->placeholderItemDocked(e->mimeData()->data(m_draggingMimeKey), RequestDockItem);
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
    const auto &itemList = m_itemController->itemList();

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

    // ensure all item is update, whatever layout is changed
    QTimer::singleShot(1, this, static_cast<void (MainPanel::*)()>(&MainPanel::update));

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
    const auto &itemList = m_itemController->itemList();

    // FSTray: FashionSystemTray
    const QSize &FSTrayTotalSize = DockSettings::Instance().fashionSystemTraySize(); // the total size of FSTray
    SystemTrayPluginItem *FSTrayItem = nullptr; // the FSTray item object
    QSize FSTraySuggestIconSize = itemSize; // the suggested size of FStray icons

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
        case DockItem::SystemTrayPlugin:
            if (m_displayMode == Fashion) {
                // 特殊处理时尚模式下的托盘插件
                if (item->itemType() == DockItem::SystemTrayPlugin) {
                    FSTrayItem = static_cast<SystemTrayPluginItem *>(item.data());
                    if (m_position == Dock::Top || m_position == Dock::Bottom) {
                        item->setFixedWidth(FSTrayTotalSize.width());
                        item->setFixedHeight(itemSize.height());
                        totalWidth += FSTrayTotalSize.width();
                        totalHeight += itemSize.height();
                    } else {
                        item->setFixedWidth(itemSize.width());
                        item->setFixedHeight(FSTrayTotalSize.height());
                        totalWidth += itemSize.width();
                        totalHeight += FSTrayTotalSize.height();
                    }
                } else {
                    item->setFixedSize(itemSize);
                    totalWidth += itemSize.width();
                    totalHeight += itemSize.height();
                    ++totalAppItemCount;
                }
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

    const int w = width() - PANEL_BORDER * 2 - PANEL_PADDING * 2 - PANEL_MARGIN * 2;
    const int h = height() - PANEL_BORDER * 2 - PANEL_PADDING * 2 - PANEL_MARGIN * 2;

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
    if (containsCompletely) {
        if (FSTrayItem) {
            FSTrayItem->setSuggestIconSize(FSTraySuggestIconSize);
        }
        return;
    }

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

    // FIXME:
    // 时尚模式下使用整形否则会出现图标大小计算不正确的问题
    // 高校模式下使用浮点数否则会出现图标背景色连到一起的问题
    const double decrease = m_displayMode == Dock::Fashion ?
                int(overflow - base) / totalAppItemCount :
                double(overflow - base) / totalAppItemCount;

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
        if (itemType == DockItem::SystemTrayPlugin) {
            if (m_displayMode == Dock::Fashion) {
                switch (m_position) {
                case Dock::Top:
                case Dock::Bottom:
                    FSTraySuggestIconSize.setWidth(itemSize.width() - decrease);
                    break;

                case Dock::Left:
                case Dock::Right:
                    FSTraySuggestIconSize.setHeight(itemSize.height() - decrease);
                    break;
                }
            }
            continue;
        }

        switch (m_position) {
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

    // 如果dock的大小已经是最大的则不再调整时尚模式托盘图标的大小,以避免递归调整dock与托盘的大小
    if (!DockSettings::Instance().isMaxSize() && FSTrayItem) {
        FSTrayItem->setSuggestIconSize(FSTraySuggestIconSize);
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
/// \brief MainPanel::itemDragStarted handle managed item dragging
///
void MainPanel::itemDragStarted()
{
    DraggingItem = qobject_cast<DockItem *>(sender());

    DockItem::ItemType draggingTyep = DraggingItem->itemType();
    if (draggingTyep == DockItem::App)
    {
        AppItem *appItem = qobject_cast<AppItem *>(DraggingItem);
        m_appDragWidget = appItem->appDragWidget();
        appItem->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));
        static_cast<QGraphicsView *>(m_appDragWidget)->viewport()->installEventFilter(this);
    }

    if (draggingTyep == DockItem::Plugins || draggingTyep == DockItem::SystemTrayPlugin)
    {
        if (static_cast<PluginsItem *>(DraggingItem)->allowContainer())
        {
            qobject_cast<PluginsItem *>(DraggingItem)->hidePopup();
            m_itemController->setDropping(true);
        }
    }

    QRect rect;
    rect.setTopLeft(mapToGlobal(pos()));
    rect.setSize(size());

    DraggingItem->setVisible(rect.contains(QCursor::pos()));
}

///
/// \brief MainPanel::itemDropped handle managed item dropped.
/// \param destnation
///
void MainPanel::itemDropped(QObject *destnation)
{
    m_itemController->setDropping(false);

    DockItem *src = qobject_cast<DockItem *>(sender());
//    DockItem *dst = qobject_cast<DockItem *>(destnation);

    if (m_displayMode == Dock::Fashion)
        return;

    if (!src)
        return;

    const bool itemIsInContainer = m_itemController->itemIsInContainer(src);

    // drag from container
    if (itemIsInContainer
            && (src->itemType() == DockItem::Plugins || src->itemType() == DockItem::SystemTrayPlugin)
            && destnation == this)
        m_itemController->itemDragOutFromContainer(src);

    // drop to container
    if (!itemIsInContainer && src->parent() == this && destnation != this)
        m_itemController->itemDroppedIntoContainer(src);

    m_itemAdjustTimer->start();
}

void MainPanel::handleDragMove(QDragMoveEvent *e, bool isFilter)
{
    e->accept();

    DockItem *dst = itemAt(isFilter ? mapFromGlobal(m_appDragWidget->mapToGlobal(e->pos())) : e->pos());

    if (!dst)
        return;

    // internal drag swap
    if (e->source())
    {
        if (dst == DraggingItem)
            return;
        if (!DraggingItem)
            return;
        if (m_itemController->itemIsInContainer(DraggingItem))
            return;

        m_itemController->itemMove(DraggingItem, dst);
    } else {
        DraggingItem = nullptr;

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

void MainPanel::checkMouseReallyLeave()
{
    if (window()->geometry().contains(QCursor::pos())) {
        return m_checkMouseLeaveTimer->start();
    }

    m_checkMouseLeaveTimer->stop();

    dragLeaveEvent(new QDragLeaveEvent);
}
