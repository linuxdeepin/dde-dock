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
#include "quicksettingitem.h"
#include "pluginsiteminterface.h"
#include "imageutil.h"

#include <DGuiApplicationHelper>
#include <DFontSizeManager>
#include <DPaletteHelper>

#include <QIcon>
#include <QPainterPath>
#include <QPushButton>
#include <QFontMetrics>

#define ICONWIDTH 24
#define ICONHEIGHT 24
#define ICONSPACE 10
#define RADIUS 8
#define FONTSIZE 10

#define BGWIDTH 128
#define BGSIZE 36
#define MARGINLEFTSPACE 10
#define OPENICONSIZE 12
#define MARGINRIGHTSPACE 9

static QSize expandSize = QSize(20, 20);

QuickSettingItem::QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, const QJsonObject &metaData, QWidget *parent)
    : DockItem(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
    , m_metaData(metaData)
    , m_iconWidgetParent(new QWidget(this))
    , m_iconWidget(new QuickIconWidget(pluginInter, itemKey, isPrimary(), m_iconWidgetParent))
    , m_textWidget(new QWidget(this))
    , m_nameLabel(new QLabel(m_textWidget))
    , m_stateLabel(new QLabel(m_textWidget))
{
    initUi();
    setAcceptDrops(true);
    this->installEventFilter(this);
}

QuickSettingItem::~QuickSettingItem()
{
}

bool QuickSettingItem::eventFilter(QObject *obj, QEvent *event)
{
    if (event->type() == QEvent::MouseButtonRelease) {
        if (obj->objectName() == "expandLabel") {
            // 如果是鼠标的按下事件
            if (isPrimary())
                Q_EMIT detailClicked(m_pluginInter);
        } else if (obj == this) {
            const QString &command = m_pluginInter->itemCommand(m_itemKey);
            if (!command.isEmpty())
                QProcess::startDetached(command);

            if (!isPrimary())
                Q_EMIT detailClicked(m_pluginInter);
        }
    } else if (event->type() == QEvent::Resize) {
        if (obj == m_nameLabel) {
            m_nameLabel->setText(QFontMetrics(m_nameLabel->font()).elidedText(m_pluginInter->pluginDisplayName(), Qt::TextElideMode::ElideRight, m_nameLabel->width()));
        } else if (obj == m_stateLabel) {
            m_stateLabel->setText(QFontMetrics(m_stateLabel->font()).elidedText(m_pluginInter->description(), Qt::TextElideMode::ElideRight, m_stateLabel->width()));
        }
    }

    return DockItem::eventFilter(obj, event);
}

PluginsItemInterface *QuickSettingItem::pluginItem() const
{
    return m_pluginInter;
}

DockItem::ItemType QuickSettingItem::itemType() const
{
    return DockItem::QuickSettingPlugin;
}

const QPixmap QuickSettingItem::dragPixmap()
{
    QPixmap pm = m_pluginInter->icon(DockPart::QuickPanel).pixmap(ICONWIDTH, ICONHEIGHT);

    QPainter pa(&pm);
    pa.setPen(foregroundColor());
    pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
    pa.fillRect(pm.rect(), foregroundColor());

    QPixmap pmRet(ICONWIDTH + ICONSPACE + FONTSIZE * 2, ICONHEIGHT + ICONSPACE + FONTSIZE * 2);
    pmRet.fill(Qt::transparent);
    QPainter paRet(&pmRet);
    paRet.drawPixmap(QPoint((ICONSPACE + FONTSIZE * 2) / 2, 0), pm);
    paRet.setPen(pa.pen());

    QFont ft;
    ft.setPixelSize(FONTSIZE);
    paRet.setFont(ft);
    QTextOption option;
    option.setAlignment(Qt::AlignTop | Qt::AlignHCenter);
    paRet.drawText(QRect(QPoint(0, ICONHEIGHT + ICONSPACE),
                           QPoint(pmRet.width(), pmRet.height())), m_pluginInter->pluginDisplayName(), option);
    return pmRet;
}

const QString QuickSettingItem::itemKey() const
{
    return m_itemKey;
}

bool QuickSettingItem::isPrimary() const
{
    if (m_metaData.contains("primary"))
        return m_metaData.value("primary").toBool();

    return false;
}

void QuickSettingItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
    QPainter painter(this);
    painter.setRenderHint(QPainter::RenderHint::Antialiasing);
    painter.setPen(foregroundColor());
    QPainterPath path;
    path.addRoundedRect(rect(), RADIUS, RADIUS);
    painter.setClipPath(path);
    // 绘制背景色
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    painter.fillRect(rect(), Qt::white);
}

QColor QuickSettingItem::foregroundColor() const
{
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return dpa.color(DPalette::ColorGroup::Active, DPalette::ColorRole::Text);

    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Deactive)
        return dpa.color(DPalette::ColorGroup::Disabled, DPalette::ColorRole::Text);

    return dpa.color(DPalette::ColorGroup::Normal, DPalette::ColorRole::Text);
}

