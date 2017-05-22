#include "_previewcontainer.h"

_PreviewContainer::_PreviewContainer(QWidget *parent)
    : QWidget(parent),

      m_wmHelper(DWindowManagerHelper::instance())
{
    m_windowListLayout = new QVBoxLayout;
    m_windowListLayout->setSpacing(0);
    m_windowListLayout->setMargin(0);

    setLayout(m_windowListLayout);
}

void _PreviewContainer::setWindowInfos(const WindowDict &infos)
{
    qDebug() << infos;
}

void _PreviewContainer::updateLayoutDirection(const Dock::Position dockPos)
{
    if (m_wmHelper->hasComposite() && (dockPos == Dock::Top || dockPos == Dock::Bottom))
        m_windowListLayout->setDirection(QBoxLayout::LeftToRight);
    else
        m_windowListLayout->setDirection(QBoxLayout::TopToBottom);
}
