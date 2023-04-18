// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "largerquickitem.h"
#include "pluginsiteminterface.h"

#include <DFontSizeManager>
#include <DGuiApplicationHelper>
#include <DPaletteHelper>

#include <QLabel>
#include <QHBoxLayout>
#include <QPainter>
#include <QMouseEvent>

#define BGSIZE 36
#define ICONWIDTH 24
#define ICONHEIGHT 24

DWIDGET_USE_NAMESPACE

static QSize expandSize = QSize(20, 20);

LargerQuickItem::LargerQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : QuickSettingItem(pluginInter, itemKey, parent)
    , m_iconWidget(nullptr)
    , m_nameLabel(nullptr)
    , m_stateLabel(nullptr)
    , m_itemWidgetParent(nullptr)
{
    initUi();
}

LargerQuickItem::~LargerQuickItem()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (itemWidget)
        itemWidget->setParent(nullptr);
}

void LargerQuickItem::doUpdate()
{
    if (m_iconWidget && m_nameLabel && m_stateLabel) {
        m_iconWidget->update();
        m_nameLabel->setText(QFontMetrics(m_nameLabel->font()).elidedText(pluginItem()->pluginDisplayName(), Qt::TextElideMode::ElideRight, m_nameLabel->width()));
        m_stateLabel->setText(QFontMetrics(m_stateLabel->font()).elidedText(pluginItem()->description(), Qt::TextElideMode::ElideRight, m_stateLabel->width()));
    } else {
        QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
        if (itemWidget) {
            // 如果插件没有返回图标的显示，则获取插件的itemWidget
            itemWidget->update();
        }
    }
}

void LargerQuickItem::detachPlugin()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (itemWidget && itemWidget->parentWidget() == this)
        itemWidget->setParent(m_itemWidgetParent);
}

QuickSettingItem::QuickItemStyle LargerQuickItem::type() const
{
    return QuickSettingItem::QuickItemStyle::Larger;
}

bool LargerQuickItem::eventFilter(QObject *obj, QEvent *event)
{
    if (!m_iconWidget)
        return QuickSettingItem::eventFilter(obj, event);

    switch (event->type()) {
    case QEvent::MouseButtonRelease: {
        QMouseEvent *mouseevent = static_cast<QMouseEvent *>(event);
        if (mouseevent->button() != Qt::LeftButton) {
            return QuickSettingItem::eventFilter(obj, event);
        }
        if (obj->objectName() == "expandLabel") {
            // 如果是鼠标的按下事件
            QWidget *widget = pluginItem()->itemPopupApplet(QUICK_ITEM_KEY);
            if (widget)
                Q_EMIT requestShowChildWidget(widget);
            break;
        }
        if (obj == this) {
            QStringList commandArgumend = pluginItem()->itemCommand(itemKey()).split(" ");
            if (commandArgumend.size() == 0)
                break;

            QString command = commandArgumend.first();
            commandArgumend.removeFirst();
            QProcess::startDetached(command, commandArgumend);
        }
        break;
    }
    case QEvent::Resize: {
        QLabel *labelWidget = qobject_cast<QLabel *>(obj);
        if (!labelWidget)
            break;

        if (labelWidget == m_nameLabel) {
            labelWidget->setText(QFontMetrics(labelWidget->font()).elidedText(pluginItem()->pluginDisplayName(), Qt::TextElideMode::ElideRight, labelWidget->width()));
            break;
        }
        if (labelWidget == m_stateLabel) {
            labelWidget->setText(QFontMetrics(labelWidget->font()).elidedText(pluginItem()->description(), Qt::TextElideMode::ElideRight, labelWidget->width()));
        }
        break;
    }
    default:
        break;
    }

    return QuickSettingItem::eventFilter(obj, event);
}

void LargerQuickItem::showEvent(QShowEvent *event)
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (pluginItem()->icon(DockPart::QuickPanel).isNull() && itemWidget) {
        itemWidget->setParent(this);
        itemWidget->setVisible(true);
    }
}

void LargerQuickItem::resizeEvent(QResizeEvent *event)
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (pluginItem()->icon(DockPart::QuickPanel).isNull() && itemWidget) {
        itemWidget->setFixedSize(size());
    }

    QuickSettingItem::resizeEvent(event);
}

void LargerQuickItem::initUi()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (pluginItem()->icon(DockPart::QuickPanel).isNull() && itemWidget) {
        m_itemWidgetParent = itemWidget->parentWidget();
        // 如果插件没有返回图标的显示，则获取插件的itemWidget
        QHBoxLayout *mainLayout = new QHBoxLayout(this);
        itemWidget->setParent(this);
        mainLayout->setContentsMargins(0, 0, 0, 0);
        mainLayout->addWidget(itemWidget);
        itemWidget->setVisible(true);
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

        m_iconWidget = new QuickIconWidget(pluginItem(), itemKey(), iconWidgetParent);
        m_iconWidget->setFixedSize(BGSIZE, BGSIZE);
        iconLayout->addWidget(m_iconWidget);
        mainLayout->addWidget(iconWidgetParent);
        mainLayout->addSpacing(10);

        // 添加中间的名称部分
        QWidget *textWidget = new QWidget(this);
        m_nameLabel = new QLabel(textWidget);
        m_stateLabel = new QLabel(textWidget);
        m_nameLabel->setObjectName("nameLabel");
        m_stateLabel->setObjectName("stateLabel");

        // 设置图标和文字的属性
        QFont nameFont = DFontSizeManager::instance()->t6();
        nameFont.setBold(true);
        QPalette pe;
        pe.setColor(QPalette::WindowText, Qt::black);
        m_nameLabel->setPalette(pe);
        m_stateLabel->setPalette(pe);
        m_nameLabel->setFont(nameFont);
        m_stateLabel->setFont(DFontSizeManager::instance()->t10());
        m_nameLabel->setText(pluginItem()->pluginDisplayName());
        m_stateLabel->setText(pluginItem()->description());
        m_nameLabel->installEventFilter(this);
        m_stateLabel->installEventFilter(this);

        QVBoxLayout *textLayout = new QVBoxLayout(textWidget);
        textLayout->setContentsMargins(0, 0, 0, 0);
        textLayout->setSpacing(0);
        textLayout->addWidget(m_nameLabel);
        textLayout->addWidget(m_stateLabel);
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
    }
}

QString LargerQuickItem::expandFileName() const
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
    // TODO 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginMode::Active)
        return dpa.color(DPalette::ColorGroup::Active, DPalette::ColorRole::Text);

    if (m_pluginInter->status() == PluginsItemInterface::PluginMode::Deactive)
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
