#include "trashplugin.h"

TrashPlugin::TrashPlugin(QObject *parent)
    : QObject(parent),
      m_trashWidget(new TrashWidget),
      m_popupApplet(new PopupControlWidget)
{
    m_popupApplet->setVisible(false);
}

const QString TrashPlugin::pluginName() const
{
    return "trash";
}

void TrashPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_proxyInter->itemAdded(this, QString());
}

QWidget *TrashPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_trashWidget;
}

QWidget *TrashPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_popupApplet;
}

const QString TrashPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString();
//    return "gvfs-open trash://";
}
