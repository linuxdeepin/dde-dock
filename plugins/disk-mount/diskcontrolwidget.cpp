#include "diskcontrolwidget.h"

DiskControlWidget::DiskControlWidget(QWidget *parent)
    : QWidget(parent),

      m_diskInter(new DBusDiskMount(this))
{

    connect(m_diskInter, &DBusDiskMount::DiskListChanged, this, &DiskControlWidget::diskListChanged);

    QMetaObject::invokeMethod(this, "diskListChanged", Qt::QueuedConnection);
}

void DiskControlWidget::diskListChanged()
{
    m_diskInfoList = m_diskInter->diskList();

    emit diskCountChanged(m_diskInfoList.count());
}
