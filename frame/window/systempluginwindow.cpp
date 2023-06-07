// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "systempluginwindow.h"
#include "docksettings.h"
#include "systempluginitem.h"
#include "quicksettingcontroller.h"
#include "utils.h"

#include <DListView>
#include <DGuiApplicationHelper>

#include <QBoxLayout>
#include <QDir>
#include <QMetaObject>
#include <QGuiApplication>
#include <QPainter>
#include <sys/types.h>

#define MAXICONSIZE 48
#define MINICONSIZE 24
#define ICONMARGIN 8

SystemPluginWindow::SystemPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_listView(new DListView(this))
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_position(Dock::Position::Bottom)
    , m_mainLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
{
    initUi();
    initConnection();
}

SystemPluginWindow::~SystemPluginWindow()
{
}

void SystemPluginWindow::setDisplayMode(const DisplayMode &displayMode)
{
    m_displayMode = displayMode;
    QList<StretchPluginsItem *> items = stretchItems();
    switch (m_position) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        for (StretchPluginsItem *item : items)
            item->setDisplayMode(displayMode);
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        for (StretchPluginsItem *item : items)
            item->setDisplayMode(displayMode);
        break;
    }
    }
}

void SystemPluginWindow::setPositon(Position position)
{
    if (m_position == position)
        return;

    m_position = position;

    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    else
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);

    StretchPluginsItem::setPosition(position);

    QObjectList childObjects = children();
    for (QObject *childObject : childObjects) {
        StretchPluginsItem *item = qobject_cast<StretchPluginsItem *>(childObject);
        if (!item)
            continue;

        item->update();
    }
}

QSize SystemPluginWindow::suitableSize() const
{
    return suitableSize(m_position);
}

QSize SystemPluginWindow::suitableSize(const Position &position) const
{
    QList<StretchPluginsItem *> items = stretchItems();
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        int itemWidth = 0;
        for (StretchPluginsItem *item : items)
            itemWidth += item->suitableSize(position).width();

        return QSize(itemWidth, QWIDGETSIZE_MAX);
    }

    int itemHeight = 0;
    for (StretchPluginsItem *item : items)
        itemHeight += item->suitableSize(position).height();

    return QSize(QWIDGETSIZE_MAX, itemHeight);
}

bool SystemPluginWindow::eventFilter(QObject *watched, QEvent *event)
{
    if (event->type() == QEvent::Drop)
        Q_EMIT requestDrop(static_cast<QDropEvent *>(event));

    return QWidget::eventFilter(watched, event);
}

void SystemPluginWindow::initUi()
{
    m_mainLayout->setContentsMargins(0, 0, 0, 0);
    m_mainLayout->setSpacing(0);
    installEventFilter(this);
}

void SystemPluginWindow::initConnection()
{
    QuickSettingController *quickController = QuickSettingController::instance();
    connect(quickController, &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::System)
            return;

        pluginAdded(itemInter);
    });

    connect(quickController, &QuickSettingController::pluginRemoved, this, &SystemPluginWindow::onPluginItemRemoved);
    connect(quickController, &QuickSettingController::pluginUpdated, this, &SystemPluginWindow::onPluginItemUpdated);

    QList<PluginsItemInterface *> plugins = quickController->pluginItems(QuickSettingController::PluginAttribute::System);
    for (int i = 0; i < plugins.size(); i++)
        pluginAdded(plugins[i]);
}

StretchPluginsItem *SystemPluginWindow::findPluginItemWidget(PluginsItemInterface *pluginItem)
{
    for (int i = 0; i < m_mainLayout->count(); i++) {
        QLayoutItem *layoutItem = m_mainLayout->itemAt(i);
        if (!layoutItem)
            continue;

        StretchPluginsItem *itemWidget = qobject_cast<StretchPluginsItem *>(layoutItem->widget());
        if (itemWidget && itemWidget->pluginInter() == pluginItem)
            return itemWidget;
    }

    return nullptr;
}

void SystemPluginWindow::pluginAdded(PluginsItemInterface *plugin)
{
    StretchPluginsItem *item = new StretchPluginsItem(plugin, QuickSettingController::instance()->itemKey(plugin));
    item->setDisplayMode(m_displayMode);
    item->setPosition(m_position);
    item->installEventFilter(this);
    item->setParent(this);
    item->show();
    m_mainLayout->addWidget(item);
    Q_EMIT itemChanged();
}

