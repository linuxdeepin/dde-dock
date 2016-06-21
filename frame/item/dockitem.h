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
    void setDockSide(const Position side);

    ItemType itemType() const;

signals:
    void dragStarted() const;

protected:
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

    void showContextMenu();
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;

protected:
    Position m_side;
    ItemType m_type;

    DBusMenuManager *m_menuManagerInter;
};

#endif // DOCKITEM_H
