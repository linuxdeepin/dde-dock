#ifndef DOCKITEM_H
#define DOCKITEM_H

#include "constants.h"
#include "util/dockpopupwindow.h"

#include <QFrame>

#include <memory>

using namespace Dock;

class DBusMenuManager;
class DockItem : public QWidget
{
    Q_OBJECT

public:
    enum ItemType {
        Launcher,
        App,
        Stretch,
        Plugins,
        Container,
    };

public:
    explicit DockItem(QWidget *parent = nullptr);
    ~DockItem();

    static void setDockPosition(const Position side);
    static void setDockDisplayMode(const DisplayMode mode);

    inline virtual ItemType itemType() const {Q_UNREACHABLE(); return App;}

public slots:
    virtual void refershIcon() {}

signals:
    void dragStarted() const;
    void itemDropped(QObject *destination) const;
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;

protected:
    void paintEvent(QPaintEvent *e);
    void moveEvent(QMoveEvent *e);
//    void mouseMoveEvent(QMouseEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);

    const QRect perfectIconRect() const;
    const QPoint popupMarkPoint();

    void hidePopup();
    void popupWindowAccept();
    void showPopupApplet(QWidget * const applet);
    void showPopupWindow(QWidget * const content, const bool model = false);
    virtual void showHoverTips();
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;
    virtual QWidget *popupTips();

protected slots:
    void showContextMenu();

private:
    void updatePopupPosition();

protected:
    bool m_hover;
    bool m_popupShown;

    QTimer *m_popupTipsDelayTimer;

    DBusMenuManager *m_menuManagerInter;

    static Position DockPosition;
    static DisplayMode DockDisplayMode;
    static std::unique_ptr<DockPopupWindow> PopupWindow;
};

#endif // DOCKITEM_H
