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

#define ITEMSIZE 22
#define ITEMSPACE 6

static QStringList fixedPluginKeys{ "network-item-key", "sound-item-key", "power" };
const int itemDataRole = Dtk::UserRole + 1;
const int itemSortRole = Dtk::UserRole + 2;

QuickPluginWindow::QuickPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_listView(new DListView(this))
    , m_model(new QStandardItemModel(this))
    , m_position(Dock::Position::Bottom)
{
    initUi();
    initConnection();

    setAcceptDrops(true);
    setMouseTracking(true);

    this->installEventFilter(this);
}

QuickPluginWindow::~QuickPluginWindow()
{
}

void QuickPluginWindow::initUi()
{
    m_listView->setModel(m_model);
    m_listView->setViewMode(QListView::IconMode);
    m_listView->setMovement(QListView::Free);
    m_listView->setWordWrap(false);
    m_listView->verticalScrollBar()->setVisible(false);
    m_listView->horizontalScrollBar()->setVisible(false);
    m_listView->setOrientation(QListView::Flow::LeftToRight, false);
    m_listView->setGridSize(QSize(ITEMSIZE + 10, ITEMSIZE + 10));
    m_listView->setSpacing(ITEMSPACE);
    m_listView->setContentsMargins(0,0,0,0);
    m_model->setSortRole(itemSortRole);

    QHBoxLayout *layout = new QHBoxLayout(this);
    layout->setContentsMargins(0,0,0,0);
    layout->setSpacing(0);
    layout->addWidget(m_listView);

    const QList<QuickSettingItem *> &items = QuickSettingController::instance()->settingItems();
    for (QuickSettingItem *settingItem : items) {
        const QString itemKey = settingItem->itemKey();
        if (!fixedPluginKeys.contains(itemKey))
            return;

        addPlugin(settingItem);
    }
}

void QuickPluginWindow::setPositon(Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        m_listView->setOrientation(QListView::Flow::LeftToRight, false);
    else
        m_listView->setOrientation(QListView::Flow::TopToBottom, false);
}

void QuickPluginWindow::addPlugin(QuickSettingItem *item)
{
    // 判断当前的插件是否存在，如果存在，则不插入
    for (int i = 0; i < m_model->rowCount(); i++) {
        QStandardItem *myItem = m_model->item(i, 0);
        QuickSettingItem *settingItem = myItem->data(itemDataRole).value<QuickSettingItem *>();
        if (settingItem == item) {
            m_model->sort(0);
            return;
        }
    }

    DStandardItem *standItem = createStandItem(item);
    if (!standItem)
        return;

    m_model->appendRow(standItem);
    resetSortRole();
    m_model->sort(0);
    Q_EMIT itemCountChanged();
}

QSize QuickPluginWindow::suitableSize()
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return QSize((ITEMSIZE + ITEMSPACE) * m_model->rowCount() + ITEMSPACE, ITEMSIZE);

    return QSize(ITEMSIZE, (ITEMSIZE + ITEMSPACE) * m_model->rowCount() + ITEMSPACE);
}

void QuickPluginWindow::removePlugin(QuickSettingItem *item)
{
    for (int i = 0; i < m_model->rowCount(); i++) {
        QModelIndex index = m_model->index(i, 0);
        QuickSettingItem *quickItem = index.data(itemDataRole).value<QuickSettingItem *>();
        if (quickItem == item) {
            m_model->removeRow(i);
            break;
        }
    }
    m_model->sort(0);

    Q_EMIT itemCountChanged();
}

void QuickPluginWindow::resetSortRole()
{
    QList<QPair<QStandardItem *, int>> fixedItems;
    QList<QStandardItem *> activeItems;
    for (int i = 0; i < m_model->rowCount(); i++) {
        QModelIndex index = m_model->index(i, 0);
        if (!index.data(itemDataRole).canConvert<QuickSettingItem *>())
            continue;

        QuickSettingItem *settingItem = index.data(itemDataRole).value<QuickSettingItem *>();
        if (fixedPluginKeys.contains(settingItem->itemKey()))
            fixedItems.push_back({ m_model->item(i, 0), fixedPluginKeys.indexOf(settingItem->itemKey()) });
        else
            activeItems << m_model->item(i, 0);
    }

    for (int i = 0; i < activeItems.size(); i++) {
        QStandardItem *item = activeItems[i];
        item->setData(i, itemSortRole);
    }

    for (QPair<QStandardItem *, int> item : fixedItems)
        item.first->setData(activeItems.size() + item.second, itemSortRole);
}

DStandardItem *QuickPluginWindow::createStandItem(QuickSettingItem *item)
{
    const QString itemKey = item->itemKey();
    QWidget *itemWidget = item->pluginItem()->itemWidget(itemKey);
    if (!itemWidget)
        return nullptr;

    itemWidget->setParent(m_listView);
    DStandardItem *standItem = new DStandardItem;
    standItem->setFlags(Qt::ItemIsEnabled);
    standItem->setBackground(Qt::transparent);
    standItem->setData(QVariant::fromValue(item), itemDataRole);

    DViewItemAction *action = new DViewItemAction(Qt::AlignCenter, QSize(ITEMSIZE, ITEMSIZE), QSize(ITEMSIZE, ITEMSIZE), true);
    action->setWidget(itemWidget);
    connect(action, &DViewItemAction::triggered, this, [ this ] {
        QPoint ptCurrent = pos();
        QWidget *callWidget = parentWidget();
        if (callWidget)
            ptCurrent = callWidget->mapToGlobal(ptCurrent);

        QuickSettingContainer::popWindow()->show(ptCurrent);
    });
    connect(action, &DViewItemAction::destroyed, this, [ itemWidget ] {
        itemWidget->setParent(nullptr);
        itemWidget->hide();
    });

    standItem->setActionList(Qt::LeftEdge, { action });
    return standItem;
}

