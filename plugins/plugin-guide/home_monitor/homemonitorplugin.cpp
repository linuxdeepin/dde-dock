#include "homemonitorplugin.h"

HomeMonitorPlugin::HomeMonitorPlugin(QObject *parent)
    : QObject(parent)
{

}

const QString HomeMonitorPlugin::pluginDisplayName() const
{
    return QString("Home Monitor");
}

const QString HomeMonitorPlugin::pluginName() const
{
    return QStringLiteral("home_monitor");
}

void HomeMonitorPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_pluginWidget = new InformationWidget;
    m_tipsWidget = new QLabel;
    m_appletWidget = new QLabel;

    // 如果插件没有被禁用则在初始化插件时才添加主控件到面板上
    if (!pluginIsDisable()) {
        m_proxyInter->itemAdded(this, pluginName());
    }
}

QWidget *HomeMonitorPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_pluginWidget;
}

QWidget *HomeMonitorPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    // 设置/刷新 tips 中的信息
    m_tipsWidget->setText(QString("Total: %1G\nAvailable: %2G")
                          .arg(qRound(m_pluginWidget->storageInfo()->bytesTotal() / qPow(1024, 3)))
                          .arg(qRound(m_pluginWidget->storageInfo()->bytesAvailable() / qPow(1024, 3))));

    return m_tipsWidget;
}

QWidget *HomeMonitorPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    m_appletWidget->setText(QString("Total: %1G\nAvailable: %2G\nDevice: %3\nVolume: %4\nLabel: %5\nFormat: %6\nAccess: %7")
                            .arg(qRound(m_pluginWidget->storageInfo()->bytesTotal() / qPow(1024, 3)))
                            .arg(qRound(m_pluginWidget->storageInfo()->bytesAvailable() / qPow(1024, 3)))
                            .arg(QString(m_pluginWidget->storageInfo()->device()))
                            .arg(m_pluginWidget->storageInfo()->displayName())
                            .arg(m_pluginWidget->storageInfo()->name())
                            .arg(QString(m_pluginWidget->storageInfo()->fileSystemType()))
                            .arg(m_pluginWidget->storageInfo()->isReadOnly() ? "ReadOnly" : "ReadWrite")
                            );

    return m_appletWidget;
}

bool HomeMonitorPlugin::pluginIsAllowDisable()
{
    // 告诉 dde-dock 本插件允许禁用
    return true;
}

bool HomeMonitorPlugin::pluginIsDisable()
{
    // 第二个参数 “disabled” 表示存储这个值的键（所有配置都是以键值对的方式存储的）
    // 第三个参数表示默认值，即默认不禁用
    return m_proxyInter->getValue(this, "disabled", false).toBool();
}

void HomeMonitorPlugin::pluginStateSwitched()
{
    // 获取当前禁用状态的反值作为新的状态值
    const bool disabledNew = !pluginIsDisable();
    // 存储新的状态值
    m_proxyInter->saveValue(this, "disabled", disabledNew);

    // 根据新的禁用状态值处理主控件的加载和卸载
    if (disabledNew) {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        m_proxyInter->itemAdded(this, pluginName());
    }
}

const QString HomeMonitorPlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> refresh;
    refresh["itemId"] = "refresh";
    refresh["itemText"] = "Refresh";
    refresh["isActive"] = true;
    items.push_back(refresh);

    QMap<QString, QVariant> open;
    open["itemId"] = "open";
    open["itemText"] = "Open Gparted";
    open["isActive"] = true;
    items.push_back(open);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    // 返回 JSON 格式的菜单数据
    return QJsonDocument::fromVariant(menu).toJson();
}

void HomeMonitorPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey);

    // 根据上面接口设置的 id 执行不同的操作
    if (menuId == "refresh") {
        m_pluginWidget->storageInfo()->refresh();
    } else if ("open") {
        QProcess::startDetached("gparted");
    }
}
