// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "previewcontainer.h"
#include "imageutil.h"
#include "utils.h"

#include <QDesktopWidget>
#include <QScreen>
#include <QApplication>
#include <QDragEnterEvent>
#include <QDesktopWidget>
#include <QCursor>
#include <QGSettings>
#include <QScrollBar>
#include <QScroller>

#define SPACING           0
#define MARGIN            0

PreviewContainer::PreviewContainer(QWidget *parent)
    : QWidget(parent)
    , m_needActivate(false)
    , m_canPreview(true)
    , m_dockSize(40)
    , m_floatingPreview(new FloatingPreview(this))
    , m_preparePreviewTimer(new QTimer(this))
    , m_mouseLeaveTimer(new QTimer(this))
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_titleMode(HoverShow)
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setSpacing(SPACING);
    m_windowListLayout->setContentsMargins(MARGIN, MARGIN, MARGIN, MARGIN);

    m_windowListWidget = new QWidget(this);
    m_windowListWidget->setAccessibleName("centralLayoutWidget");
    m_windowListWidget->setLayout(m_windowListLayout);
    m_windowListWidget->setMinimumSize(SNAP_WIDTH, SNAP_HEIGHT);
    m_windowListWidget->installEventFilter(this);
    m_windowListWidget->setAttribute(Qt::WA_TranslucentBackground);
    m_windowListWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Fixed);

    m_scrollArea = new QScrollArea(this);
    m_scrollArea->setAccessibleName("Preview_scrollArea");
    m_scrollArea->setWidgetResizable(false);
    m_scrollArea->setFrameStyle(QFrame::NoFrame);
    m_scrollArea->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_scrollArea->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_scrollArea->setContentsMargins(0, 0, 0, 0);
    m_scrollArea->setBackgroundRole(QPalette::Base);
    m_scrollArea->setWidget(m_windowListWidget);
    m_scrollArea->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);

    QVBoxLayout * mainLayout = new QVBoxLayout;
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);
    mainLayout->addWidget(m_scrollArea);
    setLayout(mainLayout);

    // 设置触摸屏手势识别
    QScroller::grabGesture(m_scrollArea->viewport(), QScroller::LeftMouseButtonGesture);
    QScroller *scroller = QScroller::scroller(m_scrollArea->viewport());
    m_sp.setScrollMetric(QScrollerProperties::VerticalOvershootPolicy, QScrollerProperties::OvershootWhenScrollable);
    m_sp.setScrollMetric(QScrollerProperties::HorizontalOvershootPolicy, QScrollerProperties::OvershootWhenScrollable);
    scroller->setScrollerProperties(m_sp);

    m_mouseLeaveTimer->setSingleShot(true);
    m_mouseLeaveTimer->setInterval(300);

    m_preparePreviewTimer->setSingleShot(true);
    m_preparePreviewTimer->setInterval(300);

    m_floatingPreview->installEventFilter(this);
    m_floatingPreview->setVisible(false);

    m_waitForShowPreviewTimer = new QTimer(this);
    m_waitForShowPreviewTimer->setSingleShot(true);
    m_waitForShowPreviewTimer->setInterval(200);

    setAcceptDrops(true);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_mouseLeaveTimer, &QTimer::timeout, this, &PreviewContainer::checkMouseLeave, Qt::QueuedConnection);
    connect(m_waitForShowPreviewTimer, &QTimer::timeout, this, &PreviewContainer::previewFloating);

    // 预览界面在滚动时，暂时不允许预览，避免在滚动时频繁切换预览造成闪屏
    connect(m_scrollArea->horizontalScrollBar(), &QScrollBar::valueChanged, this, [ = ]{
        m_floatingPreview->setVisible(false);
        m_canPreview = false;
        if (m_preparePreviewTimer->isActive()) {
            m_preparePreviewTimer->stop();
        }
        m_preparePreviewTimer->start();
    });
    connect(m_scrollArea->verticalScrollBar(), &QScrollBar::valueChanged, this, [ = ]{
        m_floatingPreview->setVisible(false);
        m_canPreview = false;
        if (m_preparePreviewTimer->isActive()) {
            m_preparePreviewTimer->stop();
        }
        m_preparePreviewTimer->start();
    });

    // 停止滚动后允许预览
    connect(m_preparePreviewTimer, &QTimer::timeout, this, [ = ]{
        m_canPreview = true;
        if (m_wmHelper->hasComposite()) {
            m_waitForShowPreviewTimer->start();
        }
    });
}

