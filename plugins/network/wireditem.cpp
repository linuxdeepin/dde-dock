#include "wireditem.h"

WiredItem::WiredItem(QWidget *parent)
    : QWidget(parent)
{

}

QSize WiredItem::sizeHint() const
{
    return QSize(24, 24);
}
