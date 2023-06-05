// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "tray_model.h"
#include "tray_monitor.h"

#include "indicatortrayitem.h"
#include "indicatorplugin.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"
#include "docksettings.h"
#include "platformutils.h"

#include <QMimeData>
#include <QIcon>
#include <QDebug>
#include <QAbstractItemModel>
#include <QDBusInterface>

#define TRAY_DRAG_FALG "tray_drag"
#define DOCKQUICKTRAYNAME "Dock_Quick_Tray_Name"

TrayModel *TrayModel::getDockModel()
{
    static TrayModel *model = nullptr;
    if (!model) {
        model = new TrayModel(false);
        TrayModel *iconModel = getIconModel();
        connect(iconModel, &TrayModel::rowsRemoved, model, [ = ] {
            model->setExpandVisible(iconModel->rowCount() > 0);
        });
        connect(iconModel, &TrayModel::rowsInserted, model, [ = ] {
            model->setExpandVisible(iconModel->rowCount() > 0);
        });
        connect(iconModel, &TrayModel::rowCountChanged, model, [ = ] {
            model->setExpandVisible(iconModel->rowCount() > 0);
        });
    }

    return model;
}

TrayModel *TrayModel::getIconModel()
{
    static TrayModel model(true);
    return &model;
}

TrayModel::TrayModel(bool isIconTray, QObject *parent)
    : QAbstractListModel(parent)
    , m_dragModelIndex(QModelIndex())
    , m_dropModelIndex(QModelIndex())
    , m_monitor(new TrayMonitor(this))
    , m_isTrayIcon(isIconTray)
{
    connect(m_monitor, &TrayMonitor::xEmbedTrayAdded, this, &TrayModel::onXEmbedTrayAdded);
    connect(m_monitor, &TrayMonitor::xEmbedTrayRemoved, this, &TrayModel::onXEmbedTrayRemoved);

    connect(m_monitor, &TrayMonitor::sniTrayAdded, this, &TrayModel::onSniTrayAdded);
    connect(m_monitor, &TrayMonitor::sniTrayRemoved, this, &TrayModel::onSniTrayRemoved);

    connect(m_monitor, &TrayMonitor::indicatorFounded, this, &TrayModel::onIndicatorFounded);

    connect(m_monitor, &TrayMonitor::systemTrayAdded, this, &TrayModel::onSystemTrayAdded);
    connect(m_monitor, &TrayMonitor::systemTrayRemoved, this, &TrayModel::onSystemTrayRemoved);

    connect(m_monitor, &TrayMonitor::requestUpdateIcon, this, &TrayModel::requestUpdateIcon);
    connect(DockSettings::instance(), &DockSettings::quickPluginsChanged, this, &TrayModel::onSettingChanged);

    m_fixedTrayNames = DockSettings::instance()->getTrayItemsOnDock();
    m_fixedTrayNames.removeDuplicates();
}

