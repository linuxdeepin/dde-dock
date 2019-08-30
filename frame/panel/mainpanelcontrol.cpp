/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     xuwenw <xuwenw@xuwenw.so>
 *
 * Maintainer: xuwenw <xuwenw@xuwenw.so>
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

#include "mainpanelcontrol.h"
#include "../item/dockitem.h"
#include "util/docksettings.h"
#include "../item/placeholderitem.h"
#include "../item/components/appdrag.h"
#include "../item/appitem.h"

#include <DAnchors>

#include <QDrag>
#include <QTimer>

DWIDGET_USE_NAMESPACE

MainPanelControl::MainPanelControl(QWidget *parent)
    : QWidget(parent)
    , m_mainPanelLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_fixedAreaWidget(new QWidget(this))
    , m_appAreaWidget(new QWidget(this))
    , m_trayAreaWidget(new QWidget(this))
    , m_pluginAreaWidget(new QWidget(this))
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_trayAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_pluginLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_position(Position::Top)
    , m_placeholderItem(nullptr)
    , m_appDragWidget(nullptr)
    , m_dislayMode(Efficient)
{
    init();
    updateMainPanelLayout();
    updateDisplayMode();
    setAcceptDrops(true);

    connect(this, SIGNAL(displayModeChanged()), this, SLOT(onDisplayModeChanged()));
    connect(this, SIGNAL(positionChanged()), this, SLOT(onPositionChanged()));
}

MainPanelControl::~MainPanelControl()
{
}

void MainPanelControl::init()
{
    m_mainPanelLayout->setMargin(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);
    m_fixedAreaLayout->setMargin(0);
    m_fixedAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_pluginLayout->setMargin(0);
    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_trayAreaLayout->setMargin(0);
    m_trayAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_appAreaSonLayout->setMargin(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_fixedAreaWidget);
    m_fixedAreaWidget->setLayout(m_fixedAreaLayout);
    m_mainPanelLayout->addWidget(m_appAreaWidget);
    m_mainPanelLayout->addWidget(m_trayAreaWidget);
    m_trayAreaWidget->setLayout(m_trayAreaLayout);
    m_mainPanelLayout->addWidget(m_pluginAreaWidget);
    m_pluginAreaWidget->setLayout(m_pluginLayout);
    m_appAreaSonWidget->setLayout(m_appAreaSonLayout);
    m_fixedAreaLayout->setSpacing(0);
    m_appAreaSonLayout->setSpacing(0);
    m_trayAreaLayout->setSpacing(0);
    m_pluginLayout->setSpacing(0);
}

void MainPanelControl::setDisplayMode(const DisplayMode mode)
{
    if (mode == m_dislayMode)
        return;
    m_dislayMode = mode;
    emit displayModeChanged();
}

void MainPanelControl::updateMainPanelLayout()
{
    switch (m_position) {
    case Position::Top:
    case Position::Bottom:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_mainPanelLayout->setDirection(QBoxLayout::LeftToRight);
        m_fixedAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_pluginLayout->setDirection(QBoxLayout::LeftToRight);
        m_trayAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_appAreaSonLayout->setDirection(QBoxLayout::LeftToRight);
        break;
    case Position::Right:
    case Position::Left:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_mainPanelLayout->setDirection(QBoxLayout::TopToBottom);
        m_fixedAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_pluginLayout->setDirection(QBoxLayout::TopToBottom);
        m_trayAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_appAreaSonLayout->setDirection(QBoxLayout::TopToBottom);
        break;
    default:
        break;
    }

    QTimer::singleShot(0, this, &MainPanelControl::updateAppAreaSonWidgetSize);
}

