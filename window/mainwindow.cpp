#include "mainwindow.h"
#include "panel/mainpanel.h"

#include <QDebug>
#include <QResizeEvent>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),
      m_mainPanel(new MainPanel(this)),

      m_settings(new DockSettings(this)),
      m_displayInter(new DBusDisplay(this)),

      m_positionUpdateTimer(new QTimer(this))
{
    setWindowFlags(Qt::X11BypassWindowManagerHint);
    setAttribute(Qt::WA_TranslucentBackground);

    initComponents();
    initConnections();
}

void MainWindow::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    m_mainPanel->setFixedSize(e->size());
}

void MainWindow::keyPressEvent(QKeyEvent *e)
{
    switch (e->key())
    {
#ifdef QT_DEBUG
    case Qt::Key_Escape:        qApp->quit();       break;
#endif
    default:;
    }
}

void MainWindow::initComponents()
{
    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(200);
    m_positionUpdateTimer->start();
}

void MainWindow::initConnections()
{
    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition);
}

void MainWindow::updatePosition()
{
    const QRect rect = m_displayInter->primaryRect();

    setFixedWidth(rect.width());
    setFixedHeight(80);

    move(0, 950);
}
