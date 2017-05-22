#include "_previewcontainer.h"

#define FIXED_WIDTH       200
#define FIXED_HEIGHT      130
#define SPACING           5
#define MARGIN            5

_PreviewContainer::_PreviewContainer(QWidget *parent)
    : QWidget(parent),

      m_wmHelper(DWindowManagerHelper::instance())
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setSpacing(SPACING);
    m_windowListLayout->setContentsMargins(MARGIN, MARGIN, MARGIN, MARGIN);

    setLayout(m_windowListLayout);
    setFixedSize(FIXED_WIDTH, FIXED_HEIGHT);
}

void _PreviewContainer::setWindowInfos(const WindowDict &infos)
{
    // check removed window
    for (auto it(m_snapshots.begin()); it != m_snapshots.end();)
    {
        if (!infos.contains(it.key()))
        {
            m_windowListLayout->removeWidget(it.value());
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

    if (!count)
        return;

    if (horizontal)
    {
        setFixedHeight(FIXED_HEIGHT + MARGIN * 2);
        setFixedWidth(FIXED_WIDTH * count + MARGIN * 2 + SPACING * (count - 1));
    } else {
        setFixedWidth(FIXED_WIDTH + MARGIN * 2);
        setFixedHeight(FIXED_HEIGHT * count + MARGIN * 2 + SPACING * (count - 1));
    }
}

void _PreviewContainer::appendSnapWidget(const WId wid)
{
    AppSnapshot *snap = new AppSnapshot(wid);

    m_windowListLayout->addWidget(snap);

    m_snapshots.insert(wid, snap);
}
