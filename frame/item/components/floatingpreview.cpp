/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "floatingpreview.h"
#include "appsnapshot.h"
#include "previewcontainer.h"

#include <QPainter>
#include <QVBoxLayout>

FloatingPreview::FloatingPreview(QWidget *parent)
    : QWidget(parent),

      m_closeBtn(new DImageButton)
{
    m_closeBtn->setFixedSize(24, 24);
    m_closeBtn->setNormalPic(":/icons/resources/close_round_normal.svg");
    m_closeBtn->setHoverPic(":/icons/resources/close_round_hover.svg");
    m_closeBtn->setPressPic(":/icons/resources/close_round_press.svg");

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addWidget(m_closeBtn);
    centralLayout->setAlignment(m_closeBtn, Qt::AlignRight | Qt::AlignTop);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_closeBtn, &DImageButton::clicked, this, &FloatingPreview::onCloseBtnClicked);
}

WId FloatingPreview::trackedWid() const
{
    Q_ASSERT(!m_tracked.isNull());

    return m_tracked->wid();
}

void FloatingPreview::trackWindow(AppSnapshot * const snap)
{
    if (!m_tracked.isNull())
        m_tracked->removeEventFilter(this);
    snap->installEventFilter(this);
    m_tracked = snap;

    const QRect r = rect();
    const QRect sr = snap->geometry();
    const QPoint offset = sr.center() - r.center();

    emit requestMove(offset);
}

void FloatingPreview::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    if (m_tracked.isNull())
        return;

    const QImage &snapshot = m_tracked->snapshot();
    if (snapshot.isNull())
        return;

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    const QRect r = rect().marginsRemoved(QMargins(8, 8, 8, 8));
    const auto ratio = devicePixelRatioF();

    QImage im = snapshot.scaled(r.size() * ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    im.setDevicePixelRatio(ratio);

    const QRect ir = im.rect();
    const int offset_x = r.x() + r.width() / 2 - ir.width() / ratio / 2;
    const int offset_y = r.y() + r.height() / 2 - ir.height() / ratio / 2;
    const int radius = 4;

    // draw background
    painter.setPen(Qt::NoPen);
    painter.setBrush(QColor(255, 255, 255, 255 * 0.3));
    painter.drawRoundedRect(r, radius, radius);

    // draw preview image
    painter.drawImage(offset_x, offset_y, im);

    // bottom black background
    QRect bgr = r;
    bgr.setTop(bgr.bottom() - 25);

    QRect bgre = bgr;
    bgre.setTop(bgr.top() - radius);

    painter.save();
    painter.setClipRect(bgr);
    painter.setPen(Qt::NoPen);
    painter.setBrush(QColor(0, 0, 0, 255 * 0.3));
    painter.drawRoundedRect(bgre, radius, radius);
    painter.restore();

    // bottom title
    painter.setPen(Qt::white);
    painter.drawText(bgr, Qt::AlignCenter, m_tracked->title());
}

void FloatingPreview::mouseReleaseEvent(QMouseEvent *e)
{
    QWidget::mouseReleaseEvent(e);

    emit m_tracked->clicked(m_tracked->wid());
}

bool FloatingPreview::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_tracked && event->type() == QEvent::Destroy)
        hide();

    return QWidget::eventFilter(watched, event);
}

void FloatingPreview::onCloseBtnClicked()
{
    Q_ASSERT(!m_tracked.isNull());

    m_tracked->closeWindow();
}
