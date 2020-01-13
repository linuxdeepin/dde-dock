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
#define SHOW_DATE_MIN_HEIGHT 40
#define TIME_FONT DFontSizeManager::instance()->t4()
#define DATE_FONT DFontSizeManager::instance()->t10()

DWIDGET_USE_NAMESPACE

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
{
    QFontMetrics fm_time(TIME_FONT);
    int timeHeight =  fm_time.boundingRect("88:88").height();

    QFontMetrics fm_date(DATE_FONT);
    int dateHeight =  fm_date.boundingRect("8888/88/88").height();

    m_timeOffset = (timeHeight - dateHeight) / 2;

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
    } else {
        QSize dateSize = QFontMetrics(m_dateFont).boundingRect("0000/00/00").size();
        if (timeSize.width() < dateSize.width() && rect().height() >= SHOW_DATE_MIN_HEIGHT)
            timeSize.setWidth(dateSize.width());
    }

    if (position == Dock::Bottom || position == Dock::Top) {
        if (rect().height() >= SHOW_DATE_MIN_HEIGHT) {
            QStringList SL = timeString.split("\n");
            int date_dec=m_timeFont.pixelSize()-m_dateFont.pixelSize();
            while (QFontMetrics(m_timeFont).boundingRect(timeString).size().height() + 11 > height()) {
                m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
                m_dateFont.setPixelSize(m_timeFont.pixelSize() - date_dec);
            }
        } else {
            while (QFontMetrics(m_timeFont).boundingRect(timeString).size().height() > height() + 10) {
                m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            }
            timeSize.setWidth(QFontMetrics(m_timeFont).boundingRect(timeString).size().width());
        }

        return QSize(timeSize.width(), height());
    } else {
        if (width() < timeSize.width()) {
            if (timeString.contains("\n")) {
                QStringList SL = timeString.split("\n");
                while (QFontMetrics(m_timeFont).boundingRect(SL.at(0)).size().width() > width()) {
                    m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
                }
            } else {
                while (QFontMetrics(m_timeFont).boundingRect(timeString).size().width() > width()) {
                    m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
                }
            }

            int timeHeight = QFontMetrics(m_timeFont).boundingRect(timeString).size().height();
            if (format.contains("\n")) {
                QStringList SL = format.split("\n");
                timeHeight = QFontMetrics(m_timeFont).boundingRect(SL.at(0)).size().height() + QFontMetrics(m_timeFont).boundingRect(SL.at(1)).size().height();
            }

            return QSize(width(), std::max(timeHeight, PLUGIN_BACKGROUND_MIN_SIZE));

        } else {
            return QSize(width(), std::max(timeSize.height(), SHOW_DATE_MIN_HEIGHT));
        }
    }
}

QSize DatetimeWidget::sizeHint() const
{
    return curTimeSize();
}

void DatetimeWidget::resizeEvent(QResizeEvent *e)
{
    setMaximumSize(curTimeSize() + QSize(1, 1));

    QWidget::resizeEvent(e);
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

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

    if (rect().height() >= SHOW_DATE_MIN_HEIGHT) {

        QRect timeRect = rect();
        timeRect.setBottom(rect().center().y() + m_timeOffset);

        if (position == Dock::Top || position == Dock::Bottom) {
            painter.drawText(timeRect, Qt::AlignBottom | Qt::AlignHCenter, current.toString(format));

            QRect dateRect = rect();
            dateRect.setTop(timeRect.bottom());
            format = "yyyy/MM/dd";
            painter.setFont(m_dateFont);
            painter.drawText(dateRect, Qt::AlignTop | Qt::AlignHCenter, current.toString(format));
        } else {
            painter.drawText(rect(), Qt::AlignVCenter | Qt::AlignHCenter, current.toString(format));
        }
    } else {
        painter.drawText(rect(), Qt::AlignCenter, current.toString(format));
    }
}
