/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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
#ifndef UTILS
#define UTILS
#include <QPixmap>
#include <QImageReader>
#include <QApplication>
#include <QScreen>
#include <QGSettings>

namespace Utils {

#define ICBC_CONF_FILE "/etc/deepin/icbc.conf"

inline QPixmap renderSVG(const QString &path, const QSize &size, const qreal devicePixelRatio) {
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

inline QScreen * screenAt(const QPoint &point) {
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
inline QScreen * screenAtByScaled(const QPoint &point) {
    for (QScreen *screen : qApp->screens()) {
        const QRect r { screen->geometry() };
        QRect rect { r.topLeft(), r.size() * screen->devicePixelRatio() };
        if (rect.contains(point)) {
            return screen;
        }
    }

    return nullptr;
}
    
inline bool isSettingConfigured(const QString& id, const QString& path, const QString& keyName) {
    if (!QGSettings::isSchemaInstalled(id.toUtf8())) {
        return false;
    }
    QGSettings setting(id.toUtf8(), path.toUtf8());
    QVariant v = setting.get(keyName);
    if (!v.isValid()) {
        return false;
    }
    return v.toBool();
}
}

#endif // UTILS
