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
#ifndef DATETIMEDISPLAYER_H
#define DATETIMEDISPLAYER_H

#include "constants.h"

#include <QWidget>
#include <QFont>

#include <com_deepin_daemon_timedate.h>

using Timedate = com::deepin::daemon::Timedate;

class DateTimeDisplayer : public QWidget
{
    Q_OBJECT

private:
    struct DateTimeInfo {
        QString m_time;
        QString m_date;
        QRect m_timeRect;
        QRect m_dateRect;
    };

public:
    explicit DateTimeDisplayer(QWidget *parent = nullptr);
    ~DateTimeDisplayer() override;
    void setPositon(Dock::Position position);
    QSize suitableSize();

protected:
    void mouseReleaseEvent(QMouseEvent *event) override;
    void paintEvent(QPaintEvent *e) override;

private:
    void setCurrentPolicy();
    DateTimeInfo dateTimeInfo();

    QString getDateFormat() const;
    QString getTimeFormat() const;

private:
    Timedate *m_timedateInter;
    Dock::Position m_position;
    mutable QFont m_timeFont;
    mutable QFont m_dateFont;
};

#endif // DATETIMEDISPLAYER_H