void TrayModel::dropSwap(int newPos)
{
    if (!m_dragModelIndex.isValid())
        return;

    int row = m_dragModelIndex.row();

    if (row < m_winInfos.size())
        m_dragInfo = m_winInfos.takeAt(row);

    WinInfo name = m_dragInfo;
    m_winInfos.insert(newPos, name);

    emit QAbstractItemModel::dataChanged(m_dragModelIndex, m_dropModelIndex);
    requestRefreshEditor();
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

void TrayModel::setExpandVisible(bool visible, bool openExpand)
{
    // 如果是托盘，不支持展开图标
    if (m_isTrayIcon)
        return;

    if (visible) {
        // 如果展开图标已经存在，则不添加,
        for (WinInfo &winInfo : m_winInfos) {
            if (winInfo.type == TrayIconType::ExpandIcon) {
                winInfo.expand = openExpand;
                return;
            }
        }
        // 如果是任务栏图标，则添加托盘展开图标
        beginInsertRows(QModelIndex(), rowCount(), rowCount());
        WinInfo info;
        info.type = TrayIconType::ExpandIcon;
        info.expand = openExpand;
        m_winInfos.insert(0, info);  // 展开图标始终显示在第一个
        endInsertRows();

        Q_EMIT requestRefreshEditor();
        Q_EMIT rowCountChanged();
    } else {
        // 如果隐藏，则直接从列表中移除
        bool rowChanged = false;
        beginResetModel();
        for (const WinInfo &winInfo : m_winInfos) {
            if (winInfo.type == TrayIconType::ExpandIcon) {
                m_winInfos.removeOne(winInfo);
                rowChanged = true;
            }
        }
        endResetModel();
        if (rowChanged)
            Q_EMIT rowCountChanged();
    }
}

void TrayModel::updateOpenExpand(bool openExpand)
{
    for (WinInfo &winInfo : m_winInfos) {
        if (winInfo.type == TrayIconType::ExpandIcon)
            winInfo.expand = openExpand;
    }
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
        mime->setData("itemKey", info.itemKey.toLatin1());
        mime->setData("winId", QByteArray::number(info.winId));
        mime->setData("servicePath", info.servicePath.toLatin1());
        mime->setData("isTypeWritting", info.isTypeWriting ? "1" : "0");
        mime->setData("expand", info.expand ? "1" : "0");
        mime->setImageData(QVariant::fromValue((qulonglong)(info.pluginInter)));

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
    case Role::ExpandRole:
        return info.expand;
    case Role::ItemKeyRole:
        return info.itemKey;
    case Role::Blank:
        return indexDragging(index);
    default:
        return QVariant();
    }
}

bool TrayModel::removeRows(int row, int count, const QModelIndex &parent)
{
    Q_UNUSED(count);

    if (m_winInfos.size() - 1 < row)
        return false;

    beginRemoveRows(parent, row, row);
    m_dragInfo = m_winInfos.takeAt(row);
    endRemoveRows();

    Q_EMIT rowCountChanged();

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
    Q_EMIT requestOpenEditor(index);

    return (defaultFlags | Qt::ItemIsEditable |  Qt::ItemIsDragEnabled | Qt::ItemIsDropEnabled) & ~Qt::ItemIsSelectable;
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

bool TrayModel::hasExpand() const
{
    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.type == TrayIconType::ExpandIcon)
            return true;
    }

    return false;
}

bool TrayModel::isEmpty() const
{
    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.type != TrayIconType::ExpandIcon)
            return false;
    }

    return true;
}

void TrayModel::clear()
{
    beginResetModel();
    m_winInfos.clear();
    endResetModel();

    Q_EMIT rowCountChanged();
}

WinInfo TrayModel::getWinInfo(const QModelIndex &index)
{
    int row = index.row();
    if (row < 0 || row >= m_winInfos.size())
        return WinInfo();

    return m_winInfos[row];
}

