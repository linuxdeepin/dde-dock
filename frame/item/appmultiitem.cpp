/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#include "appitem.h"
#include "appmultiitem.h"
#include "imageutil.h"
#include "themeappicon.h"

#include <QBitmap>
#include <QMenu>
#include <QPixmap>
#include <QX11Info>

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>
#include <sys/shm.h>

struct SHMInfo {
    long shmid;
    long width;
    long height;
    long bytesPerLine;
    long format;

    struct Rect {
        long x;
        long y;
        long width;
        long height;
    } rect;
};

AppMultiItem::AppMultiItem(AppItem *appItem, WId winId, const WindowInfo &windowInfo, QWidget *parent)
    : DockItem(parent)
    , m_appItem(appItem)
    , m_windowInfo(windowInfo)
    , m_entryInter(appItem->itemEntryInter())
    , m_winId(winId)
    , m_menu(new QMenu(this))
{
    initMenu();
    initConnection();
}

AppMultiItem::~AppMultiItem()
{
}

QSize AppMultiItem::suitableSize(int size) const
{
    return QSize(size, size);
}

AppItem *AppMultiItem::appItem() const
{
    return m_appItem;
}

quint32 AppMultiItem::winId() const
{
    return m_winId;
}

const WindowInfo &AppMultiItem::windowInfo() const
{
    return m_windowInfo;
}

DockItem::ItemType AppMultiItem::itemType() const
{
    return DockItem::AppMultiWindow;
}

bool AppMultiItem::isKWinAvailable() const
{
    if (QDBusConnection::sessionBus().interface()->isServiceRegistered(QStringLiteral("org.kde.KWin"))) {
        QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/Effects"), QStringLiteral("org.kde.kwin.Effects"));
        QDBusReply<bool> reply = interface.call(QStringLiteral("isEffectLoaded"), "screenshot");

        return reply.value();
    }
    return false;
}

QImage AppMultiItem::snapImage() const
{
    // 优先使用窗管进行窗口截图
    if (isKWinAvailable()) {
        QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/Screenshot"), QStringLiteral("org.kde.kwin.Screenshot"));

        QList<QVariant> args;
        args << QVariant::fromValue(m_winId);
        args << QVariant::fromValue(quint32(width() - 20));
        args << QVariant::fromValue(quint32(height() - 20));

        QImage image;
        QDBusReply<QString> reply = interface.callWithArgumentList(QDBus::Block, QStringLiteral("screenshotForWindowExtend"), args);
        if(reply.isValid()){
            const QString tmpFile = reply.value();
            if (QFile::exists(tmpFile)) {
                image.load(tmpFile);
                qDebug() << "reply: " << tmpFile;
                QFile::remove(tmpFile);
            } else {
                qDebug() << "get current workspace bckground error, file does not exist : " << tmpFile;
            }
        } else {
            qDebug() << "get current workspace bckground error: "<< reply.error().message();
        }
        return image;
    }

    // get window image from shm(only for deepin app)
    SHMInfo *info = getImageDSHM();
    QImage image;
    uchar *image_data = 0;
    if (info) {
        qDebug() << "get Image from dxcbplugin SHM...";
        image_data = (uchar *)shmat(info->shmid, 0, 0);
        if ((qint64)image_data != -1)
            return QImage(image_data, info->width, info->height, info->bytesPerLine, (QImage::Format)info->format);

        qDebug() << "invalid pointer of shm!";
        image_data = nullptr;
    }

    QImage qimage;
    XImage *ximage;
    if (!image_data || qimage.isNull()) {
        ximage = getImageXlib();
        if (!ximage)
            return QImage();

        qimage = QImage((const uchar *)(ximage->data), ximage->width, ximage->height, ximage->bytes_per_line, QImage::Format_RGB32);
    }

    return image;
}

