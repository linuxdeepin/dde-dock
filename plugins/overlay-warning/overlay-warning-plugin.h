// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef OVERLAY_WARNING_PLUGIN_H
#define OVERLAY_WARNING_PLUGIN_H

#include "pluginsiteminterface.h"
#include "overlaywarningwidget.h"
#include "../widgets/tipswidget.h"

#include <QLabel>

namespace Dtk {
    namespace Widget {
        class DDialog;
    }
}

class OverlayWarningPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "overlay-warning.json")

public:
    explicit OverlayWarningPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return false; }
    bool pluginIsDisable() override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    PluginFlags flags() const override;

private:
    void loadPlugin();
    bool isOverlayRoot();

private slots:
    void showCloseOverlayDialogPre();
    void showCloseOverlayDialog();

private:
    bool m_pluginLoaded;

    QScopedPointer<OverlayWarningWidget> m_warningWidget;
    QTimer *m_showDisableOverlayDialogTimer;
};

#endif // OVERLAY_WARNING_PLUGIN_H
