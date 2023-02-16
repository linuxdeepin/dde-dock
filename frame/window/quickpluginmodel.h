// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef QUICKPLUGINMODEL_H
#define QUICKPLUGINMODEL_H

#include <QObject>
#include <QMap>

class PluginsItemInterface;
enum class DockPart;
/**
 *  这是一个独立的Model，用来记录显示在任务栏下方的快捷插件
 * @brief The QuickPluginModel class
 */

class QuickPluginModel : public QObject
{
    Q_OBJECT

public:
    static QuickPluginModel *instance();

    void addPlugin(PluginsItemInterface *itemInter, int index = -1);
    void removePlugin(PluginsItemInterface *itemInter);

    QList<PluginsItemInterface *> dockedPluginItems() const;
    bool isDocked(PluginsItemInterface *itemInter) const;
    bool isFixed(PluginsItemInterface *itemInter) const;

Q_SIGNALS:
    void requestUpdate();
    void requestUpdatePlugin(PluginsItemInterface *, const DockPart &);

protected:
    explicit QuickPluginModel(QObject *parent = nullptr);

private Q_SLOTS:
    void onPluginRemoved(PluginsItemInterface *itemInter);

private:
    void initConnection();
    void initConfig();
    void saveConfig();
    int getCurrentIndex(PluginsItemInterface *itemInter);
    int generaIndex(int sourceIndex, int oldIndex);

private:
    QList<PluginsItemInterface *> m_dockedPluginsItems;
    QMap<QString, int> m_dockedPluginIndex;
};

#endif // QUICKPLUGINMODEL_H
