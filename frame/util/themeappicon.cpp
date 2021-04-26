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

#include "themeappicon.h"
#include "imageutil.h"

#include <QIcon>
#include <QFile>
#include <QDebug>
#include <QApplication>
#include <QPixmapCache>
#include <QCryptographicHash>
#include <QHBoxLayout>
#include <QLabel>
#include <QDate>
#include <QPainter>

#include <private/qguiapplication_p.h>
#include <private/qiconloader_p.h>
#include <qpa/qplatformtheme.h>

ThemeAppIcon::ThemeAppIcon(QObject *parent) : QObject(parent)
{

}

ThemeAppIcon::~ThemeAppIcon()
{

}

/**
 * @brief ThemeAppIcon::getIcon 根据传入的\a name 参数重新从系统主题中获取一次图标
 * @param name 图标名
 * @return 获取到的图标
 * @note 之所以不使用QIcon::fromTheme是因为这个函数中有缓存机制，获取系统主题中的图标的时候，第一次获取不到，下一次也是获取不到
 */
QIcon ThemeAppIcon::getIcon(const QString &name)
{
    QIcon icon;

    QPlatformTheme * const platformTheme = QGuiApplicationPrivate::platformTheme();
    bool hasUserTheme = QIconLoader::instance()->hasUserTheme();

    if (!platformTheme || hasUserTheme)
        return QIcon::fromTheme(name);

    QIconEngine * const engine = platformTheme->createIconEngine(name);
    QIcon *cachedIcon  = new QIcon(engine);
    icon = *cachedIcon;
    return icon;
}

bool ThemeAppIcon::getIcon(QPixmap &pix, const QString iconName, const int size, bool reObtain)
{
    QString key;
    QIcon icon;
    bool ret = true;
    // 把size改为小于size的最大偶数 :)
    const int s = int(size * qApp->devicePixelRatio()) & ~1;

    if (iconName == "dde-calendar") {
        const double iconZoom =  s / 256.0;
        QDate const date(QDate::currentDate());

        QPixmap pixmap(":/indicator/resources/calendar_bg.svg");
        pixmap = pixmap.scaled(s, s, Qt::KeepAspectRatio, Qt::SmoothTransformation);

        QPainter painter(&pixmap);
        painter.setRenderHints(QPainter::Antialiasing | QPainter::TextAntialiasing | QPainter::SmoothPixmapTransform);

        //根据不同日期显示不同日历图表
        int tw = pixmap.rect().width();
        int th = pixmap.rect().height();
        int tx = pixmap.rect().x();
        int ty = pixmap.rect().y();

        //绘制月份
        QRectF rcMonth(tx + (tw / 3.4), ty + (th / 5.4), 80 * iconZoom, 40 * iconZoom);
        painter.drawPixmap(rcMonth.topLeft(), ImageUtil::loadSvg(QString(":/icons/resources/month%1.svg").arg(date.month()), rcMonth.size().toSize()));
        //绘制日
        QRectF rcDay(tx + (tw / 3.5), ty + th / 3.1, 112 * iconZoom, 104 * iconZoom);
        painter.drawPixmap(rcDay.topLeft(), ImageUtil::loadSvg(QString(":/icons/resources/day%1.svg").arg(date.day()), rcDay.size().toSize()));
        //绘制周
        QRectF rcWeek(tx + (tw / 2.3), ty + ((th / 3.9) * 2.8), 56 * iconZoom, 24 * iconZoom);
        painter.drawPixmap(rcWeek.topLeft(), ImageUtil::loadSvg(QString(":/icons/resources/week%1.svg").arg(date.dayOfWeek()), rcWeek.size().toSize()));

        pix = pixmap;
        pix.setDevicePixelRatio(qApp->devicePixelRatio());

        return ret;
    }

    do {
        // load pixmap from our Cache
        if (iconName.startsWith("data:image/")) {
            key = QCryptographicHash::hash(iconName.toUtf8(), QCryptographicHash::Md5).toHex();

            // FIXME(hualet): The cache can reduce memory usage,
            // that is ~2M on HiDPI enabled machine with 9 icons loaded,
            // but I don't know why since QIcon has its own cache and all of the
            // icons loaded are loaded by QIcon::fromTheme, really strange here.
            if (QPixmapCache::find(key, &pix))
                break;
        }

        // load pixmap from Byte-Data
        if (iconName.startsWith("data:image/")) {
            const QStringList strs = iconName.split("base64,");
            if (strs.size() == 2)
                pix.loadFromData(QByteArray::fromBase64(strs.at(1).toLatin1()));

            if (!pix.isNull())
                break;
        }

        // load pixmap from File
        if (QFile::exists(iconName)) {
            pix = QPixmap(iconName);
            if (!pix.isNull())
                break;
        }

        // 重新从主题中获取一次
        if (reObtain)
            icon = getIcon(iconName);
        else
            icon = QIcon::fromTheme(iconName);

        if(icon.isNull()) {
            icon = QIcon::fromTheme("application-x-desktop");
            ret = false;
        }

        // load pixmap from Icon-Theme
        const int fakeSize = std::max(48, s); // cannot use 16x16, cause 16x16 is label icon
        pix = icon.pixmap(QSize(fakeSize, fakeSize));
        if (!pix.isNull())
            break;

        // fallback to a Default pixmap
        pix = QPixmap(":/icons/resources/application-x-desktop.svg");
        if (!pix.isNull())
            break;

        Q_UNREACHABLE();

    } while (false);

    if (!key.isEmpty()) {
        QPixmapCache::insert(key, pix);
    }

    if (pix.size().width() != s) {
        pix = pix.scaled(s, s, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    }
    pix.setDevicePixelRatio(qApp->devicePixelRatio());

    return ret;
}