void QuickPluginWindow::mouseReleaseEvent(QMouseEvent *event)
{
    QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
    QModelIndex selectedIndex = m_listView->indexAt(mouseEvent->pos());
    if (!selectedIndex.isValid())
        return;

    QuickSettingItem *moveItem = selectedIndex.data(itemDataRole).value<QuickSettingItem *>();
    if (!moveItem)
        return;

    if (fixedPluginKeys.contains(moveItem->itemKey())) {
        QPoint currentPoint = pos();
        QWidget *callWidget = parentWidget();
        if (callWidget)
            currentPoint = callWidget->mapToGlobal(currentPoint);

        QuickSettingContainer::popWindow()->show(currentPoint);
    }
}

void QuickPluginWindow::mousePressEvent(QMouseEvent *event)
{
    QModelIndex selectedIndex = m_listView->indexAt(event->pos());
    if (!selectedIndex.isValid()) {
        QWidget::mousePressEvent(event);
        return;
    }

    QuickSettingItem *moveItem = selectedIndex.data(itemDataRole).value<QuickSettingItem *>();
    if (!moveItem) {
        QWidget::mousePressEvent(event);
        return;
    }

    if (fixedPluginKeys.contains(moveItem->itemKey())) {
        QWidget::mousePressEvent(event);
        return;
    }

    startDrag(moveItem);
}

void QuickPluginWindow::startDrag(QuickSettingItem *moveItem)
{
    AppDrag *drag = new AppDrag(this, new QuickDragWidget);
    CustomMimeData *mimedata = new CustomMimeData;
    mimedata->setData(moveItem);
    drag->setMimeData(mimedata);
    drag->appDragWidget()->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));
    QPixmap dragPixmap = moveItem->pluginItem()->icon()->pixmap(QSize(ITEMSIZE, ITEMSIZE));
    drag->setPixmap(dragPixmap);
    drag->setHotSpot(QPoint(dragPixmap.width() / 2, dragPixmap.height() / 2));

    connect(drag->appDragWidget(), &AppDragWidget::requestRemoveItem, this, [ this, moveItem ] {
        removePlugin(moveItem);
    });

    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDropItem, this, [ this, moveItem ](){
        addPlugin(moveItem);
    });
    connect(static_cast<QuickDragWidget *>(drag->appDragWidget()), &QuickDragWidget::requestDragMove, this, [ this ](QDragMoveEvent *eve){
        QPoint ptCurrent = m_listView->mapFromGlobal(QCursor::pos());
        QModelIndex index = m_listView->indexAt(ptCurrent);
        if (!index.isValid())
            return;

        CustomMimeData *data = const_cast<CustomMimeData *>(qobject_cast<const CustomMimeData *>(eve->mimeData()));
        if (!data)
            return;

        QuickSettingItem *sourceItem = static_cast<QuickSettingItem *>(data->data());
        if (!sourceItem)
            return;

        QuickSettingItem *targetItem = index.data(itemDataRole).value<QuickSettingItem *>();
        if (!targetItem || fixedPluginKeys.contains(targetItem->itemKey()) || sourceItem == targetItem)
            return;

        // recall all sortroles
        QList<QPair<QModelIndex, QuickSettingItem *>> allItems;
        for (int i = 0; i < m_model->rowCount(); i++) {
            QModelIndex rowIndex = m_model->index(i, 0);
            allItems.push_back({ rowIndex, rowIndex.data(itemDataRole).value<QuickSettingItem *>() });
        }
        auto findIndex = [ allItems ](QuickSettingItem *item) {
            for (int i =  0; i < allItems.size(); i++) {
                const QPair<QModelIndex, QuickSettingItem *> &rowItem = allItems[i];
                if (rowItem.second == item)
                    return i;
            }
            return -1;
        };
        int sourceIndex = findIndex(sourceItem);
        int targetIndex = findIndex(targetItem);
        if (sourceIndex < 0 || targetIndex < 0 || sourceIndex == targetIndex)
            return;

        allItems.move(sourceIndex, targetIndex);

        for (int i = 0; i < allItems.size(); i++) {
            const QPair<QModelIndex, QuickSettingItem *> &rowItem = allItems[i];
            m_model->setData(rowItem.first, i, itemSortRole);
        }

        eve->accept();
    });

    drag->exec(Qt::MoveAction | Qt::CopyAction);
}

void QuickPluginWindow::initConnection()
{
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInsert, this, [ this ](QuickSettingItem * settingItem) {
        const QString itemKey = settingItem->itemKey();
        if (!fixedPluginKeys.contains(itemKey))
            return;

        addPlugin(settingItem);
    });

    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemove, this, [ this ](QuickSettingItem *settingItem) {
        removePlugin(settingItem);
    });
}

int QuickPluginWindow::fixedItemCount()
{
    int count = 0;
    for (int i = 0; i < m_model->rowCount(); i++) {
        QModelIndex index = m_model->index(i, 0);
        QuickSettingItem *item = index.data(itemDataRole).value<QuickSettingItem *>();
        if (item && fixedPluginKeys.contains(item->itemKey()))
            count++;
    }

    return count;
}
