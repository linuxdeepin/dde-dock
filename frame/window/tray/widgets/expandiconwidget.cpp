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
{
    connect(m_regionInter, &DRegionMonitor::buttonPress, this, [ = ](const QPoint &mousePos, const int flag) {
        TrayGridWidget *gridView = popupTrayView();
        // 如果当前是隐藏，那么在点击任何地方都隐藏
        if (!isVisible()) {
            gridView->hide();
            return;
        }

        if ((flag != DRegionMonitor::WatchedFlags::Button_Left) && (flag != DRegionMonitor::WatchedFlags::Button_Right))
            return;

        QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
        const QRect rect = QRect(ptPos, size());
        if (rect.contains(mousePos))
            return;

        const QRect rctView(gridView->pos(), gridView->size());
        if (rctView.contains(mousePos))
            return;

        gridView->hide();
    });
}

ExpandIconWidget::~ExpandIconWidget()
{
    TrayGridWidget *gridView = popupTrayView();
    gridView->setOwnerWidget(nullptr);
    setTrayPanelVisible(false);
}

void ExpandIconWidget::setPositon(Dock::Position position)
{
    if (m_position != position)
        m_position = position;

    TrayGridWidget::setPosition(position);
}

void ExpandIconWidget::sendClick(uint8_t mouseButton, int x, int y)
{
    Q_UNUSED(x);
    Q_UNUSED(y);

    // 如果当前图标不可见，则不让展开托盘列表
    if (popupTrayView()->trayView()->model()->rowCount() == 0)
        return;

    if (mouseButton != XCB_BUTTON_INDEX_1)
        return;

    QWidget *gridParentView = popupTrayView();
    setTrayPanelVisible(!gridParentView->isVisible());
}

void ExpandIconWidget::setTrayPanelVisible(bool visible)
{
    TrayGridWidget *gridParentView = popupTrayView();
    if (visible) {
        gridParentView->resetPosition();
        gridParentView->show();
        m_regionInter->registerRegion();
    } else {
        gridParentView->hide();
        m_regionInter->unregisterRegion();
    }
}

QPixmap ExpandIconWidget::icon()
{
    return QPixmap(dropIconFile());
}

void ExpandIconWidget::paintEvent(QPaintEvent *event)
{
    TrayGridWidget *gridView = popupTrayView();
    if (gridView->trayView()->model()->rowCount() == 0)
        return BaseTrayWidget::paintEvent(event);

    QPainter painter(this);
    QPixmap pixmap = ImageUtil::loadSvg(dropIconFile(), QSize(ICON_SIZE, ICON_SIZE));
    QRect rectOfPixmap(rect().x() + (rect().width() - ICON_SIZE) / 2,
                    rect().y() + (rect().height() - ICON_SIZE) / 2,
                    ICON_SIZE, ICON_SIZE);

    painter.drawPixmap(rectOfPixmap, pixmap);

    gridView->setOwnerWidget(this);
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

TrayGridWidget *ExpandIconWidget::popupTrayView()
{
    static TrayGridWidget *gridParentView = nullptr;
    if (gridParentView)
        return gridParentView;

    gridParentView = new TrayGridWidget(nullptr);
    TrayGridView *trayView = new TrayGridView(gridParentView);
    TrayDelegate *trayDelegate = new TrayDelegate(trayView, trayView);
    TrayModel *trayModel = TrayModel::getIconModel();
    gridParentView->setTrayGridView(trayView);

    gridParentView->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool);
    trayView->setModel(trayModel);
    trayView->setItemDelegate(trayDelegate);
    trayView->setSpacing(ITEM_SPACING);
    trayView->setDragDistance(2);

    QVBoxLayout *layout = new QVBoxLayout(gridParentView);
    layout->setContentsMargins(ITEM_SPACING, ITEM_SPACING, ITEM_SPACING, ITEM_SPACING);
    layout->setSpacing(0);
    layout->addWidget(trayView);

    auto rowCountChanged = [ = ] {
        int count = trayModel->rowCount();
        if (count > 0)
            gridParentView->resetPosition();
        else if (gridParentView->isVisible())
            gridParentView->hide();
    };

    connect(trayModel, &TrayModel::rowCountChanged, gridParentView, rowCountChanged);

    connect(trayDelegate, &TrayDelegate::removeRow, trayView, [ = ](const QModelIndex &index) {
        trayView->model()->removeRow(index.row(),index.parent());
    });
    connect(trayView, &TrayGridView::requestRemove, trayModel, &TrayModel::removeRow);
    connect(trayModel, &TrayModel::requestOpenEditor, trayView, [ trayView ](const QModelIndex &index) {
        trayView->openPersistentEditor(index);
    });

    QMetaObject::invokeMethod(gridParentView, rowCountChanged, Qt::QueuedConnection);

    return gridParentView;
}

/**
 * @brief 圆角窗体的绘制
 * @param parent
 */

Dock::Position TrayGridWidget::m_position = Dock::Position::Bottom;

TrayGridWidget::TrayGridWidget(QWidget *parent)
    : QWidget (parent)
    , m_dockInter(new DockInter(dockServiceName(), dockServicePath(), QDBusConnection::sessionBus(), this))
    , m_trayGridView(nullptr)
    , m_ownerWidget(nullptr)
{
    setAttribute(Qt::WA_TranslucentBackground);
}

void TrayGridWidget::setPosition(const Dock::Position &position)
{
    m_position = position;
}

void TrayGridWidget::setTrayGridView(TrayGridView *trayView)
{
    m_trayGridView = trayView;
}

void TrayGridWidget::setOwnerWidget(QWidget *widget)
{
    // 设置所属的Widget，目的是为了计算当前窗体的具体位置
    m_ownerWidget = widget;
}

TrayGridView *TrayGridWidget::trayView() const
{
    return m_trayGridView;
}

void TrayGridWidget::resetPosition()
{
    // 如果没有设置所属窗体，则无法计算位置
    if (!m_ownerWidget || !m_ownerWidget->parentWidget())
        return;

    QWidget *topWidget = m_ownerWidget->topLevelWidget();
    QPoint ptPos = m_ownerWidget->parentWidget()->mapToGlobal(m_ownerWidget->pos());
    switch (m_position) {
    case Dock::Position::Bottom: {
        ptPos.setX(ptPos.x() - width());
        ptPos.setY(topWidget->y() - height());
        break;
    }
    case Dock::Position::Top: {
        ptPos.setY(topWidget->y() + topWidget->height());
        ptPos.setX(ptPos.x() - width());
        break;
    }
    case Dock::Position::Left: {
        ptPos.setX(topWidget->x() + topWidget->width());
        break;
    }
    case Dock::Position::Right: {
        ptPos.setX(topWidget->x() - width());
        break;
    }
    }
    m_trayGridView->setFixedSize(m_trayGridView->suitableSize());
    setFixedSize(m_trayGridView->size() + QSize(ITEM_SPACING * 2, ITEM_SPACING * 2));
    move(ptPos);
}

void TrayGridWidget::paintEvent(QPaintEvent *event)
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

QColor TrayGridWidget::maskColor() const
{
    QColor color = DGuiApplicationHelper::standardPalette(DGuiApplicationHelper::instance()->themeType()).window().color();
    int maskAlpha(static_cast<int>(255 * m_dockInter->opacity()));
    color.setAlpha(maskAlpha);
    return color;
}
