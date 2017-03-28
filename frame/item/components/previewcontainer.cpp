#include "previewcontainer.h"

#include <QLabel>
#include <QWindow>
#include <QX11Info>

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>

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
        XWindowAttributes attrs;
        XGetWindowAttributes(QX11Info::display(), it.key(), &attrs);
        XImage *ximage = XGetImage(QX11Info::display(), it.key(), 0, 0, attrs.width, attrs.height, AllPlanes, ZPixmap);
        const QImage qimage((uchar*)(ximage->data), attrs.width, attrs.height, QImage::Format_ARGB32);
        XDestroyImage(ximage);

        QLabel *l = new QLabel;
        l->setFixedSize(250, 200);
        l->setPixmap(QPixmap::fromImage(qimage).scaled(250, 200, Qt::KeepAspectRatio, Qt::SmoothTransformation));
        m_windowListLayout->addWidget(l);
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
