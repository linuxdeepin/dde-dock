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

#include "dockitem.h"

class PluginsItemInterface;
class QuickIconWidget;

class QuickSettingItem : public DockItem
{
    Q_OBJECT

Q_SIGNALS:
    void detailClicked(PluginsItemInterface *);

public:
    QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, const QJsonObject &metaData, QWidget *parent = nullptr);
    ~QuickSettingItem() override;
    PluginsItemInterface *pluginItem() const;
    ItemType itemType() const override;
    const QPixmap dragPixmap();
    const QString itemKey() const;
    bool isPrimary() const;

protected:

    bool eventFilter(QObject *obj, QEvent *event) override;
    void paintEvent(QPaintEvent *e) override;
    QColor foregroundColor() const;

private:
    void initUi();
    QString expandFileName();
    QPixmap pluginIcon() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
    QJsonObject m_metaData;
    QWidget *m_iconWidgetParent;
    QuickIconWidget *m_iconWidget;
    QWidget *m_textWidget;
    QLabel *m_nameLabel;
    QLabel *m_stateLabel;
};

/**
 * @brief The QuickIconWidget class
 * 图标的Widget
 */
class QuickIconWidget : public QWidget
{
    Q_OBJECT

public:
    explicit QuickIconWidget(PluginsItemInterface *pluginInter, const QString &itemKey, bool isPrimary, QWidget *parent = Q_NULLPTR);

protected:
    void paintEvent(QPaintEvent *event) override;

private:
    QColor foregroundColor() const;
    QPixmap pluginIcon() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
    bool m_isPrimary;
};

#endif // QUICKSETTINGITEM_H
