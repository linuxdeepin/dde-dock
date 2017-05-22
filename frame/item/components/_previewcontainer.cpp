#include "_previewcontainer.h"

#define FIXED_WIDTH       200
#define FIXED_HEIGHT      130

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
    // check removed window
    for (auto it(m_snapshots.begin()); it != m_snapshots.end();)
    {
        if (!infos.contains(it.key()))
        {
            it.value()->deleteLater();
            it = m_snapshots.erase(it);
        } else {
            ++it;
        }
    }

    for (auto it(infos.cbegin()); it != infos.cend(); ++it)
    {
        if (!m_snapshots.contains(it.key()))
            appendSnapWidget(it.key());
    }

    adjustSize();
}

void _PreviewContainer::updateLayoutDirection(const Dock::Position dockPos)
{
    if (m_wmHelper->hasComposite() && (dockPos == Dock::Top || dockPos == Dock::Bottom))
        m_windowListLayout->setDirection(QBoxLayout::LeftToRight);
    else
        m_windowListLayout->setDirection(QBoxLayout::TopToBottom);

    adjustSize();
}

void _PreviewContainer::adjustSize()
{
    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    const int count = m_snapshots.size();

    if (horizontal)
    {
        setFixedHeight(FIXED_HEIGHT);
        setFixedWidth(FIXED_WIDTH * count);
    } else {
        setFixedWidth(FIXED_WIDTH);
        setFixedHeight(FIXED_HEIGHT * count);
    }
}

void _PreviewContainer::appendSnapWidget(const WId wid)
{
    AppSnapshot *snap = new AppSnapshot;

    m_windowListLayout->addWidget(snap);

    m_snapshots.insert(wid, snap);
}
