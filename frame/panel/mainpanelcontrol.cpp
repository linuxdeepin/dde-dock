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
#include "../item/pluginsitem.h"
#include "../item/traypluginitem.h"

#include <QDrag>
#include <QTimer>
#include <QStandardPaths>
#include <QString>
#include <QApplication>

#include <DGuiApplicationHelper>
#include <DWindowManagerHelper>

#define SPLITER_SIZE 2
#define TRASH_MARGIN 20
#define PLUGIN_MAX_SIZE  40
#define PLUGIN_MIN_SIZE  20
#define DESKTOP_SIZE  10

DWIDGET_USE_NAMESPACE

MainPanelControl::MainPanelControl(QWidget *parent)
    : QWidget(parent)
    , m_mainPanelLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_fixedAreaWidget(new QWidget(this))
    , m_appAreaWidget(new QWidget(this))
    , m_trayAreaWidget(new QWidget(this))
    , m_pluginAreaWidget(new QWidget(this))
    , m_desktopWidget(new QWidget(this))
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_trayAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_pluginLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_position(Position::Bottom)
    , m_placeholderItem(nullptr)
    , m_appDragWidget(nullptr)
    , m_dislayMode(Efficient)
    , m_fixedSpliter(new QLabel(this))
    , m_appSpliter(new QLabel(this))
    , m_traySpliter(new QLabel(this))
    , m_isHover(false)
{
    init();
    updateMainPanelLayout();
    setAcceptDrops(true);
    setMouseTracking(true);
    m_desktopWidget->setMouseTracking(true);

    m_appAreaWidget->installEventFilter(this);
    m_appAreaSonWidget->installEventFilter(this);
    m_trayAreaWidget->installEventFilter(this);
    m_desktopWidget->installEventFilter(this);
}

MainPanelControl::~MainPanelControl()
{
}

void MainPanelControl::init()
{
    // 主窗口
    m_mainPanelLayout->addWidget(m_fixedAreaWidget);
    m_mainPanelLayout->addWidget(m_fixedSpliter);
    m_mainPanelLayout->addWidget(m_appAreaWidget);
    m_mainPanelLayout->addWidget(m_appSpliter);
    m_mainPanelLayout->addWidget(m_trayAreaWidget);
    m_mainPanelLayout->addWidget(m_traySpliter);
    m_mainPanelLayout->addWidget(m_pluginAreaWidget);

    m_mainPanelLayout->setMargin(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->setSpacing(0);
    m_mainPanelLayout->setAlignment(m_fixedSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_appSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_traySpliter, Qt::AlignCenter);

    // 固定区域
    m_fixedAreaWidget->setLayout(m_fixedAreaLayout);
    m_fixedAreaLayout->setMargin(0);
    m_fixedAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_fixedAreaLayout->setSpacing(0);

    // 应用程序
    m_appAreaSonWidget->setLayout(m_appAreaSonLayout);
    m_appAreaSonLayout->setMargin(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);
    m_appAreaSonLayout->setSpacing(0);

    // 托盘
    m_trayAreaWidget->setLayout(m_trayAreaLayout);
    m_trayAreaLayout->setMargin(0);
    m_trayAreaLayout->setContentsMargins(0, 10, 0, 10);
    m_trayAreaLayout->setSpacing(0);

    // 插件
    m_pluginAreaWidget->setLayout(m_pluginLayout);
    m_pluginLayout->setMargin(0);
    m_pluginLayout->setContentsMargins(10, 10, 10, 10);
    m_pluginLayout->setSpacing(10);

    //桌面
    m_mainPanelLayout->addWidget(m_desktopWidget);

}

void MainPanelControl::setDisplayMode(DisplayMode mode)
{
    if (mode == m_dislayMode)
        return;
    m_dislayMode = mode;
    updateDisplayMode();
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
        m_trayAreaLayout->setContentsMargins(0, 10, 0, 10);
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
        m_trayAreaLayout->setContentsMargins(10, 0, 10, 0);
        break;
    }
}

void MainPanelControl::addFixedAreaItem(int index, QWidget *wdg)
{
    m_fixedAreaLayout->insertWidget(index, wdg);
    resizeDockIcon();
}

void MainPanelControl::addAppAreaItem(int index, QWidget *wdg)
{
    m_appAreaSonLayout->insertWidget(index, wdg);
    resizeDockIcon();
}

void MainPanelControl::addTrayAreaItem(int index, QWidget *wdg)
{
    m_tray = static_cast<TrayPluginItem *>(wdg);
    m_trayAreaLayout->insertWidget(index, wdg);
    resizeDockIcon();
}

