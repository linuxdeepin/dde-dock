#ifndef APPITEM_H
#define APPITEM_H

#include <QMap>
#include <QList>
#include <QDrag>
#include <QRectF>
#include <QImage>
#include <QPixmap>
#include <QWidget>
#include <QObject>
#include <QMimeData>
#include <QPushButton>
#include <QJsonObject>
#include <QMouseEvent>
#include <QJsonDocument>

#include "appicon.h"
#include "apppreviews.h"
#include "appbackground.h"
#include "abstractdockitem.h"
#include "dbus/dbusentryproxyer.h"
#include "dbus/dbusclientmanager.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"

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
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    void initClientManager();
    void initBackground();
    void initPreview();
    void initAppIcon();
    void initTitle();
    void initData();

    void updateIcon();
    void updateTitle();
    void updateState();
    void updateXids();
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
    void setActived(bool value);
    void invokeMenuItem(QString id,bool);
    QString getMenuContent();

private:
    AppItemData m_itemData;
    DockModeData *m_dockModeData = DockModeData::instance();
    DBusClientManager *m_clientmanager = NULL;
    DBusEntryProxyer *m_entryProxyer = NULL;
    AppBackground * m_appBackground = NULL;
    AppPreviews *m_preview = NULL;
    AppIcon * m_appIcon = NULL;
    QLabel * m_appTitle = NULL;

    const QEasingCurve MOVE_ANIMATION_CURVE = QEasingCurve::OutCubic;
};

#endif // APPITEM_H
