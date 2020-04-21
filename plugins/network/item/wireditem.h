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

#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "deviceitem.h"

#include <QWidget>
#include <QLabel>
#include <QTimer>

#include <WiredDevice>
#include <QIcon>

class TipsWidget;
class WiredItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WiredItem(dde::network::WiredDevice *device);

    QWidget *itemTips() override;
    const QString itemCommand() const override;

protected:
    void paintEvent(QPaintEvent *e) override;
    void resizeEvent(QResizeEvent *e) override;

private slots:
    void refreshIcon() override;
    void reloadIcon();
    void deviceStateChanged();
    void refreshTips();

private:
    QPixmap m_icon;

    TipsWidget *m_itemTips;
    QTimer *m_delayTimer;
};

#endif // WIREDITEM_H