void TrayModel::onXEmbedTrayAdded(quint32 winId)
{
    if (!xembedCanExport(winId))
        return;

    for (const WinInfo &info : m_winInfos) {
        if (info.winId == winId)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = XEmbed;
    info.key = "wininfo:" + QString::number(winId);
    info.itemKey = xembedItemKey(winId);
    info.winId = winId;
    m_winInfos.append(info);
    sortItems();
    endInsertRows();

    Q_EMIT rowCountChanged();
}

void TrayModel::onXEmbedTrayRemoved(quint32 winId)
{
    for (auto info : m_winInfos) {
        if (info.winId == winId)  {
            int index = m_winInfos.indexOf(info);
            beginRemoveRows(QModelIndex(),  index, index);
            m_winInfos.removeOne(info);
            endRemoveRows();

            Q_EMIT rowCountChanged();
            return;
        }
    }
}

QString TrayModel::fileNameByServiceName(const QString &serviceName) const
{
    QStringList serviceInfo = serviceName.split("/");
    if (serviceInfo.size() <= 0)
        return QString();

    QDBusInterface dbsInterface("org.freedesktop.DBus", "/org/freedesktop/DBus",
                                "org.freedesktop.DBus", QDBusConnection::sessionBus());
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

bool TrayModel::isTypeWriting(const QString &servicePath) const
{
    const QString appFilePath = fileNameByServiceName(servicePath);
    return (appFilePath.startsWith("/usr/bin/fcitx") || appFilePath.endsWith("chinime-qim"));
}

void TrayModel::saveConfig(int index, const WinInfo &winInfo)
{
    if (m_isTrayIcon) {
        // 如果是从任务栏将图标移动到托盘，就从配置中移除
        if (!m_fixedTrayNames.contains(winInfo.itemKey))
            return;

        m_fixedTrayNames.removeOne(winInfo.itemKey);
    } else {
        // 如果是将图标从托盘移到任务栏上面，就增加到配置中
        if (m_fixedTrayNames.contains(winInfo.itemKey))
            return;

        if (index >= 0 && index < m_fixedTrayNames.size()) {
            m_fixedTrayNames.insert(index, winInfo.itemKey);
        } else {
            m_fixedTrayNames << winInfo.itemKey;
        }
    }

    DockSettings::instance()->updateTrayItemsOnDock(m_fixedTrayNames);
}

void TrayModel::removeWinInfo(WinInfo winInfo)
{
    for (const WinInfo &info : m_winInfos) {
        if (winInfo == info) {
            int index = m_winInfos.indexOf(info);
            beginRemoveRows(QModelIndex(),  index, index);
            m_winInfos.removeOne(info);
            endRemoveRows();

            Q_EMIT rowCountChanged();
            break;
        }
    }
}

bool TrayModel::inTrayConfig(const QString itemKey) const
{
    if (m_isTrayIcon) {
        // 如果是托盘区域，显示所有不在配置中的应用
        return !m_fixedTrayNames.contains(itemKey);
    }
    // 如果是任务栏区域，显示所有在配置中的应用
    return m_fixedTrayNames.contains(itemKey);
}

QString TrayModel::xembedItemKey(quint32 winId) const
{
    return QString("embed:%1").arg(PlatformUtils::getAppNameForWindow(winId));
}

bool TrayModel::xembedCanExport(quint32 winId) const
{
    return inTrayConfig(xembedItemKey(winId));
}

QString TrayModel::sniItemKey(const QString &servicePath) const
{
    if (isTypeWriting(servicePath))
        return "fcitx";

    QString fileName = fileNameByServiceName(servicePath);
    return QString("sni:%1").arg(fileName.mid(fileName.lastIndexOf("/") + 1));
}

bool TrayModel::sniCanExport(const QString &servicePath) const
{
    return inTrayConfig(sniItemKey(servicePath));
}

bool TrayModel::indicatorCanExport(const QString &indicatorName) const
{
    return inTrayConfig(IndicatorTrayItem::toIndicatorKey(indicatorName));
}

QString TrayModel::systemItemKey(const QString &pluginName) const
{
    return QString("systemItem:%1").arg(pluginName);
}

bool TrayModel::systemItemCanExport(const QString &pluginName) const
{
    return inTrayConfig(systemItemKey(pluginName));
}

void TrayModel::sortItems()
{
    // 如果当前是展开托盘的内容，则无需排序
    if (m_isTrayIcon)
        return;

    // 数据排列，展开按钮始终排在最前面，输入法始终排在最后面
    WinInfos expandWin;
    WinInfos inputMethodWin;
    // 从列表中获取输入法和展开按钮
    for (const WinInfo &winInfo : m_winInfos) {
        switch (winInfo.type) {
        case TrayIconType::ExpandIcon: {
            expandWin << winInfo;
            break;
        }
        case TrayIconType::Sni: {
            if (winInfo.isTypeWriting)
                inputMethodWin << winInfo;
            break;
        }
        default:
            break;
        }
    }
    // 从列表中移除展开按钮
    for (const WinInfo &winInfo : expandWin)
        m_winInfos.removeOne(winInfo);

    // 从列表中移除输入法
    for (const WinInfo &winInfo : inputMethodWin)
        m_winInfos.removeOne(winInfo);

    // 将展开按钮添加到列表的最前面
    for (int i = expandWin.size() - 1; i >= 0; i--)
        m_winInfos.push_front(expandWin[i]);

    // 将输入法添加到列表的最后面
    for (int i = 0; i < inputMethodWin.size(); i++)
        m_winInfos.push_back(inputMethodWin[i]);
}

void TrayModel::onSniTrayAdded(const QString &servicePath)
{
    if (!sniCanExport(servicePath))
        return;

    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.servicePath == servicePath)
            return;
    }

    bool typeWriting = isTypeWriting(servicePath);

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = Sni;
    info.key = "sni:" + servicePath;
    info.itemKey = sniItemKey(servicePath);
    info.servicePath = servicePath;
    info.isTypeWriting = typeWriting;    // 是否为输入法
    m_winInfos.append(info);

    sortItems();
    endInsertRows();

    Q_EMIT rowCountChanged();
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

                Q_EMIT rowCountChanged();
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
    if (!indicatorCanExport(indicatorName))
        return;

    const QString &itemKey = IndicatorTrayItem::toIndicatorKey(indicatorName);
    for (const WinInfo &info : m_winInfos) {
        if (info.itemKey == itemKey)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());
    WinInfo info;
    info.type = Incicator;
    info.key = itemKey;
    info.itemKey = itemKey;
    m_winInfos.append(info);

    sortItems();
    endInsertRows();

    Q_EMIT rowCountChanged();
}