void MainPanelControl::addPluginAreaItem(int index, QWidget *wdg)
{
    m_pluginLayout->insertWidget(index, wdg);
    resizeDockIcon();
    QTimer::singleShot(50, this, [ = ] {m_pluginAreaWidget->adjustSize();});
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
    if (m_position == Position::Right || m_position == Position::Left)
        m_desktopWidget->setFixedSize(width(), DESKTOP_SIZE);
    else
       m_desktopWidget->setFixedSize(DESKTOP_SIZE, height());

    if (DisplayMode::Fashion == m_dislayMode)
        m_desktopWidget->setFixedSize(0, 0);

    resizeDockIcon();
    return QWidget::resizeEvent(event);
}

void MainPanelControl::updateAppAreaSonWidgetSize()
{
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        m_appAreaSonWidget->setMaximumHeight(this->height());
        m_appAreaSonWidget->setMaximumWidth(m_appAreaWidget->width());
    } else {
        m_appAreaSonWidget->setMaximumWidth(this->width());
        m_appAreaSonWidget->setMaximumHeight(m_appAreaWidget->height());
    }

    m_appAreaSonWidget->adjustSize();

    moveAppSonWidget();
}

void MainPanelControl::setPositonValue(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    updateMainPanelLayout();
}

void MainPanelControl::insertItem(int index, DockItem *item)
{
    item->installEventFilter(this);

    switch (item->itemType()) {
    case DockItem::Launcher:
        addFixedAreaItem(0, item);
        break;
    case DockItem::FixedPlugin:
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
}

void MainPanelControl::removeItem(DockItem *item)
{
    switch (item->itemType()) {
    case DockItem::Launcher:
    case DockItem::FixedPlugin:
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
    else if (targetItem->itemType() == DockItem::FixedPlugin)
        idx = m_fixedAreaLayout->indexOf(targetItem);
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

    //如果当前从桌面拖拽的的app是trash，则不能放入app任务栏中
    QString str = "file://";
    //启动器
    QString str_t = "";

    str.append(QStandardPaths::locate(QStandardPaths::DesktopLocation, "dde-trash.desktop"));
    str_t.append(QStandardPaths::locate(QStandardPaths::ApplicationsLocation, "dde-trash.desktop"));

    if ((str == e->mimeData()->data(m_draggingMimeKey)) || (str_t == e->mimeData()->data(m_draggingMimeKey)))
        return;

    if (m_delegate && m_delegate->appIsOnDock(e->mimeData()->data(m_draggingMimeKey)))
        return;

    e->accept();
}

void MainPanelControl::dragLeaveEvent(QDragLeaveEvent *e)
{
    Q_UNUSED(e);
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

            if (m_position == Dock::Top || m_position == Dock::Bottom) {
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
    if (watched == m_appAreaSonWidget) {
        if (event->type() == QEvent::LayoutRequest) {
            m_appAreaSonWidget->adjustSize();
            resizeDockIcon();

        } else {
            moveAppSonWidget();
        }

        if (event->type() == QEvent::Resize) {
            moveAppSonWidget();
        }
    }

    if (watched == m_desktopWidget){
        if (event->type() == QEvent::Enter){
            m_isHover = true;
            update();
        } else if (event->type() == QEvent::Leave){
            m_isHover = false;
            update();
        }
    }

    if (watched == m_appAreaWidget) {
        if (event->type() == QEvent::Resize)
            updateAppAreaSonWidgetSize();

        if (event->type() == QEvent::Move)
            moveAppSonWidget();
    }

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

    if (item->itemType() != DockItem::App && item->itemType() != DockItem::Plugins && item->itemType() != DockItem::FixedPlugin)
        return false;

    const QPoint pos = mouseEvent->globalPos();
    const QPoint distance = pos - m_mousePressPos;
    if (distance.manhattanLength() < QApplication::startDragDistance())
        return false;

    startDrag(item);

    return QWidget::eventFilter(watched, event);
}

void MainPanelControl::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton) {
        m_mousePressPos = e->globalPos();

        QRect rect(m_desktopWidget->pos(),m_desktopWidget->size());
        if (rect.contains(e->pos()))
            QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
    }

    QWidget::mousePressEvent(e);
}

void MainPanelControl::startDrag(DockItem *item)
{
    const QPixmap pixmap = item->grab();

    item->setDraging(true);
    item->update();

    QDrag *drag = nullptr;
    if (item->itemType() == DockItem::App) {
        AppDrag *appDrag = new AppDrag(item);

        m_appDragWidget = appDrag->appDragWidget();

        connect(m_appDragWidget, &AppDragWidget::destroyed, this, [ = ] {
            m_appDragWidget = nullptr;
        });

        appDrag->appDragWidget()->setOriginPos((m_appAreaSonWidget->mapToGlobal(item->pos())));
        appDrag->appDragWidget()->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));

        if (DWindowManagerHelper::instance()->hasComposite()) {
            appDrag->setPixmap(pixmap);
            m_appDragWidget->show();

            static_cast<QGraphicsView *>(m_appDragWidget)->viewport()->installEventFilter(this);
        } else {
            const QPixmap &dragPix = qobject_cast<AppItem *>(item)->appIcon();

            appDrag->QDrag::setPixmap(dragPix);
            appDrag->setHotSpot(dragPix.rect().center() / dragPix.devicePixelRatioF());
        }

        drag = appDrag;
    } else {
        drag = new QDrag(item);
        drag->setPixmap(pixmap);
    }
    drag->setHotSpot(pixmap.rect().center() / pixmap.devicePixelRatioF());
    drag->setMimeData(new QMimeData);
    drag->exec(Qt::MoveAction);

    // app关闭特效情况下移除
    if (item->itemType() == DockItem::App && !DWindowManagerHelper::instance()->hasComposite()) {
        if (m_appDragWidget->isRemoveAble())
            qobject_cast<AppItem *>(item)->undock();
    }

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
        case DockItem::FixedPlugin:
            parentWidget = m_fixedAreaWidget;
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
        if (dockItem->itemType() == DockItem::Plugins) {
            if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
                rect.setSize(QSize(PLUGIN_MAX_SIZE, height()));
            } else {
                rect.setSize(QSize(width(), PLUGIN_MAX_SIZE));
            }
        } else {
            rect.setSize(dockItem->size());
        }
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

        if (m_position == Dock::Top || m_position == Dock::Bottom) {

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
    moveAppSonWidget();
}

