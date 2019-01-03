/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "snitraywidget.h"
#include "util/themeappicon.h"

#include <QPainter>
#include <QApplication>

#include <xcb/xproto.h>

#define IconSize 16

const QStringList ItemCategoryList {"ApplicationStatus" , "Communications" , "SystemServices", "Hardware"};
const QStringList ItemStatusList {"ApplicationStatus" , "Communications" , "SystemServices", "Hardware"};
const QStringList LeftClickInvalidIdList {"sogou-qimpanel",};

SNITrayWidget::SNITrayWidget(const QString &sniServicePath, QWidget *parent)
    : AbstractTrayWidget(parent),
        m_dbusMenuImporter(nullptr),
        m_menu(nullptr),
        m_updateTimer(new QTimer(this))
{
    if (sniServicePath.startsWith("/") || !sniServicePath.contains("/")) {
        return;
    }

    QPair<QString, QString> pair = serviceAndPath(sniServicePath);
    m_dbusService = pair.first;
    m_dbusPath = pair.second;

    m_sniInter = new StatusNotifierItem(m_dbusService, m_dbusPath, QDBusConnection::sessionBus(), this);

    if (!m_sniInter->isValid()) {
        qDebug() << "SNI dbus interface is invalid!" << m_dbusService << m_dbusPath << m_sniInter->lastError();
        return;
    }

    m_updateTimer->setInterval(100);
    m_updateTimer->setSingleShot(true);

    connect(m_updateTimer, &QTimer::timeout, this, &SNITrayWidget::refreshIcon);
    connect(m_sniInter, &StatusNotifierItem::NewIcon, m_updateTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(m_sniInter, &StatusNotifierItem::NewOverlayIcon, this, &SNITrayWidget::refreshOverlayIcon);
    connect(m_sniInter, &StatusNotifierItem::NewAttentionIcon, this, &SNITrayWidget::refreshAttentionIcon);

    QTimer::singleShot(0, this, &SNITrayWidget::refreshIcon);
}

SNITrayWidget::~SNITrayWidget()
{
}

void SNITrayWidget::setActive(const bool active)
{
}

void SNITrayWidget::updateIcon()
{
    m_updateTimer->start();
}

void SNITrayWidget::sendClick(uint8_t mouseButton, int x, int y)
{
    switch (mouseButton) {
        case XCB_BUTTON_INDEX_1:
            // left button click invalid
            if (LeftClickInvalidIdList.contains(m_sniInter->id())) {
                showContextMenu(x, y);
            } else {
                m_sniInter->Activate(x, y);
            }
            break;
        case XCB_BUTTON_INDEX_2:
            m_sniInter->SecondaryActivate(x, y);
            break;
        case XCB_BUTTON_INDEX_3:
            showContextMenu(x, y);
            break;
        default:
            qDebug() << "unknown mouse button key";
            break;
    }
}

const QImage SNITrayWidget::trayImage()
{
    return m_pixmap.toImage();
}

bool SNITrayWidget::isValid()
{
    return m_sniInter->isValid();
}

QString SNITrayWidget::toSNIKey(const QString &sniServicePath)
{
    QString key;

    do {
        const QPair<QString, QString> &sap = serviceAndPath(sniServicePath);
        key = QDBusInterface(sap.first, sap.second).property("Id").toString();
        if (!key.isEmpty()) {
            break;
        }

        key = sniServicePath;
    } while (false);

    return QString("sni:%1").arg(key);
}

bool SNITrayWidget::isSNIKey(const QString &itemKey)
{
    return itemKey.startsWith("sni:");
}

QPair<QString, QString> SNITrayWidget::serviceAndPath(const QString &servicePath)
{
    QStringList list = servicePath.split("/");
    QPair<QString, QString> pair;
    pair.first = list.takeFirst();

    for (auto i : list) {
        pair.second.append("/");
        pair.second.append(i);
    }

    return pair;
}

void SNITrayWidget::initMenu()
{
    qDebug() << "using sni service path:" << m_dbusService;

    const QString &menuPath = m_sniInter->menu().path();

    qDebug() << "using sni menu path:" << menuPath;

    m_dbusMenuImporter = new DBusMenuImporter(m_dbusService, menuPath, ASYNCHRONOUS, this);

    qDebug() << "generate the sni menu object";

    m_menu = m_dbusMenuImporter->menu();

    qDebug() << "the sni menu obect is:" << m_menu;
}

/*
 *ItemCategory SNITrayWidget::category()
 *{
 *    const QString &category = m_sniInter->category();
 *    if (!ItemCategoryList.contains(category)) {
 *        return UnknownCategory;
 *    }
 *
 *    return static_cast<ItemCategory>(ItemCategoryList.indexOf(category));
 *}
 *
 *ItemStatus SNITrayWidget::status()
 *{
 *    const QString &status = m_sniInter->status();
 *    if (!ItemStatusList.contains(status)) {
 *        return UnknownStatus;
 *    }
 *
 *    return static_cast<ItemStatus>(ItemStatusList.indexOf(status));
 *}
 *
 */

void SNITrayWidget::refreshIcon()
{
    QPixmap pix = newIconPixmap(Icon);
    if (pix.isNull()) {
        return;
    }

    m_pixmap = pix;
    update();
    Q_EMIT iconChanged();

    if (!isVisible()) {
        Q_EMIT needAttention();
    }
}

void SNITrayWidget::refreshOverlayIcon()
{
    QPixmap pix = newIconPixmap(OverlayIcon);
    if (pix.isNull()) {
        return;
    }

    m_overlayPixmap = pix;
    update();
    Q_EMIT iconChanged();

    if (!isVisible()) {
        Q_EMIT needAttention();
    }
}

void SNITrayWidget::refreshAttentionIcon()
{
    /* TODO: A new approach may be needed to deal with attentionIcon */
    QPixmap pix = newIconPixmap(AttentionIcon);
    if (pix.isNull()) {
        return;
    }

    m_pixmap = pix;
    update();
    Q_EMIT iconChanged();

    if (!isVisible()) {
        Q_EMIT needAttention();
    }
}

void SNITrayWidget::showContextMenu(int x, int y)
{
    // ContextMenu does not work
    if (m_sniInter->menu().path().startsWith("/NO_DBUSMENU")) {
        m_sniInter->ContextMenu(x, y);
    } else {
        if (!m_menu) {
            qDebug() << "context menu has not be ready, init menu";
            initMenu();
        }
        m_menu->popup(QPoint(x, y));
    }
}

QSize SNITrayWidget::sizeHint() const
{
    return QSize(26, 26);
}

void SNITrayWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);
    if (m_pixmap.isNull())
        return;

    QPainter painter;
    painter.begin(this);
    painter.setRenderHint(QPainter::Antialiasing);
