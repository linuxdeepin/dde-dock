#include "launcheritem.h"
#include "controller/signalmanager.h"

LauncherItem::LauncherItem(QWidget *parent) : AbstractDockItem(parent)
{
    resize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &LauncherItem::changeDockMode);

    m_appIcon = new AppIcon(this);
    m_appIcon->resize(height(), height());
    connect(m_appIcon, &AppIcon::mousePress, this, &LauncherItem::slotMousePress);
    connect(m_appIcon, &AppIcon::mouseRelease, this, &LauncherItem::slotMouseRelease);

    m_launcherProcess = new QProcess();

    //TODO icon not show on init
    QTimer::singleShot(20, this, SLOT(updateIcon()));
    connect(SignalManager::instance(), &SignalManager::requestAppIconUpdate, this, &LauncherItem::updateIcon);
}

void LauncherItem::enterEvent(QEvent *)
{
    if (m_dockModeData->getHideState() != Dock::HideStateShown)
        return;

    emit mouseEntered();

    showPreview();
}

void LauncherItem::leaveEvent(QEvent *)
{
    emit mouseExited();

    hidePreview();
}

void LauncherItem::mousePressEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMousePress(event);
    else
        AbstractDockItem::mousePressEvent(event);
}

void LauncherItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMouseRelease(event);
    else
        AbstractDockItem::mouseReleaseEvent(event);
}

void LauncherItem::slotMousePress(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    emit mousePress(event);

    hidePreview();
}

void LauncherItem::slotMouseRelease(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    emit mouseRelease(event);

    m_launcherProcess->startDetached("dde-launcher",QStringList());
}

void LauncherItem::changeDockMode(Dock::DockMode, Dock::DockMode)
{
    resize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    updateIcon();
}

void LauncherItem::updateIcon()
{
    m_appIcon->setIcon("deepin-launcher");
    m_appIcon->resize(m_dockModeData->getAppIconSize(), m_dockModeData->getAppIconSize());
    reanchorIcon();
}

void LauncherItem::reanchorIcon()
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

LauncherItem::~LauncherItem()
{

}

