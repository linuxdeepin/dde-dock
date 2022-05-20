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

namespace Dock { class TipsWidget; }
class DockPopupWindow;

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
    void setOneRow(bool oneRow);
    QSize suitableSize();

Q_SIGNALS:
    void sizeChanged();         // 当日期时间格式发生变化的时候，需要通知外面来更新窗口尺寸

protected:
    void mouseReleaseEvent(QMouseEvent *event) override;
    void paintEvent(QPaintEvent *e) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    void updatePolicy();
    DateTimeInfo dateTimeInfo();
    void updateLastData(const DateTimeInfo &info);

    QString getTimeString() const;
    QString getDateString() const;

    QPoint tipsPoint() const;
    QFont timeFont() const;

private Q_SLOTS:
    void onTimeChanged();
    void onDateTimeFormatChanged();

private:
    Timedate *m_timedateInter;
    Dock::Position m_position;
    QFont m_dateFont;
    Dock::TipsWidget *m_tipsWidget;
    QSharedPointer<DockPopupWindow> m_tipPopupWindow;
    QTimer *m_tipsTimer;
    QString m_lastDateString;
    QString m_lastTimeString;
    int m_currentSize;
    bool m_oneRow;
};

#endif // DATETIMEDISPLAYER_H
