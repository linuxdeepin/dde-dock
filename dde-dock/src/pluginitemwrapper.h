#ifndef PLUGINITEMWRAPPER_H
#define PLUGINITEMWRAPPER_H

#include "abstractdockitem.h"
#include "dockplugininterface.h"

class QMouseEvent;
class PluginItemWrapper : public AbstractDockItem
{
    Q_OBJECT
public:
    PluginItemWrapper(DockPluginInterface *plugin, QString uuid, QWidget * parent = 0);
    virtual ~PluginItemWrapper();

    QString uuid() const;

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
    QWidget *m_pluginItemContents = NULL;
    DockPluginInterface * m_plugin;
    QString m_uuid;
};

#endif // PLUGINITEMWRAPPER_H
