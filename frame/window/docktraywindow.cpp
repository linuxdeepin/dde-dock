// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "docktraywindow.h"
#include "datetimedisplayer.h"
#include "systempluginwindow.h"
#include "quickpluginwindow.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "tray_delegate.h"
#include "quicksettingcontroller.h"
#include "pluginsitem.h"
#include "expandiconwidget.h"
#include "quickdragcore.h"
#include "desktop_widget.h"

#include <DGuiApplicationHelper>

#include <QBoxLayout>
#include <QLabel>
#include <QPainter>

#define FRONTSPACING 18
#define SPLITERSIZE 2
#define SPLITESPACE 5

DockTrayWindow::DockTrayWindow(QWidget *parent)
    : QWidget(parent)
    , m_position(Dock::Position::Bottom)
    , m_displayMode(Dock::DisplayMode::Efficient)
    , m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::RightToLeft, this))
    , m_showDesktopWidget(new DesktopWidget(this))
    , m_toolWidget(new QWidget(this))
    , m_toolLayout(new QBoxLayout(QBoxLayout::RightToLeft, m_toolWidget))
    , m_toolLineLabel(new QLabel(this))
    , m_dateTimeWidget(new DateTimeDisplayer(true, this))
    , m_systemPuginWidget(new SystemPluginWindow(this))
    , m_quickIconWidget(new QuickPluginWindow(Dock::DisplayMode::Efficient, this))
    , m_trayView(TrayGridView::getDockTrayGridView(this))
    , m_model(TrayModel::getDockModel())
    , m_delegate(TrayDelegate::getDockTrayDelegate(m_trayView, this))
    , m_toolFrontSpaceWidget(new QWidget(this))
    , m_toolBackSpaceWidget(new QWidget(this))
    , m_dateTimeSpaceWidget(new QWidget(this))
{
    initUi();
    initConnection();
    initAttribute();
}

void DockTrayWindow::setPositon(const Dock::Position &position)
{
    m_position = position;
    m_dateTimeWidget->setPositon(position);
    m_systemPuginWidget->setPositon(position);
    m_quickIconWidget->setPositon(position);
    m_trayView->setPosition(position);
    m_delegate->setPositon(position);
    // 改变位置的时候，需要切换编辑器，以适应正确的位置
    m_trayView->onUpdateEditorView();
    updateLayout(position);
    onUpdateComponentSize();
}

void DockTrayWindow::setDisplayMode(const Dock::DisplayMode &displayMode)
{
    m_displayMode = displayMode;
    moveToolPlugin();
    updateToolWidget();
    // 如果当前模式为高效模式，则设置当前的trayView为其计算位置的参照
    if (displayMode == Dock::DisplayMode::Efficient) {
        ExpandIconWidget::popupTrayView()->setReferGridView(m_trayView);
        // TODO: reuse QuickPluginWindow, SystemPluginWindow
        auto stretch = m_mainBoxLayout->takeAt(m_mainBoxLayout->count()-1);
        m_mainBoxLayout->addWidget(m_trayView);
        m_mainBoxLayout->addItem(stretch);
    } else {
        m_mainBoxLayout->removeWidget(m_trayView);
    }
        
}

QSize DockTrayWindow::suitableSize(const Dock::Position &position, const int &, const double &) const
{
    if (position == Dock::Position::Left || position == Dock::Position::Right) {
        // 左右的尺寸
        int height = m_showDesktopWidget->height()
                + m_toolFrontSpaceWidget->height()
                + m_toolWidget->height()
                + m_toolBackSpaceWidget->height()
                + m_toolLineLabel->height()
                + m_dateTimeSpaceWidget->height()
                + m_dateTimeWidget->suitableSize(position).height()
                + m_systemPuginWidget->suitableSize(position).height()
                + m_quickIconWidget->suitableSize(position).height()
                + m_trayView->suitableSize(position).height();

        return QSize(-1, height);
    }
    // 上下的尺寸
    int width = m_showDesktopWidget->width()
            + m_toolFrontSpaceWidget->width()
            + m_toolWidget->width()
            + m_toolBackSpaceWidget->width()
            + m_toolLineLabel->width()
            + m_dateTimeSpaceWidget->width()
            + m_dateTimeWidget->width()
            + m_systemPuginWidget->width()
            + m_quickIconWidget->width()
            + m_trayView->width();
    return QSize(width, -1);
}

