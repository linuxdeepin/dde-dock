﻿/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "quickpluginwindow.h"
#include "quicksettingcontroller.h"
#include "quicksettingitem.h"
#include "pluginsiteminterface.h"
#include "quicksettingcontainer.h"
#include "appdrag.h"
#include "quickpluginmodel.h"
#include "quickdragcore.h"

#include <DStyleOption>
#include <DStandardItem>
#include <DGuiApplicationHelper>

#include <QDrag>
#include <QScrollBar>
#include <QStringList>
#include <QSize>
#include <QMouseEvent>
#include <QBoxLayout>
#include <QGuiApplication>
#include <QMenu>
#include <QDragLeaveEvent>

#define ITEMSIZE 22
#define ITEMSPACE 6
#define ICONWIDTH 18
#define ICONHEIGHT 16

typedef struct DragInfo{
    QPoint dragPoint;
    QuickDockItem *dockItem = nullptr;

    void reset() {
        dockItem = nullptr;
        dragPoint.setX(0);
        dragPoint.setY(0);
    }

    bool isNull() const {
        return (!dockItem);
    }

    bool canDrag(QPoint currentPoint) const {
        if (dragPoint.isNull())
            return false;

        if (!dragPixmap())
            return false;

        return (qAbs(currentPoint.x() - dragPoint.x()) >= 1 ||
                qAbs(currentPoint.y() - dragPoint.y()) >= 1);
    }

    QPixmap dragPixmap() const {
        if (!dockItem)
            return QPixmap();

        QPixmap pixmap = dockItem->pluginItem()->icon(DockPart::QuickShow).pixmap(QSize(ITEMSIZE, ITEMSIZE));
        if (!pixmap.isNull())
            return pixmap;

        QString itemKey = QuickSettingController::instance()->itemKey(dockItem->pluginItem());
        QWidget *itemWidget = dockItem->pluginItem()->itemWidget(itemKey);
        if (!itemWidget)
            return QPixmap();

        return itemWidget->grab();
    }
} DragInfo;

QuickPluginWindow::QuickPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QBoxLayout(QBoxLayout::RightToLeft, this))
    , m_position(Dock::Position::Bottom)
    , m_dragInfo(new DragInfo)
    , m_dragEnterMimeData(nullptr)
{
    initUi();
    initConnection();

    topLevelWidget()->installEventFilter(this);
    installEventFilter(this);
    setAcceptDrops(true);
    setMouseTracking(true);
}

QuickPluginWindow::~QuickPluginWindow()
{
    delete m_dragInfo;
}

void QuickPluginWindow::initUi()
{
    setAcceptDrops(true);
    m_mainLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    m_mainLayout->setDirection(QBoxLayout::RightToLeft);
    m_mainLayout->setContentsMargins(ITEMSPACE, 0, ITEMSPACE, 0);
    m_mainLayout->setSpacing(ITEMSPACE);
}

void QuickPluginWindow::setPositon(Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    QuickSettingContainer::setPosition(position);
    QuickDockItem::setPosition(position);
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        m_mainLayout->setDirection(QBoxLayout::RightToLeft);
        m_mainLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    } else {
        m_mainLayout->setDirection(QBoxLayout::BottomToTop);
        m_mainLayout->setAlignment(Qt::AlignTop | Qt::AlignHCenter);
    }
}

void QuickPluginWindow::dragPlugin(PluginsItemInterface *item)
{
    QuickPluginModel *quickModel = QuickPluginModel::instance();
    QPoint itemPoint = mapFromGlobal(QCursor::pos());
    // 查找移动后的位置，如果移动后的插件找不到，就直接放到最后
    int index = -1;
    QuickDockItem *targetWidget = qobject_cast<QuickDockItem *>(childAt(itemPoint));
    if (targetWidget) {
        // 如果是拖动到固定插件区域，也放到最后
        QList<PluginsItemInterface *> pluginItems = quickModel->dockedPluginItems();
        for (int i = 0; i < pluginItems.size(); i++) {
            PluginsItemInterface *plugin = pluginItems[i];
            if (quickModel->isFixed(plugin))
                continue;

            if (targetWidget->pluginItem() == plugin) {
                index = i;
                break;
            }
        }
    }

    quickModel->addPlugin(item, index);
}

