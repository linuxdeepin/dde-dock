#include "diskcontrolitem.h"

#include <QVBoxLayout>
#include <QIcon>

DWIDGET_USE_NAMESPACE

DiskControlItem::DiskControlItem(const DiskInfo &info, QWidget *parent)
    : QWidget(parent),

      m_diskIcon(new QLabel),
      m_diskName(new QLabel),
      m_diskCapacity(new QLabel),
      m_capacityValueBar(new QProgressBar),
      m_unmountButton(new DImageButton)
{
    QIcon::setThemeName("deepin");

    m_diskName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_diskName->setStyleSheet("color:white;");

    m_diskCapacity->setStyleSheet("color:white;");

    m_capacityValueBar->setTextVisible(false);
    m_capacityValueBar->setFixedHeight(3);
    m_capacityValueBar->setStyleSheet("QProgressBar {"
                                      "border:none;"
                                      "background-color:rgba(255, 255, 255, .3);"
                                      "}"
                                      "QProgressBar::chunk {"
                                      "background-color:white;"
                                      "}");

    m_unmountButton->setNormalPic(":/icons/resources/unmount-normal.png");
    m_unmountButton->setHoverPic(":/icons/resources/unmount-hover.png");
    m_unmountButton->setPressPic(":/icons/resources/unmount-press.png");

    QVBoxLayout *infoLayout = new QVBoxLayout;
    infoLayout->addWidget(m_diskName);
    infoLayout->addWidget(m_diskCapacity);
    infoLayout->setSpacing(0);
    infoLayout->setMargin(0);

    QHBoxLayout *unmountLayout = new QHBoxLayout;
    unmountLayout->addLayout(infoLayout);
    unmountLayout->addWidget(m_unmountButton);
    unmountLayout->setSpacing(0);
    unmountLayout->setMargin(0);

    QVBoxLayout *progressLayout = new QVBoxLayout;
    progressLayout->addLayout(unmountLayout);
    progressLayout->addWidget(m_capacityValueBar);
    progressLayout->setSpacing(0);
    progressLayout->setMargin(0);

    QHBoxLayout *centeralLayout = new QHBoxLayout;
    centeralLayout->addWidget(m_diskIcon);
    centeralLayout->addLayout(progressLayout);
    centeralLayout->setSpacing(0);
    centeralLayout->setMargin(0);

    setLayout(centeralLayout);

    connect(m_unmountButton, &DImageButton::clicked, [this] {emit requestUnmount(m_info.m_id);});

    updateInfo(info);
}

void DiskControlItem::updateInfo(const DiskInfo &info)
{
    m_info = info;

    m_diskIcon->setPixmap(QIcon::fromTheme(info.m_icon).pixmap(32, 32));
    m_diskName->setText(info.m_name);
    m_diskCapacity->setText(QString("%1/%2").arg(formatDiskSize(info.m_usedSize)).arg(formatDiskSize(info.m_totalSize)));
    m_capacityValueBar->setMinimum(0);
    m_capacityValueBar->setMaximum(info.m_totalSize);
    m_capacityValueBar->setValue(info.m_usedSize);
}

const QString DiskControlItem::formatDiskSize(const quint64 size) const
{
    const quint64 mSize = 1000;
    const quint64 gSize = mSize * 1000;
    const quint64 tSize = gSize * 1000;

    if (size >= tSize)
        return QString::number(double(size) / tSize, 'f', 2) + 'T';
    else if (size >= gSize)
        return QString::number(double(size) / gSize, 'f', 2) + "G";
    else if (size >= mSize)
        return QString::number(double(size) / mSize, 'f', 1) + "M";
    else
        return QString::number(size) + "K";
}
