#include "mainpanel.h"
#include "item/appitem.h"

#include <QBoxLayout>
#include <QDragEnterEvent>

DockItem *MainPanel::DragingItem = nullptr;

MainPanel::MainPanel(QWidget *parent)
    : QFrame(parent),
      m_position(Dock::Top),
      m_itemLayout(new QBoxLayout(QBoxLayout::LeftToRight, this)),

      m_itemController(DockItemController::instance(this))
{
    m_itemLayout->setSpacing(0);
    m_itemLayout->setContentsMargins(0, 0, 0, 0);

    setAcceptDrops(true);
    setObjectName("MainPanel");
    setStyleSheet("QWidget #MainPanel {"
                  "border:none;"
                  "background-color:green;"
                  "border-radius:5px 5px 5px 5px;"
                  "}");

    connect(m_itemController, &DockItemController::itemInserted, this, &MainPanel::itemInserted);
    connect(m_itemController, &DockItemController::itemRemoved, this, &MainPanel::itemRemoved);

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        initItemConnection(item);
        m_itemLayout->addWidget(item);
    }

    setLayout(m_itemLayout);
}

void MainPanel::updateDockPosition(const Position dockPosition)
{
    m_position = dockPosition;

    switch (m_position)
    {
    case Position::Top:
    case Position::Bottom:          m_itemLayout->setDirection(QBoxLayout::LeftToRight);    break;
    case Position::Left:
    case Position::Right:           m_itemLayout->setDirection(QBoxLayout::TopToBottom);    break;
    }

    adjustItemSize();
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

    if (qobject_cast<DockItem *>(e->source()))
        DragingItem->show();
}

void MainPanel::dragMoveEvent(QDragMoveEvent *e)
{
    DockItem *item = itemAt(e->pos());
    if (item == DragingItem)
        return;

    m_itemController->itemMove(DragingItem, item);
}

void MainPanel::dragLeaveEvent(QDragLeaveEvent *e)
{
    Q_UNUSED(e)

    DragingItem->hide();
}

void MainPanel::dropEvent(QDropEvent *e)
{
    Q_UNUSED(e)
//    qDebug() << e;
}

void MainPanel::initItemConnection(DockItem *item)
{
    connect(item, &DockItem::dragStarted, this, &MainPanel::itemDragStarted);
}

DockItem *MainPanel::itemAt(const QPoint &point)
{
    const QList<DockItem *> itemList = m_itemController->itemList();

    for (auto item : itemList)
    {
        QRect rect;
        rect.setTopLeft(item->pos());
        rect.setSize(item->size());

        if (rect.contains(point))
            return item;
    }

    return nullptr;
}

void MainPanel::adjustItemSize()
{
    QSize itemSize;
    switch (m_position)
    {
    case Top:
    case Bottom:
        itemSize.setHeight(height());
        itemSize.setWidth(AppItem::itemBaseWidth());
        break;

    case Left:
    case Right:
        itemSize.setHeight(AppItem::itemBaseHeight());
        itemSize.setWidth(width());
        break;

    default:
        Q_ASSERT(false);
    }

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
    {
        switch (item->itemType())
        {
        case DockItem::Launcher:
        case DockItem::App:     item->setFixedSize(itemSize);    break;
        default:;
        }
    }
}

void MainPanel::itemInserted(const int index, DockItem *item)
{
    initItemConnection(item);
    m_itemLayout->insertWidget(index, item);

    adjustSize();
}

void MainPanel::itemRemoved(DockItem *item)
{
    m_itemLayout->removeWidget(item);
}

void MainPanel::itemDragStarted()
{
    DragingItem = qobject_cast<DockItem *>(sender());

    QRect rect;
    rect.setTopLeft(mapToGlobal(pos()));
    rect.setSize(size());

    DragingItem->setVisible(rect.contains(QCursor::pos()));
}
