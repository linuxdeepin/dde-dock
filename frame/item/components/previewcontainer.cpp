#include "previewcontainer.h"
#include "previewwidget.h"

#include <QLabel>
#include <QWindow>

PreviewContainer::PreviewContainer(QWidget *parent)
    : QWidget(parent)
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setMargin(0);
    m_windowListLayout->setSpacing(5);

    setLayout(m_windowListLayout);
}

void PreviewContainer::setWindowInfos(const WindowDict &infos)
{
    // TODO: optimize
    while (QLayoutItem *item = m_windowListLayout->takeAt(0))
    {
        item->widget()->deleteLater();
        delete item;
    }

    for (auto it(infos.cbegin()); it != infos.cend(); ++it)
    {
        PreviewWidget *w = new PreviewWidget(it.key());

        connect(w, &PreviewWidget::requestActivateWindow, this, &PreviewContainer::requestActivateWindow);

        m_windowListLayout->addWidget(w);
    }
}

void PreviewContainer::updateLayoutDirection(const Dock::Position dockPos)
{
    switch (dockPos)
    {
    case Dock::Top:
    case Dock::Bottom:
        m_windowListLayout->setDirection(QBoxLayout::LeftToRight);
        break;

    case Dock::Left:
    case Dock::Right:
        m_windowListLayout->setDirection(QBoxLayout::TopToBottom);
        break;
    }
}