QSize QuickPluginWindow::suitableSize() const
{
    return suitableSize(m_position);
}

QSize QuickPluginWindow::suitableSize(const Dock::Position &position) const
{
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        int itemWidth = 0;
        for (int i = 0; i < m_mainLayout->count(); i++) {
            QWidget *itemWidget = m_mainLayout->itemAt(i)->widget();
            if (itemWidget)
                itemWidth += itemWidget->width() + ITEMSPACE;
        }
        itemWidth += ITEMSPACE;

        return QSize(itemWidth, ITEMSIZE);
    }

    int itemHeight = 0;
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QWidget *itemWidget = m_mainLayout->itemAt(i)->widget();
        if (itemWidget)
            itemHeight += itemWidget->height() + ITEMSPACE;
    }
    itemHeight += ITEMSPACE;

    return QSize(ITEMSIZE, itemHeight);
}

PluginsItemInterface *QuickPluginWindow::findQuickSettingItem(const QPoint &mousePoint, const QList<PluginsItemInterface *> &settingItems)
{
    QuickDockItem *selectWidget = qobject_cast<QuickDockItem *>(childAt(mousePoint));
    if (!selectWidget)
        return nullptr;

    for (int i = 0; i < settingItems.size(); i++) {
        PluginsItemInterface *settingItem = settingItems[i];
        if (selectWidget->pluginItem() == settingItem)
            return settingItem;
    }

    return nullptr;
}

bool QuickPluginWindow::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == topLevelWidget()) {
        switch (event->type()) {
        case QEvent::DragEnter: {
            QDragEnterEvent *dragEvent = static_cast<QDragEnterEvent *>(event);
            dragEnterEvent(dragEvent);
            break;
        }
        case QEvent::DragLeave: {
            QDragLeaveEvent *dragEvent = static_cast<QDragLeaveEvent *>(event);
            dragLeaveEvent(dragEvent);
            break;
        }
        default:
            break;
        }
    }
    switch (event->type()) {
    case QEvent::MouseButtonPress: {
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        if (mouseEvent->button() != Qt::LeftButton)
            break;

        QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(watched);
        if (!dockItem)
            break;

        m_dragInfo->dockItem = dockItem;
        m_dragInfo->dragPoint = mouseEvent->pos();
        break;
    }
    case QEvent::MouseButtonRelease: {
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        if (mouseEvent->button() != Qt::LeftButton)
            break;

        if (m_dragInfo->isNull())
            break;

        do {
            if (m_dragInfo->canDrag(mouseEvent->pos()))
                break;

            showPopup(m_dragInfo->dockItem);
        } while (false);
        m_dragInfo->reset();

        break;
    }
    case QEvent::MouseMove: {
        if (m_dragInfo->isNull())
            break;

        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        if (m_dragInfo->canDrag(mouseEvent->pos()))
            startDrag();

        m_dragInfo->reset();
        break;
    }
    case QEvent::Drop: {
        m_dragEnterMimeData = nullptr;
        QDropEvent *dropEvent = static_cast<QDropEvent *>(event);
        if (qobject_cast<QuickSettingContainer *>(dropEvent->source())) {
            const QuickPluginMimeData *mimeData = qobject_cast<const QuickPluginMimeData *>(dropEvent->mimeData());
            if (mimeData)
                dragPlugin(mimeData->pluginItemInterface());
        }
        break;
    }
    default:
        break;
    }
    return QWidget::eventFilter(watched, event);
}

void QuickPluginWindow::dragEnterEvent(QDragEnterEvent *event)
{
    m_dragEnterMimeData = const_cast<QuickPluginMimeData *>(qobject_cast<const QuickPluginMimeData *>(event->mimeData()));
    if (m_dragEnterMimeData) {
        PluginsItemInterface *plugin = m_dragEnterMimeData->pluginItemInterface();
        QIcon icon = plugin->icon(DockPart::QuickShow);
        if (icon.isNull()) {
            QWidget *widget = plugin->itemWidget(QuickSettingController::instance()->itemKey(plugin));
            if (widget)
                icon = widget->grab();
        }
        QuickIconDrag *drag = qobject_cast<QuickIconDrag *>(m_dragEnterMimeData->drag());
        if (drag && !icon.isNull()) {
            QPixmap pixmap = icon.pixmap(QSize(16, 16));
            drag->updatePixmap(pixmap);
        }
        event->accept();
    } else {
        event->ignore();
    }
}

