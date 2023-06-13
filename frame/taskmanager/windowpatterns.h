// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWPATTERNS_H
#define WINDOWPATTERNS_H

#include "windowinfox.h"

#include <QString>
#include <QVector>


struct RuleValueParse {
    RuleValueParse();
    bool match(const WindowInfoX *winInfo);
    static QString parseRuleKey(WindowInfoX *winInfo, const QString &ruleKey);
    QString key;
    bool negative;
    bool (*fn)(QString, QString);
    uint8_t type;
    uint flags;
    QString original;
    QString value;
};

class WindowPatterns
{
    // 窗口类型匹配
    struct WindowPattern {
        QVector<QVector<QString>> rules;    // rules
        QString result;                     // ret
        QVector<RuleValueParse> parseRules;
    };

public:
    WindowPatterns();

    QString match(WindowInfoX *winInfo);

private:
    void loadWindowPatterns();
    RuleValueParse parseRule(QVector<QString> rule);

private:
    QVector<WindowPattern> m_patterns;

};

#endif // WINDOWPATTERNS_H
