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
#define ICONSIZE 16

static QStringList fixedPluginKeys{ "network-item-key", "sound-item-key", "power" };
const int itemDataRole = Dtk::UserRole + 1;
const int itemSortRole = Dtk::UserRole + 2;

QuickPluginWindow::QuickPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QBoxLayout(QBoxLayout::LeftToRight, this))
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
    m_mainLayout->setDirection(QBoxLayout::Direction::RightToLeft);
    m_mainLayout->setContentsMargins(ITEMSPACE, 0, ITEMSPACE, 0);
    m_mainLayout->setSpacing(ITEMSPACE);
    const QList<QuickSettingItem *> &items = QuickSettingController::instance()->settingItems();
    for (QuickSettingItem *settingItem : items) {
        const QString itemKey = settingItem->itemKey();
        if (!fixedPluginKeys.contains(itemKey))
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
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        m_mainLayout->setDirection(QBoxLayout::RightToLeft);
    else
        m_mainLayout->setDirection(QBoxLayout::BottomToTop);
}

int QuickPluginWindow::findActiveTargetIndex(QWidget *widget)
{
    for (int i = 0; i < m_activeSettingItems.size(); i++) {
        QuickSettingItem *settingItem = m_activeSettingItems[i];
        if (settingItem->pluginItem()->itemWidget(settingItem->itemKey()) == widget)
            return i;
    }
    return -1;
}

void QuickPluginWindow::dragPlugin(QuickSettingItem *item)
{
    // 释放插件，一般是从快捷设置面板区域移动到这里的，固定插件不支持拖动
    if (fixedPluginKeys.contains(item->itemKey()))
        return;

    QPoint itemPoint = mapFromGlobal(QCursor::pos());
    // 查找移动后的位置，如果移动后的插件找不到，就直接放到最后
    QWidget *targetWidget = childAt(itemPoint);
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
        QWidget *widget = m_mainLayout->itemAt(i)->widget();
        if (item == widget) {
            resetPluginDisplay();
            return;
        }
    }
    QWidget *widget = item->pluginItem()->itemWidget(item->itemKey());
    if (!widget)
        return;

    widget->installEventFilter(this);
    if (fixedPluginKeys.contains(item->itemKey())) {
        // 新插入的插件如果是固定插件,则将其插入到固定插件列表中，并对其进行排序
        m_fixedSettingItems << item;
        qSort(m_fixedSettingItems.begin(), m_fixedSettingItems.end(), [](QuickSettingItem *item1, QuickSettingItem *item2) {
            int index1 = fixedPluginKeys.indexOf(item1->itemKey());
            int index2 = fixedPluginKeys.indexOf(item2->itemKey());
            return index1 < index2;
        });
    } else {
        // 如果是非固定插件，则直接插入到末尾
        m_activeSettingItems << item;
    }
    resetPluginDisplay();
    Q_EMIT itemCountChanged();
}

QSize QuickPluginWindow::suitableSize()
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return QSize((ITEMSPACE + ICONSIZE) * m_mainLayout->count() + ITEMSIZE, ITEMSIZE);

    return QSize(ITEMSIZE, (ITEMSIZE + ITEMSPACE) * m_mainLayout->count() + ITEMSPACE);
}

void QuickPluginWindow::removePlugin(QuickSettingItem *item)
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

QuickSettingItem *QuickPluginWindow::findQuickSettingItem(const QPoint &mousePoint, const QList<QuickSettingItem *> &settingItems)
{
    QWidget *selectWidget = childAt(mousePoint);
    if (!selectWidget)
        return nullptr;

    for (int i = 0; i < settingItems.size(); i++) {
        QuickSettingItem *settingItem = settingItems[i];
        QWidget *widget = settingItem->pluginItem()->itemWidget(settingItem->itemKey());
        if (selectWidget == widget)
            return settingItem;
    }

    return nullptr;
}

void QuickPluginWindow::mousePressEvent(QMouseEvent *event)
{
    // 查找非固定的图标，然后执行拖动
    QuickSettingItem *quickItem = findQuickSettingItem(event->pos(), m_activeSettingItems);
    if (!quickItem)
        return;

    // 如果不是固定图标，则让其拖动
    startDrag(quickItem);
}

