/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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

#include "diskcontrolitem.h"

#include <QVBoxLayout>
#include <QIcon>

DWIDGET_USE_NAMESPACE

DiskControlItem::DiskControlItem(const DiskInfo &info, QWidget *parent)
    : QFrame(parent),

      m_unknowIcon(":/icons/resources/unknown.svg"),

      m_diskIcon(new QLabel),
      m_diskName(new QLabel),
      m_diskCapacity(new QLabel),
      m_capacityValueBar(new QProgressBar),
      m_unmountButton(new DImageButton)
{
//    QIcon::setThemeName("deepin");

    m_diskName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_diskName->setStyleSheet("color:white;");

    m_diskCapacity->setStyleSheet("color:rgba(255, 255, 255, .6);");

    m_capacityValueBar->setTextVisible(false);
    m_capacityValueBar->setFixedHeight(2);
    m_capacityValueBar->setStyleSheet("QProgressBar {"
                                      "border:none;"
                                      "background-color:rgba(255, 255, 255, .1);"
                                      "}"
                                      "QProgressBar::chunk {"
                                      "background-color:rgba(255, 255, 255, .8);"
                                      "}");

    m_unmountButton->setNormalPic(":/icons/resources/unmount-normal.png");
    m_unmountButton->setHoverPic(":/icons/resources/unmount-hover.png");
    m_unmountButton->setPressPic(":/icons/resources/unmount-press.png");
    m_unmountButton->setStyleSheet("margin-top:12px;");

    QVBoxLayout *infoLayout = new QVBoxLayout;
    infoLayout->addWidget(m_diskName);
    infoLayout->addWidget(m_diskCapacity);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(3, 6, 0, 8);

    QHBoxLayout *unmountLayout = new QHBoxLayout;
    unmountLayout->addLayout(infoLayout);
    unmountLayout->addWidget(m_unmountButton);
    unmountLayout->setSpacing(0);
    unmountLayout->setMargin(0);

    QVBoxLayout *progressLayout = new QVBoxLayout;
    progressLayout->addLayout(unmountLayout);
    progressLayout->addWidget(m_capacityValueBar);
    progressLayout->setSpacing(0);
    progressLayout->setContentsMargins(10, 0, 0, 5);

    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->addWidget(m_diskIcon);
    centralLayout->addLayout(progressLayout);
    centralLayout->setSpacing(0);
    centralLayout->setContentsMargins(0, 0, 5, 0);

    setLayout(centralLayout);
    setObjectName("DiskItem");
    setStyleSheet("QFrame #DiskItem:hover {"
                  "background-color:rgba(255, 255, 255, .1);"
                  "border-radius:4px;"
                  "}");

    connect(m_unmountButton, &DImageButton::clicked, [this] {emit requestUnmount(m_info.m_id);});

    updateInfo(info);
}

void DiskControlItem::updateInfo(const DiskInfo &info)
{
    m_info = info;

    m_diskIcon->setPixmap(QIcon::fromTheme(info.m_icon, m_unknowIcon).pixmap(48, 48));
    if (!info.m_name.isEmpty())
        m_diskName->setText(info.m_name);
    else
        m_diskName->setText(tr("Unknown device"));
    if (info.m_totalSize)
        m_diskCapacity->setText(QString("%1/%2").arg(formatDiskSize(info.m_usedSize)).arg(formatDiskSize(info.m_totalSize)));
    else if (info.m_name.isEmpty())
        m_diskCapacity->clear();
    else
        m_diskCapacity->setText(tr("Unknown volume"));
    m_capacityValueBar->setMinimum(0);
    m_capacityValueBar->setMaximum(std::max(1ull, info.m_totalSize));
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
