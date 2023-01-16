// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKITEM_H
#define DOCKITEM_H

#include "constants.h"
#include "dockpopupwindow.h"

#include <QFrame>
#include <QPointer>
#include <QGestureEvent>
#include <QMenu>

#include <memory>

using namespace Dock;

class DockItem : public QWidget
{
    Q_OBJECT

public:
    enum ItemType {
        Launcher,
        App,
        Plugins,
        FixedPlugin,
        Placeholder,
        OverflowIcon,
        TrayPlugin,
    };

public:
    explicit DockItem(QWidget *parent = nullptr);
    ~DockItem() override;

    static void setDockPosition(const Position side);
    static void setDockSize(const int size);
    static void setDockDisplayMode(const DisplayMode mode);

    inline virtual ItemType itemType() const {return App;}

    QSize sizeHint() const override;
    virtual QString accessibleName();

    inline int getClickedCount() const {
        return m_clickedcount;
    }

public slots:
    virtual void refreshIcon() {}

    void showPopupApplet(QWidget *const applet);
    void hidePopup();
    virtual void setDraging(bool bDrag);
    virtual void checkEntry() {}

    bool isDragging();
signals:
    void dragStarted() const;
    void itemDropped(QObject *destination, const QPoint &dropPoint) const;
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefreshWindowVisible() const;

protected:
    bool event(QEvent *event) override;
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;

    const QRect perfectIconRect() const;
    const QPoint popupMarkPoint() ;
    const QPoint topleftPoint() const;

    void hideNonModel();
    void popupWindowAccept();
    virtual void showPopupWindow(QWidget *const content, const bool model = false, const int radius = 6);
    virtual void showHoverTips();
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;
    virtual QWidget *popupTips();

    bool checkAndResetTapHoldGestureState();
    virtual void gestureEvent(QGestureEvent *event);

protected slots:
    void showContextMenu();
    void onContextMenuAccepted();

private:
    void updatePopupPosition();
    void menuActionClicked(QAction *action);

private:
    int  m_clickedcount = 0; // caculate the time be clicked

protected:
    bool m_hover;
    bool m_popupShown;
    bool m_tapAndHold;
    bool m_draging;
    QMenu m_contextMenu;

    QPointer<QWidget> m_lastPopupWidget;

    QTimer *m_popupTipsDelayTimer;
    QTimer *m_popupAdjustDelayTimer;

    static Position DockPosition;
    static int DockSize;
    static DisplayMode DockDisplayMode;
    static QPointer<DockPopupWindow> PopupWindow;
};

#endif // DOCKITEM_H
