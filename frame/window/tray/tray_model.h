// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRAYMODEL_H
#define TRAYMODEL_H

#include <QAbstractListModel>
#include <QObject>
#include <QListView>

class TrayMonitor;
class IndicatorPlugin;
class IndicatorTrayItem;
class PluginsItemInterface;

enum TrayIconType {
    UnKnow,
    XEmbed,
    Sni,
    Incicator,
    ExpandIcon,
    SystemItem
};

struct WinInfo {
    TrayIconType type;
    QString key;
    QString itemKey;
    quint32 winId;
    QString servicePath;
    bool isTypeWriting;
    bool expand;
    PluginsItemInterface *pluginInter;

    WinInfo() : type(UnKnow)
      , key(QString())
      , itemKey(QString())
      , winId(0)
      , servicePath(QString())
      , isTypeWriting(false)
      , expand(false)
      , pluginInter(nullptr) {}

    bool operator==(const WinInfo &other) {
        return this->type == other.type
                && this->key == other.key
                && this->winId == other.winId
                && this->servicePath == other.servicePath
                && this->itemKey == other.itemKey
                && this->isTypeWriting == other.isTypeWriting
                && this->pluginInter == other.pluginInter;
    }
};

class TrayModel : public QAbstractListModel
{
    Q_OBJECT

public:
    enum Role {
        TypeRole = Qt::UserRole + 1,
        KeyRole,
        WinIdRole,
        ServiceRole,
        PluginInterfaceRole,
        ExpandRole,
        ItemKeyRole,
        Blank
    };

    typedef QList<WinInfo> WinInfos;

    static TrayModel *getDockModel();
    static TrayModel *getIconModel();

    void dropSwap(int newPos);

    void clearDragDropIndex();
    void setDragingIndex(const QModelIndex index);
    void setDragDropIndex(const QModelIndex index);
    void setExpandVisible(bool visible, bool openExpand = false);
    void updateOpenExpand(bool openExpand);

    void setDragKey(const QString &key);

    bool indexDragging(const QModelIndex &index) const;

    IndicatorTrayItem *indicatorWidget(const QString &indicatorName) const;

    int rowCount(const QModelIndex &parent = QModelIndex()) const Q_DECL_OVERRIDE;
    bool isIconTray();
    bool hasExpand() const;
    bool isEmpty() const;

    void clear();
    WinInfo getWinInfo(const QModelIndex &index);
    void saveConfig(int index, const WinInfo &winInfo);
    void removeWinInfo(WinInfo winInfo);

Q_SIGNALS:
    void requestUpdateIcon(quint32);
    void requestOpenEditor(const QModelIndex &index, bool isOpen = true) const;
    void rowCountChanged();
    void requestRefreshEditor();

public Q_SLOTS:
    void addRow(WinInfo info);
    void insertRow(int index, WinInfo info);

protected:
    TrayModel(bool isIconTray, QObject *parent = Q_NULLPTR);

private Q_SLOTS:
    void onXEmbedTrayAdded(quint32 winId);
    void onXEmbedTrayRemoved(quint32 winId);
    void onSniTrayAdded(const QString &servicePath);
    void onSniTrayRemoved(const QString &servicePath);

    void onIndicatorFounded(const QString &indicatorName);
    void onIndicatorAdded(const QString &indicatorName);
    void onIndicatorRemoved(const QString &indicatorName);

    void onSystemTrayAdded(PluginsItemInterface *itemInter);
    void onSystemTrayRemoved(PluginsItemInterface *itemInter);

    void onSettingChanged(const QStringList &value);

protected:
    QMimeData *mimeData(const QModelIndexList &indexes) const Q_DECL_OVERRIDE;
    QVariant data(const QModelIndex &index, int role) const Q_DECL_OVERRIDE;
    bool removeRows(int row, int count, const QModelIndex &parent) Q_DECL_OVERRIDE;
    bool canDropMimeData(const QMimeData *data, Qt::DropAction action,
                         int row, int column, const QModelIndex &parent) const Q_DECL_OVERRIDE;
    Qt::ItemFlags flags(const QModelIndex &index) const Q_DECL_OVERRIDE;

private:
    void removeRow(const QString &itemKey);
    bool exist(const QString &itemKey);
    QString fileNameByServiceName(const QString &serviceName) const;
    bool isTypeWriting(const QString &servicePath) const;

    bool inTrayConfig(const QString itemKey) const;
    QString xembedItemKey(quint32 winId) const;
    bool xembedCanExport(quint32 winId) const;
    QString sniItemKey(const QString &servicePath) const;
    bool sniCanExport(const QString &servicePath) const;
    bool indicatorCanExport(const QString &indicatorName) const;
    QString systemItemKey(const QString &pluginName) const;
    bool systemItemCanExport(const QString &pluginName) const;
    void sortItems();

private:
    WinInfos m_winInfos;

    QModelIndex m_dragModelIndex;
    QModelIndex m_dropModelIndex;
    WinInfo m_dragInfo;
    TrayMonitor *m_monitor;

    QString m_dragKey;

    QMap<QString, IndicatorPlugin *> m_indicatorMap;
    QStringList m_fixedTrayNames;
    bool m_isTrayIcon;
};

Q_DECLARE_METATYPE(TrayIconType);

#endif // TRAYMODEL_H
