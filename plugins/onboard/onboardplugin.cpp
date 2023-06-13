// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "onboardplugin.h"
#include "../widgets/tipswidget.h"

#include <DGuiApplicationHelper>

#include <QIcon>
#include <QSettings>
#include <QPainter>

#define PLUGIN_STATE_KEY    "enable"

DGUI_USE_NAMESPACE

using namespace Dock;
OnboardPlugin::OnboardPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_startupState(false)
    , m_onboardItem(nullptr)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setText(tr("Onboard"));
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setAccessibleName("Onboard");
}

const QString OnboardPlugin::pluginName() const
{
    return "onboard";
}

const QString OnboardPlugin::pluginDisplayName() const
{
    return tr("Onboard");
}

QWidget *OnboardPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == pluginName())
        return m_onboardItem.data();

    return nullptr;
}

QWidget *OnboardPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_tipsLabel.data();
}

void OnboardPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void OnboardPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool OnboardPlugin::pluginIsDisable()
{
    return !(m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());
}

const QString OnboardPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=org.onboard.Onboard /org/onboard/Onboard/Keyboard org.onboard.Onboard.Keyboard.ToggleVisible");
}

void OnboardPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId != "onboard-settings")
        return;

    if(!m_startupState) {
        QProcess *process = new QProcess;
        connect(process,&QProcess::started, this, [ = ] {
            m_startupState = true;
        });

        connect(process, QOverload<int, QProcess::ExitStatus>::of(&QProcess::finished),
              [ = ](int exitCode, QProcess::ExitStatus exitStatus){
            Q_UNUSED(exitCode)
            Q_UNUSED(exitStatus)

            m_startupState = false;
            process->close();
            process->deleteLater();
        });
        process->start("onboard-settings", QStringList());
    }

}

void OnboardPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (!pluginIsDisable()) {
        m_onboardItem->update();
    }
}

int OnboardPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 3).toInt();
}

void OnboardPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

void OnboardPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

QIcon OnboardPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    if (dockPart == DockPart::DCCSetting) {
        if (themeType == DGuiApplicationHelper::ColorType::LightType)
            return QIcon(":/icons/icon/dcc_keyboard.svg");

        QPixmap pixmap(":/icons/icon/dcc_keyboard.svg");
        QPainter pa(&pixmap);
        pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
        pa.fillRect(pixmap.rect(), Qt::white);
        return pixmap;
    }

    if (dockPart == DockPart::QuickPanel)
        return m_onboardItem->iconPixmap(QSize(24, 24), themeType);

    return m_onboardItem->iconPixmap(QSize(18, 16), themeType);
}

PluginsItemInterface::PluginMode OnboardPlugin::status() const
{
    return PluginsItemInterface::PluginMode::Active;
}

QString OnboardPlugin::description() const
{
    return pluginDisplayName();
}

void OnboardPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "onboard plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_onboardItem.reset(new OnboardItem);

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());
}

void OnboardPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
    {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    }
}
