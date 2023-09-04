// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "appitem.h"
#include "docksettings.h"
#include "taskmanager/windowinfobase.h"
#include "themeappicon.h"
#include "xcb_misc.h"
#include "appswingeffectbuilder.h"
#include "utils.h"
#include "screenspliter.h"

#include <X11/X.h>
#include <X11/Xlib.h>

#include <QPainter>
#include <QDrag>
#include <QMouseEvent>
#include <QApplication>
#include <QHBoxLayout>
#include <QGraphicsScene>
#include <QTimeLine>
#include <QX11Info>
#include <QGSettings>

#include <DGuiApplicationHelper>
#include <DPlatformTheme>
#include <DConfig>

#include <cstdint>
#include <sys/types.h>

DGUI_USE_NAMESPACE
DCORE_USE_NAMESPACE

#define APP_DRAG_THRESHOLD      20

QPoint AppItem::MousePressPos;

AppItem::AppItem(const QGSettings *appSettings, const QGSettings *activeAppSettings, const QGSettings *dockedAppSettings, const Entry *entry, QWidget *parent)
    : DockItem(parent)
    , m_appSettings(appSettings)
    , m_activeAppSettings(activeAppSettings)
    , m_dockedAppSettings(dockedAppSettings)
    , m_itemEntry(const_cast<Entry*>(entry))
    , m_appPreviewTips(nullptr)
    , m_swingEffectView(nullptr)
    , m_itemAnimation(nullptr)
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_drag(nullptr)
    , m_retryTimes(0)
    , m_iconValid(true)
    , m_lastclickTimes(0)
    , m_showMultiWindow(DockSettings::instance()->showMultiWindow())
    , m_appIcon(QPixmap())
    , m_activeColor(DGuiApplicationHelper::instance()->systemTheme()->activeColor())
    , m_updateIconGeometryTimer(new QTimer(this))
    , m_retryObtainIconTimer(new QTimer(this))
    , m_refershIconTimer(new QTimer(this))
    , m_themeType(DGuiApplicationHelper::instance()->themeType())
    , m_createMSecs(QDateTime::currentMSecsSinceEpoch())
    , m_screenSpliter(ScreenSpliterFactory::createScreenSpliter(this))
{
    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setAcceptDrops(true);
    setLayout(centralLayout);

    m_id = m_itemEntry->getId();
    m_active = m_itemEntry->getIsActive();
    m_name = m_itemEntry->getName();
    m_icon = m_itemEntry->getIcon();
    m_mode = m_itemEntry->mode();
    m_isDocked = m_itemEntry->getIsDocked();
    m_menu = m_itemEntry->getMenu();

    setObjectName(m_name);

    m_updateIconGeometryTimer->setInterval(500);
    m_updateIconGeometryTimer->setSingleShot(true);

    m_retryObtainIconTimer->setInterval(3000);
    m_retryObtainIconTimer->setSingleShot(true);

    m_refershIconTimer->setInterval(1000);
    m_refershIconTimer->setSingleShot(false);

    connect(m_itemEntry, &Entry::isActiveChanged, this, &AppItem::activeChanged);
    connect(m_itemEntry, &Entry::isActiveChanged, this, static_cast<void (AppItem::*)()>(&AppItem::update));
    connect(m_itemEntry, &Entry::windowInfosChanged, this, &AppItem::updateWindowInfos, Qt::QueuedConnection);
    connect(m_itemEntry, &Entry::iconChanged, this, [=](QString icon) { 
        if (!icon.isEmpty() && icon != m_icon)
            m_icon = icon;
    });
    connect(m_itemEntry, &Entry::iconChanged, this, &AppItem::refreshIcon);
    connect(m_itemEntry, &Entry::modeChanged, this, [=] (int32_t mode) { m_mode = mode; Q_EMIT modeChanged(m_mode);});
    connect(m_updateIconGeometryTimer, &QTimer::timeout, this, &AppItem::updateWindowIconGeometries, Qt::QueuedConnection);
    connect(m_retryObtainIconTimer, &QTimer::timeout, this, &AppItem::refreshIcon, Qt::QueuedConnection);
    connect(DockSettings::instance(), &DockSettings::showMultiWindowChanged, this, [=] (bool show) {
        m_showMultiWindow = show;
    });
    connect(m_itemEntry, &Entry::nameChanged, this, [=](const QString& name){ m_name = name; });
    connect(m_itemEntry, &Entry::desktopFileChanged, this, [=](const QString& desktopfile){ m_desktopfile = desktopfile; });
    connect(m_itemEntry, &Entry::isDockedChanged, this, [=](bool docked){ m_isDocked = docked; });
    connect(m_itemEntry, &Entry::menuChanged, this, [=](const QString& menu){ m_menu = menu; });
    connect(m_itemEntry, &Entry::currentWindowChanged, this, [=](uint32_t currentWindow){ 
        m_currentWindow = currentWindow;
        Q_EMIT onCurrentWindowChanged(m_currentWindow);
    });

    connect(this, &AppItem::requestUpdateEntryGeometries, this, &AppItem::updateWindowIconGeometries);

    updateWindowInfos(m_itemEntry->getExportWindowInfos());

    if (m_appSettings)
        connect(m_appSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);
    if (m_dockedAppSettings)
        connect(m_dockedAppSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);
    if (m_activeAppSettings)
        connect(m_activeAppSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &AppItem::onThemeTypeChanged);
    connect(DGuiApplicationHelper::instance()->systemTheme(), &DPlatformTheme::activeColorChanged, this,
            [this](const auto &color) { m_activeColor = color; });

    /** 日历 1S定时判断是否刷新icon的处理 */
    connect(m_refershIconTimer, &QTimer::timeout, this, &AppItem::onRefreshIcon);
}

