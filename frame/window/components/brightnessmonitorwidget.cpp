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
#include "brightnessmonitorwidget.h"
#include "brightnessmodel.h"
#include "customslider.h"
#include "settingdelegate.h"

#include <DListView>
#include <DDBusSender>

#include <QScrollBar>
#include <QLabel>
#include <QVBoxLayout>
#include <QEvent>
#include <QProcess>
#include <QDBusInterface>
#include <QDBusConnection>

DWIDGET_USE_NAMESPACE

#define ITEMSPACE 16

BrightnessMonitorWidget::BrightnessMonitorWidget(BrightnessModel *model, QWidget *parent)
    : QWidget(parent)
    , m_sliderWidget(new QWidget(this))
    , m_sliderLayout(new QVBoxLayout(m_sliderWidget))
    , m_descriptionLabel(new QLabel(tr("Collaboration"), this))
    , m_deviceList(new DListView(this))
    , m_brightModel(model)
    , m_model(new QStandardItemModel(this))
    , m_delegate(new SettingDelegate(m_deviceList))
{
    initUi();
    initConnection();
    reloadMonitor();

    QMetaObject::invokeMethod(this, [ this ]{
        resetHeight();
    }, Qt::QueuedConnection);
}

BrightnessMonitorWidget::~BrightnessMonitorWidget()
{
}

void BrightnessMonitorWidget::initUi()
{
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->setContentsMargins(10, 12, 10, 12);
    layout->setSpacing(6);

    m_sliderLayout->setContentsMargins(0, 0, 0, 0);
    m_sliderLayout->setSpacing(5);

    QList<BrightMonitor *> monitors = m_brightModel->monitors();
    for (BrightMonitor *monitor : monitors) {
        SliderContainer *container = new SliderContainer(CustomSlider::Normal, m_sliderWidget);
        container->setTitle(monitor->name());
        container->slider()->setIconSize(QSize(20, 20));
        container->slider()->setLeftIcon(QIcon(":/icons/resources/brightnesslow"));
        container->slider()->setRightIcon(QIcon(":/icons/resources/brightnesshigh"));
        container->setFixedHeight(50);
        m_sliderLayout->addWidget(container);

        m_sliderContainers << qMakePair(monitor, container);
    }

    layout->addSpacing(ITEMSPACE - layout->spacing());
    layout->addWidget(m_sliderWidget);
    layout->addSpacing(4);
    layout->addWidget(m_descriptionLabel);

    m_deviceList->setContentsMargins(0, 0, 0, 0);
    m_deviceList->setModel(m_model);
    m_deviceList->setViewMode(QListView::ListMode);
    m_deviceList->setMovement(QListView::Free);
    m_deviceList->setItemRadius(12);
    m_deviceList->setWordWrap(false);
    m_deviceList->verticalScrollBar()->setVisible(false);
    m_deviceList->horizontalScrollBar()->setVisible(false);
    m_deviceList->setOrientation(QListView::Flow::TopToBottom, false);
    layout->addWidget(m_deviceList);
    layout->addStretch();
    m_deviceList->setSpacing(10);

    m_deviceList->setItemDelegate(m_delegate);
}

void BrightnessMonitorWidget::initConnection()
{
    connect(m_delegate, &SettingDelegate::selectIndexChanged, this, [ this ](const QModelIndex &index) {
        BrightMonitor *monitor = index.data(itemDataRole).value<BrightMonitor *>();
        if (monitor) {
            m_deviceList->update();
            // 更新滚动条的内容
            onBrightChanged(monitor);
        } else {
            DDBusSender().service("com.deepin.dde.ControlCenter")
                    .path("/com/deepin/dde/ControlCenter")
                    .interface("com.deepin.dde.ControlCenter")
                    .method("ShowModule").arg(QString("display")).call();
            hide();
        }
    });

    for (QPair<BrightMonitor *, SliderContainer *> container : m_sliderContainers) {
        SliderContainer *slider = container.second;
        slider->slider()->setValue(container.first->brihtness());
        connect(slider->slider(), &CustomSlider::valueChanged, this, [ = ](int value) {
            m_brightModel->setBrightness(container.first, value);
        });
    }

    connect(m_brightModel, &BrightnessModel::brightnessChanged, this, &BrightnessMonitorWidget::onBrightChanged);
}

void BrightnessMonitorWidget::reloadMonitor()
{
    m_model->clear();
    // 跨端协同列表，后续会新增该功能
    // 显示设置
    DStandardItem *settingItem = new DStandardItem;
    settingItem->setIcon(QIcon(""));
    settingItem->setText(tr("Display settings"));
    settingItem->setFlags(Qt::NoItemFlags);
    settingItem->setData(false, itemCheckRole);
    m_model->appendRow(settingItem);
}

void BrightnessMonitorWidget::onBrightChanged(BrightMonitor *monitor)
{
    for (QPair<BrightMonitor *, SliderContainer *> container : m_sliderContainers) {
        SliderContainer *slider = container.second;
        if (container.first == monitor) {
            slider->slider()->blockSignals(true);
            slider->slider()->setValue(monitor->brihtness());
            slider->slider()->blockSignals(false);
        }
    }
}

void BrightnessMonitorWidget::resetHeight()
{
    int viewHeight = 0;
    for (int i = 0; i < m_model->rowCount(); i++) {
        QRect indexRect = m_deviceList->visualRect(m_model->index(i, 0));
        viewHeight += indexRect.height();
        // 上下间距
        viewHeight += m_deviceList->spacing() * 2;
    }
    // 设置列表的高度
    m_deviceList->setFixedHeight(viewHeight);
    QMargins sliderMargin = m_sliderLayout->contentsMargins();
    int sliderHeight = sliderMargin.top() + sliderMargin.bottom();
    for (QPair<BrightMonitor *, SliderContainer *> container : m_sliderContainers) {
        SliderContainer *slider = container.second;
        sliderHeight += slider->height();
    }

    m_sliderWidget->setFixedHeight(sliderHeight);
    QMargins m = layout()->contentsMargins();
    int space1 = ITEMSPACE - layout()->spacing();
    int space2 = 4;
    int height = m.top() + m.bottom() + sliderHeight + space1
            + m_descriptionLabel->height() + space2 + m_deviceList->height();
    setFixedHeight(height);
}
