/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef THEMEAPPICON_H
#define THEMEAPPICON_H

#include <QObject>
#include <QIcon>
#include <QMap>

class ThemeAppIcon : public QObject
{
    Q_OBJECT
public:
    explicit ThemeAppIcon(QObject *parent = nullptr);
    ~ThemeAppIcon();

    static void insertCache(const QString& iconName, const QIcon& icon) {
        if(!m_iconCache.contains(iconName)) m_iconCache.insert(iconName, icon);
    }

    static void removeCache(const QString& iconName) {
        if(m_iconCache.contains(iconName)) m_iconCache.remove(iconName);
    }

    static const QPixmap getIcon(const QString iconName, const int size, const qreal ratio);

private:
    static QMap<QString, QIcon> m_iconCache;
};

#endif // THEMEAPPICON_H
