// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "expandiconwidget.h"
#include "taskmanager/taskmanager.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "tray_delegate.h"
#include "dockpopupwindow.h"
#include "imageutil.h"
#include "systempluginitem.h"
#include "docksettings.h"

#include <DGuiApplicationHelper>
#include <DRegionMonitor>
#include <QBitmap>
#include <QPainter>
#include <QPainterPath>

#include <qobjectdefs.h>
#include <xcb/xproto.h>

DGUI_USE_NAMESPACE

using RegionMonitor = Dtk::Gui::DRegionMonitor;

ExpandIconWidget::ExpandIconWidget(QWidget *parent, Qt::WindowFlags f)
    : BaseTrayWidget(parent, f)
    , m_position(Dock::Position::Bottom)
{
}

ExpandIconWidget::~ExpandIconWidget()
{
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
    } else {
        gridParentView->hide();
    }
}

QPixmap ExpandIconWidget::icon()
{
    return QPixmap(dropIconFile());
}

void ExpandIconWidget::paintEvent(QPaintEvent *event)
{
    QPainter painter(this);
    QPixmap pixmap = ImageUtil::loadSvg(dropIconFile(), QSize(ICON_SIZE, ICON_SIZE));
    QRect rectOfPixmap(rect().x() + (rect().width() - ICON_SIZE) / 2,
                    rect().y() + (rect().height() - ICON_SIZE) / 2,
                    ICON_SIZE, ICON_SIZE);

    painter.drawPixmap(rectOfPixmap, pixmap);
}

