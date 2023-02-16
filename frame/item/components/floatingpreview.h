// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef FLOATINGPREVIEW_H
#define FLOATINGPREVIEW_H

#include <QWidget>
#include <QPointer>

#include <DIconButton>
#include <DPushButton>

DWIDGET_USE_NAMESPACE

class AppSnapshot;
class FloatingPreview : public QWidget
{
    Q_OBJECT

public:
    explicit FloatingPreview(QWidget *parent = 0);

    WId trackedWid() const;
    AppSnapshot *trackedWindow();
    void setFloatingTitleVisible(bool bVisible);

public slots:
    void trackWindow(AppSnapshot *const snap);

private:
    void paintEvent(QPaintEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private slots:
    void onCloseBtnClicked();

private:
    QPointer<AppSnapshot> m_tracked;

    DIconButton *m_closeBtn3D;
    DPushButton *m_titleBtn;
};

#endif // FLOATINGPREVIEW_H
