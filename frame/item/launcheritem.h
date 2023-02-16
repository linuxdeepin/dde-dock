// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef LAUNCHERITEM_H
#define LAUNCHERITEM_H

#include "dockitem.h"

namespace Dock {
class TipsWidget;
}

class QGSettings;
class LauncherItem : public DockItem
{
    Q_OBJECT

public:
    explicit LauncherItem(QWidget *parent = nullptr);

    inline ItemType itemType() const override {return Launcher;}

    void refreshIcon() override;

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
    const QGSettings *m_gsettings;
    QSharedPointer<TipsWidget> m_tips;
};

#endif // LAUNCHERITEM_H
