#ifndef DOCKITEM_H
#define DOCKITEM_H

#include <QFrame>

class DockItem : public QWidget
{
    Q_OBJECT

    enum ItemType {
        Launcher,
        App,
        Plugins,
    };

public:
    explicit DockItem(QWidget *parent = nullptr);

protected:
    void paintEvent(QPaintEvent *e);
};

#endif // DOCKITEM_H
