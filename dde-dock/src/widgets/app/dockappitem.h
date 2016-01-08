#ifndef DOCKAPPITEM_H
#define DOCKAPPITEM_H

#include "dockappicon.h"
#include "dockappbg.h"
#include "../dockitem.h"
#include "interfaces/dockconstants.h"
#include "dbus/dbusdockentry.h"
#include "dbus/dbusclientmanager.h"
#include "controller/dockmodedata.h"
#include "dbus/dbusdockedappmanager.h"
#include "apppreview/apppreviewscontainer.h"

struct DockAppItemData {
    QString id;
    QString iconPath;
    QString title;
    QMap<int, QString> xidTitleMap;
    QString menuJsonString;
    bool isActived;
    bool currentOpened;
    bool isDocked;
};

class DockAppItem : public DockItem
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)
public:
    DockAppItem(QWidget *parent = 0);
    ~DockAppItem();

    DockAppItemData itemData() const;
    QWidget *getApplet();
    QString getItemId();
    QString getTitle();
    void setEntryProxyer(DBusDockEntry *entryProxyer);

    bool actived() const;
    void setActived(bool actived);

protected:
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);
    void dropEvent(QDropEvent * event);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);
    void resizeEvent(QResizeEvent *e);

private:
    void initPreviewContainer();
    void initClientManager();
    void initBackground();
    void initAppIcon();
    void initTitle();
    void initData();

    void updateIcon();
    void updateTitle();
    void updateState();
    void updateXidTitleMap();
    void updateMenuJsonString();

    void onDbusDataChanged(const QString &, const QString &);
    void onDockModeChanged(Dock::DockMode, Dock::DockMode);
    void onMousePress(QMouseEvent *event);
    void onMouseRelease(QMouseEvent *event);
    void onMouseEnter();
    void onMouseLeave();

    void resizeBackground();
    void resizeResources();
    void reanchorIcon();
    void setCurrentOpened(uint);
    void invokeMenuItem(QString id,bool);
    QString getMenuContent();

private:
    QLabel * m_appTitle;
    DockAppBG * m_appBG;
    DockAppIcon * m_appIcon;
    DockAppItemData m_itemData;
    DockModeData *m_dockModeData;
    DBusDockEntry *m_entryProxyer;
    DBusClientManager *m_clientManager;
    DBusDockedAppManager *m_appManager;
    AppPreviewsContainer *m_previewContainer;

    bool m_actived;
};

#endif // DOCKAPPITEM_H
