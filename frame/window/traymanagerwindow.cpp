#include "traymanagerwindow.h"
#include "quickpluginwindow.h"
#include "tray_gridview.h"
#include "tray_delegate.h"
#include "tray_model.h"
#include "constants.h"
#include "quicksettingcontainer.h"
#include "systempluginwindow.h"
#include "datetimedisplayer.h"

#include <DGuiApplicationHelper>
#include <DRegionMonitor>

#include <QDropEvent>
#include <QBoxLayout>
#include <QLabel>
#include <QMimeData>
#include <QDBusConnection>
#include <QPainter>

#define MAXFIXEDSIZE 999999
#define CRITLCALHEIGHT 56

TrayManagerWindow::TrayManagerWindow(QWidget *parent)
    : QWidget(parent)
    , m_appPluginDatetimeWidget(new DBlurEffectWidget(this))
    , m_systemPluginWidget(new SystemPluginWindow(this))
    , m_appPluginWidget(new QWidget(m_appPluginDatetimeWidget))
    , m_quickIconWidget(new QuickPluginWindow(m_appPluginWidget))
    , m_dateTimeWidget(new DateTimeDisplayer(m_appPluginDatetimeWidget))
    , m_appPluginLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
    , m_appDatetimeLayout(new QBoxLayout(QBoxLayout::Direction::TopToBottom, this))
    , m_mainLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this))
    , m_trayView(new TrayGridView(this))
    , m_model(new TrayModel(m_trayView, false, true))
    , m_postion(Dock::Position::Bottom)
{
    initUi();
    initConnection();

    setAcceptDrops(true);
    setMouseTracking(true);
}

TrayManagerWindow::~TrayManagerWindow()
{
}

void TrayManagerWindow::setPositon(Dock::Position position)
{
    if (m_postion == position)
        return;

    m_postion = position;
    QModelIndex index = m_model->index(0, 0);
    m_trayView->closePersistentEditor(index);
    TrayDelegate *delegate = static_cast<TrayDelegate *>(m_trayView->itemDelegate());
    delegate->setPositon(position);
    m_trayView->openPersistentEditor(index);

    m_quickIconWidget->setPositon(position);
    m_dateTimeWidget->setPositon(position);
    m_systemPluginWidget->setPositon(position);

    QTimer::singleShot(0, this, [ this ]{
        if (showSingleRow())
            resetSingleDirection();
        else
            resetMultiDirection();

        resetChildWidgetSize();
    });
}

int TrayManagerWindow::appDatetimeSize()
{
    int count = m_trayView->model()->rowCount();
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom) {
        QMargins m = m_appDatetimeLayout->contentsMargins();
        int trayWidth = count * ITEM_SIZE + m_trayView->spacing() * (count - 1) + 5;
        int topWidth = trayWidth + m_quickIconWidget->suitableSize().width() + m.left() + m.right() + m_appPluginLayout->spacing();
        int spacing = m.left() + m.right() + m_appPluginLayout->spacing();
        if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::LeftToRight)
            return topWidth + m_appDatetimeLayout->spacing() + m_dateTimeWidget->suitableSize().width() + m_appDatetimeLayout->spacing() + 10;

        int bottomWidth = m_dateTimeWidget->suitableSize().width();
        return (topWidth > bottomWidth ? topWidth : bottomWidth) + m_appDatetimeLayout->spacing() + spacing + 10;
    }

    int trayHeight = count * ITEM_SIZE + m_trayView->spacing() * (count - 1) + 5;
    int datetimeHeight = m_dateTimeWidget->suitableSize().height();
    QMargins m = m_appDatetimeLayout->contentsMargins();
    int traypluginHeight = trayHeight + m_quickIconWidget->suitableSize().height() + m.top() + m.bottom() + m_appPluginLayout->spacing();
    if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::TopToBottom)
        return traypluginHeight + m_appDatetimeLayout->spacing() + m_dateTimeWidget->suitableSize().height() + 10;
    return (traypluginHeight > datetimeHeight ? traypluginHeight : datetimeHeight) + 10;
}

QSize TrayManagerWindow::suitableSize()
{
    QMargins m = m_mainLayout->contentsMargins();
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom) {
        return QSize(appDatetimeSize() + m_appDatetimeLayout->spacing() +
                     m_systemPluginWidget->suitableSize().width() + m_mainLayout->spacing() +
                     m.left() + m.right(), height());
    }

    return QSize(width(), appDatetimeSize() + m_appDatetimeLayout->spacing() +
                 m_systemPluginWidget->suitableSize().height() + m_mainLayout->spacing() +
                 m.top() + m.bottom());
}

void TrayManagerWindow::resizeEvent(QResizeEvent *event)
{
    Q_UNUSED(event);

    if (showSingleRow())
        resetSingleDirection();
    else
        resetMultiDirection();

    resetChildWidgetSize();
}

