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

#include "previewcontainer.h"

#include <QDesktopWidget>
#include <QScreen>
#include <QApplication>
#include <QDragEnterEvent>

#define SPACING           0
#define MARGIN            0
#define SNAP_HEIGHT_WITHOUT_COMPOSITE       30

PreviewContainer::PreviewContainer(QWidget *parent)
    : QWidget(parent),
      m_needActivate(false),

      m_floatingPreview(new FloatingPreview(this)),
      m_mouseLeaveTimer(new QTimer(this)),
      m_wmHelper(DWindowManagerHelper::instance())
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setSpacing(SPACING);
    m_windowListLayout->setContentsMargins(MARGIN, MARGIN, MARGIN, MARGIN);

    m_mouseLeaveTimer->setSingleShot(true);
    m_mouseLeaveTimer->setInterval(300);

    m_floatingPreview->setVisible(false);

    setAcceptDrops(true);
    setLayout(m_windowListLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_mouseLeaveTimer, &QTimer::timeout, this, &PreviewContainer::checkMouseLeave, Qt::QueuedConnection);
    connect(m_floatingPreview, &FloatingPreview::requestMove, this, &PreviewContainer::moveFloatingPreview);
}

void PreviewContainer::setWindowInfos(const WindowInfoMap &infos, const WindowList &allowClose)
{
    // check removed window
    for (auto it(m_snapshots.begin()); it != m_snapshots.end();)
    {
        if (!infos.contains(it.key()))
        {
            m_windowListLayout->removeWidget(it.value());
            it.value()->deleteLater();
            it = m_snapshots.erase(it);
        } else {
            ++it;
        }
    }

    for (auto it(infos.cbegin()); it != infos.cend(); ++it)
    {
        const WId key = it.key();
        if (!m_snapshots.contains(key))
            appendSnapWidget(key);
        m_snapshots[key]->setWindowInfo(it.value());
        m_snapshots[key]->setCloseAble(allowClose.contains(key));
    }

    if (m_snapshots.isEmpty())
        emit requestCancelAndHidePreview();

    adjustSize();
}

void PreviewContainer::updateSnapshots()
{
    for (AppSnapshot *snap : m_snapshots)
        snap->fetchSnapshot();
}

void PreviewContainer::updateLayoutDirection(const Dock::Position dockPos)
{
    if (m_wmHelper->hasComposite() && (dockPos == Dock::Top || dockPos == Dock::Bottom))
        m_windowListLayout->setDirection(QBoxLayout::LeftToRight);
    else
        m_windowListLayout->setDirection(QBoxLayout::TopToBottom);

    adjustSize();
}

void PreviewContainer::checkMouseLeave()
{
    const bool hover = underMouse();

    if (hover)
        return;

    emit requestCancelAndHidePreview();

    m_floatingPreview->setVisible(false);

    if (m_needActivate)
    {
        m_needActivate = false;
        emit requestActivateWindow(m_floatingPreview->trackedWid());
    }
}

void PreviewContainer::prepareHide()
{
    m_mouseLeaveTimer->start();
}

void PreviewContainer::adjustSize()
{
    const int count = m_snapshots.size();
    const bool composite = m_wmHelper->hasComposite();
    if (!composite)
    {
        const int h = SNAP_HEIGHT_WITHOUT_COMPOSITE * count + MARGIN * 2 + SPACING * (count - 1);
        setFixedSize(SNAP_WIDTH, h);
        return;
    }

    const QRect r = qApp->primaryScreen()->geometry();
    const int padding = 20;

    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    if (horizontal)
    {
        const int h = SNAP_HEIGHT + MARGIN * 2;
        const int w = SNAP_WIDTH * count + MARGIN * 2 + SPACING * (count - 1);

        setFixedHeight(h);
        setFixedWidth(std::min(w, r.width() - padding));
    } else {
        const int w = SNAP_WIDTH + MARGIN * 2;
        const int h = SNAP_HEIGHT * count + MARGIN * 2 + SPACING * (count - 1);

        setFixedWidth(w);
        setFixedHeight(std::min(h, r.height() - padding));
    }
}

void PreviewContainer::appendSnapWidget(const WId wid)
{
    AppSnapshot *snap = new AppSnapshot(wid);

    connect(snap, &AppSnapshot::clicked, this, &PreviewContainer::requestActivateWindow, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::clicked, this, &PreviewContainer::requestCancelAndHidePreview, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::entered, this, &PreviewContainer::previewEntered, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::requestCheckWindow, this, &PreviewContainer::requestCheckWindows);

    m_windowListLayout->addWidget(snap);

    m_snapshots.insert(wid, snap);

    // refresh if visible
    if (isVisible())
        snap->fetchSnapshot();
}

void PreviewContainer::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    m_needActivate = false;
    m_mouseLeaveTimer->stop();
}

void PreviewContainer::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_floatingPreview->setVisible(false);
    m_mouseLeaveTimer->start();
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
    m_mouseLeaveTimer->start(10);
}

void PreviewContainer::previewEntered(const WId wid)
{
    if (!m_wmHelper->hasComposite())
        return;

    AppSnapshot *snap = static_cast<AppSnapshot *>(sender());
    if (!snap) {
        return;
    }
    snap->setContentsMargins(100, 0, 100, 0);

    AppSnapshot *preSnap = m_floatingPreview->trackedWindow();
    if (preSnap && preSnap != snap) {
        preSnap->setContentsMargins(0, 0, 0, 0);
    }

    QTimer::singleShot(0, [=] {
        m_floatingPreview->trackWindow(snap);
        m_floatingPreview->setVisible(true);
        m_floatingPreview->raise();
    });

    emit requestPreviewWindow(wid);
}

void PreviewContainer::moveFloatingPreview(const QPoint &p)
{
    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    const QRect r = rect();

    if (horizontal)
    {
        if (p.x() < r.left())
            m_floatingPreview->move(MARGIN, p.y());
        else if (p.x() + m_floatingPreview->width() > r.right())
            m_floatingPreview->move(r.right() - m_floatingPreview->width() - MARGIN + 1, p.y());
        else
            m_floatingPreview->move(p);
    } else {
        if (p.y() < r.top())
            m_floatingPreview->move(p.x(), MARGIN);
        else if (p.y() + m_floatingPreview->height() > r.bottom())
            m_floatingPreview->move(p.x(), r.bottom() - m_floatingPreview->height() - MARGIN + 1);
        else
            m_floatingPreview->move(p);
    }
}
