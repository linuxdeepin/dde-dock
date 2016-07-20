#ifndef DISKCONTROLITEM_H
#define DISKCONTROLITEM_H

#include "dbus/dbusdiskmount.h"

#include <dimagebutton.h>

#include <QWidget>
#include <QLabel>
#include <QProgressBar>

class DiskControlItem : public QWidget
{
    Q_OBJECT

public:
    explicit DiskControlItem(const DiskInfo &info, QWidget *parent = 0);

private slots:
    void updateInfo(const DiskInfo &info);
    const QString formatDiskSize(const quint64 size) const;

private:
    DiskInfo m_info;

    QLabel *m_diskIcon;
    QLabel *m_diskName;
    QLabel *m_diskCapacity;
    QProgressBar *m_capacityValueBar;
    Dtk::Widget::DImageButton *m_unmountButton;
};

#endif // DISKCONTROLITEM_H
