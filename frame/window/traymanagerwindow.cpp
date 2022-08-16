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
#include "traymanagerwindow.h"
#include "quickpluginwindow.h"
#include "tray_gridview.h"
#include "tray_delegate.h"
#include "tray_model.h"
#include "constants.h"
#include "quicksettingcontainer.h"
#include "systempluginwindow.h"
#include "datetimedisplayer.h"

#include <DGuiApplicationHelper>
#include <DRegionMonitor>

#include <QDropEvent>
#include <QBoxLayout>
#include <QLabel>
#include <QMimeData>
#include <QDBusConnection>
#include <QPainter>
#include <QPainterPath>

#define MAXFIXEDSIZE 999999
#define CRITLCALHEIGHT 42
#define CONTENTSPACE 7
// 高度小于等于这个值的时候，间距最小值
#define MINHIGHT 46
// 最小值与最大值的差值
#define MINSPACE 4
// 当前高度与最小高度差值大于这个值的时候，间距使用最大值
#define MAXDIFF 3

TrayManagerWindow::TrayManagerWindow(QWidget *parent)
    : QWidget(parent)
    , m_appPluginDatetimeWidget(new QWidget(this))
    , m_systemPluginWidget(new SystemPluginWindow(this))
    , m_appPluginWidget(new QWidget(m_appPluginDatetimeWidget))
    , m_quickIconWidget(new QuickPluginWindow(m_appPluginWidget))
    , m_dateTimeWidget(new DateTimeDisplayer(m_appPluginDatetimeWidget))
    , m_appPluginLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
    , m_appDatetimeLayout(new QBoxLayout(QBoxLayout::Direction::TopToBottom, this))
    , m_mainLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
    , m_trayView(new TrayGridView(this))
    , m_model(new TrayModel(m_trayView, false, true))
    , m_delegate(new TrayDelegate(m_trayView, m_trayView))
    , m_postion(Dock::Position::Bottom)
    , m_splitLine(new QLabel(m_appPluginDatetimeWidget))
{
    initUi();
    initConnection();

    setAcceptDrops(true);
    setMouseTracking(true);
}

TrayManagerWindow::~TrayManagerWindow()
{
}

void TrayManagerWindow::updateLayout()
{
    bool showSingle = true;
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
        showSingle = (topLevelWidget()->height() <= CRITLCALHEIGHT);

    QBoxLayout::Direction lastDirection = m_appDatetimeLayout->direction();

    if (showSingle)
        resetSingleDirection();
    else
        resetMultiDirection();

    resetChildWidgetSize();
    // 当尺寸发生变化的时候，通知托盘区域刷新尺寸，让托盘图标始终保持居中显示
    Q_EMIT m_delegate->sizeHintChanged(m_model->index(0, 0));

    // 当插件区域从单行变成两行或者两行变成单行的时候，发送该信号，通知外部重新调整区域大小
    if (lastDirection != m_appDatetimeLayout->direction())
        Q_EMIT requestUpdate();
}

void TrayManagerWindow::setPositon(Dock::Position position)
{
    if (m_postion == position)
        return;

    m_postion = position;

    if (position == Dock::Position::Top || position == Dock::Position::Bottom)
        m_trayView->setOrientation(QListView::Flow::LeftToRight, false);
    else
        m_trayView->setOrientation(QListView::Flow::TopToBottom, false);

    QModelIndex index = m_model->index(0, 0);
    m_trayView->closePersistentEditor(index);
    TrayDelegate *delegate = static_cast<TrayDelegate *>(m_trayView->itemDelegate());
    delegate->setPositon(position);
    m_trayView->openPersistentEditor(index);

    m_trayView->setPosition(position);
    m_quickIconWidget->setPositon(position);
    m_dateTimeWidget->setPositon(position);
    m_systemPluginWidget->setPositon(position);
}

int TrayManagerWindow::appDatetimeSize()
{
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom) {
        // 如果是一行
        if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::LeftToRight) {
            return m_trayView->suitableSize().width() + m_quickIconWidget->suitableSize().width()
                    + m_dateTimeWidget->suitableSize().width() + m_appDatetimeLayout->spacing();
        }
        //如果是两行
        int topWidth = m_trayView->suitableSize().width() + m_appDatetimeLayout->spacing() + m_quickIconWidget->width();
        int bottomWidth = m_dateTimeWidget->suitableSize().width();
        return qMax(topWidth, bottomWidth);
    }

    int trayHeight = m_trayView->suitableSize().height();
    QMargins m = m_appDatetimeLayout->contentsMargins();
    int traypluginHeight = trayHeight + m_quickIconWidget->suitableSize().height() + m.top() + m.bottom() + m_appPluginLayout->spacing();
    return traypluginHeight + m_appDatetimeLayout->spacing() + m_dateTimeWidget->suitableSize().height() + 10;
}