void MainPanelControl::moveAppSonWidget()
{
    QRect rect(QPoint(0, 0), m_appAreaSonWidget->size());
    if (DisplayMode::Efficient == m_dislayMode) {
        switch (m_position) {
        case Top:
        case Bottom :
            rect.moveTo(m_appAreaWidget->pos());
            break;
        case Right:
        case Left:
            rect.moveTo(m_appAreaWidget->pos());
            break;
        }
    } else {
        switch (m_position) {
        case Top:
        case Bottom :
            rect.moveCenter(this->rect().center());
            if (rect.right() > m_appAreaWidget->geometry().right()) {
                rect.moveRight(m_appAreaWidget->geometry().right());
            }

            break;
        case Right:
        case Left:
            rect.moveCenter(this->rect().center());
            if (rect.bottom() > m_appAreaWidget->geometry().bottom()) {
                rect.moveBottom(m_appAreaWidget->geometry().bottom());
            }

            break;
        }
    }

    m_appAreaSonWidget->move(rect.x(), rect.y());
}

void MainPanelControl::itemUpdated(DockItem *item)
{
    item->parentWidget()->adjustSize();
}

void MainPanelControl::paintEvent(QPaintEvent *event)
{
    QPainter painter(this);
    QColor color;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        color = Qt::black;
        painter.setOpacity(0.5);
    } else {
        color = Qt::white;
        painter.setOpacity(0.1);
    }

    painter.fillRect(m_fixedSpliter->geometry(), color);
    painter.fillRect(m_appSpliter->geometry(), color);
    painter.fillRect(m_traySpliter->geometry(), color);

    //描绘桌面区域的颜色
    painter.setOpacity(1);
    QPen pen;
    QColor penColor(0, 0, 0, 25);
    pen.setWidth(2);
    pen.setColor(penColor);
    painter.setPen(pen);
    painter.drawRect(m_desktopWidget->geometry());
    if (m_isHover){
        painter.fillRect(m_desktopWidget->geometry(), QColor(255, 255, 255, 51));
        return;
    }
     painter.fillRect(m_desktopWidget->geometry(), QColor(255, 255, 255, 25));
}

