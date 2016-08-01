#include "popupcontrolwidget.h"

#include <QVBoxLayout>
#include <QProcess>
#include <QDebug>
#include <QMouseEvent>

DWIDGET_USE_NAMESPACE

PopupControlWidget::PopupControlWidget(QWidget *parent)
    : QWidget(parent),

      m_openBtn(new DLinkButton(tr("Open"), this)),
      m_clearBtn(new DLinkButton(tr("Clear"), this))
{
    QVBoxLayout *centeralLayout = new QVBoxLayout;
    centeralLayout->addWidget(m_openBtn);
    centeralLayout->addWidget(m_clearBtn);

    connect(m_openBtn, &DLinkButton::clicked, this, &PopupControlWidget::openTrashFloder);

    setLayout(centeralLayout);
    setFixedWidth(80);
    setFixedHeight(60);
}

void PopupControlWidget::openTrashFloder()
{
    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    proc->startDetached("gvfs-open trash://");
}
