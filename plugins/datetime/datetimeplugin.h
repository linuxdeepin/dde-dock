// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DATETIMEPLUGIN_H
#define DATETIMEPLUGIN_H

#include "pluginsiteminterface.h"
#include "datetimewidget.h"

#include <QTimer>
#include <QLabel>
#include <QSettings>

namespace Dock{
class TipsWidget;
}
class QDBusInterface;
class DatetimePlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "datetime.json")

public:
    explicit DatetimePlugin(QObject *parent = nullptr);

    PluginSizePolicy pluginSizePolicy() const override;

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;

    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;

    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;

    void pluginSettingsChanged() override;

private slots:
    void updateCurrentTimeString();
    void refreshPluginItemsVisible();
    void propertiesChanged();

private:
    void loadPlugin();
    QDBusInterface *timedateInterface();

private:
    QScopedPointer<DatetimeWidget> m_centralWidget;
    QScopedPointer<Dock::TipsWidget> m_dateTipsLabel;
    QTimer *m_refershTimer;
    QString m_currentTimeString;
    QDBusInterface *m_interface;
    bool m_pluginLoaded;
};

#endif // DATETIMEPLUGIN_H
