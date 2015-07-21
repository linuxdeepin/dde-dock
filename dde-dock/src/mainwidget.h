#ifndef MAINWIDGET_H
#define MAINWIDGET_H

#include <QApplication>
#include <QDesktopWidget>
#include <QWidget>
#include <QScreen>
#include <QStateMachine>
#include <QState>
#include <QPropertyAnimation>
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

private:
    void showDock();
    void hideDock();

    void changeDockMode(Dock::DockMode newMode,Dock::DockMode oldMode);
private:
    Panel *mainPanel = NULL;
    bool hasHidden = false;
    DockModeData * m_dmd = DockModeData::instance();
};

#endif // MAINWIDGET_H