void ExpandIconWidget::moveEvent(QMoveEvent *event)
{
    BaseTrayWidget::moveEvent(event);
    // 当前展开按钮位置发生变化的时候，需要同时改变托盘的位置
    QMetaObject::invokeMethod(this, [] {
        TrayGridWidget *gridView = popupTrayView();
        if (gridView->isVisible())
            gridView->resetPosition();
    }, Qt::QueuedConnection);
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
    TrayGridView *trayView = TrayGridView::getIconTrayGridView(gridParentView);
    TrayDelegate *trayDelegate = TrayDelegate::getIconTrayDelegate(trayView, trayView);
    TrayModel *trayModel = TrayModel::getIconModel();
    gridParentView->setTrayGridView(trayView);

    gridParentView->setWindowFlags(Qt::FramelessWindowHint | Qt::ToolTip | Qt::WindowStaysOnTopHint | Qt::WindowDoesNotAcceptFocus);
    trayView->setModel(trayModel);
    trayView->setItemDelegate(trayDelegate);
    trayView->setSpacing(ITEM_SPACING);
    trayView->setDragDistance(2);

    QVBoxLayout *layout = new QVBoxLayout(gridParentView);
    layout->setContentsMargins(ITEM_SPACING, ITEM_SPACING, ITEM_SPACING, ITEM_SPACING);
    layout->setSpacing(0);
    layout->addWidget(trayView);

    auto rowCountChanged = [ = ] {
        if (gridParentView->isVisible()) {
            int count = trayModel->rowCount();
            if (count > 0)
                gridParentView->resetPosition();
            else
                gridParentView->hide();
        }
    };

    connect(trayModel, &TrayModel::rowCountChanged, gridParentView, rowCountChanged);
    connect(trayModel, &TrayModel::requestRefreshEditor, trayView, &TrayGridView::onUpdateEditorView);

    connect(trayDelegate, &TrayDelegate::requestHide, trayView, &TrayGridView::requestHide);
    connect(trayDelegate, &TrayDelegate::removeRow, trayView, [ = ](const QModelIndex &index) {
        trayView->model()->removeRow(index.row(),index.parent());
    });
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
    : DBlurEffectWidget (parent)
    , m_trayGridView(nullptr)
    , m_referGridView(nullptr)
    , m_regionInter(new RegionMonitor(this))
{
    initMember();
    setAttribute(Qt::WA_TranslucentBackground);
}

void TrayGridWidget::setPosition(const Dock::Position &position)
{
    m_position = position;
}

void TrayGridWidget::setTrayGridView(TrayGridView *trayView)
{
    m_trayGridView = trayView;
    connect(m_trayGridView, &TrayGridView::requestHide, this, &TrayGridWidget::hide);
}

void TrayGridWidget::setReferGridView(TrayGridView *trayView)
{
    m_referGridView = trayView;
}

TrayGridView *TrayGridWidget::trayView() const
{
    return m_trayGridView;
}

void TrayGridWidget::resetPosition()
{
    // 如果没有设置所属窗体，则无法计算位置
    ExpandIconWidget *expWidget = expandWidget();
    if (!expWidget)
        return;

    m_trayGridView->setFixedSize(m_trayGridView->suitableSize());
    setFixedSize(m_trayGridView->size() + QSize(ITEM_SPACING * 2, ITEM_SPACING * 2));

    QWidget *topWidget = expWidget->topLevelWidget();
    QPoint ptPos = expWidget->mapToGlobal(QPoint(0, 0));
    switch (m_position) {
    case Dock::Position::Bottom: {
        ptPos.setY(topWidget->y() - height());
        break;
    }
    case Dock::Position::Top: {
        ptPos.setY(topWidget->y() + topWidget->height());
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
    move(ptPos);
}

void TrayGridWidget::showEvent(QShowEvent *event)
{
    TaskManager::instance()->setTrayGridWidgetVisible(true);
    TaskManager::instance()->updateHideState(true);
    m_regionInter->registerRegion();
    DBlurEffectWidget::showEvent(event);
}

void TrayGridWidget::hideEvent(QHideEvent *event)
{
    TaskManager::instance()->setTrayGridWidgetVisible(false);
    TaskManager::instance()->updateHideState(true);
    m_regionInter->unregisterRegion();
    // 在当前托盘区域隐藏后，需要设置任务栏区域的展开按钮的托盘为隐藏状态
    TrayModel::getDockModel()->updateOpenExpand(false);
    DBlurEffectWidget::hideEvent(event);
}

void TrayGridWidget::initMember()
{
    connect(m_regionInter, &RegionMonitor::buttonPress, this, [ = ](const QPoint &mousePos, const int flag) {
        // 如果当前是隐藏，那么在点击任何地方都隐藏
        if (!isVisible()) {
            hide();
            return;
        }

        if ((flag != RegionMonitor::WatchedFlags::Button_Left) && (flag != RegionMonitor::WatchedFlags::Button_Right))
            return;

        QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
        const QRect rect = QRect(ptPos, size());
        if (rect.contains(mousePos))
            return;
        // 如果点击的是展开区域，则不做任何处理，因为点击展开区域自己会处理
        if (m_referGridView) {
            QAbstractItemModel *model = m_referGridView->model();
            for (int i = 0; i < model->rowCount(); i++) {
                ExpandIconWidget *widget = qobject_cast<ExpandIconWidget *>(m_referGridView->indexWidget(model->index(i, 0)));
                if (!widget)
                    continue;

                QRect rectExpandWidget(widget->mapToGlobal(QPoint(0, 0)), widget->size());
                if (rectExpandWidget.contains(mousePos))
                    return;
            }
        }

        const QRect rctView(pos(), size());
        if (rctView.contains(mousePos))
            return;

        // 查看是否存在SystemPluginItem插件，在此处判断的原因是因为当弹出右键菜单的时候，如果鼠标在菜单上点击
        // 刚好把托盘区域给隐藏了，导致菜单也跟着隐藏，导致点击菜单的时候不生效
        QAbstractItemModel *dataModel = m_trayGridView->model();
        for (int i = 0; i < dataModel->rowCount(); i++) {
            QModelIndex index = dataModel->index(i, 0);
            BaseTrayWidget *widget = qobject_cast<BaseTrayWidget *>(m_trayGridView->indexWidget(index));
            if (widget && widget->containsPoint(mousePos))
                return;
        }

        hide();
    });
}

QColor TrayGridWidget::maskColor() const
{
    QColor color = DGuiApplicationHelper::standardPalette(DGuiApplicationHelper::instance()->themeType()).window().color();
    color.setAlpha(0);
    return color;
}

ExpandIconWidget *TrayGridWidget::expandWidget() const
{
    if (!m_referGridView)
        return nullptr;

    QAbstractItemModel *dataModel = m_referGridView->model();
    if (!dataModel)
        return nullptr;

    for (int i = 0; i < dataModel->rowCount(); i++) {
        QModelIndex index = dataModel->index(i, 0);
        ExpandIconWidget *widget = qobject_cast<ExpandIconWidget *>(m_referGridView->indexWidget(index));
        if (widget)
            return widget;
    }

    return nullptr;
}