void TrayManagerWindow::initUi()
{
    TrayDelegate *delegate = new TrayDelegate(m_trayView);
    m_trayView->setModel(m_model);
    m_trayView->setItemDelegate(delegate);

    WinInfo info;
    info.type = TrayIconType::EXPANDICON;
    m_model->addRow(info);
    m_trayView->openPersistentEditor(m_model->index(0, 0));

    // 左侧的区域，包括应用托盘插件和下方的日期时间区域
    m_appPluginDatetimeWidget->setBlurRectXRadius(10);
    m_appPluginDatetimeWidget->setBlurRectYRadius(10);
    m_appPluginDatetimeWidget->setMaskAlpha(uint8_t(0.1 * 255));
    m_appPluginDatetimeWidget->installEventFilter(this);

    m_appPluginLayout->setSpacing(0);
    m_appPluginWidget->setLayout(m_appPluginLayout);
    m_appPluginLayout->addWidget(m_trayView);
    m_appPluginLayout->addWidget(m_quickIconWidget);

    m_appPluginDatetimeWidget->setLayout(m_appDatetimeLayout);
    m_appDatetimeLayout->setContentsMargins(0, 0, 0, 0);
    m_appDatetimeLayout->setSpacing(3);
    m_appDatetimeLayout->addWidget(m_appPluginWidget);
    m_appDatetimeLayout->addWidget(m_dateTimeWidget);

    m_systemPluginWidget->setBlurRectXRadius(10);
    m_systemPluginWidget->setBlurRectYRadius(10);
    m_systemPluginWidget->installEventFilter(this);
    m_systemPluginWidget->setMaskAlpha(uint8_t(0.1 * 255));

    setLayout(m_mainLayout);
    m_mainLayout->setContentsMargins(8, 8, 8, 8);
    m_mainLayout->setSpacing(10);
    m_mainLayout->addWidget(m_appPluginDatetimeWidget);
    m_mainLayout->addWidget(m_systemPluginWidget);
}

void TrayManagerWindow::initConnection()
{
    connect(m_trayView, &TrayGridView::requestRemove, m_model, &TrayModel::removeRow);
    connect(m_trayView, &TrayGridView::rowCountChanged, this, &TrayManagerWindow::sizeChanged);
    connect(m_quickIconWidget, &QuickPluginWindow::itemCountChanged, this, [ this ] {
        m_quickIconWidget->setFixedSize(QWIDGETSIZE_MAX, QWIDGETSIZE_MAX);
        if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
            m_quickIconWidget->setFixedWidth(m_quickIconWidget->suitableSize().width());
        else
            m_quickIconWidget->setFixedHeight(m_quickIconWidget->suitableSize().height());

        Q_EMIT sizeChanged();
    });

    connect(m_systemPluginWidget, &SystemPluginWindow::pluginSizeChanged, this, [ this ] {
        m_systemPluginWidget->setFixedSize(QWIDGETSIZE_MAX, QWIDGETSIZE_MAX);
        if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
            m_systemPluginWidget->setFixedWidth(m_systemPluginWidget->suitableSize().width());
        else
            m_systemPluginWidget->setFixedHeight(m_systemPluginWidget->suitableSize().height());

        Q_EMIT sizeChanged();
    });

    TrayDelegate *trayDelegate = static_cast<TrayDelegate *>(m_trayView->itemDelegate());
    connect(trayDelegate, &TrayDelegate::visibleChanged, this, [ this ](const QModelIndex &index, bool visible) {
        m_trayView->setRowHidden(index.row(), !visible);
        resetChildWidgetSize();
        Q_EMIT sizeChanged();
    });

    connect(m_trayView, &TrayGridView::dragLeaved, trayDelegate, [ trayDelegate ]{
        Q_EMIT trayDelegate->requestDrag(true);
    });
    connect(m_trayView, &TrayGridView::dragEntered, trayDelegate, [ trayDelegate ]{
        Q_EMIT trayDelegate->requestDrag(false);
    });
    connect(m_model, &TrayModel::requestUpdateWidget, this, [ this ](const QList<int> &idxs) {
        for (int i = 0; i < idxs.size(); i++) {
             int idx = idxs[i];
             if (idx < m_model->rowCount()) {
                 QModelIndex index = m_model->index(idx);
                 m_trayView->closePersistentEditor(index);
                 m_trayView->openPersistentEditor(index);
             }
        }
    });

    m_trayView->installEventFilter(this);
    m_quickIconWidget->installEventFilter(this);
    installEventFilter(this);
    QMetaObject::invokeMethod(this, &TrayManagerWindow::resetChildWidgetSize, Qt::QueuedConnection);
}

bool TrayManagerWindow::showSingleRow()
{
    if (m_postion == Dock::Position::Top || m_postion == Dock::Position::Bottom)
        return height() < CRITLCALHEIGHT;

    return true;
}

