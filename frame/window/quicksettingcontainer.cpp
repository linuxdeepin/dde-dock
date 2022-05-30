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
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"
#include "quicksettingitem.h"
#include "mediawidget.h"
#include "dockpopupwindow.h"
#include "brightnesswidget.h"
#include "volumewidget.h"
#include "volumedeviceswidget.h"
#include "brightnessmonitorwidget.h"
#include "pluginchildpage.h"

#include <DListView>
#include <DStyle>
#include <QDrag>

#include <QVBoxLayout>
#include <QMetaObject>
#include <QStackedLayout>

DWIDGET_USE_NAMESPACE

static const int QuickItemRole = Dtk::UserRole + 10;

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
    , m_volumnWidget(new VolumeWidget(m_componentWidget))
    , m_brihtnessWidget(new BrightnessWidget(m_componentWidget))
    , m_volumeSettingWidget(new VolumeDevicesWidget(m_volumnWidget->model(), this))
    , m_brightSettingWidget(new BrightnessMonitorWidget(m_brihtnessWidget->model(), this))
    , m_childPage(new PluginChildPage(this))
    , m_dragPluginPosition(QPoint(0, 0))
{
    initUi();
    initConnection();
    m_childPage->installEventFilter(this);
    setMouseTracking(true);
}

QuickSettingContainer::~QuickSettingContainer()
{
}

void QuickSettingContainer::showHomePage()
{
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

void QuickSettingContainer::initQuickItem(QuickSettingItem *quickItem)
{
    quickItem->setParent(m_pluginWidget);
    quickItem->setMouseTracking(true);
    quickItem->installEventFilter(this);
    connect(quickItem, &QuickSettingItem::detailClicked, this, &QuickSettingContainer::onItemDetailClick);
}

void QuickSettingContainer::onItemDetailClick(PluginsItemInterface *pluginInter)
{
    QuickSettingItem *quickItemWidget = static_cast<QuickSettingItem *>(sender());
    if (!quickItemWidget)
        return;

    QWidget *widget = pluginInter->itemWidget(quickItemWidget->itemKey());
    if (!widget)
        return;

    showWidget(widget, pluginInter->pluginDisplayName());
}

bool QuickSettingContainer::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_childPage && event->type() == QEvent::Resize)
        resizeView();

    return QWidget::eventFilter(watched, event);
}

void QuickSettingContainer::showWidget(QWidget *widget, const QString &title)
{
    m_childPage->setTitle(title);
    m_childPage->pushWidget(widget);
    m_switchLayout->setCurrentWidget(m_childPage);
}

void QuickSettingContainer::onPluginInsert(QuickSettingItem *quickItem)
{
    initQuickItem(quickItem);
    updateItemLayout();
    resizeView();
}

void QuickSettingContainer::onPluginRemove(QuickSettingItem *quickItem)
{
    disconnect(quickItem, &QuickSettingItem::detailClicked, this, &QuickSettingContainer::onItemDetailClick);
    quickItem->setParent(nullptr);
    quickItem->removeEventFilter(this);
    quickItem->setMouseTracking(false);

    //调整子控件的位置
    updateItemLayout();
    resizeView();
}

void QuickSettingContainer::mousePressEvent(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return QWidget::mousePressEvent(event);

    QuickSettingItem *moveItem = qobject_cast<QuickSettingItem *>(childAt(event->pos()));
    if (!moveItem || moveItem->pluginItem()->isPrimary())
        return QWidget::mousePressEvent(event);

    m_dragPluginPosition = event->pos();
}

void QuickSettingContainer::clearDragPoint()
{
    m_dragPluginPosition.setX(0);
    m_dragPluginPosition.setY(0);
}

void QuickSettingContainer::mouseReleaseEvent(QMouseEvent *event)
{
    Q_UNUSED(event);
    clearDragPoint();
}

void QuickSettingContainer::mouseMoveEvent(QMouseEvent *event)
{
    if (m_dragPluginPosition.isNull())
        return;

    QuickSettingItem *moveItem = qobject_cast<QuickSettingItem *>(childAt(m_dragPluginPosition));
    if (!moveItem) {
        clearDragPoint();
        return;
    }

    QPoint pointCurrent = event->pos();
    if (qAbs(m_dragPluginPosition.x() - pointCurrent.x()) > 5
            || qAbs(m_dragPluginPosition.y() - pointCurrent.y()) > 5) {
        clearDragPoint();
        QDrag *drag = new QDrag(this);
        QuickPluginMimeData *mimedata = new QuickPluginMimeData(moveItem->pluginItem());
        drag->setMimeData(mimedata);
        QPixmap dragPixmap = moveItem->dragPixmap();
        drag->setPixmap(dragPixmap);
        drag->setHotSpot(QPoint(dragPixmap.width() / 2, dragPixmap.height() / 2));

        drag->exec(Qt::MoveAction | Qt::CopyAction);
    }
}

