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
#include "../item/placeholderitem.h"
#include "../item/components/appdrag.h"
#include "../item/appitem.h"
#include "../item/pluginsitem.h"
#include "../item/traypluginitem.h"
#include "../controller/dockitemmanager.h"
#include "util/touchsignalmanager.h"

#include <QDrag>
#include <QTimer>
#include <QStandardPaths>
#include <QString>
#include <QApplication>
#include <QGSettings>
#include <QPointer>

#include <DGuiApplicationHelper>
#include <DWindowManagerHelper>

#define SPLITER_SIZE 2
#define TRASH_MARGIN 20
#define PLUGIN_MAX_SIZE  40
#define PLUGIN_MIN_SIZE  20
#define DESKTOP_SIZE  10

DWIDGET_USE_NAMESPACE

static QGSettings *GSettingsByLaunch()
{
    static QGSettings settings("com.deepin.dde.dock.module.launcher");
    return &settings;
}

MainPanelControl::MainPanelControl(QWidget *parent)
    : QWidget(parent)
    , m_mainPanelLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_fixedAreaWidget(new QWidget(this))
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_fixedSpliter(new QLabel(this))
    , m_appAreaWidget(new QWidget(this))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_appSpliter(new QLabel(this))
    , m_trayAreaWidget(new QWidget(this))
    , m_trayAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_traySpliter(new QLabel(this))
    , m_pluginAreaWidget(new QWidget(this))
    , m_pluginLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_desktopWidget(new DesktopWidget(this))
    , m_position(Position::Bottom)
    , m_placeholderItem(nullptr)
    , m_appDragWidget(nullptr)
    , m_dislayMode(Efficient)
    , m_isHover(false)
    , m_needRecoveryWin(false)
    , m_isEnableLaunch(true)
{
    initUi();
    updateMainPanelLayout();
    setAcceptDrops(true);
    setMouseTracking(true);
    m_desktopWidget->setMouseTracking(true);
    m_desktopWidget->setObjectName("showdesktoparea");

    m_appAreaWidget->installEventFilter(this);
    m_appAreaSonWidget->installEventFilter(this);
    m_trayAreaWidget->installEventFilter(this);
    m_desktopWidget->installEventFilter(this);

    //在设置每条线大小前，应该设置fixedsize(0,0)
    //应为paintEvent函数会先调用设置背景颜色，大小为随机值
    m_fixedSpliter->setFixedSize(0,0);
    m_appSpliter ->setFixedSize(0,0);
    m_traySpliter->setFixedSize(0,0);
}

MainPanelControl::~MainPanelControl()
{
}

void MainPanelControl::initUi()
{
    /* 固定区域 */
    m_fixedAreaWidget->setObjectName("fixedarea");
    m_fixedAreaWidget->setLayout(m_fixedAreaLayout);
    m_fixedAreaLayout->setSpacing(0);
    m_fixedAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_fixedAreaWidget);

    m_fixedSpliter->setObjectName("spliter_fix");
    m_mainPanelLayout->addWidget(m_fixedSpliter, Qt::AlignCenter);

    /* 应用程序区域 */
    m_appAreaWidget->setAccessibleName("AppFullArea");
    m_mainPanelLayout->addWidget(m_appAreaWidget);
    m_appAreaSonLayout->setSpacing(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);
    m_appAreaSonWidget->setObjectName("apparea");
    m_appAreaSonWidget->setLayout(m_appAreaSonLayout);
    m_appAreaSonLayout->setSpacing(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);

    m_appSpliter->setObjectName("spliter_app");
    m_mainPanelLayout->addWidget(m_appSpliter, Qt::AlignCenter);

    /* 托盘区域 */
    m_trayAreaWidget->setObjectName("trayarea");
    m_trayAreaWidget->setLayout(m_trayAreaLayout);
    m_trayAreaLayout->setSpacing(0);
    m_trayAreaLayout->setContentsMargins(0, 10, 0, 10);
    m_mainPanelLayout->addWidget(m_trayAreaWidget);

    m_traySpliter->setObjectName("spliter_tray");
    m_mainPanelLayout->addWidget(m_traySpliter, Qt::AlignCenter);

    /* 插件区域 */
    m_pluginAreaWidget->setObjectName("pluginarea");
    m_pluginAreaWidget->setLayout(m_pluginLayout);
    m_pluginLayout->setSpacing(10);
    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_pluginAreaWidget);

    /* 桌面预览 */
    m_mainPanelLayout->addWidget(m_desktopWidget);

    m_mainPanelLayout->setSpacing(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);

    connect(GSettingsByLaunch(), &QGSettings::changed, this, &MainPanelControl::onGSettingsChanged);
}

