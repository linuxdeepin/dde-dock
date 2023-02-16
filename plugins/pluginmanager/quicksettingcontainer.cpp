// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "quicksettingcontainer.h"
#include "dockplugincontroller.h"
#include "pluginsiteminterface.h"
#include "quicksettingitem.h"
#include "pluginchildpage.h"
#include "utils.h"
#include "quickdragcore.h"

#include <DListView>
#include <DStyle>

#include <QDrag>
#include <QVBoxLayout>
#include <QMetaObject>
#include <QStackedLayout>
#include <QMouseEvent>
#include <QLabel>
#include <QBitmap>
#include <QPainterPath>

DWIDGET_USE_NAMESPACE

struct QuickDragInfo {
    QPoint dragPosition;
    QWidget *dragItem = nullptr;
    PluginsItemInterface *pluginInter = nullptr;
    void reset() {
        dragPosition.setX(0);
        dragPosition.setY(0);
        dragItem = nullptr;
        pluginInter = nullptr;
    }
    bool isNull() {
        return !dragItem;
    }
} QuickDragInfo;

#define ITEMWIDTH 70
#define ITEMHEIGHT 60
#define ITEMSPACE 10
#define COLUMNCOUNT 4

QuickSettingContainer::QuickSettingContainer(DockPluginController *pluginController, QWidget *parent)
    : QWidget(parent)
    , m_switchLayout(new QStackedLayout(this))
    , m_mainWidget(new QWidget(this))
    , m_pluginWidget(new QWidget(m_mainWidget))
    , m_pluginLayout(new QGridLayout(m_pluginWidget))
    , m_componentWidget(new QWidget(m_mainWidget))
    , m_mainlayout(new QVBoxLayout(m_mainWidget))
    , m_pluginController(pluginController)
    , m_childPage(new PluginChildPage(this))
    , m_dragInfo(new struct QuickDragInfo)
    , m_childShowPlugin(nullptr)
{
    initUi();
    initConnection();
    m_childPage->installEventFilter(this);
    setMouseTracking(true);
}

QuickSettingContainer::~QuickSettingContainer()
{
    delete m_dragInfo;
}

void QuickSettingContainer::showPage(QWidget *widget, PluginsItemInterface *pluginInter)
{
    if (widget && pluginInter && widget != m_mainWidget) {
        m_childShowPlugin = pluginInter;
        m_childPage->setTitle(pluginInter->pluginDisplayName());
        m_childPage->pushWidget(widget);
        m_switchLayout->setCurrentWidget(m_childPage);
    } else {
        m_childShowPlugin = nullptr;
        m_switchLayout->setCurrentIndex(0);
    }

    onResizeView();
}

bool QuickSettingContainer::eventFilter(QObject *watched, QEvent *event)
{
    switch (event->type()) {
    case QEvent::Resize: {
        onResizeView();
        break;
    }
    case QEvent::MouseButtonPress: {
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        QuickSettingItem *item = qobject_cast<QuickSettingItem *>(watched);
        if (item && item->pluginItem()->flags() & PluginFlag::Attribute_CanDrag) {
            m_dragInfo->dragPosition = mouseEvent->pos();
            m_dragInfo->dragItem = item;
            m_dragInfo->pluginInter = item->pluginItem();
        }
        break;
    }
    case QEvent::MouseButtonRelease: {
        m_dragInfo->reset();
        break;
    }
    default:
        break;
    }

    return QWidget::eventFilter(watched, event);
}

void QuickSettingContainer::showEvent(QShowEvent *event)
{
    // 当面板显示的时候，直接默认显示快捷面板的窗口
    QWidget::showEvent(event);
    if (m_switchLayout->currentWidget() != m_mainWidget) {
        m_childPage->pushWidget(nullptr);
        m_switchLayout->setCurrentWidget(m_mainWidget);
        onResizeView();
    }
}

void QuickSettingContainer::appendPlugin(PluginsItemInterface *itemInter, QString itemKey, bool needLayout)
{
    QuickSettingItem *quickItem = QuickSettingFactory::createQuickWidget(itemInter, itemKey);
    if (!quickItem)
        return;

    quickItem->setParent(m_pluginWidget);
    quickItem->setMouseTracking(true);
    quickItem->installEventFilter(this);
    connect(quickItem, &QuickSettingItem::requestShowChildWidget, this, &QuickSettingContainer::onShowChildWidget);
    m_quickSettings << quickItem;
    if (quickItem->type() == QuickSettingItem::QuickItemStyle::Line) {
        // 插件位置占据整行，例如声音、亮度和音乐等
        m_componentWidget->layout()->addWidget(quickItem);
        updateFullItemLayout();
    } else if (needLayout) {
        // 插件占据两行或者一行
        updateItemLayout();
    }

    onResizeView();
}

