#include "mainpanel.h"

#include <QBoxLayout>
#include <QDragEnterEvent>

MainPanel::MainPanel(QWidget *parent)
    : QFrame(parent),
      m_itemLayout(new QBoxLayout(QBoxLayout::LeftToRight, this)),

      m_itemController(DockItemController::instance(this))
{
    setAcceptDrops(true);
    setObjectName("MainPanel");
    setStyleSheet("QWidget #MainPanel {"
                  "border:none;"
                  "background-color:red;"
                  "border-radius:5px 5px 5px 5px;"
                  "}");

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
        m_itemLayout->addWidget(item);

    setLayout(m_itemLayout);
}

void MainPanel::updateDockSide(const DockSettings::DockSide dockSide)
{
    switch (dockSide)
    {
    case DockSettings::Top:
    case DockSettings::Bottom:          m_itemLayout->setDirection(QBoxLayout::LeftToRight);    break;
    case DockSettings::Left:
    case DockSettings::Right:           m_itemLayout->setDirection(QBoxLayout::TopToBottom);    break;
    }
}

void MainPanel::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    adjustItemSize();
}

void MainPanel::dragEnterEvent(QDragEnterEvent *e)
{
    // TODO: check
    e->accept();
}

void MainPanel::dragMoveEvent(QDragMoveEvent *e)
{
    qDebug() << e;
}

void MainPanel::dropEvent(QDropEvent *e)
{
    qDebug() << e;
}

void MainPanel::adjustItemSize()
{
    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        switch (item->itemType())
        {
        case DockItem::App:     item->setFixedWidth(80);    break;
        default:;
        }
    }
}
