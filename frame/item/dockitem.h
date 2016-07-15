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
        Placeholder,
        Plugins,
    };

public:
    explicit DockItem(const ItemType type, QWidget *parent = nullptr);
    static void setDockPosition(const Position side);
    static void setDockDisplayMode(const DisplayMode mode);

    ItemType itemType() const;

signals:
    void dragStarted() const;
    void menuUnregistered() const;

protected:
    void paintEvent(QPaintEvent *e);
    void mouseMoveEvent(QMouseEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);

    const QRect perfectIconRect() const;
    const QPoint popupMarkPoint();

    void showContextMenu();
    void showHoverTips();
    void showPopupWindow(QWidget *content, const bool model = false);
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;
    virtual QWidget *popupTips();

protected:
    ItemType m_type;
    bool m_hover;

    QTimer *m_popupTipsDelayTimer;

    DBusMenuManager *m_menuManagerInter;

    static Position DockPosition;
    static DisplayMode DockDisplayMode;
    static std::unique_ptr<DockPopupWindow> PopupWindow;
};

#endif // DOCKITEM_H
