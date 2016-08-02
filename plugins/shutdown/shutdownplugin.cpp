#include "shutdownplugin.h"

#include <QIcon>

#define POWER_KEY       "power"
#define SHUTDOWN_KEY    "shutdown"

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent),

      m_shutdownWidget(new PluginWidget),
      m_powerStatusWidget(new PowerStatusWidget),
      m_tipsLabel(new QLabel),

      m_powerInter(new DBusPower(this))
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("power");
    m_tipsLabel->setAlignment(Qt::AlignCenter);
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding:5px 10px;");

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &ShutdownPlugin::updateBatteryVisible);
}

const QString ShutdownPlugin::pluginName() const
{
    return "shutdown";
}

QWidget *ShutdownPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == SHUTDOWN_KEY)
        return m_shutdownWidget;
    if (itemKey == POWER_KEY)
        return m_powerStatusWidget;

    return nullptr;
}

QWidget *ShutdownPlugin::itemTipsWidget(const QString &itemKey)
{
    const Dock::DisplayMode mode = displayMode();

    if (mode == Dock::Efficient && itemKey == SHUTDOWN_KEY)
        return nullptr;

    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty())
        return nullptr;

    m_tipsLabel->setText(QString("%1%").arg(data.value("Display"), 1, 'f', 1));

    return m_tipsLabel;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, SHUTDOWN_KEY);

    displayModeChanged(displayMode());
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    const Dock::DisplayMode mode = displayMode();

    if (mode == Dock::Efficient && itemKey != SHUTDOWN_KEY)
        return QString();

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    m_shutdownWidget->update();

    updateBatteryVisible();
}

void ShutdownPlugin::updateBatteryVisible()
{
    const bool exist = !m_powerInter->batteryPercentage().isEmpty();

    if (!exist || displayMode() == Dock::Fashion)
        m_proxyInter->itemRemoved(this, POWER_KEY);
    else if (exist)
        m_proxyInter->itemAdded(this, POWER_KEY);
}
