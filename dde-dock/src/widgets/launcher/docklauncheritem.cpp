#include <QTimer>
#include <QProcess>
#include "docklauncheritem.h"
#include "controller/signalmanager.h"

DockLauncherItem::DockLauncherItem(QWidget *parent)
    : DockItem(parent)
{
    setFixedSize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockLauncherItem::changeDockMode);

    m_appIcon = new DockAppIcon(this);
    m_appIcon->resize(height(), height());
    connect(m_appIcon, &DockAppIcon::mousePress, this, &DockLauncherItem::slotMousePress);
    connect(m_appIcon, &DockAppIcon::mouseRelease, this, &DockLauncherItem::slotMouseRelease);

    m_launcherProcess = new QProcess();

    //TODO icon not show on init
    QTimer::singleShot(20, this, SLOT(updateIcon()));
    connect(SignalManager::instance(), &SignalManager::requestAppIconUpdate, this, &DockLauncherItem::updateIcon);
}

void DockLauncherItem::enterEvent(QEvent *)
{
    if (!hoverable())
        return;

    showPreview();
}

void DockLauncherItem::leaveEvent(QEvent *)
{

    hidePreview();
}

void DockLauncherItem::mousePressEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMousePress(event);
    else
        DockItem::mousePressEvent(event);
}

void DockLauncherItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMouseRelease(event);
    else
        DockItem::mouseReleaseEvent(event);
}

void DockLauncherItem::slotMousePress(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    hidePreview();
}

void DockLauncherItem::slotMouseRelease(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    m_launcherProcess->startDetached("dde-launcher",QStringList());
}

void DockLauncherItem::changeDockMode(Dock::DockMode, Dock::DockMode)
{
    setFixedSize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    updateIcon();
}

void DockLauncherItem::updateIcon()
{
    m_appIcon->setIcon("deepin-launcher");
    m_appIcon->resize(m_dockModeData->getAppIconSize(), m_dockModeData->getAppIconSize());
    reanchorIcon();
}

void DockLauncherItem::reanchorIcon()
{
    switch (m_dockModeData->getDockMode()) {
    case Dock::FashionMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, 0);
        break;
    case Dock::EfficientMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, (height() - m_appIcon->height()) / 2);
        break;
    case Dock::ClassicMode:
        m_appIcon->move((height() - m_appIcon->height()) / 2, (height() - m_appIcon->height()) / 2);
    default:
        break;
    }
}

DockLauncherItem::~DockLauncherItem()
{

}