void MainPanelControl::onGSettingsChanged(const QString &key)
{
    if (key != "enable") {
        return;
    }

    QGSettings *setting = GSettingsByLaunch();

    if (setting->keys().contains("enable")) {
        const bool isEnable = GSettingsByLaunch()->keys().contains("enable") && GSettingsByLaunch()->get("enable").toBool();
        if (isEnable && setting->get("enable").toBool()) {
            m_fixedAreaWidget->setVisible(true);
            m_fixedSpliter->setVisible(true);
            m_isEnableLaunch = true;
        } else {
            m_fixedAreaWidget->setVisible(false);
            m_fixedSpliter->setVisible(false);
            m_isEnableLaunch = false;
        }
    }
}

void MainPanelControl::setDisplayMode(DisplayMode dislayMode)
{
    if (dislayMode == m_dislayMode)
        return;
    m_dislayMode = dislayMode;
    updateDisplayMode();
}

void MainPanelControl::updateMainPanelLayout()
{
    switch (m_position) {
    case Position::Top:
    case Position::Bottom:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Fixed);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_mainPanelLayout->setDirection(QBoxLayout::LeftToRight);
        m_fixedAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_pluginLayout->setDirection(QBoxLayout::LeftToRight);
        m_trayAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_appAreaSonLayout->setDirection(QBoxLayout::LeftToRight);
        m_trayAreaLayout->setContentsMargins(0, 10, 0, 10);
        m_pluginLayout->setContentsMargins(10, 0, 10, 0);
        break;
    case Position::Right:
    case Position::Left:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_pluginAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Fixed);
        m_trayAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_mainPanelLayout->setDirection(QBoxLayout::TopToBottom);
        m_fixedAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_pluginLayout->setDirection(QBoxLayout::TopToBottom);
        m_trayAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_appAreaSonLayout->setDirection(QBoxLayout::TopToBottom);
        m_trayAreaLayout->setContentsMargins(10, 0, 10, 0);
        m_pluginLayout->setContentsMargins(0, 10, 0, 10);
        break;
    }

    resizeDesktopWidget();
    resizeDockIcon();
}

void MainPanelControl::addFixedAreaItem(int index, QWidget *wdg)
{
    if(m_position == Position::Top || m_position == Position::Bottom){
        wdg->setMaximumSize(height(),height());
    } else {
        wdg->setMaximumSize(width(),width());
    }
    m_fixedAreaLayout->insertWidget(index, wdg);
}

void MainPanelControl::addAppAreaItem(int index, QWidget *wdg)
{
    if(m_position == Position::Top || m_position == Position::Bottom){
        wdg->setMaximumSize(height(),height());
    } else {
        wdg->setMaximumSize(width(),width());
    }
    m_appAreaSonLayout->insertWidget(index, wdg);
}

void MainPanelControl::addTrayAreaItem(int index, QWidget *wdg)
{
    m_tray = static_cast<TrayPluginItem *>(wdg);
    m_trayAreaLayout->insertWidget(index, wdg);
}

void MainPanelControl::addPluginAreaItem(int index, QWidget *wdg)
{
    //因为日期时间插件和其他插件的大小有异，为了方便设置边距，在插件区域布局再添加一层布局设置边距
    //因此在处理插件图标时，需要通过两层布局判断是否为需要的插件，例如拖动插件位置等判断
    QBoxLayout * boxLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    boxLayout->addWidget(wdg, 0, Qt::AlignCenter);
    m_pluginLayout->insertLayout(index, boxLayout, 0);
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
    //因为日期时间插件大小和其他插件有异，为了方便设置边距，各插件中增加一层布局
    //因此remove插件图标时，需要从多的一层布局中取widget进行判断是否需要移除的插件
    for (int i = 0; i < m_pluginLayout->count(); ++i) {
        QLayoutItem *layoutItem = m_pluginLayout->itemAt(i);
        QLayout *boxLayout = layoutItem->layout();
        if (boxLayout && boxLayout->itemAt(0)->widget() == wdg) {
            boxLayout->removeWidget(wdg);
            m_pluginLayout->removeItem(layoutItem);
        }
    }
}

