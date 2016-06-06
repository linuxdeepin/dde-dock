#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "xcb/xcb_misc.h"
#include "dbus/dbusdisplay.h"
#include "util/docksettings.h"

#include <QWidget>
#include <QTimer>

class MainPanel;
class MainWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);

private:
    void resizeEvent(QResizeEvent *e);
    void keyPressEvent(QKeyEvent *e);
    void initComponents();
    void initConnections();

private slots:
    void updatePosition();

private:
    MainPanel *m_mainPanel;

    DockSettings *m_settings;
    DBusDisplay *m_displayInter;

    QTimer *m_positionUpdateTimer;
};

#endif // MAINWINDOW_H
