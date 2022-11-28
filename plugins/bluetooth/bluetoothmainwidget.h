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
#ifndef BLUETOOTHMAINWIDGET_H
#define BLUETOOTHMAINWIDGET_H

#include <QWidget>

class AdaptersManager;
class QLabel;
class Adapter;

class BluetoothMainWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BluetoothMainWidget(AdaptersManager *adapterManager, QWidget *parent = nullptr);
    ~BluetoothMainWidget();

Q_SIGNALS:
    void requestExpand();

protected:
    bool eventFilter(QObject *watcher, QEvent *event) override;

private:
    void initUi();
    void initConnection();

    void updateExpandIcon();

    bool isOpen() const;
    QString bluetoothIcon(bool isOpen) const;

private Q_SLOTS:
    void onAdapterChanged();

private:
    AdaptersManager *m_adapterManager;
    QWidget *m_iconWidget;
    QLabel *m_nameLabel;
    QLabel *m_stateLabel;
    QLabel *m_expandLabel;
};

#endif // BLUETOOTHMAINWIDGET_H
