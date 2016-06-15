#include "placeholderitem.h"

PlaceholderItem::PlaceholderItem(QWidget *parent)
    : DockItem(Placeholder, parent)
{
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
}
