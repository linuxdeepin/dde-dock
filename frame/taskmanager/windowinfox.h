// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWINFOX_H
#define WINDOWINFOX_H

#include "windowinfobase.h"
#include "xcbutils.h"

#include <QVector>
#include <qobject.h>
#include <qobjectdefs.h>

class AppInfo;

// X11下窗口信息 在明确X11环境下使用
class WindowInfoX: public WindowInfoBase
{
    Q_OBJECT
public:
    WindowInfoX(XWindow _xid = 0, QObject *parent = nullptr);
    virtual ~WindowInfoX() override;

    virtual bool shouldSkip() override;
    virtual QString getIcon() override;
    virtual QString getTitle() override;
    virtual bool isDemandingAttention() override;
    virtual void close(uint32_t timestamp) override;
    virtual void activate() override;
    virtual void minimize() override;
    virtual bool isMinimized() override;
    virtual int64_t getCreatedTime() override;
    virtual QString getDisplayName() override;
    virtual QString getWindowType() override;
    virtual bool allowClose() override;
    virtual void update() override;
    virtual void killClient() override;
    virtual QString uuid() override;

    QString genInnerId(WindowInfoX *winInfo);
    QString getGtkAppId();
    QString getFlatpakAppId();
    QString getWmRole();
    WMClass getWMClass();
    QString getWMName();
    void updateProcessInfo();
    bool getUpdateCalled();
    void setInnerId(QString _innerId);
    ConfigureEvent *getLastConfigureEvent();
    void setLastConfigureEvent(ConfigureEvent *event);
    bool isGeometryChanged(int _x, int _y, int _width, int _height);
    void setGtkAppId(QString _gtkAppId);

    /************************更新XCB窗口属性*********************/
    void updateWmWindowType();
    void updateWmAllowedActions();
    void updateWmState();
    void updateWmClass();
    void updateMotifWmHints();
    void updateWmName();
    void updateIcon();
    void updateHasXEmbedInfo();
    void updateHasWmTransientFor();

private:
    QString getIconFromWindow();
    bool isActionMinimizeAllowed();
    bool hasWmStateDemandsAttention();
    bool hasWmStateSkipTaskBar();
    bool hasWmStateModal();
    bool isValidModal();
    bool shouldSkipWithWMClass();

private:
    int16_t m_x;
    int16_t m_y;
    uint16_t m_width;
    uint16_t m_height;
    QVector<XCBAtom> m_wmState;
    QVector<XCBAtom> m_wmWindowType;
    QVector<XCBAtom> m_wmAllowedActions;
    bool m_hasWMTransientFor;
    WMClass m_wmClass;
    QString m_wmName;
    bool m_hasXEmbedInfo;

    // 自定义atom属性
    QString m_gtkAppId;
    QString m_flatpakAppId;
    QString m_wmRole;
    MotifWMHints m_motifWmHints;

    bool m_updateCalled;
    ConfigureEvent *m_lastConfigureNotifyEvent;
};

#endif // WINDOWINFOX_H
