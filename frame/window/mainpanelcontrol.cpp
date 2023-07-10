// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainpanelcontrol.h"
#include "constants.h"
#include "dockitem.h"
#include "placeholderitem.h"
#include "components/appdrag.h"
#include "appitem.h"
#include "pluginsitem.h"
#include "traypluginitem.h"
#include "dockitemmanager.h"
#include "touchsignalmanager.h"
#include "utils.h"
#include "desktop_widget.h"
#include "imageutil.h"
#include "multiscreenworker.h"
#include "displaymanager.h"
#include "recentapphelper.h"
#include "toolapphelper.h"
#include "multiwindowhelper.h"
#include "mainwindow.h"
#include "appmultiitem.h"
#include "dockscreen.h"
#include "docksettings.h"
#include "docktraywindow.h"
#include "quicksettingcontroller.h"

#include <QDrag>
#include <QUrl>
#include <QTimer>
#include <QStandardPaths>
#include <QString>
#include <QApplication>
#include <QPointer>
#include <QBoxLayout>
#include <QLabel>
#include <QPixmap>
#include <QtConcurrent/QtConcurrentRun>
#include <QX11Info>

#include <qpa/qplatformnativeinterface.h>
#include <qpa/qplatformintegration.h>

#include <DGuiApplicationHelper>
#include <DWindowManagerHelper>

#include <X11/Xlib.h>

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
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_fixedSpliter(new QLabel(this))
    , m_appAreaWidget(new QWidget(this))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_appSpliter(new QLabel(this))
    , m_recentAreaWidget(new QWidget(this))
    , m_recentLayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_recentSpliter(new QLabel(this))
    , m_toolAreaWidget(new QWidget(this))
    , m_toolAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight, m_toolAreaWidget))
    , m_multiWindowWidget(new QWidget(m_toolAreaWidget))
    , m_multiWindowLayout(new QBoxLayout(QBoxLayout::LeftToRight, m_multiWindowWidget))
    , m_toolSonAreaWidget(new QWidget(m_toolAreaWidget))
    , m_toolSonLayout(new QBoxLayout(QBoxLayout::LeftToRight, m_toolSonAreaWidget))
    , m_position(Position::Bottom)
    , m_placeholderItem(nullptr)
    , m_appDragWidget(nullptr)
    , m_displayMode(Efficient)
    , m_tray(new DockTrayWindow(this))
    , m_recentHelper(new RecentAppHelper(m_appAreaSonWidget, m_recentAreaWidget, this))
    , m_toolHelper(new ToolAppHelper(m_toolSonAreaWidget, this))
    , m_multiHelper(new MultiWindowHelper(m_appAreaSonWidget, m_multiWindowWidget, this))
    , m_showRecent(DockSettings::instance()->showRecent())
{
    initUI();
    initConnection();
    updateMainPanelLayout();
    updateModeChange();
    setAcceptDrops(true);
    setMouseTracking(true);

    m_appAreaWidget->installEventFilter(this);
    m_appAreaSonWidget->installEventFilter(this);
    m_fixedAreaWidget->installEventFilter(this);
    m_tray->installEventFilter(this);

    // 在设置每条线大小前，应该设置fixedsize(0,0)
    // 应为paintEvent函数会先调用设置背景颜色，大小为随机值
    m_fixedSpliter->setFixedSize(0, 0);
    m_appSpliter ->setFixedSize(0, 0);
    m_recentSpliter->setFixedSize(0, 0);
}

