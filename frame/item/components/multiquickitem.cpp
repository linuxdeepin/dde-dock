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
#include "multiquickitem.h"
#include "pluginsiteminterface.h"

#include <DFontSizeManager>
#include <DGuiApplicationHelper>
#include <DPaletteHelper>

#define BGSIZE 36
#define ICONWIDTH 24
#define ICONHEIGHT 24

static QSize expandSize = QSize(20, 20);

MultiQuickItem::MultiQuickItem(PluginsItemInterface *const pluginInter, QWidget *parent)
    : QuickSettingItem(pluginInter, parent)
    , m_selfDefine(false)
{
    initUi();
}

MultiQuickItem::~MultiQuickItem()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (itemWidget)
        itemWidget->setParent(nullptr);
}

QuickSettingItem::QuickSettingType MultiQuickItem::type() const
{
    return QuickSettingItem::QuickSettingType::Multi;
}

bool MultiQuickItem::eventFilter(QObject *obj, QEvent *event)
{
    if (m_selfDefine) {
        if (event->type() == QEvent::MouseButtonRelease) {
            if (obj->objectName() == "expandLabel") {
                // 如果是鼠标的按下事件
                QWidget *widget = pluginItem()->itemPopupApplet(QUICK_ITEM_KEY);
                if (!widget)
                    return QuickSettingItem::eventFilter(obj, event);

                Q_EMIT requestShowChildWidget(widget);

            } else if (obj == this) {
                const QString &command = pluginItem()->itemCommand(itemKey());
                if (!command.isEmpty())
                    QProcess::startDetached(command, QStringList());
            }
        } else if (event->type() == QEvent::Resize) {
            QLabel *labelWidget = qobject_cast<QLabel *>(obj);
            if (!labelWidget)
                return QuickSettingItem::eventFilter(obj, event);

            if (labelWidget->objectName() == "nameLabel") {
                labelWidget->setText(QFontMetrics(labelWidget->font()).elidedText(pluginItem()->pluginDisplayName(), Qt::TextElideMode::ElideRight, labelWidget->width()));
            } else if (labelWidget->objectName() == "stateLabel") {
                labelWidget->setText(QFontMetrics(labelWidget->font()).elidedText(pluginItem()->description(), Qt::TextElideMode::ElideRight, labelWidget->width()));
            }
        }
    }

    return QuickSettingItem::eventFilter(obj, event);
}

void MultiQuickItem::initUi()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (pluginItem()->icon(DockPart::QuickPanel).isNull() && itemWidget) {
        // 如果插件没有返回图标的显示，则获取插件的itemWidget
        QHBoxLayout *mainLayout = new QHBoxLayout(this);
        itemWidget->setParent(this);
        mainLayout->setContentsMargins(0, 0, 0, 0);
        mainLayout->addWidget(itemWidget);
    } else {
        // 如果插件获取到插件区域的图标，则让其按照图标来组合显示
        // 如果是占用两排的插件，则用横向Layout
        QHBoxLayout *mainLayout = new QHBoxLayout(this);
        mainLayout->setContentsMargins(10, 0, 10, 0);
        mainLayout->setSpacing(0);
        mainLayout->addStretch(10);
        mainLayout->setAlignment(Qt::AlignCenter);

        // 添加图标
        QWidget *iconWidgetParent = new QWidget(this);
        QVBoxLayout *iconLayout = new QVBoxLayout(iconWidgetParent);
        iconLayout->setContentsMargins(0, 0, 0, 0);
        iconLayout->setSpacing(0);
        iconLayout->setAlignment(Qt::AlignCenter);

        QWidget *iconWidget = new QuickIconWidget(pluginItem(), itemKey(), iconWidgetParent);
        iconWidget->setFixedSize(BGSIZE, BGSIZE);
        iconLayout->addWidget(iconWidget);
        mainLayout->addWidget(iconWidgetParent);
        mainLayout->addSpacing(10);

        // 添加中间的名称部分
        QWidget *textWidget = new QWidget(this);
        QLabel *nameLabel = new QLabel(textWidget);
        QLabel *stateLabel = new QLabel(textWidget);
        nameLabel->setObjectName("nameLabel");
        stateLabel->setObjectName("stateLabel");

        // 设置图标和文字的属性
        QFont nameFont = DFontSizeManager::instance()->t6();
        nameFont.setBold(true);
        QPalette pe;
        pe.setColor(QPalette::WindowText, Qt::black);
        nameLabel->setPalette(pe);
        stateLabel->setPalette(pe);
        nameLabel->setFont(nameFont);
        stateLabel->setFont(DFontSizeManager::instance()->t10());
        nameLabel->setText(pluginItem()->pluginDisplayName());
        stateLabel->setText(pluginItem()->description());
        nameLabel->installEventFilter(this);
        stateLabel->installEventFilter(this);

        QVBoxLayout *textLayout = new QVBoxLayout(textWidget);
        textLayout->setContentsMargins(0, 0, 0, 0);
        textLayout->setSpacing(0);
        textLayout->addWidget(nameLabel);
        textLayout->addWidget(stateLabel);
        textLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
        mainLayout->addWidget(textWidget);

        // 添加右侧的展开按钮
        QWidget *expandWidgetParent = new QWidget(this);
        QVBoxLayout *expandLayout = new QVBoxLayout(expandWidgetParent);
        expandLayout->setSpacing(0);
        QLabel *expandLabel = new QLabel(expandWidgetParent);
        expandLabel->setObjectName("expandLabel");
        expandLabel->setPixmap(QPixmap(expandFileName()));
        expandLabel->setFixedSize(expandSize);
        expandLabel->setAutoFillBackground(true);
        expandLabel->installEventFilter(this);
        expandLayout->addWidget(expandLabel);
        pe.setBrush(QPalette::Window, Qt::transparent);
        expandLabel->setPalette(pe);
        mainLayout->addWidget(expandWidgetParent);
        m_selfDefine = true;
    }
}

