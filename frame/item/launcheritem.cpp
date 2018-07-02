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

DCORE_USE_NAMESPACE

LauncherItem::LauncherItem(QWidget *parent)
    : DockItem(parent),

      m_tips(new TipsWidget(this))
{
    setAccessibleName("Launcher");
    m_tips->setVisible(false);
    m_tips->setObjectName("launcher");
    m_tips->setText(tr("Launcher"));
}

void LauncherItem::refershIcon()
{
    const int iconSize = qMin(width(), height());
    if (DockDisplayMode == Efficient)
    {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.7);
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.9);
    } else {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.6);
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.8);
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

    const auto ratio = qApp->devicePixelRatio();
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

    if (e->button() == Qt::RightButton/* && !perfectIconRect().contains(e->pos())*/)
        return QWidget::mousePressEvent(e);

    if (e->button() != Qt::LeftButton)
        return;

    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    DDBusSender()
            .service("com.deepin.dde.Launcher")
            .interface("com.deepin.dde.Launcher")
            .path("/com/deepin/dde/Launcher")
            .method("Toggle")
            .call();
}

QWidget *LauncherItem::popupTips()
{
    return m_tips;
}