void MainPanelControl::initUI()
{
    /* 固定区域 */
    m_fixedAreaWidget->setObjectName("fixedarea");
    m_fixedAreaWidget->setLayout(m_fixedAreaLayout);
    m_fixedAreaLayout->setSpacing(0);
    m_fixedAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_fixedAreaWidget);

    m_fixedSpliter->setObjectName("spliter_fix");
    m_mainPanelLayout->addWidget(m_fixedSpliter);

    /* 应用程序区域 */
    m_appAreaWidget->setAccessibleName("AppFullArea");
    m_mainPanelLayout->addWidget(m_appAreaWidget);
    m_appAreaSonWidget->setObjectName("apparea");
    m_appAreaSonWidget->setLayout(m_appAreaSonLayout);
    m_appAreaSonLayout->setSpacing(0);
    m_appAreaSonLayout->setContentsMargins(0, 0, 0, 0);

    m_appSpliter->setObjectName("spliter_app");
    m_mainPanelLayout->addWidget(m_appSpliter);

    /* 最近打开应用 */
    m_recentAreaWidget->setObjectName("recentarea");
    m_recentAreaWidget->setAccessibleName("recentarea");
    m_recentAreaWidget->setLayout(m_recentLayout);
    m_recentLayout->setSpacing(0);
    m_recentLayout->setContentsMargins(0, 0, 0, 0);
    m_recentLayout->setAlignment(Qt::AlignCenter);
    m_mainPanelLayout->addWidget(m_recentAreaWidget);

    m_recentSpliter->setObjectName("spliter_recent");
    m_mainPanelLayout->addWidget(m_recentSpliter);

    /* 工具应用 */
    // 包含窗口多开和工具组合
    m_toolAreaWidget->setObjectName("toolArea");
    m_toolAreaWidget->setAccessibleName("toolArea");
    m_toolAreaWidget->setLayout(m_toolAreaLayout);
    m_toolAreaLayout->setContentsMargins(0, 0, 0, 0);
    m_toolAreaLayout->setSpacing(0);
    m_mainPanelLayout->addWidget(m_toolAreaWidget);
    // 多开窗口区域
    m_multiWindowWidget->setObjectName("multiWindow");
    m_multiWindowWidget->setAccessibleName("multiWindow");
    m_multiWindowWidget->setLayout(m_multiWindowLayout);
    m_multiWindowLayout->setContentsMargins(0, 2, 0, 2);
    m_multiWindowLayout->setSpacing(0);
    m_toolAreaLayout->addWidget(m_multiWindowWidget);
    // 工具应用区域-包含打开窗口区域和回收站区域
    m_toolSonAreaWidget->setObjectName("toolsonarea");
    m_toolSonAreaWidget->setAccessibleName("toolsonarea");
    m_toolSonAreaWidget->setLayout(m_toolSonLayout);
    m_toolSonLayout->setSpacing(0);
    m_toolSonLayout->setContentsMargins(0, 0, 0, 0);
    m_toolSonLayout->setAlignment(Qt::AlignCenter);
    m_toolAreaLayout->addWidget(m_toolSonAreaWidget);

    // 添加托盘区域（包括托盘图标和插件）等
    m_tray->setObjectName("tray");
    m_mainPanelLayout->addWidget(m_tray);

    m_mainPanelLayout->setSpacing(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->setAlignment(m_fixedSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_appSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_recentSpliter, Qt::AlignCenter);
}

void MainPanelControl::initConnection()
{
    connect(m_recentHelper, &RecentAppHelper::requestUpdate, this, &MainPanelControl::requestUpdate);
    connect(m_recentHelper, &RecentAppHelper::recentVisibleChanged, this, &MainPanelControl::onRecentVisibleChanged);
    connect(m_recentHelper, &RecentAppHelper::dockAppVisibleChanged, this, &MainPanelControl::onDockAppVisibleChanged);
    connect(m_toolHelper, &ToolAppHelper::requestUpdate, this, &MainPanelControl::requestUpdate);
    connect(m_toolHelper, &ToolAppHelper::toolVisibleChanged, this, &MainPanelControl::onToolVisibleChanged);
    connect(m_multiHelper, &MultiWindowHelper::requestUpdate, this, &MainPanelControl::requestUpdate);
    connect(m_tray, &DockTrayWindow::requestUpdate, this, &MainPanelControl::onTrayRequestUpdate);
    connect(DockSettings::instance(), &DockSettings::showRecentChanged, this, [=] (bool show) {
        m_showRecent = show;
    });
}

/**
 * @brief MainPanelControl::setDisplayMode 根据任务栏显示模式更新界面显示，如果是时尚模式，没有‘显示桌面'区域，否则就有
 * @param dislayMode 任务栏显示模式
 */
void MainPanelControl::setDisplayMode(DisplayMode dislayMode)
{
    m_displayMode = dislayMode;
    m_recentHelper->setDisplayMode(dislayMode);
    m_tray->setDisplayMode(dislayMode);
    m_toolHelper->setDisplayMode(dislayMode);
    m_multiHelper->setDisplayMode(dislayMode);
    updateDisplayMode();
}

/**根据任务栏在屏幕上的位置，更新任务栏各控件布局
 * @brief MainPanelControl::updateMainPanelLayout
 */
