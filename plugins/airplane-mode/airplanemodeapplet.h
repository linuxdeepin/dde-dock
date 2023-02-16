// Copyright (C) 2020 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef AIRPLANEMODEAPPLET_H
#define AIRPLANEMODEAPPLET_H

#include <QWidget>

namespace Dtk {
namespace Widget {
class DSwitchButton;
}
}

class QLabel;
class AirplaneModeApplet : public QWidget
{
    Q_OBJECT

public:
    explicit AirplaneModeApplet(QWidget *parent = nullptr);
    void setEnabled(bool enable);

signals:
    void enableChanged(bool enable);

private:
    Dtk::Widget::DSwitchButton *m_switchBtn;
};

#endif // AIRPLANEMODEAPPLET_H
