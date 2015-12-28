#ifndef DOCKPLUGINITEM_H
#define DOCKPLUGINITEM_H

#include <QJsonObject>

#include "widgets/dockitem.h"
#include "interfaces/dockplugininterface.h"
#include "dbus/dbusdisplay.h"

class QMouseEvent;
class DockPluginItem : public DockItem
{
    Q_OBJECT
public:
    DockPluginItem(DockPluginInterface *plugin, QString id, QWidget * parent = 0);
    virtual ~DockPluginItem();

    QString id() const;

    QString getTitle() Q_DECL_OVERRIDE;
    QWidget * getApplet() Q_DECL_OVERRIDE;

    QString getMenuContent() Q_DECL_OVERRIDE;
    void invokeMenuItem(QString itemId, bool checked) Q_DECL_OVERRIDE;

protected:
    void enterEvent(QEvent * event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent * event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent * event) Q_DECL_OVERRIDE;
//    void mouseReleaseEvent(QMouseEvent * event) Q_DECL_OVERRIDE;

private:
    DBusDisplay *m_display = NULL;
    QWidget *m_pluginItemContents = NULL;
    DockPluginInterface * m_plugin;
    QString m_id;

    QJsonObject createMenuItem(QString itemId, QString itemName, bool checkable, bool checked);
    const int DOCK_PREVIEW_MARGIN = 7;
};

#endif // DOCKPLUGINITEM_H
