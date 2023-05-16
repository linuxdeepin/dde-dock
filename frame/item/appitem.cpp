// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "appitem.h"
#include "themeappicon.h"
#include "xcb_misc.h"
#include "appswingeffectbuilder.h"
#include "utils.h"

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
#include <DConfig>

DGUI_USE_NAMESPACE
DCORE_USE_NAMESPACE

#define APP_DRAG_THRESHOLD      20

QPoint AppItem::MousePressPos;

AppItem::AppItem(const QGSettings *appSettings, const QGSettings *activeAppSettings, const QGSettings *dockedAppSettings, const QDBusObjectPath &entry, QWidget *parent)
    : DockItem(parent)
    , m_appSettings(appSettings)
    , m_activeAppSettings(activeAppSettings)
    , m_dockedAppSettings(dockedAppSettings)
    , m_appPreviewTips(nullptr)
    , m_itemEntryInter(new DockEntryInter("com.deepin.dde.daemon.Dock", entry.path(), QDBusConnection::sessionBus(), this))
    , m_swingEffectView(nullptr)
    , m_itemAnimation(nullptr)
    , m_wmHelper(DWindowManagerHelper::instance())
    , m_drag(nullptr)
    , m_dragging(false)
    , m_retryTimes(0)
    , m_iconValid(true)
    , m_lastclickTimes(0)
    , m_appIcon(QPixmap())
    , m_updateIconGeometryTimer(new QTimer(this))
    , m_retryObtainIconTimer(new QTimer(this))
    , m_refershIconTimer(new QTimer(this))
    , m_themeType(DGuiApplicationHelper::instance()->themeType())
{
    QHBoxLayout *centralLayout = new QHBoxLayout;
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setObjectName(m_itemEntryInter->name());
    setAcceptDrops(true);
    setLayout(centralLayout);

    m_id = m_itemEntryInter->id();
    m_active = m_itemEntryInter->isActive();
    m_currentWindowId = m_itemEntryInter->currentWindow();

    m_updateIconGeometryTimer->setInterval(500);
    m_updateIconGeometryTimer->setSingleShot(true);

    m_retryObtainIconTimer->setInterval(3000);
    m_retryObtainIconTimer->setSingleShot(true);

    m_refershIconTimer->setInterval(1000);
    m_refershIconTimer->setSingleShot(false);

    connect(m_itemEntryInter, &DockEntryInter::IsActiveChanged, this, &AppItem::activeChanged);
    connect(m_itemEntryInter, &DockEntryInter::IsActiveChanged, this, static_cast<void (AppItem::*)()>(&AppItem::update));
    connect(m_itemEntryInter, &DockEntryInter::WindowInfosChanged, this, &AppItem::updateWindowInfos, Qt::QueuedConnection);
    connect(m_itemEntryInter, &DockEntryInter::IconChanged, this, &AppItem::refreshIcon);

    connect(m_retryObtainIconTimer, &QTimer::timeout, this, &AppItem::refreshIcon, Qt::QueuedConnection);
    connect(m_updateIconGeometryTimer, &QTimer::timeout, this, &AppItem::updateWindowIconGeometries, Qt::QueuedConnection);
    connect(this, &AppItem::requestUpdateEntryGeometries, this, &AppItem::updateWindowIconGeometries);

    updateWindowInfos(m_itemEntryInter->windowInfos());
    refreshIcon();

    if (m_appSettings)
        connect(m_appSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);
    if (m_dockedAppSettings)
        connect(m_dockedAppSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);
    if (m_activeAppSettings)
        connect(m_activeAppSettings, &QGSettings::changed, this, &AppItem::onGSettingsChanged);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &AppItem::onThemeTypeChanged);

    /** 日历 1S定时判断是否刷新icon的处理 */
    connect(m_refershIconTimer, &QTimer::timeout, this, &AppItem::onRefreshIcon);
}

/**将属于同一个应用的窗口合并到同一个应用图标
 * @brief AppItem::checkEntry
 */
void AppItem::checkEntry()
{
    m_itemEntryInter->Check();
}

const QString AppItem::appId() const
{
    return m_id;
}

bool AppItem::isValid() const
{
    return m_itemEntryInter->isValid() && !m_itemEntryInter->id().isEmpty();
}

