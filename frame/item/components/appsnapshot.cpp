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

#include "appsnapshot.h"
#include "previewcontainer.h"

#include <DStyle>

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>
#include <sys/shm.h>

#include <QX11Info>
#include <QPainter>
#include <QVBoxLayout>
#include <QSizeF>
#include <QTimer>

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

AppSnapshot::AppSnapshot(const WId wid, const QDBusObjectPath &entry,QWidget *parent)
    : QWidget(parent)
    , m_wid(wid)
    , m_title(new TipsWidget)
    , m_waitLeaveTimer(new QTimer(this))
    , m_closeBtn2D(new DImageButton)
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_itemEntryInter(new DockEntryInter("com.deepin.dde.daemon.Dock", entry.path(), QDBusConnection::sessionBus(), this))
{
    m_closeBtn2D->setFixedSize(24, 24);
    m_closeBtn2D->setObjectName("Btn_closebutton-2d");
    m_closeBtn2D->setNormalPic(":/icons/resources/close_round_normal.svg");
    m_closeBtn2D->setHoverPic(":/icons/resources/close_round_hover.svg");
    m_closeBtn2D->setPressPic(":/icons/resources/close_round_press.svg");
    m_closeBtn2D->setVisible(false);
    m_title->setObjectName("AppSnapshotTitle");

    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->addWidget(m_title);
    centralLayout->addWidget(m_closeBtn2D);
    centralLayout->setSpacing(5);
    centralLayout->setMargin(0);

    centralLayout->setAlignment(m_closeBtn2D, Qt::AlignRight);

    setLayout(centralLayout);
    setAcceptDrops(true);
    resize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_closeBtn2D, &DImageButton::clicked, this, &AppSnapshot::closeWindow, Qt::QueuedConnection);
    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &AppSnapshot::compositeChanged, Qt::QueuedConnection);
    QTimer::singleShot(1, this, &AppSnapshot::compositeChanged);
}

void AppSnapshot::closeWindow() const
{
    const auto display = QX11Info::display();

    XEvent e;

    memset(&e, 0, sizeof(e));
    e.xclient.type = ClientMessage;
    e.xclient.window = m_wid;
    e.xclient.message_type = XInternAtom(display, "WM_PROTOCOLS", true);
    e.xclient.format = 32;
    e.xclient.data.l[0] = XInternAtom(display, "WM_DELETE_WINDOW", false);
    e.xclient.data.l[1] = CurrentTime;

    XSendEvent(display, m_wid, false, NoEventMask, &e);
    XFlush(display);
}

void AppSnapshot::compositeChanged() const
{
    const bool composite = m_wmHelper->hasComposite();

    m_title->setVisible(!composite);

    QTimer::singleShot(1, this, &AppSnapshot::fetchSnapshot);
}

void AppSnapshot::setWindowInfo(const WindowInfo &info)
{
    m_windowInfo = info;
    QFontMetrics fm(m_title->font());
    QString strTtile = m_title->fontMetrics().elidedText(m_windowInfo.title, Qt::ElideRight, width());
    m_title->setText(strTtile);
}

void AppSnapshot::dragEnterEvent(QDragEnterEvent *e)
{
    QWidget::dragEnterEvent(e);

    if (m_wmHelper->hasComposite())
        emit entered(m_wid);
}

