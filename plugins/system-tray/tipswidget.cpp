#include "tipswidget.h"
#include "traywidget.h"

TipsWidget::TipsWidget(QWidget *parent)
    : QWidget(parent),
      m_mainLayout(new QHBoxLayout)
{
    setLayout(m_mainLayout);
    setFixedHeight(32);
}

void TipsWidget::clear()
{
    QLayoutItem *item = nullptr;
    while ((item = m_mainLayout->takeAt(0)) != nullptr)
    {
        if (item->widget())
            item->widget()->setParent(nullptr);
        delete item;
    }
}

void TipsWidget::addWidgets(QList<TrayWidget *> widgets)
{
    for (auto w : widgets)
    {
        w->setVisible(true);
        m_mainLayout->addWidget(w);
    }
    setFixedWidth(widgets.size() * 20 + 20);
}
