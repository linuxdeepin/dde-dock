/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *               2016 ~ 2018 dragondjf
 *
 * Author:     sbw <sbw@sbw.so>
 *             dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
 *             Tangtong<tangtong@deepin.com>
 *
 * Maintainer: dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
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

#include "constants.h"
#include "trashwidget.h"

#include <QPainter>
#include <QIcon>
#include <QApplication>
#include <QDragEnterEvent>
#include <QJsonDocument>
#include <QApplication>

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent)
    , m_popupApplet(new PopupControlWidget(this))
{
    m_popupApplet->setVisible(false);

    connect(m_popupApplet, &PopupControlWidget::emptyChanged, this, &TrashWidget::updateIconAndRefresh);

    setAcceptDrops(true);

    m_defaulticon = QIcon::fromTheme(":/icons/user-trash.svg");

    setMinimumSize(PLUGIN_ICON_MIN_SIZE, PLUGIN_ICON_MIN_SIZE);
}

QWidget *TrashWidget::popupApplet()
{
    return m_popupApplet;
}

const QString TrashWidget::contextMenu() const
{
    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> open;
    open["itemId"] = "open";
    open["itemText"] = tr("Open");
    open["isActive"] = true;
    items.push_back(open);

    if (!m_popupApplet->empty()) {
        QMap<QString, QVariant> empty;
        empty["itemId"] = "empty";
        empty["itemText"] = tr("Empty");
        empty["isActive"] = true;
        items.push_back(empty);
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

int TrashWidget::trashItemCount() const
{
    return m_popupApplet->trashItems();
}

void TrashWidget::invokeMenuItem(const QString &menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == "open")
        m_popupApplet->openTrashFloder();
    else if (menuId == "empty")
        m_popupApplet->clearTrashFloder();
}

void TrashWidget::dragEnterEvent(QDragEnterEvent *e)
{
    if (!e->mimeData()->hasUrls())
        return e->ignore();

    if (e->mimeData()->hasFormat("RequestDock")) {
        // accept prevent the event from being propgated to the dock main panel
        // which also takes drag event;
        e->accept();

        if (!e->mimeData()->hasFormat("Removable")) {
            // show the forbit dropping cursor.
            e->setDropAction(Qt::IgnoreAction);
        }

        return;
    }

    e->setDropAction(Qt::MoveAction);

    if (e->dropAction() != Qt::MoveAction) {
        e->ignore();
    } else {
        e->accept();
    }
}

void TrashWidget::dragMoveEvent(QDragMoveEvent *e)
{
    if (!e->mimeData()->hasUrls())
        return;

    e->setDropAction(Qt::MoveAction);

    if (e->dropAction() != Qt::MoveAction) {
        e->ignore();
    } else {
        e->accept();
    }
}

void TrashWidget::dropEvent(QDropEvent *e)
{
    if (e->mimeData()->hasFormat("RequestDock"))
        return removeApp(e->mimeData()->data("AppKey"));

    if (!e->mimeData()->hasUrls()) {
        return e->ignore();
    }

    e->setDropAction(Qt::MoveAction);

    if (e->dropAction() != Qt::MoveAction) {
        return e->ignore();
    }

    const QMimeData *mime = e->mimeData();
    for (auto url : mime->urls())
        moveToTrash(url);
}

void TrashWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    updateIcon();

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_icon.rect());
    painter.drawPixmap(rf.center() - rfp.center() / devicePixelRatioF(), m_icon);
}

void TrashWidget::updateIcon()
{
//    Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    Dock::DisplayMode displayMode = Dock::Fashion;

    QString iconString = "user-trash";
    if (!m_popupApplet->empty())
        iconString.append("-full");
    if (displayMode == Dock::Efficient)
        iconString.append("-symbolic");

    int size = std::min(width(), height());
    if (size < PLUGIN_ICON_MIN_SIZE)
        size = PLUGIN_ICON_MIN_SIZE;
    if (size > PLUGIN_BACKGROUND_MAX_SIZE)
    {
        size *= ((Dock::Fashion == qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>()) ? 0.8 : 0.7);
        if(size < PLUGIN_BACKGROUND_MAX_SIZE)
            size = PLUGIN_BACKGROUND_MAX_SIZE;
    }


    QIcon icon = QIcon::fromTheme(iconString, m_defaulticon);

    const auto ratio = devicePixelRatioF();
    m_icon = icon.pixmap(size * ratio, size * ratio);
    m_icon.setDevicePixelRatio(ratio);
}

void TrashWidget::updateIconAndRefresh()
{
    updateIcon();
    update();
}

void TrashWidget::removeApp(const QString &appKey)
{
    const QString cmd("dbus-send --print-reply --dest=com.deepin.dde.Launcher /com/deepin/dde/Launcher com.deepin.dde.Launcher.UninstallApp string:\"" + appKey + "\"");

    QProcess *proc = new QProcess;
    proc->start(cmd);
    proc->waitForFinished();

    proc->deleteLater();
}

void TrashWidget::moveToTrash(const QUrl &url)
{
    const QFileInfo info = url.toLocalFile();

    QProcess::startDetached("gio", QStringList() << "trash" << "-f" << info.absoluteFilePath());
}
