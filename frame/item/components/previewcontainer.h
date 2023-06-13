// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PREVIEWCONTAINER_H
#define PREVIEWCONTAINER_H

#include <QWidget>
#include <QBoxLayout>
#include <QTimer>

#include "constants.h"
#include "appsnapshot.h"
#include "floatingpreview.h"

#include <DWindowManagerHelper>

DWIDGET_USE_NAMESPACE
typedef QList<quint32> WindowList;
class PreviewContainer : public QWidget
{
    Q_OBJECT

public:
    explicit PreviewContainer(QWidget *parent = 0);

    enum TitleDisplayMode {
        HoverShow       = 0,
        AlwaysShow      = 1,
        AlwaysHide      = 2,
    };

signals:
    void requestActivateWindow(const WId wid) const;
    void requestPreviewWindow(const WId wid) const;
    void requestCheckWindows() const;
    void requestCancelPreviewWindow() const;
    void requestHidePopup() const;

public:
    void setWindowInfos(const WindowInfoMap &infos, const WindowList &allowClose);
    void setTitleDisplayMode(int mode);

public slots:
    void updateLayoutDirection(const Dock::Position dockPos);
    void checkMouseLeave();
    void prepareHide();

private:
    void adjustSize(bool composite);
    void appendSnapWidget(const WId wid);

    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);

private slots:
    void onSnapshotClicked(const WId wid);
    void previewEntered(const WId wid);
    void previewFloating();
    void onRequestCloseAppSnapshot();

private:
    bool m_needActivate;
    QMap<WId, AppSnapshot *> m_snapshots;

    FloatingPreview *m_floatingPreview;
    QBoxLayout *m_windowListLayout;

    QTimer *m_mouseLeaveTimer;
    DWindowManagerHelper *m_wmHelper;
    QTimer *m_waitForShowPreviewTimer;
    WId m_currentWId;
    TitleDisplayMode m_titleMode;
};

#endif // PREVIEWCONTAINER_H
