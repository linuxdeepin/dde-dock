#include "expandiconwidget.h"
#include "tray_gridview.h"
#include "tray_model.h"
#include "tray_delegate.h"
#include "dockpopupwindow.h"

#include <DGuiApplicationHelper>
#include <DRegionMonitor>
#include <QPainter>

#include <xcb/xproto.h>

DGUI_USE_NAMESPACE

ExpandIconWidget::ExpandIconWidget(QWidget *parent, Qt::WindowFlags f)
    : BaseTrayWidget(parent, f)
    , m_regionInter(new DRegionMonitor(this))
    , m_position(Dock::Position::Bottom)
    , m_trayView(nullptr)
{
    connect(m_regionInter, &DRegionMonitor::buttonPress, this, &ExpandIconWidget::onGlobMousePress);
}

ExpandIconWidget::~ExpandIconWidget()
{
}

void ExpandIconWidget::setPositonValue(Dock::Position position)
{
    m_position = position;
}

void ExpandIconWidget::sendClick(uint8_t mouseButton, int x, int y)
{
    Q_UNUSED(x);
    Q_UNUSED(y);

    if (mouseButton != XCB_BUTTON_INDEX_1)
        return;

    TrayGridView *trayIcon = popupTrayView();
    setTrayPanelVisible(!trayIcon->isVisible());
}

void ExpandIconWidget::setTrayPanelVisible(bool visible)
{
    TrayGridView *trayIcon = popupTrayView();
    if (visible) {
        resetPosition();
        trayIcon->show();
        m_regionInter->registerRegion();
    } else {
        trayIcon->hide();
        m_regionInter->unregisterRegion();
    }
}

QPixmap ExpandIconWidget::icon()
{
    return QPixmap(dropIconFile());
}

void ExpandIconWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);
    QPixmap pixmap(dropIconFile());
    painter.drawPixmap(0, 0, pixmap);
}

const QString ExpandIconWidget::dropIconFile() const
{
    QString arrow;
    switch (m_position) {
    case Dock::Position::Bottom: {
        arrow = "up";
        break;
    }
    case Dock::Position::Top: {
        arrow = "down";
        break;
    }
    case Dock::Position::Left: {
        arrow = "right";
        break;
    }
    case Dock::Position::Right: {
        arrow = "left";
        break;
    }
    }

    QString iconFile = QString(":/icons/resources/arrow-%1").arg(arrow);
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconFile += QString("-dark");

    return iconFile + ".svg";
}

TrayGridView *ExpandIconWidget::popupTrayView()
{
    if (m_trayView)
        return m_trayView;

    m_trayView = new TrayGridView(nullptr);
    TrayModel *model = new TrayModel(m_trayView, true, false);
    TrayDelegate *trayDelegate = new TrayDelegate(m_trayView);
    m_trayView->setWindowFlags(Qt::FramelessWindowHint | Qt::Popup);
    m_trayView->setModel(model);
    m_trayView->setItemDelegate(trayDelegate);
    m_trayView->setSpacing(3);
    m_trayView->setDragDistance(5);

    connect(m_trayView, &TrayGridView::rowCountChanged, this, [ this ] {
        int count = m_trayView->model()->rowCount();
        if (count > 0) {
            int lineCount = (count % 3) != 0 ? (count / 3 + 1) : (count / 3);
            // 如果只有一行，则根据实际的数量显示宽度
            int columnCount = qMin(count, 3);
            int width = ITEM_SIZE * columnCount + m_trayView->spacing() * 2;;
            int height = lineCount * ITEM_SIZE + m_trayView->spacing() * (lineCount - 1) + ITEM_SPACING;
            m_trayView->setFixedSize(width, height);
            resetPosition();
        } else if (m_trayView->isVisible()) {
            m_trayView->hide();
        }
        Q_EMIT trayVisbleChanged(count > 0);
    });

    connect(trayDelegate, &TrayDelegate::removeRow, this, [ = ](const QModelIndex &index) {
        QAbstractItemModel *abModel = model;
        abModel->removeRow(index.row(),index.parent());
    });
    connect(m_trayView, &TrayGridView::requestRemove, model, &TrayModel::removeRow);
    return m_trayView;
}

void ExpandIconWidget::resetPosition()
{
    if (!parentWidget())
        return;

    TrayGridView *trayView = popupTrayView();
    QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
    ptPos.setY(ptPos.y() - trayView->height());
    ptPos.setX(ptPos.x() - trayView->width());
    trayView->move(ptPos);
}

void ExpandIconWidget::onGlobMousePress(const QPoint &mousePos, const int flag)
{
    if (!isVisible() || !((flag == DRegionMonitor::WatchedFlags::Button_Left) || (flag == DRegionMonitor::WatchedFlags::Button_Right)))
        return;

    TrayGridView *trayView = popupTrayView();
    QPoint ptPos = parentWidget()->mapToGlobal(this->pos());
    const QRect rect = QRect(ptPos, size());
    if (rect.contains(mousePos))
        return;

    const QRect rctView(trayView->pos(), trayView->size());
    if (rctView.contains(mousePos))
        return;

    trayView->hide();
}
