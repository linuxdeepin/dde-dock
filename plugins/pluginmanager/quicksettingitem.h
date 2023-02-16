// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef QUICKSETTINGITEM_H
#define QUICKSETTINGITEM_H

#include <QWidget>

class PluginsItemInterface;
class QuickIconWidget;

class QuickSettingItem : public QWidget
{
    Q_OBJECT

public:
    enum class QuickItemStyle {
        Standard = 1,          // 插件的UI显示单列
        Larger,                // 插件的UI显示双列，例如网络和蓝牙等
        Line                   // 插件的UI整行显示，例如声音，亮度、音乐播放等
    };

public:
    QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~QuickSettingItem() override;
    PluginsItemInterface *pluginItem() const;
    virtual const QPixmap dragPixmap();
    virtual void doUpdate() {}
    virtual void detachPlugin() {}
    const QString itemKey() const;

    virtual QuickItemStyle type() const = 0;

Q_SIGNALS:
    void requestShowChildWidget(QWidget *);

protected:
    void paintEvent(QPaintEvent *e) override;
    QColor foregroundColor() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
};

class QuickSettingFactory
{
public:
    static QuickSettingItem *createQuickWidget(PluginsItemInterface *const pluginInter, const QString &itemKey);
};

#endif // QUICKSETTINGITEM_H
