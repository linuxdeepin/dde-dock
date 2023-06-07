// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "appitem.h"
#include "appmultiitem.h"
#include "imageutil.h"
#include "screenspliter.h"
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
    connect(m_appItem, &AppItem::onCurrentWindowChanged, this, &AppMultiItem::onCurrentWindowChanged);
}

void AppMultiItem::onOpen()
{
    m_appItem->activeWindow(m_winId);
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

    if (m_pixmap.isNull())
        m_pixmap = ImageUtil::loadWindowThumb(Utils::IS_WAYLAND_DISPLAY ? m_windowInfo.uuid : QString::number(m_winId));

    DStyleHelper dstyle(style());
    const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);
    QRect itemRect = rect();
    itemRect.marginsRemoved(QMargins(6, 6, 6, 6));
    QPainterPath path;
    path.addRoundedRect(rect(), radius, radius);
    painter.fillPath(path, Qt::transparent);

    if (m_appItem->currentWindow() == m_winId) {
        QColor backColor = Qt::black;
        backColor.setAlpha(255 * 0.8);
        painter.fillPath(path, backColor);
    }

    itemRect = m_pixmap.rect();
    int itemWidth = itemRect.width();
    int itemHeight = itemRect.height();
    int x = (rect().width() - itemWidth) / 2;
    int y = (rect().height() - itemHeight) / 2;
    painter.drawPixmap(QRect(x, y, itemWidth, itemHeight), m_pixmap);

    QPixmap pixmapAppIcon;
    ThemeAppIcon::getIcon(pixmapAppIcon, m_appItem->appId(), qMin(width(), height()) * 0.8);
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
        m_appItem->activeWindow(m_winId);
    } else {
        QPoint currentPoint = QCursor::pos();
        m_menu->exec(currentPoint);
    }
}
