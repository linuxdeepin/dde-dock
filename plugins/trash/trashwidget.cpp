// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "constants.h"
#include "trashwidget.h"

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
        // accept prevent the event from being propgated to the dock main panel
        // which also takes drag event;

        if (!e->mimeData()->hasFormat("Removable")) {
            // show the forbit dropping cursor.
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
        return removeApp(e->mimeData()->data("AppKey"));

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
    // 这里需要取到所有url然后一次性移动到回收站。
    moveToTrash(mime->urls());
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
    if (size > PLUGIN_BACKGROUND_MAX_SIZE) {
        size *= ((Dock::Fashion == qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>()) ? 0.8 : 0.7);
        if (size < PLUGIN_BACKGROUND_MAX_SIZE)
            size = PLUGIN_BACKGROUND_MAX_SIZE;
    }


    QIcon icon = QIcon::fromTheme(iconString, m_defaulticon);

    const auto ratio = devicePixelRatioF();
    int pixmapSize = QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? size : int(size * ratio);
    m_icon = icon.pixmap(pixmapSize, pixmapSize);
    m_icon.setDevicePixelRatio(ratio);
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
    const QString cmd("dbus-send --print-reply --dest=com.deepin.dde.Launcher /com/deepin/dde/Launcher com.deepin.dde.Launcher.UninstallApp string:\"" + appKey + "\"");

    QProcess *proc = new QProcess;
    proc->start(cmd);
    proc->waitForFinished();

    proc->deleteLater();
}

void TrashWidget::moveToTrash(const QList<QUrl> &urlList)
{
    QStringList argumentList;
    for (const QUrl &url : urlList) {
        const QFileInfo& info = url.toLocalFile();
        argumentList << info.absoluteFilePath();
    }

    m_fileManagerInter->Trash(argumentList);
}