/**将属于同一个应用的窗口合并到同一个应用图标
 * @brief AppItem::checkEntry
 */
void AppItem::checkEntry()
{
    m_itemEntry->check();
}

const QString AppItem::appId() const
{
    return m_id;
}

QString AppItem::name() const
{
    return m_name;
}

bool AppItem::isValid() const
{
    return m_itemEntry->isValid() && !m_id.isEmpty();
}

// Update _NET_WM_ICON_GEOMETRY property for windows that every item
// that manages, so that WM can do proper animations for specific
// window behaviors like minimization.
void AppItem::updateWindowIconGeometries()
{
    // wayland没做处理
    if (Utils::IS_WAYLAND_DISPLAY)
        return;

    const QRect r(mapToGlobal(QPoint(0, 0)),
                  mapToGlobal(QPoint(width(), height())));
    if (!QX11Info::connection()) {
        qWarning() << "QX11Info::connection() is 0x0";
        return;
    }

    auto *xcb_misc = XcbMisc::instance();

    for (auto it(m_windowInfos.cbegin()); it != m_windowInfos.cend(); ++it)
        xcb_misc->set_window_icon_geometry(it.key(), r);
}

/**取消驻留在dock上的应用
 * @brief AppItem::undock
 */
void AppItem::undock()
{
    m_itemEntry->requestUndock();
}

QWidget *AppItem::appDragWidget()
{
    if (m_drag) {
        return m_drag->appDragWidget();
    }

    return nullptr;
}

void AppItem::setDockInfo(Dock::Position dockPosition, const QRect &dockGeometry)
{
    if (m_drag) {
        m_drag->appDragWidget()->setDockInfo(dockPosition, dockGeometry);
    }
}

void AppItem::setDraging(bool drag)
{
    if (drag == isDragging())
        return;

    DockItem::setDraging(drag);
    if (!drag)
        m_screenSpliter->releaseSplit();
}

void AppItem::startSplit(const QRect &rect)
{
    m_screenSpliter->startSplit(rect);
}

bool AppItem::supportSplitWindow()
{
    return m_screenSpliter->suportSplitScreen();
}

bool AppItem::splitWindowOnScreen(ScreenSpliter::SplitDirection direction)
{
    return m_screenSpliter->split(direction);
}

int AppItem::mode() const
{
    return m_mode;
}


Entry *AppItem::itemEntry() const
{
    return m_itemEntry;
}

QString AppItem::accessibleName()
{
    return m_name;
}

void AppItem::requestDock()
{
    m_itemEntry->requestDock();
}

bool AppItem::isDocked() const
{
    return m_isDocked;
}

qint64 AppItem::appOpenMSecs() const
{
    return m_createMSecs;
}

void AppItem::updateMSecs()
{
    m_createMSecs = QDateTime::currentMSecsSinceEpoch();
}

const WindowInfoMap &AppItem::windowsInfos() const
{
    return m_windowInfos;
}

void AppItem::moveEvent(QMoveEvent *e)
{
    DockItem::moveEvent(e);

    if (m_drag) {
        m_drag->appDragWidget()->setOriginPos(mapToGlobal(appIconPosition()));
    }

    m_updateIconGeometryTimer->start();
}

void AppItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (isDragging() || (m_swingEffectView != nullptr && DockDisplayMode != Fashion))
        return;

    QPainter painter(this);
    if (!painter.isActive())
        return;
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setRenderHint(QPainter::SmoothPixmapTransform, true);

    const QRectF itemRect = rect();

    if (DockDisplayMode == Efficient) {
        // draw background
        qreal min = qMin(itemRect.width(), itemRect.height());
        QRectF backgroundRect = QRectF(itemRect.x(), itemRect.y(), min, min);
        backgroundRect = backgroundRect.marginsRemoved(QMargins(2, 2, 2, 2));
        backgroundRect.moveCenter(itemRect.center());

        QPainterPath path;
        path.addRoundedRect(backgroundRect, 8, 8);

        // 在没有开启窗口多开的情况下，显示背景色
        if (!m_showMultiWindow) {
            if (m_active) {
                QColor color = Qt::black;
                color.setAlpha(255 * 0.8);
                painter.fillPath(path, color);
            } else if (!m_windowInfos.isEmpty()) {
                if (hasAttention()) {
                    painter.fillPath(path, QColor(241, 138, 46, 255 * .8));
                } else {
                    QColor color = Qt::black;
                    color.setAlpha(255 * 0.3);
                    painter.fillPath(path, color);
                }
            }
        }
    } else {
        if (!m_windowInfos.isEmpty()) {
            QPoint p;
            QPixmap pixmap;
            QPixmap activePixmap;
            if (DGuiApplicationHelper::DarkType == m_themeType) {
                m_horizontalIndicator = QPixmap(":/indicator/resources/indicator_dark.svg");
                m_verticalIndicator = QPixmap(":/indicator/resources/indicator_dark_ver.svg");
            } else {
                m_horizontalIndicator = QPixmap(":/indicator/resources/indicator.svg");
                m_verticalIndicator = QPixmap(":/indicator/resources/indicator_ver.svg");
            }
            m_activeHorizontalIndicator = QPixmap(":/indicator/resources/indicator_active.svg");
            m_activeHorizontalIndicator.fill(m_activeColor);
            m_activeVerticalIndicator = QPixmap(":/indicator/resources/indicator_active_ver.svg");
            m_activeVerticalIndicator.fill(m_activeColor);
            switch (DockPosition) {
            case Top:
                pixmap = m_horizontalIndicator;
                activePixmap = m_activeHorizontalIndicator;
                p.setX((itemRect.width() - pixmap.width()) / 2);
                p.setY(1);
                break;
            case Bottom:
                pixmap = m_horizontalIndicator;
                activePixmap = m_activeHorizontalIndicator;
                p.setX((itemRect.width() - pixmap.width()) / 2);
                p.setY(itemRect.height() - pixmap.height() - 1);
                break;
            case Left:
                pixmap = m_verticalIndicator;
                activePixmap = m_activeVerticalIndicator;
                p.setX(1);
                p.setY((itemRect.height() - pixmap.height()) / 2);
                break;
            case Right:
                pixmap = m_verticalIndicator;
                activePixmap = m_activeVerticalIndicator;
                p.setX(itemRect.width() - pixmap.width() - 1);
                p.setY((itemRect.height() - pixmap.height()) / 2);
                break;
            }

            if (m_active)
                painter.drawPixmap(p, activePixmap);
            else
                painter.drawPixmap(p, pixmap);
        }
    }

    if (m_swingEffectView != nullptr)
        return;

    // icon
    if (m_appIcon.isNull())
        return;

    painter.drawPixmap(appIconPosition(), m_appIcon);
}

void AppItem::mouseReleaseEvent(QMouseEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }

    // 获取时间戳qint64转quint64，是不存在任何问题的
    quint64 curTimestamp = QDateTime::currentDateTime().toMSecsSinceEpoch();
    if ((curTimestamp - m_lastclickTimes) < 300)
        return;

    m_lastclickTimes = curTimestamp;

    // 鼠标在图标外边松开时，没必要响应点击操作
    const QRect rect { QPoint(0, 0), size()};
    if (!rect.contains(mapFromGlobal(QCursor::pos())))
        return;

    if (e->button() == Qt::MiddleButton) {
        m_itemEntry->newInstance(QX11Info::getTimestamp());

        // play launch effect
        if (m_windowInfos.isEmpty())
            playSwingEffect();

    } else if (e->button() == Qt::LeftButton) {
        if (checkAndResetTapHoldGestureState() && e->source() == Qt::MouseEventSynthesizedByQt) {
            qDebug() << "tap and hold gesture detected, ignore the synthesized mouse release event";
            return;
        }

        if (m_showMultiWindow) {
            // 如果开启了多窗口显示，则直接新建一个窗口
            m_itemEntry->newInstance(QX11Info::getTimestamp());
        } else {
            // 如果没有开启新窗口显示，则
            m_itemEntry->active(QX11Info::getTimestamp());
            // play launch effect
            if (m_windowInfos.isEmpty() && DGuiApplicationHelper::isSpecialEffectsEnvironment())
                playSwingEffect();
        }
    }
}

