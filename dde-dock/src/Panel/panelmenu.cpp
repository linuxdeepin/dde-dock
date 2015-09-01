#include "panelmenu.h"


PanelMenu * PanelMenu::m_panelMenu = NULL;
PanelMenu * PanelMenu::instance()
{
    if (!m_panelMenu)
        m_panelMenu = new PanelMenu();
    return m_panelMenu;
}

PanelMenu::PanelMenu(QObject *parent) : QObject(parent)
{
    m_menuManager = new DBusMenuManager(this);
}

void PanelMenu::showMenu(int x, int y)
{
    if (m_menuManager && m_menuManager->isValid()){
        QDBusPendingReply<QDBusObjectPath> pr = m_menuManager->RegisterMenu();
        if (pr.count() == 1)
        {
            QDBusObjectPath op = pr.argumentAt(0).value<QDBusObjectPath>();
            m_menuInterfacePath = op.path();
            DBusMenu *m_menu = new DBusMenu(m_menuInterfacePath,this);
            connect(m_menu,&DBusMenu::MenuUnregistered,m_menu,&DBusMenu::deleteLater);
            connect(m_menu,&DBusMenu::ItemInvoked,this,&PanelMenu::slotItemInvoked);

            QJsonObject targetObj;
            targetObj.insert("x",QJsonValue(x));
            targetObj.insert("y",QJsonValue(y));
            targetObj.insert("isDockMenu",QJsonValue(false));

            QJsonArray contentArry;
            contentArry.append(createRadioItemObj("Fashion mode",ToFashionMode,DockModeGroup,dockCons->getDockMode() == Dock::FashionMode));
            contentArry.append(createRadioItemObj("Efficient mode",ToEfficientMode,DockModeGroup,dockCons->getDockMode() == Dock::EfficientMode));
            contentArry.append(createRadioItemObj("Classic mode",ToClassicMode,DockModeGroup,dockCons->getDockMode() == Dock::ClassicMode));
            contentArry.append(createItemObj("",OperationType(-1)));
            contentArry.append(createRadioItemObj("Keep showing",ToKeepShowing,HideModeGroup,dockCons->getHideMode() == Dock::KeepShowing));
            contentArry.append(createRadioItemObj("Keep hidden",ToKeepHidden,HideModeGroup,dockCons->getHideMode() == Dock::KeepHidden));
            contentArry.append(createRadioItemObj("Smart hide",ToSmartHide,HideModeGroup,dockCons->getHideMode() == Dock::SmartHide));
            contentArry.append(createItemObj("",OperationType(-1)));
            contentArry.append(createItemObj("Notification area setting",ToPluginSetting));

            QJsonObject contentObj;
            contentObj.insert("items",contentArry);
            targetObj.insert("menuJsonContent",QString(QJsonDocument(contentObj).toJson()));

            m_menu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));
        }
    }
}

void PanelMenu::slotItemInvoked(const QString &itemId, bool result)
{
    if (itemId.split(":").length() < 1)
        return;

    OperationType tt = OperationType(itemId.split(":").at(0).toInt());
    switch (tt)
    {
    case ToFashionMode:
        changeToFashionMode();
        break;
    case ToEfficientMode:
        changeToEfficientMode();
        break;
    case ToClassicMode:
        changeToClassicMode();
        break;
    case ToKeepShowing:
        changeToKeepShowing();
        break;
    case ToKeepHidden:
        changeToKeepHidden();
        break;
    case ToSmartHide:
        changeToSmartHide();
        break;
    case ToPluginSetting:
        emit settingPlugin();
        break;
    default:
        break;
    }

}

void PanelMenu::changeToFashionMode()
{
    qWarning() << "Change to fashion mode...";
    dockCons->setDockMode(Dock::FashionMode);
}

void PanelMenu::changeToEfficientMode()
{
    qWarning() << "Change to efficient mode...";
    dockCons->setDockMode(Dock::EfficientMode);
}

void PanelMenu::changeToClassicMode()
{
    qWarning() << "Change to classic mode...";
    dockCons->setDockMode(Dock::ClassicMode);
}

void PanelMenu::changeToKeepShowing()
{
    qWarning() << "Change to keep showing mode...";
    dockCons->setHideMode(Dock::KeepShowing);
}

void PanelMenu::changeToKeepHidden()
{
    qWarning() << "Change to keep hidden mode...";
    dockCons->setHideMode(Dock::KeepHidden);
}

void PanelMenu::changeToSmartHide()
{
    qWarning() << "Change to smart hide mode...";
    dockCons->setHideMode(Dock::SmartHide);
}

QJsonObject PanelMenu::createItemObj(const QString &itemName, OperationType type)
{
    QJsonObject itemObj;
    itemObj.insert("itemId",QString::number(type));
    itemObj.insert("itemText",itemName);
    itemObj.insert("itemIcon","");
    itemObj.insert("itemIconHover","");
    itemObj.insert("itemIconInactive","");
    itemObj.insert("itemExtra","");
    itemObj.insert("isActive",true);
    itemObj.insert("checked",false);
    itemObj.insert("itemSubMenu",QJsonObject());

    return itemObj;
}

QJsonObject PanelMenu::createRadioItemObj(const QString &itemName, OperationType type, MenuGroup group, bool check)
{
    QJsonObject itemObj;
    itemObj.insert("itemId",QString::number(type) + ":radio:" + QString::number(group));
    itemObj.insert("itemText",itemName);
    itemObj.insert("itemIcon","");
    itemObj.insert("itemIconHover","");
    itemObj.insert("itemIconInactive","");
    itemObj.insert("itemExtra","");
    itemObj.insert("isActive",true);
    itemObj.insert("checked",check);
    itemObj.insert("itemSubMenu",QJsonObject());

    return itemObj;
}
