#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "xcb/xcb_misc.h"
#include "panel/mainpanel.h"
#include "controller/dockitemcontroller.h"

#include <QWidget>

class MainWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);

private:
    void resizeEvent(QResizeEvent *e);

    MainPanel *m_mainPanel;
    DockItemController *m_itemController;
};

#endif // MAINWINDOW_H
