#include <QDBusConnection>
#include <QWidget>

#include "systrayplugin.h"
#include "compositetrayitem.h"

static const QString CompositeItemKey = "composite_item_key";

SystrayPlugin::~SystrayPlugin()
{
    this->clearItems();
}

void SystrayPlugin::init(DockPluginProxyInterface * proxier)
{
    m_proxier = proxier;

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

    foreach (uint trayIcon, trayIcons) {
        onAdded(trayIcon);
    }
}

QStringList SystrayPlugin::uuids()
{
    return m_items.keys();
}

QWidget * SystrayPlugin::getItem(QString uuid)
{
    return m_items.value(uuid);
}

QString SystrayPlugin::name()
{
    return QString("System Tray");
}

void SystrayPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    m_mode = newMode;

    CompositeTrayItem * compositeItem = NULL;

    QWidget * widget = m_items.value(CompositeItemKey);
    if (!widget)
    {
        compositeItem = new CompositeTrayItem;
        m_items[CompositeItemKey] = compositeItem;
    }
    else
        compositeItem = qobject_cast<CompositeTrayItem*>(widget);

    compositeItem->resize(Dock::APPLET_FASHION_ITEM_WIDTH,Dock::APPLET_FASHION_ITEM_HEIGHT);

    if (oldMode == Dock::FashionMode && newMode != Dock::FashionMode) {

        compositeItem->setParent(NULL);

        qDebug() << "SystrayPlugin change mode to other mode.";
        foreach (QWidget * widget, m_items) {
            if (widget != compositeItem) {
                compositeItem->removeWidget(widget);
                m_proxier->itemAddedEvent(m_items.key(widget));
            }
        }

        m_proxier->itemRemovedEvent(CompositeItemKey);
    } else if (newMode == Dock::FashionMode && oldMode != Dock::FashionMode) {

        qDebug() << "SystrayPlugin change mode to fashion mode.";
        foreach (QWidget * widget, m_items) {
            if (widget != compositeItem) {
                compositeItem->addWidget(widget);
                m_proxier->itemRemovedEvent(m_items.key(widget));
            }
        }

        m_proxier->itemAddedEvent(CompositeItemKey);
    }
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
    } else {
        m_proxier->itemAddedEvent(uuid);
    }
}

void SystrayPlugin::removeItem(QString uuid)
{
     QWidget * item = m_items[uuid];

    if (m_mode == Dock::FashionMode) {
        CompositeTrayItem * compositeItem = qobject_cast<CompositeTrayItem*>(m_items.value(CompositeItemKey));
        compositeItem->removeWidget(item);
    }

    m_items.remove(uuid);
    item->deleteLater();

    m_proxier->itemRemovedEvent(uuid);
}

// private slots
void SystrayPlugin::onAdded(WId winId)
{
    QString key = QString::number(winId);

    QWidget *item = new QWidget;
    item->setFixedSize(16, 16);

    QWindow *win = QWindow::fromWinId(winId);
    QWidget *w = QWidget::createWindowContainer(win, item);
    w->setFixedSize(item->size());

    addItem(key, item);
}

void SystrayPlugin::onRemoved(WId winId)
{
    QString key = QString::number(winId);

    removeItem(key);
}
