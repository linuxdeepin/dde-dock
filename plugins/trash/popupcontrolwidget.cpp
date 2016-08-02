#include "popupcontrolwidget.h"

#include <QVBoxLayout>
#include <QProcess>
#include <QDebug>
#include <QDir>

DWIDGET_USE_NAMESPACE

const QString TrashDir = QDir::homePath() + "/.local/share/Trash/files";

PopupControlWidget::PopupControlWidget(QWidget *parent)
    : QWidget(parent),

      m_openBtn(new DLinkButton(tr("Open"), this)),
      m_clearBtn(new DLinkButton(tr("Clear"), this)),

      m_fsWatcher(new QFileSystemWatcher(this))
{
    m_fsWatcher->addPath(TrashDir);

    QVBoxLayout *centeralLayout = new QVBoxLayout;
    centeralLayout->addWidget(m_openBtn);
    centeralLayout->addWidget(m_clearBtn);

    connect(m_openBtn, &DLinkButton::clicked, this, &PopupControlWidget::openTrashFloder);
    connect(m_clearBtn, &DLinkButton::clicked, this, &PopupControlWidget::clearTrashFloder);
    connect(m_fsWatcher, &QFileSystemWatcher::directoryChanged, this, &PopupControlWidget::trashStatusChanged, Qt::QueuedConnection);

    setLayout(centeralLayout);
    setFixedWidth(80);
    setFixedHeight(60);
}

bool PopupControlWidget::empty() const
{
    return m_empty;
}

void PopupControlWidget::openTrashFloder()
{
    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    proc->startDetached("gvfs-open trash://");
}

void PopupControlWidget::clearTrashFloder()
{
    for (auto item : QDir(TrashDir).entryInfoList())
    {
        if (item.isFile())
            QFile(item.fileName()).remove();
        else if (item.isDir())
            QDir(item.path()).removeRecursively();
    }
}

void PopupControlWidget::trashStatusChanged()
{
    const bool empty = QDir(TrashDir).entryList().count() == 2;

    if (m_empty == empty)
        return;

    m_empty = empty;
    emit emptyChanged(m_empty);
}
