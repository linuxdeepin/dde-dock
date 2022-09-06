// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "multitaskingplugin.h"
#include "../widgets/tipswidget.h"

#include <DWindowManagerHelper>
#include <DDBusSender>

#include <QIcon>

#define PLUGIN_STATE_KEY    "enable"
DGUI_USE_NAMESPACE

using namespace Dock;
MultitaskingPlugin::MultitaskingPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_multitaskingWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("multitasking");

    connect(DWindowManagerHelper::instance(), &DWindowManagerHelper::hasCompositeChanged, this, [ = ] {
        if (!m_proxyInter || !m_pluginLoaded)
            return;

        if (DWindowManagerHelper::instance()->hasComposite() && !pluginIsDisable())
            m_proxyInter->itemAdded(this, PLUGIN_KEY);
        else
            m_proxyInter->itemRemoved(this, PLUGIN_KEY);
    });
}

const QString MultitaskingPlugin::pluginName() const
{
    return "multitasking";
}

const QString MultitaskingPlugin::pluginDisplayName() const
{
    return tr("Multitasking View");
}

QWidget *MultitaskingPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_multitaskingWidget.data();
}

QWidget *MultitaskingPlugin::itemTipsWidget(const QString &itemKey)
{
    m_tipsLabel->setObjectName(itemKey);

    m_tipsLabel->setText(pluginDisplayName());

    return m_tipsLabel.data();
}

void MultitaskingPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void MultitaskingPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool MultitaskingPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString MultitaskingPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == PLUGIN_KEY)
        return "dbus-send --session --dest=com.deepin.wm --print-reply /com/deepin/wm com.deepin.wm.PerformAction int32:1";

    return "";
}

const QString MultitaskingPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey != PLUGIN_KEY) {
        return QString();
    }

    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> desktop;
    desktop["itemId"] = "multitasking";
    desktop["itemText"] = tr("Multitasking View");
    desktop["isActive"] = true;
    items.push_back(desktop);

    QMap<QString, QVariant> power;
    power["itemId"] = "remove";
    power["itemText"] = tr("Undock");
    power["isActive"] = true;
    items.push_back(power);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void MultitaskingPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "multitasking") {
        DDBusSender()
        .service("com.deepin.wm")
        .interface("com.deepin.wm")
        .path("/com/deepin/wm")
        .method(QString("PerformAction"))
        .arg(1)
        .call();
    } else if (menuId == "remove") {
        pluginStateSwitched();
    }
}

void MultitaskingPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == PLUGIN_KEY) {
        m_multitaskingWidget->refreshIcon();
    }
}

int MultitaskingPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 2).toInt();
}

void MultitaskingPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void MultitaskingPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

PluginsItemInterface::PluginType MultitaskingPlugin::type()
{
    return PluginType::Fixed;
}

void MultitaskingPlugin::updateVisible()
{
    if (pluginIsDisable() || !DWindowManagerHelper::instance()->hasComposite()) {
        m_proxyInter->itemRemoved(this, PLUGIN_KEY);
    } else {
        m_proxyInter->itemAdded(this, PLUGIN_KEY);
    }
}

void MultitaskingPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        return;
    }

    m_pluginLoaded = true;

    m_multitaskingWidget.reset(new MultitaskingWidget);

    m_proxyInter->itemAdded(this, pluginName());

    updateVisible();
}

void MultitaskingPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, PLUGIN_KEY);
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        updateVisible();
    }
}
