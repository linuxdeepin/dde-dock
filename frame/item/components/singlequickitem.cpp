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
#include "singlequickitem.h"
#include "pluginsiteminterface.h"

#include <DFontSizeManager>
#include <DGuiApplicationHelper>

#define ICONHEIGHT 24
#define ICONWIDTH 24
#define TEXTHEIGHT 11

SingleQuickItem::SingleQuickItem(PluginsItemInterface *const pluginInter, QWidget *parent)
    : QuickSettingItem(pluginInter, parent)
    , m_itemParentWidget(nullptr)
{
    initUi();
}

SingleQuickItem::~SingleQuickItem()
{
}

QuickSettingItem::QuickSettingType SingleQuickItem::type() const
{
    return QuickSettingItem::QuickSettingType::Single;
}

void SingleQuickItem::mouseReleaseEvent(QMouseEvent *event)
{
    Q_UNUSED(event);
    QStringList commandArgument = pluginItem()->itemCommand(itemKey()).split(" ");
    if (commandArgument.size() > 0) {
        QString command = commandArgument.first();
        commandArgument.removeFirst();
        QProcess::startDetached(command, commandArgument);
    }
}

void SingleQuickItem::resizeEvent(QResizeEvent *event)
{
    updateShow();
    QuickSettingItem::resizeEvent(event);
}

void SingleQuickItem::initUi()
{
    QWidget *topWidget = iconWidget(this);
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->addWidget(topWidget);
    installEventFilter(this);
}

QWidget *SingleQuickItem::iconWidget(QWidget *parent)
{
    // 显示图标的窗体
    QWidget *widget = new QWidget(parent);
    bool childIsEmpty = true;
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
            childIsEmpty = false;
        }
    }

    if (childIsEmpty) {
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
        labelText->setFixedHeight(TEXTHEIGHT);
        updatePluginName(labelText);
        labelText->setAlignment(Qt::AlignCenter);
        labelText->setFont(DFontSizeManager::instance()->t10());
        layout->addWidget(imageLabel);
        layout->addSpacing(7);
        layout->addWidget(labelText);
    }

    setProperty("paint", childIsEmpty);

    return widget;
}

QPixmap SingleQuickItem::pixmap() const
{
    // 如果快捷面板区域的图标为空，那么就获取itemWidget的截图
    QIcon icon = pluginItem()->icon(DockPart::QuickPanel);
    if (icon.isNull())
        return QPixmap();

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

QLabel *SingleQuickItem::findChildLabel(QWidget *parent, const QString &childObjectName) const
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

void SingleQuickItem::updatePluginName(QLabel *textLabel)
{
    if (!textLabel)
        return;

    QString text = pluginItem()->description();
    if (text.isEmpty())
        text = pluginItem()->pluginDisplayName();
    QFontMetrics ftm(textLabel->font());
    text = ftm.elidedText(text, Qt::TextElideMode::ElideRight, textLabel->width());
    textLabel->setText(text);
}

void SingleQuickItem::updateShow()
{
    if (property("paint").toBool()) {
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

void SingleQuickItem::detachPlugin()
{
    QWidget *itemWidget = pluginItem()->itemWidget(QUICK_ITEM_KEY);
    if (itemWidget && !property("paint").toBool())
        itemWidget->setParent(m_itemParentWidget);
}