void PreviewContainer::setWindowInfos(const WindowInfoMap &infos, const WindowList &allowClose)
{
    // check removed window
    for (auto it(m_snapshots.begin()); it != m_snapshots.end();) {
        if (!infos.contains(it.key())) {
            m_windowListLayout->removeWidget(it.value());
            it.value()->deleteLater();
            it = m_snapshots.erase(it);
        } else {
            ++it;
        }
    }

    for (auto it(infos.cbegin()); it != infos.cend(); ++it) {
        const WId key = it.key();
        if (!m_snapshots.contains(key))
            appendSnapWidget(key);
        m_snapshots[key]->setWindowInfo(it.value());
        m_snapshots[key]->setCloseAble(allowClose.contains(key));
    }

    if (m_snapshots.isEmpty()) {
        emit requestCancelPreviewWindow();
        emit requestHidePopup();
    }

    adjustSize(m_wmHelper->hasComposite());
}

void PreviewContainer::setTitleDisplayMode(int mode)
{
    m_titleMode = static_cast<TitleDisplayMode>(mode);

    if (!m_wmHelper->hasComposite())
        return;

    m_floatingPreview->setFloatingTitleVisible(m_titleMode == HoverShow);

    for (AppSnapshot *snap : m_snapshots) {
        snap->setTitleVisible(m_titleMode == AlwaysShow);
    }
}

void PreviewContainer::updateLayoutDirection(const Dock::Position dockPos)
{
    m_dockPos = dockPos;
    if (m_wmHelper->hasComposite() && (dockPos == Dock::Top || dockPos == Dock::Bottom)) {
        m_windowListLayout->setDirection(QBoxLayout::LeftToRight);

        m_scrollArea->setHorizontalScrollBarPolicy(Qt::ScrollBarAsNeeded);
        m_scrollArea->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);

        m_sp.setScrollMetric(QScrollerProperties::VerticalOvershootPolicy, QScrollerProperties::OvershootAlwaysOff);
        m_sp.setScrollMetric(QScrollerProperties::HorizontalOvershootPolicy, QScrollerProperties::OvershootWhenScrollable);
    } else {
        m_scrollArea->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
        m_scrollArea->setVerticalScrollBarPolicy(Qt::ScrollBarAsNeeded);
        m_windowListLayout->setDirection(QBoxLayout::TopToBottom);

        m_sp.setScrollMetric(QScrollerProperties::VerticalOvershootPolicy, QScrollerProperties::OvershootWhenScrollable);
        m_sp.setScrollMetric(QScrollerProperties::HorizontalOvershootPolicy, QScrollerProperties::OvershootAlwaysOff);
    }

    adjustSize(m_wmHelper->hasComposite());
}

void PreviewContainer::updateDockSize(const int size)
{
    m_dockSize = size;
}

void PreviewContainer::checkMouseLeave()
{
    const bool hover = underMouse();

    if (hover)
        return;

    m_floatingPreview->setVisible(false);

    if (m_wmHelper->hasComposite()) {
        if (m_needActivate) {
            m_needActivate = false;
            emit requestActivateWindow(m_floatingPreview->trackedWid());
        } else {
            Q_EMIT requestHidePopup();
            emit requestCancelPreviewWindow();
        }
    }

    emit requestHidePopup();
}

void PreviewContainer::prepareHide()
{
    m_mouseLeaveTimer->start();
}

