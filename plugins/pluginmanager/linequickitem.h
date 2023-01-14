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
#ifndef LINEQUICKITEM_H
#define LINEQUICKITEM_H

#include "quicksettingitem.h"

#include <DGuiApplicationHelper>

namespace Dtk {
namespace Widget {
class DBlurEffectWidget;
}
}

DGUI_USE_NAMESPACE

// 插件在快捷面板中的展示的样式，这个为整行显示的插件，例如声音，亮度调整和音乐播放等
class LineQuickItem : public QuickSettingItem
{
    Q_OBJECT

public:
    LineQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~LineQuickItem() override;
    void doUpdate() override;
    void detachPlugin() override;

    QuickItemStyle type() const override;

protected:
    bool eventFilter(QObject *obj, QEvent *event) override;

private:
    void initUi();
    void initConnection();
    void resizeSelf();

private Q_SLOTS:
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    QWidget *m_centerWidget;
    QWidget *m_centerParentWidget;
    Dtk::Widget::DBlurEffectWidget *m_effectWidget;
};

#endif // FULLQUICKITEM_H
