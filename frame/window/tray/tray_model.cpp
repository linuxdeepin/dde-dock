/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
#include "tray_model.h"
#include "tray_monitor.h"

#include "indicatortrayitem.h"
#include "indicatorplugin.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"

#include <QMimeData>
#include <QIcon>
#include <QDebug>
#include <QAbstractItemModel>
#include <QDBusInterface>

#define TRAY_DRAG_FALG "tray_drag"

TrayModel::TrayModel(QListView *view, bool isIconTray, bool hasInputMethod, QObject *parent)
    : QAbstractListModel(parent)
    , m_dragModelIndex(QModelIndex())
    , m_dropModelIndex(QModelIndex())
    , m_view(view)
    , m_monitor(new TrayMonitor(this))
    , m_isTrayIcon(isIconTray)
    , m_hasInputMethod(hasInputMethod)
{
    Q_ASSERT(m_view);

    if (isIconTray) {
        connect(m_monitor, &TrayMonitor::xEmbedTrayAdded, this, &TrayModel::onXEmbedTrayAdded);
        connect(m_monitor, &TrayMonitor::indicatorFounded, this, &TrayModel::onIndicatorFounded);
        connect(QuickSettingController::instance(), &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute &pluginAttr) {
            if (pluginAttr != QuickSettingController::PluginAttribute::System)
                return;

            systemItemAdded(itemInter);
        });

        connect(QuickSettingController::instance(), &QuickSettingController::pluginRemoved, this, &TrayModel::onSystemItemRemoved);
        QMetaObject::invokeMethod(this, [ = ] {
            QList<PluginsItemInterface *> systemPlugins = QuickSettingController::instance()->pluginItems(QuickSettingController::PluginAttribute::System);
            for (PluginsItemInterface *plugin : systemPlugins)
                systemItemAdded(plugin);
        }, Qt::QueuedConnection);
    }
    connect(m_monitor, &TrayMonitor::xEmbedTrayRemoved, this, &TrayModel::onXEmbedTrayRemoved);
    connect(m_monitor, &TrayMonitor::requestUpdateIcon, this, &TrayModel::requestUpdateIcon);
    connect(m_monitor, &TrayMonitor::sniTrayAdded, this, &TrayModel::onSniTrayAdded);
    connect(m_monitor, &TrayMonitor::sniTrayRemoved, this, &TrayModel::onSniTrayRemoved);
}

void TrayModel::dropSwap(int newPos)
{
    if (!m_dragModelIndex.isValid())
        return;

    removeRows(m_dragModelIndex.row(), 1, QModelIndex());
    dropInsert(newPos);

    emit QAbstractItemModel::dataChanged(m_dragModelIndex, m_dropModelIndex);
}

void TrayModel::dropInsert(int newPos)
{
    beginInsertRows(QModelIndex(), newPos, newPos);
    WinInfo name = m_dragInfo;
    m_winInfos.insert(newPos, name);
    // 更新输入法的位置
    endInsertRows();
}

void TrayModel::clearDragDropIndex()
{
    const QModelIndex startIndex = m_dragModelIndex;
    const QModelIndex endIndex = m_dropModelIndex;

    m_dragModelIndex = m_dropModelIndex = QModelIndex();

    emit QAbstractItemModel::dataChanged(startIndex, endIndex);
    emit QAbstractItemModel::dataChanged(endIndex, startIndex);
}

void TrayModel::setDragingIndex(const QModelIndex index)
{
    m_dragModelIndex = index;
    m_dropModelIndex = index;

    emit QAbstractListModel::dataChanged(index, index);
}

void TrayModel::setDragDropIndex(const QModelIndex index)
{
    if (m_dragModelIndex == index)
        return;

    m_dropModelIndex = index;

    emit QAbstractListModel::dataChanged(m_dragModelIndex, index);
    emit QAbstractListModel::dataChanged(index, m_dragModelIndex);
}

void TrayModel::setDragKey(const QString &key)
{
    m_dragKey = key;
}

bool TrayModel::indexDragging(const QModelIndex &index) const
{
    if (index.isValid() && index.data(Role::KeyRole).toString() == m_dragKey)
        return true;

    if (!m_dragModelIndex.isValid() || !m_dropModelIndex.isValid())
        return false;

    const int start = m_dragModelIndex.row();
    const int end = m_dropModelIndex.row();
    const int current = index.row();

    return (start <= end && current >= start && current <= end)
            || (start >= end && current <= start && current >= end);
}

IndicatorTrayItem *TrayModel::indicatorWidget(const QString &indicatorName) const
{
    QString indicatorKey = indicatorName;
    indicatorKey = indicatorKey.remove(0, QString("indicator:").length());
    if (m_indicatorMap.contains(indicatorKey))
        return m_indicatorMap.value(indicatorKey)->widget();

    return nullptr;
}

