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
    ~MainWindow();

private:
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void keyPressEvent(QKeyEvent *e);
    void setFixedSize(const QSize &size);
    void move(int x, int y);
    void initComponents();
    void initConnections();

private slots:
    void updatePosition();
    void updateGeometry();
    void clearStrutPartial();
    void setStrutPartial();

private:
    MainPanel *m_mainPanel;

    QTimer *m_positionUpdateTimer;
    QPropertyAnimation *m_sizeChangeAni;
    QPropertyAnimation *m_posChangeAni;

    XcbMisc *m_xcbMisc;
    DockSettings *m_settings;
};

#endif // MAINWINDOW_H