QSize DockTrayWindow::suitableSize() const
{
    return suitableSize(m_position, 0, 0);
}

void DockTrayWindow::layoutWidget()
{
    resizeTool();
}

void DockTrayWindow::resizeEvent(QResizeEvent *event)
{
    // 当尺寸发生变化的时候，通知托盘区域刷新尺寸，让托盘图标始终保持居中显示
    Q_EMIT m_delegate->sizeHintChanged(m_model->index(0, 0));
    QWidget::resizeEvent(event);
    onUpdateComponentSize();
}

void DockTrayWindow::paintEvent(QPaintEvent *event)
{
    QWidget::paintEvent(event);
    if (!m_toolLineLabel->isVisible())
        return;

    QPainter painter(this);
    QColor color;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        color = Qt::black;
        painter.setOpacity(0.5);
    } else {
        color = Qt::white;
        painter.setOpacity(0.1);
    }

    painter.fillRect(m_toolLineLabel->geometry(), color);
}

bool DockTrayWindow::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == this || watched == m_toolWidget || watched == m_dateTimeWidget
            || watched == m_trayView) {
        switch (event->type()) {
        case QEvent::Drop: {
                QDropEvent *dropEvent = static_cast<QDropEvent *>(event);
                onDropIcon(dropEvent);
                break;
            }
        case QEvent::DragEnter: {
            QDragEnterEvent *dragEnterEvent = static_cast<QDragEnterEvent *>(event);
            dragEnterEvent->setDropAction(Qt::CopyAction);
            dragEnterEvent->accept();
            return true;
        }
        case QEvent::DragMove: {
            QDragMoveEvent *dragMoveEvent = static_cast<QDragMoveEvent *>(event);
            dragMoveEvent->setDropAction(Qt::CopyAction);
            dragMoveEvent->accept();
            return true;
        }
        case QEvent::DragLeave: {
            QDragLeaveEvent *dragLeaveEvent = static_cast<QDragLeaveEvent *>(event);
            dragLeaveEvent->accept();
            break;
        }
        default:
            break;
        }
    }

    return QWidget::eventFilter(watched, event);
}

/** 根据任务栏的位置来更新布局的方向
 * @brief DockTrayWindow::updateLayout
 * @param position
 */
void DockTrayWindow::updateLayout(const Dock::Position &position)
{
    switch (position) {
    case Dock::Position::Left:
    case Dock::Position::Right: {
        m_mainBoxLayout->setDirection(QBoxLayout::BottomToTop);
        m_toolLayout->setDirection(QBoxLayout::BottomToTop);
        break;
    }
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        m_mainBoxLayout->setDirection(QBoxLayout::RightToLeft);
        m_toolLayout->setDirection(QBoxLayout::RightToLeft);
        break;
    }
    }
}

void DockTrayWindow::resizeTool() const
{
    int toolSize = 0;
    int size = 0;
    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right)
        size = width();
    else
        size = height();

    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *toolWidget = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!toolWidget)
            continue;

        toolWidget->setFixedSize(size, size);
        toolSize += size;
    }

    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right)
        m_toolWidget->setFixedSize(QWIDGETSIZE_MAX, toolSize);
    else
        m_toolWidget->setFixedSize(toolSize, QWIDGETSIZE_MAX);
}

bool DockTrayWindow::pluginExists(PluginsItemInterface *itemInter) const
{
    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginItem = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginItem)
            continue;

        if (pluginItem->pluginItem() == itemInter)
            return true;
    }

    return false;
}

void DockTrayWindow::moveToolPlugin()
{
    for (int i = m_toolLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginWidget = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginWidget)
            continue;

        m_toolLayout->removeWidget(pluginWidget);
    }
    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 如果当前是高效模式，则将所有的工具插件移动到当前的工具区域
        QuickSettingController *quickController = QuickSettingController::instance();
        QList<PluginsItemInterface *> plugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tool);
        for (PluginsItemInterface *pluginInter : plugins) {
            PluginsItem *pluginWidget = quickController->pluginItemWidget(pluginInter);
            m_toolLayout->addWidget(pluginWidget);
        }
    }
}

