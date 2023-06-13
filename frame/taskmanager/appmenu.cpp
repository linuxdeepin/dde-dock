// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "appmenu.h"

#include <QJsonArray>
#include <QJsonDocument>

AppMenu::AppMenu()
 : m_itemCount(0)
 , m_dirty(false)
 , m_checkableMenu(false)
 , m_singleCheck(false)
{

}

/**
 * @brief AppMenu::appendItem 增加菜单选项
 * @param item
 */
void AppMenu::appendItem(AppMenuItem item)
{
    if (!item.text.isEmpty()) {
        item.id = allocateId();
        m_items.push_back(item);
    }
}

/**
 * @brief AppMenu::handleAction 响应应用菜单项
 * @param timestamp
 * @param itemId
 */
void AppMenu::handleAction(uint32_t timestamp, QString itemId)
{
    for (auto &item : m_items) {
        if (item.id == itemId) {
            item.action(timestamp);
            break;
        }
    }
}

void AppMenu::setDirtyStatus(bool isDirty)
{
    m_dirty = isDirty;
}

QString AppMenu::getMenuJsonStr()
{
    QJsonObject obj;
    QJsonArray array;
    for (auto item : m_items) {
        QJsonObject objItem;
        objItem["itemId"] = item.id;
        objItem["itemText"] = item.text;
        objItem["isActive"] = item.isActive;
        objItem["isCheckable"] = item.isCheckable;
        objItem["checked"] = item.checked;
        objItem["itemIcon"] = item.icon;
        objItem["itemIconHover"] = item.iconHover;
        objItem["itemIconInactive"] = item.iconInactive;
        objItem["showCheckMark"] = item.showCheckMark;
        objItem["itemSubMenu"] = item.subMenu ? item.subMenu->getMenuJsonStr() : "";
        array.push_back(objItem);
    }
    obj["items"] = QJsonValue(array);
    obj["checkableMenu"] = m_checkableMenu;
    obj["singleCheck"] = m_singleCheck;

    QString ret = QJsonDocument(obj).toJson();
    return ret;
}

QString AppMenu::allocateId()
{
    return QString::number(m_itemCount++);
}
