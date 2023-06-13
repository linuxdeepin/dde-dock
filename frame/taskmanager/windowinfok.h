// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWINFOK_H
#define WINDOWINFOK_H

#include "windowinfobase.h"
#include "org_deepin_dde_kwayland_plasmawindow.h"

#include <QString>
#include <qobject.h>
#include <qscopedpointer.h>

using PlasmaWindow = org::deepin::dde::kwayland1::PlasmaWindow;

class Entry;
class ProcessInfo;

// wayland下窗口信息
class WindowInfoK: public WindowInfoBase
{
public:
    explicit WindowInfoK(PlasmaWindow *window, XWindow _xid = 0, QObject *parent = nullptr);
    virtual ~WindowInfoK() override;

    virtual bool shouldSkip() override;
    virtual QString getIcon() override;
    virtual QString getTitle() override;
    virtual bool isDemandingAttention() override;
    virtual bool allowClose() override;
    virtual void close(uint32_t timestamp) override;
    virtual void activate() override;
    virtual void minimize() override;
    virtual bool isMinimized() override;
    virtual int64_t getCreatedTime() override;
    virtual QString getDisplayName() override;
    virtual QString getWindowType() override;
    virtual void update() override;
    virtual void killClient() override;
    virtual QString uuid() override;
    QString getInnerId() override;

    QString getAppId();
    void setAppId(QString _appId);
    bool changeXid(XWindow _xid);
    PlasmaWindow *getPlasmaWindow();
    bool updateGeometry();
    void updateTitle();
    void updateDemandingAttention();
    void updateIcon();
    void updateAppId();
    void updateInternalId();
    void updateCloseable();
    void updateProcessInfo();
    DockRect getGeometry();

private:
    bool m_updateCalled;
    QString m_appId;
    uint32_t m_internalId;
    bool m_demaningAttention;
    bool m_closeable;
    DockRect m_geometry;
    QScopedPointer<PlasmaWindow> m_plasmaWindow;
};

#endif // WINDOWINFOK_H
