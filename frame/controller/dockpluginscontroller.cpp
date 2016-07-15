#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dockitemcontroller.h"

#include <QDebug>
#include <QDir>

DockPluginsController::DockPluginsController(DockItemController *itemControllerInter)
    : QObject(itemControllerInter),
      m_itemControllerInter(itemControllerInter)
{
    qApp->installEventFilter(this);

    QMetaObject::invokeMethod(this, "loadPlugins", Qt::QueuedConnection);
}

DockPluginsController::~DockPluginsController()
{
}

void DockPluginsController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = new PluginsItem(itemInter, itemKey);

    // check if same item added
    if (m_pluginList.contains(itemInter))
        Q_ASSERT(!m_pluginList[itemInter].contains(itemKey));

    m_pluginList[itemInter][itemKey] = item;

    emit pluginItemInserted(item);
}

void DockPluginsController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    Q_ASSERT(item);

    item->update();
}

void DockPluginsController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(item);

    m_pluginList[itemInter].remove(itemKey);
    item->deleteLater();
}

//void DockPluginsController::requestPopupApplet(PluginsItemInterface * const itemInter, const QString &itemKey)
//{
//    PluginsItem *item = pluginItemAt(itemInter, itemKey);

//    Q_ASSERT(item);
//    item->showPopupApplet();
//}

void DockPluginsController::loadPlugins()
{
//    Q_ASSERT(m_pluginLoaderList.isEmpty());
//    Q_ASSERT(m_pluginsInterfaceList.isEmpty());

#ifdef QT_DEBUG
    const QDir pluginsDir("plugins");
#else
    const QDir pluginsDir("../lib/dde-dock/plugins");
#endif
    const QStringList plugins = pluginsDir.entryList(QDir::Files);

    for (const QString file : plugins)
    {
        if (!QLibrary::isLibrary(file))
            continue;

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-"))
            continue;

        const QString pluginFilePath = pluginsDir.absoluteFilePath(file);
        QPluginLoader *pluginLoader = new QPluginLoader(pluginFilePath, this);
        PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
        if (!interface)
        {
            pluginLoader->unload();
            pluginLoader->deleteLater();
            continue;
        }

//        interface->init(this);
        // delay load
        QTimer::singleShot(100, [=] {interface->init(this);});
    }
}

void DockPluginsController::displayModeChanged()
{
    const DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    for (auto inter : m_pluginList.keys())
        inter->displayModeChanged(displayMode);
}

void DockPluginsController::positionChanged()
{
    const Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    for (auto inter : m_pluginList.keys())
        inter->positionChanged(position);
}

bool DockPluginsController::eventFilter(QObject *o, QEvent *e)
{
    if (o != qApp)
        return false;
    if (e->type() != QEvent::DynamicPropertyChange)
        return false;

    QDynamicPropertyChangeEvent * const dpce = static_cast<QDynamicPropertyChangeEvent *>(e);
    const QString propertyName = dpce->propertyName();

    if (propertyName == PROP_POSITION)
        positionChanged();
    else if (propertyName == PROP_DISPLAY_MODE)
        displayModeChanged();

    return false;
}

PluginsItem *DockPluginsController::pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const
{
    if (!m_pluginList.contains(itemInter))
        return nullptr;

    return m_pluginList[itemInter][itemKey];
}
