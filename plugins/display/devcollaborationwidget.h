// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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

signals:
    void sizeChanged();

protected:
    void showEvent(QShowEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

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
