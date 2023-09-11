// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "datetimedisplayer.h"
#include "tipswidget.h"
#include "dockpopupwindow.h"
#include "utils.h"
#include "dbusutil.h"

#include <DFontSizeManager>
#include <DDBusSender>
#include <DGuiApplicationHelper>

#include <QHBoxLayout>
#include <QPainter>
#include <QFont>
#include <QMenu>
#include <QPainterPath>
#include <QMouseEvent>
#include <QFontMetrics>

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

#define DATETIMESIZE 40
#define ITEMSPACE 8

static QMap<int, QString> dateFormat{{ 0,"yyyy/M/d" }, { 1,"yyyy-M-d" }, { 2,"yyyy.M.d" }, { 3,"yyyy/MM/dd" },
                                     { 4,"yyyy-MM-dd" }, { 5,"yyyy.MM.dd" }, { 6,"yy/M/d" }, { 7,"yy-M-d" }, { 8,"yy.M.d" }};
static QMap<int, QString> timeFormat{{0, "h:mm"}, {1, "hh:mm"}};

DateTimeDisplayer::DateTimeDisplayer(bool showMultiRow, QWidget *parent)
    : QWidget (parent)
    , m_timedateInter(new Timedate("org.deepin.dde.Timedate1", "/org/deepin/dde/Timedate1", QDBusConnection::sessionBus(), this))
    , m_position(Dock::Position::Bottom)
    , m_dateFont(QFont())
    , m_timeFont(QFont())
    , m_tipsWidget(new Dock::TipsWidget(this))
    , m_menu(new QMenu(this))
    , m_tipsTimer(new QTimer(this))
    , m_currentSize(0)
    , m_oneRow(false)
    , m_showMultiRow(showMultiRow)
{
    m_tipPopupWindow.reset(new DockPopupWindow);
    // 日期格式变化的时候，需要重绘
    connect(m_timedateInter, &Timedate::ShortDateFormatChanged, this, &DateTimeDisplayer::onDateTimeFormatChanged);
    // 时间格式变化的时候，需要重绘
    connect(m_timedateInter, &Timedate::ShortTimeFormatChanged, this, &DateTimeDisplayer::onDateTimeFormatChanged);
    // 是否使用24小时制发生变化的时候，也需要重绘
    connect(m_timedateInter, &Timedate::Use24HourFormatChanged, this, &DateTimeDisplayer::onDateTimeFormatChanged);
    // 连接日期时间修改信号,更新日期时间插件的布局
    connect(m_timedateInter, &Timedate::TimeUpdate, this, static_cast<void (QWidget::*)()>(&DateTimeDisplayer::update));
    // 连接定时器和时间显示的tips信号,一秒钟触发一次，显示时间
    connect(m_tipsTimer, &QTimer::timeout, this, &DateTimeDisplayer::onTimeChanged);
    QMetaObject::invokeMethod(this, "onDateTimeFormatChanged");
    m_tipsTimer->setInterval(1000);
    m_tipsTimer->start();
    updatePolicy();
    createMenuItem();
    if (Utils::IS_WAYLAND_DISPLAY)
        m_tipPopupWindow->setWindowFlags(m_tipPopupWindow->windowFlags() | Qt::FramelessWindowHint);
    m_tipPopupWindow->hide();
}

DateTimeDisplayer::~DateTimeDisplayer()
{
}

void DateTimeDisplayer::setPositon(Dock::Position position)
{
    if (m_position == position)
        return;

    m_position = position;
    updatePolicy();
    update();
}

void DateTimeDisplayer::setOneRow(bool oneRow)
{
    m_oneRow = oneRow;
    update();
}

void DateTimeDisplayer::updatePolicy()
{
    switch(m_position) {
    case Dock::Position::Top:
    case Dock::Position::Bottom:
        setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Expanding);
        break;
    case Dock::Position::Left:
    case Dock::Position::Right:
        setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
        break;
    }

    m_tipPopupWindow->setPosition(m_position);
    m_tipPopupWindow->setContent(m_tipsWidget);
}