SHMInfo *AppMultiItem::getImageDSHM() const
{
    const auto display = Utils::IS_WAYLAND_DISPLAY ? XOpenDisplay(nullptr) : QX11Info::display();
    if (!display) {
        qWarning() << "Error: get display failed!";
        return nullptr;
    }

    Atom atom_prop = XInternAtom(display, "_DEEPIN_DXCB_SHM_INFO", true);
    if (!atom_prop) {
        return nullptr;
    }

    Atom actual_type_return_deepin_shm;
    int actual_format_return_deepin_shm;
    unsigned long nitems_return_deepin_shm;
    unsigned long bytes_after_return_deepin_shm;
    unsigned char *prop_return_deepin_shm;

    XGetWindowProperty(display, m_winId, atom_prop, 0, 32 * 9, false, AnyPropertyType,
                       &actual_type_return_deepin_shm, &actual_format_return_deepin_shm, &nitems_return_deepin_shm,
                       &bytes_after_return_deepin_shm, &prop_return_deepin_shm);

    return reinterpret_cast<SHMInfo *>(prop_return_deepin_shm);
}

XImage *AppMultiItem::getImageXlib() const
{
    const auto display = Utils::IS_WAYLAND_DISPLAY ? XOpenDisplay(nullptr) : QX11Info::display();
    if (!display) {
        qWarning() << "Error: get display failed!";
        return nullptr;
    }

    Window unused_window;
    int unused_int;
    unsigned unused_uint, w, h;
    XGetGeometry(display, m_winId, &unused_window, &unused_int, &unused_int, &w, &h, &unused_uint, &unused_uint);
    return XGetImage(display, m_winId, 0, 0, w, h, AllPlanes, ZPixmap);
}

void AppMultiItem::initMenu()
{
    QAction *actionOpen = new QAction(m_menu);
    actionOpen->setText(tr("Open"));
    connect(actionOpen, &QAction::triggered, this, &AppMultiItem::onOpen);
    m_menu->addAction(actionOpen);
}

void AppMultiItem::initConnection()
{
    connect(m_entryInter, &DockEntryInter::CurrentWindowChanged, this, &AppMultiItem::onCurrentWindowChanged);
}

void AppMultiItem::onOpen()
{
#ifdef USE_AM
    m_entryInter->ActiveWindow(m_winId);
#endif
}

void AppMultiItem::onCurrentWindowChanged(uint32_t value)
{
    if (value != m_winId)
        return;

    update();
}

void AppMultiItem::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setRenderHint(QPainter::SmoothPixmapTransform, true);

    if (m_snapImage.isNull()) {
        m_snapImage = snapImage();
    }

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);
    QRect itemRect = rect();
    itemRect.marginsRemoved(QMargins(6, 6, 6, 6));
    QPixmap pixmapWindowIcon = QPixmap::fromImage(m_snapImage);
    QPainterPath path;
    path.addRoundedRect(rect(), radius, radius);
    painter.fillPath(path, Qt::transparent);

    if (m_entryInter->currentWindow() == m_winId) {
        QColor backColor = Qt::black;
        backColor.setAlpha(255 * 0.8);
        painter.fillPath(path, backColor);
    }

    itemRect = m_snapImage.rect();
    int itemWidth = itemRect.width();
    int itemHeight = itemRect.height();
    int x = (rect().width() - itemWidth) / 2;
    int y = (rect().height() - itemHeight) / 2;
    painter.drawPixmap(QRect(x, y, itemWidth, itemHeight), pixmapWindowIcon);

    QPixmap pixmapAppIcon;
    ThemeAppIcon::getIcon(pixmapAppIcon, m_entryInter->icon(), qMin(width(), height()) * 0.8);
    if (!pixmapAppIcon.isNull()) {
        // 绘制下方的图标，下方的小图标大约为应用图标的三分之一的大小
        //pixmap = pixmap.scaled(pixmap.width() * 0.3, pixmap.height() * 0.3);
        QRect rectIcon = rect();
        int iconWidth = rectIcon.width() * 0.3;
        int iconHeight = rectIcon.height() * 0.3;
        rectIcon.setX((rect().width() - iconWidth) * 0.5);
        rectIcon.setY(rect().height() - iconHeight);
        rectIcon.setWidth(iconWidth);
        rectIcon.setHeight(iconHeight);
        painter.drawPixmap(rectIcon, pixmapAppIcon);
    }
}

void AppMultiItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton) {
#ifdef USE_AM
        m_entryInter->ActiveWindow(m_winId);
#endif
    } else {
        QPoint currentPoint = QCursor::pos();
        m_menu->exec(currentPoint);
    }
}
