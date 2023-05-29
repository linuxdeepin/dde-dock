// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef POPUPCONTROLWIDGET_H
#define POPUPCONTROLWIDGET_H

#include "trashhelper.h"

#include <QWidget>

class PopupControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PopupControlWidget(QWidget *parent = 0);

    bool empty() const;
    int trashItems() const;
    QSize sizeHint() const;
//    static const QString trashDir();

public slots:
    void openTrashFloder();
    void clearTrashFloder();

signals:
    void emptyChanged(const bool empty) const;

private:
    int trashItemCount() const;

private slots:
    void trashStatusChanged();

private:
    bool m_empty;
    int m_trashItemsCount;

    TrashHelper *m_trashHelper;
};

#endif // POPUPCONTROLWIDGET_H
