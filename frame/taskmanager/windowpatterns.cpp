// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "windowpatterns.h"
#include "processinfo.h"

#include <QJsonDocument>
#include <QStandardPaths>
#include <QJsonObject>
#include <QJsonValue>
#include <QJsonArray>
#include <QVariant>
#include <QVariantMap>
#include <QDebug>
#include <QRegExp>
#include <QFile>
#include <QFileInfo>

const int parsedFlagNegative = 0x001;
const int parsedFlagIgnoreCase = 0x010;

const QString getWindowPatternsFile(){
    for (auto dataLocation : QStandardPaths::standardLocations(QStandardPaths::GenericDataLocation)) {
        QString targetFilePath = dataLocation.append("/dde-dock/window_patterns.json");
        if (QFile::exists(targetFilePath)) return targetFilePath;
    }

    return QString();
} 

bool contains(QString key, QString value) {
    return key.contains(value);
}

bool containsIgnoreCase(QString key, QString value) {
    QString _key = key.toLower();
    QString _value = value.toLower();
    return _key.contains(_value);
}

bool equal(QString key, QString value) {
    return key == value;
}

bool equalIgnoreCase(QString key, QString value) {
    return key.toLower() == value.toLower();
}

bool regexMatch(QString key, QString value) {
    QRegExp ruleRegex(value);
    bool ret = ruleRegex.exactMatch(key);

    // 配置中\.exe$ 在V20中go代码可以匹配以.exe结尾的字符串, Qt中使用\.*exe$匹配以.exe结尾字符串失败，暂时做兼容处理
    if (!ret && value == "\\.exe$") {
        ret = key.endsWith(".exe");
    }
    return ret;
}

bool regexMatchIgnoreCase(QString key, QString value) {
    QRegExp ruleRegex(value, Qt::CaseInsensitive);
    bool ret = ruleRegex.exactMatch(key);

    // 配置中\.exe$ 在V20中go代码可以匹配以.exe结尾的字符串, Qt中使用\.*exe$匹配以.exe结尾字符串失败，暂时做兼容处理
    if (!ret && value == "\\.exe$") {
        QString _key = key.toLower();
        ret = _key.endsWith(".exe");
    }
    return ret;
}

RuleValueParse::RuleValueParse()
 : negative(false)
 , type(0)
 , flags(0)
{
}

bool RuleValueParse::match(const WindowInfoX *winInfo)
{
    QString parsedKey = parseRuleKey(const_cast<WindowInfoX *>(winInfo), key);
    if (!fn)
        return false;

    bool ret = fn(parsedKey, value);
    return negative ? !ret : ret;
}

QString RuleValueParse::parseRuleKey(WindowInfoX *winInfo, const QString &ruleKey)
{
    ProcessInfo * process = winInfo->getProcess();
    if (ruleKey == "hasPid") {
        if (process && process->initWithPid()) {
            return "t";
        }
        return "f";
    } else if (ruleKey == "exec") {
        if (process) {
            // 返回执行文件baseName
            auto baseName = QFileInfo(process->getExe()).completeBaseName();
            return baseName.isEmpty() ? "" : baseName;
        }
    } else if (ruleKey == "arg") {
        if (process) {
            // 返回命令行参数
            auto ret = process->getArgs().join("");
            return ret.isEmpty() ? "" : ret;
        }
    } else if (ruleKey == "wmi") {
        // 窗口实例
        auto wmClass = winInfo->getWMClass();
        if (!wmClass.instanceName.empty())
            return wmClass.instanceName.c_str();
    } else if (ruleKey == "wmc") {
        // 窗口类型
        auto wmClass = winInfo->getWMClass();
        if (!wmClass.className.empty())
            return wmClass.className.c_str();
    } else if (ruleKey == "wmn") {
        // 窗口名称
        return winInfo->getWMName();
    } else if (ruleKey == "wmrole") {
        // 窗口角色
        return winInfo->getWmRole();
    } else {
        const QString envPrefix = "env.";
        if (ruleKey.startsWith(envPrefix)) {
            QString envName = ruleKey.mid(envPrefix.size());
            if (winInfo->getProcess()) {
                auto ret = process->getEnv(envName);
                return ret.isEmpty() ? "" : ret;
            }
        }
    }

    return "";
}


