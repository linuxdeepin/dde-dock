#include "_previewcontainer.h"

#include <QDesktopWidget>
#include <QScreen>
#include <QApplication>

#define SPACING           0
#define MARGIN            0
#define SNAP_HEIGHT_WITHOUT_COMPOSITE       30

_PreviewContainer::_PreviewContainer(QWidget *parent)
    : QWidget(parent),

      m_floatingPreview(new FloatingPreview(this)),
      m_mouseLeaveTimer(new QTimer(this)),
      m_wmHelper(DWindowManagerHelper::instance())
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setSpacing(SPACING);
    m_windowListLayout->setContentsMargins(MARGIN, MARGIN, MARGIN, MARGIN);

    m_mouseLeaveTimer->setSingleShot(true);
    m_mouseLeaveTimer->setInterval(300);

    m_floatingPreview->setVisible(false);

    setLayout(m_windowListLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);

    connect(m_mouseLeaveTimer, &QTimer::timeout, this, &_PreviewContainer::checkMouseLeave, Qt::QueuedConnection);
    connect(m_floatingPreview, &FloatingPreview::requestMove, this, &_PreviewContainer::moveFloatingPreview);
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
        m_snapshots[it.key()]->setWindowTitle(it.value());
    }

    if (m_snapshots.isEmpty())
    {
        emit requestCancelPreview();
        emit requestHidePreview();
    } else {
        adjustSize();
    }
}

void _PreviewContainer::updateSnapshots()
{
    for (AppSnapshot *snap : m_snapshots)
        snap->fetchSnapshot();
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
    const bool hover = underMouse();

    if (!hover)
    {
        m_floatingPreview->setVisible(false);

        if (!isVisible())
            return;

        emit requestCancelPreview();
        emit requestHidePreview();
    }
}

void _PreviewContainer::adjustSize()
{
    const int count = m_snapshots.size();
    const bool composite = m_wmHelper->hasComposite();
    if (!composite)
    {
        const int h = SNAP_HEIGHT_WITHOUT_COMPOSITE * count + MARGIN * 2 + SPACING * (count - 1);
        setFixedSize(SNAP_WIDTH, h);
        return;
    }

    const QRect r = qApp->primaryScreen()->geometry();
    const int padding = 20;

    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
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
    connect(snap, &AppSnapshot::clicked, this, &_PreviewContainer::requestCancelPreview, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::clicked, this, &_PreviewContainer::requestHidePreview, Qt::QueuedConnection);
    connect(snap, &AppSnapshot::entered, this, &_PreviewContainer::previewEntered, Qt::QueuedConnection);

    m_windowListLayout->addWidget(snap);

    m_snapshots.insert(wid, snap);
}

void _PreviewContainer::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    m_mouseLeaveTimer->stop();
}

void _PreviewContainer::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_mouseLeaveTimer->start();
}

void _PreviewContainer::previewEntered(const WId wid)
{
    if (!m_wmHelper->hasComposite())
        return;

    AppSnapshot *snap = static_cast<AppSnapshot *>(sender());

    m_floatingPreview->trackWindow(snap);
    m_floatingPreview->setVisible(true);
    m_floatingPreview->raise();

    emit requestPreviewWindow(wid);
}

void _PreviewContainer::moveFloatingPreview(const QPoint &p)
{
    const bool horizontal = m_windowListLayout->direction() == QBoxLayout::LeftToRight;
    const QRect r = rect();

    if (horizontal)
    {
        if (p.x() < r.left())
            m_floatingPreview->move(MARGIN, p.y());
        else if (p.x() + m_floatingPreview->width() > r.right())
            m_floatingPreview->move(r.right() - m_floatingPreview->width() - MARGIN + 1, p.y());
        else
            m_floatingPreview->move(p);
    } else {
        if (p.y() < r.top())
            m_floatingPreview->move(p.x(), MARGIN);
        else if (p.y() + m_floatingPreview->height() > r.bottom())
            m_floatingPreview->move(p.x(), r.bottom() - m_floatingPreview->height() - MARGIN + 1);
        else
            m_floatingPreview->move(p);
    }
}
