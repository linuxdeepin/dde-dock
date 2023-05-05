// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later
#include "soundplugin.h"
#include "soundaccessible.h"
#include "soundwidget.h"
#include "sounddeviceswidget.h"

#include <DDBusSender>

#include <QDebug>
#include <QAccessible>

#define STATE_KEY  "enable"
#define SOUND_KEY "sound-item-key"

SoundPlugin::SoundPlugin(QObject *parent)
    : QObject(parent)
    , m_soundWidget(nullptr)
    , m_soundDeviceWidget(nullptr)
{
    QAccessible::installFactory(soundAccessibleFactory);
}

const QString SoundPlugin::pluginName() const
{
    return "sound";
}

const QString SoundPlugin::pluginDisplayName() const
{
    return tr("Sound");
}

void SoundPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (m_soundWidget) return;

    m_soundWidget.reset(new SoundWidget);
    m_soundWidget->setFixedHeight(60);

    m_soundDeviceWidget.reset(new SoundDevicesWidget);

    if (!pluginIsDisable()) {
        m_proxyInter->itemAdded(this, SOUND_KEY);
        connect(m_soundWidget.data(), &SoundWidget::rightIconClick, this, [ this, proxyInter ] {
            proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, true);
        });
    }

    connect(m_soundDeviceWidget.data(), &SoundDevicesWidget::enableChanged, m_soundWidget.data(), &SoundWidget::setEnabled);
    connect(m_soundDeviceWidget.data(), &SoundDevicesWidget::requestHide, this, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, false);
    });

    connect(m_soundDeviceWidget.data(), &SoundDevicesWidget::iconChanged, this, [=] {
        m_proxyInter->updateDockInfo(this, DockPart::QuickPanel);
        m_proxyInter->updateDockInfo(this, DockPart::QuickShow);
        m_proxyInter->itemUpdate(this, SOUND_KEY);
    });
}

void SoundPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool SoundPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, STATE_KEY, true).toBool();
}

QWidget *SoundPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY)
        return m_soundWidget.data();

    return nullptr;
}

QWidget *SoundPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == SOUND_KEY)
        return m_soundDeviceWidget->tipsWidget();

    return nullptr;
}

QWidget *SoundPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY)
        return m_soundDeviceWidget.data();

    return nullptr;
}

int SoundPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 2).toInt();
}

void SoundPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void SoundPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

QIcon SoundPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    switch (dockPart) {
    case DockPart::QuickShow:
        return m_soundDeviceWidget->pixmap(themeType, 18, 16);
    case DockPart::DCCSetting:
        return m_soundDeviceWidget->pixmap(themeType, 18, 18);
    default:
        break;
    }
    return QIcon();
}

PluginsItemInterface::PluginMode SoundPlugin::status() const
{
    return SoundPlugin::Active;
}

PluginFlags SoundPlugin::flags() const
{
    return PluginFlag::Type_Common
            | PluginFlag::Quick_Full
            | PluginFlag::Attribute_CanDrag
            | PluginFlag::Attribute_CanInsert
            | PluginFlag::Attribute_CanSetting;
}

bool SoundPlugin::eventHandler(QEvent *event)
{
    // 当前只处理鼠标滚轮事件
    if (event->type() != QEvent::Wheel)
        return PluginsItemInterface::eventHandler(event);

    // 获取当前默认的声音设备
    QDBusPendingCall defaultSinkCall = DDBusSender().service("org.deepin.dde.Audio1")
            .path("/org/deepin/dde/Audio1")
            .interface("org.deepin.dde.Audio1")
            .property("DefaultSink").get();
    defaultSinkCall.waitForFinished();
    QDBusReply<QVariant> path = defaultSinkCall.reply();
    const QString defaultSinkPath = path.value().value<QDBusObjectPath>().path();
    if (defaultSinkPath.isNull())
        return false;

    // 获取当前默认声音设备的音量
    DDBusSender sinkDBus = DDBusSender().service("org.deepin.dde.Audio1")
            .path(defaultSinkPath).interface("org.deepin.dde.Audio1.Sink");
    QDBusPendingCall volumeCall = sinkDBus.property("Volume").get();
    volumeCall.waitForFinished();
    QDBusReply<QVariant> volumePath = volumeCall.reply();
    double volume = volumePath.value().value<double>();

    // 获取当前默认声音设备的最大音量
    DDBusSender audioDBus = DDBusSender().service("org.deepin.dde.Audio1")
            .path("/org/deepin/dde/Audio1").interface("org.deepin.dde.Audio1");
    QDBusPendingCall call = audioDBus.property("MaxUIVolume").get();
    call.waitForFinished();
    QDBusReply<QVariant> maxVolumeReply = call.reply();
    double maxVolume = maxVolumeReply.value().value<double>();

    // 根据滚轮的动作来增加音量或者减小音量
    QWheelEvent *wheelEvent = static_cast<QWheelEvent *>(event);
    if (wheelEvent->angleDelta().y() > 0) {
        // 向上滚动，增大音量
        if (volume < maxVolume)
            sinkDBus.method("SetVolume").arg(qMin(volume + 0.02, maxVolume)).arg(true).call();
    } else {
        // 向下滚动，调小音量
        if (volume > 0)
            sinkDBus.method("SetVolume").arg(qMax(volume - 0.02, 0.0)).arg(true).call();
    }

    return true;
}

void SoundPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, SOUND_KEY);
    else
        m_proxyInter->itemAdded(this, SOUND_KEY);
}
