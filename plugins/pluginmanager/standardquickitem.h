// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
