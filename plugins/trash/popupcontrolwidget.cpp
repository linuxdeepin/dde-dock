#include "popupcontrolwidget.h"

#include <QVBoxLayout>
#include <QProcess>
#include <QDebug>
#include <QDir>

#include <ddialog.h>

DWIDGET_USE_NAMESPACE

const QString TrashDir = QDir::homePath() + "/.local/share/Trash";

PopupControlWidget::PopupControlWidget(QWidget *parent)
    : QWidget(parent),

      m_empty(false),

//      m_openBtn(new DLinkButton(tr("Run"), this)),
//      m_clearBtn(new DLinkButton(tr("Empty Trash"), this)),

      m_fsWatcher(new QFileSystemWatcher(this))
{
//    QVBoxLayout *centeralLayout = new QVBoxLayout;
//    centeralLayout->addWidget(m_openBtn);
//    centeralLayout->addWidget(m_clearBtn);

//    connect(m_openBtn, &DLinkButton::clicked, this, &PopupControlWidget::openTrashFloder);
//    connect(m_clearBtn, &DLinkButton::clicked, this, &PopupControlWidget::clearTrashFloder);
    connect(m_fsWatcher, &QFileSystemWatcher::directoryChanged, this, &PopupControlWidget::trashStatusChanged, Qt::QueuedConnection);

//    setLayout(centeralLayout);
    setObjectName("trash");
    setFixedWidth(80);

    trashStatusChanged();
}

bool PopupControlWidget::empty() const
{
    return m_empty;
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

    proc->startDetached("gvfs-open trash://");
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

    for (auto item : QDir(TrashDir).entryInfoList())
    {
        if (item.fileName() == "." || item.fileName() == "..")
            continue;

        if (item.isFile())
            QFile(item.fileName()).remove();
        else if (item.isDir())
            QDir(item.absoluteFilePath()).removeRecursively();
    }
}

int PopupControlWidget::trashItemCount() const
{
    return QDir(TrashDir + "/info").entryInfoList().count() - 2;
}

void PopupControlWidget::trashStatusChanged()
{
    const bool files = QDir(TrashDir + "/files").exists();
    const bool info = QDir(TrashDir + "/info").exists();

    // add monitor paths
    m_fsWatcher->addPath(TrashDir);
    if (files)
        m_fsWatcher->addPath(TrashDir + "/files");
    if (info)
        m_fsWatcher->addPath(TrashDir + "/info");

    // check empty
    bool empty;
    if ((!info || QDir(TrashDir + "/info").entryList().count() == 2) &&
        (!files || QDir(TrashDir + "/files").entryList().count() == 2))
        empty = true;
    else
        empty = false;

    if (m_empty == empty)
        return;

//    m_clearBtn->setVisible(!empty);
    m_empty = empty;

    setFixedHeight(sizeHint().height());

    emit emptyChanged(m_empty);
}
