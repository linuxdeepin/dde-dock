#ifndef TRASHPLUGIN_H
#define TRASHPLUGIN_H

#include <QMap>
#include <QLabel>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QIcon>
#include <QDebug>
#include "mainitem.h"
#include "dock/dockconstants.h"
#include "dock/dockplugininterface.h"
#include "dock/dockpluginproxyinterface.h"

class TrashPlugin : public QObject, public DockPluginInterface
{
    Q_OBJECT
#if QT_VERSION >= 0x050000
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "dde-dock-trash-plugin.json")
#endif // QT_VERSION >= 0x050000
    Q_INTERFACES(DockPluginInterface)

public:
    TrashPlugin();
    ~TrashPlugin();

    void init(DockPluginProxyInterface *proxy) Q_DECL_OVERRIDE;

    QString name() Q_DECL_OVERRIDE;

    QStringList uuids() Q_DECL_OVERRIDE;
    QString getTitle(QString uuid) Q_DECL_OVERRIDE;
    QWidget * getItem(QString uuid) Q_DECL_OVERRIDE;
    QWidget * getApplet(QString uuid) Q_DECL_OVERRIDE;
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) Q_DECL_OVERRIDE;

    QString getMenuContent(QString uuid) Q_DECL_OVERRIDE;
    void invokeMenuItem(QString uuid, QString itemId, bool checked) Q_DECL_OVERRIDE;

signals:
    void menuItemInvoked();

private:
    QList<MainItem *> m_itemList;
    QString m_uuid = "trash_plugin";
    DockPluginProxyInterface * m_proxy;

    Dock::DockMode m_mode = Dock::EfficientMode;

private:
    void setMode(Dock::DockMode mode);
    QJsonObject createMenuItem(QString itemId,
                               QString itemName,
                               bool checkable = false,
                               bool checked = false);
};

#endif // TRASHPLUGIN_H