void MainPanelControl::resizeEvent(QResizeEvent *event)
{
    resizeDesktopWidget();
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
    QTimer::singleShot(0, this, [=] {
        updateMainPanelLayout();
    });
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
    resizeDockIcon();
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
    resizeDockIcon();
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
    else if (targetItem->itemType() == DockItem::Plugins){
        //因为日期时间插件大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
        //因此有拖动图标时，需要从多的一层布局中判断是否相同插件而获取插件位置顺序
        for (int i = 0; i < m_pluginLayout->count(); ++i) {
            QLayout *layout = m_pluginLayout->itemAt(i)->layout();
            if (layout && layout->itemAt(0)->widget() == targetItem) {
                idx = i;
                break;
            }
        }
    } else if (targetItem->itemType() == DockItem::FixedPlugin)
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
    //拖拽图标到任务栏时，如果拖拽到垃圾箱插件图标widget上，则默认不允许拖拽，其他位置默认为允许拖拽
    QWidget *widget = QApplication::widgetAt(QCursor::pos());
    //"trash-centralwidget"名称是在PluginsItem类中m_centralWidget->setObjectName(pluginInter->pluginName() + "-centralwidget");
    if (widget && widget->objectName() == "trash-centralwidget") {
        return;
    }

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
    DockItem *sourceItem = qobject_cast<DockItem *>(e->source());
    if (sourceItem) {
        handleDragMove(e, false);
        return;
    }

    // 拖app到dock上
    const char *RequestDockKey = "RequestDock";
    const char *RequestDockKeyFallback = "text/plain";
    const char *DesktopMimeType = "application/x-desktop";
    auto DragmineData = e->mimeData();

    m_draggingMimeKey = DragmineData->formats().contains(RequestDockKey) ? RequestDockKey : RequestDockKeyFallback;

    // dragging item is NOT a desktop file
    if (QMimeDatabase().mimeTypeForFile(DragmineData->data(m_draggingMimeKey)).name() != DesktopMimeType) {
        m_draggingMimeKey.clear();
        e->setAccepted(false);
        qDebug() << "dragging item is NOT a desktop file";
        return;
    }

    //如果当前从桌面拖拽的的app是trash，则不能放入app任务栏中
    QString str = "file://";
    //启动器
    QString str_t = "";

    str.append(QStandardPaths::locate(QStandardPaths::DesktopLocation, "dde-trash.desktop"));
    str_t.append(QStandardPaths::locate(QStandardPaths::ApplicationsLocation, "dde-trash.desktop"));

    if ((str == DragmineData->data(m_draggingMimeKey)) || (str_t == DragmineData->data(m_draggingMimeKey))) {
        e->setAccepted(false);
        return;
    }

    if (m_delegate && m_delegate->appIsOnDock(DragmineData->data(m_draggingMimeKey))) {
        e->setAccepted(false);
        return;
    }

    handleDragMove(e, false);
}

