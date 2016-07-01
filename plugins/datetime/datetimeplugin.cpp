#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),

      m_calendar(new DCalendar(nullptr)),

      m_refershTimer(new QTimer(this))
{
    m_calendar->setFixedSize(300, 300);

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centeralWidget = new DatetimeWidget;

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_centeralWidget;
}

const QString DatetimePlugin::pluginName() const
{
    return "datetime";
}

PluginsItemInterface::ItemType DatetimePlugin::pluginType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Complex;
}

PluginsItemInterface::ItemType DatetimePlugin::tipsType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Complex;
}

void DatetimePlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, QString());
}

int DatetimePlugin::itemSortKey(const QString &itemKey) const
{
    Q_UNUSED(itemKey);

    return -1;
}

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_centeralWidget;
}

QWidget *DatetimePlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_calendar;
}

void DatetimePlugin::updateCurrentTimeString()
{
    const QString currentString = QTime::currentTime().toString("mm");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    m_centeralWidget->update();
}
