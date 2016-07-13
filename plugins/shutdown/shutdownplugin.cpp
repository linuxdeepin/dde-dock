#include "shutdownplugin.h"

#include <QIcon>

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent),

      m_pluginWidget(new PluginWidget),

      m_powerInter(new DBusPower(this))
{
}

const QString ShutdownPlugin::pluginName() const
{
    return "shutdown";
}

PluginsItemInterface::ItemType ShutdownPlugin::pluginType(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return Complex;
}

 QWidget *ShutdownPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_pluginWidget;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, QString());

    displayModeChanged(qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>());
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    if (displayMode == Dock::Fashion)
        m_icon.addFile(":/icons/resources/icons/fashion.svg");
    else
        m_icon.addFile(":/icons/resources/icons/normal.svg");
}

const QIcon ShutdownPlugin::itemIcon(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_icon;
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

const QString ShutdownPlugin::itemTipsString(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty())
        return QString();

    double percentage = 0.0;
    for (auto percent : data.values())
        percentage += percent;

    if (percentage >= 99.0)
        return QString("100%");
    return QString("%1%").arg(percentage / data.values().count(), 1, 'f', 1);
}
