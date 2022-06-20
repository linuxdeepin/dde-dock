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

class QuickSettingItem : public DockItem
{
    Q_OBJECT

    friend class QuickSettingController;

Q_SIGNALS:
    void detailClicked(PluginsItemInterface *);

public:
    PluginsItemInterface *pluginItem() const;
    ItemType itemType() const override;
    const QPixmap dragPixmap();
    const QString itemKey() const;

protected:
    QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~QuickSettingItem() override;

    void paintEvent(QPaintEvent *e) override;
    QRect iconRect();
    QColor foregroundColor() const;

    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    int yMarginSpace();
    QString expandFileName();

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
};

#endif // QUICKSETTINGITEM_H
