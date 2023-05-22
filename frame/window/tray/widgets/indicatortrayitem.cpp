// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "indicatortrayitem.h"

#include <cmath>

#include <QLabel>
#include <QBoxLayout>

#include <QDBusConnection>
#include <QDBusInterface>
#include <QFile>
#include <QDebug>
#include <QPoint>
#include <QPainter>
#include <QJsonDocument>
#include <QJsonObject>
#include <QDBusInterface>
#include <QDBusReply>
#include <QApplication>
#include <thread>
#include <DFontSizeManager>

DWIDGET_USE_NAMESPACE

IndicatorTrayItem::IndicatorTrayItem(const QString &indicatorName, QWidget *parent, Qt::WindowFlags f)
    : BaseTrayWidget(parent, f)
    , m_indicatorName(indicatorName)
    , m_enableClick(true)
{
    setAttribute(Qt::WA_TranslucentBackground);
    QPalette p = palette();
    p.setColor(QPalette::WindowText, Qt::white);
    p.setColor(QPalette::Window, Qt::transparent);
    setPalette(p);

    QFont qf = DFontSizeManager::instance()->t10();
    // make indicator font size fixed, 16 x 16 is a standard icon size
    qf.setPixelSize(16);
    setFont(qf);

    // register dbus
    auto path = QString("/org/deepin/dde/Dock1/Indicator/") + m_indicatorName;
    auto interface =  QString("org.deepin.dde.Dock1.Indicator.") + m_indicatorName;
    auto sessionBus = QDBusConnection::sessionBus();
    sessionBus.registerObject(path,
                              interface,
                              this,
                              QDBusConnection::ExportScriptableSlots);
}

IndicatorTrayItem::~IndicatorTrayItem()
{
}

QString IndicatorTrayItem::itemKeyForConfig()
{
    return toIndicatorKey(m_indicatorName);
}

void IndicatorTrayItem::updateIcon()
{

}

void IndicatorTrayItem::sendClick(uint8_t buttonIndex, int x, int y)
{
    if (m_enableClick)
        Q_EMIT clicked(buttonIndex, x, y);
}

void IndicatorTrayItem::enableLabel(bool enable)
{
    QPalette p = palette();
    if (!enable) {
        m_enableClick = false;
        p.setColor(QPalette::Disabled, QPalette::WindowText, Qt::lightGray);
        p.setColor(QPalette::Disabled, QPalette::Window, Qt::transparent);
        setEnabled(enable);
    } else {
        m_enableClick = true;
        p.setColor(QPalette::Normal, QPalette::BrightText, Qt::white);
        p.setColor(QPalette::Normal, QPalette::Window, Qt::transparent);
        setEnabled(enable);
    }

    setPalette(p);
    update();
}

QPixmap IndicatorTrayItem::icon()
{
    auto rawPixmap = QPixmap::fromImage(QImage::fromData(m_pixmapData));
    rawPixmap.setDevicePixelRatio(devicePixelRatioF());
    return (rawPixmap);
}

const QByteArray &IndicatorTrayItem::pixmapData() const
{
    return m_pixmapData;
}

const QString IndicatorTrayItem::text() const
{
    return m_text;
}

void IndicatorTrayItem::setPixmapData(const QByteArray &data)
{
    m_pixmapData = data;
}

void IndicatorTrayItem::setText(const QString &text)
{
    m_text = text;
    Q_EMIT textChanged(m_text);
    update();
}

void IndicatorTrayItem::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    QFontMetrics qfm = QFontMetrics(font());
    QRect tightTextRect = qfm.tightBoundingRect(m_text);
    QRect textRect = qfm.boundingRect(m_text);
    QPoint topLeft = textRect.topLeft() - tightTextRect.topLeft();
    QPoint bottom = textRect.bottomRight() - tightTextRect.bottomRight();
    QPoint center = QPoint(std::floor((rect().width() - textRect.width()) / 2), std::floor((rect().height() - textRect.height()) / 2)) // this make textRect in center
                        + (topLeft + bottom) / 2; // this adjust make tightTextRect in center
    painter.drawText(QRect(center.x(), center.y(), textRect.width() + 1, textRect.height() + 1), (m_text));

    if (m_pixmapData != nullptr) {
        painter.drawPixmap(rect(), icon());
    }
}
