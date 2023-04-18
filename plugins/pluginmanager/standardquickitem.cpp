// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "standardquickitem.h"
#include "pluginsiteminterface.h"

#include <DFontSizeManager>
#include <DGuiApplicationHelper>

#include <QLabel>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QMouseEvent>

static constexpr int ICONHEIGHT = 24;
static constexpr int ICONWIDTH = 24;
static constexpr int TEXTHEIGHT = 20;

DWIDGET_USE_NAMESPACE

StandardQuickItem::StandardQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : QuickSettingItem(pluginInter, itemKey, parent)
    , m_itemParentWidget(nullptr)
    , m_needPaint(true)
{
    initUi();
}

StandardQuickItem::~StandardQuickItem()
{
}

QuickSettingItem::QuickItemStyle StandardQuickItem::type() const
{
    return QuickSettingItem::QuickItemStyle::Standard;
}

void StandardQuickItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton) {
        return;
    }
    QStringList commandArgument = pluginItem()->itemCommand(itemKey()).split(" ");
    if (commandArgument.size() > 0) {
        QString command = commandArgument.first();
        commandArgument.removeFirst();
        QProcess::startDetached(command, commandArgument);
    }
}

void StandardQuickItem::resizeEvent(QResizeEvent *event)
{
    doUpdate();
    QuickSettingItem::resizeEvent(event);
}

void StandardQuickItem::initUi()
{
    QWidget *topWidget = iconWidget(this);
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->addWidget(topWidget);
    installEventFilter(this);
}

QWidget *StandardQuickItem::iconWidget(QWidget *parent)
{
    // 显示图标的窗体
    QWidget *widget = new QWidget(parent);
    m_needPaint = true;
    QIcon icon = pluginItem()->icon(DockPart::QuickPanel);
    if (icon.isNull()) {
        // 如果图标为空，则将获取itemWidget作为它的显示
        QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
        if (itemWidget) {
            m_itemParentWidget = itemWidget->parentWidget();
            QHBoxLayout *layout = new QHBoxLayout(widget);
            layout->setContentsMargins(0, 0, 0 ,0);
            itemWidget->setParent(widget);
            layout->addWidget(itemWidget);
            itemWidget->setVisible(true);
            m_needPaint = false;
        }
    }

    if (m_needPaint) {
        // 如果没有子窗体，则需要添加下方的文字
        QVBoxLayout *layout = new QVBoxLayout(widget);
        layout->setAlignment(Qt::AlignVCenter);
        layout->setContentsMargins(0, 0, 0, 0);
        layout->setSpacing(0);
        QLabel *imageLabel = new QLabel(widget);
        imageLabel->setObjectName("imageLabel");
        imageLabel->setFixedHeight(ICONHEIGHT);
        imageLabel->setAlignment(Qt::AlignCenter);

        QLabel *labelText = new QLabel(widget);
        labelText->setObjectName("textLabel");
        labelText->setAlignment(Qt::AlignCenter);
        labelText->setFont(DFontSizeManager::instance()->t10());
        labelText->setFixedHeight(TEXTHEIGHT);
        labelText->setFixedWidth(70);
        updatePluginName(labelText);
        layout->addWidget(imageLabel);
        layout->addSpacing(7);
        layout->addWidget(labelText);
    }

    return widget;
}

QPixmap StandardQuickItem::pixmap() const
{
    // 如果快捷面板区域的图标为空，那么就获取itemWidget的截图
    QIcon icon = pluginItem()->icon(DockPart::QuickPanel);
    if (icon.isNull())
        return QPixmap();

    int pixmapWidth = ICONWIDTH;
    int pixmapHeight = ICONHEIGHT;
    QList<QSize> iconSizes = icon.availableSizes();
    if (iconSizes.size() > 0) {
        QSize size = iconSizes[0];
        if (size.isValid() && !size.isEmpty() && !size.isNull()) {
            pixmapWidth = size.width();
            pixmapHeight = size.height();
        }
    }

    return icon.pixmap(pixmapWidth / qApp->devicePixelRatio(), pixmapHeight / qApp->devicePixelRatio());
}

QLabel *StandardQuickItem::findChildLabel(QWidget *parent, const QString &childObjectName) const
{
    QList<QObject *> childrends = parent->children();
    for (QObject *child : childrends) {
        QWidget *widget = qobject_cast<QWidget *>(child);
        if (!widget)
            continue;

        QLabel *label = qobject_cast<QLabel *>(child);
        if (label && widget->objectName() == childObjectName)
            return label;

        label = findChildLabel(widget, childObjectName);
        if (label)
            return label;
    }

    return nullptr;
}

void StandardQuickItem::updatePluginName(QLabel *textLabel)
{
    if (!textLabel)
        return;
    QString text = pluginItem()->description();
    if (text.isEmpty())
        text = pluginItem()->pluginDisplayName();
    QFontMetrics ftm(textLabel->font());
    if (ftm.boundingRect(text).width() > 70) {
        this->setToolTip(text);
    } else {
        this->setToolTip("");
    }
    text = ftm.elidedText(text, Qt::TextElideMode::ElideMiddle, 70);
    textLabel->setText(text);
    qInfo() << "text update to: " << text;
}

void StandardQuickItem::doUpdate()
{
    if (m_needPaint) {
        QLabel *imageLabel = findChildLabel(this, "imageLabel");
        if (imageLabel) {
            // 更新图像
            imageLabel->setPixmap(pixmap());
        }
        updatePluginName(findChildLabel(this, "textLabel"));
    } else {
        QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
        if (itemWidget)
            itemWidget->update();
    }
}

void StandardQuickItem::detachPlugin()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (itemWidget && !m_needPaint)
        itemWidget->setParent(m_itemParentWidget);
}