void MainPanelControl::addFixedAreaItem(const int index, QWidget *wdg)
{
    m_fixedAreaLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addAppAreaItem(const int index, QWidget *wdg)
{
    m_appAreaSonLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addTrayAreaItem(const int index, QWidget *wdg)
{
    m_trayAreaLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::addPluginAreaItem(const int index, QWidget *wdg)
{
    m_pluginLayout->insertWidget(index, wdg, 0, Qt::AlignCenter);
}

void MainPanelControl::removeFixedAreaItem(QWidget *wdg)
{
    m_fixedAreaLayout->removeWidget(wdg);
}

void MainPanelControl::removeAppAreaItem(QWidget *wdg)
{
    m_appAreaSonLayout->removeWidget(wdg);
}

void MainPanelControl::removeTrayAreaItem(QWidget *wdg)
{
    m_trayAreaLayout->removeWidget(wdg);
}

void MainPanelControl::removePluginAreaItem(QWidget *wdg)
{
    m_pluginLayout->removeWidget(wdg);
}

void MainPanelControl::resizeEvent(QResizeEvent *event)
{
    updateAppAreaSonWidgetSize();

    return QWidget::resizeEvent(event);
}

void MainPanelControl::updateAppAreaSonWidgetSize()
{
    for (int i = 0; i < m_appAreaSonLayout->count(); ++i) {
        QWidget *w = m_appAreaSonLayout->itemAt(i)->widget();
        if (w) {
            if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
                w->setMaximumSize(height(), height());
            } else {
                w->setMaximumSize(width(), width());
            }
        }
    }

    for (int i = 0; i < m_fixedAreaLayout->count(); ++i) {
        QWidget *w = m_fixedAreaLayout->itemAt(i)->widget();
        if (w) {
            if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
                w->setMaximumSize(height(), height());
            } else {
                w->setMaximumSize(width(), width());
            }
        }
    }

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        m_appAreaSonWidget->setMaximumWidth(qMin((m_appAreaWidget->geometry().right() - width() / 2) * 2, m_appAreaWidget->width()));
    } else {
        m_appAreaSonWidget->setMaximumHeight(qMin((m_appAreaWidget->geometry().bottom() - height() / 2) * 2, m_appAreaWidget->height()));
    }

    m_appAreaSonWidget->adjustSize();
}

void MainPanelControl::setPositonValue(const Position position)
{
    if (m_position == position)
        return;
    m_position = position;
    emit positionChanged();
}

void MainPanelControl::insertItem(const int index, DockItem *item)
{
    item->installEventFilter(this);

    switch (item->itemType()) {
    case DockItem::Launcher:
        addFixedAreaItem(index, item);
        break;
    case DockItem::App:
    case DockItem::Placeholder:
        addAppAreaItem(index, item);
        break;
    case DockItem::TrayPlugin:
        addTrayAreaItem(index, item);
        break;
    case DockItem::Plugins:
        addPluginAreaItem(index, item);
        break;
    default:
        break;
    }

    updateMainPanelLayout();
    updateAppAreaSonWidgetSize();
}

void MainPanelControl::removeItem(DockItem *item)
{
    switch (item->itemType()) {
    case DockItem::Launcher:
        removeFixedAreaItem(item);
        break;
    case DockItem::App:
    case DockItem::Placeholder:
        removeAppAreaItem(item);
        break;
    case DockItem::TrayPlugin:
        removeTrayAreaItem(item);
        break;
    case DockItem::Plugins:
        removePluginAreaItem(item);
        break;
    default:
        break;
    }

    updateAppAreaSonWidgetSize();
}

MainPanelDelegate *MainPanelControl::delegate() const
{
    return m_delegate;
}

void MainPanelControl::setDelegate(MainPanelDelegate *delegate)
{
    m_delegate = delegate;
}

void MainPanelControl::moveItem(DockItem *sourceItem, DockItem *targetItem)
{
    // get target index
    int idx = -1;
    if (targetItem->itemType() == DockItem::App)
        idx = m_appAreaSonLayout->indexOf(targetItem);
    else if (targetItem->itemType() == DockItem::Plugins)
        idx = m_pluginLayout->indexOf(targetItem);
    else
        return;

    // remove old item
    removeItem(sourceItem);

    // insert new position
    insertItem(idx, sourceItem);
}

void MainPanelControl::dragEnterEvent(QDragEnterEvent *e)
{
    DockItem *sourceItem = qobject_cast<DockItem *>(e->source());
    if (sourceItem) {
        e->accept();
        return;
    }

    // 拖app到dock上
    const char *RequestDockKey = "RequestDock";
    const char *RequestDockKeyFallback = "text/plain";
    const char *DesktopMimeType = "application/x-desktop";

    m_draggingMimeKey = e->mimeData()->formats().contains(RequestDockKey) ? RequestDockKey : RequestDockKeyFallback;

    // dragging item is NOT a desktop file
    if (QMimeDatabase().mimeTypeForFile(e->mimeData()->data(m_draggingMimeKey)).name() != DesktopMimeType) {
        m_draggingMimeKey.clear();
        qDebug() << "dragging item is NOT a desktop file";
        return;
    }

    if (m_delegate && m_delegate->appIsOnDock(e->mimeData()->data(m_draggingMimeKey)))
        return;

    e->accept();
}

