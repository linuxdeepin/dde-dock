#ifndef PANELMENU_H
#define PANELMENU_H

#include <QWidget>
#include <QLabel>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QDebug>
#include "DBus/dbusmenumanager.h"
#include "DBus/dbusmenu.h"
#include "Controller/dockmodedata.h"
#include "../dockconstants.h"

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

private slots:
    void slotItemInvoked(const QString &itemId,bool result);

private:
    explicit PanelMenu(QObject *parent = 0);

    void changeToFashionMode();
    void changeToEfficientMode();
    void changeToClassicMode();
    void changeToKeepShowing();
    void changeToKeepHidden();
    void changeToSmartHide();

    QJsonObject createItemObj(const QString &itemName,OperationType type);
    QJsonObject createRadioItemObj(const QString &itemName,OperationType type,MenuGroup group,bool check);

private:
    static PanelMenu * m_panelMenu;
    DockModeData *dockCons = DockModeData::instance();
    QString m_menuInterfacePath = "";
    DBusMenuManager *m_menuManager = NULL;
    DBusDockSetting m_dockSetting;

};

#endif // PANELMENU_H
