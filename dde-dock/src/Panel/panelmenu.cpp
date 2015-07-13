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
            contentArry.append(createItemObj("Fashion Mode",ToFashionMode));
            contentArry.append(createItemObj("Efficient Mode",ToEfficientMode));
            contentArry.append(createItemObj("Classic Mode",ToClassicMode));

            QJsonObject contentObj;
            contentObj.insert("items",contentArry);
            targetObj.insert("menuJsonContent",QString(QJsonDocument(contentObj).toJson()));

            m_menu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));
        }
    }
}

void PanelMenu::slotItemInvoked(const QString &itemId, bool result)
{
    OperationType tt = OperationType(itemId.toInt());
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
    default:
        break;
    }

    qWarning() << itemId << result << tt;
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