void QuickSettingContainer::onPluginRemove(PluginsItemInterface *itemInter)
{
    QList<QuickSettingItem *>::Iterator removeItemIter = std::find_if(m_quickSettings.begin(), m_quickSettings.end(), [ = ](QuickSettingItem *item)->bool {
        return item->pluginItem() == itemInter;
    });

    if (removeItemIter == m_quickSettings.end())
        return;

    QuickSettingItem *removeItem = *removeItemIter;
    removeItem->detachPlugin();

    if (removeItem->type() == QuickSettingItem::QuickItemStyle::Line)
        m_componentWidget->layout()->removeWidget(removeItem);
    else
        m_pluginLayout->removeWidget(removeItem);

    m_quickSettings.removeOne(removeItem);
    removeItem->deleteLater();
    if (m_childShowPlugin == itemInter)
        showPage(nullptr);

    updateItemLayout();
    updateFullItemLayout();
    onResizeView();
}

void QuickSettingContainer::onShowChildWidget(QWidget *childWidget)
{
    QuickSettingItem *quickWidget = qobject_cast<QuickSettingItem *>(sender());
    if (!quickWidget)
        return;

    showPage(childWidget, quickWidget->pluginItem());
}

void QuickSettingContainer::mouseMoveEvent(QMouseEvent *event)
{
    if (m_dragInfo->isNull())
        return;

    QPoint pointCurrent = event->pos();
    if (qAbs(m_dragInfo->dragPosition.x() - pointCurrent.x()) > 5
            || qAbs(m_dragInfo->dragPosition.y() - pointCurrent.y()) > 5) {
        QuickSettingItem *moveItem = qobject_cast<QuickSettingItem *>(m_dragInfo->dragItem);
        QuickIconDrag *drag = new QuickIconDrag(this, moveItem->dragPixmap());
        QuickPluginMimeData *mimedata = new QuickPluginMimeData(m_dragInfo->pluginInter, drag);
        drag->setMimeData(mimedata);
        drag->setDragHotPot(m_dragInfo->dragPosition);

        m_dragInfo->reset();
        drag->exec(Qt::CopyAction);
    }
}

void QuickSettingContainer::updateItemLayout()
{
    // 清空之前的控件，重新添加
    while (m_pluginLayout->count() > 0)
        m_pluginLayout->takeAt(0);

    // 将插件按照两列和一列的顺序来进行排序
    QMap<QuickSettingItem::QuickItemStyle, QList<QuickSettingItem *>> quickSettings;
    QMap<QuickSettingItem::QuickItemStyle, QMap<QuickSettingItem *, int>> orderQuickSettings;
    for (QuickSettingItem *item : m_quickSettings) {
        QuickSettingItem::QuickItemStyle type = item->type();
        if (type == QuickSettingItem::QuickItemStyle::Line)
            continue;

        QJsonObject metaData = m_pluginController->metaData(item->pluginItem());
        if (metaData.contains("order"))
            orderQuickSettings[type][item] = metaData.value("order").toInt();
        else
            quickSettings[type] << item;
    }
    // 将需要排序的插件按照顺序插入到原来的数组中
    for (auto itQuick = orderQuickSettings.begin(); itQuick != orderQuickSettings.end(); itQuick++) {
        QuickSettingItem::QuickItemStyle type = itQuick.key();
        QMap<QuickSettingItem *, int> &orderQuicks = itQuick.value();
        for (auto it = orderQuicks.begin(); it != orderQuicks.end(); it++) {
            int index = it.value();
            if (index >= 0 && index < quickSettings[type].size())
                quickSettings[type][index] = it.key();
            else
                quickSettings[type] << it.key();
        }
    }
    auto insertQuickSetting = [ quickSettings, this ](QuickSettingItem::QuickItemStyle type, int &row, int &column) {
        if (!quickSettings.contains(type))
            return;

        int usedColumn = (type == QuickSettingItem::QuickItemStyle::Larger ? 2 : 1);
        QList<QuickSettingItem *> quickPlugins = quickSettings[type];
        for (QuickSettingItem *quickItem : quickPlugins) {
            quickItem->setVisible(true);
            m_pluginLayout->addWidget(quickItem, row, column, 1, usedColumn);
            column += usedColumn;
            if (column >= COLUMNCOUNT) {
                row++;
                column = 0;
            }
        }
    };

    int row = 0;
    int column = 0;
    insertQuickSetting(QuickSettingItem::QuickItemStyle::Larger, row, column);
    insertQuickSetting(QuickSettingItem::QuickItemStyle::Standard, row, column);
}

void QuickSettingContainer::updateFullItemLayout()
{
    while (m_componentWidget->layout()->count() > 0)
        m_componentWidget->layout()->takeAt(0);

    QList<QuickSettingItem *> fullItems;
    QMap<QuickSettingItem *, int> fullItemOrder;
    for (QuickSettingItem *item : m_quickSettings) {
        if (item->type() != QuickSettingItem::QuickItemStyle::Line)
            continue;

        fullItems << item;
        int order = -1;
        QJsonObject metaData = m_pluginController->metaData(item->pluginItem());
        if (metaData.contains("order"))
            order = metaData.value("order").toInt();

        fullItemOrder[item] = order;
    }

    std::sort(fullItems.begin(), fullItems.end(), [ fullItemOrder, fullItems ](QuickSettingItem *item1, QuickSettingItem *item2) {
        int order1 = fullItemOrder.value(item1, -1);
        int order2 = fullItemOrder.value(item2, -1);
        if (order1 == order2) {
            // 如果两个值相等，就根据他们的加载顺序进行排序
            return fullItems.indexOf(item1) < fullItems.indexOf(item2);
        }
        if (order1 == -1)
            return false;
        if (order2 == -1)
            return true;

        return order1 < order2;
    });

    for (QuickSettingItem *item : fullItems) {
        item->setVisible(true);
        m_componentWidget->layout()->addWidget(item);
    }
}

