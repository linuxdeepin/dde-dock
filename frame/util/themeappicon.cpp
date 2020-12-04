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

#define APPICONPATH1 "/usr/share/icons/hicolor/"
#define APPICONPATH2 "/usr/share/icons/"

QMap<QString, QIcon> ThemeAppIcon::m_iconCache;

ThemeAppIcon::ThemeAppIcon(QObject *parent) : QObject(parent)
{

}

ThemeAppIcon::~ThemeAppIcon()
{

}

QIcon ThemeAppIcon::_getIconFromDir(QString dirName, QString subdirName, QString iconName, int fakeSize)
{
    QIcon icon;
    QString filename;
    int sizelist[] = {24,32,48,64,96,128,256,512};
    int sl = 8;
    for (int i = 0; i < sl; i++) {
        if (fakeSize <= sizelist[i]) {
            filename.clear();
            filename = dirName;
            if (!subdirName.isEmpty()) {
                filename.append(QString::number(sizelist[i]) + 'x' + QString::number(sizelist[i]));
                filename.append(subdirName);
            }
            filename.append(iconName);
            if (!QFile::exists(filename)) {
                continue;
            } else {
                icon.addFile(filename,QSize(sizelist[i],sizelist[i]));
                break;
            }
        }
    }
    return icon;
}

//手动查找第三方应用的图标
QPixmap ThemeAppIcon::manualGetAppIcon(QString iconName,int fakeSize)
{
    QIcon icon;
    QPixmap pixmap;

    if (iconName.isEmpty()) {
        qDebug() << "-->Get Icon is Empty";
        return pixmap;
    }
    iconName.append(".png");

    //例在/usr/share/icons/hicolor/48x48/apps/目录下查找图片
    //例在/usr/share/icons/hicolor/48x48/mimetypes/目录下查找图片
    //例在/usr/share/icons/目录下查找图片
    icon = _getIconFromDir(APPICONPATH1,"/apps/",iconName,fakeSize);
    if (icon.isNull()) {
        icon = _getIconFromDir(APPICONPATH1,"/mimetypes/",iconName,fakeSize);
        if (icon.isNull()) {
            icon = _getIconFromDir(APPICONPATH2,"",iconName,fakeSize);
        }
    }

    if (icon.isNull()) {
        qDebug() << "-->Cannot find:"<< iconName << ", used default icon :application-x-desktop";
    } else {
        qDebug() << "-->Find:"<< iconName << ", and setSize=" << fakeSize;
        pixmap = icon.pixmap(QSize(fakeSize,fakeSize));
    }
    return pixmap;
}

const QPixmap ThemeAppIcon::getIcon(const QString iconName, const int size, const qreal ratio)
{
    QPixmap pixmap;
    QString key;
    QIcon icon;
    // 把size改为小于size的最大偶数 :)
    const int s = int(size * ratio) & ~1;
    const float iconZoom = size / 64.0 * 0.8;

    if (iconName == "dde-calendar") {
        QDate const date(QDate::currentDate());

        auto calendar = new QWidget() ;
        calendar->setFixedSize(s, s);

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
        auto monthPix = ImageUtil::loadSvg(QString(":/icons/resources/month%1.svg").arg(date.month()), QSize(36, 16)*iconZoom, ratio);
        month->setPixmap(monthPix.scaled(monthPix.width()*ratio,monthPix.height()*ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation));
        month->setFixedHeight(month->pixmap()->height());
        month->setAlignment(Qt::AlignCenter);
        month->setFixedWidth(s - 5 * iconZoom);
        layout->addWidget(month, Qt::AlignVCenter);

        auto day = new QLabel();
        auto dayPix =ImageUtil::loadSvg(QString(":/icons/resources/day%1.svg").arg(date.day()), QSize(32, 30)*iconZoom, ratio);
        day->setPixmap(dayPix.scaled(dayPix.width()*ratio,dayPix.height()*ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation));
        day->setAlignment(Qt::AlignCenter);
        day->setFixedHeight(day->pixmap()->height()/ratio);
        day->raise();
        layout->addWidget(day, Qt::AlignVCenter);

        auto week = new QLabel();
        auto weekPix = ImageUtil::loadSvg(QString(":/icons/resources/week%1.svg").arg(date.dayOfWeek()), QSize(26, 13)*iconZoom, ratio);
        week->setPixmap(weekPix.scaled(weekPix.width()*ratio,weekPix.height()*ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation));
        week->setFixedHeight(week->pixmap()->height());
        week->setAlignment(Qt::AlignCenter);
        week->setFixedWidth(s + 5 * iconZoom);
        layout->addWidget(week, Qt::AlignVCenter);
        layout->setSpacing(0);
        layout->setContentsMargins(0, 10 * iconZoom, 0, 10 * iconZoom);
        calendar->setLayout(layout);
        pixmap = calendar->grab(calendar->rect());
        if (pixmap.size().width() != s) {
            pixmap = pixmap.scaled(s, s, Qt::KeepAspectRatio, Qt::SmoothTransformation);
        }
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
        if (iconName.startsWith("data:image/")) {
            const QStringList strs = iconName.split("base64,");
            if (strs.size() == 2)
                pixmap.loadFromData(QByteArray::fromBase64(strs.at(1).toLatin1()));

            if (!pixmap.isNull())
                break;
        }

        // load pixmap from File
        if (QFile::exists(iconName)) {
            pixmap = QPixmap(iconName);
            if (!pixmap.isNull())
                break;
        }

        icon = QIcon::fromTheme(iconName);
        if (icon.isNull()) {
            icon = QIcon::fromTheme("deepinwine-" + iconName);
        } else {
            icon = QIcon::fromTheme(iconName, QIcon::fromTheme("application-x-desktop"));
        }

        const int fakeSize = std::max(48, s);
        if (icon.isNull()) {
            qDebug() << "-->Qt Function Cannot Find Icon:" << iconName;
            pixmap = manualGetAppIcon(iconName, fakeSize);
            if (!pixmap.isNull())
                break;
        } else {
            // cannot use 16x16, cause 16x16 is label icon
            pixmap = icon.pixmap(QSize(fakeSize, fakeSize));
            if (!pixmap.isNull())
                break;
        }

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
