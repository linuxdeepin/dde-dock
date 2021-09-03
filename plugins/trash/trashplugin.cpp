/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *               2016 ~ 2018 dragondjf
 *
 * Author:     sbw <sbw@sbw.so>
 *             dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
 *             Tangtong<tangtong@deepin.com>
 *
 * Maintainer: dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
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

#include "trashplugin.h"
#include "../../widgets/tipswidget.h"

#include <DApplication>
#include <DDesktopServices>

#define PLUGIN_STATE_KEY    "enable"

DWIDGET_USE_NAMESPACE

using namespace Dock;

TrashPlugin::TrashPlugin(QObject *parent)
    : QObject(parent)
    , m_trashWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setObjectName("trash");
}

const QString TrashPlugin::pluginName() const
{
    return "trash";
}

const QString TrashPlugin::pluginDisplayName() const
{
    return tr("Trash");
}

void TrashPlugin::init(PluginProxyInterface *proxyInter)
{
    // transfex config
    QSettings settings("deepin", "dde-dock-trash");
    if (QFile::exists(settings.fileName())) {
        const QString &key = QString("pos_%1_%2").arg(pluginName()).arg(displayMode());
        proxyInter->saveValue(this, key, settings.value(key));

        QFile::remove(settings.fileName());
    }

    // blumia: we are using i10n translation from DFM so...
    QString applicationName = qApp->applicationName();
    qApp->setApplicationName("dde-file-manager");
    qDebug() << qApp->loadTranslator();
    qApp->setApplicationName(applicationName);

    m_proxyInter = proxyInter;

    if (m_trashWidget.isNull())
        m_trashWidget.reset(new TrashWidget);

    displayModeChanged(displayMode());
}

QWidget *TrashPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_trashWidget.data();
}

QWidget *TrashPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    if (m_trashWidget->getDragging()) {
        m_tipsLabel->setText(tr("Move to trash"));
        return m_tipsLabel.data();
    }

    const int count = m_trashWidget->trashItemCount();
    if (count < 2)
        m_tipsLabel->setText(tr("Trash - %1 file").arg(count));
    else
        m_tipsLabel->setText(tr("Trash - %1 files").arg(count));

    return m_tipsLabel.data();
}

QWidget *TrashPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

const QString TrashPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);
    DDesktopServices::showFolder(QUrl("trash:///"));

    return QString();
    // return "gio open trash:///";
}

const QString TrashPlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_trashWidget->contextMenu();
}

void TrashPlugin::refreshIcon(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    m_trashWidget->updateIcon();
}

void TrashPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey);

    m_trashWidget->invokeMenuItem(menuId, checked);
}

bool TrashPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

void TrashPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, pluginName());
        return;
    }

    if (m_trashWidget) {
        m_proxyInter->itemAdded(this, pluginName());
    }
}

int TrashPlugin::itemSortKey(const QString &itemKey)
{
    const QString &key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 7).toInt();
}

void TrashPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString &key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

void TrashPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (pluginIsDisable()) {
        return;
    }

    m_proxyInter->itemAdded(this, pluginName());
}

void TrashPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void TrashPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, pluginName());
        return;
    }

    if (m_trashWidget) {
        m_proxyInter->itemAdded(this, pluginName());
    }
}

