// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINPROXYINTERFACE_H
#define PLUGINPROXYINTERFACE_H

#include "constants.h"

#include <QtCore>

class PluginsItemInterface;
enum class DockPart;

class PluginProxyInterface
{
public:
    ///
    /// \brief itemAdded
    /// add a new dock item
    /// if itemkey of this plugin inter already exist, the new item
    /// will be ignored, so if you need to add multiple item, you need
    /// to ensure all itemKey is different.
    /// \param itemInter
    /// your plugin interface
    /// \param itemKey
    /// your item unique key
    ///
    virtual void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;
    ///
    /// \brief itemUpdate
    /// update(repaint) spec item
    /// \param itemInter
    /// \param itemKey
    ///
    virtual void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;
    ///
    /// \brief itemRemoved
    /// remove spec item, if spec item is not exist, dock will to nothing.
    /// dock will NOT delete your object, you should manage memory by your self.
    /// \param itemInter
    /// \param itemKey
    ///
    virtual void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;

    ///
    /// \brief requestContextMenu
    /// request show context menu
    ///
    //virtual void requestContextMenu(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;

    virtual void requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide) = 0;
    virtual void requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;

    virtual void requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible) = 0;

    ///
    /// \brief saveValue
    /// save module config to .config/deepin/dde-dock.conf
    /// all key-values of all plugins will be save to that file
    /// and grouped by the returned value of pluginName() function which is defined in PluginsItemInterface
    /// \param itemInter the plugin object
    /// \param key the key of data
    /// \param value the data
    ///
    virtual void saveValue(PluginsItemInterface * const itemInter, const QString &key, const QVariant &value) = 0;

    ///
    /// \brief getValue
    /// SeeAlse: saveValue
    /// return value from .config/deepin/dde-dock.conf
    ///
    virtual const QVariant getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant()) = 0;

    ///
    /// \brief removeValue
    /// remove the values specified by keyList
    /// remove all values of itemInter if keyList is empty
    /// SeeAlse: saveValue
    ///
    virtual void removeValue(PluginsItemInterface *const itemInter, const QStringList &keyList) = 0;

    ///
    /// update display or information
    ///
    ///
    virtual void updateDockInfo(PluginsItemInterface *const, const DockPart &) {}
};

#endif // PLUGINPROXYINTERFACE_H
