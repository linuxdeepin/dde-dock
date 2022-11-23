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
#include "quicksettingcontainer.h"
#include "brightnessmodel.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"
#include "quicksettingitem.h"
#include "mediawidget.h"
#include "dockpopupwindow.h"
#include "brightnesswidget.h"
#include "slidercontainer.h"
#include "pluginchildpage.h"
#include "utils.h"
#include "displaysettingwidget.h"

#include <DListView>
#include <DStyle>
#include <QDrag>

#include <QVBoxLayout>
#include <QMetaObject>
#include <QStackedLayout>
#include <QMouseEvent>

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

DockPopupWindow *QuickSettingContainer::m_popWindow = nullptr;
Dock::Position QuickSettingContainer::m_position = Dock::Position::Bottom;

QuickSettingContainer::QuickSettingContainer(QWidget *parent)
    : QWidget(parent)
    , m_switchLayout(new QStackedLayout(this))
    , m_mainWidget(new QWidget(this))
    , m_pluginWidget(new QWidget(m_mainWidget))
    , m_pluginLayout(new QGridLayout(m_pluginWidget))
    , m_componentWidget(new QWidget(m_mainWidget))
    , m_mainlayout(new QVBoxLayout(m_mainWidget))
    , m_pluginLoader(QuickSettingController::instance())
    , m_playerWidget(new MediaWidget(m_componentWidget))
    , m_brightnessModel(new BrightnessModel(this))
    , m_brihtnessWidget(new BrightnessWidget(m_brightnessModel, m_componentWidget))
    , m_displaySettingWidget(new DisplaySettingWidget(this))
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

void QuickSettingContainer::showHomePage()
{
    m_childShowPlugin = nullptr;
    m_switchLayout->setCurrentIndex(0);
}

// 根据位置获取箭头的方向
static DArrowRectangle::ArrowDirection getDirection(const Dock::Position &position)
{
    switch (position) {
    case Dock::Position::Top:
        return DArrowRectangle::ArrowDirection::ArrowTop;
    case Dock::Position::Left:
        return DArrowRectangle::ArrowDirection::ArrowLeft;
    case Dock::Position::Right:
        return DArrowRectangle::ArrowDirection::ArrowRight;
    default:
        return DArrowRectangle::ArrowDirection::ArrowBottom;
    }

    return DArrowRectangle::ArrowDirection::ArrowBottom;
}

DockPopupWindow *QuickSettingContainer::popWindow()
{
    if (m_popWindow) {
        QuickSettingContainer *container = static_cast<QuickSettingContainer *>(m_popWindow->getContent());
        container->showHomePage();
        return m_popWindow;
    }

    m_popWindow = new DockPopupWindow;
    m_popWindow->setWindowFlag(Qt::Popup);
    m_popWindow->setShadowBlurRadius(20);
    m_popWindow->setRadius(18);
    m_popWindow->setShadowYOffset(2);
    m_popWindow->setShadowXOffset(0);
    m_popWindow->setArrowWidth(18);
    m_popWindow->setArrowHeight(10);
    m_popWindow->setArrowDirection(getDirection(m_position));
    m_popWindow->setContent(new QuickSettingContainer(m_popWindow));
    if (Utils::IS_WAYLAND_DISPLAY)
        m_popWindow->setWindowFlags(m_popWindow->windowFlags() | Qt::FramelessWindowHint | Qt::Popup);
    return m_popWindow;
}

