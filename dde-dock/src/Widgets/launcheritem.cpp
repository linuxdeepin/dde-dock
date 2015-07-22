#include "launcheritem.h"

LauncherItem::LauncherItem(QWidget *parent) : AbstractDockItem(parent)
{
    resize(m_dmd->getNormalItemWidth(), m_dmd->getItemHeight());
    connect(m_dmd, &DockModeData::dockModeChanged, this, &LauncherItem::changeDockMode);

    m_appIcon = new AppIcon(this);
    m_appIcon->resize(height(), height());

    m_launcherProcess = new QProcess();

    //TODO icon not show on init
    QTimer::singleShot(20, this, &LauncherItem::updateIcon);
}

void LauncherItem::mousePressEvent(QMouseEvent *)
{
    hidePreview();
}

void LauncherItem::mouseReleaseEvent(QMouseEvent *)
{
    m_launcherProcess->start("dde-launcher",QStringList());
}

void LauncherItem::enterEvent(QEvent *)
{
    showPreview();
}

void LauncherItem::leaveEvent(QEvent *)
{
    hidePreview();
}

void LauncherItem::changeDockMode(Dock::DockMode, Dock::DockMode)
{
    resize(m_dmd->getNormalItemWidth(), m_dmd->getItemHeight());
    updateIcon();
}

void LauncherItem::updateIcon()
{
    m_appIcon->setIcon("deepin-launcher");
    m_appIcon->move((width() - m_appIcon->width()) / 2, (height() - m_appIcon->height()) / 2);
}

LauncherItem::~LauncherItem()
{

}