void DockTrayWindow::updateToolWidget()
{
    m_toolWidget->setVisible(m_toolLayout->count() > 0);
    m_toolLineLabel->setVisible(m_toolLayout->count() > 0);
    m_toolFrontSpaceWidget->setVisible(m_toolLayout->count() > 0);
    m_toolBackSpaceWidget->setVisible(m_toolLayout->count() > 0);
    m_dateTimeSpaceWidget->setVisible(m_toolLayout->count() > 0);
}

void DockTrayWindow::initUi()
{
    m_toolLayout->setContentsMargins(0, 0, 0, 0);
    m_toolLayout->setSpacing(0);

    m_systemPuginWidget->setDisplayMode(Dock::DisplayMode::Efficient);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(0);
    m_mainBoxLayout->addWidget(m_showDesktopWidget);
    m_mainBoxLayout->addWidget(m_toolFrontSpaceWidget);
    m_mainBoxLayout->addWidget(m_toolWidget);
    m_mainBoxLayout->addWidget(m_toolBackSpaceWidget);
    m_mainBoxLayout->addWidget(m_toolLineLabel);
    m_mainBoxLayout->addWidget(m_dateTimeSpaceWidget);
    m_mainBoxLayout->addWidget(m_dateTimeWidget);
    m_mainBoxLayout->addWidget(m_systemPuginWidget);
    m_mainBoxLayout->addWidget(m_quickIconWidget);
    m_mainBoxLayout->addWidget(m_trayView);
    m_mainBoxLayout->setAlignment(m_toolLineLabel, Qt::AlignCenter);

    m_toolLineLabel->setFixedSize(0, 0);

    m_mainBoxLayout->addStretch();
    updateToolWidget();
}

void DockTrayWindow::initConnection()
{
    connect(m_systemPuginWidget, &SystemPluginWindow::itemChanged, this, &DockTrayWindow::onUpdateComponentSize);
    connect(m_dateTimeWidget, &DateTimeDisplayer::requestUpdate, this, &DockTrayWindow::onUpdateComponentSize);
    connect(m_quickIconWidget, &QuickPluginWindow::itemCountChanged, this, &DockTrayWindow::onUpdateComponentSize);
    connect(m_systemPuginWidget, &SystemPluginWindow::requestDrop, this, &DockTrayWindow::onDropIcon);
    connect(m_model, &TrayModel::rowCountChanged, this, &DockTrayWindow::onUpdateComponentSize);
    connect(m_model, &TrayModel::rowCountChanged, m_trayView, &TrayGridView::onUpdateEditorView);
    connect(m_model, &TrayModel::requestRefreshEditor, m_trayView, &TrayGridView::onUpdateEditorView);
    connect(m_trayView, &TrayGridView::dragFinished, this, [ this ] {
        // 如果拖拽结束，则隐藏托盘
        Q_EMIT m_delegate->requestDrag(false);
    });

    connect(m_trayView, &TrayGridView::dragLeaved, m_delegate, [ this ] {
        Q_EMIT m_delegate->requestDrag(true);
    });
    connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ this ] (PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        switch (pluginAttr) {
        case QuickSettingController::PluginAttribute::Tool:
            // 下方只处理回收站等插件
            onItemAdded(itemInter);
            break;
        default:
            break;
        }
    });

    connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, &DockTrayWindow::onItemRemove);
}

void DockTrayWindow::initAttribute()
{
    setAcceptDrops(true);
    setMouseTracking(true);

    m_trayView->setModel(m_model);
    m_trayView->setItemDelegate(m_delegate);
    m_trayView->setDragDistance(2);
    m_trayView->setDragEnabled(true);

    installEventFilter(this);
    m_toolWidget->installEventFilter(this);
    m_dateTimeWidget->installEventFilter(this);
    m_systemPuginWidget->installEventFilter(this);
    m_quickIconWidget->installEventFilter(this);
    m_trayView->installEventFilter(this);
}

