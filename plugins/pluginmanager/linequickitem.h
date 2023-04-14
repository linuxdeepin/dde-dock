// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef LINEQUICKITEM_H
#define LINEQUICKITEM_H

#include "quicksettingitem.h"

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
    void resizeSelf();

private:
    QWidget *m_centerWidget;
    QWidget *m_centerParentWidget;
};

#endif // FULLQUICKITEM_H
