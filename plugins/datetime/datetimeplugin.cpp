#include "datetimeplugin.h"

#include <QLabel>

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),

      m_dateTipsLabel(new QLabel),

      m_refershTimer(new QTimer(this))
{
    m_dateTipsLabel->setObjectName("datetime");
    m_dateTipsLabel->setStyleSheet("color:white;"
                                   "padding:6px 10px;");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centeralWidget = new DatetimeWidget;

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_centeralWidget;
    delete m_dateTipsLabel;
}

const QString DatetimePlugin::pluginName() const
{
    return "datetime";
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

    return m_dateTipsLabel;
}

const QString DatetimePlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return "dde-calendar";
}

void DatetimePlugin::updateCurrentTimeString()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();

    m_dateTipsLabel->setText(currentDateTime.toString(tr("MM/dd/yyyy ddd HH:mm:ss")));

    const QString currentString = currentDateTime.toString("mm");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    m_centeralWidget->update();
}