void QuickPluginWindow::dragLeaveEvent(QDragLeaveEvent *event)
{
    if (m_dragEnterMimeData) {
        QPoint mousePos = topLevelWidget()->mapFromGlobal(QCursor::pos());
        QuickIconDrag *drag = qobject_cast<QuickIconDrag *>(m_dragEnterMimeData->drag());
        if (!topLevelWidget()->rect().contains(mousePos) && drag) {
            drag->useSourcePixmap();
        }
        m_dragEnterMimeData = nullptr;
    }
    event->accept();
}

void QuickPluginWindow::onRequestUpdate()
{
    bool countChanged = false;
    QuickPluginModel *model = QuickPluginModel::instance();
    QList<PluginsItemInterface *> plugins = model->dockedPluginItems();
    // 先删除所有的widget
    QMap<PluginsItemInterface *, QuickDockItem *> pluginItems;
    for (int i = m_mainLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (!layoutItem)
            continue;

        QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
        if (!dockItem)
            continue;

        dockItem->setParent(nullptr);
        m_mainLayout->removeItem(layoutItem);
        if (plugins.contains(dockItem->pluginItem())) {
            // 如果该插件在任务栏上，则先将其添加到临时列表中
            pluginItems[dockItem->pluginItem()] = dockItem;
        } else {
            // 如果该插件不在任务栏上，则先删除
            dockItem->deleteLater();
            countChanged = true;
        }
    }

    // 将列表中所有的控件按照顺序添加到布局上
    QuickSettingController *quickController = QuickSettingController::instance();
    for (PluginsItemInterface *item : plugins) {
        QuickDockItem *itemWidget = nullptr;
        if (pluginItems.contains(item)) {
            itemWidget = pluginItems[item];
        } else {
            itemWidget = new QuickDockItem(item, quickController->itemKey(item), this);
            itemWidget->setFixedSize(itemWidget->suitableSize());
            itemWidget->installEventFilter(this);
            itemWidget->setMouseTracking(true);
            countChanged = true;
        }
        itemWidget->setParent(this);
        m_mainLayout->addWidget(itemWidget);
    }

    if (countChanged)
        Q_EMIT itemCountChanged();
}

QPoint QuickPluginWindow::popupPoint(QWidget *widget) const
{
    QWidget *itemWidget = widget;
    if (!itemWidget && m_mainLayout->count() > 0)
        itemWidget = m_mainLayout->itemAt(0)->widget();

    if (!itemWidget)
        return QPoint();

    QPoint pointCurrent = itemWidget->mapToGlobal(QPoint(0, 0));
    switch (m_position) {
    case Dock::Position::Bottom: {
        // 在下方的时候，Y坐标设置在顶层窗口的y值，保证下方对齐
        pointCurrent.setX(pointCurrent.x() + itemWidget->width() / 2);
        pointCurrent.setY(topLevelWidget()->y());
        break;
    }
    case Dock::Position::Top: {
        // 在上面的时候，Y坐标设置为任务栏的下方，保证上方对齐
        pointCurrent.setX(pointCurrent.x() + itemWidget->width() / 2);
        pointCurrent.setY(topLevelWidget()->y() + topLevelWidget()->height());
        break;
    }
    case Dock::Position::Left: {
        // 在左边的时候，X坐标设置在顶层窗口的最右侧，保证左对齐
        pointCurrent.setX(topLevelWidget()->x() + topLevelWidget()->width());
        pointCurrent.setY(pointCurrent.y() + itemWidget->height() / 2);
        break;
    }
    case Dock::Position::Right: {
        // 在右边的时候，X坐标设置在顶层窗口的最左侧，保证右对齐
        pointCurrent.setX(topLevelWidget()->x());
        pointCurrent.setY(pointCurrent.y() + itemWidget->height() / 2);
    }
    }
    return pointCurrent;
}

