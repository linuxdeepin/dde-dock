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
#include "proxyplugincontroller.h"
#include "quickpluginmodel.h"

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

        return (qAbs(currentPoint.x() - dragPoint.x()) >= 5 ||
                qAbs(currentPoint.y() - dragPoint.y()) >= 5);
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

static QStringList fixedPluginNames{ "network", "sound", "power" };

QuickPluginWindow::QuickPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QBoxLayout(QBoxLayout::RightToLeft, this))
    , m_position(Dock::Position::Bottom)
    , m_dragInfo(new DragInfo)
{
    initUi();
    initConnection();

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
    if (position == Dock::Position::Top || position == Dock::Position::Bottom)
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

            showPopup(qobject_cast<QuickDockItem *>(watched));
        } while (false);
        m_dragInfo->reset();

        break;
    }
    case QEvent::MouseMove: {
        if (m_dragInfo->isNull())
            break;

        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        if (m_dragInfo->canDrag(mouseEvent->pos())) {
            startDrag();
            m_dragInfo->reset();
        }
        break;
    }
    case QEvent::Drop: {
        Q_EMIT requestDrop(static_cast<QDropEvent *>(event));
        break;
    }
    default:
        break;
    }
    return QWidget::eventFilter(watched, event);
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
            QJsonObject metaData;
            QPluginLoader *pluginLoader = ProxyPluginController::instance(PluginType::QuickPlugin)->pluginLoader(item);
            if (pluginLoader)
                metaData = pluginLoader->metaData().value("MetaData").toObject();

            itemWidget = new QuickDockItem(item, metaData, quickController->itemKey(item), this);
            itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
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
    if (!widget)
        return pos();

    QPoint pointCurrent = widget->mapToGlobal(QPoint(0, 0));
    switch (m_position) {
    case Dock::Position::Bottom: {
        // 在下方的时候，Y坐标设置在顶层窗口的y值，保证下方对齐
        pointCurrent.setX(pointCurrent.x() + widget->width() / 2);
        pointCurrent.setY(topLevelWidget()->y());
        break;
    }
    case Dock::Position::Top: {
        // 在上面的时候，Y坐标设置为任务栏的下方，保证上方对齐
        pointCurrent.setX(pointCurrent.x() + widget->width() / 2);
        pointCurrent.setY(topLevelWidget()->y() + topLevelWidget()->height());
        break;
    }
    case Dock::Position::Left: {
        // 在左边的时候，X坐标设置在顶层窗口的最右侧，保证左对齐
        pointCurrent.setX(topLevelWidget()->x() + topLevelWidget()->width());
        pointCurrent.setY(pointCurrent.y() + widget->height() / 2);
        break;
    }
    case Dock::Position::Right: {
        // 在右边的时候，X坐标设置在顶层窗口的最左侧，保证右对齐
        pointCurrent.setX(topLevelWidget()->x());
        pointCurrent.setY(pointCurrent.y() + widget->height() / 2);
    }
    }
    return pointCurrent;
}

void QuickPluginWindow::onUpdatePlugin(PluginsItemInterface *itemInter, const DockPart &dockPart)
{
    //update plugin status
    if (dockPart != DockPart::QuickShow)
        return;

    QuickDockItem *dockItem = getDockItemByPlugin(itemInter);
    if (dockItem)
        dockItem->update();
}

void QuickPluginWindow::onRequestAppletShow(PluginsItemInterface *itemInter, const QString &itemKey)
{
    showPopup(getDockItemByPlugin(itemInter), itemInter->itemPopupApplet(itemKey));
}

