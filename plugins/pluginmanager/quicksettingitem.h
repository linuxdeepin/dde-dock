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