#ifdef QT_DEBUG
//    painter.fillRect(rect(), Qt::red);
#endif

    const QRectF &rf = QRect(rect());
    const QRectF &rfp = QRect(m_pixmap.rect());
    const QPointF &p = rf.center() - rfp.center() / m_pixmap.devicePixelRatioF();
    painter.drawPixmap(p, m_pixmap);

    if (!m_overlayPixmap.isNull()) {
        painter.drawPixmap(p, m_overlayPixmap);
    }

    painter.end();
}

QPixmap SNITrayWidget::newIconPixmap(IconType iconType)
{
    QPixmap pixmap;
    if (iconType == UnknownIconType) {
        return pixmap;
    }

    QString iconName;
    DBusImageList dbusImageList;

    QString iconThemePath = m_sniInter->iconThemePath();

    switch (iconType) {
        case Icon:
            iconName = m_sniInter->iconName();
            dbusImageList = m_sniInter->iconPixmap();
            break;
        case OverlayIcon:
            iconName = m_sniInter->overlayIconName();
            dbusImageList = m_sniInter->overlayIconPixmap();
            break;
        case AttentionIcon:
            iconName = m_sniInter->attentionIconName();
            dbusImageList = m_sniInter->attentionIconPixmap();
            break;
        case AttentionMovieIcon:
            iconName = m_sniInter->attentionMovieName();
            break;
        default:
            break;
    }

    const auto ratio = qApp->devicePixelRatio();
    const int iconSizeScaled = IconSize * ratio;
    do {
        // load icon from sni dbus
        if (!dbusImageList.isEmpty()) {
            for (DBusImage dbusImage : dbusImageList) {
                char *image_data = dbusImage.pixels.data();

                if (QSysInfo::ByteOrder == QSysInfo::LittleEndian) {
                    for (int i = 0; i < dbusImage.pixels.size(); i += 4) {
                        *(qint32*)(image_data + i) = qFromBigEndian(*(qint32*)(image_data + i));
                    }
                }

                QImage image((const uchar*)dbusImage.pixels.constData(), dbusImage.width, dbusImage.height, QImage::Format_ARGB32);
                pixmap = QPixmap::fromImage(image.scaled(iconSizeScaled, iconSizeScaled, Qt::KeepAspectRatio, Qt::SmoothTransformation));
                pixmap.setDevicePixelRatio(ratio);
                if (!pixmap.isNull()) {
                    break;
                }
            }
        }

        // load icon from specified file
        if (!iconThemePath.isEmpty() && !iconName.isEmpty()) {
            QDirIterator it(iconThemePath, QDirIterator::Subdirectories);
            while (it.hasNext()) {
                it.next();
                if (it.fileName().startsWith(iconName, Qt::CaseInsensitive)) {
                    QImage image(it.filePath());
                    pixmap = QPixmap::fromImage(image.scaled(iconSizeScaled, iconSizeScaled, Qt::KeepAspectRatio, Qt::SmoothTransformation));
                    pixmap.setDevicePixelRatio(ratio);
                    if (!pixmap.isNull()) {
                        break;
                    }
                }
            }
            if (!pixmap.isNull()) {
                break;
            }
        }

        // load icon from theme
        // Note: this will ensure return a None-Null pixmap
        // so, it should be the last fallback
        if (!iconName.isEmpty()) {
            // ThemeAppIcon::getIcon 会处理高分屏缩放问题
            pixmap = ThemeAppIcon::getIcon(iconName, IconSize);
            if (!pixmap.isNull()) {
                break;
            }
        }

        if (pixmap.isNull()) {
            qDebug() << "get icon faild!";
        }
    } while (false);

    return pixmap;
}