bool MainPanelControl::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_appAreaSonWidget) {
        switch (event->type()) {
        case QEvent::LayoutRequest:
            m_appAreaSonWidget->adjustSize();
            resizeDockIcon();
            break;
        case QEvent::Resize:
            resizeDockIcon();
            break;
        default:
            moveAppSonWidget();
            break;
        }
    }

    if (watched == m_desktopWidget) {
        if (event->type() == QEvent::Enter) {
            if (checkNeedShowDesktop()) {
                m_needRecoveryWin = true;
                QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
            }
            m_isHover = true;
            update();
        } else if (event->type() == QEvent::Leave) {
            // 鼠标移入时隐藏了窗口，移出时恢复
            if (m_needRecoveryWin) {
                QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
            }
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

    QMouseEvent *mouseEvent = dynamic_cast<QMouseEvent *>(event);
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

    // source为MouseEventSynthesizedByQt时，事件由触屏事件转换而来，触屏没有收到后端的延迟触屏信号时不进行拖动
    if (mouseEvent->source() == Qt::MouseEventSynthesizedByQt && !TouchSignalManager::instance()->isDragIconPress()) {
        return false;
    }

    startDrag(item);

    return QWidget::eventFilter(watched, event);
}

void MainPanelControl::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton) {
        m_mousePressPos = e->globalPos();

        QRect rect(m_desktopWidget->pos(), m_desktopWidget->size());
        if (rect.contains(e->pos())) {
            if (m_needRecoveryWin) {
                // 手动点击 显示桌面窗口 后，鼠标移出时不再调用显/隐窗口进程，以手动点击设置为准
                m_needRecoveryWin = false;
            } else {
                // 需求调整，鼠标移入，预览桌面时再点击显示桌面保持显示桌面状态，再点击才切换桌面显、隐状态
                QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
            }
        }
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
        const QPixmap &dragPix = qobject_cast<AppItem *>(item)->appIcon();

        appDrag->setPixmap(dragPix);
        m_appDragWidget->show();

        if (DWindowManagerHelper::instance()->hasComposite()) {
            static_cast<QGraphicsView *>(m_appDragWidget)->viewport()->installEventFilter(this);
        } else {
            appDrag->QDrag::setPixmap(dragPix);
       }

        drag = appDrag;
        drag->setHotSpot(dragPix.rect().center() / dragPix.devicePixelRatioF());
    } else {
        drag = new QDrag(item);
        drag->setPixmap(pixmap);
        drag->setHotSpot(pixmap.rect().center() / pixmap.devicePixelRatioF());
    }
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

        DockItem *dockItem = nullptr;
        if (parentWidget == m_pluginAreaWidget) {
            QLayout *layout = layoutItem->layout();
            if (layout) {
                dockItem = qobject_cast<DockItem *>(layout->itemAt(0)->widget());
            }
        } else{
            dockItem = qobject_cast<DockItem *>(layoutItem->widget());
        }

        if (!dockItem)
            continue;

        QRect rect(dockItem->pos(), dockItem->size());
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
    resizeDesktopWidget();
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
            if (rect.left() < m_appAreaWidget->geometry().left()) {
                rect.moveLeft(m_appAreaWidget->geometry().left());
            }
            break;
        case Right:
        case Left:
            rect.moveCenter(this->rect().center());
            if (rect.bottom() > m_appAreaWidget->geometry().bottom()) {
                rect.moveBottom(m_appAreaWidget->geometry().bottom());
            }
            if (rect.top() < m_appAreaWidget->geometry().top()) {
                rect.moveTop(m_appAreaWidget->geometry().top());
            }
            break;
        }
    }

    m_appAreaSonWidget->move(rect.x(), rect.y());
}

void MainPanelControl::itemUpdated(DockItem *item)
{
    item->parentWidget()->adjustSize();
    resizeDockIcon();
}

void MainPanelControl::showEvent(QShowEvent *event)
{
    QTimer::singleShot(0, this, [ = ] {
        onGSettingsChanged("enable");
    });

    return QWidget::showEvent(event);
}

void MainPanelControl::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);

    QPainter painter(this);
    QColor color;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        color = Qt::black;
        painter.setOpacity(0.5);
    } else {
        color = Qt::white;
        painter.setOpacity(0.1);
    }

    if (m_isEnableLaunch)
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
    if (m_isHover) {
        painter.fillRect(m_desktopWidget->geometry(), QColor(255, 255, 255, 51));
        return;
    }
    painter.fillRect(m_desktopWidget->geometry(), QColor(255, 255, 255, 25));
}