QList<StretchPluginsItem *> SystemPluginWindow::stretchItems() const
{
    QList<StretchPluginsItem *> items;
    QObjectList childObjects = children();
    for (QObject *childObject : childObjects) {
        StretchPluginsItem *item = qobject_cast<StretchPluginsItem *>(childObject);
        if (!item)
            continue;

        items << item;
    }
    return items;
}

void SystemPluginWindow::onPluginItemRemoved(PluginsItemInterface *pluginItem)
{
    StretchPluginsItem *item = findPluginItemWidget(pluginItem);
    if (item) {
        item->setParent(nullptr);
        item->hide();
        m_mainLayout->removeWidget(item);
        Q_EMIT itemChanged();
    }
}

void SystemPluginWindow::onPluginItemUpdated(PluginsItemInterface *pluginItem)
{
    StretchPluginsItem *item = findPluginItemWidget(pluginItem);
    if (item)
        item->update();
}

// 图标的尺寸
#define ICONSIZE 20
#define ICONTEXTSPACE 6
#define PLUGIN_ITEM_DRAG_THRESHOLD 20

Dock::Position StretchPluginsItem::m_position = Dock::Position::Bottom;

StretchPluginsItem::StretchPluginsItem(PluginsItemInterface * const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_windowSizeFashion(DockSettings::instance()->getWindowSizeFashion())
{
    connect(DockSettings::instance(), &DockSettings::windowSizeFashionChanged, this, [=](uint size) {
        m_windowSizeFashion = size;
    });
}

StretchPluginsItem::~StretchPluginsItem()
{
}

void StretchPluginsItem::setDisplayMode(const DisplayMode &displayMode)
{
    m_displayMode = displayMode;
}

void StretchPluginsItem::setPosition(Position position)
{
    m_position = position;
}

QString StretchPluginsItem::itemKey() const
{
    return m_itemKey;
}

QSize StretchPluginsItem::suitableSize() const
{
    return suitableSize(m_position);
}

PluginsItemInterface *StretchPluginsItem::pluginInter() const
{
    return m_pluginInter;
}

void StretchPluginsItem::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);
    QPainter painter(this);
    QIcon icon = m_pluginInter->icon(DockPart::SystemPanel);

    QRect rctPixmap(rect());
    if (needShowText()) {
        int textHeight = QFontMetrics(textFont()).height();
        // 文本与图标的间距为6
        int iconTop = (height() - textHeight - ICONSIZE - ICONTEXTSPACE) / 2;
        rctPixmap.setX((width() - ICONSIZE) / 2);
        rctPixmap.setY(iconTop);
        rctPixmap.setWidth(ICONSIZE);
        rctPixmap.setHeight(ICONSIZE);
        // 先绘制下面的文本
        painter.setFont(textFont());
        painter.drawText(QRect(0, iconTop + ICONSIZE + ICONTEXTSPACE, width(), textHeight), Qt::AlignCenter, m_pluginInter->pluginDisplayName());
    } else {
        rctPixmap.setX((width() - ICONSIZE) / 2);
        rctPixmap.setY((height() - ICONSIZE) / 2);
        rctPixmap.setWidth(ICONSIZE);
        rctPixmap.setHeight(ICONSIZE);
    }

    // 绘制图标
    int iconSize = static_cast<int>(ICONSIZE * (QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? 1 : qApp->devicePixelRatio()));
    painter.drawPixmap(rctPixmap, icon.pixmap(iconSize, iconSize));
}

QSize StretchPluginsItem::suitableSize(const Position &position) const
{
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        int textWidth = 0;
        if (needShowText())
            textWidth = QFontMetrics(textFont(position)).boundingRect(m_pluginInter->pluginDisplayName()).width();
        return QSize(qMax(textWidth, ICONSIZE) + (m_displayMode == Dock::DisplayMode::Efficient ? 5 : 10) * 2, -1);
    }

    int height = 6;                                                 // 图标上边距6
    height += ICONSIZE;                                             // 图标尺寸20
    if (m_displayMode == Dock::DisplayMode::Fashion) {
        height += ICONTEXTSPACE;                                    // 图标与文字间距6
        if (needShowText())                                         // 只有在显示文本的时候才计算文本的高度
            height += QFontMetrics(textFont(position)).height();    // 文本高度
    }
    height += 4;                                                    // 下间距4
    return QSize(-1, height);
}

