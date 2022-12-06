// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainpanelcontrol.h"
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

#include <QDrag>
#include <QTimer>
#include <QStandardPaths>
#include <QString>
#include <QApplication>
#include <QPointer>
#include <QBoxLayout>
#include <QLabel>
#include <QPixmap>
#include <QtConcurrent/QtConcurrentRun>
#include <qpa/qplatformnativeinterface.h>
#include <qpa/qplatformintegration.h>
#define protected public
#include <QtGui/private/qsimpledrag_p.h>
#undef protected
#include <QtGui/private/qshapedpixmapdndwindow_p.h>
#include <QtGui/private/qguiapplication_p.h>

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
    , m_fixedAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_fixedSpliter(new QLabel(this))
    , m_appAreaWidget(new QWidget(this))
    , m_appAreaSonWidget(new QWidget(this))
    , m_appAreaSonLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_appSpliter(new QLabel(this))
    , m_trayAreaWidget(new QWidget(this))
    , m_trayAreaLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_traySpliter(new QLabel(this))
    , m_pluginAreaWidget(new QWidget(this))
    , m_pluginLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
    , m_desktopWidget(new DesktopWidget(this))
    , m_position(Position::Bottom)
    , m_placeholderItem(nullptr)
    , m_appDragWidget(nullptr)
    , m_dislayMode(Efficient)
    , m_tray(nullptr)
    , m_trashItem(nullptr)
{
    initUI();
    updateMainPanelLayout();
    setAcceptDrops(true);
    setMouseTracking(true);

    m_appAreaWidget->installEventFilter(this);
    m_appAreaSonWidget->installEventFilter(this);
    m_trayAreaWidget->installEventFilter(this);
    m_pluginAreaWidget->installEventFilter(this);

    //在设置每条线大小前，应该设置fixedsize(0,0)
    //应为paintEvent函数会先调用设置背景颜色，大小为随机值
    m_fixedSpliter->setFixedSize(0,0);
    m_appSpliter ->setFixedSize(0,0);
    m_traySpliter->setFixedSize(0,0);
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

    /* 托盘区域 */
    m_trayAreaWidget->setObjectName("trayarea");
    m_trayAreaWidget->setLayout(m_trayAreaLayout);
    m_trayAreaLayout->setSpacing(0);
    m_trayAreaLayout->setContentsMargins(0, 10, 0, 10);
    m_mainPanelLayout->addWidget(m_trayAreaWidget);

    m_traySpliter->setObjectName("spliter_tray");
    m_mainPanelLayout->addWidget(m_traySpliter);

    /* 插件区域 */
    m_pluginAreaWidget->setObjectName("pluginarea");
    m_pluginAreaWidget->setLayout(m_pluginLayout);
    m_pluginLayout->setSpacing(10);
    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->addWidget(m_pluginAreaWidget, 0, Qt::AlignCenter);

    /* 桌面预览 */
    m_desktopWidget->setObjectName("showdesktoparea");
    m_mainPanelLayout->addWidget(m_desktopWidget);

    m_mainPanelLayout->setSpacing(0);
    m_mainPanelLayout->setContentsMargins(0, 0, 0, 0);
    m_mainPanelLayout->setAlignment(m_fixedSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_appSpliter, Qt::AlignCenter);
    m_mainPanelLayout->setAlignment(m_traySpliter, Qt::AlignCenter);
}

/**
 * @brief MainPanelControl::setDisplayMode 根据任务栏显示模式更新界面显示，如果是时尚模式，没有‘显示桌面'区域，否则就有
 * @param dislayMode 任务栏显示模式
 */
void MainPanelControl::setDisplayMode(DisplayMode dislayMode)
{
    if (dislayMode == m_dislayMode)
        return;
    m_dislayMode = dislayMode;
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

    // 显示桌面的区域
    resizeDesktopWidget();

    // 设置任务栏各区域图标大小
    resizeDockIcon();
}

/**往固定区域添加应用
 * @brief MainPanelControl::addFixedAreaItem
 * @param index　位置索引，如果为负数则插入到最后，为正则插入到指定位置
 * @param wdg　应用指针对象
 */
void MainPanelControl::addFixedAreaItem(int index, QWidget *wdg)
{
    if(m_position == Position::Top || m_position == Position::Bottom){
        wdg->setMaximumSize(height(),height());
    } else {
        wdg->setMaximumSize(width(),width());
    }
    m_fixedAreaLayout->insertWidget(index, wdg);
}

