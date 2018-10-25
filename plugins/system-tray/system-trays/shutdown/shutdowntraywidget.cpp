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

#include "shutdowntraywidget.h"
#include "dbus/dbusaccount.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>

ShutdownTrayWidget::ShutdownTrayWidget(QWidget *parent)
    : AbstractSystemTrayWidget(parent),
      m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setText(tr("Shut down"));
    m_tipsLabel->setVisible(false);

    updateIcon();
}

void ShutdownTrayWidget::setActive(const bool active)
{

}

void ShutdownTrayWidget::updateIcon()
{
    const auto ratio = qApp->devicePixelRatio();

    QPixmap pixmap(QSize(16, 16) * ratio);
    QSvgRenderer renderer(QString(":/icons/system-trays/shutdown/resources/icons/normal.svg"));
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    pixmap.setDevicePixelRatio(ratio);

    m_pixmap = pixmap;

    update();
}

const QImage ShutdownTrayWidget::trayImage()
{
    return m_pixmap.toImage();
}

QWidget *ShutdownTrayWidget::trayTipsWidget()
{
    return m_tipsLabel;
}

const QString ShutdownTrayWidget::trayClickCommand()
{
    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

const QString ShutdownTrayWidget::contextMenu() const
{
    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> shutdown;
    shutdown["itemId"] = "Shutdown";
    shutdown["itemText"] = tr("Shut down");
    shutdown["isActive"] = true;
    items.push_back(shutdown);

    QMap<QString, QVariant> reboot;
    reboot["itemId"] = "Restart";
    reboot["itemText"] = tr("Restart");
    reboot["isActive"] = true;
    items.push_back(reboot);

    QMap<QString, QVariant> suspend;
    suspend["itemId"] = "Suspend";
    suspend["itemText"] = tr("Suspend");
    suspend["isActive"] = true;
    items.push_back(suspend);

    QMap<QString, QVariant> lock;
    lock["itemId"] = "Lock";
    lock["itemText"] = tr("Lock");
    lock["isActive"] = true;
    items.push_back(lock);

    QMap<QString, QVariant> logout;
    logout["itemId"] = "Logout";
    logout["itemText"] = tr("Log out");
    logout["isActive"] = true;
    items.push_back(logout);

    if (DBusAccount().userList().count() > 1)
    {
        QMap<QString, QVariant> switchUser;
        switchUser["itemId"] = "SwitchUser";
        switchUser["itemText"] = tr("Switch account");
        switchUser["isActive"] = true;
        items.push_back(switchUser);
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void ShutdownTrayWidget::invokedMenuItem(const QString &menuId, const bool checked)
{
    Q_UNUSED(checked)

    if (menuId == "Lock")
        QProcess::startDetached("dbus-send", QStringList() << "--print-reply"
                                                           << "--dest=com.deepin.dde.lockFront"
                                                           << "/com/deepin/dde/lockFront"
                                                           << QString("com.deepin.dde.lockFront.Show"));
    else
        QProcess::startDetached("dbus-send", QStringList() << "--print-reply"
                                                           << "--dest=com.deepin.dde.shutdownFront"
                                                           << "/com/deepin/dde/shutdownFront"
                                                           << QString("com.deepin.dde.shutdownFront.%1").arg(menuId));
}

QSize ShutdownTrayWidget::sizeHint() const
{
    return QSize(26, 26);
}

void ShutdownTrayWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_pixmap.rect().center() / qApp->devicePixelRatio(), m_pixmap);
}