void QuickPluginWindow::onUpdatePlugin(PluginsItemInterface *itemInter, const DockPart &dockPart)
{
    //update plugin status
    if (dockPart != DockPart::QuickShow)
        return;

    QuickDockItem *quickDockItem = getDockItemByPlugin(itemInter);
    if (quickDockItem) {
        quickDockItem->setFixedSize(quickDockItem->suitableSize());
        quickDockItem->update();
    }
}

void QuickPluginWindow::onRequestAppletVisible(PluginsItemInterface *itemInter, const QString &itemKey, bool visible)
{
    if (visible)
        showPopup(getDockItemByPlugin(itemInter), itemInter, itemInter->itemPopupApplet(itemKey), false);
}

void QuickPluginWindow::startDrag()
{
    if (!m_dragInfo->dockItem)
        return;

    PluginsItemInterface *moveItem = m_dragInfo->dockItem->pluginItem();
    //AppDrag *drag = new AppDrag(this, new QuickDragWidget);
    QDrag *drag = new QDrag(this);
    QuickPluginMimeData *mimedata = new QuickPluginMimeData(moveItem, drag);
    drag->setMimeData(mimedata);
    QPixmap dragPixmap = m_dragInfo->dragPixmap();
    drag->setPixmap(dragPixmap);

    drag->setHotSpot(dragPixmap.rect().center());

    drag->exec(Qt::CopyAction);
    // 获取当前鼠标在任务栏快捷图标区域的位置
    QPoint currentPoint = mapFromGlobal(QCursor::pos());
    // 获取区域图标插入的位置
    QuickPluginModel::instance()->addPlugin(mimedata->pluginItemInterface(), getDropIndex(currentPoint));
}

QuickDockItem *QuickPluginWindow::getDockItemByPlugin(PluginsItemInterface *item)
{
    if (!item)
        return nullptr;

    for (int i = 0; i < m_mainLayout->count(); i++) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (!layoutItem)
            continue;

        QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
        if (!dockItem)
            continue;

        if (dockItem->pluginItem() == item)
            return dockItem;
    }

    return nullptr;
}

QuickDockItem *QuickPluginWindow::getActiveDockItem(QPoint point) const
{
    QuickDockItem *selectWidget = qobject_cast<QuickDockItem *>(childAt(point));
    if (!selectWidget)
        return nullptr;

    // 如果当前图标是固定插件，则不让插入
    if (QuickPluginModel::instance()->isFixed(selectWidget->pluginItem()))
        return nullptr;

    return selectWidget;
}

void QuickPluginWindow::showPopup(QuickDockItem *item, PluginsItemInterface *itemInter, QWidget *childPage, bool isClicked)
{
    if (!isVisible())
        return;

    bool canBack = true;
    DockPopupWindow *popWindow = QuickSettingContainer::popWindow();
    if (isClicked && popWindow->isVisible()) {
        // 如果是点击插件，并且该插件曾经打开快捷面板且已经是显示状态，那么就直接隐藏快捷面板
        popWindow->hide();
        return;
    }

    if (!popWindow->isVisible()) {
        if (Utils::IS_WAYLAND_DISPLAY) {
            // TODO: 临时解决方案，如果是wayland环境，toolTip没有消失，因此，此处直接调用接口来隐藏
            for (int i = m_mainLayout->count() - 1; i >= 0; i--) {
                QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
                if (!layoutItem)
                    continue;

                QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
                if (!dockItem)
                    continue;

                dockItem->hideToolTip();
            }
        }

        popWindow->setExtendWidget(item);
        popWindow->show(popupPoint(item), true);
        canBack = false;
    }

    QuickSettingContainer *container = static_cast<QuickSettingContainer *>(popWindow->getContent());
    container->showPage(childPage, itemInter, canBack);
    popWindow->raise();
}