QSize TrayManagerWindow::suitableSize()
{
    QMargins m = m_mainLayout->contentsMargins();
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom) {
        return QSize(appDatetimeSize() + m_appDatetimeLayout->spacing() +
                     m_systemPluginWidget->suitableSize().width() + m_mainLayout->spacing() +
                     m.left() + m.right(), QWIDGETSIZE_MAX);
    }

    return QSize(QWIDGETSIZE_MAX, appDatetimeSize() + m_appDatetimeLayout->spacing() +
                 m_systemPluginWidget->suitableSize().height() + m_mainLayout->spacing() +
                 m.top() + m.bottom());
}

// 用于返回需要绘制的圆形区域
QPainterPath TrayManagerWindow::roundedPaths()
{
    int topLevelRadius = qApp->property("EffectBorderRadius").toInt();
    QMargins mainMargin = m_mainLayout->contentsMargins();
    int radius = topLevelRadius - mainMargin.top();
    QPainterPath path;
    if ((m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
        && (m_appDatetimeLayout->direction() == QBoxLayout::Direction::LeftToRight)) {
        // 如果是上下方向，且只有一行
        // 计算托盘和快捷插件区域
        QPoint pointPlugin(mainMargin.left(), mainMargin.top());
        QRect rctPlugin(QPoint(mainMargin.left(), mainMargin.top()), m_appPluginWidget->size());
        path.addRoundedRect(rctPlugin, radius, radius);

        // 计算日期时间区域
        QPoint pointDateTime = m_dateTimeWidget->pos();
        pointDateTime = m_dateTimeWidget->parentWidget()->mapTo(this, pointDateTime);
        QRect rctDatetime(pointDateTime, m_dateTimeWidget->size());
        path.addRoundedRect(rctDatetime, radius, radius);
        // 计算系统插件区域
        path.addRoundedRect(m_systemPluginWidget->geometry(), radius, radius);
    } else {
        // 添加系统插件区域的位置信息
        path.addRoundedRect(m_appPluginDatetimeWidget->geometry(), radius, radius);
        path.addRoundedRect(m_systemPluginWidget->geometry(), radius, radius);
    }

    return path;
}

void TrayManagerWindow::resizeEvent(QResizeEvent *event)
{
    Q_UNUSED(event);

    updateLayout();
}

void TrayManagerWindow::initUi()
{
    m_trayView->setModel(m_model);
    m_trayView->setItemDelegate(m_delegate);
    m_trayView->setDragDistance(2);

    m_splitLine->setFixedHeight(1);
    QPalette pal;
    QColor lineColor(Qt::black);
    lineColor.setAlpha(static_cast<int>(255 * 0.1));
    pal.setColor(QPalette::Background, lineColor);
    m_splitLine->setAutoFillBackground(true);
    m_splitLine->setPalette(pal);

    WinInfo info;
    info.type = TrayIconType::EXPANDICON;
    m_model->addRow(info);
    m_trayView->openPersistentEditor(m_model->index(0, 0));

    // 左侧的区域，包括应用托盘插件和下方的日期时间区域
    m_appPluginLayout->setContentsMargins(0, 0, 0, 0);
    m_appPluginLayout->setSpacing(0);
    m_appPluginWidget->setLayout(m_appPluginLayout);
    m_appPluginLayout->addWidget(m_trayView);
    m_appPluginLayout->addWidget(m_quickIconWidget);

    m_appPluginDatetimeWidget->setLayout(m_appDatetimeLayout);
    m_appDatetimeLayout->setContentsMargins(0, 0, 0, 0);
    m_appDatetimeLayout->setSpacing(0);
    m_appDatetimeLayout->addWidget(m_appPluginWidget);
    m_appDatetimeLayout->addWidget(m_splitLine);
    m_appDatetimeLayout->addWidget(m_dateTimeWidget);

    setLayout(m_mainLayout);
    // 通用情况下，设置边距和间距都为7
    m_mainLayout->setContentsMargins(CONTENTSPACE, CONTENTSPACE, CONTENTSPACE, CONTENTSPACE);
    m_mainLayout->setSpacing(CONTENTSPACE);
    m_mainLayout->addWidget(m_appPluginDatetimeWidget);
    m_mainLayout->addWidget(m_systemPluginWidget);
}

void TrayManagerWindow::initConnection()
{
    connect(m_trayView, &TrayGridView::requestRemove, m_model, &TrayModel::removeRow);
    connect(m_trayView, &TrayGridView::rowCountChanged, this, [ this ] {
        if (m_quickIconWidget->x() == 0) {
            // 在加载界面的时候，会出现快捷设置区域的图标和左侧的托盘图标挤在一起(具体原因未知)，此时需要延时50毫秒重新刷新界面来保证界面布局正常(临时解决方案)
            QTimer::singleShot(50, this, [ this ] {
                resetChildWidgetSize();
                Q_EMIT requestUpdate();
            });
        } else {
            resetChildWidgetSize();
            Q_EMIT requestUpdate();
        }
    });
    connect(m_quickIconWidget, &QuickPluginWindow::itemCountChanged, this, [ this ] {
        // 当插件数量发生变化的时候，需要调整尺寸
        m_quickIconWidget->setFixedSize(QWIDGETSIZE_MAX, QWIDGETSIZE_MAX);
        if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
            m_quickIconWidget->setFixedWidth(m_quickIconWidget->suitableSize().width());
        else
            m_quickIconWidget->setFixedHeight(m_quickIconWidget->suitableSize().height());

        Q_EMIT requestUpdate();
    });

    connect(m_systemPluginWidget, &SystemPluginWindow::itemChanged, this, [ this ] {
        // 当系统插件发生变化的时候，同样需要调整尺寸
        m_systemPluginWidget->setFixedSize(QWIDGETSIZE_MAX, QWIDGETSIZE_MAX);
        if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
            m_systemPluginWidget->setFixedWidth(m_systemPluginWidget->suitableSize().width());
        else
            m_systemPluginWidget->setFixedHeight(m_systemPluginWidget->suitableSize().height());

        Q_EMIT requestUpdate();
    });

    connect(m_delegate, &TrayDelegate::visibleChanged, this, [ this ](const QModelIndex &index, bool visible) {
        m_trayView->setRowHidden(index.row(), !visible);
        resetChildWidgetSize();
        Q_EMIT requestUpdate();
    });

    connect(m_trayView, &TrayGridView::dragLeaved, m_delegate, [ this ]{
        Q_EMIT m_delegate->requestDrag(true);
    });
    connect(m_trayView, &TrayGridView::dragEntered, m_delegate, [ this ]{
        Q_EMIT m_delegate->requestDrag(false);
    });
    connect(m_model, &TrayModel::requestUpdateWidget, this, [ this ](const QList<int> &idxs) {
        for (int i = 0; i < idxs.size(); i++) {
             int idx = idxs[i];
             if (idx < m_model->rowCount()) {
                 QModelIndex index = m_model->index(idx);
                 m_trayView->closePersistentEditor(index);
                 m_trayView->openPersistentEditor(index);
             }
        }
    });
    connect(m_dateTimeWidget, &DateTimeDisplayer::requestUpdate, this, &TrayManagerWindow::requestUpdate);

    m_trayView->installEventFilter(this);
    m_quickIconWidget->installEventFilter(this);
    installEventFilter(this);
    QMetaObject::invokeMethod(this, &TrayManagerWindow::requestUpdate, Qt::QueuedConnection);
}

void TrayManagerWindow::resetChildWidgetSize()
{
    int count = 0;
    for (int i = 0; i < m_model->rowCount(); i++) {
        if (!m_trayView->isRowHidden(i))
            count++;
    }

    switch (m_postion) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        int trayWidth = m_trayView->suitableSize().width();
        QMargins m = m_appPluginLayout->contentsMargins();
        if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::LeftToRight) {
            // 单行显示
            int trayHeight = m_appPluginDatetimeWidget->height() - m.top() - m.bottom();
            m_trayView->setFixedSize(trayWidth, trayHeight);
            m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), trayHeight);
            m_appPluginWidget->setFixedSize(trayWidth + m_quickIconWidget->suitableSize().width(), trayHeight);
            m_dateTimeWidget->setFixedSize(m_dateTimeWidget->suitableSize().width(), trayHeight);
            // 设置右侧的电源按钮的尺寸
            m_systemPluginWidget->setFixedSize(m_systemPluginWidget->suitableSize());
            m_mainLayout->setContentsMargins(4, 4, 4 ,4);
            m_mainLayout->setSpacing(4);
            m_appDatetimeLayout->setSpacing(4);
        } else {
            // 多行显示
            int trayHeight = m_appPluginDatetimeWidget->height() / 2 + 4 - m.top() - m.bottom();
            m_trayView->setFixedSize(trayWidth, trayHeight);
            m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), trayHeight);
            m_appPluginWidget->setFixedSize(trayWidth + m_quickIconWidget->suitableSize().width(), m_appPluginDatetimeWidget->height() / 2 + 4);
            // 因为是两行，所以对于时间控件的尺寸，只能设置最小值
            int dateTimeWidth = qMax(m_appPluginWidget->width(), m_dateTimeWidget->suitableSize().width());
            m_dateTimeWidget->setMinimumSize(dateTimeWidth, QWIDGETSIZE_MAX);
            m_systemPluginWidget->setFixedSize(m_systemPluginWidget->suitableSize());
            int contentSpace = qMin(MAXDIFF, qMax(height() - MINHIGHT, 0)) + MINSPACE;
            m_mainLayout->setContentsMargins(contentSpace, contentSpace, contentSpace, contentSpace);
            m_appDatetimeLayout->setSpacing(0);
            m_mainLayout->setSpacing(contentSpace);
        }
        QMargins margin = m_mainLayout->contentsMargins();
        int appDateHeight = height() - margin.top() - margin.bottom();
        m_appPluginDatetimeWidget->setFixedSize(appDatetimeSize(), appDateHeight);
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        int trayHeight = m_trayView->suitableSize().height();
        int quickAreaHeight = m_quickIconWidget->suitableSize().height();
        QMargins m = m_appPluginLayout->contentsMargins();
        // 左右方向始终只有一列
        int datetimeHeight = m_dateTimeWidget->suitableSize().height();
        int sizeWidth = m_appPluginDatetimeWidget->width() - m.left() - m.right();
        m_trayView->setFixedSize(sizeWidth, trayHeight);
        m_quickIconWidget->setFixedSize(sizeWidth, quickAreaHeight);
        m_dateTimeWidget->setFixedSize(sizeWidth, datetimeHeight);
        m_appPluginWidget->setFixedSize(sizeWidth, trayHeight + quickAreaHeight);
        m_systemPluginWidget->setFixedSize(m_systemPluginWidget->suitableSize());

        int contentSpace = qMin(MAXDIFF, qMax(width() - MINHIGHT, 0)) + MINSPACE;
        m_mainLayout->setContentsMargins(contentSpace, contentSpace, contentSpace, contentSpace);
        m_appDatetimeLayout->setSpacing(0);
        m_mainLayout->setSpacing(contentSpace);

        int appDateWidth = width() - (contentSpace * 2);
        m_appPluginDatetimeWidget->setFixedSize(appDateWidth, appDatetimeSize());
        break;
    }
    }
}

