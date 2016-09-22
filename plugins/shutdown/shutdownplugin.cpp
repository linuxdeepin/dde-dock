#include "shutdownplugin.h"
#include "dbus/dbusaccount.h"

#include <QIcon>

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
    connect(m_shutdownWidget, &PluginWidget::requestContextMenu, this, &ShutdownPlugin::requestContextMenu);
    connect(m_powerStatusWidget, &PowerStatusWidget::requestContextMenu, this, &ShutdownPlugin::requestContextMenu);
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
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty() || (itemKey == SHUTDOWN_KEY && displayMode() == Dock::Efficient))
    {
        m_tipsLabel->setText(tr("Shut down"));
        return m_tipsLabel;
    }

    const QString value = QString("%1%").arg(data.value("Display"), 0, 'f', 1);
    const bool charging = !m_powerInter->onBattery();
    if (!charging)
        m_tipsLabel->setText(tr("Remaining Capacity %1").arg(value));
    else
    {
        if (m_powerInter->batteryState()["Display"] == 4)
            m_tipsLabel->setText(tr("Charged %1").arg(value));
        else
            m_tipsLabel->setText(tr("Charging %1").arg(value));
    }

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
    if (itemKey == SHUTDOWN_KEY)
        return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
    if (itemKey == POWER_KEY)
        return QString("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");

    return QString();
}

const QString ShutdownPlugin::itemContextMenu(const QString &itemKey)
{
    QList<QVariant> items;
    items.reserve(6);

    const Dock::DisplayMode mode = displayMode();

    if (mode == Dock::Fashion || itemKey == SHUTDOWN_KEY)
    {
        QMap<QString, QVariant> shutdown;
        shutdown["itemId"] = "shutdown";
        shutdown["itemText"] = tr("Shut down");
        shutdown["isActive"] = true;
        items.push_back(shutdown);

        QMap<QString, QVariant> reboot;
        reboot["itemId"] = "reboot";
        reboot["itemText"] = tr("Restart");
        reboot["isActive"] = true;
        items.push_back(reboot);

        QMap<QString, QVariant> logout;
        logout["itemId"] = "logout";
        logout["itemText"] = tr("Log out");
        logout["isActive"] = true;
        items.push_back(logout);

        QMap<QString, QVariant> suspend;
        suspend["itemId"] = "suspend";
        suspend["itemText"] = tr("Suspend");
        suspend["isActive"] = true;
        items.push_back(suspend);

        DBusAccount *accountInter = new DBusAccount(this);
        if (accountInter->userList().count() > 1)
        {
            QMap<QString, QVariant> switchUser;
            switchUser["itemId"] = "switchUser";
            switchUser["itemText"] = tr("Switch account");
            switchUser["isActive"] = true;
            items.push_back(switchUser);
        }
        accountInter->deleteLater();
    }

    if (mode == Dock::Fashion || itemKey == POWER_KEY)
    {
        QMap<QString, QVariant> power;
        power["itemId"] = "power";
        power["itemText"] = tr("Power settings");
        power["isActive"] = true;
        items.push_back(power);
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
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

void ShutdownPlugin::requestContextMenu(const QString &itemKey)
{
    m_proxyInter->requestContextMenu(this, itemKey);
}
