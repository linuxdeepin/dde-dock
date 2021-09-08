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

#ifndef DATETIMEWIDGET_H
#define DATETIMEWIDGET_H

#include <com_deepin_daemon_timedate.h>

#include <QWidget>

using Timedate = com::deepin::daemon::Timedate;

class DatetimeWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DatetimeWidget(QWidget *parent = 0);

    bool is24HourFormat() const { return m_24HourFormat; }
    QSize sizeHint() const;

protected:
    void resizeEvent(QResizeEvent *event);
    void paintEvent(QPaintEvent *e);

signals:
    void requestUpdateGeometry() const;

public slots:
    void set24HourFormat(const bool value);

private Q_SLOTS:
    void setShortDateFormat(int type);
    void setShortTimeFormat(int type);

private:
    QSize curTimeSize() const;

private:
    bool m_24HourFormat;
    mutable QFont m_timeFont;
    mutable QFont m_dateFont;
    mutable int m_timeOffset;
    Timedate *m_timedateInter;
    QString m_shortDateFormat;
    QString m_shortTimeFormat;
};

#endif // DATETIMEWIDGET_H
