#include "previewcontainer.h"
#include "previewwidget.h"

#include <QLabel>
#include <QWindow>
#include <QDebug>

PreviewContainer::PreviewContainer(QWidget *parent)
    : QWidget(parent),

      m_mouseLeaveTimer(new QTimer(this))
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setMargin(5);
    m_windowListLayout->setSpacing(3);

    m_mouseLeaveTimer->setSingleShot(true);
    m_mouseLeaveTimer->setInterval(100);

    setLayout(m_windowListLayout);

    connect(m_mouseLeaveTimer, &QTimer::timeout, this, &PreviewContainer::checkMouseLeave, Qt::QueuedConnection);
}

void PreviewContainer::setWindowInfos(const WindowDict &infos)
{
    if (infos.isEmpty())
    {
        emit requestCancelPreview();
        emit requestHidePreview();

        return;
    }

    QList<WId> removedWindows;

    // remove desroyed window
    for (auto it(m_windows.cbegin()); it != m_windows.cend(); ++it)
    {
        if (infos.contains(it.key()))
            continue;

        removedWindows << it.key();
        m_windowListLayout->removeWidget(it.value());
        it.value()->deleteLater();
    }
    for (auto id : removedWindows)
        m_windows.remove(id);

    for (auto it(infos.cbegin()); it != infos.cend(); ++it)
    {
        if (m_windows.contains(it.key()))
            continue;

        PreviewWidget *w = new PreviewWidget(it.key());
        w->setTitle(it.value());

        connect(w, &PreviewWidget::requestActivateWindow, this, &PreviewContainer::requestActivateWindow);
        connect(w, &PreviewWidget::requestPreviewWindow, this, &PreviewContainer::requestPreviewWindow);
        connect(w, &PreviewWidget::requestCancelPreview, this, &PreviewContainer::requestCancelPreview);
        connect(w, &PreviewWidget::requestHidePreview, this, &PreviewContainer::requestHidePreview);

        m_windowListLayout->addWidget(w);
        m_windows.insert(it.key(), w);
    }

    // update geometry
    QMetaObject::invokeMethod(this, "updateContainerSize", Qt::QueuedConnection);
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

void PreviewContainer::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_mouseLeaveTimer->start();
}

void PreviewContainer::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    m_mouseLeaveTimer->stop();
}

void PreviewContainer::updateContainerSize()
{
    resize(sizeHint());
}

void PreviewContainer::checkMouseLeave()
{
    const QPoint p = mapFromGlobal(QCursor::pos());

    if (!rect().contains(p))
        emit requestCancelPreview();
}