void QuickPluginWindow::startDrag()
{
    if (!m_dragInfo->dockItem)
        return;

    PluginsItemInterface *moveItem = m_dragInfo->dockItem->pluginItem();
    AppDrag *drag = new AppDrag(this, new QuickDragWidget);
    QuickPluginMimeData *mimedata = new QuickPluginMimeData(moveItem);
    drag->setMimeData(mimedata);
    drag->appDragWidget()->setDockInfo(m_position, QRect(mapToGlobal(pos()), size()));
    QPixmap dragPixmap = m_dragInfo->dragPixmap();
    drag->setPixmap(dragPixmap);

    drag->setHotSpot(QPoint(0, 0));

    connect(drag->appDragWidget(), &AppDragWidget::requestSplitWindow, this, [ this, moveItem ] {
        QuickPluginModel::instance()->removePlugin(moveItem);
        Q_EMIT itemCountChanged();
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

void QuickPluginWindow::showPopup(QuickDockItem *item, QWidget *childPage)
{
    if (!isVisible() || !item)
        return;

    bool canBack = true;
    DockPopupWindow *popWindow = QuickSettingContainer::popWindow();
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

        popWindow->show(popupPoint(item), true);
        canBack = false;
    }

    QuickSettingContainer *container = static_cast<QuickSettingContainer *>(popWindow->getContent());
    container->showPage(childPage, item->pluginItem(), canBack);
}

int QuickPluginWindow::getDropIndex(QPoint point)
{
    QuickDockItem *targetItem = getActiveDockItem(point);
    if (targetItem) {
        for (int i = 0; i < m_mainLayout->count(); i++) {
            QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
            if (!layoutItem)
                continue;

            if (layoutItem->widget() == targetItem)
                return i;
        }

        return -1;
    }

    // 上下方向从右向左排列
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        for (int i = 0; i < m_mainLayout->count() - 1; i++) {
            QLayoutItem *layoutBefore = m_mainLayout->itemAt(i);
            QLayoutItem *layoutItem = m_mainLayout->itemAt(i + 1);
            if (!layoutBefore || !layoutItem)
                continue;

            QuickDockItem *dockBeforeItem = qobject_cast<QuickDockItem *>(layoutBefore->widget());
            QuickDockItem *dockItem = qobject_cast<QuickDockItem *>(layoutItem->widget());
            if (dockItem->isPrimary())
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
        if (dockBeforeItem->isPrimary())
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
    QuickPluginModel::instance()->addPlugin(data->pluginItemInterface(), getDropIndex(currentPoint));
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
    QuickDockItem *targetItem = getActiveDockItem(currentPoint);
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
/*    int sourceIndex = m_activeSettingItems.indexOf(sourceItem);
    int targetIndex = m_activeSettingItems.indexOf(targetItem->pluginItem());
    if (sourceIndex >= 0)
        m_activeSettingItems.move(sourceIndex, targetIndex);
    else
        m_activeSettingItems.insert(targetIndex, sourceItem);
*/
    event->accept();
}

void QuickPluginWindow::initConnection()
{
    QuickPluginModel *model = QuickPluginModel::instance();
    connect(model, &QuickPluginModel::requestUpdate, this, &QuickPluginWindow::onRequestUpdate);
    connect(model, &QuickPluginModel::requestUpdatePlugin, this, &QuickPluginWindow::onUpdatePlugin);
    connect(QuickSettingController::instance(), &QuickSettingController::requestAppletShow, this, &QuickPluginWindow::onRequestAppletShow);
}

/**
 * @brief QuickDockItem::QuickDockItem
 * @param pluginItem
 * @param parent
 */
QuickDockItem::QuickDockItem(PluginsItemInterface *pluginItem, const QJsonObject &metaData, const QString itemKey, QWidget *parent)
    : QWidget(parent)
    , m_pluginItem(pluginItem)
    , m_metaData(metaData)
    , m_itemKey(itemKey)
    , m_position(Dock::Position::Bottom)
    , m_popupWindow(new DockPopupWindow)
    , m_contextMenu(new QMenu(this))
    , m_tipParent(nullptr)
    , m_mainLayout(nullptr)
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

PluginsItemInterface *QuickDockItem::pluginItem()
{
    return m_pluginItem;
}

bool QuickDockItem::isPrimary() const
{
    if (m_metaData.contains("primary"))
        return m_metaData.value("primary").toBool();

    return false;
}

void QuickDockItem::hideToolTip()
{
    m_popupWindow->hide();
}

void QuickDockItem::paintEvent(QPaintEvent *event)
{
    if (!m_pluginItem)
        return QWidget::paintEvent(event);

    QPixmap pixmap = iconPixmap();
    QRect pixmapRect = QRect((rect().width() - ICONHEIGHT) / 2, (rect().height() - ICONHEIGHT) / 2,
                             ICONHEIGHT, ICONHEIGHT);

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
        itemWidget->setFixedSize(ICONWIDTH - 2, ICONHEIGHT - 2);
        m_mainLayout->addWidget(itemWidget);
    }
}

void QuickDockItem::hideEvent(QHideEvent *event)
{
    if (!m_mainLayout)
        return QWidget::hideEvent(event);

    QWidget *itemWidget = m_pluginItem->itemWidget(m_itemKey);
    if (itemWidget && m_mainLayout->indexOf(itemWidget) >= 0)
        m_mainLayout->removeWidget(m_pluginItem->itemWidget(m_itemKey));
}

QPixmap QuickDockItem::iconPixmap() const
{
    int pixmapSize = static_cast<int>(ICONHEIGHT * qApp->devicePixelRatio());
    QIcon icon = m_pluginItem->icon(DockPart::QuickShow, DGuiApplicationHelper::instance()->themeType());
    if (!icon.isNull())
        return icon.pixmap(pixmapSize, pixmapSize);

    return QPixmap();
}

void QuickDockItem::initUi()
{
    QPixmap pixmap = iconPixmap();
    if (!pixmap.isNull())
        return;

    m_mainLayout = new QHBoxLayout(this);
    m_mainLayout->setContentsMargins(0, 0, 0, 0);
    m_pluginItem->itemWidget(m_itemKey)->installEventFilter(this);
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
