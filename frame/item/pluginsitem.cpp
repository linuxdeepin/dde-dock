/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "constants.h"
#include "pluginsitem.h"
#include "pluginsiteminterface.h"

#include "util/imagefactory.h"

#include <QPainter>
#include <QBoxLayout>
#include <QMouseEvent>
#include <QDrag>
#include <QMimeData>

#define PLUGIN_ITEM_DRAG_THRESHOLD      20

QPoint PluginsItem::MousePressPoint = QPoint();

PluginsItem::PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(parent),
      m_pluginInter(pluginInter),
      m_centralWidget(m_pluginInter->itemWidget(itemKey)),
      m_itemKey(itemKey),
      m_dragging(false),
      m_hover(false)
{
    qDebug() << "load plugins item: " << pluginInter->pluginName() << itemKey << m_centralWidget;

    m_centralWidget->setParent(this);
    m_centralWidget->setVisible(true);
    m_centralWidget->installEventFilter(this);

    QBoxLayout *hLayout = new QHBoxLayout;
    hLayout->addWidget(m_centralWidget);
    hLayout->setSpacing(0);
    hLayout->setMargin(0);

    setLayout(hLayout);
    setAccessibleName(pluginInter->pluginName() + "-" + m_itemKey);
    setAttribute(Qt::WA_TranslucentBackground);
}

PluginsItem::~PluginsItem()
{
}

int PluginsItem::itemSortKey() const
{
    return m_pluginInter->itemSortKey(m_itemKey);
}

void PluginsItem::setItemSortKey(const int order) const
{
    m_pluginInter->setSortKey(m_itemKey, order);
}

void PluginsItem::detachPluginWidget()
{
    QWidget *widget = m_pluginInter->itemWidget(m_itemKey);
    if (widget)
        widget->setParent(nullptr);
}

bool PluginsItem::allowContainer() const
{
    if (DockDisplayMode == Dock::Fashion)
        return false;

    return m_pluginInter->itemAllowContainer(m_itemKey);
}

bool PluginsItem::isInContainer() const
{
    if (DockDisplayMode == Dock::Fashion)
        return false;

    return m_pluginInter->itemIsInContainer(m_itemKey);
}

void PluginsItem::setInContainer(const bool container)
{
    m_pluginInter->setItemIsInContainer(m_itemKey, container);
}

QString PluginsItem::pluginName() const
{
    return m_pluginInter->pluginName();
}

QSize PluginsItem::sizeHint() const
{
    return m_centralWidget->sizeHint();
}

void PluginsItem::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);

    DisplayMode displayMode = m_pluginInter->displayMode();

    if (displayMode == Dock::DisplayMode::Fashion) {
        return;
    }
    if (!m_hover || m_dragging) {
        return;
    }

    // draw hover background
    QRect destRect;
    destRect.setSize(m_centralWidget->sizeHint());
    destRect.moveCenter(rect().center());

    QPainterPath path;
    path.addRoundedRect(destRect, 6, 6);

    QColor color;
    color = QColor::fromRgb(255, 255, 255);

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setOpacity(0.1);

    painter.fillPath(path, color);
}

void PluginsItem::refershIcon()
{
    m_pluginInter->refreshIcon(m_itemKey);
}

QWidget *PluginsItem::centralWidget() const
{
    return m_centralWidget;
}

void PluginsItem::mousePressEvent(QMouseEvent *e)
{
    m_hover = false;
    update();

    if (!isInContainer() && PopupWindow->isVisible())
        hideNonModel();

    if (e->button() == Qt::LeftButton)
        MousePressPoint = e->pos();

    // context menu will handle in DockItem
    DockItem::mousePressEvent(e);
}

void PluginsItem::mouseMoveEvent(QMouseEvent *e)
{
    if (e->buttons() != Qt::LeftButton)
        return DockItem::mouseMoveEvent(e);

    e->accept();

    const QPoint distance = e->pos() - MousePressPoint;
    if (distance.manhattanLength() > PLUGIN_ITEM_DRAG_THRESHOLD)
        startDrag();
}

void PluginsItem::mouseReleaseEvent(QMouseEvent *e)
{
    DockItem::mouseReleaseEvent(e);

    if (e->button() != Qt::LeftButton)
        return;

    if (checkAndResetTapHoldGestureState()&& e->source() == Qt::MouseEventSynthesizedByQt) {
        qDebug() << "tap and hold gesture detected, ignore the synthesized mouse release event";
        return;
    }

    e->accept();

    const QPoint distance = e->pos() - MousePressPoint;
    if (distance.manhattanLength() < PLUGIN_ITEM_DRAG_THRESHOLD)
        mouseClicked();
}

void PluginsItem::enterEvent(QEvent *event)
{
    m_hover = true;
    update();

    DockItem::enterEvent(event);
}

void PluginsItem::leaveEvent(QEvent *event)
{
    // Note:
    // here we should check the mouse position to ensure the mouse is really leaved
    // because this leaveEvent will also be called if setX11PassMouseEvent(false) is invoked
    // in XWindowTrayWidget::sendHoverEvent()
    if (!rect().contains(mapFromGlobal(QCursor::pos()))) {
        m_hover = false;
        update();
    }

    DockItem::leaveEvent(event);
}

bool PluginsItem::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_centralWidget) {
        if (event->type() == QEvent::MouseButtonRelease) {
            m_hover = false;
            update();
        }
    }

    return false;
}

void PluginsItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    m_pluginInter->invokedMenuItem(m_itemKey, itemId, checked);
}

void PluginsItem::showPopupWindow(QWidget * const content, const bool model)
{
    if (isInContainer())
        return;

    DockItem::showPopupWindow(content, model);
}

const QString PluginsItem::contextMenu() const
{
    return m_pluginInter->itemContextMenu(m_itemKey);
}

QWidget *PluginsItem::popupTips()
{
    return m_pluginInter->itemTipsWidget(m_itemKey);
}

void PluginsItem::startDrag()
{
    const QPixmap pixmap = grab();

    m_dragging = true;
    m_centralWidget->setVisible(false);
    update();

    QMimeData *mime = new QMimeData;
    mime->setData(DOCK_PLUGIN_MIME, m_itemKey.toStdString().c_str());

    QDrag *drag = new QDrag(this);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center() / pixmap.devicePixelRatioF());
    drag->setMimeData(mime);

    emit dragStarted();
    const Qt::DropAction result = drag->exec(Qt::MoveAction);
    Q_UNUSED(result);
    emit itemDropped(drag->target(), QCursor::pos());

    m_dragging = false;
    m_hover = false;
    m_centralWidget->setVisible(true);
    setVisible(true);
    update();
}

void PluginsItem::mouseClicked()
{
    const QString command = m_pluginInter->itemCommand(m_itemKey);
    if (!command.isEmpty())
    {
        QProcess *proc = new QProcess(this);

        connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

        proc->startDetached(command);
        return;
    }

    // request popup applet
    QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey);
    if (w)
        showPopupApplet(w);
}
