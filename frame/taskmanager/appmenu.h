// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef APPMENU_H
#define APPMENU_H

#include <QString>
#include <QJsonObject>
#include <QVector>

#include <memory>
#include <vector>
#include <functional>

typedef std::function<void(uint32_t)> AppMenuAction;

class AppMenu;

// 应用菜单选项
struct AppMenuItem
{
    AppMenuItem()
        : isActive(true)
        , hint(0)
    {
    }

    QString id;
    QString text;
    QString isCheckable;
    QString checked;
    QString icon;
    QString iconHover;
    QString iconInactive;
    QString showCheckMark;
    std::shared_ptr<AppMenu> subMenu;

    bool isActive;
    int hint;
    AppMenuAction action;
};

// 应用菜单类
class AppMenu
{
public:
    AppMenu();

    void appendItem(AppMenuItem item);
    void setDirtyStatus(bool isDirty);
    void handleAction(uint32_t timestamp, QString itemId);

    QString getMenuJsonStr();

private:
    QString allocateId();

private:
    int m_itemCount;
    bool m_dirty;
    bool m_checkableMenu;               // json:"checkableMenu"
    bool m_singleCheck;                 // json:"singleCheck"
    QVector<AppMenuItem> m_items;       // json:"items"
};

#endif // APPMENU_H