void AppItem::mousePressEvent(QMouseEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }
    m_updateIconGeometryTimer->stop();
    hidePopup();

    if (e->button() == Qt::LeftButton)
        MousePressPos = e->pos();

    // context menu will handle in DockItem
    DockItem::mousePressEvent(e);
}

void AppItem::mouseMoveEvent(QMouseEvent *e)
{
    e->accept();

    // handle preview
    //    if (e->buttons() == Qt::NoButton)
    //        return showPreview();

    // handle drag
    if (e->buttons() != Qt::LeftButton)
        return;

    const QPoint pos = e->pos();
    if (!rect().contains(pos))
        return;

    const QPoint distance = pos - MousePressPos;
    if (distance.manhattanLength() > APP_DRAG_THRESHOLD)
        return startDrag();
}

void AppItem::wheelEvent(QWheelEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }

    QWidget::wheelEvent(e);

    if (qAbs(e->angleDelta().y()) > 20) {
        m_itemEntry->presentWindows();
    }
}

void AppItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    refreshIcon();
}

void AppItem::dragEnterEvent(QDragEnterEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }

    // ignore drag from panel
    if (e->source()) {
        return e->ignore();
    }

    // ignore request dock event
    QString draggingMimeKey = e->mimeData()->formats().contains("RequestDock") ? "RequestDock" : "text/plain";
    if (QMimeDatabase().mimeTypeForFile(e->mimeData()->data(draggingMimeKey)).name() == "application/x-desktop") {
        return e->ignore();
    }

    e->accept();
}

void AppItem::dragMoveEvent(QDragMoveEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }

    DockItem::dragMoveEvent(e);

    if (m_windowInfos.isEmpty())
        return;

    if (!PopupWindow->isVisible() || !m_appPreviewTips)
        showPreview();
}

void AppItem::dropEvent(QDropEvent *e)
{
    QStringList uriList;
    for (auto uri : e->mimeData()->urls()) {
        uriList << uri.toEncoded();
    }

    qDebug() << "accept drop event with URIs: " << uriList;
    m_itemEntry->handleDragDrop(QX11Info::getTimestamp(), uriList);
}

void AppItem::leaveEvent(QEvent *e)
{
    DockItem::leaveEvent(e);

    if (m_appPreviewTips) {
        if (m_appPreviewTips->isVisible()) {
            m_appPreviewTips->prepareHide();
        }
    }
}

void AppItem::showHoverTips()
{
    if (checkGSettingsControl()) {
        return;
    }

    if (!m_windowInfos.isEmpty())
        return showPreview();

    DockItem::showHoverTips();
}

void AppItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    Q_UNUSED(checked);

    m_itemEntry->handleMenuItem(QX11Info::getTimestamp(), itemId);
}

const QString AppItem::contextMenu() const
{
    return m_menu;
}

QWidget *AppItem::popupTips()
{
    if (checkGSettingsControl())
        return nullptr;

    if (isDragging())
        return nullptr;

    static TipsWidget appNameTips(topLevelWidget());
    appNameTips.setAccessibleName("tip");
    appNameTips.setObjectName(m_name);

    if (!m_windowInfos.isEmpty()) {
        const quint32 currentWindow = m_currentWindow;
        Q_ASSERT(m_windowInfos.contains(currentWindow));
        appNameTips.setText(m_windowInfos[currentWindow].title.simplified());
    } else {
        appNameTips.setText(m_name.simplified());
    }

    return &appNameTips;
}

