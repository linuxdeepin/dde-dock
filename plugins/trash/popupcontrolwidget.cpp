/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "popupcontrolwidget.h"

#include <QVBoxLayout>
#include <QProcess>
#include <QDebug>
#include <QDir>

#include <ddialog.h>

#include <com_deepin_daemon_soundeffect.h>

DWIDGET_USE_NAMESPACE

const QString TrashDir = QDir::homePath() + "/.local/share/Trash";

using SoundEffectInter = com::deepin::daemon::SoundEffect;

PopupControlWidget::PopupControlWidget(QWidget *parent)
    : QWidget(parent),

      m_empty(false),

      m_fsWatcher(new QFileSystemWatcher(this))
{
    connect(m_fsWatcher, &QFileSystemWatcher::directoryChanged, this, &PopupControlWidget::trashStatusChanged, Qt::QueuedConnection);

    setObjectName("trash");
    setFixedWidth(80);

    trashStatusChanged();
}

bool PopupControlWidget::empty() const
{
    return m_empty;
}

int PopupControlWidget::trashItems() const
{
    return m_trashItemsCount;
}

QSize PopupControlWidget::sizeHint() const
{
    return QSize(width(), m_empty ? 30 : 60);
}

const QString PopupControlWidget::trashDir()
{
    return TrashDir;
}

void PopupControlWidget::openTrashFloder()
{
    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    proc->startDetached("gvfs-open trash:///");
}

void PopupControlWidget::clearTrashFloder()
{
    // show confrim dialog
    bool accept = false;
    const int itemCount = trashItemCount();
    const QStringList btns = {tr("Cancel"), tr("Empty")};

    DDialog *dialog = new DDialog(nullptr);
    dialog->addButtons(btns);
    dialog->setIconPixmap(QIcon::fromTheme("user-trash-full").pixmap(48, 48));
    dialog->setMessage(tr("This action cannot be restored"));
    if (itemCount == 1)
        dialog->setTitle(tr("Are you sure to empty 1 item ?"));
    else
        dialog->setTitle(tr("Are you sure to empty %1 items ?").arg(itemCount));

    connect(dialog, &DDialog::buttonClicked, [&] (const int index) {
        accept = index;
    });
    dialog->exec();
    dialog->deleteLater();

    if (!accept)
        return;

    QProcess::startDetached("gvfs-trash", QStringList() << "-f" << "--empty");

    // play sound effects
    SoundEffectInter sei("com.deepin.daemon.SoundEffect", "/com/deepin/daemon/SoundEffect", QDBusConnection::sessionBus());
    sei.PlaySystemSound("trash-empty");

//    for (auto item : QDir(TrashDir).entryInfoList())
//    {
//        if (item.fileName() == "." || item.fileName() == "..")
//            continue;

//        if (item.isFile())
//            QFile(item.fileName()).remove();
//        else if (item.isDir())
//            QDir(item.absoluteFilePath()).removeRecursively();
//    }
}

int PopupControlWidget::trashItemCount() const
{
    return QDir(TrashDir + "/info").entryInfoList().count() - 2;
}

void PopupControlWidget::trashStatusChanged()
{
    const bool files = QDir(TrashDir + "/files").exists();
//    const bool info = QDir(TrashDir + "/info").exists();

    // add monitor paths
    m_fsWatcher->addPath(TrashDir);
    if (files)
        m_fsWatcher->addPath(TrashDir + "/files");
//    if (info)
//        m_fsWatcher->addPath(TrashDir + "/info");

    // check empty
    if (!files)
        m_trashItemsCount = 0;
    else
        m_trashItemsCount = QDir(TrashDir + "/files").entryList().count() - 2;

    const bool empty = m_trashItemsCount == 0;
    if (m_empty == empty)
        return;

//    m_clearBtn->setVisible(!empty);
    m_empty = empty;

    setFixedHeight(sizeHint().height());

    emit emptyChanged(m_empty);
}