void QuickPluginWindow::mouseReleaseEvent(QMouseEvent *event)
{
    // 查找固定团图标，然后点击弹出快捷面板
    QuickSettingItem *quickItem = findQuickSettingItem(event->pos(), m_fixedSettingItems);
    if (!quickItem)
        return;

    // 如果是固定图标，则让其弹出列表
    QPoint currentPoint = pos();
    QWidget *callWidget = parentWidget();
    if (callWidget) {
        currentPoint = callWidget->mapToGlobal(currentPoint);
        // 如果是向上的方向，则需要减去高度
        if (m_position == Dock::Position::Top)
            currentPoint.setY(currentPoint.y() + callWidget->height());
        else if (m_position == Dock::Position::Left)
            currentPoint.setX(currentPoint.x() + callWidget->width());
    }

    DockPopupWindow *popWindow = QuickSettingContainer::popWindow();
    popWindow->show(currentPoint);
}

void QuickPluginWindow::startDrag(QuickSettingItem *moveItem)
{
    AppDrag *drag = new AppDrag(this, new QuickDragWidget);
    QuickPluginMimeData *mimedata = new QuickPluginMimeData(moveItem);
    drag->setMimeData(mimedata);
    drag->appDragWidget()->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));
    QPixmap dragPixmap = moveItem->pluginItem()->icon()->pixmap(QSize(ITEMSIZE, ITEMSIZE));
    drag->setPixmap(dragPixmap);
    drag->setHotSpot(QPoint(0, 0));

    connect(drag->appDragWidget(), &AppDragWidget::requestRemoveItem, this, [ this, moveItem ] {
        removePlugin(moveItem);
    });

    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDropItem, this, [ this] {
        resetPluginDisplay();
        Q_EMIT itemCountChanged();
    });
    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDragMove, this, &QuickPluginWindow::onPluginDragMove);

    drag->exec(Qt::MoveAction | Qt::CopyAction);
}

void QuickPluginWindow::onPluginDragMove(QDragMoveEvent *event)
{
    QPoint currentPoint = mapFromGlobal(QCursor::pos());
    const QuickPluginMimeData *data = qobject_cast<const QuickPluginMimeData *>(event->mimeData());
    if (!data)
        return;

    QuickSettingItem *sourceItem = data->quickSettingItem();
    if (!sourceItem)
        return;

    QWidget *sourceMoveWidget = sourceItem->pluginItem()->itemWidget(sourceItem->itemKey());
    QuickSettingItem *targetItem = findQuickSettingItem(currentPoint, m_activeSettingItems);
    // 如果未找到要移动的目标位置，或者移动的目标位置是固定插件，或者原插件和目标插件是同一个插件，则不做任何操作
    if (!sourceMoveWidget || !targetItem || sourceItem == targetItem)
        return;

    // 重新对所有的插件进行排序
    QMap<QWidget *, int> allItems;
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QWidget *childWidget = m_mainLayout->itemAt(i)->widget();
        allItems[childWidget] = i;
    }
    // 调整列表中的位置
    int sourceIndex = m_activeSettingItems.indexOf(sourceItem);
    int targetIndex = m_activeSettingItems.indexOf(targetItem);
    if (sourceIndex >= 0)
        m_activeSettingItems.move(sourceIndex, targetIndex);
    else
        m_activeSettingItems.insert(targetIndex, sourceItem);

    event->accept();
}

QList<QuickSettingItem *> QuickPluginWindow::settingItems()
{
    QList<QuickSettingItem *> items;
    for (int i = 0; i < m_mainLayout->count(); i++) {
        qInfo() << m_mainLayout->itemAt(i)->widget();
        QuickSettingItem *item = qobject_cast<QuickSettingItem *>(m_mainLayout->itemAt(i)->widget());
        if (item)
            items << item;
    }
    return items;
}

void QuickPluginWindow::resetPluginDisplay()
{
    // 先删除所有的widget
    for (int i = m_mainLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (layoutItem) {
            layoutItem->widget()->setParent(nullptr);
            m_mainLayout->removeItem(layoutItem);
        }
    }
    // 将列表中所有的控件按照顺序添加到布局上
    auto addWidget = [ this ](const QList<QuickSettingItem *> &items) {
        for (QuickSettingItem *item : items) {
            QWidget *itemWidget = item->pluginItem()->itemWidget(item->itemKey());
            itemWidget->setParent(this);
            m_mainLayout->addWidget(itemWidget);
        }
    };

    addWidget(m_activeSettingItems);
    addWidget(m_fixedSettingItems);
}

void QuickPluginWindow::initConnection()
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ this ](QuickSettingItem * settingItem) {
        const QString itemKey = settingItem->itemKey();
        if (!fixedPluginKeys.contains(itemKey))
            return;

        addPlugin(settingItem);
    });

    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, &QuickPluginWindow::removePlugin);
}
