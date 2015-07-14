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
#include "DBus/dbushidestatemanager.h"
#include "Controller/dockmodedata.h"
#include "Panel/panel.h"

class MainWidget : public QWidget
{
    Q_OBJECT
    Q_PROPERTY(QSize size READ size WRITE resize)
    Q_PROPERTY(QRect geometry READ geometry WRITE setGeometry)
    Q_PROPERTY(QPoint pos READ pos WRITE move)

public:
    MainWidget(QWidget *parent = 0);
    ~MainWidget();

public slots:
    void slotDockModeChanged(Dock::DockMode newMode,Dock::DockMode oldMode);

signals:
    void startShow();
    void startHide();

private:
    void hasShown();
    void hasHidden();
    void hideStateChanged(int value);
    void initHSManager();
    void initState();

private:
    Panel *mainPanel = NULL;
    DBusHideStateManager * m_HSManager = NULL;
};

#endif // MAINWIDGET_H
