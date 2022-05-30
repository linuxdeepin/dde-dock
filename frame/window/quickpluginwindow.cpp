/*
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

#include <DStyleOption>
#include <DStandardItem>

#include <QDrag>
#include <QScrollBar>
#include <QStringList>
#include <QSize>
#include <QMouseEvent>
#include <QBoxLayout>

#define ITEMSIZE 22
#define ITEMSPACE 6
#define ICONWIDTH 18
#define ICONHEIGHT 16

static QStringList fixedPluginNames{ "network", "sound", "power" };

const int itemDataRole = Dtk::UserRole + 1;
const int itemSortRole = Dtk::UserRole + 2;

QuickPluginWindow::QuickPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QBoxLayout(QBoxLayout::RightToLeft, this))
    , m_position(Dock::Position::Bottom)
{
    initUi();
    initConnection();

    setAcceptDrops(true);
    setMouseTracking(true);
}

QuickPluginWindow::~QuickPluginWindow()
{
}

void QuickPluginWindow::initUi()
{
    setAcceptDrops(true);
    m_mainLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    m_mainLayout->setDirection(QBoxLayout::RightToLeft);
    m_mainLayout->setContentsMargins(ITEMSPACE, 0, ITEMSPACE, 0);
    m_mainLayout->setSpacing(ITEMSPACE);
    const QList<QuickSettingItem *> &items = QuickSettingController::instance()->settingItems();
    for (QuickSettingItem *settingItem : items) {
        const QString pluginName = settingItem->pluginItem()->pluginName();
        if (!fixedPluginNames.contains(pluginName))
            continue;

        addPlugin(settingItem);
    }
}

void QuickPluginWindow::setPositon(Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    QuickSettingContainer::setPosition(position);
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        m_mainLayout->setDirection(QBoxLayout::RightToLeft);
        m_mainLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    } else {
        m_mainLayout->setDirection(QBoxLayout::BottomToTop);
        m_mainLayout->setAlignment(Qt::AlignTop | Qt::AlignHCenter);
    }
}

int QuickPluginWindow::findActiveTargetIndex(QuickDockItem *widget)
{
    for (int i = 0; i < m_activeSettingItems.size(); i++) {
        PluginsItemInterface *settingItem = m_activeSettingItems[i];
        if (widget->pluginItem() == settingItem)
            return i;
    }

    return -1;
}

void QuickPluginWindow::dragPlugin(PluginsItemInterface *item)
{
    // 释放插件，一般是从快捷设置面板区域移动到这里的，固定插件不支持拖动
    if (fixedPluginNames.contains(item->pluginName()))
        return;

    QPoint itemPoint = mapFromGlobal(QCursor::pos());
    // 查找移动后的位置，如果移动后的插件找不到，就直接放到最后
    QuickDockItem *targetWidget = qobject_cast<QuickDockItem *>(childAt(itemPoint));
    if (!targetWidget) {
        m_activeSettingItems << item;
    } else {
        // 如果是拖动到固定插件区域，也放到最后
        int targetIndex = findActiveTargetIndex(targetWidget);
        if (targetIndex < 0)
            m_activeSettingItems << item;
        else
            m_activeSettingItems.insert(targetIndex, item);
    }

    //排序插入到当前窗体
    resetPluginDisplay();
    Q_EMIT itemCountChanged();
}

void QuickPluginWindow::addPlugin(QuickSettingItem *item)
{
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (!layoutItem)
            continue;

        QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
        if (!dockItem)
            continue;

        if (item->pluginItem() == dockItem->pluginItem()) {
            resetPluginDisplay();
            return;
        }
    }
    if (fixedPluginNames.contains(item->pluginItem()->pluginName())) {
        // 新插入的插件如果是固定插件,则将其插入到固定插件列表中，并对其进行排序
        m_fixedSettingItems << item->pluginItem();
        qSort(m_fixedSettingItems.begin(), m_fixedSettingItems.end(), [](PluginsItemInterface *item1, PluginsItemInterface *item2) {
            int index1 = fixedPluginNames.indexOf(item1->pluginName());
            int index2 = fixedPluginNames.indexOf(item2->pluginName());
            return index1 < index2;
        });
    } else {
        // 如果是非固定插件，则直接插入到末尾
        m_activeSettingItems << item->pluginItem();
    }
    resetPluginDisplay();
    Q_EMIT itemCountChanged();
}

QSize QuickPluginWindow::suitableSize()
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return QSize((ITEMSPACE + ICONWIDTH) * m_mainLayout->count() + ITEMSPACE, ITEMSIZE);

    int height = 0;
    int itemCount = m_mainLayout->count();
    if (itemCount > 0) {
        // 每个图标占据的高度
        height += ICONHEIGHT * itemCount;
        // 图标间距占据的高度
        height += ITEMSPACE * itemCount;
    }

    return QSize(ITEMSIZE, height);
}

void QuickPluginWindow::removePlugin(PluginsItemInterface *item)
{
    if (m_fixedSettingItems.contains(item))
        m_fixedSettingItems.removeOne(item);
    else if (m_activeSettingItems.contains(item))
        m_activeSettingItems.removeOne(item);
    else
        return;

    resetPluginDisplay();
    Q_EMIT itemCountChanged();
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

void QuickPluginWindow::mousePressEvent(QMouseEvent *event)
{
    // 查找非固定的图标，然后执行拖动
    PluginsItemInterface *quickItem = findQuickSettingItem(event->pos(), m_activeSettingItems);
    if (!quickItem)
        return;

    // 如果不是固定图标，则让其拖动
    startDrag(quickItem);
}

QPoint QuickPluginWindow::popupPoint() const
{
    if (!parentWidget())
        return pos();

    QPoint pointCurrent = parentWidget()->mapToGlobal(pos());
    switch (m_position) {
    case Dock::Position::Bottom: {
        // 在下方的时候，Y坐标设置在顶层窗口的y值，保证下方对齐
        pointCurrent.setY(topLevelWidget()->y());
        break;
    }
    case Dock::Position::Top: {
        // 在上面的时候，Y坐标设置为任务栏的下方，保证上方对齐
        pointCurrent.setY(topLevelWidget()->y() + topLevelWidget()->height());
        break;
    }
    case Dock::Position::Left: {
        // 在左边的时候，X坐标设置在顶层窗口的最右侧，保证左对齐
        pointCurrent.setX(topLevelWidget()->x() + topLevelWidget()->width());
        break;
    }
    case Dock::Position::Right: {
        // 在右边的时候，X坐标设置在顶层窗口的最左侧，保证右对齐
        pointCurrent.setX(topLevelWidget()->x());
    }
    }
    return pointCurrent;
}

void QuickPluginWindow::onFixedClick()
{
    // 查找固定团图标，然后点击弹出快捷面板
    QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(sender());
    if (!dockItem || !fixedPluginNames.contains(dockItem->pluginItem()->pluginName()))
        return;

    // 弹出快捷设置面板
    DockPopupWindow *popWindow = QuickSettingContainer::popWindow();
    popWindow->show(popupPoint());
}

void QuickPluginWindow::startDrag(PluginsItemInterface *moveItem)
{
    AppDrag *drag = new AppDrag(this, new QuickDragWidget);
    QuickPluginMimeData *mimedata = new QuickPluginMimeData(moveItem);
    drag->setMimeData(mimedata);
    drag->appDragWidget()->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));
    QPixmap dragPixmap = moveItem->icon(DockPart::QuickPanel).pixmap(QSize(ITEMSIZE, ITEMSIZE));
    drag->setPixmap(dragPixmap);
    drag->setHotSpot(QPoint(0, 0));

    connect(drag->appDragWidget(), &AppDragWidget::requestRemoveItem, this, [ this, moveItem ] {
        removePlugin(moveItem);
    });

    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDropItem, this, &QuickPluginWindow::onPluginDropItem);
    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDragMove, this, &QuickPluginWindow::onPluginDragMove);

    drag->exec(Qt::MoveAction | Qt::CopyAction);
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

int QuickPluginWindow::getDropIndex(QPoint point)
{
    QuickDockItem *targetItem = getDockItemByPlugin(findQuickSettingItem(point, m_activeSettingItems));
    if (targetItem)
        return m_activeSettingItems.indexOf(targetItem->pluginItem());

    // 上下方向从右向左排列
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        for (int i = 0; i < m_mainLayout->count() - 1; i++) {
            QLayoutItem *layoutBefore = m_mainLayout->itemAt(i);
            QLayoutItem *layoutItem = m_mainLayout->itemAt(i + 1);
            if (!layoutBefore || !layoutItem)
                continue;

            QuickDockItem *dockBeforeItem = qobject_cast<QuickDockItem *>(layoutBefore->widget());
            QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
            if (dockItem->pluginItem()->isPrimary())
                continue;

            if (dockBeforeItem->geometry().x() > point.x() && dockItem->geometry().right() < point.x())
                return i;
        }
    }
    for (int i = 0; i < m_mainLayout->count() - 1; i++) {
        QLayoutItem *layoutBefore = m_mainLayout->itemAt(i);
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i + 1);
        if (!layoutBefore || !layoutItem)
            continue;

        QuickDockItem *dockBeforeItem = qobject_cast<QuickDockItem *>(layoutBefore->widget());
        if (dockBeforeItem->pluginItem()->isPrimary())
            break;

        QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());

        // 从上向下排列
        if (dockBeforeItem->geometry().bottom() < point.y() && dockItem->geometry().top() > point.y())
            return i;
    }
    // 如果都没有找到，直接插入到最后
    return -1;
}

void QuickPluginWindow::onPluginDropItem(QDropEvent *event)
{
    const QuickPluginMimeData *data = qobject_cast<const QuickPluginMimeData *>(event->mimeData());
    if (!data)
        return;

    // 获取当前鼠标在任务栏快捷图标区域的位置
    QPoint currentPoint = mapFromGlobal(QCursor::pos());
    // 获取区域图标插入的位置
    int dropIndex = getDropIndex(currentPoint);
    if (dropIndex >= 0)
        m_activeSettingItems.insert(dropIndex, data->pluginItemInterface());
    else
        m_activeSettingItems << data->pluginItemInterface();

    resetPluginDisplay();
    Q_EMIT itemCountChanged();
}

void QuickPluginWindow::onPluginDragMove(QDragMoveEvent *event)
{
    QPoint currentPoint = mapFromGlobal(QCursor::pos());
    const QuickPluginMimeData *data = qobject_cast<const QuickPluginMimeData *>(event->mimeData());
    if (!data)
        return;

    // 查找移动的
    PluginsItemInterface *sourceItem = data->pluginItemInterface();
    if (!sourceItem)
        return;

    QuickDockItem *sourceMoveWidget = getDockItemByPlugin(sourceItem);
    QuickDockItem *targetItem = getDockItemByPlugin(findQuickSettingItem(currentPoint, m_activeSettingItems));
    // 如果未找到要移动的目标位置，或者移动的目标位置是固定插件，或者原插件和目标插件是同一个插件，则不做任何操作
    if (!sourceMoveWidget || !targetItem || sourceMoveWidget == targetItem)
        return;

    // 重新对所有的插件进行排序
    QMap<QWidget *, int> allItems;
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QWidget *childWidget = m_mainLayout->itemAt(i)->widget();
        allItems[childWidget] = i;
    }
    // 调整列表中的位置
    int sourceIndex = m_activeSettingItems.indexOf(sourceItem);
    int targetIndex = m_activeSettingItems.indexOf(targetItem->pluginItem());
    if (sourceIndex >= 0)
        m_activeSettingItems.move(sourceIndex, targetIndex);
    else
        m_activeSettingItems.insert(targetIndex, sourceItem);

    event->accept();
}

void QuickPluginWindow::resetPluginDisplay()
{
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
        pluginItems[dockItem->pluginItem()] = dockItem;
    }
    // 将列表中所有的控件按照顺序添加到布局上
    auto addWidget = [ = ](const QList<PluginsItemInterface *> &items) {
        for (PluginsItemInterface *item : items) {
            QuickDockItem *itemWidget = nullptr;
            if (pluginItems.contains(item)) {
                itemWidget = pluginItems[item];
            } else {
                itemWidget = new QuickDockItem(item, this);
                itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
            }
            connect(itemWidget, &QuickDockItem::clicked, this, &QuickPluginWindow::onFixedClick);
            itemWidget->setParent(this);
            m_mainLayout->addWidget(itemWidget);
        }
    };

    addWidget(m_fixedSettingItems);
    addWidget(m_activeSettingItems);
}

void QuickPluginWindow::initConnection()
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ this ](QuickSettingItem * settingItem) {
        const QString pluginName = settingItem->pluginItem()->pluginName();
        if (!fixedPluginNames.contains(pluginName))
            return;

        addPlugin(settingItem);
    });

    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, [ this ] (QuickSettingItem *settingItem){
        removePlugin(settingItem->pluginItem());
    });
}

/**
 * @brief QuickDockItem::QuickDockItem
 * @param pluginItem
 * @param parent
 */
QuickDockItem::QuickDockItem(PluginsItemInterface *pluginItem, QWidget *parent)
    : QWidget(parent)
    , m_pluginItem(pluginItem)
{
}

QuickDockItem::~QuickDockItem()
{
}

PluginsItemInterface *QuickDockItem::pluginItem()
{
    return m_pluginItem;
}

void QuickDockItem::paintEvent(QPaintEvent *event)
{
    if (!m_pluginItem)
        return QWidget::paintEvent(event);

    QPixmap pixmap = m_pluginItem->icon(DockPart::QuickPanel).pixmap(ICONHEIGHT, ICONHEIGHT);
    QRect pixmapRect = QRect((rect().width() - ICONHEIGHT) / 2, (rect().height() - ICONHEIGHT) / 2,
                             ICONHEIGHT, ICONHEIGHT);
    QPainter painter(this);
    painter.drawPixmap(pixmapRect, pixmap);
}

void QuickDockItem::mouseReleaseEvent(QMouseEvent *event)
{
    Q_EMIT clicked();
    QWidget::mouseReleaseEvent(event);
}
