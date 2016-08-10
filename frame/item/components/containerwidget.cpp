#include "containerwidget.h"

#include <QDebug>

#define ITEM_HEIGHT         30
#define ITEM_WIDTH          30

ContainerWidget::ContainerWidget(QWidget *parent)
    : QWidget(parent),

      m_centeralLayout(new QHBoxLayout)
{
    m_centeralLayout->setSpacing(0);
    m_centeralLayout->setMargin(0);

    setLayout(m_centeralLayout);
    setFixedHeight(ITEM_HEIGHT);
    setFixedWidth(ITEM_WIDTH);
}

void ContainerWidget::addWidget(QWidget * const w)
{
    w->setFixedSize(ITEM_WIDTH, ITEM_HEIGHT);
    m_centeralLayout->addWidget(w);
    m_itemList.append(w);

    setFixedWidth(std::max(ITEM_WIDTH, ITEM_WIDTH * m_itemList.size()));
}