void PreviewContainer::adjustSize(bool composite)
{
    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    const int count = m_snapshots.size();

    const int screenWidth = QDesktopWidget().screenGeometry(this).width();
    const int screenHeight = QDesktopWidget().screenGeometry(this).height();

    //先根据屏幕宽高计算出能预览的最大数量,然后根据数量计算界面宽高
    if (composite) {
        if (horizontal) {
            const int h = SNAP_HEIGHT + MARGIN * 2;
            const int w = SNAP_WIDTH * count + MARGIN * 2 + SPACING * (count - 1);

            m_windowListWidget->setFixedSize(w, h);
            setFixedSize(qMin(w, screenWidth), h);
        } else {
            const int h = SNAP_HEIGHT * count + MARGIN * 2 + SPACING * (count - 1);
            const int w = SNAP_WIDTH + MARGIN * 2;

            m_windowListWidget->setFixedSize(w, h);
            setFixedSize(w, qMin(h, screenHeight));
        }
    } else if (m_windowListLayout->count()) {
        // 2D
        const int h = SNAP_HEIGHT_WITHOUT_COMPOSITE * count + MARGIN * 2 + SPACING * (count - 1);

        int titleWidth = 0;
        for (int i = 0; i < m_windowListLayout->count(); i++) {
            auto snapshotWidget = static_cast<AppSnapshot *>(m_windowListLayout->itemAt(i)->widget());
            titleWidth = qMax(titleWidth,  snapshotWidget->minimumWidth());
        }

        m_windowListWidget->setFixedSize(titleWidth, h);

        if (m_dockPos == Dock::Top || m_dockPos == Dock::Bottom) {
            // 滚动区域高度 = 屏幕高度 - 任务栏高度 - 箭头高度
            setFixedSize(titleWidth, qMin(h, screenHeight - m_dockSize - 20));
        } else {
            // 滚动区域高度 = 屏幕高度
            setFixedSize(titleWidth, qMin(h, screenHeight));
        }
    }
}

void PreviewContainer::appendSnapWidget(const WId wid)
{
    //创建预览界面,默认不显示,等计算出显示数量后再加入布局并显示
    AppSnapshot *snap = new AppSnapshot(wid);

    connect(snap, &AppSnapshot::clicked, this, &PreviewContainer::onSnapshotClicked, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::entered, this, &PreviewContainer::previewEntered, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::requestCheckWindow, this, &PreviewContainer::requestCheckWindows, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::requestCloseAppSnapshot, this, &PreviewContainer::onRequestCloseAppSnapshot);

    m_snapshots.insert(wid, snap);
    m_windowListLayout->addWidget(snap);
}

void PreviewContainer::enterEvent(QEvent *e)
{
    if (Utils::IS_WAYLAND_DISPLAY) {
        Utils::updateCursor(this);
    }

    QWidget::enterEvent(e);

    m_needActivate = false;
    m_mouseLeaveTimer->stop();

    if (m_wmHelper->hasComposite()) {
        m_waitForShowPreviewTimer->start();
    }
}

void PreviewContainer::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_mouseLeaveTimer->start();
    m_waitForShowPreviewTimer->stop();
}

void PreviewContainer::dragEnterEvent(QDragEnterEvent *e)
{
    if (!m_wmHelper->hasComposite())
        return;

    e->accept();

    m_needActivate = false;
    m_mouseLeaveTimer->stop();
}

void PreviewContainer::dragLeaveEvent(QDragLeaveEvent *e)
{
    e->ignore();

    m_needActivate = true;
    m_mouseLeaveTimer->start();
}

