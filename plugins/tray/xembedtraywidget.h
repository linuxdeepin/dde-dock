// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef XEMBEDTRAYWIDGET_H
#define XEMBEDTRAYWIDGET_H

#include "abstracttraywidget.h"

#include <QWidget>
#include <QTimer>

#include <xcb/xcb.h>

typedef struct _XDisplay Display;

class XEmbedTrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    explicit XEmbedTrayWidget(quint32 winId, xcb_connection_t *cnn = nullptr, Display *disp = nullptr, QWidget *parent = nullptr);
    ~XEmbedTrayWidget();

    QString itemKeyForConfig() override;
    void updateIcon() override;
    void sendClick(uint8_t mouseButton, int x, int y) override;

    static QString getWindowProperty(quint32 winId, QString propName);
    static QString toXEmbedKey(quint32 winId);
    static uint getWindowPID(quint32 winId);
    static bool isXEmbedKey(const QString &itemKey);
    virtual bool isValid() override {return m_valid;}

private:
    void showEvent(QShowEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void configContainerPosition();

    void wrapWindow();
    void sendHoverEvent();
    void refershIconImage();

    static QString getAppNameForWindow(quint32 winId);

private slots:
    void setX11PassMouseEvent(const bool pass);
    void setWindowOnTop(const bool top);
    bool isBadWindow();

private:
    // Direct client关注xevent，使用xevent来处理button事件等
    // XTest client不关注xevent，使用xtest extension处理
    enum InjectMode {
        Direct,
        XTest,
    };

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
    InjectMode m_injectMode;
};

#endif // XEMBEDTRAYWIDGET_H
