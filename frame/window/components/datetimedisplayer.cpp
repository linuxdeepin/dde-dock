/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "datetimedisplayer.h"

#include <DFontSizeManager>
#include <DDBusSender>

#include <QHBoxLayout>
#include <QPainter>
#include <QFont>

DWIDGET_USE_NAMESPACE

#define DATETIMESIZE 40
#define ITEMSPACE 8

static QMap<int, QString> dateFormat{{ 0,"yyyy/M/d" }, { 1,"yyyy-M-d" }, { 2,"yyyy.M.d" }, { 3,"yyyy/MM/dd" },
                                     { 4,"yyyy-MM-dd" }, { 5,"yyyy.MM.dd" }, { 6,"yy/M/d" }, { 7,"yy-M-d" }, { 8,"yy.M.d" }};
static QMap<int, QString> timeFormat{{0, "h:mm"}, {1, "hh:mm"}};

DateTimeDisplayer::DateTimeDisplayer(QWidget *parent)
    : QWidget (parent)
    , m_timedateInter(new Timedate("com.deepin.daemon.Timedate", "/com/deepin/daemon/Timedate", QDBusConnection::sessionBus(), this))
    , m_position(Dock::Position::Bottom)
    , m_timeFont(DFontSizeManager::instance()->t6())
    , m_dateFont(DFontSizeManager::instance()->t10())
{
    // 日期格式变化的时候，需要重绘
    connect(m_timedateInter, &Timedate::ShortDateFormatChanged, this, [ this ] { update(); });
    // 时间格式变化的时候，需要重绘
    connect(m_timedateInter, &Timedate::ShortTimeFormatChanged, this, [ this ] { update(); });
    // 连接日期时间修改信号,更新日期时间插件的布局
    connect(m_timedateInter, &Timedate::TimeUpdate, this, [ this ] { update(); });
}

DateTimeDisplayer::~DateTimeDisplayer()
{
}

void DateTimeDisplayer::setPositon(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    setCurrentPolicy();
    update();
}

void DateTimeDisplayer::setCurrentPolicy()
{
    switch (m_position) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        break;
    }
    }
}

QSize DateTimeDisplayer::suitableSize()
{
    DateTimeInfo info = dateTimeInfo();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        int width = info.m_timeRect.width() + info.m_dateRect.width() + 16;
        return QSize(width, height());
    }

    return QSize(width(), info.m_timeRect.height() + info.m_dateRect.height());
}

void DateTimeDisplayer::mouseReleaseEvent(QMouseEvent *event)
{
    Q_UNUSED(event);

    DDBusSender().service("com.deepin.Calendar")
            .path("/com/deepin/Calendar")
            .interface("com.deepin.Calendar")
            .method("RaiseWindow").call();
}

DateTimeDisplayer::DateTimeInfo DateTimeDisplayer::dateTimeInfo()
{
    DateTimeInfo info;
    const QDateTime current = QDateTime::currentDateTime();

    info.m_timeRect = rect();
    info.m_dateRect = rect();

    QString format = getTimeFormat();
    if (!m_timedateInter->use24HourFormat()) {
        if (m_position == Dock::Top || m_position == Dock::Bottom)
            format = format.append(" AP");
        else
            format = format.append("\nAP");
    }

    info.m_time = current.toString(format);
    info.m_date = current.toString(getDateFormat());

    if (m_position == Dock::Top || m_position == Dock::Bottom) {
        int timeWidth = QFontMetrics(m_timeFont).boundingRect(info.m_time).width() + 10;
        int dateWidth = QFontMetrics(m_dateFont).boundingRect(info.m_date).width() + 2;
        info.m_timeRect = QRect(ITEMSPACE, 0, timeWidth, height());
        int dateX = rect().width() - QFontMetrics(m_dateFont).width(info.m_date) - 2 - ITEMSPACE;
        // 如果时间的X坐标小于日期的X坐标，需要手动设置坐标在日期坐标的右侧
        if (dateX < info.m_timeRect.right())
            dateX = info.m_timeRect.right();
        info.m_dateRect = QRect(dateX, 0, dateWidth, height());
    } else {
        int textWidth = rect().width();
        info.m_timeRect = QRect(0, 0, textWidth, DATETIMESIZE / 2);
        info.m_dateRect = QRect(0, DATETIMESIZE / 2 + 1, textWidth, DATETIMESIZE / 2);
    }
    return info;
}

QString DateTimeDisplayer::getDateFormat() const
{
    int type = m_timedateInter->shortDateFormat();
    QString shortDateFormat = "yyyy-MM-dd";
    if (dateFormat.contains(type))
        shortDateFormat = dateFormat.value(type);
    // 如果是左右方向，则不显示年份
    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right) {
        static QStringList yearStrList{"yyyy/", "yyyy-", "yyyy.", "yy/", "yy-", "yy."};
        for (int i = 0; i < yearStrList.size() ; i++) {
            const QString &yearStr = yearStrList[i];
            if (shortDateFormat.contains(yearStr)) {
                shortDateFormat = shortDateFormat.remove(yearStr);
                break;
            }
        }
    }

    return shortDateFormat;
}

QString DateTimeDisplayer::getTimeFormat() const
{
    int type = m_timedateInter->shortTimeFormat();
    if (timeFormat.contains(type))
        return timeFormat[type];

    return QString("hh:mm");
}

void DateTimeDisplayer::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    DateTimeInfo info = dateTimeInfo();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    painter.setPen(QPen(palette().brightText(), 1));

    int timeTextFlag = Qt::AlignCenter;
    int dateTextFlag = Qt::AlignCenter;
    if (m_position == Dock::Top || m_position == Dock::Bottom) {
        timeTextFlag = Qt::AlignLeft | Qt::AlignVCenter;
        dateTextFlag = Qt::AlignRight | Qt::AlignVCenter;
    }
    painter.setFont(m_timeFont);
    painter.drawText(info.m_timeRect, timeTextFlag, info.m_time);
    painter.setFont(m_dateFont);
    painter.drawText(info.m_dateRect, dateTextFlag, info.m_date);
}
