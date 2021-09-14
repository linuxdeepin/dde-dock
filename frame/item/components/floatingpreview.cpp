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

FloatingPreview::FloatingPreview(QWidget *parent)
    : QWidget(parent)
    , m_closeBtn3D(new DIconButton(this))
    , m_titleBtn(new DPushButton(this))
{
    m_closeBtn3D->setObjectName("closebutton-3d");
    m_closeBtn3D->setFixedSize(24, 24);
    m_closeBtn3D->setIconSize(QSize(24, 24));
    m_closeBtn3D->setIcon(QIcon(":/icons/resources/close_round_normal.svg"));
    m_closeBtn3D->setFlat(true);
    m_closeBtn3D->installEventFilter(this);

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

    connect(m_closeBtn3D, &DIconButton::clicked, this, &FloatingPreview::onCloseBtnClicked);
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

void FloatingPreview::setFloatingTitleVisible(bool bVisible)
{
    m_titleBtn->setVisible(bVisible);
}

void FloatingPreview::trackWindow(AppSnapshot *const snap)
{
    if (!snap)
        return;

    if (!m_tracked.isNull())
        m_tracked->removeEventFilter(this);

    snap->installEventFilter(this);
    m_tracked = snap;

    m_closeBtn3D->setVisible(m_tracked->closeAble());

    // 显示此标题的前提条件：配置了标题跟随鼠标显示
    // 此对象是共用的，鼠标移动到哪个预览图，title就移动到哪里显示，所以他的text统一snap获取，不再重复计算显示长度
    m_titleBtn->setText(snap->appTitle());

    QTimer::singleShot(0, this, [ = ] {
        // 此处获取的snap->geometry()有可能是错误的，所以做个判断并且在resizeEvent中也做处理
        if(snap->width() == SNAP_WIDTH)
            setGeometry(snap->geometry());
    });
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

    const QRectF r = rect().marginsRemoved(QMargins(BORDER_MARGIN, BORDER_MARGIN, BORDER_MARGIN, BORDER_MARGIN));

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);

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
    if(watched == m_closeBtn3D) {
        if(watched == m_closeBtn3D && (event->type() == QEvent::HoverEnter || event->type() == QEvent::HoverMove)) {
            m_closeBtn3D->setIcon(QIcon(":/icons/resources/close_round_hover.svg"));
        }
        else if (watched == m_closeBtn3D && event->type() == QEvent::HoverLeave) {
            m_closeBtn3D->setIcon(QIcon(":/icons/resources/close_round_normal.svg"));
        }
        else if (watched == m_closeBtn3D && event->type() == QEvent::MouseButtonPress) {
            m_closeBtn3D->setIcon(QIcon(":/icons/resources/close_round_press.svg"));
        }
    }

    if (watched == m_tracked) {
        if (event->type() == QEvent::Destroy) {
            // 此处需要置空，否则当Destroy事件响应结束后，会在FloatingPreview::hideEvent使用m_tracked野指针
            m_tracked = nullptr;
            hide();
        }

        if (event->type() == QEvent::Resize && m_tracked->width() == SNAP_WIDTH)
            setGeometry(m_tracked->geometry());
    }

    return QWidget::eventFilter(watched, event);
}

void FloatingPreview::hideEvent(QHideEvent *event)
{
    if (m_tracked) {
        m_tracked->setContentsMargins(0, 0, 0, 0);
        m_tracked->setWindowState();
    }

    QWidget::hideEvent(event);
}

void FloatingPreview::onCloseBtnClicked()
{
    Q_ASSERT(!m_tracked.isNull());

    m_tracked->closeWindow();
}