/**往应用区域添加应用
 * @brief MainPanelControl::addAppAreaItem
 * @param index　位置索引，如果为负数则插入到最后，为正则插入到指定位置
 * @param wdg　应用指针对象
 */
void MainPanelControl::addAppAreaItem(int index, QWidget *wdg)
{
    if(m_position == Position::Top || m_position == Position::Bottom){
        wdg->setMaximumSize(height(),height());
    } else {
        wdg->setMaximumSize(width(),width());
    }
    m_appAreaSonLayout->insertWidget(index, wdg);
}

/**往托盘插件区域添加应用
 * @brief MainPanelControl::addTrayAreaItem
 * @param index　位置索引，如果为负数则插入到最后，为正则插入到指定位置
 * @param wdg　应用指针对象
 */
void MainPanelControl::addTrayAreaItem(int index, QWidget *wdg)
{
    m_tray = static_cast<TrayPluginItem *>(wdg);
    m_trayAreaLayout->insertWidget(index, wdg);
}

/**往插件区域添加应用，保存回收站插件指针对象
 * @brief MainPanelControl::addPluginAreaItem
 * @param index　位置索引，如果为负数则插入到最后，为正则插入到指定位置
 * @param wdg　应用指针对象
 */
void MainPanelControl::addPluginAreaItem(int index, QWidget *wdg)
{
    //因为日期时间插件和其他插件的大小有异，为了方便设置边距，在插件区域布局再添加一层布局设置边距
    //因此在处理插件图标时，需要通过两层布局判断是否为需要的插件，例如拖动插件位置等判断
    QBoxLayout * boxLayout = new QBoxLayout(QBoxLayout::LeftToRight, this);
    boxLayout->addWidget(wdg, 0, Qt::AlignCenter);
    m_pluginLayout->insertLayout(index, boxLayout, 0);

    // 保存垃圾箱插件指针
    PluginsItem *pluginsItem = qobject_cast<PluginsItem *>(wdg);
    if (pluginsItem && pluginsItem->pluginName() == "trash")
        m_trashItem = pluginsItem;
}

/**移除固定区域某一应用
 * @brief MainPanelControl::removeFixedAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removeFixedAreaItem(QWidget *wdg)
{
    m_fixedAreaLayout->removeWidget(wdg);
}

/**移除应用区域某一应用
 * @brief MainPanelControl::removeAppAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removeAppAreaItem(QWidget *wdg)
{
    m_appAreaSonLayout->removeWidget(wdg);
}

/**移除托盘插件区域某一应用
 * @brief MainPanelControl::removeTrayAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removeTrayAreaItem(QWidget *wdg)
{
    m_trayAreaLayout->removeWidget(wdg);
}

/**移除插件区域某一应用
 * @brief MainPanelControl::removePluginAreaItem
 * @param wdg 应用指针对象
 */
void MainPanelControl::removePluginAreaItem(QWidget *wdg)
{
    //因为日期时间插件大小和其他插件有异，为了方便设置边距，各插件中增加一层布局
    //因此remove插件图标时，需要从多的一层布局中取widget进行判断是否需要移除的插件
    // 清空保存的垃圾箱插件指针
    PluginsItem *pluginsItem = qobject_cast<PluginsItem *>(wdg);
    if (pluginsItem && pluginsItem->pluginName() == "trash")
        m_trashItem = nullptr;

    for (int i = 0; i < m_pluginLayout->count(); ++i) {
        QLayoutItem *layoutItem = m_pluginLayout->itemAt(i);
        QLayout *boxLayout = layoutItem->layout();
        if (boxLayout && boxLayout->itemAt(0)->widget() == wdg) {
            boxLayout->removeWidget(wdg);
            m_pluginLayout->removeItem(layoutItem);
            delete layoutItem;
            layoutItem = nullptr;
        }
    }
}

void MainPanelControl::resizeEvent(QResizeEvent *event)
{
    //先通过消息循环让各部件调整好size后再计算图标大小
    //避免因为部件size没有调整完导致计算的图标大小不准确
    //然后重复触发m_pluginAreaWidget的reszie事件并重复计算，造成任务栏图标抖动问题
    QWidget::resizeEvent(event);
    resizeDesktopWidget();
    resizeDockIcon();
}

/**根据任务栏所在位置， 设置应用区域控件的大小
 * @brief MainPanelControl::updateAppAreaSonWidgetSize
 */
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

