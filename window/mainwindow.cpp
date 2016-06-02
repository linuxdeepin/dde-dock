#include "mainwindow.h"

#include <QDebug>
#include <QResizeEvent>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),
      m_mainPanel(new MainPanel(this))
{
}

void MainWindow::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    m_mainPanel->setFixedSize(e->size());
}
