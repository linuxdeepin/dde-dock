// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MEDIAPLUGIN_H
#define MEDIAPLUGIN_H

#include "pluginsiteminterface.h"

namespace Dock{
class TipsWidget;
}
class MediaWidget;
class MediaPlayerModel;

class MediaPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "media.json")

public:
    explicit MediaPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;

    PluginFlags flags() const override;

private:
    QScopedPointer<MediaWidget> m_mediaWidget;
    QScopedPointer<MediaPlayerModel> m_model;
};

#endif // DATETIMEPLUGIN_H
