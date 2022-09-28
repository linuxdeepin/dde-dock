// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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

DWIDGET_USE_NAMESPACE

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
    , m_24HourFormat(false)
    , m_longDateFormatType(0)
    , m_weekdayFormatType(0)
    , m_timeOffset(false)
    , m_timedateInter(new Timedate("com.deepin.daemon.Timedate", "/com/deepin/daemon/Timedate", QDBusConnection::sessionBus(), this))
    , m_shortDateFormat("yyyy-MM-dd")
    , m_shortTimeFormat("hh:mm")
    , m_longTimeFormat(" hh:mm:ss")
{
    auto font = this->font();
    font.setPixelSize(10);
    this->setFont(font);
    m_dateFont = font;
    m_timeSize = this->font();
    m_timeSize.setPixelSize(20);
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);
    setShortDateFormat(m_timedateInter->shortDateFormat());
    setShortTimeFormat(m_timedateInter->shortTimeFormat());
    setWeekdayFormat(m_timedateInter->weekdayFormat());
    setLongDateFormat(m_timedateInter->longDateFormat());
    setLongTimeFormat(m_timedateInter->longTimeFormat());
    set24HourFormat(m_timedateInter->use24HourFormat());
    updateDateTimeString();

    connect(m_timedateInter, &Timedate::ShortDateFormatChanged, this, &DatetimeWidget::setShortDateFormat);
    connect(m_timedateInter, &Timedate::ShortTimeFormatChanged, this, &DatetimeWidget::setShortTimeFormat);
    connect(m_timedateInter, &Timedate::LongDateFormatChanged, this, &DatetimeWidget::setLongDateFormat);
    connect(m_timedateInter, &Timedate::WeekdayFormatChanged, this, &DatetimeWidget::setWeekdayFormat);
    connect(m_timedateInter, &Timedate::LongTimeFormatChanged, this, &DatetimeWidget::setLongTimeFormat);
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
    updateLongTimeFormat();
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
    case 6: m_shortDateFormat = "MM.dd.yyyy"; break;
    case 7: m_shortDateFormat = "dd.MM.yyyy"; break;
    case 8: m_shortDateFormat = "yy/M/d"; break;
    case 9: m_shortDateFormat = "yy-M-d"; break;
    case 10: m_shortDateFormat = "yy.M.d"; break;
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

/**
 * @brief DatetimeWidget::setLongDateFormat 根据类型设置长时间显示格式
 * @param type 自定义类型
 */
void DatetimeWidget::setLongDateFormat(int type)
{
    if (m_longDateFormatType == type)
        return;

    m_longDateFormatType = type;
    updateDateTimeString();
}

/**
 * @brief DatetimeWidget::setWeekdayFormat 根据类型设置周显示格式
 * @param type 自定义类型
 */
void DatetimeWidget::setWeekdayFormat(int type)
{
    if (m_weekdayFormatType == type)
        return;

    m_weekdayFormatType = type;
    updateWeekdayFormat();
    updateDateTimeString();
}

/**
 * @brief DatetimeWidget::setLongTimeFormat 根据类型设置长时间的显示格式
 * @param type 自定义类型
 */
void DatetimeWidget::setLongTimeFormat(int type)
{
    if (m_longTimeFormatType == type)
        return;

    m_longTimeFormatType = type;
    updateLongTimeFormat();
    updateDateTimeString();
}

/**
 * @brief DatetimeWidget::updateWeekdayFormat 更新周的显示格式
 */
void DatetimeWidget::updateWeekdayFormat()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();
    auto dayOfWeek = currentDateTime.date().dayOfWeek();

    if (0 == m_weekdayFormatType) {
        switch (dayOfWeek) {
        case 1:
            m_weekFormat = tr("Monday"); //星期一
            break;
        case 2:
            m_weekFormat = tr("Tuesday"); //星期二
            break;
        case 3:
            m_weekFormat = tr("Wednesday"); //星期三
            break;
        case 4:
            m_weekFormat = tr("Thursday"); //星期四
            break;
        case 5:
            m_weekFormat = tr("Friday"); //星期五
            break;
        case 6:
            m_weekFormat = tr("Saturday"); //星期六
            break;
        case 7:
            m_weekFormat = tr("Sunday"); //星期天
            break;
        default:
            m_weekFormat = tr("Monday"); //星期一
            break;
        }
    } else {
        switch (dayOfWeek) {
        case 1:
            m_weekFormat = tr("monday"); //周一
            break;
        case 2:
            m_weekFormat = tr("tuesday"); //周二
            break;
        case 3:
            m_weekFormat = tr("wednesday"); //周三
            break;
        case 4:
            m_weekFormat = tr("thursday"); //周四
            break;
        case 5:
            m_weekFormat = tr("friday"); //周五
            break;
        case 6:
            m_weekFormat = tr("saturday"); //周六
            break;
        case 7:
            m_weekFormat = tr("sunday"); //周天
            break;
        default:
            m_weekFormat = tr("monday"); //周一
            break;
        }
    }
}