void QuickSettingItem::initUi()
{
    if (isPrimary()) {
        // 如果是占用两排的插件，则用横向Layout
        QHBoxLayout *mainLayout = new QHBoxLayout(this);
        mainLayout->setContentsMargins(10, 0, 10, 0);
        mainLayout->setSpacing(0);
        mainLayout->addStretch(10);
        mainLayout->setAlignment(Qt::AlignCenter);
        // 添加图标
        QVBoxLayout *iconLayout = new QVBoxLayout(m_iconWidgetParent);
        iconLayout->setContentsMargins(0, 0, 0, 0);
        iconLayout->setSpacing(0);
        iconLayout->setAlignment(Qt::AlignCenter);
        m_iconWidget->setFixedSize(BGSIZE, BGSIZE);
        iconLayout->addWidget(m_iconWidget);
        mainLayout->addWidget(m_iconWidgetParent);
        mainLayout->addSpacing(10);
        // 添加中间的名称部分
        QFont nameFont = DFontSizeManager::instance()->t6();
        nameFont.setBold(true);
        QPalette pe;
        pe.setColor(QPalette::WindowText, Qt::black);
        m_nameLabel->setPalette(pe);
        m_stateLabel->setPalette(pe);
        m_nameLabel->setFont(nameFont);
        m_stateLabel->setFont(DFontSizeManager::instance()->t10());
        m_nameLabel->setText(m_pluginInter->pluginDisplayName());
        m_stateLabel->setText(m_pluginInter->description());
        m_nameLabel->installEventFilter(this);
        m_stateLabel->installEventFilter(this);
        QVBoxLayout *textLayout = new QVBoxLayout(m_textWidget);
        textLayout->setContentsMargins(0, 0, 0, 0);
        textLayout->setSpacing(0);
        textLayout->addWidget(m_nameLabel);
        textLayout->addWidget(m_stateLabel);
        textLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
        mainLayout->addWidget(m_textWidget);

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
    } else {
        QHBoxLayout *iconLayout = new QHBoxLayout(m_iconWidgetParent);
        iconLayout->setContentsMargins(0, 0, 0, 0);
        iconLayout->setSpacing(0);
        iconLayout->setAlignment(Qt::AlignHCenter);

        m_iconWidgetParent->setFixedHeight(ICONHEIGHT);
        m_iconWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
        iconLayout->addWidget(m_iconWidget);

        QVBoxLayout *mainLayout = new QVBoxLayout(this);
        mainLayout->setContentsMargins(0, 10, 0, 10);
        mainLayout->setSpacing(7);
        mainLayout->setAlignment(Qt::AlignCenter);
        // 添加上方的图标
        mainLayout->addWidget(m_iconWidgetParent);

        // 添加下方的文字
        QHBoxLayout *textLayout = new QHBoxLayout(m_textWidget);
        textLayout->setAlignment(Qt::AlignCenter);
        textLayout->setContentsMargins(0, 0, 0, 0);
        textLayout->setSpacing(0);
        QFont nameFont = DFontSizeManager::instance()->t10();
        QPalette pe;
        pe.setColor(QPalette::WindowText, Qt::black);
        m_nameLabel->setFont(nameFont);
        m_nameLabel->setPalette(pe);
        m_nameLabel->setText(m_pluginInter->pluginDisplayName());
        textLayout->addWidget(m_nameLabel);
        m_stateLabel->setVisible(false);
        m_textWidget->setFixedHeight(11);
        mainLayout->addWidget(m_textWidget);
        installEventFilter(this);
    }
}

QString QuickSettingItem::expandFileName()
{
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        return QString(":/icons/resources/arrow-right-dark.svg");

    return QString(":/icons/resources/arrow-right.svg");
}

QPixmap QuickSettingItem::pluginIcon() const
{
    QIcon icon = m_pluginInter->icon(DockPart::QuickPanel);
    if (icon.isNull()) {
        // 如果图标为空，就使用itemWidget的截图作为它的图标，这种一般是适用于老版本插件或者没有实现v23接口的插件
        QWidget *itemWidget = m_pluginInter->itemWidget(m_itemKey);
        itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
        QPixmap grabPixmap = itemWidget->grab();
        return grabPixmap;
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

/**
 * @brief QuickIconWidget::QuickIconWidget
 * @param pluginInter
 * @param parent
 * 图标的widget
 */
QuickIconWidget::QuickIconWidget(PluginsItemInterface *pluginInter, const QString &itemKey, bool isPrimary, QWidget *parent)
    : QWidget(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
    , m_isPrimary(isPrimary)
{
}

void QuickIconWidget::paintEvent(QPaintEvent *event)
{
    QWidget::paintEvent(event);
    QPixmap pm = pluginIcon();

    QPainter painter(this);
    painter.setRenderHint(QPainter::RenderHint::Antialiasing);
    painter.setPen(foregroundColor());

    if (m_isPrimary) {
        DPalette dpa = DPaletteHelper::instance()->palette(this);
        QPainter pa(&pm);
        pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
        pa.fillRect(pm.rect(), painter.pen().brush());
        // 如果是主图标，则显示阴影背景
        painter.save();
        painter.setPen(Qt::NoPen);
        painter.setBrush(dpa.brush(DPalette::ColorRole::Midlight));
        painter.drawEllipse(rect());
        painter.restore();
        QRect rctIcon((rect().width() - pm.width()) / 2, (rect().height() - pm.height()) / 2, pm.width(), pm.height());
        painter.drawPixmap(rctIcon, pm);
    } else {
        QRect rctIcon(0, 0, pm.width(), pm.height());
        painter.drawPixmap(rctIcon, pm);
    }
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

QPixmap QuickIconWidget::pluginIcon() const
{
    QIcon icon = m_pluginInter->icon(DockPart::QuickPanel);
    if (icon.isNull()) {
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