QMimeData *TrayModel::mimeData(const QModelIndexList &indexes) const
{
    Q_ASSERT(indexes.size() == 1);

    QMimeData *mime = new QMimeData;
    mime->setData(TRAY_DRAG_FALG, QByteArray());
    for (auto index : indexes) {
        if (!index.isValid())
            continue;

        int itemIndex = index.row();
        auto info = m_winInfos.at(itemIndex);
        mime->setData("type", QByteArray::number(static_cast<int>(info.type)));
        mime->setData("key", info.key.toLatin1());
        mime->setData("winId", QByteArray::number(info.winId));
        mime->setData("servicePath", info.servicePath.toLatin1());

        //TODO 支持多个index的数据，待支持
    }
    return mime;
}

QVariant TrayModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return QVariant();

    int itemIndex = index.row();
    const WinInfo &info = m_winInfos[itemIndex];

    switch (role) {
    case Role::TypeRole:
        return info.type;
    case Role::KeyRole:
        return info.key;
    case Role::WinIdRole:
        return info.winId;
    case Role::ServiceRole:
        return info.servicePath;
    case Role::PluginInterfaceRole:
        return (qulonglong)(info.pluginInter);
    case Role::Blank:
        return indexDragging(index);
    default:
        return QVariant();
    }
}

bool TrayModel::removeRows(int row, int count, const QModelIndex &parent)
{
    Q_UNUSED(count);
    Q_UNUSED(parent);

    if (m_winInfos.size() - 1 < row)
        return false;

    beginRemoveRows(parent, row, row);
    m_dragInfo = m_winInfos.takeAt(row);
    endRemoveRows();

    return true;
}

bool TrayModel::canDropMimeData(const QMimeData *data, Qt::DropAction action, int row, int column, const QModelIndex &parent) const
{
    Q_UNUSED(action)
    Q_UNUSED(row)
    Q_UNUSED(column)

    TrayIconType iconType = parent.data(TrayModel::Role::TypeRole).value<TrayIconType>();
    if (iconType == TrayIconType::ExpandIcon)
        return false;

    return data->formats().contains(TRAY_DRAG_FALG);
}

Qt::ItemFlags TrayModel::flags(const QModelIndex &index) const
{
    const Qt::ItemFlags defaultFlags = QAbstractListModel::flags(index);
    m_view->openPersistentEditor(index);

    return defaultFlags | Qt::ItemIsEditable |  Qt::ItemIsDragEnabled | Qt::ItemIsDropEnabled;
}

int TrayModel::rowCount(const QModelIndex &parent) const
{
    Q_UNUSED(parent);
    return m_winInfos.size();
}

bool TrayModel::isIconTray()
{
    return m_isTrayIcon;
}

void TrayModel::clear()
{
    beginResetModel();
    m_winInfos.clear();
    endResetModel();
}

void TrayModel::onXEmbedTrayAdded(quint32 winId)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.winId == winId)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = XEmbed;
    info.key = "wininfo:" + QString::number(winId);
    info.winId = winId;
    m_winInfos.append(info);
    endInsertRows();
}

void TrayModel::onXEmbedTrayRemoved(quint32 winId)
{
    for (auto info : m_winInfos) {
        if (info.winId == winId)  {
            int index = m_winInfos.indexOf(info);
            beginRemoveRows(QModelIndex(),  index, index);
            m_winInfos.removeOne(info);
            endRemoveRows();
            return;
        }
    }
}

QString TrayModel::fileNameByServiceName(const QString &serviceName)
{
    QStringList serviceInfo = serviceName.split("/");
    if (serviceInfo.size() <= 0)
        return QString();

    QDBusInterface dbsInterface("org.freedesktop.DBus", "/org/freedesktop/DBus",
                                "org.freedesktop.DBus", QDBusConnection::sessionBus(), this);
    QDBusMessage msg = dbsInterface.call("GetConnectionUnixProcessID", serviceInfo[0] );
    QList<QVariant> arguments = msg.arguments();
    if (arguments.size() == 0)
        return QString();

    QVariant v = arguments.at(0);
    uint pid = v.toUInt();
    QString path = QString("/proc/%1/cmdline").arg(pid);
    QFile file(path);
    if (file.open(QIODevice::ReadOnly)) {
        const QString fileName = file.readAll();
        file.close();
        return fileName;
    }

    return QString();
}

bool TrayModel::isTypeWriting(const QString &servicePath)
{
    const QString appFilePath = fileNameByServiceName(servicePath);
    return (appFilePath.startsWith("/usr/bin/fcitx") || appFilePath.endsWith("chinime-qim"));
}

void TrayModel::systemItemAdded(PluginsItemInterface *itemInter)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.pluginInter == itemInter)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());

    WinInfo info;
    info.type = SystemItem;
    info.pluginInter = itemInter;
    m_winInfos.append(info);

    endInsertRows();
}