void DatetimeWidget::updateLongTimeFormat()
{
    if (m_24HourFormat) {
        switch (m_longTimeFormatType) {
        case 0: m_longTimeFormat = " h:mm:ss"; break;
        case 1: m_longTimeFormat = " hh:mm:ss";  break;
        default: m_longTimeFormat = " hh:mm:ss"; break;
        }
    } else {
        switch (m_longTimeFormatType) {
        case 0: m_longTimeFormat = " h:mm:ss A"; break;
        case 1: m_longTimeFormat = " hh:mm:ss A";  break;
        default: m_longTimeFormat = " hh:mm:ss A"; break;
        }
    }
}

/**
 * @brief DatetimeWidget::updateWeekdayTimeString 更新任务栏时间标签的显示
 */
void DatetimeWidget::updateDateTimeString()
{
    QString longTimeFormat("");
    const QDateTime currentDateTime = QDateTime::currentDateTime();
    int year = currentDateTime.date().year();
    int month = currentDateTime.date().month();
    int day = currentDateTime.date().day();

    auto lang = QLocale::system().language();
    bool isZhLocale = lang == QLocale::Chinese || lang == QLocale::Tibetan || lang == QLocale::Uighur;

    // 根据相应语言去显示对应的格式
    // 中文： 格式为xxxx年xx月xx日 星期x hh:mm:ss,如:2022年7月25日 星期- 12：00：00
    // 英文： 格式为x x，xxxx，x hh:mm:ss, 如：July 25，2022，Monday 12:00:00
    // 其他语言：按照国际当地长时间格式显示
    if (isZhLocale) {
        longTimeFormat = QString(tr("%1year%2month%3day")).arg(year).arg(month).arg(day);

        // 实时更新周的日期显示
        updateWeekdayFormat();

        switch (m_longDateFormatType) {
        case 0:
            m_dateTime = longTimeFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        case 1:
            m_dateTime = longTimeFormat + QString(" ") + m_weekFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        case 2:
            m_dateTime = m_weekFormat + QString(" ") + longTimeFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        default:
            m_dateTime = longTimeFormat + QString(" ") + m_weekFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        }
    } else if (lang == QLocale::English) {
        auto longDateString = currentDateTime.date().toString(Qt::SystemLocaleLongDate);
        auto week = longDateString.split(",").at(0);
        // 获取英文的日期格式字符串，-2是去掉","和" "
        auto longDateTimeFormat = longDateString.right(longDateString.size() - week.size() - 2);

        switch (m_longDateFormatType) {
        case 0:
            m_dateTime = longDateTimeFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        case 1:
            m_dateTime = longDateTimeFormat + QString(", ") + week + currentDateTime.toString(m_longTimeFormat);
            break;
        case 2:
            m_dateTime = week + QString(", ") + longDateTimeFormat + currentDateTime.toString(m_longTimeFormat);
            break;
        default:
            m_dateTime = longDateTimeFormat + QString(", ") + week + currentDateTime.toString(m_longTimeFormat);
            break;
        }
    } else {
        m_dateTime = currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(m_longTimeFormat);
    }
}

/**
 * @brief DatetimeWidget::curTimeSize 调整时间日期字体大小
 * @return 返回时间和日期绘制的区域大小
 */
