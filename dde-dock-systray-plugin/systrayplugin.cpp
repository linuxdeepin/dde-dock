#include <QDBusConnection>
#include <QWidget>

#include "systrayplugin.h"
#include "compositetrayitem.h"
#include "trayicon.h"
#include "../dde-dock/src/dbus/dbusentrymanager.h"

static const QString CompositeItemKey = "composite_item_key";

SystrayPlugin::SystrayPlugin()
{
    m_compositeItem = new CompositeTrayItem;

    connect(m_compositeItem, &CompositeTrayItem::sizeChanged,
            this, &SystrayPlugin::onCompositeItemSizeChanged);
}

SystrayPlugin::~SystrayPlugin()
{
    m_compositeItem->deleteLater();
}

void SystrayPlugin::init(DockPluginProxyInterface * proxy)
{
    m_proxy = proxy;
    m_compositeItem->setMode(proxy->dockMode());

    if (!m_dbusTrayManager) {
        m_dbusTrayManager = new com::deepin::dde::TrayManager("com.deepin.dde.TrayManager",
                                                              "/com/deepin/dde/TrayManager",
                                                              QDBusConnection::sessionBus(),
                                                              this);
        connect(m_dbusTrayManager, &TrayManager::TrayIconsChanged, this, &SystrayPlugin::onTrayIconsChanged);
        connect(m_dbusTrayManager, &TrayManager::Changed, m_compositeItem, &CompositeTrayItem::handleTrayiconDamage);
    }

    DBusEntryManager *entryManager = new DBusEntryManager(this);
    connect(entryManager, &DBusEntryManager::TrayInited, this, &SystrayPlugin::onTrayInit);

    initTrayIcons();

    if (m_compositeItem->parentWidget()) {
        //wait for parentWidget() is valuable to set eventFilter
        m_compositeItem->parentWidget()->installEventFilter(m_compositeItem);
    }
}

QString SystrayPlugin::getPluginName()
{
    return "System Tray";
}

QStringList SystrayPlugin::ids()
{
    return QStringList(CompositeItemKey);
}

QString SystrayPlugin::getTitle(QString)
{
    return "";
}

QString SystrayPlugin::getName(QString)
{
    return getPluginName();
}

QString SystrayPlugin::getCommand(QString)
{
    return "";
}

bool SystrayPlugin::configurable(const QString &)
{
    return false;
}

bool SystrayPlugin::enabled(const QString &)
{
    return true;
}

void SystrayPlugin::setEnabled(const QString &, bool)
{

}

QWidget * SystrayPlugin::getItem(QString)
{
    return m_compositeItem;
}

QWidget * SystrayPlugin::getApplet(QString)
{
    return NULL;
}

void SystrayPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode)
{
    m_compositeItem->setMode(newMode);
    m_proxy->infoChangedEvent(DockPluginInterface::InfoTypeItemSize, CompositeItemKey);
}

QString SystrayPlugin::getMenuContent(QString)
{
    return "";
}

void SystrayPlugin::invokeMenuItem(QString, QString, bool)
{

}

void SystrayPlugin::initTrayIcons()
{
    m_compositeItem->clear();

    m_dbusTrayManager->RetryManager();
    QList<uint> trayIcons = m_dbusTrayManager->trayIcons();
    qDebug() << "Init trayicons, Found trayicons: " <<m_dbusTrayManager->isValid() << trayIcons << m_dbusTrayManager->property("TrayIcons");

    foreach (uint trayIcon, trayIcons) {
        addTrayIcon(trayIcon);
    }

    m_proxy->itemAddedEvent(CompositeItemKey);
}

// private slots
void SystrayPlugin::addTrayIcon(WId winId)
{
    QString key = QString::number(winId);
    if (m_compositeItem->exist(key))
        return;
    qWarning() << "Systray add:" << winId;

    TrayIcon * icon = new TrayIcon(winId);

    m_compositeItem->addTrayIcon(key, icon);

    m_proxy->infoChangedEvent(DockPluginInterface::InfoTypeItemSize, CompositeItemKey);
}

void SystrayPlugin::removeTrayIcon(WId winId)
{
    qWarning() << "Systray remove:" << winId;
    QString key = QString::number(winId);

    m_compositeItem->remove(key);

    m_proxy->infoChangedEvent(DockPluginInterface::InfoTypeItemSize, CompositeItemKey);
}

void SystrayPlugin::onTrayIconsChanged()
{
    QList<uint> icons = m_dbusTrayManager->trayIcons();
    QStringList ids = m_compositeItem->trayIds();
    qDebug() << "TrayIconsChanged:" << icons;
    foreach (uint id, icons) {   //add news
        if (ids.indexOf(QString::number(id)) == -1) {
            addTrayIcon(id);
        }
    }
    foreach (QString id, ids) { //remove olds
        if (icons.indexOf(id.toUInt()) == -1) {
            removeTrayIcon(id.toUInt());
        }
    }
}

void SystrayPlugin::onTrayInit()
{

}

void SystrayPlugin::onCompositeItemSizeChanged()
{
    m_proxy->infoChangedEvent(DockPluginInterface::InfoTypeItemSize, CompositeItemKey);
}
