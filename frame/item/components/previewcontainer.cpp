#include "previewcontainer.h"
#include "previewwidget.h"

#include <QLabel>
#include <QWindow>
#include <QDebug>

PreviewContainer::PreviewContainer(QWidget *parent)
    : QWidget(parent)
{
    m_windowListLayout = new QBoxLayout(QBoxLayout::LeftToRight);
    m_windowListLayout->setMargin(5);
    m_windowListLayout->setSpacing(3);

    setLayout(m_windowListLayout);
    setMouseTracking(true);
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
        w->setTitle(it.value());

        connect(w, &PreviewWidget::requestActivateWindow, this, &PreviewContainer::requestActivateWindow);
        connect(w, &PreviewWidget::requestPreviewWindow, this, &PreviewContainer::requestPreviewWindow);
        connect(w, &PreviewWidget::requestCancelPreview, this, &PreviewContainer::requestCancelPreview);
        connect(w, &PreviewWidget::requestHidePreview, this, &PreviewContainer::requestHidePreview);

        m_windowListLayout->addWidget(w);
    }

    if (!isVisible())
        return;

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

    const QPoint p = mapFromGlobal(QCursor::pos()) + pos();

    if (!rect().contains(p))
        emit requestCancelPreview();
}

void PreviewContainer::updateContainerSize()
{
    resize(sizeHint());
}
