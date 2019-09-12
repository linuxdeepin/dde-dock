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

#define PLUGIN_STATE_KEY    "enable"
#define SHOW_DATE_MIN_HEIGHT 45
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

    if (isVisible()) {
        emit requestUpdateGeometry();
    }
}

QSize DatetimeWidget::curTimeSize() const
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

    m_timeFont = TIME_FONT;
    QFontMetrics fm(m_timeFont);
    QString format;
    if (m_24HourFormat)
        format = "hh:mm";
    else
        format = "hh:mm AP";

    QSize timeSize = fm.boundingRect(QDateTime::currentDateTime().toString(format)).size();
    QSize dateSize = QFontMetrics(DATE_FONT).boundingRect("0000/00/00").size();
    if (timeSize.width() < dateSize.width())
        timeSize.setWidth(dateSize.width());

    if (position == Dock::Bottom || position == Dock::Top) {
        return QSize(timeSize.width(), DOCK_MAX_SIZE);
    } else {
        // 宽度不够显示日期则隐藏日期，宽度不够显示时间，则缩小时间字体适应dock
        if (width() < dateSize.width()) {
            QString timeString = QDateTime::currentDateTime().toString(format);

            while (QFontMetrics(m_timeFont).boundingRect(timeString).size().width() > width()) {
                m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            }

            int timeHeight = QFontMetrics(m_timeFont).boundingRect(timeString).size().height();

            return QSize(DOCK_MAX_SIZE, std::max(timeHeight, PLUGIN_BACKGROUND_MIN_SIZE));

        } else {
            return QSize(DOCK_MAX_SIZE, timeSize.height() + dateSize.height());
        }
    }
}

QSize DatetimeWidget::sizeHint() const
{
    return curTimeSize();
}

void DatetimeWidget::resizeEvent(QResizeEvent *e)
{
    setMaximumSize(curTimeSize());

    QWidget::resizeEvent(e);
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
    else
        format = "hh:mm AP";

    painter.setPen(Qt::white);
    painter.setFont(m_timeFont);

    if (rect().height() > SHOW_DATE_MIN_HEIGHT) {
        QRect timeRect = rect();
        timeRect.setBottom(rect().center().y() + m_timeOffset);
        painter.drawText(timeRect, Qt::AlignBottom | Qt::AlignHCenter, current.toString(format));

        QRect dateRect = rect();
        dateRect.setTop(timeRect.bottom());
        format = "yyyy/MM/dd";
        painter.setFont(DATE_FONT);
        painter.drawText(dateRect, Qt::AlignTop | Qt::AlignHCenter, current.toString(format));

    } else {
        painter.drawText(rect(), Qt::AlignCenter, current.toString(format));
    }
}