// Update _NET_WM_ICON_GEOMETRY property for windows that every item
// that manages, so that WM can do proper animations for specific
// window behaviors like minimization.
void AppItem::updateWindowIconGeometries()
{
    const QRect r(mapToGlobal(QPoint(0, 0)),
                  mapToGlobal(QPoint(width(), height())));

    if (Utils::IS_WAYLAND_DISPLAY){
        Q_EMIT requestUpdateItemMinimizedGeometry(r);
        return;
    }

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
    m_itemEntryInter->RequestUndock();
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

QString AppItem::accessibleName()
{
    return m_itemEntryInter->name();
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
    if (m_draging)
        return;

    if (m_dragging || (m_swingEffectView != nullptr && DockDisplayMode != Fashion))
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

        if (m_active) {
            painter.fillPath(path, QColor(0, 0, 0, 255 * 0.8));
        } else if (!m_windowInfos.isEmpty()) {
            if (hasAttention())
                painter.fillPath(path, QColor(241, 138, 46, 255 * .8));
            else
                painter.fillPath(path, QColor(0, 0, 0, 255 * 0.3));
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
            m_activeVerticalIndicator = QPixmap(":/indicator/resources/indicator_active_ver.svg");
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
        m_itemEntryInter->NewInstance(QX11Info::getTimestamp());

        // play launch effect
        if (m_windowInfos.isEmpty())
            playSwingEffect();

    } else if (e->button() == Qt::LeftButton) {
        if (checkAndResetTapHoldGestureState() && e->source() == Qt::MouseEventSynthesizedByQt) {
            qDebug() << "tap and hold gesture detected, ignore the synthesized mouse release event";
            return;
        }

        qDebug() << "app item clicked, name:" << m_itemEntryInter->name()
                 << "id:" << m_itemEntryInter->id() << "my-id:" << m_id << "icon:" << m_itemEntryInter->icon();

        m_itemEntryInter->Activate(QX11Info::getTimestamp());

        // play launch effect
        if (m_windowInfos.isEmpty() && DGuiApplicationHelper::isSpecialEffectsEnvironment())
            playSwingEffect();
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

    // handle drag
    if (e->buttons() != Qt::LeftButton)
        return;

    const QPoint pos = e->pos();
    if (!rect().contains(pos))
        return;
}

void AppItem::wheelEvent(QWheelEvent *e)
{
    if (checkGSettingsControl()) {
        return;
    }

    QWidget::wheelEvent(e);

    if (qAbs(e->angleDelta().y()) > 20) {
        m_itemEntryInter->PresentWindows();
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
    m_itemEntryInter->HandleDragDrop(QX11Info::getTimestamp(), uriList);
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

    m_itemEntryInter->HandleMenuItem(QX11Info::getTimestamp(), itemId);
}

const QString AppItem::contextMenu() const
{
    return m_itemEntryInter->menu();
}

QWidget *AppItem::popupTips()
{
    if (checkGSettingsControl())
        return nullptr;

    if (m_dragging)
        return nullptr;

    static TipsWidget appNameTips(topLevelWidget());
    appNameTips.setAccessibleName("tip");
    appNameTips.setObjectName(m_itemEntryInter->name());

    if (!m_windowInfos.isEmpty()) {
        Q_ASSERT(m_windowInfos.contains(m_currentWindowId));
        appNameTips.setText(m_windowInfos[m_currentWindowId].title.simplified());
    } else {
        appNameTips.setText(m_itemEntryInter->name().simplified());
    }

    return &appNameTips;
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
    if (info.size() <= 0)
        return;
    m_windowInfos = info;
    m_currentWindowId = info.firstKey();
    if (m_appPreviewTips)
        m_appPreviewTips->setWindowInfos(m_windowInfos, m_itemEntryInter->GetAllowedCloseWindows().value());
    m_updateIconGeometryTimer->start();

    // process attention effect
    if (hasAttention()) {
        if (DockDisplayMode == DisplayMode::Fashion)
            playSwingEffect();
    } else {
        stopSwingEffect();
    }

    update();
}

void AppItem::refreshIcon()
{
    if (!isVisible())
        return;

    const QString icon = m_itemEntryInter->icon();
    const int iconSize = qMin(width(), height());

    if (DockDisplayMode == Efficient)
        m_iconValid = ThemeAppIcon::getIcon(m_appIcon, icon, iconSize * 0.7, !m_iconValid);
    else
        m_iconValid = ThemeAppIcon::getIcon(m_appIcon, icon, iconSize * 0.8, !m_iconValid);

    if (!m_refershIconTimer->isActive() && m_itemEntryInter->icon() == "dde-calendar") {
        m_refershIconTimer->start();
    }

    if (!m_iconValid) {
        if (m_retryTimes < 10) {
            m_retryTimes++;
            qDebug() << m_itemEntryInter->name() << "obtain app icon(" << icon << ")failed, retry times:" << m_retryTimes;
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
    } else if (m_retryTimes > 0) {
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
    m_appPreviewTips->updateDockSize(DockSize);
    m_appPreviewTips->setWindowInfos(m_windowInfos, m_itemEntryInter->GetAllowedCloseWindows().value());
    m_appPreviewTips->updateLayoutDirection(DockPosition);

    connect(m_appPreviewTips, &PreviewContainer::requestActivateWindow, this, &AppItem::requestActivateWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestPreviewWindow, this, &AppItem::requestPreviewWindow, Qt::QueuedConnection);
    connect(m_appPreviewTips, &PreviewContainer::requestCancelPreviewWindow, this, &AppItem::requestCancelPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestHidePopup, this, &AppItem::hidePopup);
    connect(m_appPreviewTips, &PreviewContainer::requestCheckWindows, m_itemEntryInter, &DockEntryInter::Check);

    connect(m_appPreviewTips, &PreviewContainer::requestActivateWindow, this, &AppItem::onResetPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestCancelPreviewWindow, this, &AppItem::onResetPreview);
    connect(m_appPreviewTips, &PreviewContainer::requestHidePopup, this, &AppItem::onResetPreview);

    // 预览标题显示方式的配置
    DConfig *config = DConfig::create("org.deepin.dde.dock", "org.deepin.dde.dock");
    if (config->isValid() && config->keyList().contains("showWindowName"))
        m_appPreviewTips->setTitleDisplayMode(config->value("showWindowName").toInt());
    delete config;

    // 设置预览界面是否开启左右两边的圆角
    if (!PopupWindow.isNull() && m_wmHelper->hasComposite()) {
        PopupWindow->setLeftRightRadius(true);
    } else {
        PopupWindow->setLeftRightRadius(false);
    }

    showPopupWindow(m_appPreviewTips, true, 18);
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

    const QGSettings *setting = m_itemEntryInter->isDocked()
            ? m_dockedAppSettings
            : m_activeAppSettings;

    if (setting && setting->keys().contains("enable")) {
        const bool isEnable = !m_appSettings || (m_appSettings->keys().contains("enable") && m_appSettings->get("enable").toBool());
        setVisible(isEnable && setting->get("enable").toBool());
    }
}

bool AppItem::checkGSettingsControl() const
{
    const QGSettings *setting = m_itemEntryInter->isDocked()
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
