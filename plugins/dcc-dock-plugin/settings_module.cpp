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
#include "config_watcher.h"

#include <QLayout>

#include <DApplication>
#include <DConfig>

DWIDGET_USE_NAMESPACE
DCORE_USE_NAMESPACE

SettingsModule::SettingsModule()
    : QObject()
    , ModuleInterface()
    , m_moduleWidget(nullptr)
    , m_config(new DConfig("org.deepin.dde.dock.plugin", QString(), this))
{
    QTranslator *translator = new QTranslator(this);
    translator->load(QString("/usr/share/dcc-dock-plugin/translations/dcc-dock-plugin_%1.qm").arg(QLocale::system().name()));
    QCoreApplication::installTranslator(translator);
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
    return QStringList() << "Dock";
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
    return QString(":/translations/dcc-dock-plugin_%1.ts");
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
    return QStringLiteral("Dock");
}

void SettingsModule::showPage(const QString &pageName)
{
    Q_UNUSED(pageName);
}

void SettingsModule::addChildPageTrans() const
{
    if (!m_frameProxy)
        return;

    m_frameProxy->addChildPageTrans("Dock", tr("Dock"));
}

void SettingsModule::initSearchData()
{
    onStatusChanged();

    if (m_config->isValid())
        connect(m_config, &DConfig::valueChanged, this, &SettingsModule::onStatusChanged);
}

void SettingsModule::preInitialize(bool sync, FrameProxyInterface::PushType)
{
    Q_UNUSED(sync);
    addChildPageTrans();
    initSearchData();
}

void SettingsModule::onStatusChanged()
{
    if (!m_frameProxy)
        return;

    // 模块名称
    const QString &module = m_frameProxy->moduleDisplayName(PERSONALIZATION);

    // 子模块名称
    const QString &dock = tr("Dock");

    // 二级菜单显示状态设置
    m_frameProxy->setWidgetVisible(module, dock, true);

    auto visibleState = [ = ](const QString &key) {
        return (m_config->value(QString("%1").arg(key)).toString() == "Enabled");
    };

    // 三级菜单显示状态设置
    m_frameProxy->setDetailVisible(module, dock, tr("Mode"), visibleState("dockModel"));
    m_frameProxy->setDetailVisible(module, dock, tr("Location"), visibleState("dockLocation"));
    m_frameProxy->setDetailVisible(module, dock, tr("Status"), visibleState("dockState"));
    m_frameProxy->setDetailVisible(module, dock, tr("Size"), visibleState("dockSize"));
    m_frameProxy->setDetailVisible(module, dock, tr("Show Dock"), visibleState("multiscreen"));
    m_frameProxy->setDetailVisible(module, dock, tr("Plugin Area"), visibleState("dockPlugins"));
    m_frameProxy->updateSearchData(module);
}
