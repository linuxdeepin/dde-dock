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
#ifndef STANDARDQUICKITEM_H
#define STANDARDQUICKITEM_H

#include "quicksettingitem.h"

class QLabel;

// 插件在快捷面板中展示的样式，这个为默认，展示一行一列的那种
class StandardQuickItem : public QuickSettingItem
{
    Q_OBJECT

public:
    StandardQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~StandardQuickItem() override;

    QuickItemStyle type() const override;
    void doUpdate() override;
    void detachPlugin() override;

protected:
    void mouseReleaseEvent(QMouseEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    void initUi();
    QWidget *iconWidget(QWidget *parent);
    QPixmap pixmap() const;
    QLabel *findChildLabel(QWidget *parent, const QString &childObjectName) const;
    void updatePluginName(QLabel *textLabel);

private:
    QWidget *m_itemParentWidget;
    bool m_needPaint;
};

#endif // SINGLEQUICKITEM_H
