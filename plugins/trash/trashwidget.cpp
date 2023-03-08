// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "constants.h"
#include "trashwidget.h"
#include "imageutil.h"

#include <DDBusSender>

#include <QPainter>
#include <QIcon>
#include <QApplication>
#include <QDragEnterEvent>
#include <QJsonDocument>
#include <QApplication>
#include <QDBusConnection>

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent)
    , m_popupApplet(new PopupControlWidget(this))
    , m_fileManagerInter(new DBusFileManager1("org.freedesktop.FileManager1",
                                              "/org/freedesktop/FileManager1",
                                              QDBusConnection::sessionBus(),
                                              this))
    , m_dragging(false)
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
    if (e->mimeData()->hasFormat("RequestDock")) {
        // accept prevent the event from being propagated to the dock main panel
        // which also takes drag event;

        if (!e->mimeData()->hasFormat("Removable")) {
            // show the forbid dropping cursor.
            e->setDropAction(Qt::IgnoreAction);
        } else {
            e->setDropAction(Qt::MoveAction);
            e->accept();
        }

        return;
    }

    if (!e->mimeData()->hasUrls())
        return e->ignore();

    e->setDropAction(Qt::MoveAction);

    if (e->dropAction() != Qt::MoveAction) {
        e->ignore();
    } else {
        // 设置item是否拖入回收站的状态，给DockItem发送鼠标进入事件
        setDragging(true);
        qApp->postEvent(this->parent(), new QEnterEvent(e->pos(), mapToParent(e->pos()), mapToGlobal(e->pos())));
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

void TrashWidget::dragLeaveEvent(QDragLeaveEvent *e)
{
    Q_UNUSED(e);

    // 设置item是否拖入回收站的状态，给DockItem发送鼠标离开事件
    setDragging(false);
    qApp->postEvent(this->parent(), new QEvent(QEvent::Leave));
}

void TrashWidget::dropEvent(QDropEvent *e)
{
    if (e->mimeData()->hasFormat("RequestDock"))
        return removeApp(e->mimeData()->data("DesktopPath"));

    if (!e->mimeData()->hasUrls()) {
        return e->ignore();
    }

    e->setDropAction(Qt::MoveAction);

    if (e->dropAction() != Qt::MoveAction) {
        return e->ignore();
    }

    // 设置item是否拖入回收站的状态，给DockItem发送鼠标离开事件
    setDragging(false);
    qApp->postEvent(this->parent(), new QEvent(QEvent::Leave));

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
    Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QString iconString = "user-trash";
    if (!m_popupApplet->empty())
        iconString.append("-full");

    int size = qMin(width(), height());
    size = qMax(PLUGIN_ICON_MIN_SIZE, static_cast<int>(size * ((Dock::Efficient == displayMode) ? 0.7 : 0.8)));

    m_icon = ImageUtil::loadSvg(iconString, QSize(size, size), devicePixelRatioF());
}

void TrashWidget::updateIconAndRefresh()
{
    updateIcon();
    update();
}

bool TrashWidget::getDragging() const
{
    return m_dragging;
}

void TrashWidget::setDragging(bool state)
{
    m_dragging = state;
}

void TrashWidget::removeApp(const QString &appKey)
{
    DDBusSender().service("org.deepin.dde.Launcher1")
            .path("/org/deepin/dde/Launcher1")
            .interface("org.deepin.dde.Launcher1")
            .method("UninstallApp")
            .arg(appKey)
            .call();
}

void TrashWidget::moveToTrash(const QUrl &url)
{
    const QFileInfo info = url.toLocalFile();

    QStringList argumentList;
    argumentList << info.absoluteFilePath();

    m_fileManagerInter->Trash(argumentList);
}