void TrayModel::onIndicatorRemoved(const QString &indicatorName)
{
    const QString &itemKey = IndicatorTrayItem::toIndicatorKey(indicatorName);
    removeRow(itemKey);
}

void TrayModel::onSystemTrayAdded(PluginsItemInterface *itemInter)
{
    if (!systemItemCanExport(itemInter->pluginName()))
        return;

    for (const WinInfo &info : m_winInfos) {
        if (info.pluginInter == itemInter)
            return;
    }

    beginInsertRows(QModelIndex(), rowCount(), rowCount());

    WinInfo info;
    info.type = SystemItem;
    info.pluginInter = itemInter;
    info.itemKey = systemItemKey(itemInter->pluginName());
    m_winInfos.append(info);

    sortItems();
    endInsertRows();

    Q_EMIT rowCountChanged();
}

void TrayModel::onSystemTrayRemoved(PluginsItemInterface *itemInter)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.pluginInter != itemInter)
            continue;

        beginInsertRows(QModelIndex(), rowCount(), rowCount());
        m_winInfos.removeOne(info);
        endInsertRows();

        Q_EMIT rowCountChanged();
        break;
    }
}

void TrayModel::onSettingChanged(const QStringList &value)
{
    // 先将其转换为任务栏上的图标列表
    m_fixedTrayNames = value;
}

void TrayModel::removeRow(const QString &itemKey)
{
    for (const WinInfo &info : m_winInfos) {
        if (info.itemKey == itemKey) {
            int index = m_winInfos.indexOf(info);
            beginRemoveRows(QModelIndex(),  index, index);
            m_winInfos.removeOne(info);
            endRemoveRows();

            Q_EMIT rowCountChanged();
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

    beginResetModel();
    m_winInfos.append(info);
    sortItems();
    endResetModel();

    Q_EMIT requestRefreshEditor();
    Q_EMIT rowCountChanged();
}

void TrayModel::insertRow(int index, WinInfo info)
{
    for (int i = 0; i < m_winInfos.size(); i++) {
        const WinInfo &wininfo = m_winInfos[i];
        if (wininfo.key == info.key) {
            beginResetModel();
            m_winInfos.swapItemsAt(index, i);
            endResetModel();
            return;
        }
    }
    beginInsertRows(QModelIndex(), index, index);
    m_winInfos.insert(index, info);
    sortItems();

    endInsertRows();

    Q_EMIT requestRefreshEditor();
    Q_EMIT rowCountChanged();
}

bool TrayModel::exist(const QString &itemKey)
{
    for (const WinInfo &winInfo : m_winInfos) {
        if (winInfo.key == itemKey)
            return true;
    }

    return false;
}