void QuickSettingContainer::initUi()
{
    m_mainlayout->setSpacing(ITEMSPACE);
    m_mainlayout->setContentsMargins(ITEMSPACE, ITEMSPACE, ITEMSPACE, ITEMSPACE);

    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_pluginLayout->setSpacing(ITEMSPACE);
    m_pluginLayout->setAlignment(Qt::AlignLeft);
    for (int i = 0; i < COLUMNCOUNT; i++)
        m_pluginLayout->setColumnMinimumWidth(i, ITEMWIDTH);

    m_pluginWidget->setLayout(m_pluginLayout);
    m_mainlayout->addWidget(m_pluginWidget);

    QVBoxLayout *ctrlLayout = new QVBoxLayout(m_componentWidget);
    ctrlLayout->setContentsMargins(0, 0, 0, 0);
    ctrlLayout->setSpacing(ITEMSPACE);
    ctrlLayout->setDirection(QBoxLayout::BottomToTop);

    m_mainlayout->addWidget(m_componentWidget);
    // 加载所有的可以在快捷面板显示的插件
    QList<PluginsItemInterface *> plugins = m_pluginController->currentPlugins();
    for (PluginsItemInterface *plugin : plugins) {
         appendPlugin(plugin, m_pluginController->itemKey(plugin), false);
    }

    m_switchLayout->addWidget(m_mainWidget);
    m_switchLayout->addWidget(m_childPage);

    setMouseTracking(true);
    setAcceptDrops(true);

    QMetaObject::invokeMethod(this, [ = ] {
        if (plugins.size() > 0) {
            updateItemLayout();
            updateFullItemLayout();
        }
        // 设置当前窗口的大小
        onResizeView();
        setFixedWidth(ITEMWIDTH * 4 + (ITEMSPACE * 5));
    }, Qt::QueuedConnection);
}

void QuickSettingContainer::initConnection()
{
    connect(m_pluginController, &DockPluginController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, QString itemKey) {
        appendPlugin(itemInter, itemKey);
    });
    connect(m_pluginController, &DockPluginController::pluginRemoved, this, &QuickSettingContainer::onPluginRemove);
    connect(m_pluginController, &DockPluginController::pluginUpdated, this, &QuickSettingContainer::onPluginUpdated);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &QuickSettingContainer::onThemeTypeChanged);

    connect(m_childPage, &PluginChildPage::back, this, [ this ] {
        showPage(m_mainWidget);
    });
}

// 调整尺寸
void QuickSettingContainer::onResizeView()
{
    if (m_switchLayout->currentWidget() == m_mainWidget) {
        int selfPluginCount = 0;
        int fullItemHeight = 0;
        int widgetCount = 0;
        for (QuickSettingItem *item : m_quickSettings) {
            item->setFixedHeight(ITEMHEIGHT);
            if (item->type() == QuickSettingItem::QuickItemStyle::Line) {
                fullItemHeight += item->height();
                widgetCount++;
                continue;
            }
            // 如果是置顶的插件，则认为它占用两个普通插件的位置
            int increCount = (item->type() == QuickSettingItem::QuickItemStyle::Larger ? 2 : 1);
            selfPluginCount += increCount;
        }

        int rowCount = selfPluginCount / COLUMNCOUNT;
        if (selfPluginCount % COLUMNCOUNT > 0)
            rowCount++;

        m_pluginWidget->setFixedHeight(ITEMHEIGHT * rowCount + ITEMSPACE * (rowCount - 1));
        m_componentWidget->setFixedHeight(fullItemHeight + (widgetCount - 1) * ITEMSPACE);

        setFixedHeight(ITEMSPACE * 3 + m_pluginWidget->height() + m_componentWidget->height());
    } else if (m_switchLayout->currentWidget() == m_childPage) {
        setFixedHeight(m_childPage->height());
    }
}

void QuickSettingContainer::onPluginUpdated(PluginsItemInterface *itemInter, const DockPart dockPart)
{
    if (dockPart != DockPart::QuickPanel)
        return;

    for (QuickSettingItem *settingItem : m_quickSettings) {
        if (settingItem->pluginItem() != itemInter)
            continue;

        settingItem->doUpdate();
        break;
    }
}

void QuickSettingContainer::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    for (QuickSettingItem *settingItem : m_quickSettings)
        settingItem->doUpdate();
}
