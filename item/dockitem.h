#ifndef DOCKITEM_H
#define DOCKITEM_H

#include <QFrame>

#include "util/docksettings.h"

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
    void setDockSide(const DockSettings::DockSide side);

    ItemType itemType() const;

protected:
    void paintEvent(QPaintEvent *e);

protected:
    DockSettings::DockSide m_side;
    ItemType m_type;
};

#endif // DOCKITEM_H
