#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),
      m_timeLabel(new QLabel),
      m_refershTimer(new QTimer(this))
{
    m_timeLabel->setAlignment(Qt::AlignCenter);
    m_timeLabel->setStyleSheet("color:white;"
                               "background-color:black;"
                               "font-size:12px;");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::refershTime);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_timeLabel;
}

const QString DatetimePlugin::name()
{
    return "datetime";
}

QWidget *DatetimePlugin::centeralWidget()
{
    return m_timeLabel;
}

void DatetimePlugin::refershTime()
{
    const QString text = QTime::currentTime().toString(tr("HH:mm"));

    if (m_timeLabel->text() != text)
        m_timeLabel->setText(text);
}
