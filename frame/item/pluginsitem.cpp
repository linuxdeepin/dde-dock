#include "pluginsitem.h"

#include "pluginsiteminterface.h"

#include <QPainter>
#include <QBoxLayout>

PluginsItem::PluginsItem(PluginsItemInterface* const inter, QWidget *parent)
    : DockItem(Plugins, parent),
      m_inter(inter)
{
    QBoxLayout *layout = new QHBoxLayout;
    layout->addWidget(m_inter->centeralWidget());
    layout->setSpacing(0);
    layout->setMargin(0);

    setLayout(layout);
}
