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

    emit diskCountChanged(m_diskInfoList.count());

    for (auto info : m_diskInfoList)
    {
        DiskControlItem *item = new DiskControlItem(info, this);
        m_centeralLayout->addWidget(item);
    }

    const int contentHeight = m_diskInfoList.count() * 70;
    const int maxHeight = std::min(contentHeight, MAX_HEIGHT);

    m_centeralWidget->setFixedHeight(contentHeight);
    setFixedHeight(maxHeight);
}
