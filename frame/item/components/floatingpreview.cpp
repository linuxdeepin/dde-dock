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

#include "floatingpreview.h"
#include "appsnapshot.h"
#include "previewcontainer.h"

#include <DStyle>

#include <QGraphicsEffect>
#include <QPainter>
#include <QVBoxLayout>

#define BORDER_MARGIN 8
#define TITLE_MARGIN 20

FloatingPreview::FloatingPreview(QWidget *parent)
    : QWidget(parent)
    , m_closeBtn3D(new DImageButton)
    , m_titleBtn(new DPushButton)
{
    m_closeBtn3D->setFixedSize(24, 24);
    m_closeBtn3D->setNormalPic(":/icons/resources/close_round_normal.svg");
    m_closeBtn3D->setHoverPic(":/icons/resources/close_round_hover.svg");
    m_closeBtn3D->setPressPic(":/icons/resources/close_round_press.svg");

    m_titleBtn->setBackgroundRole(QPalette::Base);
    m_titleBtn->setForegroundRole(QPalette::Text);
    m_titleBtn->setFocusPolicy(Qt::NoFocus);
    m_titleBtn->setAttribute(Qt::WA_TransparentForMouseEvents);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addWidget(m_closeBtn3D);
    centralLayout->setAlignment(m_closeBtn3D, Qt::AlignRight | Qt::AlignTop);
    centralLayout->addWidget(m_titleBtn);
    centralLayout->setAlignment(m_titleBtn, Qt::AlignCenter | Qt::AlignBottom);
    centralLayout->addSpacing(TITLE_MARGIN);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_closeBtn3D, &DImageButton::clicked, this, &FloatingPreview::onCloseBtnClicked);
}

WId FloatingPreview::trackedWid() const
{
    Q_ASSERT(!m_tracked.isNull());

    return m_tracked->wid();
}

AppSnapshot *FloatingPreview::trackedWindow()
{
    return m_tracked;
}

void FloatingPreview::trackWindow(AppSnapshot *const snap)
{
    if (!m_tracked.isNull())
        m_tracked->removeEventFilter(this);

    snap->installEventFilter(this);
    m_tracked = snap;

    m_closeBtn3D->setVisible(m_tracked->closeAble());
    m_titleBtn->setText(m_tracked->title());

    QFontMetrics fm(m_titleBtn->font());
    int textWidth = fm.width(m_tracked->title()) + 10;
    int titleWidth = width() - (TITLE_MARGIN * 2  + BORDER_MARGIN);

    if (textWidth  < titleWidth) {
        m_titleBtn->setFixedWidth(textWidth);
    } else {
        m_titleBtn->setFixedWidth(titleWidth);
    }

    QTimer::singleShot(0, this, [ = ] {
        setGeometry(snap->geometry());
    });
}

void FloatingPreview::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    if (m_tracked.isNull())
        return;

    const QImage &snapshot = m_tracked->snapshot();
    const QRectF &snapshot_geometry = m_tracked->snapshotGeometry();

    if (snapshot.isNull())
        return;

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    const QRectF r = rect().marginsRemoved(QMargins(BORDER_MARGIN, BORDER_MARGIN, BORDER_MARGIN, BORDER_MARGIN));
    const auto ratio = devicePixelRatioF();

    const qreal offset_x = width() / 2.0 - snapshot_geometry.width() / ratio / 2 - snapshot_geometry.left();
    const qreal offset_y = height() / 2.0 - snapshot_geometry.height() / ratio / 2 - snapshot_geometry.top();

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);

    // 预览图
    QBrush brush;
    brush.setTextureImage(snapshot);
    painter.setBrush(brush);
    painter.setPen(Qt::NoPen);
    painter.translate(QPoint(offset_x, offset_y));
    painter.drawRoundedRect(snapshot_geometry, radius, radius);
    painter.translate(QPoint(-offset_x, -offset_y));

    // 选中外框
    QPen pen;
    pen.setColor(palette().highlight().color());
    pen.setWidth(dstyle.pixelMetric(DStyle::PM_FocusBorderWidth));
    painter.setPen(pen);
    painter.setBrush(Qt::NoBrush);
    painter.drawRoundedRect(r, radius, radius);
}

void FloatingPreview::mouseReleaseEvent(QMouseEvent *e)
{
    QWidget::mouseReleaseEvent(e);

    if (m_tracked) {
        emit m_tracked->clicked(m_tracked->wid());
    }
}

bool FloatingPreview::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_tracked && event->type() == QEvent::Destroy)
        hide();

    return QWidget::eventFilter(watched, event);
}

void FloatingPreview::hideEvent(QHideEvent *event)
{
    if (m_tracked) {
        m_tracked->setContentsMargins(0, 0, 0, 0);
    }

    QWidget::hideEvent(event);
}

void FloatingPreview::onCloseBtnClicked()
{
    Q_ASSERT(!m_tracked.isNull());

    m_tracked->closeWindow();
}
