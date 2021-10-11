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
#include "settings_module.h"
#include "module_widget.h"

#include <QLayout>

SettingsModule::SettingsModule()
    : QObject()
    , ModuleInterface()
    , m_moduleWidget(nullptr)
{

}

SettingsModule::~SettingsModule()
{

}

void SettingsModule::initialize()
{

}

void SettingsModule::active()
{
    m_moduleWidget = new ModuleWidget;

    m_frameProxy->pushWidget(this, m_moduleWidget);
    m_moduleWidget->setVisible(true);
}

QStringList SettingsModule::availPage() const
{
    return QStringList() << tr("Dock");
}

const QString SettingsModule::displayName() const
{
    return tr("Dock");
}

QIcon SettingsModule::icon() const
{
    return QIcon::fromTheme("icon_dock");
}

QString SettingsModule::translationPath() const
{
    return QString("/usr/share/dde-dock/translations");
}

QString SettingsModule::path() const
{
    return PERSONALIZATION;
}

QString SettingsModule::follow() const
{
    return "10";
}

const QString SettingsModule::name() const
{
    return tr("Dock");
}

void SettingsModule::showPage(const QString &pageName)
{
    Q_UNUSED(pageName);
}
