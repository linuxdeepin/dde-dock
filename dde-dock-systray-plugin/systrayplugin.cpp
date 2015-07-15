#include <QDBusConnection>
#include <QWidget>

#include "systrayplugin.h"
#include "compositetrayitem.h"

static const QString CompositeItemKey = "composite_item_key";

SystrayPlugin::SystrayPlugin()
{
    m_items[CompositeItemKey] = new CompositeTrayItem;
}

SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

void SystrayPlugin::init(DockPluginProxyInterface * proxy)
{
    m_proxy = proxy;
    m_mode = proxy->dockMode();

    if (!m_dbusTrayManager) {
        m_dbusTrayManager = new com::deepin::dde::TrayManager("com.deepin.dde.TrayManager",
                                                              "/com/deepin/dde/TrayManager",
                                                              QDBusConnection::sessionBus(),
                                                              this);
        connect(m_dbusTrayManager, &TrayManager::Added, this, &SystrayPlugin::onAdded);
        connect(m_dbusTrayManager, &TrayManager::Removed, this, &SystrayPlugin::onRemoved);
    }

    QList<uint> trayIcons = m_dbusTrayManager->trayIcons();
    qDebug() << "Found trayicons: " << trayIcons;

    if (m_mode == Dock::FashionMode) {
        m_proxy->itemAddedEvent(CompositeItemKey);
    }

    foreach (uint trayIcon, trayIcons) {
        onAdded(trayIcon);
    }
}

QString SystrayPlugin::name()
{
    return QString("System Tray");
}

QStringList SystrayPlugin::uuids()
{
    return m_items.keys();
}

QString SystrayPlugin::getTitle(QString)
{
    return "";
}

QWidget * SystrayPlugin::getItem(QString uuid)
{
    return m_items.value(uuid);
}

QWidget * SystrayPlugin::getApplet(QString)
{
    return NULL;
}

void SystrayPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    QWidget * widget = m_items.value(CompositeItemKey);
    CompositeTrayItem * compositeItem = qobject_cast<CompositeTrayItem*>(widget);

    if (oldMode == Dock::FashionMode && newMode != Dock::FashionMode) {

        qDebug() << "SystrayPlugin change mode to other mode.";
        foreach (QWidget * widget, m_items) {
            if (widget != compositeItem) {
                m_proxy->itemAddedEvent(m_items.key(widget));
                compositeItem->removeWidget(widget);
            }
        }

        compositeItem->setParent(NULL);
        m_proxy->itemRemovedEvent(CompositeItemKey);
    } else if (newMode == Dock::FashionMode && oldMode != Dock::FashionMode) {

        qDebug() << "SystrayPlugin change mode to fashion mode.";
        foreach (QWidget * widget, m_items) {
            if (widget != compositeItem) {
                compositeItem->addWidget(widget);
                m_proxy->itemRemovedEvent(m_items.key(widget));
            }
        }

        m_proxy->itemAddedEvent(CompositeItemKey);
        m_proxy->itemSizeChangedEvent(CompositeItemKey);
    }
}

QString SystrayPlugin::getMenuContent(QString)
{
    return "";
}

void SystrayPlugin::invokeMenuItem(QString, QString, bool)
{

}

// private methods
void SystrayPlugin::clearItems()
{
    foreach (QWidget * item, m_items) {
        item->deleteLater();
    }
    m_items.clear();
}

void SystrayPlugin::addItem(QString uuid, QWidget * item)
{
    m_items[uuid] = item;

    if (m_mode == Dock::FashionMode) {
        CompositeTrayItem * compositeItem = qobject_cast<CompositeTrayItem*>(m_items.value(CompositeItemKey));
        compositeItem->addWidget(item);

        m_proxy->itemSizeChangedEvent(CompositeItemKey);
    } else {
        m_proxy->itemAddedEvent(uuid);
    }
}

void SystrayPlugin::removeItem(QString uuid)
{
    QWidget * item = m_items[uuid];

    if (m_mode == Dock::FashionMode) {
        CompositeTrayItem * compositeItem = qobject_cast<CompositeTrayItem*>(m_items.value(CompositeItemKey));
        compositeItem->removeWidget(item);

        m_proxy->itemSizeChangedEvent(CompositeItemKey);

        m_items.remove(uuid);
        item->deleteLater();
    } else {
        m_items.remove(uuid);
        item->deleteLater();

        m_proxy->itemRemovedEvent(uuid);
    }
}

// private slots
void SystrayPlugin::onAdded(WId winId)
{
    QString key = QString::number(winId);

    QWidget *item = new QWidget;
    item->setStyleSheet("QWidget { background-color: green }");
    item->resize(Dock::APPLET_CLASSIC_ICON_SIZE, Dock::APPLET_CLASSIC_ICON_SIZE);

//    QWindow * win = QWindow::fromWinId(winId);
//    QWidget * winItem = QWidget::createWindowContainer(win, item);
//    winItem->resize(item->size());

    addItem(key, item);
}

void SystrayPlugin::onRemoved(WId winId)
{
    QString key = QString::number(winId);

    removeItem(key);
}
