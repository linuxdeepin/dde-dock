#include "pluginitemwrapper.h"

PluginItemWrapper::PluginItemWrapper(DockPluginInterface *plugin,
                                     QString uuid, QWidget * parent) :
    AbstractDockItem(parent)
{
    m_plugin = plugin;
    m_uuid = uuid;

    if (m_plugin) {
        QWidget * item = m_plugin->getItem(uuid);

        if (item) {
            setFixedSize(item->size());

            item->setParent(this);
        }
    }
}

QString PluginItemWrapper::uuid() const
{
    return m_uuid;
}

