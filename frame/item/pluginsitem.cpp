#include "pluginsitem.h"

#include "pluginsiteminterface.h"

#include <QPainter>
#include <QBoxLayout>

PluginsItem::PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(Plugins, parent),
      m_pluginInter(pluginInter),
      m_itemKey(itemKey)
{
    m_type = pluginInter->pluginType(itemKey);

    if (m_type == PluginsItemInterface::Simple)
        return;

    // construct complex widget layout
    QBoxLayout *layout = new QHBoxLayout;
    layout->addWidget(m_pluginInter->itemWidget(itemKey));
    layout->setSpacing(0);
    layout->setMargin(0);
    setLayout(layout);
}
