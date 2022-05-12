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

    private:
        void initUi();
        void setCurrentPolicy();
        DateTimeInfo dateTimeInfo();

        void paintEvent(QPaintEvent *e) override;

    private Q_SLOTS:
        void setShortDateFormat(int type);
        void setShortTimeFormat(int type);

    private:
        Timedate *m_timedateInter;
        QString m_shortDateFormat;
        QString m_shortTimeFormat;
        Dock::Position m_position;
        mutable QFont m_timeFont;
        mutable QFont m_dateFont;
};

#endif // DATETIMEDISPLAYER_H
