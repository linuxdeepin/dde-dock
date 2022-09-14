// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "fcitxinputmethoditem.h"

#include <QDBusArgument>
#include <QDBusMetaType>

bool FcitxQtInputMethodItem::enabled() const
{
    return m_enabled;
}

const QString& FcitxQtInputMethodItem::langCode() const
{
    return m_langCode;
}

const QString& FcitxQtInputMethodItem::name() const
{
    return m_name;
}

const QString& FcitxQtInputMethodItem::uniqueName() const
{
    return m_uniqueName;
}

void FcitxQtInputMethodItem::setEnabled(bool enable)
{
    m_enabled = enable;
}

void FcitxQtInputMethodItem::setLangCode(const QString& lang)
{
    m_langCode = lang;
}

void FcitxQtInputMethodItem::setName(const QString& name)
{
    m_name = name;
}

void FcitxQtInputMethodItem::setUniqueName(const QString& name)
{
    m_uniqueName = name;
}

void FcitxQtInputMethodItem::registerMetaType()
{
    qRegisterMetaType<FcitxQtInputMethodItem>("FcitxQtInputMethodItem");
    qDBusRegisterMetaType<FcitxQtInputMethodItem>();
    qRegisterMetaType<FcitxQtInputMethodItemList>("FcitxQtInputMethodItemList");
    qDBusRegisterMetaType<FcitxQtInputMethodItemList>();
}

QDBusArgument& operator<<(QDBusArgument& argument, const FcitxQtInputMethodItem& im)
{
    argument.beginStructure();
    argument << im.name();
    argument << im.uniqueName();
    argument << im.langCode();
    argument << im.enabled();
    argument.endStructure();
    return argument;
}

const QDBusArgument& operator>>(const QDBusArgument& argument, FcitxQtInputMethodItem& im)
{
    QString name;
    QString uniqueName;
    QString langCode;
    bool enabled;
    argument.beginStructure();
    argument >> name >> uniqueName >> langCode >> enabled;
    argument.endStructure();
    im.setName(name);
    im.setUniqueName(uniqueName);
    im.setLangCode(langCode);
    im.setEnabled(enabled);
    return argument;
}