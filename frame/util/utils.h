// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef UTILS
#define UTILS
#include <QPixmap>
#include <QImageReader>
#include <QApplication>
#include <QScreen>
#include <QGSettings>
#include <QDebug>

#include "imageutil.h"

namespace Utils {

#define ICBC_CONF_FILE "/etc/deepin/icbc.conf"

const bool IS_WAYLAND_DISPLAY = !qgetenv("WAYLAND_DISPLAY").isEmpty();

inline bool isDraging()
{
    if (!qApp->property("isDraging").isValid())
        return false;

    return qApp->property("isDraging").toBool();
}

inline void setIsDraging(bool isDraging)
{
    qApp->setProperty("isDraging", isDraging);
}

/**
 * @brief SettingsPtr 根据给定信息返回一个QGSettings指针
 * @param schema_id The id of the schema
 * @param path If non-empty, specifies the path for a relocatable schema
 * @param parent 创建指针的付对象
 * @return
 */
inline QGSettings *SettingsPtr(const QString &schema_id, const QByteArray &path = QByteArray(), QObject *parent = nullptr)
{
    if (QGSettings::isSchemaInstalled(schema_id.toUtf8())) {
        QGSettings *settings = new QGSettings(schema_id.toUtf8(), path, parent);
        return settings;
    }
    qDebug() << "Cannot find gsettings, schema_id:" << schema_id;
    return nullptr;
}

/**
 * @brief SettingsPtr 根据给定信息返回一个QGSettings指针
 * @param module 传入QGSettings构造函数时，会添加"com.deepin.dde.dock.module."前缀
 * @param path If non-empty, specifies the path for a relocatable schema
 * @param parent 创建指针的付对象
 * @return
 */
inline const QGSettings *ModuleSettingsPtr(const QString &module, const QByteArray &path = QByteArray(), QObject *parent = nullptr)
{
    return SettingsPtr("com.deepin.dde.dock.module." + module, path, parent);
}

/* convert 'some-key' to 'someKey' or 'SomeKey'.
 * the second form is needed for appending to 'set' for 'setSomeKey'
 */
inline QString qtify_name(const char *name)
{
    bool next_cap = false;
    QString result;

    while (*name) {
        if (*name == '-') {
            next_cap = true;
        } else if (next_cap) {
            result.append(QChar(*name).toUpper().toLatin1());
            next_cap = false;
        } else {
            result.append(*name);
        }

        name++;
    }

    return result;
}

/**
 * @brief SettingValue 根据给定信息返回获取的值
 * @param schema_id The id of the schema
 * @param path If non-empty, specifies the path for a relocatable schema
 * @param key 对应信息的key值
 * @param fallback 如果找不到信息，返回此默认值
 * @return
 */
inline const QVariant SettingValue(const QString &schema_id, const QByteArray &path = QByteArray(), const QString &key = QString(), const QVariant &fallback = QVariant())
{
    const QGSettings *settings = SettingsPtr(schema_id, path);

    if (settings && ((settings->keys().contains(key)) || settings->keys().contains(qtify_name(key.toUtf8().data())))) {
        QVariant v = settings->get(key);
        delete settings;
        return v;
    } else{
        qDebug() << "Cannot find gsettings, schema_id:" << schema_id
                 << " path:" << path << " key:" << key
                 << "Use fallback value:" << fallback;
        // 如果settings->keys()不包含key则会存在内存泄露，所以需要释放
        if (settings)
            delete settings;

        return fallback;
    }
}

inline bool SettingSaveValue(const QString &schema_id, const QByteArray &path, const QString &key, const QVariant &value)
{
    QGSettings *settings = SettingsPtr(schema_id, path);

    if (settings && ((settings->keys().contains(key)) || settings->keys().contains(qtify_name(key.toUtf8().data())))) {
        settings->set(key, value);
        delete settings;
        return true;
    } else{
        qDebug() << "Cannot find gsettings, schema_id:" << schema_id
                 << " path:" << path << " key:" << key;
        if (settings)
            delete settings;

        return false;
    }
}

inline QPixmap renderSVG(const QString &path, const QSize &size, const qreal devicePixelRatio)
{
    QImageReader reader;
    QPixmap pixmap;
    reader.setFileName(path);
    if (reader.canRead()) {
        reader.setScaledSize(size * devicePixelRatio);
        pixmap = QPixmap::fromImage(reader.read());
        pixmap.setDevicePixelRatio(devicePixelRatio);
    }
    else {
        pixmap.load(path);
    }

    return pixmap;
}

inline QScreen *screenAt(const QPoint &point) {
    for (QScreen *screen : qApp->screens()) {
        const QRect r { screen->geometry() };
        const QRect rect { r.topLeft(), r.size() * screen->devicePixelRatio() };
        if (rect.contains(point)) {
            return screen;
        }
    }

    return nullptr;
}

//!!! 注意:这里传入的QPoint是未计算缩放的
inline QScreen *screenAtByScaled(const QPoint &point) {
    for (QScreen *screen : qApp->screens()) {
        const QRect r { screen->geometry() };
        QRect rect { r.topLeft(), r.size() * screen->devicePixelRatio() };
        if (rect.contains(point)) {
            return screen;
        }
    }

    return nullptr;
}

/**
* @brief 比较两个插件版本号的大小
* @param pluginApi1 第一个插件版本号
* @param pluginApi2 第二个插件版本号
* @return 0:两个版本号相等,1:第一个版本号大,-1:第二个版本号大
*/
inline int comparePluginApi(const QString &pluginApi1, const QString &pluginApi2)
{
    // 版本号相同
    if (pluginApi1 == pluginApi2)
        return 0;

    // 拆分版本号
    QStringList subPluginApis1 = pluginApi1.split(".", Qt::SkipEmptyParts, Qt::CaseSensitive);
    QStringList subPluginApis2 = pluginApi2.split(".", Qt::SkipEmptyParts, Qt::CaseSensitive);
    for (int i = 0; i < subPluginApis1.size(); ++i) {
        auto subPluginApi1 = subPluginApis1[i];
        if (subPluginApis2.size() > i) {
            auto subPluginApi2 = subPluginApis2[i];

            // 相等判断下一个子版本号
            if (subPluginApi1 == subPluginApi2)
                continue;

            // 转成整形比较
            if (subPluginApi1.toInt() > subPluginApi2.toInt()) {
                return 1;
            } else {
                return -1;
            }
        }
    }

    // 循环结束但是没有返回,说明子版本号个数不同,且前面的子版本号都相同
    // 子版本号多的版本号大
    if (subPluginApis1.size() > subPluginApis2.size()) {
        return 1;
    } else {
        return -1;
    }
}

inline void updateCursor(QWidget *w)
{
    static QString lastCursorTheme;
    static int lastCursorSize = 0;
    QString theme = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-name", "bloom").toString();
    int cursorSize = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-size", 24).toInt();
    if (theme != lastCursorTheme || cursorSize != lastCursorSize) {
        QCursor *cursor = ImageUtil::loadQCursorFromX11Cursor(theme.toStdString().c_str(), "left_ptr", cursorSize);
        if (!cursor)
            return;
        lastCursorTheme = theme;
        lastCursorSize = cursorSize;
        w->setCursor(*cursor);
        static QCursor *lastArrowCursor = nullptr;
        if (lastArrowCursor)
            delete lastArrowCursor;

        lastArrowCursor = cursor;
    }
}

}

#endif // UTILS
