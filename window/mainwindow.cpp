#include "mainwindow.h"

#include <QDebug>
#include <QResizeEvent>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),
      m_position(BOTTOM),

      m_mainPanel(new MainPanel(this)),

      m_displayInter(new DBusDisplay(this)),

      m_positionUpdateTimer(new QTimer(this))
{
    setWindowFlags(Qt::FramelessWindowHint);
    setAttribute(Qt::WA_TranslucentBackground);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition);

    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(200);
    m_positionUpdateTimer->start();
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

void MainWindow::updatePosition()
{
    const QRect rect = m_displayInter->primaryRect();

    setFixedWidth(rect.width());
    setFixedHeight(80);

    move(0, 0);
}
