// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "quickdragcore.h"

#include <QWidget>
#include <QTimer>
#include <QPainter>
#include <QPainterPath>
#include <QBitmap>
#include <QEvent>
#include <QDebug>
#include <QCoreApplication>
#include <QDragEnterEvent>
#include <QGuiApplication>

QuickPluginMimeData::QuickPluginMimeData(PluginsItemInterface *item, QDrag *drag)
    : QMimeData()
    , m_item(item)
    , m_drag(drag)
{
}

QuickPluginMimeData::~QuickPluginMimeData()
{
}

PluginsItemInterface *QuickPluginMimeData::pluginItemInterface() const
{
    return m_item;
}

QDrag *QuickPluginMimeData::drag() const
{
    return m_drag;
}

/**
 * @brief 拖动图标的窗口，可以根据实际情况设置动态图标
 * @param dragSource
 */
QuickIconDrag::QuickIconDrag(QObject *dragSource, const QPixmap &pixmap)
    : QDrag(dragSource)
    , m_imageWidget(new QWidget)
    , m_timer(new QTimer(this))
    , m_sourcePixmap(pixmap)
    , m_hotPoint(QPoint(0, 0))
{
    m_timer->setInterval(10);
    connect(m_timer, &QTimer::timeout, this, &QuickIconDrag::onDragMove);
    m_timer->start();

    m_imageWidget->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool | Qt::WindowDoesNotAcceptFocus);
    m_imageWidget->setAttribute(Qt::WA_TransparentForMouseEvents);
    m_imageWidget->installEventFilter(this);
    useSourcePixmap();
}

QuickIconDrag::~QuickIconDrag()
{
    m_imageWidget->deleteLater();
}

void QuickIconDrag::updatePixmap(QPixmap pixmap)
{
    if (m_sourcePixmap == pixmap)
        return;

    m_pixmap = pixmap;
    m_useSourcePixmap = false;
    m_imageWidget->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool | Qt::WindowDoesNotAcceptFocus | Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint);
    m_imageWidget->setFixedSize(pixmap.size());
    m_imageWidget->show();
    m_imageWidget->raise();
    m_imageWidget->update();
}

void QuickIconDrag::useSourcePixmap()
{
    m_useSourcePixmap = true;
    m_imageWidget->setFixedSize(m_sourcePixmap.size() / qApp->devicePixelRatio());
    m_imageWidget->show();
    m_imageWidget->raise();
    m_imageWidget->update();
}

void QuickIconDrag::setDragHotPot(QPoint point)
{
    m_hotPoint = point;
    m_imageWidget->update();
}

bool QuickIconDrag::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_imageWidget) {
        switch (event->type()) {
        case QEvent::Paint: {
            QPixmap pixmap = m_useSourcePixmap ? m_sourcePixmap : m_pixmap;
            QPainter painter(m_imageWidget);
            painter.drawPixmap(QPoint(0, 0), pixmap);

            QPixmap pixmapMask(m_imageWidget->size());
            pixmapMask.fill(Qt::transparent);
            QPainter painterMask(&pixmapMask);
            QPainterPath path;
            path.addRoundedRect(pixmapMask.rect(), 8, 8);
            painterMask.fillPath(path, Qt::white);
            painterMask.setRenderHint(QPainter::Antialiasing, true);
            painterMask.setCompositionMode(QPainter::CompositionMode_Source);
            painterMask.drawPixmap(0, 0, pixmap);
            painterMask.setCompositionMode(QPainter::CompositionMode_DestinationIn);
            QColor maskColor(Qt::black);
            maskColor.setAlpha(150);
            painterMask.fillRect(pixmapMask.rect(), maskColor);
            painterMask.end();

            // 绘制圆角
            QBitmap radiusMask(m_imageWidget->size());
            radiusMask.fill();
            QPainter radiusPainter(&radiusMask);
            radiusPainter.setPen(Qt::NoPen);
            radiusPainter.setBrush(Qt::black);
            radiusPainter.setRenderHint(QPainter::Antialiasing);
            radiusPainter.drawRoundedRect(radiusMask.rect(), 8, 8);
            m_imageWidget->setMask(radiusMask);

            painter.end();
            break;
        }
        default:
            break;
        }
    }
    return QDrag::eventFilter(watched, event);
}

QPoint QuickIconDrag::currentPoint() const
{
    QPoint mousePos = QCursor::pos();
    if (m_useSourcePixmap)
        return mousePos - m_hotPoint;

    QSize pixmapSize = m_pixmap.size();
    return (mousePos - QPoint(pixmapSize.width() * (m_hotPoint.x() / m_sourcePixmap.width())
                              , pixmapSize.height() * (m_hotPoint.y() / m_sourcePixmap.height())));
}

void QuickIconDrag::onDragMove()
{
    m_imageWidget->move(currentPoint());
}