/**
 * @brief setPositonValue 根据传入的位置更新界面布局，比如任务栏在左，布局应该是上下布局，任务栏在下，应该是左右布局
 * @param position 任务栏的位置
 */
void MainPanelControl::setPositonValue(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    QTimer::singleShot(0, this, &MainPanelControl::updateMainPanelLayout);
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
        addAppAreaItem(index, item);
        break;
    case DockItem::TrayPlugin: // 此处只会有一个tray系统托盘插件，微信、声音、网络蓝牙等等，都在系统托盘插件中处理的
        addTrayAreaItem(index, item);
        break;
    case DockItem::Plugins:
        addPluginAreaItem(index, item);
        break;
    }

    // 同removeItem处 注意:不能屏蔽此接口，否则会造成插件插入时无法显示
    if (item->itemType() != DockItem::App)
        resizeDockIcon();

    item->checkEntry();
    QTimer::singleShot(0, this, &MainPanelControl::updatePluginsLayout);
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
        removeAppAreaItem(item);
        break;
    case DockItem::TrayPlugin:
        removeTrayAreaItem(item);
        break;
    case DockItem::Plugins:
        removePluginAreaItem(item);
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
    int idx = getItemIndex(targetItem);
    if (-1 == idx)
        return;

    // remove old item
    removeItem(sourceItem);

    // insert new position
    if (sourceItem->isDragging()) {
        m_dragIndex = idx;
    }
    insertItem(idx, sourceItem);
}

int MainPanelControl::getItemIndex(DockItem *targetItem) const
{
    if (!targetItem)
        return -1;

    if (targetItem->itemType() == DockItem::App)
        return m_appAreaSonLayout->indexOf(targetItem);

    if (targetItem->itemType() == DockItem::Plugins){
        //因为日期时间插件大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
        //因此有拖动图标时，需要从多的一层布局中判断是否相同插件而获取插件位置顺序
        for (int i = 0; i < m_pluginLayout->count(); ++i) {
            QLayout *layout = m_pluginLayout->itemAt(i)->layout();
            if (layout && layout->itemAt(0)->widget() == targetItem) {
                return i;
            }
        }
    }

    if (targetItem->itemType() == DockItem::FixedPlugin)
        return m_fixedAreaLayout->indexOf(targetItem);

    return -1;
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

    if (appIsOnDock(DragmineData->data(m_draggingMimeKey))) {
        e->setAccepted(false);
        return;
    }

    handleDragMove(e, false);
}

