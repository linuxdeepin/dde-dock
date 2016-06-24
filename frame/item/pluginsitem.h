#ifndef PLUGINSITEM_H
#define PLUGINSITEM_H

#include "dockitem.h"
#include "pluginsiteminterface.h"

class PluginsItem : public DockItem
{
    Q_OBJECT

public:
    explicit PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent = 0);

    int itemSortKey() const;

private:
    void paintEvent(QPaintEvent *e);
    QSize sizeHint() const;

private:
    PluginsItemInterface * const m_pluginInter;
    const QString m_itemKey;

    PluginsItemInterface::PluginType m_type;
};

#endif // PLUGINSITEM_H
