/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#ifndef FULLQUICKITEM_H
#define FULLQUICKITEM_H

#include "quicksettingitem.h"

class FullQuickItem : public QuickSettingItem
{
    Q_OBJECT

public:
    FullQuickItem(PluginsItemInterface *const pluginInter, QWidget *parent = nullptr);
    ~FullQuickItem() override;

    QuickSettingType type() const override;

protected:
    bool eventFilter(QObject *obj, QEvent *event) override;

private:
    void initUi();

private:
    QWidget *m_centerWidget;
    DBlurEffectWidget *m_effectWidget;
};

#endif // FULLQUICKITEM_H
