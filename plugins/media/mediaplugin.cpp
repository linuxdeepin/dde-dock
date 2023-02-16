// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mediaplugin.h"
#include "mediawidget.h"
#include "mediaplayermodel.h"

MediaPlugin::MediaPlugin(QObject *parent)
    : QObject(parent)
    , m_mediaWidget(nullptr)
    , m_model(nullptr)
{
}

const QString MediaPlugin::pluginName() const
{
    return "media";
}

const QString MediaPlugin::pluginDisplayName() const
{
    return "Media";
}

void MediaPlugin::init(PluginProxyInterface *proxyInter)
{
    if (m_proxyInter == proxyInter)
        return;

    m_proxyInter = proxyInter;

    m_model.reset(new MediaPlayerModel);
    m_mediaWidget.reset(new MediaWidget(m_model.data()));
    m_mediaWidget->setFixedHeight(60);
    m_mediaWidget->setVisible(m_model->isActived());

    if (m_model->isActived())
        m_proxyInter->itemAdded(this, pluginName());

    connect(m_model.data(), &MediaPlayerModel::startStop, this, [ this ](bool visible) {
        if (visible)
            m_proxyInter->itemAdded(this, pluginName());
        else
            m_proxyInter->itemRemoved(this, pluginName());
    });
}

QWidget *MediaPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY)
        return m_mediaWidget.data();

    return nullptr;
}

QWidget *MediaPlugin::itemTipsWidget(const QString &itemKey)
{
    return nullptr;
}

QWidget *MediaPlugin::itemPopupApplet(const QString &itemKey)
{
    return nullptr;
}

PluginFlags MediaPlugin::flags() const
{
    return PluginFlag::Type_Common | PluginFlag::Quick_Full;
}
