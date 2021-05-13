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
    , m_24HourFormat(false)
    , m_timeOffset(false)
    , m_timedateInter(new Timedate("com.deepin.daemon.Timedate", "/com/deepin/daemon/Timedate", QDBusConnection::sessionBus(), this))
    , m_shortDateFormat("yyyy-MM-dd")
    , m_shortTimeFormat("hh:mm")
{
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);
    setShortDateFormat(m_timedateInter->shortDateFormat());
    setShortTimeFormat(m_timedateInter->shortTimeFormat());

    connect(m_timedateInter, &Timedate::ShortDateFormatChanged, this, &DatetimeWidget::setShortDateFormat);
    connect(m_timedateInter, &Timedate::ShortTimeFormatChanged, this, &DatetimeWidget::setShortTimeFormat);
    //连接日期时间修改信号,更新日期时间插件的布局
    connect(m_timedateInter, &Timedate::TimeUpdate, this, [ = ]{
        if (isVisible()) {
            emit requestUpdateGeometry();
        }
    });
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

/**
 * @brief DatetimeWidget::setShortDateFormat 根据类型设置时间显示格式
 * @param type 自定义类型
 */
void DatetimeWidget::setShortDateFormat(int type)
{
    switch (type) {
    case 0: m_shortDateFormat = "yyyy/M/d";  break;
    case 1: m_shortDateFormat = "yyyy-M-d"; break;
    case 2: m_shortDateFormat = "yyyy.M.d"; break;
    case 3: m_shortDateFormat = "yyyy/MM/dd"; break;
    case 4: m_shortDateFormat = "yyyy-MM-dd"; break;
    case 5: m_shortDateFormat = "yyyy.MM.dd"; break;
    case 6: m_shortDateFormat = "yy/M/d"; break;
    case 7: m_shortDateFormat = "yy-M-d"; break;
    case 8: m_shortDateFormat = "yy.M.d"; break;
    default: m_shortDateFormat = "yyyy-MM-dd"; break;
    }
    update();

    if (isVisible()) {
        emit requestUpdateGeometry();
    }
}

/**
 * @brief DatetimeWidget::setShortTimeFormat 根据类型设置短时间显示格式
 * @param type 自定义类型
 */
void DatetimeWidget::setShortTimeFormat(int type)
{
    switch (type) {
    case 0: m_shortTimeFormat = "h:mm"; break;
    case 1: m_shortTimeFormat = "hh:mm";  break;
    default: m_shortTimeFormat = "hh:mm"; break;
    }
    update();

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
    QString format = m_shortTimeFormat;
    if (!m_24HourFormat) {
        if (position == Dock::Top || position == Dock::Bottom)
            format = format.append(" AP");
        else
            format = format.append("\nAP");
    }

    QString timeString = QDateTime::currentDateTime().toString(format);
    QSize timeSize = fm.boundingRect(timeString).size();
    if (timeString.contains("\n")) {
        QStringList SL = timeString.split("\n");
        timeSize = QSize(fm.boundingRect(SL.at(0)).width(), fm.boundingRect(SL.at(0)).height() + fm.boundingRect(SL.at(1)).height());
    }

    QSize dateSize = QFontMetrics(m_dateFont).boundingRect("0000/00/00").size();

    if (position == Dock::Bottom || position == Dock::Top) {
        while (QFontMetrics(m_timeFont).boundingRect(timeString).height() + QFontMetrics(m_dateFont).boundingRect("0000/00/00").height() > height()) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            timeSize.setWidth(QFontMetrics(m_timeFont).boundingRect(timeString).size().width());
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1) {
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width());
            }
        }
        return QSize(std::max(timeSize.width(), dateSize.width()), timeSize.height() + dateSize.height());
    } else {
        while (std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().width()) > (width() - 4)) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            if (m_24HourFormat) {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height());
            } else {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height() * 2);
            }
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1) {
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect("0000/00/00").size().height());
            }
        }
        m_timeOffset = (timeSize.height() - dateSize.height()) / 2 ;
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

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    painter.setPen(QPen(palette().brightText(), 1));

    QRect timeRect = rect();
    QRect dateRect = rect();

    QString format = m_shortTimeFormat;
    if (!m_24HourFormat) {
        if (position == Dock::Top || position == Dock::Bottom)
            format = format.append(" AP");
        else
            format = format.append("\nAP");
    }
    QString timeStr = current.toString(format);

    format = m_shortDateFormat;
    QString dateStr = current.toString(format);

    if (position == Dock::Top || position == Dock::Bottom) {
        timeRect.setBottom(rect().top() + QFontMetrics(m_timeFont).boundingRect(timeStr).height() + 6);
        dateRect.setTop(timeRect.bottom() - 4);
    } else {
        timeRect.setBottom(rect().center().y() + m_timeOffset);
        dateRect.setTop(timeRect.bottom());
    }
    painter.setFont(m_timeFont);
    painter.drawText(timeRect, Qt::AlignCenter, timeStr);

    painter.setFont(m_dateFont);
    painter.drawText(dateRect, Qt::AlignCenter, dateStr);
}