void MainPanelControl::resizeDockIcon()
{
    // 插件有点特殊，因为会引入第三方的插件，并不会受dock的缩放影响，我们只能限制我们自己的插件，否则会导致显示错误。
    // 以下是受控制的插件
    PluginsItem *trashPlugin = nullptr;
    PluginsItem *shutdownPlugin = nullptr;
    PluginsItem *keyboardPlugin = nullptr;
    PluginsItem *notificationPlugin = nullptr;

    //因为日期时间大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
    //因此需要通过多一层布局来获取各插件
    for (int i = 0; i < m_pluginLayout->count(); ++ i) {
        QLayout *layout = m_pluginLayout->itemAt(i)->layout();
        if (layout) {
            PluginsItem *w = static_cast<PluginsItem *>(m_pluginLayout->itemAt(i)->widget());
            if (w) {
                if (w->pluginName() == "trash") {
                    trashPlugin = w;
                } else if (w->pluginName() == "shutdown") {
                    shutdownPlugin = w;
                } else if (w->pluginName() == "onboard") {
                    keyboardPlugin = w;
                } else if (w->pluginName() == "notifications") {
                    notificationPlugin = w;
                }
            }
        }
    }

    // 总宽度
    int totalLength = ((m_position == Position::Top) || (m_position == Position::Bottom)) ? width() : height();
    // 减去托盘间隔区域
    if (m_tray) {
        totalLength -= (m_tray->trayVisableItemCount() + 1) * 10;
    }
    // 减去3个分割线的宽度
    totalLength -= 3 * SPLITER_SIZE;

    // 减去所有插件宽度，加上参与计算的4个插件宽度
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        totalLength -= m_pluginAreaWidget->width();
        if (trashPlugin) totalLength += trashPlugin->width();
        if (shutdownPlugin) totalLength += shutdownPlugin->width();
        if (keyboardPlugin) totalLength += keyboardPlugin->width();
        if (notificationPlugin) totalLength += notificationPlugin->width();
        totalLength -= m_desktopWidget->width();
    } else {
        totalLength -= m_pluginAreaWidget->height();
        if (trashPlugin) totalLength += trashPlugin->height();
        if (shutdownPlugin) totalLength += shutdownPlugin->height();
        if (keyboardPlugin) totalLength += keyboardPlugin->height();
        if (notificationPlugin) totalLength += notificationPlugin->height();
        totalLength -= m_desktopWidget->height();
    }

    if (totalLength < 0)
        return;

    // 参与计算的插件的个数（包含托盘和插件，垃圾桶，关机，屏幕键盘）
    int pluginCount = 0;
    if (m_tray) {
        pluginCount = m_tray->trayVisableItemCount() + (trashPlugin ? 1 : 0) + (shutdownPlugin ? 1 : 0) + (keyboardPlugin ? 1 : 0) + (notificationPlugin ? 1 : 0);
    } else {
        pluginCount = (trashPlugin ? 1 : 0) + (shutdownPlugin ? 1 : 0) + (keyboardPlugin ? 1 : 0) + (notificationPlugin ? 1 : 0);
    }
    // icon个数
    int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count() + pluginCount;
    // 余数
    int yu = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    int iconSize = (totalLength - yu) / iconCount;
    //计算插件图标的最大或最小值
    int tray_item_size = qBound(20, iconSize, 40);
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        tray_item_size = qMin(tray_item_size,height());
        tray_item_size = std::min(tray_item_size, height() - 20);
    } else {
        tray_item_size = qMin(tray_item_size,width());
        tray_item_size = std::min(tray_item_size, width() - 20);
    }
    //减去插件图标的大小后重新计算固定图标和应用图标的平均大小
    totalLength -= tray_item_size * pluginCount;
    iconCount -= pluginCount;
    // 余数
    yu = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    iconSize = (totalLength - yu) / iconCount;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        if (iconSize >= height()) {
            calcuDockIconSize(height(), height(), trashPlugin, shutdownPlugin, keyboardPlugin, notificationPlugin);
        } else {
            calcuDockIconSize(iconSize, height(), trashPlugin, shutdownPlugin, keyboardPlugin, notificationPlugin);
        }
    } else {
        if (iconSize >= width()) {
            calcuDockIconSize(width(), width(), trashPlugin, shutdownPlugin, keyboardPlugin, notificationPlugin);
        } else {
            calcuDockIconSize(width(), iconSize, trashPlugin, shutdownPlugin, keyboardPlugin, notificationPlugin);
        }
    }
}

