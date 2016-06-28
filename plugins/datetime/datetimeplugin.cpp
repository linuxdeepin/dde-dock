#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),
      m_refershTimer(new QTimer(this))
{
    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centeralWidget = new DatetimeWidget;

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_centeralWidget;
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

void DatetimePlugin::updateCurrentTimeString()
{
    const QString currentString = QTime::currentTime().toString("mm");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    m_centeralWidget->update();
}
