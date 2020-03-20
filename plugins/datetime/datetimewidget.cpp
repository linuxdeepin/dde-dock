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

#include "datetimewidget.h"
#include "constants.h"

#include <QApplication>
#include <QPainter>
#include <QDebug>
#include <QSvgRenderer>
#include <QMouseEvent>

#define PLUGIN_STATE_KEY "enable"

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
{
}

void DatetimeWidget::set24HourFormat(const bool value)
{
    if (m_24HourFormat == value)
    {
        return;
    }

    m_24HourFormat = value;

    m_cachedTime.clear();
    update();

    if (isVisible())
    {
        emit requestUpdateGeometry();
    }
}

QSize DatetimeWidget::sizeHint() const
{
    QFontMetrics fm(qApp->font());

    if (m_24HourFormat)
        return fm.boundingRect(this->rect(), Qt::AlignCenter, "88:88:88\n88/88/8888").size() + QSize(20, 10);
    else
        return fm.boundingRect(this->rect(), Qt::AlignCenter, "88:88:88 A.A.\n88/88/8888").size() + QSize(20, 20);
}

void DatetimeWidget::resizeEvent(QResizeEvent *e)
{
    m_cachedTime.clear();

    QWidget::resizeEvent(e);
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const auto ratio = devicePixelRatioF();
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    if (displayMode == Dock::Efficient)
    {
        QString format;
        if (m_24HourFormat)
            format = "hh:mm:ss\nd/M/yyyy";
        else
        {
            if (position == Dock::Top || position == Dock::Bottom)
                format = "hh:mm:ss AP\nd/M/yyyy";
            else
                format = "hh:mm\nAP";
        }

        painter.setPen(QPen(palette().brightText(), 1));
        painter.drawText(rect(), Qt::AlignCenter, current.toString(format));
        return;
    }

    // use language Chinese to fix can not find image resources which will be drawn
    const QString currentTimeString = current.toString(m_24HourFormat ? "hhmmss" : "hhmmssa");

    // check cache valid
    if (m_cachedTime != currentTimeString)
    {
        m_cachedTime = currentTimeString;

        // draw new pixmap
        m_cachedTime = currentTimeString;
        m_cachedIcon = QPixmap(size() * ratio);
        m_cachedIcon.fill(Qt::transparent);
        m_cachedIcon.setDevicePixelRatio(ratio);
        QPainter p(&m_cachedIcon);

        // draw fashion mode datetime plugin
        const int perfectIconSize = qMin(width(), height()) * 0.8;
        const QRect r = rect();

        // draw background
        QPixmap background = loadSvg(":/icons/resources/icons/background.svg", QSize(perfectIconSize, perfectIconSize));
        const QPoint backgroundOffset = r.center() - background.rect().center() / ratio;
        p.drawPixmap(backgroundOffset, background);

        const int bigNumHeight = perfectIconSize / 2.5;
        const int bigNumWidth = double(bigNumHeight) * 8 / 18;
        const int smallNumHeight = bigNumHeight / 2;
        const int smallNumWidth = double(smallNumHeight) * 5 / 9;

        // draw big num 1
        const QString bigNum1Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[0]);
        const QPixmap bigNum1 = loadSvg(bigNum1Path, QSize(bigNumWidth, bigNumHeight));
        const QPoint bigNum1Offset = backgroundOffset + QPoint(perfectIconSize / 2 - bigNumWidth * 2 + 1, perfectIconSize / 2 - bigNumHeight / 2);
        p.drawPixmap(bigNum1Offset, bigNum1);

        // draw big num 2
        const QString bigNum2Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[1]);
        const QPixmap bigNum2 = loadSvg(bigNum2Path, QSize(bigNumWidth, bigNumHeight));
        const QPoint bigNum2Offset = bigNum1Offset + QPoint(bigNumWidth + 1, 0);
        p.drawPixmap(bigNum2Offset, bigNum2);

        if (!m_24HourFormat)
        {
            // draw small num 1
            const QString smallNum1Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[2]);
            const QPixmap smallNum1 = loadSvg(smallNum1Path, QSize(smallNumWidth, smallNumHeight));
            const QPoint smallNum1Offset = bigNum2Offset + QPoint(bigNumWidth + 2, 1);
            p.drawPixmap(smallNum1Offset, smallNum1);

            // draw small num 2
            const QString smallNum2Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[3]);
            const QPixmap smallNum2 = loadSvg(smallNum2Path, QSize(smallNumWidth, smallNumHeight));
            const QPoint smallNum2Offset = smallNum1Offset + QPoint(smallNumWidth + 1, 0);
            p.drawPixmap(smallNum2Offset, smallNum2);

            // draw am/pm tips
            const int tips_width = (smallNumWidth * 2 + 2) & ~0x1;
            const int tips_height = tips_width / 2;

            QPixmap tips;
            if (current.time().hour() > 11)
                tips = loadSvg(":/icons/resources/icons/tips-pm.svg", QSize(tips_width, tips_height));
            else
                tips = loadSvg(":/icons/resources/icons/tips-am.svg", QSize(tips_width, tips_height));

            const QPoint tipsOffset = bigNum2Offset + QPoint(bigNumWidth + 2, bigNumHeight - tips_height);
            p.drawPixmap(tipsOffset, tips);
        }
        else
        {
            // draw small num 1
            const QString smallNum1Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[2]);
            const QPixmap smallNum1 = loadSvg(smallNum1Path, QSize(smallNumWidth, smallNumHeight));
            const QPoint smallNum1Offset = bigNum2Offset + QPoint(bigNumWidth + 2, smallNumHeight);
            p.drawPixmap(smallNum1Offset, smallNum1);

            // draw small num 2
            const QString smallNum2Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[3]);
            const QPixmap smallNum2 = loadSvg(smallNum2Path, QSize(smallNumWidth, smallNumHeight));
            const QPoint smallNum2Offset = smallNum1Offset + QPoint(smallNumWidth + 1, 0);
            p.drawPixmap(smallNum2Offset, smallNum2);
        }
    }

    // draw cached fashion mode time item
    painter.drawPixmap(rect().center() - m_cachedIcon.rect().center() / ratio, m_cachedIcon);
}

const QPixmap DatetimeWidget::loadSvg(const QString &fileName, const QSize size)
{
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap(size * ratio);
    QSvgRenderer renderer(fileName);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
