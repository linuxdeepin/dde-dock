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

#include <QIcon>
#include <QFile>
#include <QDebug>
#include <QApplication>
#include <QPixmapCache>
#include <QCryptographicHash>
#include <QHBoxLayout>
#include <QLabel>
#include <QDate>

#include "../frame/util/imageutil.h"

ThemeAppIcon::ThemeAppIcon(QObject *parent) : QObject(parent)
{

}

ThemeAppIcon::~ThemeAppIcon()
{

}

const QPixmap ThemeAppIcon::getIcon(const QString iconName, const int size, const qreal ratio)
{
    QPixmap pixmap;
    QString key;


    // 把size改为小于size的最大偶数 :)
    const int s = int(size * ratio) & ~1;
    const float iconZoom = size /64.0*0.8;

    if(iconName == "dde-calendar"){
        QDate const date(QDate::currentDate());

        auto pixday  =  ImageUtil::loadSvg(":/indicator/resources/calendar_bg.svg","",size,ratio);
        auto calendar = new QWidget() ;
        calendar->setFixedSize(s,s);

        calendar->setAutoFillBackground(true);
            QPalette palette = calendar->palette();
            palette.setBrush(QPalette::Window,
                    QBrush(QPixmap(":/indicator/resources/calendar_bg.svg").scaled(
                        calendar->size(),
                        Qt::IgnoreAspectRatio,
                        Qt::SmoothTransformation)));
         calendar->setPalette(palette);

        QVBoxLayout *layout = new QVBoxLayout;
        layout->setSpacing(0);
        auto month = new QLabel();
        month->setPixmap( ImageUtil::loadSvg(QString(":/icons/resources/month%1.svg").arg(date.month()),QSize(36,16)*iconZoom,ratio));
        month->setFixedHeight(month->pixmap()->height());
        month->setAlignment(Qt::AlignCenter);
        month->setFixedWidth(s-5*iconZoom);
        layout->addWidget(month,Qt::AlignVCenter);

        auto day = new QLabel();
        day->setPixmap( ImageUtil::loadSvg(QString(":/icons/resources/day%1.svg").arg(date.day()),QSize(32,30)*iconZoom,ratio));
        day->setAlignment(Qt::AlignCenter);
        day->setFixedHeight(day->pixmap()->height());
        day->raise();
        layout->addWidget(day,Qt::AlignVCenter);

        auto week = new QLabel();
        week->setPixmap( ImageUtil::loadSvg(QString(":/icons/resources/week%1.svg").arg(date.dayOfWeek()),QSize(26,13)*iconZoom,ratio));
        week->setFixedHeight(week->pixmap()->height());
        week->setAlignment(Qt::AlignCenter);
        week->setFixedWidth(s+5*iconZoom);
        layout->addWidget(week,Qt::AlignVCenter);
        layout->setSpacing(0);
        layout->setContentsMargins(0,10*iconZoom,0,10*iconZoom);
        calendar->setLayout(layout);
        pixmap = calendar->grab(calendar->rect());

        return pixmap;
    }


    do {

        // load pixmap from our Cache
        if (iconName.startsWith("data:image/")) {
            key = QCryptographicHash::hash(iconName.toUtf8(), QCryptographicHash::Md5).toHex();

            // FIXME(hualet): The cache can reduce memory usage,
            // that is ~2M on HiDPI enabled machine with 9 icons loaded,
            // but I don't know why since QIcon has its own cache and all of the
            // icons loaded are loaded by QIcon::fromTheme, really strange here.
            if (QPixmapCache::find(key, &pixmap))
                break;
        }

        // load pixmap from Byte-Data
        if (iconName.startsWith("data:image/"))
        {
            const QStringList strs = iconName.split("base64,");
            if (strs.size() == 2)
                pixmap.loadFromData(QByteArray::fromBase64(strs.at(1).toLatin1()));

            if (!pixmap.isNull())
                break;
        }

        // load pixmap from File
        if (QFile::exists(iconName))
        {
            pixmap = QPixmap(iconName);
            if (!pixmap.isNull())
                break;
        }

        // load pixmap from Icon-Theme
        const QIcon icon = QIcon::fromTheme(iconName, QIcon::fromTheme("application-x-desktop"));
        const int fakeSize = std::max(48, s); // cannot use 16x16, cause 16x16 is label icon
        pixmap = icon.pixmap(QSize(fakeSize, fakeSize));
        if (!pixmap.isNull())
            break;

        // fallback to a Default pixmap
        pixmap = QPixmap(":/icons/resources/application-x-desktop.svg");
        if (!pixmap.isNull())
            break;

        Q_UNREACHABLE();

    } while (false);

    if (!key.isEmpty()) {
        QPixmapCache::insert(key, pixmap);
    }

    if (pixmap.size().width() != s) {
        pixmap = pixmap.scaled(s, s, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    }
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
