#include "_previewcontainer.h"

#include <QDesktopWidget>
#include <QScreen>
#include <QApplication>

#define SPACING           5
#define MARGIN            5

_PreviewContainer::_PreviewContainer(QWidget *parent)
    : QWidget(parent),

      m_mouseLeaveTimer(new QTimer(this)),
      m_wmHelper(DWindowManagerHelper::instance())
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setSpacing(SPACING);
    m_windowListLayout->setContentsMargins(MARGIN, MARGIN, MARGIN, MARGIN);

    m_mouseLeaveTimer->setSingleShot(true);
    m_mouseLeaveTimer->setInterval(10);

    setLayout(m_windowListLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_mouseLeaveTimer, &QTimer::timeout, this, &_PreviewContainer::checkMouseLeave, Qt::QueuedConnection);
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

void _PreviewContainer::checkMouseLeave()
{
    if (!underMouse())
    {
        emit requestCancelPreview();
        emit requestHidePreview();
    }
}

void _PreviewContainer::adjustSize()
{
    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    const int count = m_snapshots.size();

    if (!count)
        return;

    const QRect r = qApp->primaryScreen()->geometry();
    const int padding = 20;

    if (horizontal)
    {
        const int h = SNAP_HEIGHT + MARGIN * 2;
        const int w = SNAP_WIDTH * count + MARGIN * 2 + SPACING * (count - 1);

        setFixedHeight(h);
        setFixedWidth(std::min(w, r.width() - padding));
    } else {
        const int w = SNAP_WIDTH + MARGIN * 2;
        const int h = SNAP_HEIGHT * count + MARGIN * 2 + SPACING * (count - 1);

        setFixedWidth(w);
        setFixedHeight(std::min(h, r.height() - padding));
    }
}

void _PreviewContainer::appendSnapWidget(const WId wid)
{
    AppSnapshot *snap = new AppSnapshot(wid);

    connect(snap, &AppSnapshot::clicked, this, &_PreviewContainer::requestActivateWindow, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::clicked, this, &_PreviewContainer::requestCancelPreview);
    connect(snap, &AppSnapshot::clicked, this, &_PreviewContainer::requestHidePreview);
    connect(snap, &AppSnapshot::entered, this, &_PreviewContainer::requestPreviewWindow);

    m_windowListLayout->addWidget(snap);

    m_snapshots.insert(wid, snap);
}

void _PreviewContainer::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    m_mouseLeaveTimer->start();
}

void _PreviewContainer::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_mouseLeaveTimer->start();
}
