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
#ifndef BLUETOOTHWIDGET_H
#define BLUETOOTHWIDGET_H

#include <QWidget>

class QLabel;
class AdaptersManager;
class Adapter;
class QVBoxLayout;

namespace Dtk { namespace Widget { class DListView; class DSwitchButton; } }

using namespace Dtk::Widget;

class BluetoothWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BluetoothWidget(AdaptersManager *adapterManager, QWidget *parent = nullptr);
    ~BluetoothWidget() override;

protected Q_SLOTS:
    void onAdapterIncreased(Adapter *adapter);
    void onAdapterDecreased(Adapter *adapter);
    void onCheckedChanged(bool checked);

private:
    void initUi();
    void initConnection();
    void updateCheckStatus();
    void adjustHeight();

private:
    DSwitchButton *m_switchButton;
    QWidget *m_headerWidget;
    QWidget *m_adapterWidget;
    AdaptersManager *m_adaptersManager;
    QVBoxLayout *m_adapterLayout;
};

#endif // BLUETOOTHWIDGET_H