void MainPanelControl::dragLeaveEvent(QDragLeaveEvent *e)
{
    if (m_placeholderItem) {
        const QRect r(static_cast<QWidget *>(parent())->pos(), size());
        const QPoint p(QCursor::pos());

        // remove margins to fix a touch screen bug:
        // the mouse point position will stay on this rect's margins after
        // drag move to the edge of screen
        if (r.marginsRemoved(QMargins(1, 10, 1, 1)).contains(p))
            return;

        removeAppAreaItem(m_placeholderItem);
        m_placeholderItem->deleteLater();
        updateMainPanelLayout();
    }
}

void MainPanelControl::dropEvent(QDropEvent *e)
{
    if (m_placeholderItem) {

        emit itemAdded(e->mimeData()->data(m_draggingMimeKey), m_appAreaSonLayout->indexOf(m_placeholderItem));

        removeAppAreaItem(m_placeholderItem);
        m_placeholderItem->deleteLater();
    }
}

void MainPanelControl::handleDragMove(QDragMoveEvent *e, bool isFilter)
{
    if (!e->source()) {
        // 应用程序拖到dock上
        e->accept();

        DockItem *insertPositionItem = dropTargetItem(nullptr, e->pos());

        if (m_placeholderItem.isNull()) {

            m_placeholderItem = new PlaceholderItem;

            if (m_position == Qt::TopEdge || m_position == Qt::BottomEdge) {
                if (m_appAreaSonWidget->mapFromParent(e->pos()).x() > m_appAreaSonWidget->rect().right()) {
                    // 插入到最右侧
                    insertPositionItem = nullptr;
                }
            } else {
                if (m_appAreaSonWidget->mapFromParent(e->pos()).y() > m_appAreaSonWidget->rect().bottom()) {
                    // 插入到最下测
                    insertPositionItem = nullptr;
                }
            }

            insertItem(m_appAreaSonLayout->indexOf(insertPositionItem), m_placeholderItem);

        } else if (insertPositionItem && m_placeholderItem != insertPositionItem) {
            moveItem(m_placeholderItem, insertPositionItem);
        }

        return;
    }

    DockItem *sourceItem = qobject_cast<DockItem *>(e->source());

    if (!sourceItem) {
        e->ignore();
        return;
    }

    DockItem *targetItem = nullptr;

    if (isFilter) {
        // appItem调整顺序或者移除驻留
        targetItem = dropTargetItem(sourceItem, mapFromGlobal(m_appDragWidget->mapToGlobal(e->pos())));

        if (targetItem) {
            m_appDragWidget->setOriginPos((m_appAreaSonWidget->mapToGlobal(targetItem->pos())));
        } else {
            targetItem = sourceItem;
        }

    } else {
        // other dockItem调整顺序
        targetItem = dropTargetItem(sourceItem, e->pos());
    }

    if (!targetItem) {
        e->ignore();
        return;
    }

    e->accept();

    if (targetItem == sourceItem)
        return;

    moveItem(sourceItem, targetItem);
    emit itemMoved(sourceItem, targetItem);
}

void MainPanelControl::dragMoveEvent(QDragMoveEvent *e)
{
    handleDragMove(e, false);
}

bool MainPanelControl::eventFilter(QObject *watched, QEvent *event)
{
    if (m_appDragWidget && watched == static_cast<QGraphicsView *>(m_appDragWidget)->viewport()) {
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
        return false;
    }

    if (event->type() != QEvent::MouseMove)
        return false;

    QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
    if (!mouseEvent || mouseEvent->buttons() != Qt::LeftButton)
        return false;

    DockItem *item = qobject_cast<DockItem *>(watched);
    if (!item)
        return false;

    if (item->itemType() != DockItem::App && item->itemType() != DockItem::Plugins)
        return false;

    startDrag(item);

    return true;
}