void TrayModel::onSniTrayAdded(const QString &servicePath)
{
    bool typeWriting = isTypeWriting(servicePath);
    if ((m_hasInputMethod && !typeWriting) || (!m_hasInputMethod && typeWriting))
        return;

    int citxIndex = -1;
    for (int i = 0; i < m_winInfos.size(); i++) {
        WinInfo info = m_winInfos[i];
        if (info.servicePath == servicePath)
            return;

        if (typeWriting && info.isTypeWriting)
            citxIndex = i;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = Sni;
    info.key = "sni:" + servicePath;
    info.servicePath = servicePath;
    info.isTypeWriting = typeWriting;    // 是否为输入法
    if (typeWriting) {
        if (citxIndex < 0) {
            m_winInfos.append(info);
        } else {
            // 如果输入法在指定位置，则将输入法移动到指定位置
            m_winInfos[citxIndex] = info;
            QTimer::singleShot(150, this, [ = ] {
                 // 对比需要变化的图标
                 emit requestUpdateWidget({ citxIndex });
            });
        }
    } else {
        m_winInfos.append(info);
    }
    endInsertRows();
}

void TrayModel::onSniTrayRemoved(const QString &servicePath)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.servicePath == servicePath)  {
            int index = m_winInfos.indexOf(info);

            // 如果为输入法，则无需立刻删除，等100毫秒后再观察是否会删除输入法(因为在100毫秒内如果是切换输入法，就会很快发送add信号)
            if (info.isTypeWriting) {
                QTimer::singleShot(100, this, [ servicePath, this ] {
                    for (WinInfo info : m_winInfos) {
                        if (info.servicePath == servicePath) {
                            int index = m_winInfos.indexOf(info);
                            beginRemoveRows(QModelIndex(), index, index);
                            m_winInfos.removeOne(info);
                            endRemoveRows();
                        }
                    }
                });
            } else {
                beginRemoveRows(QModelIndex(), index, index);
                m_winInfos.removeOne(info);
                endRemoveRows();
            }
            break;
        }
    }
}

void TrayModel::onIndicatorFounded(const QString &indicatorName)
{
    const QString &itemKey = IndicatorTrayItem::toIndicatorKey(indicatorName);
    if (exist(itemKey) || !IndicatorTrayItem::isIndicatorKey(itemKey))
        return;

    IndicatorPlugin *indicatorTray = nullptr;
    if (!m_indicatorMap.keys().contains(indicatorName)) {
        indicatorTray = new IndicatorPlugin(indicatorName, this);
        m_indicatorMap[indicatorName] = indicatorTray;
    } else {
        indicatorTray = m_indicatorMap[itemKey];
    }

    connect(indicatorTray, &IndicatorPlugin::delayLoaded, indicatorTray, [ = ] {
        onIndicatorAdded(indicatorName);
    }, Qt::UniqueConnection);

    connect(indicatorTray, &IndicatorPlugin::removed, this, [ = ] {
        onIndicatorRemoved(indicatorName);
    }, Qt::UniqueConnection);
}

void TrayModel::onIndicatorAdded(const QString &indicatorName)
{
    const QString &itemKey = IndicatorTrayItem::toIndicatorKey(indicatorName);
    for (const WinInfo &info : m_winInfos) {
        if (info.key == itemKey)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = Incicator;
    info.key = itemKey;
    m_winInfos.append(info);
    endInsertRows();
}

void TrayModel::onIndicatorRemoved(const QString &indicatorName)
{
    const QString &itemKey = IndicatorTrayItem::toIndicatorKey(indicatorName);
    removeRow(itemKey);
}

void TrayModel::onSystemItemRemoved(PluginsItemInterface *itemInter)
{
    beginInsertRows(QModelIndex(), rowCount(), rowCount());

    for (const WinInfo &info : m_winInfos) {
        if (info.pluginInter != itemInter)
            continue;

        m_winInfos.removeOne(info);
        break;
    }

    endInsertRows();
}

void TrayModel::removeRow(const QString &itemKey)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.key == itemKey) {
            int index = m_winInfos.indexOf(info);
            beginRemoveRows(QModelIndex(),  index, index);
            m_winInfos.removeOne(info);
            endRemoveRows();
            break;
        }
    }
}

void TrayModel::addRow(WinInfo info)
{
    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.key == info.key)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    m_winInfos.append(info);
    endInsertRows();
}

void TrayModel::insertRow(int index, WinInfo info)
{
    for (int i = 0; i < m_winInfos.size(); i++) {
        const WinInfo &wininfo = m_winInfos[i];
        if (wininfo.key == info.key) {
            beginResetModel();
            m_winInfos.swap(index, i);
            endResetModel();
            return;
        }
    }
    beginInsertRows(QModelIndex(), index, index);
    m_winInfos.insert(index, info);
    endInsertRows();
}

bool TrayModel::exist(const QString &itemKey)
{
    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.key == itemKey)
            return true;
    }

    return false;
}
