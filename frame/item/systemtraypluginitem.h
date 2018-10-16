#ifndef SYSTEMTRAYPLUGINITEM_H
#define SYSTEMTRAYPLUGINITEM_H

#include "pluginsitem.h"

class SystemTrayPluginItem : public PluginsItem
{
    Q_OBJECT

public:
    SystemTrayPluginItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent = 0);

    inline ItemType itemType() const Q_DECL_OVERRIDE {return ItemType::SystemTrayPlugin;}

Q_SIGNALS:
    void fashionSystemTraySizeChanged(const QSize &systemTraySize) const;

private:
    bool eventFilter(QObject *watched, QEvent *e) Q_DECL_OVERRIDE;
};

#endif // SYSTEMTRAYPLUGINITEM_H
