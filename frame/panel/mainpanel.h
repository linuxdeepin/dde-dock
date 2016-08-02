#ifndef MAINPANEL_H
#define MAINPANEL_H

#include "controller/dockitemcontroller.h"
#include "util/docksettings.h"

#include <QFrame>
#include <QTimer>
#include <QBoxLayout>

#define xstr(s) str(s)
#define str(s) #s
#define PANEL_BORDER    1
#define PANEL_PADDING   6

class MainPanel : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(int displayMode READ displayMode DESIGNABLE true)
    Q_PROPERTY(int position READ position DESIGNABLE true)

public:
    explicit MainPanel(QWidget *parent = 0);

    void updateDockPosition(const Position dockPosition);
    void updateDockDisplayMode(const Dock::DisplayMode displayMode);
    int displayMode();
    int position();

signals:
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;

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
    static PlaceholderItem *RequestDockItem;
};

#endif // MAINPANEL_H
