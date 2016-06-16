#include "pluginsitem.h"

#include "pluginsiteminterface.h"

#include <QPainter>
#include <QBoxLayout>

PluginsItem::PluginsItem(PluginsItemInterface* const inter, QWidget *parent)
    : DockItem(Plugins, parent),
      m_inter(inter)
{
//    QBoxLayout *layout = new QBoxLayout;
//    layout->addWidget(m_inter->);
}

void PluginsItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect().marginsRemoved(QMargins(3, 3, 3, 3)), Qt::cyan);
    painter.setPen(Qt::red);
    painter.drawText(rect(), m_inter->name());
}
