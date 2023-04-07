// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "systempluginitem.h"
#include "utils.h"
#include "dockpopupwindow.h"

#include <QProcess>
#include <QDebug>
#include <QPointer>
#include <QMenu>

#include <xcb/xproto.h>

Dock::Position SystemPluginItem::DockPosition = Dock::Position::Bottom;
QPointer<DockPopupWindow> SystemPluginItem::PopupWindow = nullptr;

SystemPluginItem::SystemPluginItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : BaseTrayWidget(parent)
    , m_popupShown(false)
    , m_tapAndHold(false)
    , m_contextMenu(new QMenu)
    , m_pluginInter(pluginInter)
    , m_centralWidget(m_pluginInter->itemWidget(itemKey))
    , m_popupTipsDelayTimer(new QTimer(this))
    , m_popupAdjustDelayTimer(new QTimer(this))
    , m_itemKey(itemKey)
    , m_gsettings(Utils::ModuleSettingsPtr(pluginInter->pluginName(), QByteArray(), this))
{
    qDebug() << "load tray plugins item: " << m_pluginInter->pluginName() << itemKey << m_centralWidget;

    QIcon icon = m_pluginInter->icon(DockPart::SystemPanel);
    if (m_centralWidget) {
        if (icon.isNull()) {
            // 此处先创建一个Layout，在该窗体show的时候将m_centralWidget添加到layout上面
            QBoxLayout *hLayout = new QHBoxLayout(this);
            hLayout->setSpacing(0);
            hLayout->setMargin(0);
            setLayout(hLayout);
            m_centralWidget->installEventFilter(this);
        } else {
            m_centralWidget->setVisible(false);
        }
    }

    setAccessibleName(m_itemKey);
    setAttribute(Qt::WA_TranslucentBackground);

    if (PopupWindow.isNull()) {
        DockPopupWindow *arrowRectangle = new DockPopupWindow(nullptr);
        arrowRectangle->setRadius(6);
        arrowRectangle->setObjectName("systemtraypopup");
        if (Utils::IS_WAYLAND_DISPLAY) {
            Qt::WindowFlags flags = arrowRectangle->windowFlags() | Qt::FramelessWindowHint;
            arrowRectangle->setWindowFlags(flags);
        }
        PopupWindow = arrowRectangle;
        connect(qApp, &QApplication::aboutToQuit, PopupWindow, &DockPopupWindow::deleteLater);
    }

    if (Utils::IS_WAYLAND_DISPLAY) {
        Qt::WindowFlags flags = m_contextMenu->windowFlags() | Qt::FramelessWindowHint;
        m_contextMenu->setWindowFlags(flags);
    }

    // 必须初始化父窗口，否则当主题切换之后再设置父窗口的时候palette会更改为主题切换前的palette
    if (QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey)) {
        w->setParent(PopupWindow.data());
        w->setVisible(false);
    }

    m_popupTipsDelayTimer->setInterval(500);
    m_popupTipsDelayTimer->setSingleShot(true);

    m_popupAdjustDelayTimer->setInterval(10);
    m_popupAdjustDelayTimer->setSingleShot(true);

    installEventFilter(this);

    connect(m_popupTipsDelayTimer, &QTimer::timeout, this, &SystemPluginItem::showHoverTips);
    connect(m_popupAdjustDelayTimer, &QTimer::timeout, this, &SystemPluginItem::updatePopupPosition, Qt::QueuedConnection);
    connect(m_contextMenu, &QMenu::triggered, this, &SystemPluginItem::menuActionClicked);

    if (m_gsettings)
        connect(m_gsettings, &QGSettings::changed, this, &SystemPluginItem::onGSettingsChanged);

    grabGesture(Qt::TapAndHoldGesture);
}

SystemPluginItem::~SystemPluginItem()
{
    if (m_popupShown)
        popupWindowAccept();
    m_contextMenu->deleteLater();
}

QString SystemPluginItem::itemKeyForConfig()
{
    return m_itemKey;
}

void SystemPluginItem::updateIcon()
{
    m_pluginInter->refreshIcon(m_itemKey);
}

void SystemPluginItem::sendClick(uint8_t mouseButton, int x, int y)
{
    Q_UNUSED(mouseButton);
    Q_UNUSED(x);
    Q_UNUSED(y);

    // do not process this callback
    // handle all mouse event in override mouse function
}

