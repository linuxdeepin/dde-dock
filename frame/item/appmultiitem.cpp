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
#ifdef USE_AM
        if (Utils::IS_WAYLAND_DISPLAY)
            m_snapImage = ImageUtil::loadWindowThumb(m_windowInfo.uuid, width() - 20, height() - 20);
        else
#endif
            m_snapImage = ImageUtil::loadWindowThumb(m_winId, width() - 20, height() - 20);
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