void MainPanelControl::calcuDockIconSize(int w, int h, PluginsItem *trashPlugin, PluginsItem *shutdownPlugin, PluginsItem *keyboardPlugin, PluginsItem *notificationPlugin)
{
    int appItemSize = qMin(w, h);

    for (int i = 0; i < m_fixedAreaLayout->count(); ++i) {
        m_fixedAreaLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);
    }

    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        m_fixedSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_appSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_traySpliter->setFixedSize(SPLITER_SIZE, int(w * 0.5));
    } else {
        m_fixedSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_appSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_traySpliter->setFixedSize(int(h * 0.5), SPLITER_SIZE);
    }

    for (int i = 0; i < m_appAreaSonLayout->count(); ++i) {
        m_appAreaSonLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);
    }

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

    if (m_tray) {
        m_tray->centralWidget()->setProperty("iconSize", tray_item_size);
    }

    if (shutdownPlugin)
        shutdownPlugin->setFixedSize(tray_item_size, tray_item_size);
    if (keyboardPlugin)
        keyboardPlugin->setFixedSize(tray_item_size, tray_item_size);
    if (notificationPlugin)
        notificationPlugin->setFixedSize(tray_item_size, tray_item_size);
    if (trashPlugin)
        trashPlugin->setFixedSize(tray_item_size, tray_item_size);

    //因为日期时间大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
    //因此需要通过多一层布局来获取各插件
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        // 三方插件
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            QLayout *layout = m_pluginLayout->itemAt(i)->layout();
            if (layout) {
                PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());
                if (pItem) {
                    if (pItem->pluginName() == "datetime") {
                        pItem->setFixedSize(pItem->sizeHint().width(), h);
                        //排除官方插件
                    }else if(pItem->pluginName() == "trash" || pItem->pluginName() == "shutdown" || pItem->pluginName() == "onboard" || pItem->pluginName() == "notifications"){
                        pItem->setFixedSize(tray_item_size, tray_item_size);
                    } else {
                        //对第三方插件size不做限制
                        pItem->setFixedSize(pItem->sizeHint().width(), h);
                    }
                }
            }
        }
    } else {
        // 三方插件
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            QLayout *layout =  m_pluginLayout->itemAt(i)->layout();
            if (layout) {
                PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());
                if (pItem && pItem != trashPlugin && pItem != shutdownPlugin && pItem != keyboardPlugin && pItem !=notificationPlugin) {
                    if (pItem->pluginName() == "datetime") {
                        pItem->setFixedSize(w, pItem->sizeHint().height());
                        //排除官方插件
                    }else if(pItem->pluginName() == "trash" || pItem->pluginName() == "shutdown" || pItem->pluginName() == "onboard" || pItem->pluginName() == "notifications"){
                        pItem->setFixedSize(tray_item_size, tray_item_size);
                    } else {
                        //对第三方插件size不做限制
                        pItem->setFixedSize(w, pItem->sizeHint().height());
                    }
                }
            }
        }
    }

    int appTopAndBottomMargin = 0;
    int appLeftAndRightMargin = 0;

    int trayTopAndBottomMargin = 0;
    int trayLeftAndRightMargin = 0;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        appTopAndBottomMargin = (m_fixedAreaWidget->height() - appItemSize) / 2;
        trayTopAndBottomMargin = (m_trayAreaWidget->height() - tray_item_size) / 2;
    } else {
        appLeftAndRightMargin = (m_fixedAreaWidget->width() - appItemSize) / 2;
        trayLeftAndRightMargin = (m_trayAreaWidget->width() - tray_item_size) / 2;
    }

    m_fixedAreaLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
    m_appAreaSonLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
    m_trayAreaLayout->setContentsMargins(trayLeftAndRightMargin, trayTopAndBottomMargin, trayLeftAndRightMargin, trayTopAndBottomMargin);

    //去除这个，因为我不知道到底想给那些官方插件设置边距，而且去除以后dock无影响
    //因为日期时间插件大小和其他插件大小有异，需要单独设置各插件的边距
    //而不对日期时间插件设置边距
//    for (int i = 0; i < m_pluginLayout->count(); ++ i) {
//        QLayout *layout = m_pluginLayout->itemAt(i)->layout();
//        if (layout) {
//            PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());

//            if (pItem && pItem->pluginName() != "datetime") {
//                layout->setContentsMargins(trayLeftAndRightMargin, trayTopAndBottomMargin, trayLeftAndRightMargin, trayTopAndBottomMargin);
//            }
//        }
//    }
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
}

void MainPanelControl::resizeDesktopWidget()
{
    if (m_position == Position::Right || m_position == Position::Left)
        m_desktopWidget->setFixedSize(width(), DESKTOP_SIZE);
    else
        m_desktopWidget->setFixedSize(DESKTOP_SIZE, height());

    if (DisplayMode::Fashion == m_dislayMode)
        m_desktopWidget->setFixedSize(0, 0);
}

/**
 * @brief MainPanelControl::checkNeedShowDesktop 根据窗管提供接口（当前是否显示的桌面），提供鼠标
 * 移入 显示桌面窗口 区域时，是否需要显示桌面判断依据
 * @return 窗管返回 当前是桌面 或 窗管接口查询失败 返回false，否则true
 */
bool MainPanelControl::checkNeedShowDesktop()
{
    QDBusInterface wmInter("com.deepin.wm", "/com/deepin/wm", "com.deepin.wm");
    QList<QVariant> argumentList;
    QDBusMessage reply = wmInter.callWithArgumentList(QDBus::Block, QStringLiteral("GetIsShowDesktop"), argumentList);
    if (reply.type() == QDBusMessage::ReplyMessage && reply.arguments().count() == 1) {
        return !reply.arguments().at(0).toBool();
    }

    qDebug() << "wm call GetIsShowDesktop fail, res:" << reply.type();
    return false;
}