QWidget *SystemPluginItem::trayTipsWidget()
{
    if (m_pluginInter->itemTipsWidget(m_itemKey)) {
        m_pluginInter->itemTipsWidget(m_itemKey)->setAccessibleName(m_pluginInter->pluginName());
    }

    return m_pluginInter->itemTipsWidget(m_itemKey);
}

QWidget *SystemPluginItem::trayPopupApplet()
{
    if (m_pluginInter->itemPopupApplet(m_itemKey)) {
        m_pluginInter->itemPopupApplet(m_itemKey)->setAccessibleName(m_pluginInter->pluginName());
    }

    return m_pluginInter->itemPopupApplet(m_itemKey);
}

const QString SystemPluginItem::trayClickCommand()
{
    return m_pluginInter->itemCommand(m_itemKey);
}

const QString SystemPluginItem::contextMenu() const
{
    return m_pluginInter->itemContextMenu(m_itemKey);
}

void SystemPluginItem::invokedMenuItem(const QString &menuId, const bool checked)
{
    m_pluginInter->invokedMenuItem(m_itemKey, menuId, checked);
}

QWidget *SystemPluginItem::centralWidget() const
{
    return m_centralWidget;
}

void SystemPluginItem::detachPluginWidget()
{
    QWidget *widget = m_pluginInter->itemWidget(m_itemKey);
    if (widget)
        widget->setParent(nullptr);
}

bool SystemPluginItem::event(QEvent *event)
{
    if (m_popupShown) {
        switch (event->type()) {
        case QEvent::Paint:
            if (!m_popupAdjustDelayTimer->isActive())
                m_popupAdjustDelayTimer->start();
            break;
        default:;
        }
    }

    if (event->type() == QEvent::Gesture)
        gestureEvent(static_cast<QGestureEvent *>(event));

    return BaseTrayWidget::event(event);
}

bool SystemPluginItem::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == this && event->type() == QEvent::DeferredDelete
            && m_centralWidget && m_centralWidget->parentWidget()) {
        m_centralWidget->setParent(nullptr);
        m_centralWidget->setVisible(false);
    }

    return BaseTrayWidget::eventFilter(watched, event);
}

void SystemPluginItem::enterEvent(QEvent *event)
{
    if (checkGSettingsControl()) {
        //网络需要显示Tips，需要特殊处理。
        if (m_pluginInter->pluginName() != "network")
            return;
    }

    // 触屏不显示hover效果
    if (!qApp->property(IS_TOUCH_STATE).toBool()) {
        m_popupTipsDelayTimer->start();
    }
    update();

    BaseTrayWidget::enterEvent(event);
}

void SystemPluginItem::leaveEvent(QEvent *event)
{
    m_popupTipsDelayTimer->stop();

    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();

    update();

    BaseTrayWidget::leaveEvent(event);
}

void SystemPluginItem::mousePressEvent(QMouseEvent *event)
{
    if (checkGSettingsControl()) {
        return;
    }

    m_popupTipsDelayTimer->stop();
    hideNonModel();

    if (event->button() == Qt::RightButton
            && perfectIconRect().contains(event->pos(), true)) {
        return (m_gsettings && (!m_gsettings->keys().contains("menuEnable") || m_gsettings->get("menuEnable").toBool())) ? showContextMenu() : void();
    }

    BaseTrayWidget::mousePressEvent(event);
}

void SystemPluginItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (checkGSettingsControl()) {
        return;
    }

    if (event->button() != Qt::LeftButton) {
        return;
    }

    if (checkAndResetTapHoldGestureState() && event->source() == Qt::MouseEventSynthesizedByQt) {
        qDebug() << "SystemTray: tap and hold gesture detected, ignore the synthesized mouse release event";
        return;
    }

    event->accept();

    showPopupApplet(trayPopupApplet());

    if (!trayClickCommand().isEmpty()) {
        QProcess::startDetached(trayClickCommand(), {});
    }

    BaseTrayWidget::mouseReleaseEvent(event);
}

void SystemPluginItem::showEvent(QShowEvent *event)
{
    QTimer::singleShot(0, this, [ = ] {
        onGSettingsChanged("enable");
    });

    return BaseTrayWidget::showEvent(event);
}