void QuickSettingContainer::updateItemLayout()
{
    // 清空之前的控件，重新添加
    while (m_pluginLayout->count() > 0)
        m_pluginLayout->takeAt(0);

    int row = 0;
    int column = 0;
    QList<QuickSettingItem *> quickSettings = m_pluginLoader->settingItems();
    for (QuickSettingItem *item : quickSettings) {
        int usedColumn = item->pluginItem()->isPrimary() ? 2 : 1;
        m_pluginLayout->addWidget(item, row, column, 1, usedColumn);
        column += usedColumn;
        if (column >= COLUMNCOUNT) {
            row++;
            column = 0;
        }
    }
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
    m_volumnWidget->setFixedHeight(ITEMHEIGHT);
    m_brihtnessWidget->setFixedHeight(ITEMHEIGHT);

    setWidgetStyle(m_playerWidget);
    setWidgetStyle(m_volumnWidget);
    setWidgetStyle(m_brihtnessWidget);

    m_mainlayout->setSpacing(ITEMSPACE);
    m_mainlayout->setContentsMargins(ITEMSPACE, ITEMSPACE, ITEMSPACE, ITEMSPACE);

    m_pluginLayout->setContentsMargins(0, 0, 0, 0);
    m_pluginLayout->setSpacing(ITEMSPACE);

    m_pluginWidget->setLayout(m_pluginLayout);
    m_mainlayout->addWidget(m_pluginWidget);

    QVBoxLayout *ctrlLayout = new QVBoxLayout(m_componentWidget);
    ctrlLayout->setContentsMargins(0, 0, 0, 0);
    ctrlLayout->setSpacing(ITEMSPACE);

    ctrlLayout->addWidget(m_playerWidget);
    ctrlLayout->addWidget(m_volumnWidget);
    ctrlLayout->addWidget(m_brihtnessWidget);

    m_mainlayout->addWidget(m_componentWidget);
    // 加载所有的插件
    QList<QuickSettingItem *> pluginItems = m_pluginLoader->settingItems();
    for (QuickSettingItem *quickItem: pluginItems)
        initQuickItem(quickItem);

    m_switchLayout->addWidget(m_mainWidget);
    m_switchLayout->addWidget(m_childPage);

    m_volumeSettingWidget->hide();
    m_brightSettingWidget->hide();

    setMouseTracking(true);
    setAcceptDrops(true);

    QMetaObject::invokeMethod(this, [ = ] {
        if (pluginItems.size() > 0)
            updateItemLayout();
        // 设置当前窗口的大小
        resizeView();
        setFixedWidth(ITEMWIDTH * 4 + (ITEMSPACE * 5));
    }, Qt::QueuedConnection);
}

void QuickSettingContainer::initConnection()
{
    connect(m_pluginLoader, &QuickSettingController::pluginInserted, this, &QuickSettingContainer::onPluginInsert);
    connect(m_pluginLoader, &QuickSettingController::pluginRemoved, this, &QuickSettingContainer::onPluginRemove);
    connect(m_playerWidget, &MediaWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_volumnWidget, &VolumeWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_volumnWidget, &VolumeWidget::rightIconClick, this, [ this ] {
        showWidget(m_volumeSettingWidget, tr("voice"));
        resizeView();
    });
    connect(m_brihtnessWidget, &BrightnessWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_brihtnessWidget, &BrightnessWidget::rightIconClicked, this, [ this ] {
        showWidget(m_brightSettingWidget, tr("brightness"));
        resizeView();
    });
    connect(m_childPage, &PluginChildPage::back, this, [ this ] {
        m_switchLayout->setCurrentWidget(m_mainWidget);
    });
    connect(m_childPage, &PluginChildPage::closeSelf, this, [ this ] {
        if (!m_childPage->isBack())
            topLevelWidget()->hide();
    });
}

void QuickSettingContainer::resizeView()
{
    if (m_switchLayout->currentWidget() == m_mainWidget) {
        QList<QuickSettingItem *> pluginItems = m_pluginLoader->settingItems();
        int selfPluginCount = 0;
        for (QuickSettingItem *item : pluginItems) {
            // 如果是置顶的插件，则认为它占用两个普通插件的位置
            int increCount = (item->pluginItem()->isPrimary() ? 2 : 1);
            selfPluginCount += increCount;
        }
        int rowCount = selfPluginCount / COLUMNCOUNT;
        if (selfPluginCount % COLUMNCOUNT > 0)
            rowCount++;

        m_pluginWidget->setFixedHeight(ITEMHEIGHT * rowCount + ITEMSPACE * (rowCount - 1));

        int panelCount = 0;
        if (m_playerWidget->isVisible())
            panelCount++;
        if (m_volumnWidget->isVisible())
            panelCount++;
        if (m_brihtnessWidget->isVisible())
            panelCount++;

        m_componentWidget->setFixedHeight(ITEMHEIGHT * panelCount + ITEMSPACE * (panelCount - 1));
        setFixedHeight(ITEMSPACE * 3 + m_pluginWidget->height() + m_componentWidget->height());
    } else if (m_switchLayout->currentWidget() == m_childPage) {
        setFixedHeight(m_childPage->height());
    }
}