void AppItem::startDrag()
{
    // 拖拽实现放到mainpanelcontrol

    /*
    if (!acceptDrops())
        return;

    if (checkGSettingsControl()) {
        return;
    }

    m_dragging = true;
    update();

    const QPixmap &dragPix = m_appIcon;

    m_drag = new AppDrag(this);
    m_drag->setMimeData(new QMimeData);

    // handle drag finished here
    connect(m_drag->appDragWidget(), &AppDragWidget::destroyed, this, [ = ] {
        m_dragging = false;
        m_drag.clear();
        setVisible(true);
        update();
    });

    if (m_wmHelper->hasComposite()) {
        m_drag->setPixmap(dragPix);
        m_drag->appDragWidget()->setOriginPos(mapToGlobal(appIconPosition()));
        emit dragStarted();
        m_drag->exec(Qt::MoveAction);
    } else {
        m_drag->QDrag::setPixmap(dragPix);
        m_drag->setHotSpot(dragPix.rect().center() / dragPix.devicePixelRatioF());
        emit dragStarted();
        m_drag->QDrag::exec(Qt::MoveAction);
    }

    // MainPanel will put this item to Item-Container when received this signal(MainPanel::itemDropped)
    //emit itemDropped(m_drag->target());

    if (!m_wmHelper->hasComposite()) {
        if (!m_drag->target()) {
            m_itemEntryInter->RequestUndock();
        }
    }
    */
}

bool AppItem::hasAttention() const
{
    auto it = std::find_if(m_windowInfos.constBegin(), m_windowInfos.constEnd(), [ = ] (const auto &info) {
        return info.attention;
    });

    return (it != m_windowInfos.end());
}

QPoint AppItem::appIconPosition() const
{
    const auto ratio = devicePixelRatioF();
    const QRectF itemRect = rect();
    const QRectF iconRect = m_appIcon.rect();
    const qreal iconX = itemRect.center().x() - iconRect.center().x() / ratio;
    const qreal iconY = itemRect.center().y() - iconRect.center().y() / ratio;

    return QPoint(iconX, iconY);
}

void AppItem::updateWindowInfos(const WindowInfoMap &info)
{
    // 如果是打开第一个窗口，则更新窗口时间
    if (m_windowInfos.isEmpty() && !info.isEmpty())
        updateMSecs();

    m_windowInfos = info;
    if (m_appPreviewTips)
        m_appPreviewTips->setWindowInfos(m_windowInfos, m_itemEntry->getAllowedClosedWindowIds().toList());
    m_updateIconGeometryTimer->start();

    // process attention effect
    if (hasAttention()) {
        if (DockDisplayMode == DisplayMode::Fashion)
            playSwingEffect();
    } else {
        stopSwingEffect();
    }

    update();

    // 通知外面窗体数量发生变化，需要更新多开窗口的信息
    Q_EMIT windowCountChanged();
}

void AppItem::refreshIcon()
{
    if (!isVisible())
        return;

    const int iconSize = qMin(width(), height());

    if (DockDisplayMode == Efficient)
        m_iconValid = ThemeAppIcon::getIcon(m_appIcon, m_icon, iconSize * 0.7, !m_iconValid);
    else
        m_iconValid = ThemeAppIcon::getIcon(m_appIcon, m_icon, iconSize * 0.8, !m_iconValid);

    if (!m_refershIconTimer->isActive() && m_icon == "dde-calendar") {
        m_refershIconTimer->start();
    }

    if (!m_iconValid) {
        if (m_retryTimes < 10) {
            m_retryTimes++;
            qDebug() << m_name << "obtain app icon(" << m_icon << ")failed, retry times:" << m_retryTimes;
            // Maybe the icon was installed after we loaded the caches.
            // QIcon::setThemeSearchPaths will force Qt to re-check the gtk cache validity.
            QIcon::setThemeSearchPaths(QIcon::themeSearchPaths());

            m_retryObtainIconTimer->start();
        } else {
            // 如果图标获取失败，一分钟后再自动刷新一次（如果还是显示异常，基本需要应用自身看下为什么了）
            if (!m_iconValid)
                QTimer::singleShot(60 * 1000, this, &AppItem::refreshIcon);
        }

        update();

        return;
    }

    if (m_retryTimes > 0) {
        // reset times
        m_retryTimes = 0;
    }

    update();

    m_updateIconGeometryTimer->start();
}

void AppItem::onRefreshIcon()
{
    if (QDate::currentDate() == m_curDate)
        return;

    m_curDate = QDate::currentDate();
    refreshIcon();
}

void AppItem::onResetPreview()
{
    if (m_appPreviewTips != nullptr) {
        m_appPreviewTips->deleteLater();
        m_appPreviewTips = nullptr;
    }
}

