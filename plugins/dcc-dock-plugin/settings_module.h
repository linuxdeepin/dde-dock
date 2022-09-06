// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SETTINGSMODULE_H
#define SETTINGSMODULE_H

#include <dtkcore_global.h>

#include <QObject>

#include "interface/namespace.h"
#include "interface/moduleinterface.h"
#include "interface/frameproxyinterface.h"

namespace DCC_NAMESPACE {
  class ModuleInterface;
  class FrameProxyInterface;
}

using namespace DCC_NAMESPACE;

DCORE_BEGIN_NAMESPACE
class DConfig;
DCORE_END_NAMESPACE

class ModuleWidget;
class SettingsModule : public QObject, public ModuleInterface
{
    Q_OBJECT

    Q_PLUGIN_METADATA(IID ModuleInterface_iid FILE "dock_settings.json")
    Q_INTERFACES(DCC_NAMESPACE::ModuleInterface)

public:
    explicit SettingsModule();

    ~SettingsModule() Q_DECL_OVERRIDE;

    void initialize() Q_DECL_OVERRIDE;

    QStringList availPage() const Q_DECL_OVERRIDE;

    const QString displayName() const Q_DECL_OVERRIDE;

    QIcon icon() const Q_DECL_OVERRIDE;

    QString translationPath() const Q_DECL_OVERRIDE;

    QString path() const Q_DECL_OVERRIDE;

    QString follow() const Q_DECL_OVERRIDE;

    const QString name() const Q_DECL_OVERRIDE;

    void showPage(const QString &pageName) Q_DECL_OVERRIDE;

    void addChildPageTrans() const Q_DECL_OVERRIDE;

    void initSearchData() Q_DECL_OVERRIDE;

    void preInitialize(bool sync = false,FrameProxyInterface::PushType = FrameProxyInterface::PushType::Normal) Q_DECL_OVERRIDE;

private:
    void onStatusChanged();

public Q_SLOTS:
    void active() Q_DECL_OVERRIDE;

private:
    ModuleWidget *m_moduleWidget;
    DTK_CORE_NAMESPACE::DConfig *m_config;
};

#endif // SETTINGSMODULE_H
