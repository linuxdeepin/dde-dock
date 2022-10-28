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
    quint32 winId;
    QString servicePath;
    bool isTypeWriting;
    PluginsItemInterface *pluginInter;

    WinInfo() : type(UnKnow)
      , key(QString())
      , winId(0)
      , servicePath(QString())
      , isTypeWriting(false)
      , pluginInter(nullptr) {}

    bool operator==(const WinInfo &other) {
        return this->type == other.type
                && this->key == other.key
                && this->winId == other.winId
                && this->servicePath == other.servicePath
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
        Blank
    };

    typedef QList<WinInfo> WinInfos;

    TrayModel(QListView *view, bool isIconTray, bool hasInputMethod, QObject *parent = Q_NULLPTR);

    void dropSwap(int newPos);
    void dropInsert(int newPos);

    void clearDragDropIndex();
    void setDragingIndex(const QModelIndex index);
    void setDragDropIndex(const QModelIndex index);

    void setDragKey(const QString &key);

    bool indexDragging(const QModelIndex &index) const;

    IndicatorTrayItem *indicatorWidget(const QString &indicatorName) const;

    int rowCount(const QModelIndex &parent = QModelIndex()) const Q_DECL_OVERRIDE;
    bool isIconTray();

    void clear();

Q_SIGNALS:
    void requestUpdateIcon(quint32);
    void requestUpdateWidget(const QList<int> &);

public Q_SLOTS:
    void removeRow(const QString &itemKey);
    void addRow(WinInfo info);
    void insertRow(int index, WinInfo info);

private Q_SLOTS:
    void onXEmbedTrayAdded(quint32 winId);
    void onXEmbedTrayRemoved(quint32 winId);
    void onSniTrayAdded(const QString &servicePath);
    void onSniTrayRemoved(const QString &servicePath);

    void onIndicatorFounded(const QString &indicatorName);
    void onIndicatorAdded(const QString &indicatorName);
    void onIndicatorRemoved(const QString &indicatorName);

    void onSystemItemRemoved(PluginsItemInterface *itemInter);

private:
    bool exist(const QString &itemKey);
    QString fileNameByServiceName(const QString &serviceName);
    bool isTypeWriting(const QString &servicePath);

    void systemItemAdded(PluginsItemInterface *itemInter);

protected:
    QMimeData *mimeData(const QModelIndexList &indexes) const Q_DECL_OVERRIDE;
    QVariant data(const QModelIndex &index, int role) const Q_DECL_OVERRIDE;
    bool removeRows(int row, int count, const QModelIndex &parent) Q_DECL_OVERRIDE;
    bool canDropMimeData(const QMimeData *data, Qt::DropAction action,
                         int row, int column, const QModelIndex &parent) const Q_DECL_OVERRIDE;
    Qt::ItemFlags flags(const QModelIndex &index) const Q_DECL_OVERRIDE;

private:
    WinInfos m_winInfos;

    QModelIndex m_dragModelIndex;
    QModelIndex m_dropModelIndex;
    WinInfo m_dragInfo;
    QListView *m_view;
    TrayMonitor *m_monitor;

    QString m_dragKey;

    QMap<QString, IndicatorPlugin *> m_indicatorMap;
    bool m_isTrayIcon;
    bool m_hasInputMethod;
};

Q_DECLARE_METATYPE(TrayIconType);

#endif // TRAYMODEL_H
