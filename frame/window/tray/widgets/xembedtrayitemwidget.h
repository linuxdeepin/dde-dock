// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef XEMBEDTRAYWIDGET_H
#define XEMBEDTRAYWIDGET_H

#include "basetraywidget.h"

#include <QWidget>
#include <QTimer>

#include <xcb/xcb.h>

typedef struct _XDisplay Display;

class XEmbedTrayItemWidget : public BaseTrayWidget
{
    Q_OBJECT

public:
    explicit XEmbedTrayItemWidget(quint32 winId, xcb_connection_t *cnn = nullptr, Display *disp = nullptr, QWidget *parent = nullptr);
    ~XEmbedTrayItemWidget() override;

    QString itemKeyForConfig() override;
    void updateIcon() override;
    void sendClick(uint8_t mouseButton, int x, int y) override;

    static QString toXEmbedKey(quint32 winId);
    static uint getWindowPID(quint32 winId);
    static bool isXEmbedKey(const QString &itemKey);
    virtual bool isValid() override {return m_valid;}
    QPixmap icon() override;
    bool containsPoint(const QPoint &mouse) override { return false; }

private:
    void showEvent(QShowEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void configContainerPosition();

    void wrapWindow();
    void sendHoverEvent();
    void refershIconImage();

private slots:
    void setX11PassMouseEvent(const bool pass);
    void setWindowOnTop(const bool top);
    bool isBadWindow();

private:
    bool m_active = false;
    WId m_windowId;
    WId m_containerWid;
    QImage m_image;
    QString m_appName;

    QTimer *m_updateTimer;
    QTimer *m_sendHoverEvent;
    bool m_valid;
    xcb_connection_t *m_xcbCnn;
    Display* m_display;
};

#endif // XEMBEDTRAYWIDGET_H
