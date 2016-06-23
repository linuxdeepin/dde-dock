#ifndef MAINPANEL_H
#define MAINPANEL_H

#include "controller/dockitemcontroller.h"
#include "util/docksettings.h"

#include <QFrame>
#include <QTimer>
#include <QBoxLayout>

class MainPanel : public QFrame
{
    Q_OBJECT

public:
    explicit MainPanel(QWidget *parent = 0);

    void updateDockPosition(const Position dockPosition);
    void updateDockDisplayMode(const Dock::DisplayMode displayMode);

private:
    void resizeEvent(QResizeEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);
    void dropEvent(QDropEvent *e);

    void initItemConnection(DockItem *item);
    DockItem *itemAt(const QPoint &point);

private slots:
    void adjustItemSize();
    void itemInserted(const int index, DockItem *item);
    void itemRemoved(DockItem *item);
    void itemMoved(DockItem *item, const int index);
    void itemDragStarted();

private:
    Position m_position;
    DisplayMode m_displayMode;
    QBoxLayout *m_itemLayout;

    QTimer *m_itemAdjustTimer;
    DockItemController *m_itemController;

    static DockItem *DragingItem;
};

#endif // MAINPANEL_H