void TrayManagerWindow::resetSingleDirection()
{
    switch (m_postion) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        // 应用和时间在一行显示
        m_appDatetimeLayout->setSpacing(10);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_splitLine->hide();
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right:{
        m_appDatetimeLayout->setSpacing(0);
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_splitLine->show();
        break;
    }
    }
    m_dateTimeWidget->setOneRow(true);
}

void TrayManagerWindow::resetMultiDirection()
{
    m_appDatetimeLayout->setSpacing(0);
    switch (m_postion) {
    case Dock::Position::Top: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::BottomToTop);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_splitLine->show();
        m_dateTimeWidget->setOneRow(false);
        break;
    }
    case Dock::Position::Bottom: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_splitLine->show();
        m_dateTimeWidget->setOneRow(false);
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_splitLine->hide();
        m_dateTimeWidget->setOneRow(true);
        break;
    }
    }
}

void TrayManagerWindow::dragEnterEvent(QDragEnterEvent *e)
{
    e->setDropAction(Qt::CopyAction);
    e->accept();
}

void TrayManagerWindow::dragMoveEvent(QDragMoveEvent *e)
{
    e->setDropAction(Qt::CopyAction);
    e->accept();
}

void TrayManagerWindow::dropEvent(QDropEvent *e)
{
    if (!e || !e->mimeData() || e->source() == this)
        return;

    if (qobject_cast<QuickPluginWindow *>(e->source())) {
        const QuickPluginMimeData *mimeData = qobject_cast<const QuickPluginMimeData *>(e->mimeData());
        if (!mimeData)
            return;

        PluginsItemInterface *pluginItem = static_cast<PluginsItemInterface *>(mimeData->pluginItemInterface());
        if (pluginItem)
            m_quickIconWidget->dragPlugin(pluginItem);
    } else if (qobject_cast<TrayGridView *>(e->source())) {
        // 将trayView中的dropEvent扩大到整个区域（this），这样便于随意拖动到这个区域都可以捕获。
        // m_trayView中有e->accept不会导致事件重复处理。
        m_trayView->handleDropEvent(e);
    }
}

void TrayManagerWindow::dragLeaveEvent(QDragLeaveEvent *event)
{
    event->accept();
}

void TrayManagerWindow::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);
    QPainterPath path = roundedPaths();
    QPainter painter(this);
    painter.save();
    painter.setRenderHint(QPainter::Antialiasing);
    painter.setClipPath(path);
    painter.fillRect(rect().adjusted(1, 1, -1, -1), maskColor(102));
    painter.setPen(maskColor(110));
    painter.drawPath(path);
    painter.restore();
}

QColor TrayManagerWindow::maskColor(uint8_t alpha) const
{
    QColor color = DGuiApplicationHelper::standardPalette(DGuiApplicationHelper::instance()->themeType()).window().color();
    color.setAlpha(alpha);
    return color;
}