void SystemPluginItem::hideEvent(QHideEvent *event)
{
    if (m_pluginInter->icon(DockPart::SystemPanel).isNull()
            && m_centralWidget && m_centralWidget->parentWidget() == this) {
        layout()->removeWidget(m_centralWidget);
        m_centralWidget->setParent(nullptr);
        m_centralWidget->setVisible(false);
    }
    BaseTrayWidget::hideEvent(event);
}

void SystemPluginItem::paintEvent(QPaintEvent *event)
{
    QIcon icon = m_pluginInter->icon(DockPart::SystemPanel);
    if (icon.isNull()) {
        showCentralWidget();
        return BaseTrayWidget::paintEvent(event);
    }

    int x = 0;
    int y = 0;
    QSize iconSize = this->size();
    QList<QSize> sizes = icon.availableSizes();
    if (sizes.size() > 0) {
        QSize availableSize = sizes.first();
        if (iconSize.width() > availableSize.width()) {
            x = (iconSize.width() - availableSize.width()) / 2;
            y = (iconSize.height() - availableSize.height()) / 2;
            iconSize = availableSize;
        }
    }

    QPixmap pixmap = icon.pixmap(iconSize);
    QPainter painter(this);
    painter.drawPixmap(QRect(QPoint(x, y), iconSize), pixmap);
}

const QPoint SystemPluginItem::popupMarkPoint() const
{
    QPoint p(topleftPoint());

    const QRect r = rect();
    const QRect wr = window()->rect();

    switch (DockPosition) {
    case Dock::Position::Top:
        p += QPoint(r.width() / 2, r.height() + (wr.height() - r.height()) / 2);
        break;
    case Dock::Position::Bottom:
        p += QPoint(r.width() / 2, 0 - (wr.height() - r.height()) / 2);
        break;
    case Dock::Position::Left:
        p += QPoint(r.width() + (wr.width() - r.width()) / 2, r.height() / 2);
        break;
    case Dock::Position::Right:
        p += QPoint(0 - (wr.width() - r.width()) / 2, r.height() / 2);
        break;
    }

    return p;
}

// 获取在最外层的窗口(MainWindow)中的位置
const QPoint SystemPluginItem::topleftPoint() const
{
    QPoint p;
    const QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    return p;
}

void SystemPluginItem::hidePopup()
{
    m_popupTipsDelayTimer->stop();
    m_popupAdjustDelayTimer->stop();
    m_popupShown = false;
    PopupWindow->hide();

    DockPopupWindow *popup = PopupWindow.data();
    QWidget *content = popup->getContent();
    if (content) {
        content->setVisible(false);
    }

    emit PopupWindow->accept();
    emit requestWindowAutoHide(true);
}

bool SystemPluginItem::containsPoint(const QPoint& pos)
{
    QPoint ptGlobal = mapToGlobal(QPoint(0, 0));
    QRect rectGlobal(ptGlobal, this->size());
    if (rectGlobal.contains(pos))
        return true;

    // 如果菜单列表隐藏，则认为不在区域内
    if (!m_contextMenu->isVisible())
        return false;

    // 判断鼠标是否在菜单区域
    return m_contextMenu->geometry().contains(pos);
}

void SystemPluginItem::hideNonModel()
{
    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();
}

void SystemPluginItem::popupWindowAccept()
{
    if (!PopupWindow->isVisible())
        return;

    disconnect(PopupWindow.data(), &DockPopupWindow::accept, this, &SystemPluginItem::popupWindowAccept);

    hidePopup();
}

void SystemPluginItem::showPopupApplet(QWidget *const applet)
{
    if (!applet)
        return;

    // another model popup window already exists
    if (PopupWindow->model()) {
        applet->setVisible(false);
        return;
    }

    showPopupWindow(applet, true);
}

void SystemPluginItem::showPopupWindow(QWidget *const content, const bool model)
{
    m_popupShown = true;
    m_lastPopupWidget = content;

    if (model)
        emit requestWindowAutoHide(false);

    DockPopupWindow *popup = PopupWindow.data();
    QWidget *lastContent = popup->getContent();
    if (lastContent)
        lastContent->setVisible(false);

    popup->setPosition(DockPosition);
    popup->resize(content->sizeHint());
    popup->setContent(content);

    QPoint p = popupMarkPoint();
    if (!popup->isVisible())
        QMetaObject::invokeMethod(popup, "show", Qt::QueuedConnection, Q_ARG(QPoint, p), Q_ARG(bool, model));
    else
        popup->show(p, model);

    connect(popup, &DockPopupWindow::accept, this, &SystemPluginItem::popupWindowAccept, Qt::UniqueConnection);
}