void TrayManagerWindow::resetChildWidgetSize()
{
    int count = 0;
    for (int i = 0; i < m_model->rowCount(); i++) {
        if (!m_trayView->isRowHidden(i))
            count++;
    }

    switch (m_postion) {
    case Dock::Position::Top:
    case Dock::Position::Bottom: {
        int trayWidth = count * ITEM_SIZE + m_trayView->spacing() * (count - 1) + 5;
        QMargins m = m_appPluginLayout->contentsMargins();
        m_appPluginDatetimeWidget->setFixedHeight(QWIDGETSIZE_MAX);// 取消固定高度显示
        if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::LeftToRight) {
            // 单行显示
            int trayHeight = m_appPluginDatetimeWidget->height() - m.top() - m.bottom();
            m_trayView->setFixedSize(trayWidth, trayHeight);
            m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), trayHeight);
            m_dateTimeWidget->setFixedSize(m_dateTimeWidget->suitableSize().width(), trayHeight);
        } else {
            // 多行显示
            int trayHeight = m_appPluginDatetimeWidget->height() / 2 - m.top() - m.bottom();
            m_trayView->setFixedSize(trayWidth, trayHeight);
            m_quickIconWidget->setFixedSize(m_quickIconWidget->suitableSize().width(), trayHeight);
            m_dateTimeWidget->setFixedSize(m_dateTimeWidget->suitableSize().width(), m_appPluginDatetimeWidget->height() / 2);
        }
        m_appPluginDatetimeWidget->setFixedWidth(appDatetimeSize());
        break;
    }
    case Dock::Position::Left:
    case Dock::Position::Right: {
        int trayHeight = count * ITEM_SIZE + m_trayView->spacing() * (count - 1) + 5;
        int quickAreaHeight = m_quickIconWidget->suitableSize().height();
        QMargins m = m_appPluginLayout->contentsMargins();
        m_appPluginDatetimeWidget->setFixedWidth(QWIDGETSIZE_MAX);// 取消固定宽度显示
        if (m_appDatetimeLayout->direction() == QBoxLayout::Direction::TopToBottom) {
            // 宽度较小的情况下,显示一列
            int datetimeHeight = m_dateTimeWidget->suitableSize().height();
            int sizeWidth = m_appPluginDatetimeWidget->width() - m.left() - m.right();
            m_trayView->setFixedSize(sizeWidth, trayHeight);
            m_quickIconWidget->setFixedSize(sizeWidth, quickAreaHeight);
            m_dateTimeWidget->setFixedSize(sizeWidth, datetimeHeight);
        } else {
            // 显示两列
            int trayWidth = m_appPluginDatetimeWidget->width() / 2 - m.left() - m.right();
            m_trayView->setFixedSize(trayWidth, trayHeight);
            m_quickIconWidget->setFixedSize(trayWidth, quickAreaHeight);
            m_dateTimeWidget->setFixedSize(m_appPluginDatetimeWidget->width() / 2, m_dateTimeWidget->suitableSize().height());
        }
        m_appPluginDatetimeWidget->setFixedHeight(appDatetimeSize());
        break;
    }
    }
}

void TrayManagerWindow::resetSingleDirection()
{
    switch (m_postion) {
    case Dock::Position::Top: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        // 应用和时间在一行显示
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appPluginLayout->setContentsMargins(2, 2, 2, 4);
        break;
    }
    case Dock::Position::Bottom: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appPluginLayout->setContentsMargins(2, 4, 2, 2);
        break;
    }
    case Dock::Position::Left: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appPluginLayout->setContentsMargins(2, 2, 4, 2);
        break;
    }
    case Dock::Position::Right: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appPluginLayout->setContentsMargins(4, 2, 2, 2);
        break;
    }
    }
}

void TrayManagerWindow::resetMultiDirection()
{
    switch (m_postion) {
    case Dock::Position::Top: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::BottomToTop);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appPluginLayout->setContentsMargins(2, 2, 2, 4);
        break;
    }
    case Dock::Position::Bottom: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_appPluginLayout->setContentsMargins(2, 4, 2, 2);
        break;
    }
    case Dock::Position::Left: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appPluginLayout->setContentsMargins(2, 2, 4, 2);
        break;
    }
    case Dock::Position::Right: {
        m_appPluginLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appDatetimeLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_mainLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_appPluginLayout->setContentsMargins(4, 2, 2, 2);
        break;
    }
    }
}

void TrayManagerWindow::dragEnterEvent(QDragEnterEvent *e)
{
    e->setDropAction(Qt::CopyAction);
    e->accept();
}

void TrayManagerWindow::dragMoveEvent(QDragMoveEvent *e)
{
    e->setDropAction(Qt::CopyAction);
    e->accept();
}

void TrayManagerWindow::dropEvent(QDropEvent *e)
{
    CustomMimeData *mimeData = const_cast<CustomMimeData *>(qobject_cast<const CustomMimeData *>(e->mimeData()));
    if (!mimeData)
        return;

    if (e->source() == this)
        return;

    QuickSettingItem *pluginItem = static_cast<QuickSettingItem *>(mimeData->data());
    if (pluginItem)
        m_quickIconWidget->addPlugin(pluginItem);
}

void TrayManagerWindow::dragLeaveEvent(QDragLeaveEvent *event)
{
    event->accept();
}