QSize DateTimeDisplayer::suitableSize() const
{
    return suitableSize(m_position);
}

QSize DateTimeDisplayer::suitableSize(const Dock::Position &position) const
{
    DateTimeInfo info = dateTimeInfo(position);
    if (position == Dock::Position::Left || position == Dock::Position::Right)
        return QSize(width(), info.m_timeRect.height() + info.m_dateRect.height());

    // 如果在上下显示
    if (m_showMultiRow) {
        // 如果显示多行的情况,一般是在高效模式下显示，因此，返回最大的尺寸
        return QSize(qMax(info.m_timeRect.width(), info.m_dateRect.width()), height());
    }

    return QSize(info.m_timeRect.width() + info.m_dateRect.width() + 16, height());
}

void DateTimeDisplayer::mousePressEvent(QMouseEvent *event)
{
    if ((event->button() != Qt::RightButton))
        return QWidget::mousePressEvent(event);

    m_menu->exec(QCursor::pos());
}

void DateTimeDisplayer::mouseReleaseEvent(QMouseEvent *event)
{
    Q_UNUSED(event);

    DDBusSender().service("org.deepin.dde.Widgets1")
            .path("/org/deepin/dde/Widgets1")
            .interface("org.deepin.dde.Widgets1")
            .method("Toggle").call();
}

QString DateTimeDisplayer::getTimeString(const Dock::Position &position) const
{
    QString tFormat = QString("hh:mm");
    if (timeFormat.contains(m_shortDateFormat))
        tFormat = timeFormat[m_shortDateFormat];

    if (!m_use24HourFormat) {
        if (position == Dock::Top || position == Dock::Bottom)
            tFormat = tFormat.append(" AP");
        else
            tFormat = tFormat.append("\nAP");
    }

    return QDateTime::currentDateTime().toString(tFormat);
}

QString DateTimeDisplayer::getDateString() const
{
    return getDateString(m_position);
}

QString DateTimeDisplayer::getDateString(const Dock::Position &position) const
{
    QString shortDateFormat = "yyyy-MM-dd";
    if (dateFormat.contains(m_shortDateFormat))
        shortDateFormat = dateFormat.value(m_shortDateFormat);
    // 如果是左右方向，则不显示年份
    if (position == Dock::Position::Left || position == Dock::Position::Right) {
        static QStringList yearStrList{"yyyy/", "yyyy-", "yyyy.", "yy/", "yy-", "yy."};
        for (int i = 0; i < yearStrList.size() ; i++) {
            const QString &yearStr = yearStrList[i];
            if (shortDateFormat.contains(yearStr)) {
                shortDateFormat = shortDateFormat.remove(yearStr);
                break;
            }
        }
    }

    return QDateTime::currentDateTime().toString(shortDateFormat);
}

DateTimeDisplayer::DateTimeInfo DateTimeDisplayer::dateTimeInfo(const Dock::Position &position) const
{
    DateTimeInfo info;
    info.m_timeRect = rect();
    info.m_dateRect = rect();

    info.m_time = getTimeString(position);
    info.m_date = getDateString(position);

    // 如果是左右方向
    if (position == Dock::Position::Left || position == Dock::Position::Right) {
        int textWidth = rect().width();

        int timeHeight = QFontMetrics(m_timeFont).boundingRect(info.m_time).height() * (info.m_time.count('\n') + 1);
        int dateHeight = QFontMetrics(m_dateFont).boundingRect(info.m_date).height();

        info.m_timeRect = QRect(0, 0, textWidth, timeHeight);
        info.m_dateRect = QRect(0, timeHeight, textWidth, dateHeight);
        return info;
    }
    int timeWidth = QFontMetrics(m_timeFont).boundingRect(info.m_time).width() + 2;
    int dateWidth = QFontMetrics(m_dateFont).boundingRect(info.m_date).width() + 2;

    int rHeight =  height();

    // 如果是上下方向
    if (m_showMultiRow) {
        // 日期时间多行显示（一般是高效模式下，向下和向上偏移2个像素）
        info.m_timeRect = QRect(0, 2, timeWidth, rHeight / 2);
        info.m_dateRect = QRect(0, rHeight / 2 - 2, dateWidth, rHeight / 2);
    } else {
        // 3:时间和日期3部分间隔
        if (rect().width() > (ITEMSPACE * 3 + timeWidth + dateWidth)) {
            info.m_timeRect = QRect(ITEMSPACE, 0, timeWidth, rHeight);

            int dateX = info.m_timeRect.right() + (rect().width() -(ITEMSPACE * 2 + timeWidth + dateWidth));
            info.m_dateRect = QRect(dateX, 0, dateWidth, rHeight);
        } else {
            // 宽度不满足间隔为ITEMSPACE的，需要自己计算间隔。
            int itemSpace = (rect().width() - timeWidth - dateWidth) / 3;
            info.m_timeRect = QRect(itemSpace, 0, timeWidth, rHeight);

            int dateX = info.m_timeRect.right() + itemSpace;
            info.m_dateRect = QRect(dateX, 0, dateWidth, rHeight);
        }
    }

    return info;
}