void AppSnapshot::fetchSnapshot()
{
    if (!m_wmHelper->hasComposite())
        return;

    QImage qimage;
    SHMInfo *info = nullptr;
    uchar *image_data = nullptr;
    XImage *ximage = nullptr;
    unsigned char *prop_to_return_gtk = nullptr;

    // xcb_window_t activeWindow = KWindowSystem::activeWindow();
     // if KWin is available, use the KWin DBus interfaces
     if (isKWinAvailable()) {
         QDBusConnection bus = QDBusConnection::sessionBus();
          auto  reuslt =  bus.connect(QStringLiteral("org.kde.KWin"),
                     QStringLiteral("/Screenshot"),
                     QStringLiteral("org.kde.kwin.Screenshot"),
                     QStringLiteral("screenshotCreated"),
                     this, SLOT(KWinDBusScreenshotHelper(quint64)));



         QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/Screenshot"), QStringLiteral("org.kde.kwin.Screenshot"));
         qDebug() << "windowsID:"<< m_itemEntryInter->currentWindow()<<reuslt;

         interface.call(QStringLiteral("screenshotForWindow"), (quint64)m_itemEntryInter->currentWindow());
         return;
     }

    do {
        // get window image from shm(only for deepin app)
        info = getImageDSHM();
        if (info) {
            qDebug() << "get Image from dxcbplugin SHM...";
            //qDebug() << info->shmid << info->width << info->height << info->bytesPerLine << info->format << info->rect.x << info->rect.y << info->rect.width << info->rect.height;
            image_data = (uchar *)shmat(info->shmid, 0, 0);
            if ((qint64)image_data != -1) {
                m_snapshot = QImage(image_data, info->width, info->height, info->bytesPerLine, (QImage::Format)info->format);
                m_snapshotSrcRect = QRect(info->rect.x, info->rect.y, info->rect.width, info->rect.height);
                break;
            }
            qDebug() << "invalid pointer of shm!";
            image_data = nullptr;
        }

        if (!image_data || qimage.isNull()) {
            // get window image from XGetImage(a little slow)
            qDebug() << "get Image from dxcbplugin SHM failed!";
            qDebug() << "get Image from Xlib...";
            ximage = getImageXlib();
            if (!ximage) {
                qDebug() << "get Image from Xlib failed! giving up...";
                emit requestCheckWindow();
                return;
            }
            qimage = QImage((const uchar *)(ximage->data), ximage->width, ximage->height, ximage->bytes_per_line, QImage::Format_RGB32);
        }

        Q_ASSERT(!qimage.isNull());

        // remove shadow frame
        m_snapshotSrcRect = rectRemovedShadow(qimage, prop_to_return_gtk);
        m_snapshot = qimage;
    } while (false);


   // xcb_window_t activeWindow = KWindowSystem::activeWindow();
    // if KWin is available, use the KWin DBus interfaces
    if (isKWinAvailable()) {
        QDBusConnection bus = QDBusConnection::sessionBus();
        bus.connect(QStringLiteral("org.kde.KWin"),
                    QStringLiteral("/Screenshot"),
                    QStringLiteral("org.kde.kwin.Screenshot"),
                    QStringLiteral("screenshotCreated"),
                    this, SLOT(KWinDBusScreenshotHelper(quint64)));
        QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/Screenshot"), QStringLiteral("org.kde.kwin.Screenshot"));


        //interface.call(QStringLiteral("screenshotForWindow"), (quint64)activeWindow);

        m_snapshotSrcRect = rectRemovedShadow(qimage, prop_to_return_gtk);
        m_snapshot = qimage;
        return;
    }

    QSizeF size(rect().marginsRemoved(QMargins(8, 8, 8, 8)).size());
    const auto ratio = devicePixelRatioF();
    size = m_snapshotSrcRect.size().scaled(size * ratio, Qt::KeepAspectRatio);
    qreal scale = qreal(size.width()) / m_snapshotSrcRect.width();
    m_snapshot = m_snapshot.scaled(qRound(m_snapshot.width() * scale), qRound(m_snapshot.height() * scale),
                                   Qt::IgnoreAspectRatio, Qt::SmoothTransformation);
    m_snapshotSrcRect.moveTop(m_snapshotSrcRect.top() * scale + 0.5);
    m_snapshotSrcRect.moveLeft(m_snapshotSrcRect.left() * scale + 0.5);
    m_snapshotSrcRect.setWidth(size.width() - 0.5);
    m_snapshotSrcRect.setHeight(size.height() - 0.5);
    m_snapshot.setDevicePixelRatio(ratio);

    if (image_data) shmdt(image_data);
    if (ximage) XDestroyImage(ximage);
    if (info) XFree(info);
    if (prop_to_return_gtk) XFree(prop_to_return_gtk);

    update();
}

bool AppSnapshot::isKWinAvailable()
{
    if (QDBusConnection::sessionBus().interface()->isServiceRegistered(QStringLiteral("org.kde.KWin"))) {
        QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/Effects"), QStringLiteral("org.kde.kwin.Effects"));
        QDBusReply<bool> reply = interface.call(QStringLiteral("isEffectLoaded"), "screenshot");

        return reply.value();
    }
    return false;
}

