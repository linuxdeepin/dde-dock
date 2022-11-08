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
#ifndef MULTIQUICKITEM_H
#define MULTIQUICKITEM_H

#include <quicksettingitem.h>

class MultiQuickItem : public QuickSettingItem
{
    Q_OBJECT

public:
    MultiQuickItem(PluginsItemInterface *const pluginInter, QWidget *parent = nullptr);
    ~MultiQuickItem() override;

    QuickSettingType type() const override;

protected:
    bool eventFilter(QObject *obj, QEvent *event) override;

private:
    void initUi();
    QString expandFileName() const;

private:
    bool m_selfDefine;
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