QList<QuickDockItem *> QuickPluginWindow::quickDockItems()
{
    QList<QuickDockItem *> dockItems;
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (!layoutItem)
            continue;

        QuickDockItem *dockedItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
        if (!dockedItem)
            continue;

        dockItems << dockedItem;
    }

    return dockItems;
}

int QuickPluginWindow::getDropIndex(QPoint point)
{
    QList<QuickDockItem *> dockedItems = quickDockItems();
    QuickDockItem *targetItem = getActiveDockItem(point);
    if (targetItem) {
        for (int i = 0; i < dockedItems.count(); i++) {
            if (dockedItems[i] == targetItem)
                return i;
        }

        return -1;
    }

    QList<PluginsItemInterface *> dockItemInter = QuickPluginModel::instance()->dockedPluginItems();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        // 上下方向从右向左排列
        for (int i = 0; i < dockItemInter.count() - 1; i++) {
            QuickDockItem *dockBeforeItem = dockedItems[i];
            QuickDockItem *dockItem = dockedItems[i + 1];
            if (!dockItem->canInsert())
                continue;

            if (dockBeforeItem->geometry().x() > point.x() && dockItem->geometry().right() < point.x())
                return i;
        }
    } else {
        // 左右方向从下向上排列
        for (int i = 0; i < dockItemInter.count() - 1; i++) {
            QuickDockItem *dockBeforeItem = dockedItems[i];
            QuickDockItem *dockItem = dockedItems[i + 1];
            if (!dockItem->canInsert())
                continue;

            if (dockBeforeItem->geometry().bottom() > point.y() && dockItem->geometry().top() < point.y())
                return i;
        }
    }
    // 如果都没有找到，直接插入到最后
    return -1;
}

void QuickPluginWindow::dragMoveEvent(QDragMoveEvent *event)
{
    event->accept();
}

void QuickPluginWindow::initConnection()
{
    QuickPluginModel *model = QuickPluginModel::instance();
    connect(model, &QuickPluginModel::requestUpdate, this, &QuickPluginWindow::onRequestUpdate);
    connect(model, &QuickPluginModel::requestUpdatePlugin, this, &QuickPluginWindow::onUpdatePlugin);
    connect(QuickSettingController::instance(), &QuickSettingController::requestAppletVisible, this, &QuickPluginWindow::onRequestAppletVisible);
}

/**
 * @brief QuickDockItem::QuickDockItem
 * @param pluginItem
 * @param parent
 */

Dock::Position QuickDockItem::m_position(Dock::Position::Bottom);

QuickDockItem::QuickDockItem(PluginsItemInterface *pluginItem, const QString &itemKey, QWidget *parent)
    : QWidget(parent)
    , m_pluginItem(pluginItem)
    , m_itemKey(itemKey)
    , m_popupWindow(new DockPopupWindow)
    , m_contextMenu(new QMenu(this))
    , m_tipParent(nullptr)
    , m_mainLayout(nullptr)
    , m_dockItemParent(nullptr)
{
    initUi();
    initConnection();
    initAttribute();
}

QuickDockItem::~QuickDockItem()
{
    QWidget *tipWidget = m_pluginItem->itemTipsWidget(m_itemKey);
    if (tipWidget && tipWidget->parentWidget() == m_popupWindow)
        tipWidget->setParent(m_tipParent);

    m_popupWindow->deleteLater();
}

void QuickDockItem::setPosition(Dock::Position position)
{
    m_position = position;
}

PluginsItemInterface *QuickDockItem::pluginItem()
{
    return m_pluginItem;
}

bool QuickDockItem::canInsert() const
{
    return (m_pluginItem->flags() & PluginFlag::Attribute_CanInsert);
}

void QuickDockItem::hideToolTip()
{
    m_popupWindow->hide();
}

