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
#include <DFontSizeManager>
#include <DGuiApplicationHelper>

#define PLUGIN_STATE_KEY    "enable"
#define TIME_FONT DFontSizeManager::instance()->t4()
#define DATE_FONT DFontSizeManager::instance()->t10()

DWIDGET_USE_NAMESPACE

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
{
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);
}

void DatetimeWidget::set24HourFormat(const bool value)
{
    if (m_24HourFormat == value) {
        return;
    }

    m_24HourFormat = value;
    update();

    adjustSize();
    if (isVisible()) {
        emit requestUpdateGeometry();
    }
}

QSize DatetimeWidget::curTimeSize() const
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

    m_timeFont = TIME_FONT;
    m_dateFont = DATE_FONT;
    QFontMetrics fm(m_timeFont);
    QString format;
    if (m_24HourFormat)
        format = "hh:mm";
    else {
        format = "hh:mm AP";
    }

    QString timeString = QDateTime::currentDateTime().toString(format);
    QSize timeSize = fm.boundingRect(timeString).size();
    QSize dateSize = QFontMetrics(m_dateFont).boundingRect("0000/00/00").size();

    if (position == Dock::Bottom || position == Dock::Top) {
        while (QFontMetrics(m_timeFont).boundingRect(timeString).size().height() + QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().height() > height()) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            timeSize.setWidth(QFontMetrics(m_timeFont).boundingRect(timeString).size().width());
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1){
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width());
            }
        }

        return QSize(std::max(timeSize.width(), dateSize.width()) + 2, height());
    } else {
        while (std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width()) > width()) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height());
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1){
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().height());
            }
        }

        return QSize(std::max(timeSize.width(), dateSize.width()), timeSize.height() + dateSize.height());
    }
}

QSize DatetimeWidget::sizeHint() const
{
    return curTimeSize();
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    QString format;
    if (m_24HourFormat)
        format = "hh:mm";
    else {
        format = "hh:mm AP";
    }

    painter.setFont(m_timeFont);
    painter.setPen(QPen(palette().brightText(), 1));

    //由于时间和日期字体不是同等缩小，会导致时间和日期位置不居中，需要整体往下移动几个像素，
    int offsetY = 0;
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    if (position == Dock::Bottom || position == Dock::Top) {
        if (height() >= 60)
            offsetY = 5;
        else if (height() >= 50)
            offsetY = 4;
        else if (height() >= 40)
            offsetY = 2;
        else if (height() >= 20)
            offsetY = 1;
    }
    QRect timeRect = rect();
    timeRect.setBottom(rect().center().y() + offsetY);
    painter.drawText(timeRect, Qt::AlignBottom | Qt::AlignHCenter, current.toString(format));
    QRect dateRect = rect();
    dateRect.setTop(timeRect.bottom() - 2);
    format = "yyyy/MM/dd";
    painter.setFont(m_dateFont);
    painter.drawText(dateRect, Qt::AlignTop | Qt::AlignHCenter, current.toString(format));
}