void MainPanelControl::updateMainPanelLayout()
{
    switch (m_position) {
    case Position::Top:
    case Position::Bottom:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_mainPanelLayout->setDirection(QBoxLayout::LeftToRight);
        m_fixedAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_appAreaSonLayout->setDirection(QBoxLayout::LeftToRight);
        m_recentLayout->setDirection(QBoxLayout::LeftToRight);
        m_multiWindowLayout->setDirection(QBoxLayout::LeftToRight);
        m_toolAreaLayout->setDirection(QBoxLayout::LeftToRight);
        m_toolSonLayout->setDirection(QBoxLayout::LeftToRight);
        m_multiWindowLayout->setContentsMargins(0, 2, 0, 2);
        break;
    case Position::Right:
    case Position::Left:
        m_fixedAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        m_appAreaWidget->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_mainPanelLayout->setDirection(QBoxLayout::TopToBottom);
        m_fixedAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_appAreaSonLayout->setDirection(QBoxLayout::TopToBottom);
        m_recentLayout->setDirection(QBoxLayout::TopToBottom);
        m_multiWindowLayout->setDirection(QBoxLayout::TopToBottom);
        m_toolAreaLayout->setDirection(QBoxLayout::TopToBottom);
        m_toolSonLayout->setDirection(QBoxLayout::TopToBottom);
        m_multiWindowLayout->setContentsMargins(2, 0, 2, 0);
        break;
    }

    // 设置任务栏各区域图标大小
    resizeDockIcon();

    // 调整托盘区域大小
    onTrayRequestUpdate();
}

/**往固定区域添加应用
 * @brief MainPanelControl::addFixedAreaItem
 * @param index　位置索引，如果为负数则插入到最后，为正则插入到指定位置
 * @param wdg　应用指针对象
 */
void MainPanelControl::addFixedAreaItem(int index, QWidget *wdg)
{
    if (m_position == Position::Top || m_position == Position::Bottom) {
        wdg->setMaximumSize(height(),height());
    } else {
        wdg->setMaximumSize(width(),width());
    }
    m_fixedAreaLayout->insertWidget(index, wdg);
    Q_EMIT requestUpdate();
}

/**移除固定区域某一应用
 * @brief MainPanelControl::removeFixedAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removeFixedAreaItem(QWidget *wdg)
{
    m_fixedAreaLayout->removeWidget(wdg);
    Q_EMIT requestUpdate();
}

/**移除应用区域某一应用
 * @brief MainPanelControl::removeAppAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removeAppAreaItem(QWidget *wdg)
{
    m_appAreaSonLayout->removeWidget(wdg);
    Q_EMIT requestUpdate();
}

void MainPanelControl::resizeEvent(QResizeEvent *event)
{
    // 先通过消息循环让各部件调整好size后再计算图标大小
    // 避免因为部件size没有调整完导致计算的图标大小不准确
    // 然后重复触发m_pluginAreaWidget的reszie事件并重复计算，造成任务栏图标抖动问题
    QWidget::resizeEvent(event);
    resizeDockIcon();
}

/** 当用户从最近使用区域拖动应用到左侧应用区域的时候，将该应用驻留
 * @brief MainPanelControl::dockRecentApp
 * @param dockItem
 */
void MainPanelControl::dockRecentApp(DockItem *dockItem)
{
    // 如果不是插入或者当前不是特效模式，则无需做驻留操作
    if (m_dragIndex == -1 || m_displayMode != Dock::DisplayMode::Fashion)
        return;

    AppItem *appItem = qobject_cast<AppItem *>(dockItem);
    if (!appItem)
        return;

    // 如果控制中心设置不开启最近应用，则不让其驻留
    if (!m_showRecent)
        return;

    // 如果控制中心开启了最近应用并且当前应用是未驻留应用，则可以驻留
    if (!appItem->isDocked())
        appItem->requestDock();
}

PluginsItem *MainPanelControl::trash() const
{
    QuickSettingController *quickController = QuickSettingController::instance();
    QList<PluginsItemInterface *> toolPlugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tool);
    for (PluginsItemInterface *plugin : toolPlugins) {
        if (plugin->pluginName() != "trash")
            continue;

        return quickController->pluginItemWidget(plugin);
    }

    return nullptr;
}

/**根据任务栏所在位置， 设置应用区域控件的大小
 * @brief MainPanelControl::updateAppAreaSonWidgetSize
 */
void MainPanelControl::updateAppAreaSonWidgetSize()
{
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        m_appAreaSonWidget->setMaximumHeight(height());
        m_appAreaSonWidget->setMaximumWidth(m_appAreaWidget->width());
    } else {
        m_appAreaSonWidget->setMaximumWidth(width());
        m_appAreaSonWidget->setMaximumHeight(m_appAreaWidget->height());
    }

    m_appAreaSonWidget->adjustSize();
    moveAppSonWidget();
}

/**
 * @brief setPositonValue 根据传入的位置更新界面布局，比如任务栏在左，布局应该是上下布局，任务栏在下，应该是左右布局
 * @param position 任务栏的位置
 */
void MainPanelControl::setPositonValue(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    m_tray->setPositon(position);
    m_toolHelper->setPosition(position);

    QMetaObject::invokeMethod(this, &MainPanelControl::updateMainPanelLayout, Qt::QueuedConnection);
}