QSize QuickDockItem::suitableSize() const
{
    if (m_pluginItem->pluginSizePolicy() == PluginsItemInterface::PluginSizePolicy::Custom) {
        QPixmap pixmap = iconPixmap();
        if (!pixmap.isNull())
            return pixmap.size();

        QWidget *itemWidget = m_pluginItem->itemWidget(m_itemKey);
        if (itemWidget) {
            int itemWidth = ICONWIDTH;
            int itemHeight = ICONHEIGHT;
            QSize itemSize = itemWidget->sizeHint();
            if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
                if (itemSize.width() > 0)
                    itemWidth = itemSize.width();
                if (itemSize.height() > 0 && itemSize.height() <= topLevelWidget()->height())
                    itemHeight = itemSize.height();
            } else {
                if (itemSize.width() > 0 && itemSize.width() < topLevelWidget()->width())
                    itemWidth = itemSize.width();
                if (itemSize.height() > 0 && itemSize.height() < ICONHEIGHT)
                    itemHeight = itemSize.height();
            }

            return QSize(itemWidth, itemHeight);
        }
    }

    return QSize(ICONWIDTH, ICONHEIGHT);
}

void QuickDockItem::paintEvent(QPaintEvent *event)
{
    if (!m_pluginItem)
        return QWidget::paintEvent(event);

    QPixmap pixmap = iconPixmap();
    int width = ICONWIDTH;
    int height = ICONHEIGHT;
    if (m_pluginItem->pluginSizePolicy() == PluginsItemInterface::PluginSizePolicy::Custom) {
        width = pixmap.width();
        height = pixmap.height();
    }
    QRect pixmapRect = QRect(QPoint((rect().width() - width) / 2, (rect().height() - height) / 2), pixmap.size());

    QPainter painter(this);
    painter.drawPixmap(pixmapRect, pixmap);
}

void QuickDockItem::mousePressEvent(QMouseEvent *event)
{
    if (event->button() != Qt::RightButton)
        return QWidget::mousePressEvent(event);

    if (m_contextMenu->actions().isEmpty()) {
        const QString menuJson = m_pluginItem->itemContextMenu(m_itemKey);
        if (menuJson.isEmpty())
            return;

        QJsonDocument jsonDocument = QJsonDocument::fromJson(menuJson.toLocal8Bit().data());
        if (jsonDocument.isNull())
            return;

        QJsonObject jsonMenu = jsonDocument.object();

        QJsonArray jsonMenuItems = jsonMenu.value("items").toArray();
        for (auto item : jsonMenuItems) {
            QJsonObject itemObj = item.toObject();
            QAction *action = new QAction(itemObj.value("itemText").toString());
            action->setCheckable(itemObj.value("isCheckable").toBool());
            action->setChecked(itemObj.value("checked").toBool());
            action->setData(itemObj.value("itemId").toString());
            action->setEnabled(itemObj.value("isActive").toBool());
            m_contextMenu->addAction(action);
        }
    }

    m_contextMenu->exec(QCursor::pos());
}

void QuickDockItem::enterEvent(QEvent *event)
{
    QWidget::enterEvent(event);

    QWidget *tipWidget = m_pluginItem->itemTipsWidget(m_itemKey);
    if (!tipWidget)
        return;

    // 记录下toolTip的parent，因为在调用DockPopupWindow的时候会将DockPopupWindow设置为toolTip的parent,
    // 在DockPopupWindow对象释放的时候, 会将toolTip也一起给释放
    if (tipWidget->parentWidget() != m_popupWindow)
        m_tipParent = tipWidget->parentWidget();

    switch (m_position) {
    case Top:
        m_popupWindow->setArrowDirection(DockPopupWindow::ArrowTop);
        break;
    case Bottom:
        m_popupWindow->setArrowDirection(DockPopupWindow::ArrowBottom);
        break;
    case Left:
        m_popupWindow->setArrowDirection(DockPopupWindow::ArrowLeft);
        break;
    case Right:
        m_popupWindow->setArrowDirection(DockPopupWindow::ArrowRight);
        break;
    }

    m_popupWindow->resize(tipWidget->sizeHint());
    m_popupWindow->setContent(tipWidget);

    m_popupWindow->show(popupMarkPoint());
}

void QuickDockItem::leaveEvent(QEvent *event)
{
    QWidget::leaveEvent(event);
    m_popupWindow->hide();
}

