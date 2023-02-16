// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MULTIWINDOWHELPER_H
#define MULTIWINDOWHELPER_H

#include "constants.h"

#include <QObject>

class AppMultiItem;

class MultiWindowHelper : public QObject
{
    Q_OBJECT

public:
    explicit MultiWindowHelper(QWidget *appWidget, QWidget *multiWindowWidget, QObject *parent = nullptr);

    void setDisplayMode(Dock::DisplayMode displayMode);
    void addMultiWindow(int, AppMultiItem *item);
    void removeMultiWindow(AppMultiItem *item);

Q_SIGNALS:
    void requestUpdate();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    int itemIndex(AppMultiItem *item);
    void insertChildWidget(QWidget *parentWidget, int index, AppMultiItem *item);
    void resetMultiItemPosition();

private:
    QWidget *m_appWidget;
    QWidget *m_multiWindowWidget;
    Dock::DisplayMode m_displayMode;
};

#endif // MULTIWINDOWHELPER_H
