#ifndef MAINPANEL_H
#define MAINPANEL_H

#include "controller/dockitemcontroller.h"
#include "util/docksettings.h"

#include <QFrame>
#include <QBoxLayout>

class MainPanel : public QFrame
{
    Q_OBJECT

public:
    explicit MainPanel(QWidget *parent = 0);

    void updateDockSide(const DockSettings::DockSide dockSide);

private:
    void resizeEvent(QResizeEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void dropEvent(QDropEvent *e);

    void adjustItemSize();

private:
    QBoxLayout *m_itemLayout;

    DockItemController *m_itemController;
};

#endif // MAINPANEL_H
