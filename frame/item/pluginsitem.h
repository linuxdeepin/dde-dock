#ifndef PLUGINSITEM_H
#define PLUGINSITEM_H

#include "dockitem.h"

class PluginsItemInterface;
class PluginsItem : public DockItem
{
    Q_OBJECT

public:
    explicit PluginsItem(PluginsItemInterface* const inter, QWidget *parent = 0);

private:
    PluginsItemInterface* const m_inter;
};

#endif // PLUGINSITEM_H
