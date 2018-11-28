/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "fashiontraywidgetwrapper.h"

#include <QPainter>
#include <QDebug>
#include <QMouseEvent>
#include <QDrag>
#include <QMimeData>

#define TRAY_ITEM_DRAG_THRESHOLD 20

FashionTrayWidgetWrapper::FashionTrayWidgetWrapper(const QString &itemKey, AbstractTrayWidget *absTrayWidget, QWidget *parent)
    : QWidget(parent),
      m_absTrayWidget(absTrayWidget),
      m_layout(new QVBoxLayout(this)),
      m_attention(false),
      m_dragging(false),
      m_itemKey(itemKey)

{
    setStyleSheet("background: transparent;");
    setAcceptDrops(true);

    m_absTrayWidget->setVisible(true);
    m_absTrayWidget->installEventFilter(this);

    m_layout->setSpacing(0);
    m_layout->setMargin(0);
    m_layout->setContentsMargins(0, 0, 0, 0);

    m_layout->addWidget(m_absTrayWidget);

    setLayout(m_layout);

    connect(m_absTrayWidget, &AbstractTrayWidget::needAttention, this, &FashionTrayWidgetWrapper::onTrayWidgetNeedAttention);
    connect(m_absTrayWidget, &AbstractTrayWidget::clicked, this, &FashionTrayWidgetWrapper::onTrayWidgetClicked);
}

AbstractTrayWidget *FashionTrayWidgetWrapper::absTrayWidget() const
{
    return m_absTrayWidget;
}

QString FashionTrayWidgetWrapper::itemKey() const
{
    return m_itemKey;
}

void FashionTrayWidgetWrapper::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);

    if (m_dragging) {
        return;
    }

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setOpacity(0.5);

    QPainterPath path;
    path.addRoundedRect(rect(), 10, 10);
    painter.fillPath(path, QColor::fromRgb(40, 40, 40));
}

bool FashionTrayWidgetWrapper::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_absTrayWidget) {
        if (event->type() == QEvent::Type::MouseButtonPress) {
            mousePressEvent(static_cast<QMouseEvent *>(event));
        } else if (event->type() == QEvent::Type::MouseMove) {
            handleMouseMove(static_cast<QMouseEvent *>(event));
        }
    }

    return false;
}

void FashionTrayWidgetWrapper::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::MouseButton::LeftButton) {
        MousePressPoint = event->pos();
    }

    QWidget::mousePressEvent(event);
}

void FashionTrayWidgetWrapper::mouseMoveEvent(QMouseEvent *event)
{
    return QWidget::mouseMoveEvent(event);
}

void FashionTrayWidgetWrapper::dragEnterEvent(QDragEnterEvent *event)
{
    if (event->mimeData()->hasFormat(TRAY_ITEM_DRAG_MIMEDATA)) {
        event->accept();
        Q_EMIT requestSwapWithDragging();
        return;
    }

    QWidget::dragEnterEvent(event);
}

void FashionTrayWidgetWrapper::handleMouseMove(QMouseEvent *event)
{
    if (event->buttons() != Qt::MouseButton::LeftButton) {
        return QWidget::mouseMoveEvent(event);
    }

//    if ((event->pos() - MousePressPoint).manhattanLength() < TRAY_ITEM_DRAG_THRESHOLD) {
//        return QWidget::mouseMoveEvent(event);
//    }

    event->accept();

    QDrag drag(this);
    QMimeData *mimeData = new QMimeData;
    mimeData->setData(TRAY_ITEM_DRAG_MIMEDATA, m_itemKey.toLocal8Bit());
    QPixmap pixmap = grab();

    drag.setMimeData(mimeData);
    drag.setPixmap(pixmap);
    drag.setHotSpot(pixmap.rect().center() / pixmap.devicePixelRatioF());

    m_absTrayWidget->setVisible(false);
    m_dragging = true;
    Q_EMIT dragStart();

    // start drag
    drag.exec();

    m_absTrayWidget->setVisible(true);
    m_dragging = false;
    Q_EMIT dragStop();
}

void FashionTrayWidgetWrapper::onTrayWidgetNeedAttention()
{
    setAttention(true);
}

void FashionTrayWidgetWrapper::onTrayWidgetClicked()
{
    setAttention(false);
}

bool FashionTrayWidgetWrapper::attention() const
{
    return m_attention;
}

void FashionTrayWidgetWrapper::setAttention(bool attention)
{
    m_attention = attention;

    Q_EMIT attentionChanged(m_attention);
}
