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
#include "systempluginwindow.h"
#include "systemplugincontroller.h"
#include "systempluginitem.h"
#include "dockpluginscontroller.h"
#include "fixedplugincontroller.h"

#include <DListView>
#include <QBoxLayout>
#include <QDir>
#include <QMetaObject>

#define MAXICONSIZE 48
#define MINICONSIZE 24
#define ICONMARGIN 8

SystemPluginWindow::SystemPluginWindow(QWidget *parent)
    : QWidget(parent)
    , m_pluginController(new FixedPluginController(this))
    , m_listView(new DListView(this))
    , m_position(Dock::Position::Bottom)
    , m_mainLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
{
    initUi();
    connect(m_pluginController, &FixedPluginController::pluginItemInserted, this, &SystemPluginWindow::onPluginItemAdded);
    connect(m_pluginController, &FixedPluginController::pluginItemRemoved, this, &SystemPluginWindow::onPluginItemRemoved);
    connect(m_pluginController, &FixedPluginController::pluginItemUpdated, this, &SystemPluginWindow::onPluginItemUpdated);
}

SystemPluginWindow::~SystemPluginWindow()
{
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

    QObjectList childObjects = children();
    for (QObject *childObject : childObjects) {
        StretchPluginsItem *item = qobject_cast<StretchPluginsItem *>(childObject);
        if (!item)
            continue;

        item->setPosition(m_position);
    }
}

QSize SystemPluginWindow::suitableSize()
{
    QObjectList childs = children();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        int itemWidth = 0;
        for (QObject *childObject : childs) {
            StretchPluginsItem *childItem = qobject_cast<StretchPluginsItem *>(childObject);
            if (!childItem)
                continue;

            itemWidth += childItem->suitableSize().width();
        }

        return QSize(itemWidth, QWIDGETSIZE_MAX);
    }

    int itemHeight = 0;
    for (QObject *childObject : childs) {
        StretchPluginsItem *item = qobject_cast<StretchPluginsItem *>(childObject);
        if (!item)
            continue;

        itemHeight += item->suitableSize().height();
    }

    return QSize(QWIDGETSIZE_MAX, itemHeight);
}

void SystemPluginWindow::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);
    Q_EMIT sizeChanged();
}

void SystemPluginWindow::initUi()
{
    m_mainLayout->setContentsMargins(0, 0, 0, 0);
    m_mainLayout->setSpacing(0);
}

void SystemPluginWindow::onPluginItemAdded(StretchPluginsItem *pluginItem)
{
    if (m_mainLayout->children().contains(pluginItem))
        return;

    pluginItem->setPosition(m_position);
    pluginItem->setParent(this);
    pluginItem->show();
    m_mainLayout->addWidget(pluginItem);
    Q_EMIT sizeChanged();
}

void SystemPluginWindow::onPluginItemRemoved(StretchPluginsItem *pluginItem)
{
    if (!m_mainLayout->children().contains(pluginItem))
        return;

    pluginItem->setParent(nullptr);
    pluginItem->hide();
    m_mainLayout->removeWidget(pluginItem);
    Q_EMIT sizeChanged();
}

void SystemPluginWindow::onPluginItemUpdated(StretchPluginsItem *pluginItem)
{
    pluginItem->update();
}

#define ICONSIZE 20
#define ICONTEXTSPACE 6
#define PLUGIN_ITEM_DRAG_THRESHOLD 20

StretchPluginsItem::StretchPluginsItem(PluginsItemInterface * const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
    , m_position(Dock::Position::Bottom)
{
}

StretchPluginsItem::~StretchPluginsItem()
{
}

void StretchPluginsItem::setPosition(Position position)
{
    m_position = position;
    update();
}

QString StretchPluginsItem::itemKey() const
{
    return m_itemKey;
}

PluginsItemInterface *StretchPluginsItem::pluginInter() const
{
    return m_pluginInter;
}

void StretchPluginsItem::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);
    QPainter painter(this);
    const QIcon *icon = m_pluginInter->icon();

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
    if (icon)
        painter.drawPixmap(rctPixmap, icon->pixmap(ICONSIZE, ICONSIZE));
}

QSize StretchPluginsItem::suitableSize() const
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        int textWidth = QFontMetrics(textFont()).boundingRect(m_pluginInter->pluginDisplayName()).width();
        return QSize(qMax(textWidth, ICONSIZE) + 10 * 2, -1);
    }

    int height = 6;                                 // 图标上边距6
    height += ICONSIZE;                             // 图标尺寸20
    height += ICONTEXTSPACE;                        // 图标与文字间距6
    height += QFontMetrics(textFont()).height();    // 文本高度
    height += 4;                                    // 下间距4
    return QSize(-1, height);
}

QFont StretchPluginsItem::textFont() const
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        static QList<QFont> fonts{ DFontSizeManager::instance()->t9(),
                              DFontSizeManager::instance()->t8(),
                              DFontSizeManager::instance()->t7(),
                              DFontSizeManager::instance()->t6() };
#define MINHEIGHT 50
        int index = qMin(qMax((height() - MINHEIGHT) / 2, 0), fonts.size() - 1);
        return fonts[index];
    }

    return DFontSizeManager::instance()->t10();
}

bool StretchPluginsItem::needShowText() const
{
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return height() > (ICONSIZE + QFontMetrics(textFont()).height() + ICONTEXTSPACE);

    return true;
}

const QString StretchPluginsItem::contextMenu() const
{
    return m_pluginInter->itemContextMenu(m_itemKey);
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

void StretchPluginsItem::mouseClick()
{
    const QString command = m_pluginInter->itemCommand(m_itemKey);
    if (!command.isEmpty()) {
        QProcess::startDetached(command);
        return;
    }

    // request popup applet
    if (QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey))
        showPopupApplet(w);
}