/**向任务栏插入各类应用,并将属于同一个应用的窗口合并到同一个应用图标
 * @brief MainPanelControl::insertItem
 * @param index 位置索引
 * @param item 应用指针对象
 */
void MainPanelControl::insertItem(int index, DockItem *item)
{
    if (!item)
        return;

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
        m_recentHelper->addAppItem(index, item);
        break;
    case DockItem::AppMultiWindow:
        m_multiHelper->addMultiWindow(index, static_cast<AppMultiItem *>(item));
        break;
    default:
        break;
    }

    // 同removeItem处 注意:不能屏蔽此接口，否则会造成插件插入时无法显示
    if (item->itemType() != DockItem::App)
        resizeDockIcon();

    item->checkEntry();
}

/**从任务栏移除某一应用，并更新任务栏图标大小
 * @brief MainPanelControl::removeItem
 * @param item 应用指针对象
 */
void MainPanelControl::removeItem(DockItem *item)
{
    switch (item->itemType()) {
    case DockItem::Launcher:
    case DockItem::FixedPlugin:
        removeFixedAreaItem(item);
        break;
    case DockItem::App:
    case DockItem::Placeholder:
        m_recentHelper->removeAppItem(item);
        break;
    case DockItem::AppMultiWindow:
        m_multiHelper->removeMultiWindow(static_cast<AppMultiItem *>(item));
        break;
    default:
        break;
    }

    item->removeEventFilter(this);

    /** 此处重新计算大小的时候icon的个数在原有个数上减少了一个，导致每个icon的大小跟原来大小不一致，需要重新设置setFixedSize
     *  在龙芯处理器上当app数量过多时，会导致拖动app耗时严重，造成卡顿
     *  注意:不能屏蔽此接口，否则会造成插件移除时无法更新icon大小
     */
    if (item->itemType() != DockItem::App)
        resizeDockIcon();
}

/**任务栏移动应用图标
 * @brief MainPanelControl::moveItem
 * @param sourceItem 即将插入的应用
 * @param targetItem 被移动的应用
 */
void MainPanelControl::moveItem(DockItem *sourceItem, DockItem *targetItem)
{
    // get target index
    int idx = -1;
    if (targetItem->itemType() == DockItem::App)
        idx = m_appAreaSonLayout->indexOf(targetItem);
    else if (targetItem->itemType() == DockItem::FixedPlugin)
        idx = m_fixedAreaLayout->indexOf(targetItem);
    else
        return;

    // remove old item
    removeItem(sourceItem);

    // insert new position
    if (sourceItem->isDragging()) {
        m_dragIndex = idx;
    }
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
        m_placeholderItem = nullptr;
    }
}

void MainPanelControl::dropEvent(QDropEvent *e)
{
    if (m_placeholderItem) {

        QUrl desktopPath = QUrl::fromUserInput(e->mimeData()->data(m_draggingMimeKey));

        emit itemAdded(desktopPath.toLocalFile(), m_appAreaSonLayout->indexOf(m_placeholderItem));

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

    if (appIsOnDock(DragmineData->data(m_draggingMimeKey))) {
        e->setAccepted(false);
        return;
    }

    handleDragMove(e, false);
}

bool MainPanelControl::eventFilter(QObject *watched, QEvent *event)
{
    // 在从时尚模式切换到高效模式的时候，
    // m_tray子部件会调整高度，此时会触发m_tray调整尺寸
    // 但是子部件的模式变化函数在FashionTrayItem部件中的
    // NormalContainer部件尺寸变化完成之前就已经结束，导致
    // NormalContainer没有更新自己的尺寸，引起插件区域拥挤
    //if (m_tray && watched == m_tray && event->type() == QEvent::Resize)
        //m_tray->pluginItem()->displayModeChanged(m_displayMode);

    // 更新应用区域大小和任务栏图标大小
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

    // fix:88133 在计算icon大小时m_pluginAreaWidget的数据错误
    if (watched == m_tray || watched == m_fixedAreaWidget) {
        switch (event->type()) {
        case QEvent::Resize:
            resizeDockIcon();
            break;
        default:
            break;
        }
    }

    // 更新应用区域子控件大小以及位置
    if (watched == m_appAreaWidget) {
        if (event->type() == QEvent::Resize)
            updateAppAreaSonWidgetSize();

        if (event->type() == QEvent::Move)
            moveAppSonWidget();
    }

    if (m_appDragWidget && watched == static_cast<QGraphicsView *>(m_appDragWidget)->viewport()) {
        bool isContains = rect().contains(mapFromGlobal(QCursor::pos()));
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

    static const QGSettings *g_settings = Utils::ModuleSettingsPtr("app");
    if (!g_settings || !g_settings->keys().contains("removeable") || g_settings->get("removeable").toBool())
        startDrag(item);

    return QWidget::eventFilter(watched, event);
}

void MainPanelControl::enterEvent(QEvent *event)
{
    if (Utils::IS_WAYLAND_DISPLAY) {
        Utils::updateCursor(this);
    }

    QWidget::enterEvent(event);
}

void MainPanelControl::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton) {
        m_mousePressPos = e->globalPos();
    }

    QWidget::mousePressEvent(e);
}