void DockTrayWindow::onUpdateComponentSize()
{
    switch (m_position) {
    case Dock::Position::Left:
    case Dock::Position::Right:
        m_toolLineLabel->setFixedSize(width() * 0.6, SPLITERSIZE);
        m_showDesktopWidget->setFixedSize(QWIDGETSIZE_MAX, FRONTSPACING);
        m_dateTimeWidget->setFixedSize(QWIDGETSIZE_MAX, m_dateTimeWidget->suitableSize().height());
        m_systemPuginWidget->setFixedSize(QWIDGETSIZE_MAX, m_systemPuginWidget->suitableSize().height());
        m_quickIconWidget->setFixedSize(QWIDGETSIZE_MAX, m_quickIconWidget->suitableSize().height());
        m_trayView->setFixedSize(QWIDGETSIZE_MAX, m_trayView->suitableSize().height());
        m_toolFrontSpaceWidget->setFixedSize(QWIDGETSIZE_MAX, SPLITESPACE);
        m_toolBackSpaceWidget->setFixedSize(QWIDGETSIZE_MAX, SPLITESPACE);
        m_dateTimeSpaceWidget->setFixedSize(QWIDGETSIZE_MAX, SPLITESPACE);
        break;
    case Dock::Position::Top:
    case Dock::Position::Bottom:
        m_toolLineLabel->setFixedSize(SPLITERSIZE, height() * 0.6);
        m_showDesktopWidget->setFixedSize(FRONTSPACING, QWIDGETSIZE_MAX);
        // FIXME: in some cases, m_dateTimeWidget QWIDGETSIZE_MAX get a huge height.
        m_dateTimeWidget->setFixedSize(m_dateTimeWidget->suitableSize().width(), qMin(QWIDGETSIZE_MAX, this->height()));
        m_systemPuginWidget->setFixedSize(m_systemPuginWidget->suitableSize().width(), QWIDGETSIZE_MAX);
        m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), QWIDGETSIZE_MAX);
        m_trayView->setFixedSize(m_trayView->suitableSize().width(), QWIDGETSIZE_MAX);
        m_toolFrontSpaceWidget->setFixedSize(SPLITESPACE, QWIDGETSIZE_MAX);
        m_toolBackSpaceWidget->setFixedSize(SPLITESPACE, QWIDGETSIZE_MAX);
        m_dateTimeSpaceWidget->setFixedSize(SPLITESPACE, QWIDGETSIZE_MAX);
        break;
    }
    Q_EMIT requestUpdate();
}

void DockTrayWindow::onItemAdded(PluginsItemInterface *itemInter)
{
    if (m_displayMode != Dock::DisplayMode::Efficient || pluginExists(itemInter))
        return;

    QuickSettingController *quickController = QuickSettingController::instance();
    PluginsItem *pluginItem = quickController->pluginItemWidget(itemInter);
    pluginItem->setVisible(true);

    m_toolLayout->addWidget(pluginItem);
    updateToolWidget();

    Q_EMIT requestUpdate();
}

void DockTrayWindow::onItemRemove(PluginsItemInterface *itemInter)
{
    for (int i = 0; i < m_toolLayout->count(); i++) {
        QLayoutItem *layoutItem = m_toolLayout->itemAt(i);
        if (!layoutItem)
            continue;

        PluginsItem *pluginItem = qobject_cast<PluginsItem *>(layoutItem->widget());
        if (!pluginItem || pluginItem->pluginItem() != itemInter)
            continue;

        m_toolLayout->removeWidget(pluginItem);
        updateToolWidget();

        Q_EMIT requestUpdate();
        break;
    }
}

void DockTrayWindow::onDropIcon(QDropEvent *dropEvent)
{
    if (!dropEvent || !dropEvent->mimeData() || dropEvent->source() == this)
        return;

    if (m_quickIconWidget->isQuickWindow(dropEvent->source())) {
        const QuickPluginMimeData *mimeData = qobject_cast<const QuickPluginMimeData *>(dropEvent->mimeData());
        if (!mimeData)
            return;

        PluginsItemInterface *pluginItem = static_cast<PluginsItemInterface *>(mimeData->pluginItemInterface());
        if (pluginItem)
            m_quickIconWidget->dragPlugin(pluginItem);
    } else if (qobject_cast<TrayGridView *>(dropEvent->source())) {
        // 将trayView中的dropEvent扩大到整个区域（this），这样便于随意拖动到这个区域都可以捕获。
        // m_trayView中有e->accept不会导致事件重复处理。
        m_trayView->handleDropEvent(dropEvent);
    }
}
