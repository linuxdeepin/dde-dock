// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PREVIEWCONTAINER_H
#define PREVIEWCONTAINER_H

#include <QWidget>
#include <QBoxLayout>
#include <QTimer>
#include <QScrollArea>
#include <QScrollerProperties>

#include "constants.h"
#include "appsnapshot.h"
#include "floatingpreview.h"

#include <com_deepin_dde_daemon_dock_entry.h>

#include <DWindowManagerHelper>

DWIDGET_USE_NAMESPACE

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
    void updateDockSize(const int size);
    void checkMouseLeave();
    void prepareHide();

private:
    void adjustSize(bool composite);
    void appendSnapWidget(const WId wid);

    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *e) override;
    bool eventFilter(QObject *watcher, QEvent *event) override;

private slots:
    void onSnapshotClicked(const WId wid);
    void previewEntered(const WId wid);
    void previewFloating();
    void onRequestCloseAppSnapshot();

private:
    bool m_needActivate;
    bool m_canPreview;
    int m_dockSize;
    Dock::Position m_dockPos;
    QMap<WId, AppSnapshot *> m_snapshots;

    FloatingPreview *m_floatingPreview;
    QScrollArea * m_scrollArea;
    QWidget *m_windowListWidget;
    QBoxLayout *m_windowListLayout;
    QScrollerProperties m_sp;

    QTimer *m_preparePreviewTimer;
    QTimer *m_mouseLeaveTimer;
    DWindowManagerHelper *m_wmHelper;
    QTimer *m_waitForShowPreviewTimer;
    WId m_currentWId;
    TitleDisplayMode m_titleMode;
};

#endif // PREVIEWCONTAINER_H
