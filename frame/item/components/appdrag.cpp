// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "appdrag.h"
#include "utils.h"

#include <QGSettings>
#include <QDebug>

#include <X11/Xcursor/Xcursor.h>

AppDrag::AppDrag(QObject *dragSource)
    : QDrag(dragSource)
    , m_appDragWidget(new AppDragWidget)
{
    // delete by itself
    m_appDragWidget->setVisible(false);

    if (Utils::IS_WAYLAND_DISPLAY)
        setDragMoveCursor();
}

AppDrag::~AppDrag() {
    // delete when AppDragWidget is invisible
    if (m_appDragWidget) {
        m_appDragWidget->deleteLater();
    }
}

void AppDrag::setPixmap(const QPixmap &pix)
{
    m_appDragWidget->setAppPixmap(pix);
}

QPixmap AppDrag::pixmap() const
{
    /* TODO: return pixmap */
    return QPixmap();
}

Qt::DropAction AppDrag::start(Qt::DropActions supportedActions)
{
    m_appDragWidget->show();
    return QDrag::start(supportedActions);
}

Qt::DropAction AppDrag::exec(Qt::DropActions supportedActions)
{
    m_appDragWidget->show();
    return QDrag::exec(supportedActions);
}

Qt::DropAction AppDrag::exec(Qt::DropActions supportedActions, Qt::DropAction defaultAction)
{
    m_appDragWidget->show();
    return QDrag::exec(supportedActions, defaultAction);
}

AppDragWidget *AppDrag::appDragWidget()
{
    return m_appDragWidget;
}

void AppDrag::setDragMoveCursor()
{
    QString theme = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-name", "bloom").toString();
    int cursorSize = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-size", 24).toInt();
    const char* cursorName = "dnd-move";
    XcursorImages *images = XcursorLibraryLoadImages(cursorName, theme.toStdString().c_str(), cursorSize);
    if (!images || !(images->images[0])) {
        qWarning() << "loadCursorFalied, theme =" << theme << ", cursorName=" << cursorName;
        return;
    }
    const int imgW = images->images[0]->width;
    const int imgH = images->images[0]->height;

    QImage img((const uchar*)images->images[0]->pixels, imgW, imgH, QImage::Format_ARGB32);
    QPixmap pixmap = QPixmap::fromImage(img);
    XcursorImagesDestroy(images);
    setDragCursor(pixmap, Qt::MoveAction);
}