QString MultiQuickItem::expandFileName() const
{
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        return QString(":/icons/resources/arrow-right-dark.svg");

    return QString(":/icons/resources/arrow-right.svg");
}

/**
 * @brief QuickIconWidget::QuickIconWidget
 * @param pluginInter
 * @param parent
 * 图标的widget
 */
QuickIconWidget::QuickIconWidget(PluginsItemInterface *pluginInter, const QString &itemKey, QWidget *parent)
    : QWidget(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
{
}

void QuickIconWidget::paintEvent(QPaintEvent *event)
{
    QPixmap pixmapIcon = pluginIcon();
    if (pixmapIcon.isNull())
        return QWidget::paintEvent(event);

    pixmapIcon = pluginIcon(true);
    QPainter painter(this);
    painter.setRenderHint(QPainter::RenderHint::Antialiasing);
    painter.setPen(foregroundColor());

    DPalette dpa = DPaletteHelper::instance()->palette(this);
    QPainter pa(&pixmapIcon);
    pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
    pa.fillRect(pixmapIcon.rect(), painter.pen().brush());
    // 如果是主图标，则显示阴影背景
    painter.save();
    painter.setPen(Qt::NoPen);
    painter.setBrush(dpa.brush(DPalette::ColorRole::Midlight));
    painter.drawEllipse(rect());
    painter.restore();
    QRect rctIcon((rect().width() - pixmapIcon.width()) / 2, (rect().height() - pixmapIcon.height()) / 2, pixmapIcon.width(), pixmapIcon.height());
    painter.drawPixmap(rctIcon, pixmapIcon);
}

QColor QuickIconWidget::foregroundColor() const
{
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return dpa.color(DPalette::ColorGroup::Active, DPalette::ColorRole::Text);

    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Deactive)
        return dpa.color(DPalette::ColorGroup::Disabled, DPalette::ColorRole::Text);

    return dpa.color(DPalette::ColorGroup::Normal, DPalette::ColorRole::Text);
}

QPixmap QuickIconWidget::pluginIcon(bool contailGrab) const
{
    QIcon icon = m_pluginInter->icon(DockPart::QuickPanel);
    if (icon.isNull() && contailGrab) {
        // 如果图标为空，就使用itemWidget的截图作为它的图标，这种一般是适用于老版本插件或者没有实现v23接口的插件
        QWidget *itemWidget = m_pluginInter->itemWidget(m_itemKey);
        if (itemWidget) {
            itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
            return itemWidget->grab();
        }
        return QPixmap();
    }

    // 获取icon接口返回的图标
    int pixmapWidth = width();
    int pixmapHeight = height();
    QList<QSize> iconSizes = icon.availableSizes();
    if (iconSizes.size() > 0) {
        QSize size = iconSizes[0];
        if (size.isValid() && !size.isEmpty() && !size.isNull()) {
            pixmapWidth = size.width();
            pixmapHeight = size.height();
        }
    }

    return icon.pixmap(pixmapWidth, pixmapHeight);
}
