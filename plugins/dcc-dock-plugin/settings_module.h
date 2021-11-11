/*
 * Copyright (C) 2011 ~ 2021 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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
#ifndef SETTINGSMODULE_H
#define SETTINGSMODULE_H

#include <QObject>

#include "interface/namespace.h"
#include "interface/moduleinterface.h"
#include "interface/frameproxyinterface.h"

namespace DCC_NAMESPACE {
  class ModuleInterface;
  class FrameProxyInterface;
}

using namespace DCC_NAMESPACE;

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

public Q_SLOTS:
    void active() Q_DECL_OVERRIDE;

private:
    ModuleWidget *m_moduleWidget;
};

#endif // SETTINGSMODULE_H
