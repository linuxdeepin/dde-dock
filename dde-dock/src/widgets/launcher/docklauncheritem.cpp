/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QTimer>
#include <QProcess>
#include "docklauncheritem.h"
#include "controller/signalmanager.h"

DockLauncherItem::DockLauncherItem(QWidget *parent)
    : DockItem(parent),
      m_launcherInter(new DBusLauncherController(this))
{
    setFixedSize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockLauncherItem::changeDockMode);

    m_appIcon = new DockAppIcon(this);
    m_appIcon->resize(height(), height());
//    connect(m_appIcon, &DockAppIcon::mousePress, this, &DockLauncherItem::slotMousePress);
//    connect(m_appIcon, &DockAppIcon::mouseRelease, this, &DockLauncherItem::slotMouseRelease);
    connect(m_appIcon, &DockAppIcon::mouseRelease, this, &DockLauncherItem::startupLauncher);
    connect(this, &DockLauncherItem::mouseRelease, this, &DockLauncherItem::startupLauncher);

    //TODO icon not show on init
    QTimer::singleShot(20, this, SLOT(updateIcon()));
    connect(SignalManager::instance(), &SignalManager::requestAppIconUpdate, this, &DockLauncherItem::updateIcon);
}

void DockLauncherItem::enterEvent(QEvent *)
{
    if (hoverable()) {
        showPreview();
        emit mouseEnter();
    }
}

void DockLauncherItem::leaveEvent(QEvent *)
{
    hidePreview();
    emit mouseLeave();
}

void DockLauncherItem::mousePressEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMousePress(event);
    else
        DockItem::mousePressEvent(event);

    emit mousePress();
}

void DockLauncherItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        slotMouseRelease(event);
    else
        DockItem::mouseReleaseEvent(event);

    emit mouseRelease();
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

//    startupLauncher();
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

void DockLauncherItem::startupLauncher()
{
    if (m_launcherInter->isValid())
    {
        m_launcherInter->Toggle();
        return;
    }

    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    QStringList args = QStringList() << "--print-reply"
                                     << "--dest=com.deepin.dde.Launcher"
                                     << "/com/deepin/dde/Launcher"
                                     << "com.deepin.dde.Launcher.Toggle";

    proc->start("dbus-send", args);
}

DockLauncherItem::~DockLauncherItem()
{

}

