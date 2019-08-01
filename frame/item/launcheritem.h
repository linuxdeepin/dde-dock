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

#ifndef LAUNCHERITEM_H
#define LAUNCHERITEM_H

#include "dockitem.h"
#include "../widgets/tipswidget.h"

#include <com_deepin_dde_launcher.h>


using LauncherInter = com::deepin::dde::Launcher;

class QGSettings;
class LauncherItem : public DockItem
{
    Q_OBJECT

public:
    explicit LauncherItem(QWidget *parent = 0);

    inline ItemType itemType() const {return Launcher;}

    void refershIcon();

protected:
    void showEvent(QShowEvent* event) override;

private:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

    QWidget *popupTips();

    void onGSettingsChanged(const QString& key);

    bool checkGSettingsControl() const;

private:
    QPixmap m_smallIcon;
    QPixmap m_largeIcon;
    LauncherInter *m_launcherInter;
    TipsWidget *m_tips;
    QGSettings* m_gsettings;
};

#endif // LAUNCHERITEM_H