WindowPatterns::WindowPatterns()
{
    loadWindowPatterns();
}

/**
 * @brief WindowPatterns::match 匹配窗口类型
 * @param winInfo
 * @return
 */
QString WindowPatterns::match(WindowInfoX *winInfo)
{
    for (auto pattern : m_patterns) {
        bool patternOk = true;
        for (auto rule : pattern.parseRules) {
            if (!rule.match(winInfo)) {
                patternOk = false;
                break;
            }
        }

        if (patternOk) {
            // 匹配成功
            return pattern.result;
        }
    }

    // 匹配失败
    return "";
}

void WindowPatterns::loadWindowPatterns()
{
    qInfo() << "---loadWindowPatterns";
    QFile file(getWindowPatternsFile());
    if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
        return;

     QJsonDocument doc = QJsonDocument::fromJson(file.readAll());
     file.close();
     if (!doc.isArray())
         return;

     QJsonArray arr = doc.array();
     if (arr.size() == 0)
         return;

     m_patterns.clear();
     for (auto iterp = arr.begin(); iterp != arr.end(); iterp++) {
         // 过滤非Object
        if (!(*iterp).isObject())
            continue;

        QJsonObject patternObj = (*iterp).toObject();
        QVariantMap patternMap = patternObj.toVariantMap();
        WindowPattern pattern;
        for (auto infoIter = patternMap.begin(); infoIter != patternMap.end(); infoIter++) {
            QString ret = infoIter.key();
            QVariant value = infoIter.value();

            if (ret == "ret") {
                pattern.result = value.toString();
            } else if (ret == "rules") {
                for (auto &item : value.toList()) {
                    if (!item.isValid())
                        continue;

                    if (item.toList().size() != 2)
                        continue;

                    pattern.rules.push_back({item.toList()[0].toString(), item.toList()[1].toString()});
                }
            }
        }
        qInfo() << pattern.result;
        for (const auto &item : pattern.rules) {
            qInfo() << item[0] << " " << item[1];
        }
        m_patterns.push_back(pattern);
     }

     // 解析patterns
     for (auto &pattern : m_patterns) {
        for (int i=0; i < pattern.rules.size(); i++) {
            RuleValueParse ruleValue = parseRule(pattern.rules[i]);
            pattern.parseRules.push_back(ruleValue);
        }
     }
}

// "=:XXX" equal XXX
// "=!XXX" not equal XXX

// "c:XXX" contains XXX
// "c!XXX" not contains XXX

// "r:XXX" match regexp XXX
// "r!XXX" not match regexp XXX

// e c r ignore case
// = E C R not ignore case
// 解析窗口类型规则
RuleValueParse WindowPatterns::parseRule(QVector<QString> rule)
{
    RuleValueParse ret;
    ret.key = rule[0];
    ret.original = rule[1];
    if (rule[1].size() < 2)
        return ret;

    int len = ret.original.size() + 1;
    char *orig = static_cast<char *>(calloc(1, size_t(len)));
    if (!orig)
        return ret;

    strncpy(orig, ret.original.toStdString().c_str(), size_t(len));
    switch (orig[1]) {
    case ':':
        break;
    case '!':
        ret.flags |= parsedFlagNegative;
        ret.negative = true;
        break;
    default:
        return ret;
    }

    ret.value = QString(&orig[2]);
    ret.type = uint8_t(orig[0]);
    switch (orig[0]) {
    case 'C':
        ret.fn = contains;
        break;
    case 'c':
        ret.flags |= parsedFlagIgnoreCase;
        ret.fn = containsIgnoreCase;
        break;
    case '=':
    case 'E':
        ret.fn = equal;
        break;
    case 'e':
        ret.flags |= parsedFlagIgnoreCase;
        ret.fn = equalIgnoreCase;
        break;
    case 'R':
        ret.fn = regexMatch;
        break;
    case 'r':
        ret.flags |= parsedFlagIgnoreCase;
        ret.fn = regexMatchIgnoreCase;
        break;
    default:
        break;
    }

    free(orig);
    return ret;
}
