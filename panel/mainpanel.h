#ifndef MAINPANEL_H
#define MAINPANEL_H

#include "controller/dockitemcontroller.h"

#include <QFrame>

class MainPanel : public QFrame
{
    Q_OBJECT

public:
    explicit MainPanel(QWidget *parent = 0);

private:
    DockItemController *m_itemController;
};

#endif // MAINPANEL_H