void QuickDockItem::showEvent(QShowEvent *event)
{
    if (!m_mainLayout)
        return QWidget::showEvent(event);

    QWidget *itemWidget = m_pluginItem->itemWidget(m_itemKey);
    if (itemWidget && m_mainLayout->indexOf(itemWidget) < 0) {
        itemWidget->show();
        itemWidget->setFixedSize(size());
        m_mainLayout->addWidget(itemWidget);
    }
}

void QuickDockItem::hideEvent(QHideEvent *event)
{
    if (!m_mainLayout)
        return QWidget::hideEvent(event);

    QWidget *itemWidget = m_pluginItem->itemWidget(m_itemKey);
    if (itemWidget && m_mainLayout->indexOf(itemWidget) >= 0) {
        itemWidget->setParent(m_dockItemParent);
        itemWidget->hide();
        m_mainLayout->removeWidget(itemWidget);
    }
}

bool QuickDockItem::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == this)
        return m_pluginItem->eventHandler(event);

    return QWidget::eventFilter(watched, event);
}

QPixmap QuickDockItem::iconPixmap() const
{
    QIcon icon = m_pluginItem->icon(DockPart::QuickShow);
    if (!icon.isNull()) {
        if (icon.availableSizes().size() > 0) {
            QSize size = icon.availableSizes().first();
            return icon.pixmap(size);
        }
        int pixmapSize = static_cast<int>(ICONHEIGHT * (QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? 1 : qApp->devicePixelRatio()));
        return icon.pixmap(pixmapSize, pixmapSize);
    }

    return QPixmap();
}

void QuickDockItem::initUi()
{
    QPixmap pixmap = iconPixmap();
    if (!pixmap.isNull())
        return;

    m_mainLayout = new QHBoxLayout(this);
    m_mainLayout->setContentsMargins(0, 0, 0, 0);
    QWidget *itemWidget = m_pluginItem->itemWidget(m_itemKey);
    if (itemWidget) {
        m_dockItemParent = itemWidget->parentWidget();
        itemWidget->installEventFilter(this);
    }
}

void QuickDockItem::initAttribute()
{
    m_popupWindow->setShadowBlurRadius(20);
    m_popupWindow->setRadius(6);
    m_popupWindow->setShadowYOffset(2);
    m_popupWindow->setShadowXOffset(0);
    m_popupWindow->setArrowWidth(18);
    m_popupWindow->setArrowHeight(10);
    m_popupWindow->setObjectName("quickitempopup");
    if (Utils::IS_WAYLAND_DISPLAY) {
        Qt::WindowFlags flags = m_popupWindow->windowFlags() | Qt::FramelessWindowHint;
        m_popupWindow->setWindowFlags(flags);
    }

    this->installEventFilter(this);
}

void QuickDockItem::initConnection()
{
    connect(m_contextMenu, &QMenu::triggered, this, &QuickDockItem::onMenuActionClicked);
    connect(qApp, &QApplication::aboutToQuit, m_popupWindow, &DockPopupWindow::deleteLater);
}

QPoint QuickDockItem::topleftPoint() const
{
    QPoint p = this->pos();
    /* 由于点击范围的问题，在图标的外面加了一层布局，这个布局的边距需要考虑 */
    switch (m_position) {
    case Top:
        p.setY(p.y() * 2);
        break;
    case Bottom:
        p.setY(0);
        break;
    case Left:
        p.setX(p.x() * 2);
        break;
    case Right:
        p.setX(0);
        break;
    }

    QWidget *w = qobject_cast<QWidget *>(this->parent());
    while (w) {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    }

    return p;
}

QPoint QuickDockItem::popupMarkPoint() const
{
    QPoint p(topleftPoint());
    const QRect r = rect();
    switch (m_position) {
    case Top:
        p += QPoint(r.width() / 2, r.height());
        break;
    case Bottom:
        p += QPoint(r.width() / 2, 0);
        break;
    case Left:
        p += QPoint(r.width(), r.height() / 2);
        break;
    case Right:
        p += QPoint(0, r.height() / 2);
        break;
    }
    return p;
}

void QuickDockItem::onMenuActionClicked(QAction *action)
{
    m_pluginItem->invokedMenuItem(m_itemKey, action->data().toString(), true);
}
