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

int PluginsItem::itemSortKey() const
{
    return m_pluginInter->itemSortKey(m_itemKey);
}

void PluginsItem::paintEvent(QPaintEvent *e)
{
    if (m_type == PluginsItemInterface::Complex)
        return DockItem::paintEvent(e);

    QPainter painter(this);

    const QIcon icon = m_pluginInter->itemIcon(m_itemKey);
    const QRect iconRect = perfectIconRect();
    const QPixmap pixmap = icon.pixmap(iconRect.size());
    painter.drawPixmap(iconRect, pixmap);
}

QSize PluginsItem::sizeHint() const
{
    if (m_type == PluginsItemInterface::Complex)
        return DockItem::sizeHint();

    // TODO: icon size on efficient mode
    return QSize(48, 48);
}
