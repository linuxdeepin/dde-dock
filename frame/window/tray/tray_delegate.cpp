/*
 * Copyright (C) 2018 ~ 2025 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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
#include "tray_delegate.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "widgets/xembedtrayitemwidget.h"
#include "widgets/indicatortrayitem.h"
#include "widgets/indicatorplugin.h"
#include "widgets/snitrayitemwidget.h"
#include "widgets/expandiconwidget.h"
#include "utils.h"

#include <QPointer>
#include <QDebug>
#include <QEvent>
#include <QKeyEvent>
#include <QApplication>

#include <xcb/xcb_icccm.h>
#include <X11/Xlib.h>

TrayDelegate::TrayDelegate(QObject *parent)
    : QStyledItemDelegate(parent)
    , m_position(Dock::Position::Bottom)
{
}

void TrayDelegate::setPositon(Dock::Position position)
{
    m_position = position;
}

QWidget *TrayDelegate::createEditor(QWidget *parent, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(option);

    TrayIconType type = index.data(TrayModel::TypeRole).value<TrayIconType>();
    QString key = index.data(TrayModel::KeyRole).value<QString>();
    QString servicePath = index.data(TrayModel::ServiceRole).value<QString>();
    quint32 winId = index.data(TrayModel::WinIdRole).value<quint32>();

    BaseTrayWidget *trayWidget = nullptr;
    if(type == TrayIconType::XEMBED) {
        if (Utils::IS_WAYLAND_DISPLAY) {
            trayWidget = new XEmbedTrayItemWidget(winId, nullptr, nullptr, parent);
        } else {
            int screenp = 0;
            static xcb_connection_t *xcb_connection = xcb_connect(qgetenv("DISPLAY"), &screenp);
            static Display *m_display = XOpenDisplay(nullptr);
            trayWidget = new XEmbedTrayItemWidget(winId, xcb_connection, m_display, parent) ;
        }

        const TrayModel *model = qobject_cast<const TrayModel *>(index.model());
        if (model)
            connect(model, &TrayModel::requestUpdateIcon, trayWidget, &BaseTrayWidget::updateIcon);
    } else if (type == TrayIconType::SNI) {
        trayWidget = new SNITrayItemWidget(servicePath, parent);
    } else if (type == TrayIconType::EXPANDICON) {
        ExpandIconWidget *widget = new ExpandIconWidget(parent);
        widget->setPositonValue(m_position);
        connect(widget, &ExpandIconWidget::trayVisbleChanged, this, [ = ](bool visible) {
            Q_EMIT visibleChanged(index, visible);
        });
        connect(this, &TrayDelegate::requestDrag, this, [ = ](bool on) {
            if (on) {
                widget->setTrayPanelVisible(true);
            } else {
                // 如果释放鼠标，则判断当前鼠标的位置是否在托盘内部，如果在，则无需隐藏
                QPoint currentPoint = QCursor::pos();
                TrayGridView *view = widget->popupTrayView();
                if (view->geometry().contains(currentPoint))
                    widget->setTrayPanelVisible(true);
                else
                    widget->setTrayPanelVisible(false);
            }
        });
        trayWidget = widget;
    } else if (type == TrayIconType::INDICATOR) {
        QString indicateName = key;
        int flagIndex = indicateName.indexOf("indicator:");
        if (flagIndex >= 0)
            indicateName = indicateName.right(indicateName.length() - QString("indicator:").length());
        IndicatorTrayItem *indicatorWidget = new IndicatorTrayItem(indicateName, parent);
        connect(indicatorWidget, &IndicatorTrayItem::removed, this, [ = ]{
            Q_EMIT removeRow(index);
        });
        trayWidget = indicatorWidget;
    }

    return trayWidget;
}

void TrayDelegate::setEditorData(QWidget *editor, const QModelIndex &index) const
{
    BaseTrayWidget *widget = static_cast<BaseTrayWidget *>(editor);
    if (widget) {
        widget->setNeedShow(!index.data(TrayModel::Blank).toBool());
    }
}

QSize TrayDelegate::sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(option);
    Q_UNUSED(index);

    return QSize(ITEM_SIZE, ITEM_SIZE);
}

void TrayDelegate::updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(index);
    QRect rect = option.rect;
    editor->setGeometry(rect.x() + ITEM_SPACING, rect.y() + ITEM_SPACING, ITEM_SIZE - (2 * ITEM_SPACING), ITEM_SIZE - 2 * ITEM_SPACING);
}
