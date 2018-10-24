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

#ifndef WIREDTRAYWIDGET_H
#define WIREDTRAYWIDGET_H

#include "abstractnetworktraywidget.h"

#include <QWidget>
#include <QLabel>
#include <QTimer>

#include <WiredDevice>

class TipsWidget;
class WiredTrayWidget : public AbstractNetworkTrayWidget
{
    Q_OBJECT

public:
    explicit WiredTrayWidget(dde::network::WiredDevice *device, QWidget *parent = nullptr);

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *trayTipsWidget() Q_DECL_OVERRIDE;
    const QString trayClickCommand() Q_DECL_OVERRIDE;

protected:
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *e) Q_DECL_OVERRIDE;

private:
    QPixmap m_icon;

    TipsWidget *m_itemTips;
    QTimer *m_delayTimer;
};

#endif // WIREDTRAYWIDGET_H
