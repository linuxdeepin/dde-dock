// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef APPSNAPSHOT_H
#define APPSNAPSHOT_H

#include "dbusutil.h"
#include "taskmanager/windowinfomap.h"

#include <DIconButton>
#include <DWindowManagerHelper>
#include <DPushButton>

#include <QWidget>
#include <QDebug>
#include <QTimer>

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

#define SNAP_WIDTH       200
#define SNAP_HEIGHT      130

#define SNAP_CLOSE_BTN_WIDTH     (24)
#define SNAP_CLOSE_BTN_MARGIN    (5)
// 标题到左右两边的距离
#define TITLE_MARGIN    (20)

// 标题的文本到标题背景两边的距离
#define BTN_TITLE_MARGIN (6)

// 高亮框边距
#define BORDER_MARGIN (8)

struct SHMInfo;
struct _XImage;
typedef _XImage XImage;

namespace Dock {
class TipsWidget;
}

class AppSnapshot : public QWidget
{
    Q_OBJECT

public:
    explicit AppSnapshot(const WId wid, QWidget *parent = Q_NULLPTR);

    inline WId wid() const { return m_wid; }
    inline bool attentioned() const { return m_windowInfo.attention; }
    inline bool closeAble() const { return m_closeAble; }
    inline void setCloseAble(const bool value) { m_closeAble = value; }
    inline const QPixmap snapshot() const { return m_pixmap; }
    inline const QString title() const { return m_windowInfo.title; }
    void setWindowState();
    void setTitleVisible(bool bVisible);
    QString appTitle() { return m_3DtitleBtn ? m_3DtitleBtn->text() : QString(); }
    bool isKWinAvailable();

signals:
    void entered(const WId wid) const;
    void clicked(const WId wid) const;
    void requestCheckWindow() const;
    void requestCloseAppSnapshot() const;

public slots:
    void fetchSnapshot();
    void closeWindow() const;
    void compositeChanged() const;
    void setWindowInfo(const WindowInfo &info);

private:
    void dragEnterEvent(QDragEnterEvent *e) override;
    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    bool eventFilter(QObject *watched, QEvent *e) override;
    void resizeEvent(QResizeEvent *event) override;
    SHMInfo *getImageDSHM();
    XImage *getImageXlib();
    QRect rectRemovedShadow(const QImage &qimage, unsigned char *prop_to_return_gtk);
    void getWindowState();
    void updateTitle();

private:
    const WId m_wid;
    WindowInfo m_windowInfo;

    bool m_closeAble;
    bool m_isWidowHidden;
    QPixmap m_pixmap;

    Dock::TipsWidget *m_title;
    DPushButton *m_3DtitleBtn;

    QTimer *m_waitLeaveTimer;
    DIconButton *m_closeBtn2D;
    DWindowManagerHelper *m_wmHelper;
};

#endif // APPSNAPSHOT_H
