// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef FCITX_QT_INPUT_METHOD_ITEM_H
#define FCITX_QT_INPUT_METHOD_ITEM_H

// Qt
#include <QtCore/QString>
#include <QtCore/QMetaType>
#include <QtDBus/QDBusArgument>

class FcitxQtInputMethodItem
{
public:
    const QString& name() const;
    const QString& uniqueName() const;
    const QString& langCode() const;
    bool enabled() const;

    void setName(const QString& name);
    void setUniqueName(const QString& name);
    void setLangCode(const QString& name);
    void setEnabled(bool name);
    static void registerMetaType();

    inline bool operator < (const FcitxQtInputMethodItem& im) const {
        return (m_enabled && !im.m_enabled);
    }

private:
    QString m_name;
    QString m_uniqueName;
    QString m_langCode;
    bool m_enabled;
};

typedef QList<FcitxQtInputMethodItem> FcitxQtInputMethodItemList;

QDBusArgument& operator<<(QDBusArgument& argument, const FcitxQtInputMethodItem& im);
const QDBusArgument& operator>>(const QDBusArgument& argument, FcitxQtInputMethodItem& im);

Q_DECLARE_METATYPE(FcitxQtInputMethodItem)
Q_DECLARE_METATYPE(FcitxQtInputMethodItemList)

#endif