void MainPanelControl::startDrag(DockItem *dockItem)
{
    // 每次拖动使m_dragIndex==-1, 表明当前item的位置未发生变化
    m_dragIndex = -1;
    QPointer<DockItem> item = dockItem;
    const QPixmap pixmap = item->grab();

    item->setDraging(true);
    item->update();

    QDrag *drag = nullptr;
    if (item->itemType() == DockItem::App) {
        AppDrag *appDrag = new AppDrag(item);

        m_appDragWidget = appDrag->appDragWidget();

        connect(m_appDragWidget, &AppDragWidget::requestChangedArea, this, [ = ](QRect rect) {
            // 在区域改变的时候，出现分屏提示效果
            AppItem *appItem = static_cast<AppItem *>(dockItem);
            if (appItem->supportSplitWindow())
                appItem->startSplit(rect);
        });

        connect(m_appDragWidget, &AppDragWidget::requestSplitWindow, this, [ = ](ScreenSpliter::SplitDirection dir) {
            AppItem *appItem = static_cast<AppItem *>(dockItem);
            if (appItem->supportSplitWindow())
                appItem->splitWindowOnScreen(dir);
        });

        connect(m_appDragWidget, &AppDragWidget::destroyed, this, [ = ] {
            m_appDragWidget = nullptr;

            if (!item.isNull() && qobject_cast<AppItem *>(item)->isValid()) {
                // 如果是从最近打开区域移动到应用区域的，则需要将其固定
                dockRecentApp(item);
                if (-1 == m_appAreaSonLayout->indexOf(item) && m_dragIndex != -1) {
                    insertItem(m_dragIndex, item);
                    m_dragIndex = -1;
                }
                item->setDraging(false);
                item->update();
                // 发送拖拽完成事件
                m_recentHelper->resetAppInfo();
            }
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

    // isNeedBack 保存是否需要重置垃圾箱的AcceptDrops
    // 设置垃圾箱插件AcceptDrops false
    bool isNeedBack = false;
    PluginsItem *trashItem = trash();
    if (item->itemType() == DockItem::Plugins && trashItem && dockItem != trashItem) {
        trashItem->centralWidget()->setAcceptDrops(false);
        isNeedBack = true;
    }

    drag->setMimeData(new QMimeData);
    drag->exec(Qt::MoveAction);

    if (item->itemType() == DockItem::App && m_appDragWidget) {
        // TODO AppDragWidget中偶尔会出现拖拽结束后没有触发dropEvent的情况，因此exec结束后处理dropEvent中未执行的操作(临时处理方式)
        m_appDragWidget->execFinished();
    }

    if (item->itemType() == DockItem::App) {
        // 判断是否在回收站区域, 如果在回收站区域，则移除驻留
        if (!trashItem)
            return;

        QRect trashRect = trashItem->centralWidget()->geometry();
        QPoint pointMouse = trashItem->centralWidget()->mapFromGlobal(QCursor::pos());
        if (trashRect.contains(pointMouse)) {
            AppItem *appItem = qobject_cast<AppItem *>(dockItem);
            if (!appItem)
                return;

            // 先让其设置m_dragIndex==-1，避免在后续放到任务栏
            m_dragIndex = -1;
            appItem->setDraging(false);
            appItem->undock();
        }
    } else {
        m_appDragWidget = nullptr;
        item->setDraging(false);
        item->update();

        // isNeedBack是否需要设置垃圾箱插件AcceptDrops true
        if (isNeedBack)
            trashItem->centralWidget()->setAcceptDrops(true);
    }
}

DockItem *MainPanelControl::dropTargetItem(DockItem *sourceItem, QPoint point)
{
    QWidget *parentWidget = m_appAreaSonWidget;

    if (sourceItem) {
        switch (sourceItem->itemType()) {
        case DockItem::App:
            parentWidget = m_appAreaSonWidget;
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

        QRect rect(dockItem->pos(), dockItem->size());
        if (rect.contains(point)) {
            targetItem = dockItem;
            break;
        }
    }

    if (!targetItem && parentWidget == m_appAreaSonWidget) {
        // appitem调整顺序是，判断是否拖放在两边空白区域
        targetItem = sourceItem;
    }

    return targetItem;
}

void MainPanelControl::updateDisplayMode()
{
    updateModeChange();
    moveAppSonWidget();
}

void MainPanelControl::updateModeChange()
{
    m_toolAreaWidget->setVisible(m_displayMode == DisplayMode::Fashion);
    onRecentVisibleChanged(m_recentHelper->recentIsVisible());
    onDockAppVisibleChanged(m_recentHelper->dockAppIsVisible());
    onToolVisibleChanged(m_toolHelper->toolIsVisible());
    m_tray->setVisible(m_displayMode == DisplayMode::Efficient);
}

/**把驻留应用和被打开的应用所在窗口移动到指定位置
 * @brief MainPanelControl::moveAppSonWidget
 */
void MainPanelControl::moveAppSonWidget()
{
    QRect rect(QPoint(0, 0), m_appAreaSonWidget->size());
    if (DisplayMode::Efficient == m_displayMode) {
        rect.moveTo(m_appAreaWidget->pos());
    } else {
        switch (m_position) {
        case Top:
        case Bottom:
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

QSize MainPanelControl::suitableSize(const Position &position, int screenSize, double deviceRatio) const
{
    if (screenSize <= 0)
        return QSize(-1, -1);

    double ratio = deviceRatio;
    if (ratio <= 0)
        ratio = qApp->devicePixelRatio();

    int dockSize = (m_position == Position::Top || m_position == Position::Bottom) ? height() : width();
    if (dockSize == 0) {
        int windowSize = m_displayMode == DisplayMode::Efficient ?
            DockSettings::instance()->getWindowSizeEfficient() :
            DockSettings::instance()->getWindowSizeFashion();
        dockSize = windowSize;
    }

    if (m_displayMode == DisplayMode::Efficient) {
        // 如果是高效模式
        if (position == Position::Top || position == Position::Bottom)
            return QSize(static_cast<int>(screenSize / ratio), dockSize);

        return QSize(dockSize, static_cast<int>(screenSize / ratio));
    }

    // 如果是特效模式
    int totalLength = getScreenSize();
    // 减去插件区域的尺寸
    totalLength -= trayAreaSize(ratio);

    if (m_fixedSpliter->isVisible())
        totalLength -= SPLITER_SIZE;
    if (m_appSpliter->isVisible())
        totalLength -= SPLITER_SIZE;
    if (m_recentSpliter->isVisible())
        totalLength -= SPLITER_SIZE;

    // 需要参与计算的图标的总数
    int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count() + m_recentLayout->count() + m_toolSonLayout->count();
    int multiWindowCount = m_multiWindowLayout->count();
    if (iconCount <= 0 && multiWindowCount <= 0) {
        if (position == Position::Top || position == Position::Bottom)
            return QSize((static_cast<int>(dockSize)), dockSize);

        return QSize(dockSize, static_cast<int>(dockSize));
    }

    int redundantLength = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    int iconSize = qMin((static_cast<int>((totalLength - redundantLength) / iconCount / ratio)), dockSize);

    if (position == Position::Top || position == Position::Bottom) {
        int spliterWidth = m_fixedSpliter->isVisible() ? SPLITER_SIZE : 0;
        if (m_appSpliter->isVisible())
            spliterWidth += SPLITER_SIZE;

        if (m_recentSpliter->isVisible())
            spliterWidth += SPLITER_SIZE;

        int multiSize = 0;
        // 计算每个多开窗口的尺寸
        if (multiWindowCount > 0) {
            for (int i = 0; i < multiWindowCount; i++) {
                AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(m_multiWindowLayout->itemAt(i)->widget());
                if (!multiItem)
                    continue;

                multiSize += multiItem->width();
            }
        }

        int panelWidth = qMin(iconSize * iconCount + multiSize + static_cast<int>(spliterWidth),
                              static_cast<int>((screenSize - DOCKSPACE)));

        return QSize(panelWidth, static_cast<int>(dockSize));
    }

    int spliterHeight = m_fixedSpliter->isVisible() ? SPLITER_SIZE : 0;
    if (m_appSpliter->isVisible())
        spliterHeight += SPLITER_SIZE;

    if (m_recentSpliter->isVisible())
        spliterHeight += SPLITER_SIZE;

    int multiSize = 0;
    // 计算每个多开窗口的尺寸
    if (multiWindowCount > 0) {
        for (int i = 0; i < multiWindowCount; i++) {
            AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(m_multiWindowLayout->itemAt(i)->widget());
            if (!multiItem)
                continue;

            multiSize += multiItem->height();
        }
    }

    int panelHeight = qMin(iconSize * iconCount + multiSize + static_cast<int>(spliterHeight),
                           static_cast<int>((screenSize - DOCKSPACE)));

    return QSize(dockSize, panelHeight);
}

void MainPanelControl::itemUpdated(DockItem *item)
{
    item->updateGeometry();
    resizeDockIcon();
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

    if (m_fixedSpliter->isVisible())
        painter.fillRect(m_fixedSpliter->geometry(), color);

    if (m_appSpliter->isVisible())
        painter.fillRect(m_appSpliter->geometry(), color);

    if (m_recentSpliter->isVisible())
        painter.fillRect(m_recentSpliter->geometry(), color);
}

// 获取当前屏幕的高或者宽(任务栏上下的时候获取宽，左右获取高)
int MainPanelControl::getScreenSize() const
{
    DisplayManager *displayManager = DisplayManager::instance();
    QScreen *currentScreen = displayManager->screen(displayManager->primary());
    QScreen *screen = displayManager->screen(DockScreen::instance()->current());
    if (screen)
        currentScreen = screen;

    QRect screenRect = currentScreen->handle()->geometry();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return screenRect.width();

    return screenRect.height();
}

int MainPanelControl::trayAreaSize(qreal ratio) const
{
    if (m_displayMode == Dock::DisplayMode::Efficient)
        return (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom ? m_tray->width() * ratio: m_tray->height() * ratio);

    int length = 0;
    QWidgetList topLevelWidgets = qApp->topLevelWidgets();
    for (QWidget *widget : topLevelWidgets) {
        MainWindowBase *topWindow = qobject_cast<MainWindowBase *>(widget);
        if (!topWindow)
            continue;

        if (topWindow->windowType() != MainWindowBase::DockWindowType::MainWindow) {
            length += (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom ? topWindow->width() * ratio : topWindow->height() * ratio);
        }

        length += topWindow->dockSpace() * ratio;
    }

    return length;
}

/**重新计算任务栏上应用图标、插件图标的大小，并设置
 * @brief MainPanelControl::resizeDockIcon
 */
void MainPanelControl::resizeDockIcon()
{
    int iconSize = 0;
    // 总宽度
    if (m_displayMode == DisplayMode::Fashion) {
        // 时尚模式
        int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count();
        if (m_recentAreaWidget->isVisible())
            iconCount += m_recentLayout->count();

        if (m_toolAreaWidget->isVisible())
            iconCount += m_toolSonLayout->count();

        if (iconCount <= 0)
            return;

        int totalLength = getScreenSize() - trayAreaSize(qApp->devicePixelRatio());

        if (m_fixedSpliter->isVisible())
            totalLength -= SPLITER_SIZE;
        if (m_appSpliter->isVisible())
            totalLength -= SPLITER_SIZE;
        if (m_recentSpliter->isVisible())
            totalLength -= SPLITER_SIZE;

        // 余数
        int yu = (totalLength % iconCount);
        // icon宽度 = (总宽度-余数)/icon个数
        iconSize = (totalLength - yu) / iconCount;
    } else {
        int totalLength = getScreenSize() - trayAreaSize(qApp->devicePixelRatio());
        // 减去3个分割线的宽度
        if (m_fixedSpliter->isVisible())
            totalLength -= SPLITER_SIZE;
        if (m_appSpliter->isVisible())
            totalLength -= SPLITER_SIZE;
        if (m_recentSpliter->isVisible())
            totalLength -= SPLITER_SIZE;

        // 减去插件间隔大小, 只有一个插件或没有插件都是间隔20,2个或以上每多一个插件多间隔10
        totalLength -= 20;

        if (totalLength < 0)
            return;

        // 参与计算的插件的个数包含托盘和插件
        // 需要计算的图标总数
        int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count()/* + pluginCount*/;
        if (iconCount <= 0)
            return;

        // 余数
        int yu = (totalLength % iconCount);
        // icon宽度 = (总宽度-余数)/icon个数
        iconSize = (totalLength - yu) / iconCount;

        // 余数
        yu = (totalLength % iconCount);
        // icon宽度 = (总宽度-余数)/icon个数
        iconSize = (totalLength - yu) / iconCount;
    }

    iconSize = iconSize / qApp->devicePixelRatio();
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
    int appItemSize = qMin(w, h);
    for (int i = 0; i < m_fixedAreaLayout->count(); ++i)
        m_fixedAreaLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);

    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        m_fixedSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_appSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
        m_recentSpliter->setFixedSize(SPLITER_SIZE, int(w * 0.6));
    } else {
        m_fixedSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_appSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
        m_recentSpliter->setFixedSize(int(h * 0.6), SPLITER_SIZE);
    }

    // 时尚模式下判断是否需要显示最近打开的应用区域
    if (m_displayMode == Dock::DisplayMode::Fashion) {
        for (int i = 0; i < m_appAreaSonLayout->count(); ++i)
            m_appAreaSonLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);

        if (m_recentLayout->count() > 0) {
            for (int i = 0; i < m_recentLayout->count(); ++i)
                m_recentLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);

            // 时尚模式下计算最近打开应用区域的尺寸
            if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
                m_recentAreaWidget->setFixedSize(appItemSize * m_recentLayout->count(), QWIDGETSIZE_MAX);
            else
                m_recentAreaWidget->setFixedSize(QWIDGETSIZE_MAX, appItemSize * m_recentLayout->count());
        }

        if (m_multiWindowLayout->count() > 0) {
            QList<QSize> multiSizes;
            for (int i = 0; i < m_multiWindowLayout->count(); i++) {
                // 因为多开窗口的长宽会不一样，因此，需要将当前的尺寸传入
                // 由它自己来计算自己的长宽尺寸
                AppMultiItem *appMultiItem = qobject_cast<AppMultiItem *>(m_multiWindowLayout->itemAt(i)->widget());
                if (!appMultiItem)
                    continue;

                QSize size = appMultiItem->suitableSize(appItemSize);
                appMultiItem->setFixedSize(size);
                multiSizes << size;
            }
            // 计算多开窗口的尺寸
            int totalSize = 0;
            if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
                for (QSize size : multiSizes)
                    totalSize += size.width();

                m_multiWindowWidget->setFixedSize(totalSize, appItemSize);
            } else {
                for (QSize size : multiSizes)
                    totalSize += size.height();

                m_multiWindowWidget->setFixedSize(appItemSize, totalSize);
            }
        } else {
            m_multiWindowWidget->setFixedSize(0, 0);
        }
        if (m_toolSonLayout->count() > 0) {
            for (int i = 0; i < m_toolSonLayout->count(); i++)
                m_toolSonLayout->itemAt(i)->widget()->setFixedSize(appItemSize, appItemSize);

            if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
                m_toolSonAreaWidget->setFixedSize(appItemSize * m_toolSonLayout->count(), QWIDGETSIZE_MAX);
            } else {
                m_toolSonAreaWidget->setFixedSize(QWIDGETSIZE_MAX, appItemSize * m_toolSonLayout->count());
            }
        }

        if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
            m_toolAreaWidget->setFixedSize(m_multiWindowWidget->width() + m_toolSonAreaWidget->width(), QWIDGETSIZE_MAX);
        else
            m_toolAreaWidget->setFixedSize(QWIDGETSIZE_MAX, m_multiWindowWidget->height() + m_toolSonAreaWidget->height());
    } else {
        for (int i = 0; i < m_appAreaSonLayout->count(); ++i) {
            DockItem *dockItem = qobject_cast<DockItem *>(m_appAreaSonLayout->itemAt(i)->widget());
            if (!dockItem)
                continue;
            if (dockItem->itemType() == DockItem::ItemType::AppMultiWindow) {
                AppMultiItem *appMultiItem = qobject_cast<AppMultiItem *>(dockItem);
                dockItem->setFixedSize(appMultiItem->suitableSize(appItemSize));
            } else {
                dockItem->setFixedSize(appItemSize, appItemSize);
            }
        }
    }

    int appTopAndBottomMargin = 0;
    int appLeftAndRightMargin = 0;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        appTopAndBottomMargin = (m_fixedAreaWidget->height() - appItemSize) / 2;
    } else {
        appLeftAndRightMargin = (m_fixedAreaWidget->width() - appItemSize) / 2;
    }

    m_fixedAreaLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
    m_appAreaSonLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
}

void MainPanelControl::onRecentVisibleChanged(bool visible)
{
    m_appSpliter->setVisible(visible);
}

void MainPanelControl::onDockAppVisibleChanged(bool visible)
{
    m_fixedSpliter->setVisible(visible);
}

void MainPanelControl::onToolVisibleChanged(bool visible)
{
    m_recentSpliter->setVisible(visible);
}

void MainPanelControl::onTrayRequestUpdate()
{
    m_tray->layoutWidget();
    switch (m_position) {
    case Dock::Position::Left:
    case Dock::Position::Right: {
        m_tray->setFixedSize(QWIDGETSIZE_MAX, m_tray->suitableSize().height());
        break;
    }
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        m_tray->setFixedSize(m_tray->suitableSize().width(), QWIDGETSIZE_MAX);
        break;
    }
    }
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

/**
 * @brief MainWindow::appIsOnDock 判断指定的应用（驻留和运行显示在任务栏的所有应用）是否在任务栏上
 * @param appDesktop 应用的desktop文件的完整路径
 * @return true: 在任务栏；false: 不在任务栏
 */
bool MainPanelControl::appIsOnDock(const QString &appDesktop)
{
    return DockItemManager::instance()->appIsOnDock(appDesktop);
}
