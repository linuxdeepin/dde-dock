// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "popupcontrolwidget.h"

#include <QVBoxLayout>
#include <QProcess>
#include <QDebug>
#include <QDir>

#include <ddialog.h>
#include <DTrashManager>
#include <DDesktopServices>

#include <QCoreApplication>

DWIDGET_USE_NAMESPACE
DCORE_USE_NAMESPACE

PopupControlWidget::PopupControlWidget(QWidget *parent)
    : QWidget(parent),
      m_empty(false),
      m_trashHelper(new TrashHelper(this))
{
    connect(m_trashHelper, &TrashHelper::trashAttributeChanged, this, &PopupControlWidget::trashStatusChanged, Qt::QueuedConnection);

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

void PopupControlWidget::openTrashFloder()
{
    DDesktopServices::showFolder(QUrl("trash:///"));
}

void PopupControlWidget::clearTrashFloder()
{
    QString ClearTrashMutliple = qApp->translate("DialogManager", "Are you sure you want to empty %1 items?");

    // show confirm dialog
    DDialog d;
    QStringList buttonTexts;
    buttonTexts << qApp->translate("DialogManager", "Cancel") << qApp->translate("DialogManager", "Delete");

    if (!d.parentWidget()) {
        d.setWindowFlags(d.windowFlags() | Qt::WindowStaysOnTopHint);
    }

    uint count =  m_trashHelper->trashItemCount();
    int execCode = -1;

    if (count > 0) {
        // blumia: Workaround. There is a bug with DDialog which will let DDialog always use the smallest
        //         available size of the given icon. So we create a dialogTrashFullIcon and leave a minimum
        //         64*64 pixmap size icon here.
        QIcon dialogTrashFullIcon;
        QIcon trash_full_icon = QIcon::fromTheme("user-trash-full-opened");
        dialogTrashFullIcon.addPixmap(trash_full_icon.pixmap(64));
        dialogTrashFullIcon.addPixmap(trash_full_icon.pixmap(128));

        d.setTitle(ClearTrashMutliple.arg(count));
        d.setMessage(qApp->translate("DialogManager", "This action cannot be restored"));
        d.setIcon(dialogTrashFullIcon);
        d.addButton(buttonTexts[0], true, DDialog::ButtonNormal);
        d.addButton(buttonTexts[1], false, DDialog::ButtonWarning);
        d.setDefaultButton(1);
        d.moveToCenter();
        execCode = d.exec();
    }

    if (execCode != QDialog::Accepted) {
        return;
    }

    if (m_trashHelper->emptyTrash()) {
        DDesktopServices::playSystemSoundEffect(DDesktopServices::SSE_EmptyTrash);
    } else {
        qDebug() << "Clear trash failed";
    }
}

int PopupControlWidget::trashItemCount() const
{
    return m_trashHelper->trashItemCount();
}

void PopupControlWidget::trashStatusChanged()
{
    m_trashItemsCount = m_trashHelper->trashItemCount();
    const bool empty = m_trashItemsCount == 0;
    if (m_empty == empty) {
        return;
    }
    m_empty = empty;

    setFixedHeight(sizeHint().height());

    emit emptyChanged(m_empty);
}
