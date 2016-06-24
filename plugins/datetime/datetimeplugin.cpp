#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),
      m_timeLabel(new QLabel),
      m_refershTimer(new QTimer(this))
{
    m_timeLabel->setAlignment(Qt::AlignCenter);
    m_timeLabel->setStyleSheet("color:white;"
//                               "background-color:black;"
                               "padding:5px;"
                               "font-size:12px;");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::refershTime);

    refershTime();
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_timeLabel;
}

const QString DatetimePlugin::pluginName()
{
    return "datetime";
}

PluginsItemInterface::PluginType DatetimePlugin::pluginType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Complex;
}

void DatetimePlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, QString());
}

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_timeLabel;
}

void DatetimePlugin::refershTime()
{
    const QString text = QTime::currentTime().toString(tr("HH:mm"));

    if (m_timeLabel->text() != text)
        m_timeLabel->setText(text);
}
