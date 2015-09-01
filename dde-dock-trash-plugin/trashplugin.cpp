#include <QIcon>

#include "trashplugin.h"

TrashPlugin::TrashPlugin()
{
    QIcon::setThemeName("Deepin");

    m_item = new MainItem();
    connect(this, &TrashPlugin::menuItemInvoked, m_item, &MainItem::emptyTrash);
}

void TrashPlugin::init(DockPluginProxyInterface *proxy)
{
    m_proxy = proxy;

    setMode(proxy->dockMode());
}

QString TrashPlugin::getPluginName()
{
    return "Trash plugin";
}

QStringList TrashPlugin::ids()
{
    return QStringList(m_id);
}

QString TrashPlugin::getName(QString)
{
    return getPluginName();
}

QString TrashPlugin::getTitle(QString)
{
    return getPluginName();
}

QString TrashPlugin::getCommand(QString)
{
    return "";
}

bool TrashPlugin::canDisable(QString)
{
    return false;
}

bool TrashPlugin::isDisabled(QString)
{
    return false;
}

void TrashPlugin::setDisabled(QString, bool)
{

}

QWidget * TrashPlugin::getItem(QString)
{
    return m_item;
}

QWidget * TrashPlugin::getApplet(QString)
{
    return NULL;
}

void TrashPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode != oldMode)
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
        m_proxy->itemAddedEvent(m_id);
    else{
        m_proxy->itemRemovedEvent(m_id);
        m_item->setParent(NULL);
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
