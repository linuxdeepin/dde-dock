#include "dockpluginlayout.h"

DockPluginLayout::DockPluginLayout(QWidget *parent) : MovableLayout(parent)
{

}

QSize DockPluginLayout::sizeHint() const
{
    QSize size;
    int w = 0;
    int h = 0;
    switch (direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        size.setHeight(height());
        for (QWidget * widget : widgets()) {
            w += widget->width();
        }
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(width());
        for (QWidget * widget : widgets()) {
            h += widget->height();
        }
        size.setHeight(h + getLayoutSpacing() * widgets().count());
        break;
    }

    return size;
}

