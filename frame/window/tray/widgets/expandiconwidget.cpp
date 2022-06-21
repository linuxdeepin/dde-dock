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
#include "expandiconwidget.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "tray_delegate.h"
#include "dockpopupwindow.h"
#include "imageutil.h"

#include <DGuiApplicationHelper>
#include <DRegionMonitor>
#include <QBitmap>
#include <QPainter>
#include <QPainterPath>

#include <xcb/xproto.h>

DGUI_USE_NAMESPACE

ExpandIconWidget::ExpandIconWidget(QWidget *parent, Qt::WindowFlags f)
    : BaseTrayWidget(parent, f)
    , m_regionInter(new DRegionMonitor(this))
    , m_position(Dock::Position::Bottom)
    , m_gridParentView(new RoundWidget(nullptr))
    , m_trayView(new TrayGridView(m_gridParentView))
    , m_trayDelegate(new TrayDelegate(m_trayView, m_trayView))
    , m_trayModel(new TrayModel(m_trayView, true, false))
{
    initUi();
    initConnection();
}

ExpandIconWidget::~ExpandIconWidget()
{
    m_gridParentView->deleteLater();
}

void ExpandIconWidget::setPositonValue(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
}

void ExpandIconWidget::sendClick(uint8_t mouseButton, int x, int y)
{
    Q_UNUSED(x);
    Q_UNUSED(y);

    if (mouseButton != XCB_BUTTON_INDEX_1)
        return;

    setTrayPanelVisible(!m_gridParentView->isVisible());
}

void ExpandIconWidget::setTrayPanelVisible(bool visible)
{
    if (visible) {
        resetPosition();
        m_gridParentView->show();
        m_regionInter->registerRegion();
    } else {
        m_gridParentView->hide();
        m_regionInter->unregisterRegion();
    }
}

QPixmap ExpandIconWidget::icon()
{
    return QPixmap(dropIconFile());
}

void ExpandIconWidget::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    QPixmap pixmap = ImageUtil::loadSvg(dropIconFile(), QSize(ICON_SIZE, ICON_SIZE));
    QRect rectOfPixmap(rect().x() + (rect().width() - ICON_SIZE) / 2,
                    rect().y() + (rect().height() - ICON_SIZE) / 2,
                    ICON_SIZE, ICON_SIZE);

    painter.drawPixmap(rectOfPixmap, pixmap);
}

const QString ExpandIconWidget::dropIconFile() const
{
    QString arrow;
    switch (m_position) {
    case Dock::Position::Bottom: {
        arrow = "up";
        break;
    }
    case Dock::Position::Top: {
        arrow = "down";
        break;
    }
    case Dock::Position::Left: {
        arrow = "right";
        break;
    }
    case Dock::Position::Right: {
        arrow = "left";
        break;
    }
    }

    QString iconFile = QString(":/icons/resources/arrow-%1").arg(arrow);
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconFile += QString("-dark");

    return iconFile + ".svg";
}

QWidget *ExpandIconWidget::popupTrayView()
{
    return m_gridParentView;
}

void ExpandIconWidget::resetPosition()
{
    if (!parentWidget())
        return;

    QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
    switch (m_position) {
    case Dock::Position::Bottom: {
        ptPos.setY(ptPos.y() - m_gridParentView->height());
        ptPos.setX(ptPos.x() - m_gridParentView->width());
        break;
    }
    case Dock::Position::Top: {
        ptPos.setY(ptPos.y() + m_gridParentView->height());
        ptPos.setX(ptPos.x() - m_gridParentView->width());
        break;
    }
    case Dock::Position::Left: {
        ptPos.setX(ptPos.x() + m_gridParentView->width() / 2);
        break;
    }
    case Dock::Position::Right: {
        ptPos.setX(ptPos.x() - m_gridParentView->width() / 2);
        break;
    }
    }
    m_gridParentView->move(ptPos);
}

void ExpandIconWidget::initUi()
{
    m_gridParentView->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool);
    m_trayView->setModel(m_trayModel);
    m_trayView->setItemDelegate(m_trayDelegate);
    m_trayView->setSpacing(ITEM_SPACING);
    m_trayView->setDragDistance(2);

    QVBoxLayout *layout = new QVBoxLayout(m_gridParentView);
    layout->setContentsMargins(ITEM_SPACING, ITEM_SPACING, ITEM_SPACING, ITEM_SPACING);
    layout->setSpacing(0);
    layout->addWidget(m_trayView);
}

void ExpandIconWidget::initConnection()
{
    connect(m_trayView, &TrayGridView::rowCountChanged, this, &ExpandIconWidget::onRowCountChanged);

    connect(m_trayDelegate, &TrayDelegate::removeRow, this, [ = ](const QModelIndex &index) {
        m_trayView->model()->removeRow(index.row(),index.parent());
    });
    connect(m_trayView, &TrayGridView::requestRemove, m_trayModel, &TrayModel::removeRow);
    connect(m_regionInter, &DRegionMonitor::buttonPress, this, &ExpandIconWidget::onGlobMousePress);

    QMetaObject::invokeMethod(this, &ExpandIconWidget::onRowCountChanged, Qt::QueuedConnection);
}

void ExpandIconWidget::onRowCountChanged()
{
    int count = m_trayModel->rowCount();
    m_trayView->setFixedSize(m_trayView->suitableSize());
    m_gridParentView->setFixedSize(m_trayView->size() + QSize(ITEM_SPACING * 2, ITEM_SPACING * 2));
    if (count > 0)
        resetPosition();
    else if (m_gridParentView->isVisible())
        m_gridParentView->hide();

    Q_EMIT trayVisbleChanged(count > 0);
}

void ExpandIconWidget::onGlobMousePress(const QPoint &mousePos, const int flag)
{
    // 如果当前是隐藏，那么在点击任何地方都隐藏
    if (!isVisible()) {
        m_gridParentView->hide();
        return;
    }

    if ((flag != DRegionMonitor::WatchedFlags::Button_Left) && (flag != DRegionMonitor::WatchedFlags::Button_Right))
        return;

    QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
    const QRect rect = QRect(ptPos, size());
    if (rect.contains(mousePos))
        return;

    const QRect rctView(m_gridParentView->pos(), m_gridParentView->size());
    if (rctView.contains(mousePos))
        return;

    m_gridParentView->hide();
}

/**
 * @brief 圆角窗体的绘制
 * @param parent
 */

RoundWidget::RoundWidget(QWidget *parent)
    : QWidget (parent)
    , m_dockInter(new DockInter(dockServiceName(), dockServicePath(), QDBusConnection::sessionBus(), this))
{
    setAttribute(Qt::WA_TranslucentBackground);
}

void RoundWidget::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    QPainterPath path;
    path.addRoundedRect(rect(), 18, 18);
    painter.setCompositionMode(QPainter::CompositionMode_Xor);
    painter.setClipPath(path);

    painter.fillPath(path, maskColor());
}

QColor RoundWidget::maskColor() const
{
    QColor color = DGuiApplicationHelper::standardPalette(DGuiApplicationHelper::instance()->themeType()).window().color();
    int maskAlpha(static_cast<int>(255 * m_dockInter->opacity()));
    color.setAlpha(maskAlpha);
    return color;
}
