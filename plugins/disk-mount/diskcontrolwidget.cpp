#include "diskcontrolwidget.h"
#include "diskcontrolitem.h"

#define MAX_HEIGHT      300
#define WIDTH           300

DiskControlWidget::DiskControlWidget(QWidget *parent)
    : QScrollArea(parent),

      m_centeralLayout(new QVBoxLayout),
      m_centeralWidget(new QWidget),

      m_diskInter(new DBusDiskMount(this))
{
    m_centeralWidget->setLayout(m_centeralLayout);
    m_centeralWidget->setFixedWidth(WIDTH);

    setWidget(m_centeralWidget);
    setFixedWidth(WIDTH);
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

    while (QLayoutItem *item = m_centeralLayout->takeAt(0))
    {
        delete item->widget();
        delete item;
    }

    int mountedCount = 0;
    for (auto info : m_diskInfoList)
    {
        if (info.m_mountPoint.isEmpty())
            continue;
        else
            ++mountedCount;

        DiskControlItem *item = new DiskControlItem(info, this);

        connect(item, &DiskControlItem::requestUnmount, this, &DiskControlWidget::unmountDisk);

        m_centeralLayout->addWidget(item);
    }

    emit diskCountChanged(mountedCount);

    const int contentHeight = mountedCount * 70;
    const int maxHeight = std::min(contentHeight, MAX_HEIGHT);

    m_centeralWidget->setFixedHeight(contentHeight);
    setFixedHeight(maxHeight);
}

void DiskControlWidget::unmountDisk(const QString &diskId) const
{
    m_diskInter->Unmount(diskId);
}
