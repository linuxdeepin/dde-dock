// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DATETIMEWIDGET_H
#define DATETIMEWIDGET_H

#include <com_deepin_daemon_timedate.h>

#include <QWidget>

using Timedate = com::deepin::daemon::Timedate;

class DatetimeWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DatetimeWidget(QWidget *parent = nullptr);

    QSize sizeHint() const;
    inline bool is24HourFormat() const { return m_24HourFormat; }
    inline QString getDateTime() { return m_dateTime; }

protected:
    void resizeEvent(QResizeEvent *event);
    void paintEvent(QPaintEvent *e);

signals:
    void requestUpdateGeometry() const;

public slots:
    void set24HourFormat(const bool value);
    void updateDateTimeString();

private Q_SLOTS:
    void setShortDateFormat(int type);
    void setShortTimeFormat(int type);
    void setLongDateFormat(int type);
    void setWeekdayFormat(int type);
    void setLongTimeFormat(int type);

private:
    QSize curTimeSize() const;
    void updateWeekdayFormat();
    void updateLongTimeFormat();

private:
    bool m_24HourFormat;
    int m_longDateFormatType;
    int m_longTimeFormatType;
    int m_weekdayFormatType;
    mutable QFont m_timeFont;
    mutable QFont m_dateFont;
    mutable int m_timeOffset;
    Timedate *m_timedateInter;
    QString m_shortDateFormat;
    QString m_shortTimeFormat;
    QString m_dateTime;
    QString m_weekFormat;
    QString m_longTimeFormat;
};

#endif // DATETIMEWIDGET_H
