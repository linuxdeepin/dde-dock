// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bluetoothmainwidget.h"
#include "bluetoothitem.h"
#include "adaptersmanager.h"
#include "adapter.h"

#include <DGuiApplicationHelper>
#include <DFontSizeManager>
#include <DStyle>

#include <QLabel>
#include <QHBoxLayout>
#include <QPainter>
#include <QPainterPath>
#include <QMouseEvent>

DGUI_USE_NAMESPACE
DWIDGET_USE_NAMESPACE

BluetoothMainWidget::BluetoothMainWidget(AdaptersManager *adapterManager, QWidget *parent)
    : QWidget(parent)
    , m_adapterManager(adapterManager)
    , m_iconButton(new DIconButton(this))
    , m_nameLabel(new QLabel(this))
    , m_stateLabel(new QLabel(this))
    , m_expandLabel(new QLabel(this))
    , m_mouseEnter(false)
{
    initUi();
    initConnection();
}

BluetoothMainWidget::~BluetoothMainWidget()
{
}

bool BluetoothMainWidget::eventFilter(QObject *watcher, QEvent *event)
{
    if (watcher != m_iconButton && event->type() == QEvent::MouseButtonRelease) {
        Q_EMIT requestExpand();
        return true;
    }
    if (watcher == m_nameLabel && event->type() == QEvent::Resize) {
        m_nameLabel->setText(QFontMetrics(m_nameLabel->font()).elidedText(tr("Bluetooth"), Qt::TextElideMode::ElideRight, m_nameLabel->width()));
    }
    if (watcher == m_stateLabel && event->type() == QEvent::Resize) {
        m_stateLabel->setText(QFontMetrics(m_stateLabel->font()).elidedText(m_stateLabel->text(), Qt::TextElideMode::ElideRight, m_stateLabel->width()));
    }
    if (watcher == m_iconButton && event->type() == QEvent::PaletteChange) {
        onPaletteChanged();
    }
    return QWidget::eventFilter(watcher, event);
}

void BluetoothMainWidget::initUi()
{
    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    // 添加左侧的图标
    m_iconButton->setEnabledCircle(true);
    m_iconButton->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Fixed);
    m_iconButton->setFocusPolicy(Qt::FocusPolicy::TabFocus);
    m_iconButton->setIconSize({24, 24});
    m_iconButton->setIcon(QIcon::fromTheme("bluetooth"));
    m_iconButton->setCheckable(true);
    m_iconButton->setChecked(isOpen());
    onPaletteChanged();

    // 添加中间的文本
    QWidget *textWidget = new QWidget(this);
    QVBoxLayout *textLayout = new QVBoxLayout(textWidget);
    textLayout->setContentsMargins(0, 10, 0, 10);
    textLayout->setSpacing(0);
    QFont nameFont = DFontSizeManager::instance()->t6();
    nameFont.setBold(true);

    m_nameLabel->setParent(textWidget);
    m_nameLabel->setFont(nameFont);

    m_stateLabel->setParent(textWidget);
    m_stateLabel->setFont(DFontSizeManager::instance()->t10());
    m_stateLabel->setFixedWidth(76);
    textLayout->addWidget(m_nameLabel);
    textLayout->addWidget(m_stateLabel);

    // 添加右侧的展开按钮
    QWidget *expandWidget = new QWidget(this);
    QVBoxLayout *expandLayout = new QVBoxLayout(expandWidget);
    expandLayout->setContentsMargins(0, 0, 0, 0);
    expandLayout->setSpacing(0);
    expandLayout->addWidget(m_expandLabel);

    // 设置图标和文本
    m_nameLabel->setText(QFontMetrics{m_nameLabel->font()}.elidedText(tr("Bluetooth"), Qt::TextElideMode::ElideRight, m_nameLabel->width()));
    m_stateLabel->setText(QFontMetrics{m_stateLabel->font()}.elidedText(isOpen() ? tr("Turn on") : tr("Turn off"), Qt::TextElideMode::ElideRight,m_stateLabel->width()));
    updateExpandIcon();

    // 将所有的窗体都添加到主布局中
    mainLayout->setContentsMargins(10, 0, 10, 0);
    mainLayout->setSpacing(0);
    mainLayout->addWidget(m_iconButton);
    mainLayout->addSpacing(10);
    mainLayout->addWidget(textWidget);
    mainLayout->addStretch();
    mainLayout->addWidget(expandWidget);

    m_expandLabel->installEventFilter(this);
    m_nameLabel->installEventFilter(this);
    m_iconButton->installEventFilter(this);
}

void BluetoothMainWidget::initConnection()
{
    connect(m_adapterManager, &AdaptersManager::adapterIncreased, this, &BluetoothMainWidget::onAdapterChanged);
    connect(m_adapterManager, &AdaptersManager::adapterDecreased, this, &BluetoothMainWidget::onAdapterChanged);
    connect(m_adapterManager, &AdaptersManager::adapterIncreased, this, [ = ](Adapter *adapter) {
        connect(adapter, &Adapter::poweredChanged, this, &BluetoothMainWidget::onAdapterChanged);
    });

    for (const Adapter *adapter : m_adapterManager->adapters())
        connect(adapter, &Adapter::poweredChanged, this, &BluetoothMainWidget::onAdapterChanged);

    connect(m_iconButton, &DIconButton::clicked, this, [this](){
        bool status = !(isOpen());
        for (const Adapter *adapter : m_adapterManager->adapters())
            m_adapterManager->setAdapterPowered(adapter, status);
    });

    onAdapterChanged();
}

void BluetoothMainWidget::updateExpandIcon()
{
    QPixmap enterPixmap = DStyle::standardIcon(style(), DStyle::SP_ArrowEnter).pixmap(QSize(16, 16));
    m_expandLabel->setPixmap(enterPixmap);
}

bool BluetoothMainWidget::isOpen() const
{
    QList<const Adapter *> adapters = m_adapterManager->adapters();
    for (const Adapter *adapter : adapters) {
        if (adapter->powered())
            return true;
    }

    return false;
}

QString BluetoothMainWidget::bluetoothIcon(bool isOpen) const
{
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
        return isOpen ? ":/bluetooth-active-symbolic-dark.svg" : ":/bluetooth-disable-symbolic-dark.svg";

    return isOpen ? ":/bluetooth-active-symbolic.svg" : ":/bluetooth-disable-symbolic.svg";
}

void BluetoothMainWidget::onAdapterChanged()
{
    bool bluetoothIsOpen = isOpen();
    const QString& text = bluetoothIsOpen ? tr("Turn on") : tr("Turn off");
    QFontMetrics fmt{m_stateLabel->font()};
    m_stateLabel->setText(fmt.elidedText(text, Qt::TextElideMode::ElideRight,m_stateLabel->width()));
    m_iconButton->setChecked(bluetoothIsOpen);
}

void BluetoothMainWidget::onPaletteChanged()
{
    if (!m_iconButton)
        return;

    auto pa = m_iconButton->palette();
    pa.setColor(QPalette::HighlightedText, pa.color(QPalette::Highlight));
    m_iconButton->setPalette(pa);
    m_iconButton->update();
}