QFont StretchPluginsItem::textFont() const
{
    return textFont(m_position);
}

QFont StretchPluginsItem::textFont(const Position &position) const
{
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        static QList<QFont> fonts{ DFontSizeManager::instance()->t9(),
                              DFontSizeManager::instance()->t8(),
                              DFontSizeManager::instance()->t7(),
                              DFontSizeManager::instance()->t6() };
#define MINHEIGHT 50
        // 如果当前的实际位置和请求的位置不一致，说明当前正在切换位置，此时将它的宽度作为它的高度(左到下切换的时候，左侧的宽度和下面的高度一致)
        int size = (m_position == position ? height() : width());
        int index = qMin(qMax((size - MINHEIGHT) / 2, 0), fonts.size() - 1);
        return fonts[index];
    }

    return DFontSizeManager::instance()->t10();
}

bool StretchPluginsItem::needShowText() const
{
    // 如果是高效模式，则不需要显示下面的文本
    if (m_displayMode == Dock::DisplayMode::Efficient)
        return false;

    // 图标与文本，图标距离上方和文本距离下方的尺寸
#define SPACEMARGIN 6
    // 文本的高度
#define TEXTSIZE 14
    // 当前插件父窗口与顶层窗口的上下边距
#define OUTMARGIN 7

    // 任务栏在上方或者下方显示的时候，根据设计图，只有在当前区域高度大于50的时候才同时显示文本和图标
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return ((Utils::isDraging() ? topLevelWidget()->height() : m_windowSizeFashion) >=
                                    (OUTMARGIN * 2 + SPACEMARGIN * 3 + ICONSIZE + TEXTSIZE));

    return true;
}

const QString StretchPluginsItem::contextMenu() const
{
    return m_pluginInter->itemContextMenu(m_itemKey);
}

void StretchPluginsItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    m_pluginInter->invokedMenuItem(m_itemKey, itemId, checked);
}

QWidget *StretchPluginsItem::popupTips()
{
    return m_pluginInter->itemTipsWidget(m_itemKey);
}

void StretchPluginsItem::mousePressEvent(QMouseEvent *e)
{
    m_hover = false;
    update();

    if (PopupWindow->isVisible())
        hideNonModel();

    if (e->button() == Qt::LeftButton)
        m_mousePressPoint = e->pos();

    m_popupTipsDelayTimer->stop();
    hideNonModel();

    if (e->button() == Qt::RightButton
        && perfectIconRect().contains(e->pos()))
            return showContextMenu();

    DockItem::mousePressEvent(e);
}

void StretchPluginsItem::mouseReleaseEvent(QMouseEvent *e)
{
    DockItem::mouseReleaseEvent(e);

    if (e->button() != Qt::LeftButton)
        return;

    if (checkAndResetTapHoldGestureState() && e->source() == Qt::MouseEventSynthesizedByQt)
        return;

    const QPoint distance = e->pos() - m_mousePressPoint;
    if (distance.manhattanLength() < PLUGIN_ITEM_DRAG_THRESHOLD)
        mouseClick();
}

void StretchPluginsItem::enterEvent(QEvent *event)
{
    if (auto view = qobject_cast<SystemPluginWindow *>(parentWidget()))
        view->requestDrawBackground(rect());

    update();
    DockItem::enterEvent(event);
}

void StretchPluginsItem::leaveEvent(QEvent *event)
{
    if (auto view = qobject_cast<SystemPluginWindow *>(parentWidget()))
        view->requestDrawBackground(QRect());
    update();
    DockItem::leaveEvent(event);
}

void StretchPluginsItem::mouseClick()
{
    QStringList commandArgument = m_pluginInter->itemCommand(m_itemKey).split(" ");
    if (commandArgument.size() > 0) {
        QString command = commandArgument.first();
        commandArgument.removeFirst();
        QProcess::startDetached(command, commandArgument);
        return;
    }

    // request popup applet
    if (QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey))
        showPopupApplet(w);
}
