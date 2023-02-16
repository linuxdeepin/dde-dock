// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DISPLAYPLUGIN_H
#define DISPLAYPLUGIN_H

#include "pluginsiteminterface.h"

#include <QTimer>
#include <QLabel>
#include <QSettings>

namespace Dock{
class TipsWidget;
}

class BrightnessWidget;
class BrightnessModel;
class DisplaySettingWidget;

class DisplayPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "display.json")

public:
    explicit DisplayPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;

    PluginFlags flags() const override;

private:
    QScopedPointer<BrightnessWidget> m_displayWidget;
    QScopedPointer<DisplaySettingWidget> m_displaySettingWidget;
    QScopedPointer<Dock::TipsWidget> m_displayTips;
    QScopedPointer<BrightnessModel> m_model;
};

#endif // DATETIMEPLUGIN_H
