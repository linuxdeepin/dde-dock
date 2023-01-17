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
#ifndef LARGERQUICKITEM_H
#define LARGERQUICKITEM_H

#include <quicksettingitem.h>

class QuickIconWidget;
class QWidget;
class QLabel;

// 插件在快捷面板中的展示样式，这个为大图标，展示为两列的那种，例如网络和蓝牙
class LargerQuickItem : public QuickSettingItem
{
    Q_OBJECT

public:
    LargerQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~LargerQuickItem() override;
    void doUpdate() override;
    void detachPlugin() override;

    QuickItemStyle type() const override;

protected:
    bool eventFilter(QObject *obj, QEvent *event) override;
    void showEvent(QShowEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    void initUi();
    QString expandFileName() const;

private:
    QuickIconWidget *m_iconWidget;
    QLabel *m_nameLabel;
    QLabel *m_stateLabel;
    QWidget *m_itemWidgetParent;
};

/**
 * @brief The QuickIconWidget class
 * 图标的Widget
 */
class QuickIconWidget : public QWidget
{
    Q_OBJECT

public:
    explicit QuickIconWidget(PluginsItemInterface *pluginInter, const QString &itemKey, QWidget *parent = Q_NULLPTR);

protected:
    void paintEvent(QPaintEvent *event) override;

private:
    QColor foregroundColor() const;
    QPixmap pluginIcon(bool contailGrab = false) const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
};

#endif // MULTIQUICKITEM_H