void DateTimeDisplayer::onTimeChanged()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();

    if (m_use24HourFormat)
        m_tipsWidget->setText(QLocale().toString(currentDateTime.date()) + currentDateTime.toString(" HH:mm:ss"));
    else
        m_tipsWidget->setText(QLocale().toString(currentDateTime.date()) + currentDateTime.toString(" hh:mm:ss AP"));

    // 如果时间和日期有一个不等，则实时刷新界面
    if (m_lastDateString != getDateString() || m_lastTimeString != getTimeString())
        update();
}

void DateTimeDisplayer::onDateTimeFormatChanged()
{
    m_shortDateFormat = m_timedateInter->shortDateFormat();
    m_use24HourFormat = m_timedateInter->use24HourFormat();
    // 此处需要强制重绘，因为在重绘过程中才会改变m_currentSize信息，方便在后面判断是否需要调整尺寸
    repaint();
}

void DateTimeDisplayer::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    DateTimeInfo info = dateTimeInfo(m_position);

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    

    int timeAlignFlag = Qt::AlignCenter;
    int dateAlignFlag = Qt::AlignCenter;

    if (m_showMultiRow) {
        timeAlignFlag = Qt::AlignHCenter | Qt::AlignBottom;
        dateAlignFlag = Qt::AlignHCenter | Qt::AlignTop;
    }

    painter.setFont(m_timeFont);
    painter.setPen(QPen(palette().brightText(), 2));
    painter.drawText(textRect(info.m_timeRect), timeAlignFlag, info.m_time);

    painter.setFont(m_dateFont);
    painter.setPen(QPen(palette().brightText(), 1));
    painter.drawText(textRect(info.m_dateRect), dateAlignFlag, info.m_date);

    updateLastData(info);
}

QPoint DateTimeDisplayer::tipsPoint() const
{
    QPoint pointInTopWidget = parentWidget()->mapTo(window(), pos());
    switch (m_position) {
    case Dock::Position::Left: {
        pointInTopWidget.setX(window()->x() + window()->width());
        pointInTopWidget.setY(pointInTopWidget.y() + height() / 2);
        break;
    }
    case Dock::Position::Top: {
        pointInTopWidget.setY(y() + window()->y() + window()->height());
        pointInTopWidget.setX(pointInTopWidget.x() + width() / 2);
        break;
    }
    case Dock::Position::Right: {
        pointInTopWidget.setY(pointInTopWidget.y() + height() / 2);
        pointInTopWidget.setX(pointInTopWidget.x() - width() / 2);
        break;
    }
    case Dock::Position::Bottom: {
        pointInTopWidget.setY(-POPUP_PADDING);
        pointInTopWidget.setX(pointInTopWidget.x() + width() / 2);
        break;
    }
    }
    return window()->mapToGlobal(pointInTopWidget);
}