void AppItem::activeChanged()
{
    m_active = !m_active;
}

void AppItem::showPreview()
{
    if (m_windowInfos.isEmpty())
        return;

    m_appPreviewTips = new PreviewContainer;
    m_appPreviewTips->setWindowInfos(m_windowInfos, m_itemEntry->getAllowedClosedWindowIds().toList());
    m_appPreviewTips->updateLayoutDirection(DockPosition);

    connect(m_appPreviewTips, &PreviewContainer::requestActivateWindow, this, &AppItem::activeWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestPreviewWindow, this, &AppItem::requestPreviewWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestCancelPreviewWindow, this, &AppItem::requestCancelPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestHidePopup, this, &AppItem::hidePopup);
    connect(m_appPreviewTips, &PreviewContainer::requestCheckWindows, m_itemEntry, &Entry::check);

    connect(m_appPreviewTips, &PreviewContainer::requestActivateWindow, this, &AppItem::onResetPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestCancelPreviewWindow, this, &AppItem::onResetPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestHidePopup, this, &AppItem::onResetPreview);

    // 预览标题显示方式的配置
    DConfig config(QString("com.deepin.dde.dock.dconfig"), QString());
    if (config.isValid() && config.keyList().contains("Dock_Show_Window_name"))
        m_appPreviewTips->setTitleDisplayMode(config.value("Dock_Show_Window_name").toInt());

    showPopupWindow(m_appPreviewTips, true);
}

void AppItem::playSwingEffect()
{
    // NOTE(sbw): return if animation view already playing
    if (m_swingEffectView != nullptr)
        return;

    if (rect().isEmpty())
        return checkAttentionEffect();

    stopSwingEffect();

    QPair<QGraphicsView *, QGraphicsItemAnimation *> pair =  SwingEffect(
                this, m_appIcon, rect(), devicePixelRatioF());

    m_swingEffectView = pair.first;
    m_itemAnimation = pair.second;

    QTimeLine *tl = m_itemAnimation->timeLine();
    connect(tl, &QTimeLine::stateChanged, this, [ = ](QTimeLine::State newState) {
        if (newState == QTimeLine::NotRunning) {
            m_swingEffectView->hide();
            layout()->removeWidget(m_swingEffectView);
            m_swingEffectView = nullptr;
            m_itemAnimation = nullptr;
            checkAttentionEffect();
        }
    });

    layout()->addWidget(m_swingEffectView);
    tl->start();
}

void AppItem::stopSwingEffect()
{
    if (m_swingEffectView == nullptr || m_itemAnimation == nullptr)
        return;

    // stop swing effect
    m_swingEffectView->setVisible(false);

    if (m_itemAnimation->timeLine() && m_itemAnimation->timeLine()->state() != QTimeLine::NotRunning)
        m_itemAnimation->timeLine()->stop();
}

void AppItem::checkAttentionEffect()
{
    QTimer::singleShot(1000, this, [ = ] {
        if (DockDisplayMode == DisplayMode::Fashion && hasAttention())
            playSwingEffect();
    });
}

void AppItem::onGSettingsChanged(const QString &key)
{
    if (key != "enable") {
        return;
    }

    const QGSettings *setting = m_isDocked
            ? m_dockedAppSettings
            : m_activeAppSettings;

    if (setting && setting->keys().contains("enable")) {
        const bool isEnable = !m_appSettings || (m_appSettings->keys().contains("enable") && m_appSettings->get("enable").toBool());
        setVisible(isEnable && setting->get("enable").toBool());
    }
}

bool AppItem::checkGSettingsControl() const
{
    const QGSettings *setting = m_isDocked
            ? m_dockedAppSettings
            : m_activeAppSettings;

    return ((m_appSettings && m_appSettings->keys().contains("control") && m_appSettings->get("control").toBool())
            || (setting && setting->keys().contains("control") && setting->get("control").toBool()));
}

void AppItem::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    m_themeType = themeType;
    update();
}

// 放到最下面是因为析构函数和匿名函数会影响lcov统计单元测试的覆盖率
AppItem::~AppItem()
{
    stopSwingEffect();
}

void AppItem::showEvent(QShowEvent *e)
{
    DockItem::showEvent(e);

    QTimer::singleShot(0, this, [ = ] {
        onGSettingsChanged("enable");
    });

    refreshIcon();
}

void AppItem::activeWindow(WId wid)
{
    m_itemEntry->activeWindow(wid);
}
