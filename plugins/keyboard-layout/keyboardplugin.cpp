// SPDX-FileCopyrightText: 2017 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "keyboardplugin.h"
#include "utils.h"

KeyboardPlugin::KeyboardPlugin(QObject *parent)
    : QObject(parent)
{
}

KeyboardPlugin::~KeyboardPlugin()
{
}

const QString KeyboardPlugin::pluginName() const
{
    return "keyboard";
}

const QString KeyboardPlugin::pluginDisplayName() const
{
    return "Keyboard";
}

void KeyboardPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;
    if (!m_dbusAdaptors) {

        QString serverName = "com.deepin.daemon.InputDevices";
        QDBusConnectionInterface *ifc = QDBusConnection::sessionBus().interface();

        if (!ifc->isServiceRegistered(serverName)) {
            connect(QDBusConnection::sessionBus().interface(), &QDBusConnectionInterface::serviceOwnerChanged, this,
            [ = ](const QString & name, const QString & oldOwner, const QString & newOwner) {
                Q_UNUSED(oldOwner);
                if (name == serverName && !newOwner.isEmpty()) {
                    m_dbusAdaptors = new DBusAdaptors(this);
                    disconnect(ifc);
                }
            });
        } else {
            m_dbusAdaptors = new DBusAdaptors(this);
        }

        QDBusConnection::sessionBus().registerService("com.deepin.dde.Keyboard");
        QDBusConnection::sessionBus().registerObject("/com/deepin/dde/Keyboard", "com.deepin.dde.Keyboard", this);
    }
}

QWidget *KeyboardPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

QWidget *KeyboardPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

int KeyboardPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 1).toInt();
}

void KeyboardPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}
