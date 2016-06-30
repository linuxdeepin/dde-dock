#ifndef DOCKITEM_H
#define DOCKITEM_H

#include "constants.h"

#include <QFrame>

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
    void requestWindowAutoHide(const bool autoHide) const;

protected:
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);

    const QRect perfectIconRect() const;

    void showContextMenu();
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;

protected:
    ItemType m_type;
    bool m_hover;

    DBusMenuManager *m_menuManagerInter;

    static Position DockPosition;
    static DisplayMode DockDisplayMode;
};

#endif // DOCKITEM_H