bool PreviewContainer::eventFilter(QObject *watcher, QEvent *event)
{
    // 将鼠标滚轮事件转换成水平滚动
    if (watcher == m_windowListWidget && event->type() == QEvent::Wheel) {
        QWheelEvent *wheelEvent = static_cast<QWheelEvent *>(event);

        if (m_windowListLayout->direction() == Qt::LeftToRight) {
            const int delta = wheelEvent->delta();
            const int currValue = m_scrollArea->horizontalScrollBar()->value();

            if (currValue - delta <= 0) {
                m_scrollArea->horizontalScrollBar()->setValue(0);
            } else if (currValue - delta >= m_scrollArea->horizontalScrollBar()->maximum()) {
                m_scrollArea->horizontalScrollBar()->setValue(m_scrollArea->horizontalScrollBar()->maximum());
            } else {
                m_scrollArea->horizontalScrollBar()->setValue(currValue - delta);
            }

            return true;
        }
    }

    // 在m_floatingPreview界面显示时，需要响应滚轮事件
    if (watcher == m_floatingPreview && event->type() == QEvent::Wheel) {
        QWheelEvent *wheelEvent = static_cast<QWheelEvent *>(event);

        if (m_windowListLayout->direction() == Qt::LeftToRight) {
            const int delta = wheelEvent->delta();
            const int currValue = m_scrollArea->horizontalScrollBar()->value();

            if (currValue - delta <= 0) {
                m_scrollArea->horizontalScrollBar()->setValue(0);
            } else if (currValue - delta >= m_scrollArea->horizontalScrollBar()->maximum()) {
                m_scrollArea->horizontalScrollBar()->setValue(m_scrollArea->horizontalScrollBar()->maximum());
            } else {
                m_scrollArea->horizontalScrollBar()->setValue(currValue - delta);
            }

            return true;
        } else {
            const int delta = wheelEvent->delta();
            const int currValue = m_scrollArea->verticalScrollBar()->value();

            if (currValue - delta <= 0) {
                m_scrollArea->verticalScrollBar()->setValue(0);
            } else if (currValue - delta >= m_scrollArea->verticalScrollBar()->maximum()) {
                m_scrollArea->verticalScrollBar()->setValue(m_scrollArea->verticalScrollBar()->maximum());
            } else {
                m_scrollArea->verticalScrollBar()->setValue(currValue - delta);
            }

            return true;
        }
    }

    return QWidget::eventFilter(watcher, event);
}

void PreviewContainer::onSnapshotClicked(const WId wid)
{
    if (Utils::IS_WAYLAND_DISPLAY) {
        /* BUGFIX-159303 本地发现该问题仅在wayland下出现，问题已与窗管对接沟通过，根据窗管建议，上层进
        规避，避免因底层强行修改引入新的问题 */
        Q_EMIT requestCancelPreviewWindow();
    }

    Q_EMIT requestActivateWindow(wid);
    m_needActivate = true;
    m_waitForShowPreviewTimer->stop();
    requestHidePopup();
}

void PreviewContainer::previewEntered(const WId wid)
{
    if (!m_wmHelper->hasComposite())
        return;

    AppSnapshot *snap = static_cast<AppSnapshot *>(sender());
    if (!snap) {
        return;
    }

    AppSnapshot *preSnap = m_floatingPreview->trackedWindow();
    if (preSnap && preSnap != snap) {
        preSnap->setWindowState();
    }

    m_currentWId = wid;
    m_floatingPreview->trackWindow(snap);

    if (m_waitForShowPreviewTimer->isActive()) {
        return;
    }

    previewFloating();
}

void PreviewContainer::previewFloating()
{
    if (!m_waitForShowPreviewTimer->isActive() && m_canPreview) {
        // 将单个预览界面在滚动区域的坐标转换到当前界面坐标上
        if (m_floatingPreview->trackedWindow()) {
            const QRect snapGeometry = m_floatingPreview->trackedWindow()->geometry();
            const QPoint topLeft = m_windowListWidget->mapTo(m_scrollArea, snapGeometry.topLeft());
            const QPoint bottomRight = m_windowListWidget->mapTo(m_scrollArea, snapGeometry.bottomRight());
            m_floatingPreview->setGeometry(QRect(topLeft, bottomRight));
        }

        m_floatingPreview->setVisible(true);
        m_floatingPreview->raise();

        requestPreviewWindow(m_currentWId);
    }
    return;
}

void PreviewContainer::onRequestCloseAppSnapshot()
{
    if (!m_wmHelper->hasComposite())
        return ;

    if (m_snapshots.keys().isEmpty()) {
        Q_EMIT requestHidePopup();
        Q_EMIT requestCancelPreviewWindow();
    }
}