void AppSnapshot::KWinDBusScreenshotHelper(quint64 pixmapId)
{
    qDebug() <<"============================================KWinDBusScreenshotHelper";
    // obtain width and height and grab an image (x and y are always zero for pixmaps)
    QRect pixrect = getDrawableGeometry((xcb_drawable_t)pixmapId);
    QPixmap mPixmap = getPixmapFromDrawable((xcb_drawable_t)pixmapId, pixrect);

    if (!mPixmap.isNull()) {
        m_snapshotSrcRect = pixrect;
        m_snapshot = mPixmap.toImage();
    }

    QSizeF size(rect().marginsRemoved(QMargins(8, 8, 8, 8)).size());
    const auto ratio = devicePixelRatioF();
    size = m_snapshotSrcRect.size().scaled(size * ratio, Qt::KeepAspectRatio);
    qreal scale = qreal(size.width()) / m_snapshotSrcRect.width();
    m_snapshot = m_snapshot.scaled(qRound(m_snapshot.width() * scale), qRound(m_snapshot.height() * scale),
                                   Qt::IgnoreAspectRatio, Qt::SmoothTransformation);
    m_snapshotSrcRect.moveTop(m_snapshotSrcRect.top() * scale + 0.5);
    m_snapshotSrcRect.moveLeft(m_snapshotSrcRect.left() * scale + 0.5);
    m_snapshotSrcRect.setWidth(size.width() - 0.5);
    m_snapshotSrcRect.setHeight(size.height() - 0.5);
    m_snapshot.setDevicePixelRatio(ratio);

    update();

    // Cannot retrieve pixmap from KWin, just fallback to fullscreen capture. We
    // could try to detect the original action (window under cursor or active
    // window), but that is too complex for this edge case.
}

QRect AppSnapshot::getDrawableGeometry(xcb_drawable_t drawable)
{
    xcb_connection_t *xcbConn = QX11Info::connection();

    xcb_get_geometry_cookie_t geomCookie = xcb_get_geometry_unchecked(xcbConn, drawable);
    CScopedPointer<xcb_get_geometry_reply_t> geomReply(xcb_get_geometry_reply(xcbConn, geomCookie, NULL));

    if (geomReply.isNull()) {
        return QRect();
    }
    return QRect(geomReply->x, geomReply->y, geomReply->width, geomReply->height);
}

QPixmap AppSnapshot::getPixmapFromDrawable(xcb_drawable_t drawableId, const QRect &rect)
 {
     xcb_connection_t *xcbConn = QX11Info::connection();

     // proceed to get an image based on the geometry (in device pixels)

     QScopedPointer<xcb_image_t, ScopedPointerXcbImageDeleter> xcbImage(
         xcb_image_get(
             xcbConn,
             drawableId,
             rect.x(),
             rect.y(),
             rect.width(),
             rect.height(),
             ~0,
             XCB_IMAGE_FORMAT_Z_PIXMAP
         )
     );

     // too bad, the capture failed.
     if (xcbImage.isNull()) {
         return QPixmap();
     }

     // now process the image

     QPixmap nativePixmap = convertFromNative(xcbImage.data());
     return nativePixmap;
 }


QPixmap AppSnapshot::convertFromNative(xcb_image_t *xcbImage)
{
    QImage::Format format = QImage::Format_Invalid;

    switch (xcbImage->depth) {
    case 1:
        format = QImage::Format_MonoLSB;
        break;
    case 16:
        format = QImage::Format_RGB16;
        break;
    case 24:
        format = QImage::Format_RGB32;
        break;
    case 30:
        format = QImage::Format_BGR30;
        break;
    case 32:
        format = QImage::Format_ARGB32_Premultiplied;
        break;
    default:
        return QPixmap(); // we don't know
    }

    // The RGB32 format requires data format 0xffRRGGBB, ensure that this fourth byte really is 0xff
    if (format == QImage::Format_RGB32) {
        quint32 *data = reinterpret_cast<quint32 *>(xcbImage->data);
        for (int i = 0; i < xcbImage->width * xcbImage->height; i++) {
            data[i] |= 0xff000000;
        }
    }

    QImage image(xcbImage->data, xcbImage->width, xcbImage->height, format);

    if (image.isNull()) {
        return QPixmap();
    }

    // work around an abort in QImage::color

    if (image.format() == QImage::Format_MonoLSB) {
        image.setColorCount(2);
        image.setColor(0, QColor(Qt::white).rgb());
        image.setColor(1, QColor(Qt::black).rgb());
    }

    // Image is ready. Since the backing data from xcbImage could be freed
    // before the QPixmap goes away, a deep copy is necessary.
    return QPixmap::fromImage(image).copy();
}


void AppSnapshot::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    if (!m_wmHelper->hasComposite()) {
        m_closeBtn2D->setVisible(true);
    } else {
        emit entered(wid());
    }

    update();
}

