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
#ifndef APPMULTIITEM_H
#define APPMULTIITEM_H

#include "dockitem.h"
#include "dbusutil.h"

struct SHMInfo;
struct _XImage;
typedef _XImage XImage;
class AppItem;

class AppMultiItem : public DockItem
{
    Q_OBJECT

    friend class AppItem;

public:
    AppMultiItem(AppItem *appItem, WId winId, const WindowInfo &windowInfo, QWidget *parent = Q_NULLPTR);
    ~AppMultiItem() override;

    QSize suitableSize(int size) const;
    AppItem *appItem() const;
    quint32 winId() const;
    const WindowInfo &windowInfo() const;

    ItemType itemType() const override;

protected:
    void paintEvent(QPaintEvent *) override;
    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    bool isKWinAvailable() const;
    QImage snapImage() const;
    SHMInfo *getImageDSHM() const;
    XImage *getImageXlib() const;
    void initMenu();
    void initConnection();

private Q_SLOTS:
    void onOpen();
    void onCurrentWindowChanged(uint32_t value);

private:
    AppItem *m_appItem;
    WindowInfo m_windowInfo;
    DockEntryInter *m_entryInter;
    QImage m_snapImage;
    WId m_winId;
    QMenu *m_menu;
};

#endif // APPMULTIITEM_H
