#ifndef APPITEM_H
#define APPITEM_H

#include <QObject>
#include <QWidget>
#include <QPushButton>
#include <QMouseEvent>
#include <QMimeData>
#include <QDrag>
#include <QRectF>
#include <QDrag>
#include <QMimeData>
#include <QPixmap>
#include <QImage>
#include <QList>
#include <QMap>
#include <QJsonDocument>
#include <QJsonObject>
#include "DBus/dbusentryproxyer.h"
#include "DBus/dbusclientmanager.h"
#include "DBus/dbusmenu.h"
#include "DBus/dbusmenumanager.h"
#include "Controller/dockmodedata.h"
#include "abstractdockitem.h"
#include "appicon.h"
#include "appbackground.h"
#include "../dockconstants.h"
#include "apppreviews.h"

struct AppItemData {
    QString id;
    QString iconPath;
    QString title;
    QString xidsJsonString;
    QString menuJsonString;
    bool isActived;
    bool currentOpened;
    bool isDocked;
};

class AppItem : public AbstractDockItem
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)
public:
    AppItem(QWidget *parent = 0);
    ~AppItem();

    void moveWithAnimation(QPoint targetPos, int duration = 200);
    AppItemData itemData() const;
    QWidget *getApplet();
    QString getItemId();
    QString getTitle();
    void setEntryProxyer(DBusEntryProxyer *entryProxyer);

protected:
    void dragEnterEvent(QDragEnterEvent * event);
    void dragLeaveEvent(QDragLeaveEvent * event);
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);
    void mouseMoveEvent(QMouseEvent *);
    void dropEvent(QDropEvent * event);
    void enterEvent(QEvent *event);
    void leaveEvent(QEvent *event);

private:
    void initClientManager();
    void initBackground();
    void initPreview();
    void initAppIcon();
    void initTitle();
    void initMenu();
    void initData();

    void updateIcon();
    void updateTitle();
    void updateState();
    void updateXids();
    void updateMenuJsonString();

    void onDbusDataChanged(const QString &, const QString &);
    void onDockModeChanged(Dock::DockMode, Dock::DockMode);
    void onMenuItemInvoked(QString id,bool);
    void onMousePress(QMouseEvent *event);
    void onMouseRelease(QMouseEvent *event);
    void onMouseEnter();
    void onMouseLeave();

    void resizeBackground();
    void resizeResources();
    void reanchorIcon();
    void setCurrentOpened(uint);
    void setActived(bool value);
    void showMenu();

private:
    QString m_menuInterfacePath = "";
    AppItemData m_itemData;
    DockModeData *m_dockModeData = DockModeData::instance();
    DBusClientManager *m_clientmanager = NULL;
    DBusEntryProxyer *m_entryProxyer = NULL;
    DBusMenuManager *m_menuManager = NULL;
    AppBackground * m_appBackground = NULL;
    AppPreviews *m_preview = NULL;
    AppIcon * m_appIcon = NULL;
    QLabel * m_appTitle = NULL;

    const QEasingCurve MOVE_ANIMATION_CURVE = QEasingCurve::OutCubic;
};

#endif // APPITEM_H
