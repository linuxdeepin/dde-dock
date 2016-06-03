#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "xcb/xcb_misc.h"
#include "panel/mainpanel.h"
#include "dbus/dbusdisplay.h"

#include <QWidget>
#include <QTimer>

class MainWindow : public QWidget
{
    Q_OBJECT

    enum Position {
        TOP,
        BOTTOM,
        LEFT,
        RIGHT,
    };

public:
    explicit MainWindow(QWidget *parent = 0);

private:
    void resizeEvent(QResizeEvent *e);
    void keyPressEvent(QKeyEvent *e);

private slots:
    void updatePosition();

private:
    Position m_position;

    MainPanel *m_mainPanel;

    DBusDisplay *m_displayInter;

    QTimer *m_positionUpdateTimer;
};

#endif // MAINWINDOW_H
