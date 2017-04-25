#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dockitemcontroller.h"
#include "dockpluginloader.h"

#include <QDebug>
#include <QDir>

DockPluginsController::DockPluginsController(DockItemController *itemControllerInter)
    : QObject(itemControllerInter),
      m_itemControllerInter(itemControllerInter)
{
    qApp->installEventFilter(this);

    QTimer::singleShot(1, this, &DockPluginsController::startLoader);
}

DockPluginsController::~DockPluginsController()
{
}

void DockPluginsController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // check if same item added
    if (m_pluginList.contains(itemInter))
        if (m_pluginList[itemInter].contains(itemKey))
            return;

    PluginsItem *item = new PluginsItem(itemInter, itemKey);
    item->setVisible(false);

    m_pluginList[itemInter][itemKey] = item;

    emit pluginItemInserted(item);
}

void DockPluginsController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    Q_ASSERT(item);

    item->update();

    emit pluginItemUpdated(item);
}

void DockPluginsController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(item);

    m_pluginList[itemInter].remove(itemKey);
    QTimer::singleShot(1, item, &PluginsItem::deleteLater);
}

//void DockPluginsController::requestRefershWindowVisible()
//{
//    for (auto list : m_pluginList.values())
//    {
//        for (auto item : list.values())
//        {
//            Q_ASSERT(item);
//            emit item->requestRefershWindowVisible();
//            return;
//        }
//    }
//}

void DockPluginsController::requestContextMenu(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    item->showContextMenu();
}

//void DockPluginsController::requestPopupApplet(PluginsItemInterface * const itemInter, const QString &itemKey)
//{
//    PluginsItem *item = pluginItemAt(itemInter, itemKey);

//    Q_ASSERT(item);
//    item->showPopupApplet();
//}

void DockPluginsController::startLoader()
{
    DockPluginLoader *loader = new DockPluginLoader(this);

    connect(loader, &DockPluginLoader::finished, loader, &DockPluginLoader::deleteLater, Qt::QueuedConnection);
    connect(loader, &DockPluginLoader::pluginFounded, this, &DockPluginsController::loadPlugin, Qt::QueuedConnection);

    QTimer::singleShot(1, loader, [=] { loader->start(QThread::LowestPriority); });
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

void DockPluginsController::loadPlugin(const QString &pluginFile)
{
    qDebug() << "load plugin: " << pluginFile;

    QPluginLoader *pluginLoader = new QPluginLoader(pluginFile, this);
    PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
    if (!interface)
    {
        qWarning() << "load plugin failed!!!" << pluginLoader->errorString() << pluginFile;
        pluginLoader->unload();
        pluginLoader->deleteLater();
        return;
    }

    m_pluginList.insert(interface, QMap<QString, PluginsItem *>());
    qDebug() << "init plugin: " << interface->pluginName();
    interface->init(this);
    qDebug() << "init plugin finished: " << interface->pluginName();
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
