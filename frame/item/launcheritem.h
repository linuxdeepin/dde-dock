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
    explicit LauncherItem(QWidget *parent = nullptr) ;

    inline ItemType itemType() const override {return Launcher;}

    void refershIcon() override;

protected:
    void showEvent(QShowEvent* event) override;

private:
    void paintEvent(QPaintEvent *e) override;
    void resizeEvent(QResizeEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;

    QWidget *popupTips() override;

    void onGSettingsChanged(const QString& key);

    bool checkGSettingsControl() const;

private:
    QPixmap m_icon;
    LauncherInter *m_launcherInter;
    TipsWidget *m_tips;
    QGSettings* m_gsettings;
};

#endif // LAUNCHERITEM_H
