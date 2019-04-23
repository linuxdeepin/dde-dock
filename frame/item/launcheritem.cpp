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

#include "launcheritem.h"
#include "util/themeappicon.h"
#include "util/imagefactory.h"

#include <QPainter>
#include <QProcess>
#include <QMouseEvent>
#include <DDBusSender>
#include <QApplication>

DCORE_USE_NAMESPACE

LauncherItem::LauncherItem(QWidget *parent)
    : DockItem(parent)
    , m_launcherInter(new LauncherInter("com.deepin.dde.Launcher", "/com/deepin/dde/Launcher", QDBusConnection::sessionBus(), this))
    , m_tips(new TipsWidget(this))
{
    m_launcherInter->setSync(true, false);

    setAccessibleName("Launcher");
    m_tips->setVisible(false);
    m_tips->setObjectName("launcher");
}

void LauncherItem::refershIcon()
{
    const int iconSize = qMin(width(), height());
    if (DockDisplayMode == Efficient)
    {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.7, devicePixelRatioF());
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.9, devicePixelRatioF());
    } else {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.6, devicePixelRatioF());
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.8, devicePixelRatioF());
    }

    update();
}

void LauncherItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (!isVisible())
        return;

    QPainter painter(this);

    const QPixmap pixmap = DockDisplayMode == Fashion ? m_largeIcon : m_smallIcon;

    const auto ratio = devicePixelRatioF();
    const int iconX = rect().center().x() - pixmap.rect().center().x() / ratio;
    const int iconY = rect().center().y() - pixmap.rect().center().y() / ratio;

    painter.drawPixmap(iconX, iconY, pixmap);
}

void LauncherItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    refershIcon();
}

void LauncherItem::mousePressEvent(QMouseEvent *e)
{
    hidePopup();

    return QWidget::mousePressEvent(e);
}

void LauncherItem::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() != Qt::LeftButton)
        return;

    if (!m_launcherInter->IsVisible()) {
        m_launcherInter->Show();
    }
}

QWidget *LauncherItem::popupTips()
{
    m_tips->setText(tr("Launcher"));
    return m_tips;
}