void QuickSettingContainer::setPosition(Position position)
{
    if (m_position == position)
        return;

    m_position = position;

    if (m_popWindow) {
        m_popWindow->setArrowDirection(getDirection(m_position));
        // 在任务栏位置发生变化的时候，需要将当前的content获取后，重新调用setContent接口
        // 如果不调用，那么就会出现内容在容器内部的位置错误，界面上的布局会乱
        QWidget *widget = m_popWindow->getContent();
        m_popWindow->setContent(widget);
    }
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
        if (item) {
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

void QuickSettingContainer::showWidget(QWidget *widget, const QString &title)
{
    m_childPage->setTitle(title);
    m_childPage->pushWidget(widget);
    m_switchLayout->setCurrentWidget(m_childPage);
}

QPoint QuickSettingContainer::hotSpot(const QPixmap &pixmap)
{
    if (m_position == Dock::Position::Left)
        return QPoint(0, pixmap.height());

    if (m_position == Dock::Position::Top)
        return QPoint(pixmap.width(), 0);

    return QPoint(pixmap.width(), pixmap.height());
}

void QuickSettingContainer::appendPlugin(PluginsItemInterface *itemInter, bool needLayout)
{
    QuickSettingItem *quickItem = QuickSettingFactory::createQuickWidget(itemInter);
    if (!quickItem)
        return;

    quickItem->setParent(m_pluginWidget);
    quickItem->setMouseTracking(true);
    quickItem->installEventFilter(this);
    connect(quickItem, &QuickSettingItem::requestShowChildWidget, this, &QuickSettingContainer::onShowChildWidget);
    m_quickSettings << quickItem;
    if (quickItem->type() == QuickSettingItem::QuickSettingType::Full) {
        // 插件位置占据整行，例如声音、亮度和音乐等
        m_componentWidget->layout()->addWidget(quickItem);
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

    if (removeItem->type() == QuickSettingItem::QuickSettingType::Full)
        m_componentWidget->layout()->removeWidget(removeItem);
    else
        m_pluginLayout->removeWidget(removeItem);

    m_quickSettings.removeOne(removeItem);
    removeItem->deleteLater();
    if (m_childShowPlugin == itemInter)
        showHomePage();

    updateItemLayout();
    onResizeView();
}

void QuickSettingContainer::onShowChildWidget(QWidget *childWidget)
{
    QuickSettingItem *quickWidget = qobject_cast<QuickSettingItem *>(sender());
    if (!quickWidget)
        return;

    m_childShowPlugin = quickWidget->pluginItem();
    showWidget(childWidget, m_childShowPlugin->pluginDisplayName());
    onResizeView();
}

void QuickSettingContainer::mouseMoveEvent(QMouseEvent *event)
{
    if (m_dragInfo->isNull())
        return;

    QPoint pointCurrent = event->pos();
    if (qAbs(m_dragInfo->dragPosition.x() - pointCurrent.x()) > 5
            || qAbs(m_dragInfo->dragPosition.y() - pointCurrent.y()) > 5) {

        QDrag *drag = new QDrag(this);
        QuickSettingItem *moveItem = qobject_cast<QuickSettingItem *>(m_dragInfo->dragItem);
        QuickPluginMimeData *mimedata = new QuickPluginMimeData(m_dragInfo->pluginInter);
        drag->setMimeData(mimedata);
        if (moveItem) {
            QPixmap dragPixmap = moveItem->dragPixmap();
            drag->setPixmap(dragPixmap);
            drag->setHotSpot(hotSpot(dragPixmap));
        } else {
            // 如果拖动的是声音等插件
            QPixmap dragPixmap = m_dragInfo->dragItem->grab();
            drag->setPixmap(dragPixmap);
            drag->setHotSpot(hotSpot(dragPixmap));
        }

        m_dragInfo->reset();
        drag->exec(Qt::MoveAction | Qt::CopyAction);
    }
}

void QuickSettingContainer::updateItemLayout()
{
    // 清空之前的控件，重新添加
    while (m_pluginLayout->count() > 0)
        m_pluginLayout->takeAt(0);

    // 将插件按照两列和一列的顺序来进行排序
    QMap<QuickSettingItem::QuickSettingType, QList<QuickSettingItem *>> quickSettings;
    QMap<QuickSettingItem::QuickSettingType, QMap<QuickSettingItem *, int>> orderQuickSettings;
    QuickSettingController *quickController = QuickSettingController::instance();
    for (QuickSettingItem *item : m_quickSettings) {
        QuickSettingItem::QuickSettingType type = item->type();
        if (type == QuickSettingItem::QuickSettingType::Full)
            continue;

        QJsonObject metaData = quickController->metaData(item->pluginItem());
        if (metaData.contains("order"))
            orderQuickSettings[type][item] = metaData.value("order").toInt();
        else
            quickSettings[type] << item;
    }
    // 将需要排序的插件按照顺序插入到原来的数组中
    for (auto itQuick = orderQuickSettings.begin(); itQuick != orderQuickSettings.end(); itQuick++) {
        QuickSettingItem::QuickSettingType type = itQuick.key();
        QMap<QuickSettingItem *, int> &orderQuicks = itQuick.value();
        for (auto it = orderQuicks.begin(); it != orderQuicks.end(); it++) {
            int index = it.value();
            if (index >= 0 && index < quickSettings[type].size())
                quickSettings[type][index] = it.key();
            else
                quickSettings[type] << it.key();
        }
    }
    auto insertQuickSetting = [ quickSettings, this ](QuickSettingItem::QuickSettingType type, int &row, int &column) {
        if (!quickSettings.contains(type))
            return;

        int usedColumn = (type == QuickSettingItem::QuickSettingType::Multi ? 2 : 1);
        QList<QuickSettingItem *> quickPlugins = quickSettings[type];
        for (QuickSettingItem *quickItem : quickPlugins) {
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
    insertQuickSetting(QuickSettingItem::QuickSettingType::Multi, row, column);
    insertQuickSetting(QuickSettingItem::QuickSettingType::Single, row, column);
}

void QuickSettingContainer::initUi()
{
    auto setWidgetStyle = [](DBlurEffectWidget *widget) {
        widget->setMaskColor(QColor(239, 240, 245));
        widget->setBlurRectXRadius(8);
        widget->setBlurRectYRadius(8);
    };

    // 添加音乐播放插件
    m_playerWidget->setFixedHeight(ITEMHEIGHT);
    m_brihtnessWidget->setFixedHeight(ITEMHEIGHT);

    setWidgetStyle(m_playerWidget);
    setWidgetStyle(m_brihtnessWidget);

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

    ctrlLayout->addWidget(m_playerWidget);
    ctrlLayout->addWidget(m_brihtnessWidget);

    m_mainlayout->addWidget(m_componentWidget);
    // 加载所有的插件
    QList<PluginsItemInterface *> plugins = m_pluginLoader->pluginItems(QuickSettingController::PluginAttribute::Quick);
    for (PluginsItemInterface *plugin : plugins)
        appendPlugin(plugin, false);

    m_switchLayout->addWidget(m_mainWidget);
    m_switchLayout->addWidget(m_childPage);

    setMouseTracking(true);
    setAcceptDrops(true);

    QMetaObject::invokeMethod(this, [ = ] {
        if (plugins.size() > 0)
            updateItemLayout();
        // 设置当前窗口的大小
        onResizeView();
        setFixedWidth(ITEMWIDTH * 4 + (ITEMSPACE * 5));
    }, Qt::QueuedConnection);

    m_displaySettingWidget->setVisible(false);
}

void QuickSettingContainer::initConnection()
{
    connect(m_pluginLoader, &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::Quick)
            return;

        appendPlugin(itemInter);
    });
    connect(m_pluginLoader, &QuickSettingController::pluginRemoved, this, &QuickSettingContainer::onPluginRemove);
    connect(m_pluginLoader, &QuickSettingController::pluginUpdated, this, &QuickSettingContainer::onPluginUpdated);

    connect(m_playerWidget, &MediaWidget::visibleChanged, this, &QuickSettingContainer::onResizeView);
    connect(m_brihtnessWidget, &BrightnessWidget::visibleChanged, this, &QuickSettingContainer::onResizeView);
    connect(m_brihtnessWidget->sliderContainer(), &SliderContainer::iconClicked, this, [ this ](const SliderContainer::IconPosition &iconPosition) {
        if (iconPosition == SliderContainer::RightIcon) {
            // 点击右侧的按钮，弹出具体的调节的界面
            showWidget(m_displaySettingWidget, tr("brightness"));
            onResizeView();
        }
    });
    connect(m_childPage, &PluginChildPage::back, this, [ this ] {
        m_switchLayout->setCurrentWidget(m_mainWidget);
    });
    connect(m_childPage, &PluginChildPage::closeSelf, this, [ this ] {
        if (!m_childPage->isBack())
            topLevelWidget()->hide();
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
            if (item->type() == QuickSettingItem::QuickSettingType::Full) {
                fullItemHeight += item->height();
                widgetCount++;
                continue;
            }
            // 如果是置顶的插件，则认为它占用两个普通插件的位置
            int increCount = (item->type() == QuickSettingItem::QuickSettingType::Multi ? 2 : 1);
            selfPluginCount += increCount;
        }

        int rowCount = selfPluginCount / COLUMNCOUNT;
        if (selfPluginCount % COLUMNCOUNT > 0)
            rowCount++;

        m_pluginWidget->setFixedHeight(ITEMHEIGHT * rowCount + ITEMSPACE * (rowCount - 1));

        if (m_playerWidget->isVisible()) {
            fullItemHeight += m_playerWidget->height();
            widgetCount++;
        }
        if (m_brihtnessWidget->isVisible()) {
            fullItemHeight += m_brihtnessWidget->height();
            widgetCount++;
        }

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

        settingItem->update();
        break;
    }
}