void DateTimeDisplayer::updateFont() const
{
    auto info = getTimeString(m_position);
    // "xx:xx\nAP" 获取到前 xx:xx 部分
    info = info.left(info.indexOf('\n'));
    if (m_position == Dock::Position::Left || m_position == Dock::Position::Right) {
        auto f = QFont();
        bool caled = false;
        f.setPixelSize(100);
        // 左右时根据获取可以全部显示文本的最小的宽度, 且最大只到40
        while(width() > 0 && f.pixelSize() > 2 &&
                (QFontMetrics(f).boundingRect(info).width() > qMin(DATETIMESIZE, width()) - 4)) {
            f.setPixelSize(f.pixelSize() - 1);
            caled = true;
        }
        // 经过正确的计算后才能更新字体大小
        if (caled) {
            m_timeFont.setPixelSize(f.pixelSize());
            m_dateFont.setPixelSize(f.pixelSize() - 2);
        }
        return;
    }

    if ((Dock::Position::Top == m_position || Dock::Position::Bottom == m_position )) {
        // 单行时保持高度的一半，双行时尽量和高度一致，但最大只到12。
        auto s = height() / (m_oneRow ? 2 : 1) - 2;
        m_timeFont.setPixelSize(std::min(s, 12));
        // 双行时日期比时间字体小两个像素。
        m_dateFont.setPixelSize(std::min(s, 12) - (m_oneRow ? 0 : 2));
    }
}

void DateTimeDisplayer::createMenuItem()
{
    QAction *timeFormatAction = new QAction(this);
    timeFormatAction->setText(m_use24HourFormat ? tr("12-hour time"): tr("24-hour time"));

    connect(timeFormatAction, &QAction::triggered, this, [ = ] {
        bool use24hourformat = !m_use24HourFormat;
        // 此时调用 dbus 更新时间格式但是本地 m_use24HourFormat 未更新，所以需要使用新变量，设置新格式
        m_timedateInter->setUse24HourFormat(use24hourformat);
        timeFormatAction->setText(use24hourformat ? tr("12-hour time") : tr("24-hour time"));
    });
    m_menu->addAction(timeFormatAction);

    if (!QFile::exists(ICBC_CONF_FILE)) {
        QAction *timeSettingAction = new QAction(tr("Time settings"), this);
        connect(timeSettingAction, &QAction::triggered, this, [ = ] {
            DDBusSender()
                    .service(controllCenterService)
                    .path(controllCenterPath)
                    .interface(controllCenterInterface)
                    .method(QString("ShowPage"))
                    .arg(QString("datetime"))
                    .call();
        });

        m_menu->addAction(timeSettingAction);
    }
}

QRect DateTimeDisplayer::textRect(const QRect &sourceRect) const
{
    // 如果是上下，则不做任何变化
    if (!m_showMultiRow && (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom))
        return sourceRect;

    QRect resultRect = sourceRect;
    QSize size = suitableSize();
    // 如果是左右或者上下多行显示，设置宽度
    resultRect.setWidth(size.width());
    return resultRect;
}

void DateTimeDisplayer::enterEvent(QEvent *event)
{
    Q_UNUSED(event);
    Q_EMIT requestDrawBackground(rect());
    update();
    m_tipPopupWindow->show(tipsPoint());
}

void DateTimeDisplayer::leaveEvent(QEvent *event)
{
    Q_UNUSED(event);
    Q_EMIT requestDrawBackground(QRect());
    update();
    m_tipPopupWindow->hide();
}

QString DateTimeDisplayer::getTimeString() const
{
    return getTimeString(m_position);
}

void DateTimeDisplayer::updateLastData(const DateTimeInfo &info)
{
    int lastSize = m_currentSize;
    m_lastDateString = info.m_date;
    m_lastTimeString = info.m_time;
    QSize dateTimeSize = suitableSize();
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        m_currentSize = dateTimeSize.width();
    else
        m_currentSize = dateTimeSize.height();
    // 如果日期时间的格式发生了变化，需要通知外部来调整日期时间的尺寸
    if (lastSize != m_currentSize)
        Q_EMIT requestUpdate();
}

bool DateTimeDisplayer::event(QEvent *event)
{
    if (event->type() == QEvent::Resize) {
        updateFont();
    }
    return QWidget::event(event);
}
