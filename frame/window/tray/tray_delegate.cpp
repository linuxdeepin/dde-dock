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
#include "pluginsiteminterface.h"
#include "quicksettingcontroller.h"
#include "systempluginitem.h"

#include <DGuiApplicationHelper>

#include <QPointer>
#include <QDebug>
#include <QEvent>
#include <QKeyEvent>
#include <QApplication>
#include <QPainterPath>

#include <xcb/xcb_icccm.h>
#include <X11/Xlib.h>

TrayDelegate::TrayDelegate(QListView *view, QObject *parent)
    : QStyledItemDelegate(parent)
    , m_position(Dock::Position::Bottom)
    , m_listView(view)
{
    connect(this, &TrayDelegate::requestDrag, this, &TrayDelegate::onUpdateExpand);
}

void TrayDelegate::setPositon(Dock::Position position)
{
    m_position = position;
    SNITrayItemWidget::setDockPostion(position);
    SystemPluginItem::setDockPostion(m_position);
}

QWidget *TrayDelegate::createEditor(QWidget *parent, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(option);

    TrayIconType type = index.data(TrayModel::TypeRole).value<TrayIconType>();
    QString key = index.data(TrayModel::KeyRole).value<QString>();
    QString servicePath = index.data(TrayModel::ServiceRole).value<QString>();
    quint32 winId = index.data(TrayModel::WinIdRole).value<quint32>();

    BaseTrayWidget *trayWidget = nullptr;
    if(type == TrayIconType::XEmbed) {
        if (Utils::IS_WAYLAND_DISPLAY) {
            static Display *display = XOpenDisplay(nullptr);
            static int screenp = 0;
            static xcb_connection_t *xcb_connection = xcb_connect(qgetenv("DISPLAY"), &screenp);
            trayWidget = new XEmbedTrayItemWidget(winId, xcb_connection, display, parent);
        } else {
            trayWidget = new XEmbedTrayItemWidget(winId, nullptr, nullptr, parent);
        }
        const TrayModel *model = qobject_cast<const TrayModel *>(index.model());
        if (model)
            connect(model, &TrayModel::requestUpdateIcon, trayWidget, &BaseTrayWidget::updateIcon);
    } else if (type == TrayIconType::Sni) {
        trayWidget = new SNITrayItemWidget(servicePath, parent);
    } else if (type == TrayIconType::ExpandIcon) {
        ExpandIconWidget *expandWidget = new ExpandIconWidget(parent);
        expandWidget->setPositon(m_position);
        bool openExpand = index.data(TrayModel::ExpandRole).toBool();
        if (openExpand)
            expandWidget->setTrayPanelVisible(true);

        trayWidget = expandWidget;
    } else if (type == TrayIconType::Incicator) {
        QString indicateName = key;
        int flagIndex = indicateName.indexOf("indicator:");
        if (flagIndex >= 0)
            indicateName = indicateName.right(indicateName.length() - QString("indicator:").length());
        IndicatorTrayItem *indicatorWidget = new IndicatorTrayItem(indicateName, parent);
        TrayModel *dataModel = qobject_cast<TrayModel *>(m_listView->model());
        if (IndicatorTrayItem *sourceIndicatorWidget = dataModel->indicatorWidget(key)) {
            const QByteArray pixmapData = sourceIndicatorWidget->pixmapData();
            if (!pixmapData.isEmpty())
                indicatorWidget->setPixmapData(pixmapData);
            const QString text = sourceIndicatorWidget->text();
            if (!text.isEmpty())
                indicatorWidget->setText(text);
        }
        trayWidget = indicatorWidget;
    } else if (type == TrayIconType::SystemItem) {
        PluginsItemInterface *pluginInter = (PluginsItemInterface *)(index.data(TrayModel::PluginInterfaceRole).toULongLong());
        if (pluginInter) {
            const QString itemKey = QuickSettingController::instance()->itemKey(pluginInter);
            trayWidget = new SystemPluginItem(pluginInter, itemKey, parent);
        }
    }

    if (trayWidget)
        trayWidget->setFixedSize(16, 16);

    return trayWidget;
}

void TrayDelegate::onUpdateExpand(bool on)
{
    ExpandIconWidget *expandwidget = expandWidget();

    if (on) {
        if (!expandwidget) {
            // 如果三角按钮不存在，那么就设置三角按钮可见，此时它会自动创建一个三角按钮
            TrayModel *model = qobject_cast<TrayModel *>(m_listView->model());
            if (model)
                model->setExpandVisible(true, true);
        } else {
            expandwidget->setTrayPanelVisible(true);
        }
    } else if (expandwidget) {
        // 如果释放鼠标，则判断当前鼠标的位置是否在托盘内部，如果在，则无需隐藏
        QPoint currentPoint = QCursor::pos();
        TrayGridWidget *view = ExpandIconWidget::popupTrayView();
        expandwidget->setTrayPanelVisible(view->geometry().contains(currentPoint));
    }
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

    // 如果是弹出托盘，则显示正常大小
    if (isPopupTray())
        return QSize(ITEM_SIZE, ITEM_SIZE);

    // 如果是任务栏的托盘，则高度显示为listView的高度或宽度
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        return QSize(ITEM_SIZE, m_listView->height());

    return QSize(m_listView->width(), ITEM_SIZE);
}

void TrayDelegate::updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(index);
    QRect rect = option.rect;
    // 让控件居中显示
    editor->setGeometry(rect.x() + (rect.width() - ICON_SIZE) / 2,
                        rect.y() + (rect.height() - ICON_SIZE) / 2,
                        ICON_SIZE, ICON_SIZE);
}

void TrayDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    Q_UNUSED(index);

    if (!isPopupTray())
        return;

    QColor borderColor;
    QColor backColor;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        // 白色主题的情况下
        borderColor = Qt::black;
        borderColor.setAlpha(static_cast<int>(255 * 0.05));
        backColor = Qt::white;
        backColor.setAlpha(static_cast<int>(255 * 0.4));
    } else {
        borderColor = Qt::black;
        borderColor.setAlpha(static_cast<int>(255 * 0.2));
        backColor = Qt::black;
        backColor.setAlpha(static_cast<int>(255 * 0.4));
    }

    painter->save();
    QPainterPath path;
    path.addRoundedRect(option.rect, 8, 8);
    painter->setRenderHint(QPainter::Antialiasing);
    painter->fillPath(path, backColor);
    painter->setPen(borderColor);
    painter->drawPath(path);
    painter->restore();
}

ExpandIconWidget *TrayDelegate::expandWidget()
{
    if (!m_listView)
        return nullptr;

    QAbstractItemModel *dataModel = m_listView->model();
    if (!dataModel)
        return nullptr;

    for (int i = 0; i < dataModel->rowCount() - 1; i++) {
        QModelIndex index = dataModel->index(i, 0);
        ExpandIconWidget *widget = qobject_cast<ExpandIconWidget *>(m_listView->indexWidget(index));
        if (widget)
            return widget;
    }

    return nullptr;
}

bool TrayDelegate::isPopupTray() const
{
    if (!m_listView)
        return false;

    TrayModel *dataModel = qobject_cast<TrayModel *>(m_listView->model());
    if (!dataModel)
        return false;

    return dataModel->isIconTray();
}
