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
        if (position == Dock::Top || position == Dock::Bottom)
            format = "hh:mm AP";
        else
            format = "hh:mm\nAP";
    }

    QString timeString = QDateTime::currentDateTime().toString(format);
    QSize timeSize = fm.boundingRect(timeString).size();
    if (timeString.contains("\n")) {
         QStringList SL = timeString.split("\n");
         timeSize = QSize(fm.boundingRect(SL.at(0)).width(), fm.boundingRect(SL.at(0)).height() + fm.boundingRect(SL.at(1)).height());
    }

    QSize dateSize = QFontMetrics(m_dateFont).boundingRect("0000/00/00").size();

    if (position == Dock::Bottom || position == Dock::Top) {
        while (QFontMetrics(m_timeFont).boundingRect(timeString).size().height() + QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().height() > height()) {
            // 增加对setPixelSize的参数的合法性的判断，防止出现Font的pixelSize值异常时进入死循环打印日志 QFont::setPixelSize: Pixel size <= 0 (0) 占满磁盘空间
            if (((m_timeFont.pixelSize() - 1) <= 0) || ((m_dateFont.pixelSize() - 1) <= 0)) {
                qDebug() << "Invalid pixel size value:  timeFont = " << m_timeFont.pixelSize() <<"  dateFont = " << m_dateFont.pixelSize();
                break;
            }
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            timeSize.setWidth(QFontMetrics(m_timeFont).boundingRect(timeString).size().width());
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1){
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width());
            }
        }
        return QSize(std::max(timeSize.width(), dateSize.width()) + 2, height());
    } else {
        while (std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width()) > (width() - 4)) {
            // 增加对setPixelSize的参数的合法性的判断，防止出现Font的pixelSize值异常时进入死循环打印日志 QFont::setPixelSize: Pixel size <= 0 (0) 占满磁盘空间
            if (((m_timeFont.pixelSize() - 1) <= 0) || ((m_dateFont.pixelSize() - 1) <= 0)) {
                qDebug() << "Invalid pixel size value:  timeFont = " << m_timeFont.pixelSize() <<"   dateFont = " << m_dateFont.pixelSize();
                break;
            }
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            if (m_24HourFormat) {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height());
            } else {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height() * 2);
            }
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1){
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().height());
            }
        }
        m_timeOffset = (timeSize.height() - dateSize.height()) / 2 ;
        return QSize(width(), timeSize.height() + dateSize.height());
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

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    QString format;
    if (m_24HourFormat)
        format = "hh:mm";
    else {
        if (position == Dock::Top || position == Dock::Bottom)
            format = "hh:mm AP";
        else
            format = "hh:mm\nAP";
    }

    painter.setFont(m_timeFont);
    painter.setPen(QPen(palette().brightText(), 1));

    QRect timeRect = rect();
    QRect dateRect = rect();

    if (position == Dock::Top || position == Dock::Bottom){
       timeRect.setBottom(rect().center().y() + 6);
       dateRect.setTop(timeRect.bottom() - 4);
    } else {
        timeRect.setBottom(rect().center().y() + m_timeOffset);
        dateRect.setTop(timeRect.bottom());
    }
    painter.drawText(timeRect, Qt::AlignBottom | Qt::AlignHCenter, current.toString(format));
    format = "yyyy/MM/dd";
    painter.setFont(m_dateFont);
    painter.drawText(dateRect, Qt::AlignTop | Qt::AlignHCenter, current.toString(format));
}
