#include "diskcontrolwidget.h"

DiskControlWidget::DiskControlWidget(QWidget *parent)
    : QScrollArea(parent),

      m_diskInter(new DBusDiskMount(this))
{
    setFixedWidth(300);
    setFrameStyle(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    connect(m_diskInter, &DBusDiskMount::DiskListChanged, this, &DiskControlWidget::diskListChanged);

    QMetaObject::invokeMethod(this, "diskListChanged", Qt::QueuedConnection);
}

void DiskControlWidget::diskListChanged()
{
    m_diskInfoList = m_diskInter->diskList();

    emit diskCountChanged(m_diskInfoList.count());
}