bool MainPanelControl::eventFilter(QObject *watched, QEvent *event)
{
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
    if (watched == m_pluginAreaWidget) {
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
    QPointer<DockItem> item = dockItem;
    const QPixmap pixmap = item->grab();

    item->setDraging(true);
    item->update();

    QDrag *drag = nullptr;
    if (item->itemType() == DockItem::App) {
        AppDrag *appDrag = new AppDrag(item);

        m_appDragWidget = appDrag->appDragWidget();

        if (Utils::IS_WAYLAND_DISPLAY) {
            m_dragIndex = getItemIndex(dockItem);
            connect(m_appDragWidget, &AppDragWidget::requestRemoveSelf, this, [ = ](bool needDelete) {
                m_appDragWidget = nullptr;
                if (!item.isNull() && qobject_cast<AppItem *>(item)->isValid()) {
                    if (m_dragIndex > -1 && !needDelete) {
                        insertItem(m_dragIndex, item);
                        m_dragIndex = -1;
                    } else {
                        if (-1 == m_appAreaSonLayout->indexOf(item) && m_dragIndex != -1) {
                            insertItem(m_dragIndex, item);
                            m_dragIndex = -1;
                        }
                        AppItem *app = static_cast<AppItem *>(item.data());
                        app->undock();
                    }
                    item->setDraging(false);
                    item->update();
                }
            });
        } else {
            connect(m_appDragWidget, &AppDragWidget::destroyed, this, [ = ] {
                m_appDragWidget = nullptr;
                if (!item.isNull() && qobject_cast<AppItem *>(item)->isValid()) {
                    if (-1 == m_appAreaSonLayout->indexOf(item) && m_dragIndex != -1) {
                        insertItem(m_dragIndex, item);
                        m_dragIndex = -1;
                    }
                    item->setDraging(false);
                    item->update();
                }
            });
        }

        connect(m_appDragWidget, &AppDragWidget::requestRemoveItem, this, [ = ] {
            if (-1 != m_appAreaSonLayout->indexOf(item)) {
                m_dragIndex = m_appAreaSonLayout->indexOf(item);
                removeItem(item);
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
    if (item->itemType() == DockItem::Plugins && m_trashItem && dockItem != m_trashItem) {
        m_trashItem->centralWidget()->setAcceptDrops(false);
        isNeedBack = true;
    }

    drag->setMimeData(new QMimeData);
    drag->exec(Qt::MoveAction);

    if (item->itemType() != DockItem::App || m_dragIndex == -1) {
        m_appDragWidget = nullptr;
        item->setDraging(false);
        item->update();

        // isNeedBack是否需要设置垃圾箱插件AcceptDrops true
        if (isNeedBack)
            m_trashItem->centralWidget()->setAcceptDrops(true);
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
        targetItem = sourceItem;
    }

    return targetItem;
}


void MainPanelControl::updateDisplayMode()
{
    moveAppSonWidget();
    resizeDesktopWidget();
}

/**把驻留应用和被打开的应用所在窗口移动到指定位置
 * @brief MainPanelControl::moveAppSonWidget
 */
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

/**通知布局器，控件已发生变化，需要重新设置几何位置
 * @brief MainPanelControl::updatePluginsLayout
 */
void MainPanelControl::updatePluginsLayout()
{
    for (int i = 0; i < m_pluginLayout->count(); ++i) {
        QLayout *layout = m_pluginLayout->itemAt(i)->layout();
        if (layout) {
            PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());
            if (pItem && pItem->sizeHint().width() != -1) {
                pItem->updateGeometry();
            }
        }
    }
}

void MainPanelControl::setToggleDesktopInterval(int ms)
{
    if (m_desktopWidget)
        m_desktopWidget->setToggleDesktopInterval(ms);
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

    painter.fillRect(m_fixedSpliter->geometry(), color);
    painter.fillRect(m_appSpliter->geometry(), color);
    painter.fillRect(m_traySpliter->geometry(), color);
}

/**重新计算任务栏上应用图标、插件图标的大小，并设置
 * @brief MainPanelControl::resizeDockIcon
 */
void MainPanelControl::resizeDockIcon()
{
    // 总宽度
    int totalLength = ((m_position == Position::Top) || (m_position == Position::Bottom)) ? width() : height();

    // 减去托盘间隔区域
    if (m_tray) {
        totalLength -= (m_tray->trayVisibleItemCount() + 1) * 10;
    }

    // 减去3个分割线的宽度
    totalLength -= 3 * SPLITER_SIZE;

    // 减去显示桌面图标宽度
    totalLength -= ((m_position == Position::Top) || (m_position == Position::Bottom)) ? m_desktopWidget->width() : m_desktopWidget->height();

    int pluginItemCount = 0;
    int calcPluginItemCount = 0;

    // 因为日期时间大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
    // 因此需要通过多一层布局来获取各插件
    for (int i = 0; i < m_pluginLayout->count(); ++ i) {
        QLayout *layout = m_pluginLayout->itemAt(i)->layout();
        if (layout) {
            PluginsItem *w = static_cast<PluginsItem *>(layout->itemAt(0)->widget());
            if (w) {
                // 如果插件大小由自己决定,则不参与计算需要减去其宽度,其他插件则需要参与计算并计数
                if ((m_position == Position::Top || m_position == Position::Bottom) && (w->sizeHint().height() != -1)) {
                    totalLength -= w->width();
                } else if ((m_position == Position::Top || m_position == Position::Bottom) && (w->sizeHint().width() != -1)) {
                    totalLength -= w->height();
                } else {
                    calcPluginItemCount ++;
                }

                // 所有插件个数,用于计算插件之间的间隔之和
                pluginItemCount ++;
            }
        }
    }

    // 减去插件间隔大小, 只有一个插件或没有插件都是间隔20,2个或以上每多一个插件多间隔10
    if (pluginItemCount > 1)
       totalLength -= (pluginItemCount + 1) * 10;
    else
       totalLength -= 20;

    if (totalLength < 0)
        return;

    // 参与计算的插件的个数包含托盘和插件
    int pluginCount = m_tray ? m_tray->trayVisibleItemCount() + calcPluginItemCount : calcPluginItemCount;

    // 需要计算的图标总数
    int iconCount = m_fixedAreaLayout->count() + m_appAreaSonLayout->count() + pluginCount;
    if (iconCount <= 0)
        return;

    // 余数
    int yu = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    int iconSize = (totalLength - yu) / iconCount;

    // 计算插件图标的最大或最小值
    int tray_item_size = qBound(20, iconSize, 40);
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        tray_item_size = qMin(tray_item_size,height());
        tray_item_size = std::min(tray_item_size, height() - 20);
    } else {
        tray_item_size = qMin(tray_item_size,width());
        tray_item_size = std::min(tray_item_size, width() - 20);
    }

    if (tray_item_size < 20)
        tray_item_size = 20;

    // 减去插件图标的大小后重新计算固定图标和应用图标的平均大小
    totalLength -= tray_item_size * pluginCount;
    iconCount -= pluginCount;

    // 余数
    yu = (totalLength % iconCount);
    // icon宽度 = (总宽度-余数)/icon个数
    iconSize = (totalLength - yu) / iconCount;

    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        if (iconSize >= height()) {
            calcuDockIconSize(height(), height(), tray_item_size);
        } else {
            calcuDockIconSize(iconSize, height(), tray_item_size);
        }
    } else {
        if (iconSize >= width()) {
            calcuDockIconSize(width(), width(), tray_item_size);
        } else {
            calcuDockIconSize(width(), iconSize, tray_item_size);
        }
    }
}

void MainPanelControl::calcuDockIconSize(int w, int h, int traySize)
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

    if (m_tray) {
        m_tray->centralWidget()->setProperty("iconSize", traySize);
    }

    //因为日期时间大小和其他插件大小有异，为了设置边距，在各插件中增加了一层布局
    //因此需要通过多一层布局来获取各插件
    if ((m_position == Position::Top) || (m_position == Position::Bottom)) {
        // 三方插件
        for (int i = 0; i < m_pluginLayout->count(); ++ i) {
            QLayout *layout = m_pluginLayout->itemAt(i)->layout();
            if (layout && layout->itemAt(0)) {
                PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());
                if (pItem) {
                    if (pItem->sizeHint().height() == -1) {
                        pItem->setFixedSize(traySize, traySize);
                    } else if (pItem->sizeHint().height() > height()) {
                        pItem->resize(pItem->width(), height());
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
                if (pItem) {
                    if (pItem->sizeHint().width() == -1) {
                        pItem->setFixedSize(traySize, traySize);
                    } else if (pItem->sizeHint().width() > width()) {
                        pItem->resize(width(), pItem->height());
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
        trayTopAndBottomMargin = (m_trayAreaWidget->height() - traySize) / 2;
    } else {
        appLeftAndRightMargin = (m_fixedAreaWidget->width() - appItemSize) / 2;
        trayLeftAndRightMargin = (m_trayAreaWidget->width() - traySize) / 2;
    }

    m_fixedAreaLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
    m_appAreaSonLayout->setContentsMargins(appLeftAndRightMargin, appTopAndBottomMargin, appLeftAndRightMargin, appTopAndBottomMargin);
    m_trayAreaLayout->setContentsMargins(trayLeftAndRightMargin, trayTopAndBottomMargin, trayLeftAndRightMargin, trayTopAndBottomMargin);

    //因为日期时间插件或第三方插件声明自定义大小
    //而不对自定义大小插件设置边距
    for (int i = 0; i < m_pluginLayout->count(); ++ i) {
        QLayout *layout = m_pluginLayout->itemAt(i)->layout();
        if (layout && layout->itemAt(0)) {
            PluginsItem *pItem = static_cast<PluginsItem *>(layout->itemAt(0)->widget());

            if (pItem && pItem->pluginSizePolicy() != PluginsItemInterface::Custom) {
                layout->setContentsMargins(trayLeftAndRightMargin, trayTopAndBottomMargin, trayLeftAndRightMargin, trayTopAndBottomMargin);
            }
        }
    }
}

/**时尚模式没有‘显示桌面’区域
 * @brief MainPanelControl::resizeDesktopWidget
 */
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

/**
 * @brief MainWindow::appIsOnDock 判断指定的应用（驻留和运行显示在任务栏的所有应用）是否在任务栏上
 * @param appDesktop 应用的desktop文件的完整路径
 * @return true: 在任务栏；false: 不在任务栏
 */
bool MainPanelControl::appIsOnDock(const QString &appDesktop)
{
    return DockItemManager::instance()->appIsOnDock(appDesktop);
}
