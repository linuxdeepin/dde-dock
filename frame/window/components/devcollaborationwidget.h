/*
 * Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
 *
 * Author:     zhaoyingzhen <zhaoyingzhen@uniontech.com>
 *
 * Maintainer: zhaoyingzhen <zhaoyingzhen@uniontech.com>
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
#ifndef DEVICE_COLLABORATION_WIDGET_H
#define DEVICE_COLLABORATION_WIDGET_H

#include <QWidget>
#include <DListView>

DWIDGET_USE_NAMESPACE

class CollaborationDevice;
class CollaborationDevModel;

/*!
 * \brief The DevCollaborationWidget class
 * 设备跨端协同子页面
 */
class DevCollaborationWidget : public QWidget
{
    Q_OBJECT
public:
    explicit DevCollaborationWidget(QWidget *parent = nullptr);

protected:
    void showEvent(QShowEvent *event) override;

private slots:
    void loadDevice();
    void itemClicked(const QModelIndex &index);
    void itemStatusChanged();
    void refreshViewItem();

private:
    void initUI();
    void updateDeviceListView();

    void addItem(const CollaborationDevice *device);
    void resetWidgetSize();

private:
    CollaborationDevModel *m_deviceModel;
    DListView *m_deviceListView;
    QStandardItemModel *m_viewItemModel;
    QMap<QString, QStandardItem *> m_deviceItemMap;
    QStringList m_connectingDevices;

    QTimer *m_refreshTimer;
};

#endif // DEVICE_COLLABORATION_WIDGET_H
