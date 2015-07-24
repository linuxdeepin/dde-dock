#include "trashplugin.h"


TrashPlugin::TrashPlugin()
{

}


void TrashPlugin::init(DockPluginProxyInterface *proxy)
{
    m_proxy = proxy;

    setMode(proxy->dockMode());
}

QString TrashPlugin::name()
{
    return "Trash plugin";
}

QStringList TrashPlugin::uuids()
{
    return QStringList(m_uuid);
}

QString TrashPlugin::getTitle(QString)
{
    return name();
}

QWidget * TrashPlugin::getItem(QString)
{
    MainItem * item = new MainItem();
    connect(this, &TrashPlugin::menuItemInvoked, item, &MainItem::emptyTrash);

    m_itemList.append(item);

    return item;
}

QWidget * TrashPlugin::getApplet(QString)
{
    return NULL;
}

void TrashPlugin::changeMode(Dock::DockMode newMode,
                                Dock::DockMode)
{
    setMode(newMode);
}

QString TrashPlugin::getMenuContent(QString)
{
    QJsonObject contentObj;

    QJsonArray items;

    items.append(createMenuItem("clear_trash", "Clear Trash"));

    contentObj.insert("items", items);

    return QString(QJsonDocument(contentObj).toJson());
}

void TrashPlugin::invokeMenuItem(QString, QString itemId, bool checked)
{
    qWarning() << "Menu check:" << itemId << checked;
    emit menuItemInvoked();
}

// private methods
void TrashPlugin::setMode(Dock::DockMode mode)
{
    m_mode = mode;

    if (mode == Dock::FashionMode)
        m_proxy->itemAddedEvent(m_uuid);
    else
    {
        m_itemList.clear();
        m_proxy->itemRemovedEvent(m_uuid);
    }
}

QJsonObject TrashPlugin::createMenuItem(QString itemId, QString itemName, bool checkable, bool checked)
{
    QJsonObject itemObj;

    itemObj.insert("itemId", itemId);
    itemObj.insert("itemText", itemName);
    itemObj.insert("itemIcon", "");
    itemObj.insert("itemIconHover", "");
    itemObj.insert("itemIconInactive", "");
    itemObj.insert("itemExtra", "");
    itemObj.insert("isActive", true);
    itemObj.insert("isCheckable", checkable);
    itemObj.insert("checked", checked);
    itemObj.insert("itemSubMenu", QJsonObject());

    return itemObj;
}


TrashPlugin::~TrashPlugin()
{

}

#if QT_VERSION < 0x050000
Q_EXPORT_PLUGIN2(dde-dock-trash-plugin, TrashPlugin)
#endif // QT_VERSION < 0x050000