void MainPanelControl::startDrag(DockItem *item)
{
    const QPixmap pixmap = item->grab();

    item->setDraging(true);
    item->update();

    QDrag *drag = nullptr;
    if (item->itemType() == DockItem::App) {
        AppDrag *appDrag = new AppDrag(item);
        appDrag->setPixmap(pixmap);
        m_appDragWidget = appDrag->appDragWidget();

        connect(m_appDragWidget, &AppDragWidget::destroyed, this, [ = ] {
            m_appDragWidget = nullptr;
        });

        m_appDragWidget->show();

        Dock::Position position;
        switch (m_position) {
        case Position::Top: position = Dock::Top; break;
        case Position::Bottom: position = Dock::Bottom; break;
        case Position::Left: position = Dock::Left; break;
        case Position::Right: position = Dock::Right; break;
        }

        appDrag->appDragWidget()->setOriginPos((m_appAreaSonWidget->mapToGlobal(item->pos())));
        appDrag->appDragWidget()->setDockInfo(position, QRect(mapToGlobal(pos()), size()));
        static_cast<QGraphicsView *>(m_appDragWidget)->viewport()->installEventFilter(this);

        drag = appDrag;
    } else {
        drag = new QDrag(item);
        drag->setPixmap(pixmap);
    }
    drag->setHotSpot(pixmap.rect().center() / pixmap.devicePixelRatioF());
    drag->setMimeData(new QMimeData);
    drag->exec(Qt::MoveAction);

    m_appDragWidget = nullptr;
    item->setDraging(false);
    item->update();
}

DockItem *MainPanelControl::dropTargetItem(DockItem *sourceItem, QPoint point)
{
    QWidget *parentWidget = m_appAreaSonWidget;

    if (sourceItem) {
        switch (sourceItem->itemType()) {
        case DockItem::App:
            parentWidget = m_appAreaSonWidget;
            break;
        case DockItem::Plugins:
            parentWidget = m_pluginAreaWidget;
            break;
        default:
            break;
        }
    }

    if (!parentWidget)
        return nullptr;

    point = parentWidget->mapFromParent(point);
    QLayout *parentLayout = parentWidget->layout();

    DockItem *targetItem = nullptr;

    for (int i = 0 ; i < parentLayout->count(); ++i) {
        QLayoutItem *layoutItem = parentLayout->itemAt(i);
        DockItem *dockItem = qobject_cast<DockItem *>(layoutItem->widget());
        if (!dockItem)
            continue;

        QRect rect;
        rect.setTopLeft(dockItem->pos());
        rect.setSize(dockItem->size());

        if (rect.contains(point)) {
            targetItem = dockItem;
            break;
        }
    }

    if (!targetItem && parentWidget == m_appAreaSonWidget) {
        // appitem调整顺序是，判断是否拖放在两边空白区域

        if (!m_appAreaSonLayout->count())
            return targetItem;

        DockItem *first = qobject_cast<DockItem *>(m_appAreaSonLayout->itemAt(0)->widget());
        DockItem *last = qobject_cast<DockItem *>(m_appAreaSonLayout->itemAt(m_appAreaSonLayout->count() - 1)->widget());

        if (m_position == Qt::TopEdge || m_position == Qt::BottomEdge) {

            if (point.x() < 0) {
                targetItem = first;
            } else {
                targetItem = last;
            }
        } else {

            if (point.y() < 0) {
                targetItem = first;
            } else {
                targetItem = last;
            }
        }

    }
    return targetItem;
}


void MainPanelControl::updateDisplayMode()
{
    DAnchorsBase::clearAnchors(m_appAreaSonWidget);
    DAnchors<QWidget> anchors(m_appAreaSonWidget);
    if (m_dislayMode == Dock::DisplayMode::Fashion) {
        anchors.setAnchor(Qt::AnchorHorizontalCenter, this, Qt::AnchorHorizontalCenter);
    } else {
        anchors.setAnchor(Qt::AnchorLeft, m_appAreaWidget, Qt::AnchorLeft);
    }
}

void MainPanelControl::onDisplayModeChanged()
{
    updateDisplayMode();
}
void MainPanelControl::updatePosition()
{
    onPositionChanged();
}

void MainPanelControl::onPositionChanged()
{
    updateMainPanelLayout();
}
