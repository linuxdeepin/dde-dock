// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once

#include <QWidget>
#include <QTimer>
#include <QPair>

class QDBusMessage;
class BaseTrayWidget : public QWidget
{
    Q_OBJECT

public:
    enum TrayType {
        ApplicationTray,
        SystemTray,
    };

    explicit BaseTrayWidget(QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    virtual ~BaseTrayWidget() override;

    virtual QString itemKeyForConfig() = 0;
    virtual void updateIcon() = 0;
    virtual void sendClick(uint8_t mouseButton, int x, int y) = 0;
    virtual inline TrayType trayType() const { return TrayType::ApplicationTray; } // default is ApplicationTray
    virtual bool isValid() {return true;}
    uint getOwnerPID();
    virtual bool needShow();
    virtual void setNeedShow(bool needShow);
    virtual QPixmap icon() = 0;
    virtual bool containsPoint(const QPoint& mouse) = 0;

Q_SIGNALS:
    void iconChanged();
    void clicked();
    void needAttention();
    void requestWindowAutoHide(const bool autoHide);
    void requestRefershWindowVisible();

protected:
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *e) override;

    void handleMouseRelease();
    const QRect perfectIconRect() const;
    void resizeEvent(QResizeEvent *event) override;
    void setOwnerPID(uint PID);

private:
    QTimer *m_handleMouseReleaseTimer;

    QPair<QPoint, Qt::MouseButton> m_lastMouseReleaseData;

    uint m_ownerPID;
    bool m_needShow;
};

