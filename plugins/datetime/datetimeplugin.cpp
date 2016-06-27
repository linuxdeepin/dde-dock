#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),
      m_refershTimer(new QTimer(this))
{
    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centeralWidget = new DatetimeWidget;

    connect(m_refershTimer, &QTimer::timeout, m_centeralWidget, static_cast<void (DatetimeWidget::*)()>(&DatetimeWidget::update));
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

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_centeralWidget;
}
