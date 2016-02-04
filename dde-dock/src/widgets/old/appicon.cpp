/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QFile>
#include <QPainter>
#include <QSvgRenderer>

#include "appicon.h"
#include "controller/signalmanager.h"

#undef signals
extern "C" {
  #include <gtk/gtk.h>
}
#define signals public

AppIcon::AppIcon(QWidget *parent, Qt::WindowFlags f) :
    QLabel(parent, f)
{
    // as far as I know, it's safe to call this method multiple times,
    // because it does some check work internally.
    gtk_init(NULL, NULL);
    gdk_error_trap_push();

    this->setAttribute(Qt::WA_TranslucentBackground);
    this->setAlignment(Qt::AlignCenter);

    connect(SignalManager::instance(), &SignalManager::requestAppIconUpdate, this, &AppIcon::updateIcon);
}

void AppIcon::setIcon(const QString &iconPath)
{
    m_iconPath = iconPath;

    QPixmap pixmap(48, 48);

    // iconPath is an absolute path of the system.
    if (QFile::exists(iconPath)) {
        pixmap = QPixmap(iconPath);
    } else if (iconPath.startsWith("data:image/")){
        // iconPath is a string representing an inline image.
        QStringList strs = iconPath.split("base64,");
        if (strs.length() == 2) {
            QByteArray data = QByteArray::fromBase64(strs.at(1).toLatin1());
            pixmap.loadFromData(data);
        }
    } else {
        // try to read the iconPath as a icon name.
        QString path = getThemeIconPath(iconPath);
        if (path.isEmpty())
            path = getThemeIconPath("application-x-desktop");
        if (path.endsWith(".svg")) {
            QSvgRenderer renderer(path);
            pixmap.fill(Qt::transparent);

            QPainter painter;
            painter.begin(&pixmap);

            renderer.render(&painter);

            painter.end();
        } else {
            pixmap.load(path);
        }
    }

    if (!pixmap.isNull()) {
        pixmap = pixmap.scaled(m_modeData->getAppIconSize(),
                               m_modeData->getAppIconSize(),
                               Qt::KeepAspectRatioByExpanding,
                               Qt::SmoothTransformation);

        setPixmap(pixmap);
    }
}

void AppIcon::mousePressEvent(QMouseEvent *ev)
{
    if (m_modeData->getDockMode() == Dock::FashionMode)
        emit mousePress(ev);
    else
        QLabel::mousePressEvent(ev);
}

void AppIcon::mouseReleaseEvent(QMouseEvent *ev)
{
    if (m_modeData->getDockMode() == Dock::FashionMode)
        emit mouseRelease(ev);
    else
        QLabel::mousePressEvent(ev);
}

void AppIcon::enterEvent(QEvent *)
{
    if (m_modeData->getDockMode() == Dock::FashionMode)
        emit mouseEnter();
}

void AppIcon::leaveEvent(QEvent *)
{
    if (m_modeData->getDockMode() == Dock::FashionMode)
        emit mouseLeave();
}

// iconName should be a icon name constraints to the freeedesktop standard.
QString AppIcon::getThemeIconPath(QString iconName)
{
    QByteArray bytes = iconName.toUtf8();
    const char *name = bytes.constData();

    GtkIconTheme* theme = gtk_icon_theme_get_default();

    GtkIconInfo* info = gtk_icon_theme_lookup_icon(theme, name, 48, GTK_ICON_LOOKUP_GENERIC_FALLBACK);

    if (info) {
        char* path = g_strdup(gtk_icon_info_get_filename(info));
#if GTK_MAJOR_VERSION >= 3
        g_object_unref(info);
#elif GTK_MAJOR_VERSION == 2
        gtk_icon_info_free(info);
#endif
        return QString(path);
    } else {
        return "";
    }
}

void AppIcon::updateIcon()
{
    if (!m_iconPath.isEmpty())
        setIcon(m_iconPath);
}
