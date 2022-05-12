#include "datetimedisplayer.h"

#include <DFontSizeManager>

#include <QHBoxLayout>
#include <QPainter>
#include <QFont>

DWIDGET_USE_NAMESPACE

#define DATETIMESIZE 40

DateTimeDisplayer::DateTimeDisplayer(QWidget *parent)
    : QWidget (parent)
    , m_timedateInter(new Timedate("com.deepin.daemon.Timedate", "/com/deepin/daemon/Timedate", QDBusConnection::sessionBus(), this))
    , m_position(Dock::Position::Bottom)
    , m_timeFont(DFontSizeManager::instance()->t6())
    , m_dateFont(DFontSizeManager::instance()->t10())
{
    initUi();
    setShortDateFormat(m_timedateInter->shortDateFormat());
    setShortTimeFormat(m_timedateInter->shortTimeFormat());
    connect(m_timedateInter, &Timedate::ShortDateFormatChanged, this, &DateTimeDisplayer::setShortDateFormat);
    connect(m_timedateInter, &Timedate::ShortTimeFormatChanged, this, &DateTimeDisplayer::setShortTimeFormat);
    // 连接日期时间修改信号,更新日期时间插件的布局
    connect(m_timedateInter, &Timedate::TimeUpdate, this, [ this ] {
        update();
    });
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

void DateTimeDisplayer::initUi()
{
    QHBoxLayout *layout = new QHBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setSpacing(0);
}

QSize DateTimeDisplayer::suitableSize()
{
    DateTimeInfo info = dateTimeInfo();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return QSize(info.m_dateRect.right(), height());

    return QSize(width(), info.m_dateRect.bottom());
}

DateTimeDisplayer::DateTimeInfo DateTimeDisplayer::dateTimeInfo()
{
    DateTimeInfo info;
    const QDateTime current = QDateTime::currentDateTime();

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

    info.m_timeRect = rect();
    info.m_dateRect = rect();

    QString format = m_shortTimeFormat;
    if (!m_timedateInter->use24HourFormat()) {
        if (position == Dock::Top || position == Dock::Bottom)
            format = format.append(" AP");
        else
            format = format.append("\nAP");
    }

    info.m_time = current.toString(format);
    info.m_date = current.toString(m_shortDateFormat);
    int timeWidth = QFontMetrics(m_timeFont).boundingRect(info.m_time).width() + 12 * 2;
    int dateWidth = QFontMetrics(m_dateFont).boundingRect(info.m_date).width() + 2;

    if (position == Dock::Top || position == Dock::Bottom) {
        info.m_timeRect = QRect(10, 0, timeWidth, height());
        int right = rect().width() - QFontMetrics(m_dateFont).width(info.m_date) - 2;
        info.m_dateRect = QRect(right, 0, dateWidth, height());
    } else {
        info.m_timeRect = QRect(0, 0, timeWidth, DATETIMESIZE / 2);
        info.m_dateRect = QRect(0, DATETIMESIZE / 2 + 1, dateWidth, DATETIMESIZE / 2);
    }
    return info;
}

void DateTimeDisplayer::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    DateTimeInfo info = dateTimeInfo();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    painter.setPen(QPen(palette().brightText(), 1));

    painter.drawText(info.m_timeRect, Qt::AlignLeft | Qt::AlignVCenter, info.m_time);

    painter.setFont(m_dateFont);
    painter.drawText(info.m_dateRect, Qt::AlignLeft | Qt::AlignVCenter, info.m_date);
}

void DateTimeDisplayer::setShortDateFormat(int type)
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
}

void DateTimeDisplayer::setShortTimeFormat(int type)
{
    switch (type) {
    case 0: m_shortTimeFormat = "h:mm"; break;
    case 1: m_shortTimeFormat = "hh:mm"; break;
    default: m_shortTimeFormat = "hh:mm"; break;
    }

    update();
}