void AppSnapshot::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_closeBtn2D->setVisible(false);

    update();
}

void AppSnapshot::paintEvent(QPaintEvent *e)
{
    QPainter painter(this);

    if (!m_wmHelper->hasComposite()) {
        if (underMouse())
            painter.fillRect(rect(), QColor(255, 255, 255, 255 * .2));
        return;
    }

    if (m_snapshot.isNull())
        return;

    const auto ratio = devicePixelRatioF();

    // draw attention background
    if (m_windowInfo.attention) {
        painter.setBrush(QColor(241, 138, 46, 255 * .8));
        painter.setPen(Qt::NoPen);
        painter.drawRoundedRect(rect(), 5, 5);
    }

    // draw image
    const QImage &im = m_snapshot;

    const qreal offset_x = width() / 2.0 - m_snapshotSrcRect.width() / ratio / 2 - m_snapshotSrcRect.left() / ratio;
    const qreal offset_y = height() / 2.0 - m_snapshotSrcRect.height() / ratio / 2 - m_snapshotSrcRect.top() / ratio;

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);

    QBrush brush;
    brush.setTextureImage(im);
    painter.setBrush(brush);
    painter.setPen(Qt::NoPen);
    painter.scale(1 / ratio, 1 / ratio);
    painter.translate(QPoint(offset_x * ratio, offset_y * ratio));
    painter.drawRoundedRect(m_snapshotSrcRect, radius * ratio, radius * ratio);
}

void AppSnapshot::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    QTimer::singleShot(1, this, &AppSnapshot::fetchSnapshot);
}

void AppSnapshot::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);

    emit clicked(m_wid);
}

SHMInfo *AppSnapshot::getImageDSHM()
{
    const auto display = QX11Info::display();

    Atom atom_prop = 0;
    if (!atom_prop) {
        return nullptr;
    }

    Atom actual_type_return_deepin_shm;
    int actual_format_return_deepin_shm;
    unsigned long nitems_return_deepin_shm;
    unsigned long bytes_after_return_deepin_shm;
    unsigned char *prop_return_deepin_shm;

    XGetWindowProperty(display, m_wid, atom_prop, 0, 32 * 9, false, AnyPropertyType,
                       &actual_type_return_deepin_shm, &actual_format_return_deepin_shm, &nitems_return_deepin_shm,
                       &bytes_after_return_deepin_shm, &prop_return_deepin_shm);

    //qDebug() << actual_type_return_deepin_shm << actual_format_return_deepin_shm << nitems_return_deepin_shm << bytes_after_return_deepin_shm << prop_return_deepin_shm;

    return reinterpret_cast<SHMInfo *>(prop_return_deepin_shm);
}

XImage *AppSnapshot::getImageXlib()
{
//    const auto display = QX11Info::display();
//    Window unused_window;
//    int unused_int;
//    unsigned unused_uint, w, h;
//    XGetGeometry(display, m_wid, &unused_window, &unused_int, &unused_int, &w, &h, &unused_uint, &unused_uint);
//    return XGetImage(display, m_wid, 0, 0, w, h, AllPlanes, ZPixmap);

    return nullptr;
}

QRect AppSnapshot::rectRemovedShadow(const QImage &qimage, unsigned char *prop_to_return_gtk)
{
    const auto display = QX11Info::display();

    const Atom gtk_frame_extents = XInternAtom(display, "_GTK_FRAME_EXTENTS", true);
    Atom actual_type_return_gtk;
    int actual_format_return_gtk;
    unsigned long n_items_return_gtk;
    unsigned long bytes_after_return_gtk;

    const auto r = XGetWindowProperty(display, m_wid, gtk_frame_extents, 0, 4, false, XA_CARDINAL,
                                      &actual_type_return_gtk, &actual_format_return_gtk, &n_items_return_gtk, &bytes_after_return_gtk, &prop_to_return_gtk);
    if (!r && prop_to_return_gtk && n_items_return_gtk == 4 && actual_format_return_gtk == 32) {
        qDebug() << "remove shadow frame...";
        const unsigned long *extents = reinterpret_cast<const unsigned long *>(prop_to_return_gtk);
        const int left = extents[0];
        const int right = extents[1];
        const int top = extents[2];
        const int bottom = extents[3];
        const int width = qimage.width();
        const int height = qimage.height();

        return QRect(left, top, width - left - right, height - top - bottom);
    } else {
        return QRect(0, 0, qimage.width(), qimage.height());
    }
}
