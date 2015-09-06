#ifndef PANELMENU_H
#define PANELMENU_H

#include <QWidget>
#include <QLabel>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QDebug>

#include "dbus/dbusmenumanager.h"
#include "dbus/dbusmenu.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"

class PanelMenu : public QObject
{
    Q_OBJECT
public:
    enum OperationType {
        ToFashionMode,
        ToEfficientMode,
        ToClassicMode,
        ToKeepShowing,
        ToKeepHidden,
        ToSmartHide,
        ToPluginSetting
    };
    enum MenuGroup{
        DockModeGroup,
        HideModeGroup
    };

    static PanelMenu * instance();

    void showMenu(int x,int y);

signals:
    void settingPlugin();

private:
    explicit PanelMenu(QObject *parent = 0);

    void changeToFashionMode();
    void changeToEfficientMode();
    void changeToClassicMode();
    void changeToKeepShowing();
    void changeToKeepHidden();
    void changeToSmartHide();

    void onItemInvoked(const QString &itemId, bool result);

    QJsonObject createItemObj(const QString &itemName, OperationType type);
    QJsonObject createRadioItemObj(const QString &itemName, OperationType type, MenuGroup group, bool check);

private:
    static PanelMenu * m_panelMenu;
    QString m_menuInterfacePath = "";
    DBusDockSetting m_dockSetting;
    DBusMenuManager *m_menuManager = NULL;
    DockModeData *dockCons = DockModeData::instance();

};

#endif // PANELMENU_H