void SystemPluginItem::showHoverTips()
{
    // another model popup window already exists
    if (PopupWindow->model())
        return;

    // if not in geometry area
    const QRect r(topleftPoint(), size());
    if (!r.contains(QCursor::pos()))
        return;

    QWidget *const content = trayTipsWidget();
    if (!content)
        return;

    showPopupWindow(content);
}

/*!
 * \sa DockItem::checkAndResetTapHoldGestureState
 */
bool SystemPluginItem::checkAndResetTapHoldGestureState()
{
    bool ret = m_tapAndHold;
    m_tapAndHold = false;
    return ret;
}

void SystemPluginItem::gestureEvent(QGestureEvent *event)
{
    if (!event)
        return;

    QGesture *gesture = event->gesture(Qt::TapAndHoldGesture);

    if (!gesture)
        return;

    qDebug() << "SystemTray: got TapAndHoldGesture";

    m_tapAndHold = true;
}

void SystemPluginItem::showContextMenu()
{
    const QString menuJson = contextMenu();
    if (menuJson.isEmpty())
        return;

    QJsonDocument jsonDocument = QJsonDocument::fromJson(menuJson.toLocal8Bit().data());
    if (jsonDocument.isNull())
        return;

    QJsonObject jsonMenu = jsonDocument.object();

    qDeleteAll(m_contextMenu->actions());

    QJsonArray jsonMenuItems = jsonMenu.value("items").toArray();
    for (auto item : jsonMenuItems) {
        QJsonObject itemObj = item.toObject();
        QAction *action = new QAction(itemObj.value("itemText").toString(), m_contextMenu);
        action->setCheckable(itemObj.value("isCheckable").toBool());
        action->setChecked(itemObj.value("checked").toBool());
        action->setData(itemObj.value("itemId").toString());
        action->setEnabled(itemObj.value("isActive").toBool());
        m_contextMenu->addAction(action);
    }

    hidePopup();
    emit requestWindowAutoHide(false);

    m_contextMenu->exec(QCursor::pos());

    onContextMenuAccepted();
}

void SystemPluginItem::menuActionClicked(QAction *action)
{
    invokedMenuItem(action->data().toString(), true);
    Q_EMIT execActionFinished();
}

void SystemPluginItem::showCentralWidget()
{
    if (!m_pluginInter->icon(DockPart::SystemPanel).isNull() || !m_centralWidget)
        return;

    m_centralWidget->setParent(this);
    m_centralWidget->setVisible(true);
    layout()->addWidget(m_centralWidget);
}

void SystemPluginItem::onContextMenuAccepted()
{
    emit requestRefershWindowVisible();
    emit requestWindowAutoHide(true);
}

void SystemPluginItem::updatePopupPosition()
{
    Q_ASSERT(sender() == m_popupAdjustDelayTimer);

    if (!m_popupShown || !PopupWindow->model())
        return;

    if (PopupWindow->getContent() != m_lastPopupWidget.data())
        return popupWindowAccept();

    const QPoint p = popupMarkPoint();
    PopupWindow->show(p, PopupWindow->model());
}

void SystemPluginItem::onGSettingsChanged(const QString &key) {
    if (key != "enable") {
        return;
    }

    if (m_gsettings && m_gsettings->keys().contains("enable")) {
        const bool visible = m_gsettings->get("enable").toBool();
        setVisible(visible);
        emit itemVisibleChanged(visible);
    }
}

bool SystemPluginItem::checkGSettingsControl() const
{
    // 优先判断com.deepin.dde.dock.module.systemtray的control值是否为true（优先级更高），如果不为true，再判断每一个托盘对应的gsetting配置的control值
    bool isEnable = Utils::SettingValue("com.deepin.dde.dock.module.systemtray", QByteArray(), "control", false).toBool();
    return (isEnable || (m_gsettings && m_gsettings->get("control").toBool()));
}

QPixmap SystemPluginItem::icon()
{
    return QPixmap();
}
