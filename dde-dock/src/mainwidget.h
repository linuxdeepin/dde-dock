#ifndef MAINWIDGET_H
#define MAINWIDGET_H

#include <QApplication>
#include <QDesktopWidget>
#include <QWidget>
#include <QScreen>
#include <QDebug>
#include "Controller/dockmodedata.h"
#include "Panel/panel.h"

class MainWidget : public QWidget
{
    Q_OBJECT

public:
    MainWidget(QWidget *parent = 0);
    ~MainWidget();

public slots:
    void slotDockModeChanged(Dock::DockMode newMode,Dock::DockMode oldMode);

private:
    Panel *mainPanel = NULL;
};

#endif // MAINWIDGET_H