QSize DatetimeWidget::curTimeSize() const
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

    m_timeFont = m_timeSize;
    m_dateFont = m_dateSize;
    QString timeFormat = m_shortTimeFormat;
    QString dateFormat = m_shortDateFormat;
    if (!m_24HourFormat) {
        if (position == Dock::Top || position == Dock::Bottom)
            timeFormat = timeFormat.append(" AP");
        else
            timeFormat = timeFormat.append("\nAP");
    }

    QString timeString = QDateTime::currentDateTime().toString(timeFormat);
    QString dateString = QDateTime::currentDateTime().toString(dateFormat);

    QSize timeSize = QFontMetrics(m_timeFont).boundingRect(timeString).size();
    int maxWidth = std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_timeFont).horizontalAdvance(timeString));
    timeSize.setWidth(maxWidth);

    if (timeString.contains("\n")) {
        QStringList SL = timeString.split("\n");
        maxWidth = std::max(QFontMetrics(m_timeFont).boundingRect(SL.at(0)).size().width(), QFontMetrics(m_timeFont).horizontalAdvance(SL.at(0)));
        timeSize = QSize(maxWidth, QFontMetrics(m_timeFont).boundingRect(SL.at(0)).height() + QFontMetrics(m_timeFont).boundingRect(SL.at(1)).height());
    }

    QSize dateSize = QFontMetrics(m_dateFont).boundingRect(dateString).size();
    maxWidth = std::max(QFontMetrics(m_dateFont).boundingRect(dateString).size().width(), QFontMetrics(m_dateFont).horizontalAdvance(dateString));
    dateSize.setWidth(maxWidth);

    if (position == Dock::Bottom || position == Dock::Top) {
        while (QFontMetrics(m_timeFont).boundingRect(timeString).height() + QFontMetrics(m_dateFont).boundingRect(dateString).height() > height() && m_timeFont.pixelSize() > 1) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            maxWidth = std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_timeFont).horizontalAdvance(timeString));
            timeSize.setWidth(maxWidth);
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1) {
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                maxWidth = std::max(QFontMetrics(m_dateFont).boundingRect(dateString).size().width(), QFontMetrics(m_dateFont).horizontalAdvance(dateString));
                dateSize.setWidth(maxWidth);
            }
        }
        return QSize(std::max(timeSize.width(), dateSize.width()), timeSize.height() + dateSize.height());
    } else {
        while (std::max(QFontMetrics(m_timeFont).boundingRect(timeString).size().width(), QFontMetrics(m_dateFont).boundingRect(dateString).size().width()) > (width() - 4) && m_timeFont.pixelSize() > 1) {
            m_timeFont.setPixelSize(m_timeFont.pixelSize() - 1);
            if (m_24HourFormat) {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height());
            } else {
                timeSize.setHeight(QFontMetrics(m_timeFont).boundingRect(timeString).size().height() * 2);
            }
            if (m_timeFont.pixelSize() - m_dateFont.pixelSize() == 1) {
                m_dateFont.setPixelSize(m_dateFont.pixelSize() - 1);
                dateSize.setWidth(QFontMetrics(m_dateFont).boundingRect(dateString).size().height());
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

void DatetimeWidget::resizeEvent(QResizeEvent *event)
{
    if (isVisible())
        emit requestUpdateGeometry();

    QWidget::resizeEvent(event);
}

/**
 * @brief DatetimeWidget::paintEvent 绘制任务栏时间日期
 * @param e
 */
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
        // 只处理上下位置的，特殊处理一下藏文，其他的语言如果有问题也可以类似特殊处理一下
        // Unifont字体有点特殊
        // 以下的0.23 0.18 0.2 0.13数值是测试过程中微调时间跟日期之间的间距系数，不是特别计算的精确值
        QLocale locale;
        int timeHeight = QFontMetrics(m_timeFont).boundingRect(timeStr).height() + 2;   // +2只是防止显示在边界的几个像素被截断
        int dateHeight = QFontMetrics(m_dateFont).boundingRect(dateStr).height() + 2;
        int marginH = (height() - timeHeight - dateHeight) / 2;

        if (locale.language() == QLocale::Tibetan) {
            if (m_timeFont.family() == "Noto Serif Tibetan")
                marginH = marginH + 0.23 * timeHeight;
            else if (m_timeFont.family() == "Noto Sans Tibetan")
                marginH = marginH + 0.18 * timeHeight;
            else if (m_timeFont.family() == "Tibetan Machine Uni")
                marginH = marginH + 0.2 * timeHeight;
        } else {
            if (m_timeFont.family() != "Unifont")
                marginH = marginH + 0.13 * timeHeight;
        }

        timeRect = QRect(0, marginH, width(), timeHeight);
        dateRect = QRect(0, height() - dateHeight - marginH, width(), dateHeight);
    } else {
        timeRect.setBottom(rect().center().y() + m_timeOffset);
        dateRect.setTop(timeRect.bottom());
    }
    painter.setFont(m_timeFont);
    painter.drawText(timeRect, Qt::AlignCenter, timeStr);

    painter.setFont(m_dateFont);
    painter.drawText(dateRect, Qt::AlignCenter, dateStr);
}
