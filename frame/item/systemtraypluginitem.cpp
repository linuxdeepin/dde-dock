#include "systemtraypluginitem.h"

#include <QEvent>

SystemTrayPluginItem::SystemTrayPluginItem(PluginsItemInterface * const pluginInter, const QString &itemKey, QWidget *parent)
    : PluginsItem(pluginInter, itemKey, parent)
{
}

bool SystemTrayPluginItem::eventFilter(QObject *watched, QEvent *e)
{
    // 时尚模式下
    // 监听插件Widget的"FashionSystemTraySize"属性
    // 当接收到这个属性变化的事件后，重新计算和设置dock的大小

    if (watched == centralWidget() && e->type() == QEvent::DynamicPropertyChange
            && static_cast<QDynamicPropertyChangeEvent *>(e)->propertyName() == "FashionSystemTraySize") {
        Q_EMIT fashionSystemTraySizeChanged(watched->property("FashionSystemTraySize").toSize());
    }

    return PluginsItem::eventFilter(watched, e);
}
