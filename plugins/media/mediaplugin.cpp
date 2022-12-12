/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
