#include "shutdownplugin.h"

#include <QIcon>

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent),

      m_pluginWidget(new PluginWidget),
      m_tipsLabel(new QLabel),

      m_powerInter(new DBusPower(this))
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setAlignment(Qt::AlignCenter);
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding:5px 10px;");
}

const QString ShutdownPlugin::pluginName() const
{
    return "shutdown";
}

QWidget *ShutdownPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_pluginWidget;
}

QWidget *ShutdownPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty())
        return nullptr;

    m_tipsLabel->setText(QString("%1%").arg(data.value("Display"), 1, 'f', 1));

    return m_tipsLabel;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    m_proxyInter->itemAdded(this, QString());

    displayModeChanged(qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>());
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    m_pluginWidget->update();
}