void MainPanelControl::resizeDockIcon()
{
    if (!m_tray)
        return;
    // 插件有点特殊，因为会引入第三方的插件，并不会受dock的缩放影响，我们只能限制我们自己的插件，否则会导致显示错误。
    // 以下是受控制的插件
    PluginsItem *trashPlugin = nullptr;
    int pluginTotalWidth = 0;
    int pluginTotalHeight = 0;
    int count = 0; //计算除时间和垃圾箱之外的插件个数
    for (int i = 0; i < m_pluginLayout->count(); ++ i) {
        PluginsItem *w = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
        if (w->pluginName() == "trash") {
            trashPlugin = w;
        } else if (w->pluginName() == "datetime") {
        } else {
            pluginTotalWidth += w->width();
            pluginTotalHeight += w->height();
            count++;
        }
    }

    // 总宽度
    int totalLength = ((m_position == Position::Top) || (m_position == Position::Bottom)) ? width() : height();
    // 减去托盘间隔区域
    totalLength -= (m_tray->trayVisableItemCount() + 1) * 10;
    // 减去插件间隔
    totalLength -= (m_pluginLayout->count() + 1) * 10;
    // 减去3个分割线的宽度
    totalLength -= 3 * SPLITER_SIZE;

    // 减去所有插件宽度，加上参与计算的3个插件宽度
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        totalLength -= m_pluginAreaWidget->width();
        if (trashPlugin) totalLength += trashPlugin->width();
        totalLength += pluginTotalWidth;
        totalLength -= m_desktopWidget->width();
    } else {
        totalLength -= m_pluginAreaWidget->height();
        if (trashPlugin) totalLength += trashPlugin->height();
        totalLength += pluginTotalHeight;
        totalLength -= m_desktopWidget->height();
    }

    if (totalLength < 0)
        return;

    // 参与计算的插件的个数（包含托盘和插件，垃圾桶，关机，屏幕键盘）
    int pluginCount = m_tray->trayVisableItemCount() + (trashPlugin ? 1 : 0) + count;

    // icon个数
    int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count() + pluginCount;

    int iconSize = 0;

    // 余数
    int yu = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    iconSize = (totalLength - yu) / iconCount;

    if (iconSize < 20 || iconSize > 40) {

        // 减去插件和托盘的宽度
        if (iconSize < 20)
            totalLength -= 20 * pluginCount;
        else
            totalLength -= 40 * pluginCount;

        iconCount -= pluginCount;

        // 余数
        int yu = (totalLength % iconCount);
        // icon宽度 = (总宽度-余数)/icon个数
        iconSize = (totalLength - yu) / iconCount;
    }

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        if (iconSize >= height()) {
            calcuDockIconSize(height(), height());
        } else {
            calcuDockIconSize(iconSize, height());
        }
    } else {
        if (iconSize >= width()) {
            calcuDockIconSize(width(), width());
        } else {
            calcuDockIconSize(width(), iconSize);
        }
    }
}

void MainPanelControl::calcuDockIconSize(int w, int h)
{
    for (int i = 0; i < m_fixedAreaLayout->count(); ++ i) {
        m_fixedAreaLayout->itemAt(i)->widget()->setFixedSize(w, h);
    }

    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        m_fixedSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_appSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_traySpliter->setFixedSize(SPLITER_SIZE, int(w * 0.5));

        // 垃圾桶
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
            if (pluginItem->pluginName() == "trash"){
                  pluginItem->setFixedSize(std::min(w, h - 20), h - 20);
            }
        }

        for (int i = 0; i < m_appAreaSonLayout->count(); ++ i) {
            m_appAreaSonLayout->itemAt(i)->widget()->setMaximumWidth(h);
            m_appAreaSonLayout->itemAt(i)->widget()->setMaximumHeight(QWIDGETSIZE_MAX);
        }

    } else {
        m_fixedSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_appSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_traySpliter->setFixedSize(int(h * 0.5), SPLITER_SIZE);

        // 垃圾桶
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
            if (pluginItem->pluginName() == "trash"){
                  pluginItem->setFixedSize(w - 20, std::min(w - 20, h));
            }
        }

        for (int i = 0; i < m_appAreaSonLayout->count(); ++ i) {
            m_appAreaSonLayout->itemAt(i)->widget()->setMaximumHeight(w);
            m_appAreaSonLayout->itemAt(i)->widget()->setMaximumWidth(QWIDGETSIZE_MAX);
        }
    }

    // 插件和托盘

    // 托盘上每个图标大小
    int tray_item_size = 20;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        w = qBound(20, w, 40);
        tray_item_size = std::min(w, h - 20);
    } else {
        h = qBound(20, h, 40);
        tray_item_size = std::min(w - 20, h);
    }

    if (tray_item_size < 20)
        return;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        m_tray->centralWidget()->setProperty("iconSize", tray_item_size);

        // 插件
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
            if (pluginItem->pluginName() != "trash" && pluginItem->pluginName() != "datetime"){
                  pluginItem->setFixedSize(tray_item_size, h - 20);
            }
        }
    } else {
        m_tray->centralWidget()->setProperty("iconSize", tray_item_size);
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
            if (pluginItem->pluginName() != "trash" && pluginItem->pluginName() != "datetime"){
                pluginItem->setFixedSize(w - 20, tray_item_size);
            }
        }
    }
}

void MainPanelControl::getTrayVisableItemCount()
{
    if (m_trayAreaLayout->count() > 0) {
        TrayPluginItem *w = static_cast<TrayPluginItem *>(m_trayAreaLayout->itemAt(0)->widget());
        m_trayIconCount = w->trayVisableItemCount();
    } else {
        m_trayIconCount = 0;
    }

    resizeDockIcon();

    // 模式切换时，托盘区域宽度错误，对应任务11933
    m_trayAreaWidget->adjustSize();
}
