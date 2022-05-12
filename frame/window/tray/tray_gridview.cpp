#include "tray_gridview.h"

#include <QMouseEvent>
#include <QDragEnterEvent>
#include <QDragLeaveEvent>
#include <QDropEvent>
#include <QPropertyAnimation>
#include <QLabel>
#include <QDrag>
#include <QMimeData>
#include <QApplication>
#include <QDebug>
#include <QTimer>

#include "tray_model.h"
#include "basetraywidget.h"

TrayGridView::TrayGridView(QWidget *parent)
    : DListView(parent)
    , m_aniCurveType(QEasingCurve::Linear)
    , m_aniDuringTime(250)
    , m_dragDistance(15)
    , m_aniStartTime(new QTimer(this))
    , m_pressed(false)
    , m_aniRunning(false)
{
    initUi();
}

void TrayGridView::setDragDistance(int pixel)
{
    m_dragDistance = pixel;
}

void TrayGridView::setAnimationProperty(const QEasingCurve::Type easing, const int duringTime)
{
    m_aniCurveType = easing;
    m_aniDuringTime = duringTime;
}

void TrayGridView::moveAnimation()
{
    if (m_aniRunning || m_aniStartTime->isActive())
        return;

    const QModelIndex dropModelIndex = indexAt(m_dropPos);
    if (!dropModelIndex.isValid())
        return;

    const QModelIndex dragModelIndex = indexAt(m_dragPos);
    if (dragModelIndex == dropModelIndex)
        return;

    if (!dragModelIndex.isValid()) {
        m_dragPos = indexRect(dropModelIndex).center();
        return;
    }

    TrayModel *listModel = qobject_cast<TrayModel *>(model());
    if (!listModel)
        return;

    listModel->clearDragDropIndex();
    listModel->setDragingIndex(dragModelIndex);
    listModel->setDragDropIndex(dropModelIndex);

    const int startPos = dragModelIndex.row();
    const int endPos = dropModelIndex.row();

    const bool next = startPos <= endPos;
    const int start = next ? startPos : endPos;
    const int end = !next ? startPos : endPos;

    for (int i = start + next; i <= (end - !next); i++)
        createAnimation(i, next, (i == (end - !next)));

    m_dropPos = indexRect(dropModelIndex).center();
    m_dragPos = indexRect(dropModelIndex).center();
}

const QModelIndex TrayGridView::modelIndex(const int index) const
{
    return model()->index(index, 0, QModelIndex());
}

const QRect TrayGridView::indexRect(const QModelIndex &index) const
{
    return rectForIndex(index);
}

void TrayGridView::dropSwap()
{
    qDebug() << "drop end";
    TrayModel *listModel = qobject_cast<TrayModel *>(model());
    if (!listModel)
        return;

    QModelIndex index = indexAt(m_dropPos);
    if (!index.isValid())
        return;

    listModel->dropSwap(index.row());
    clearDragModelIndex();
    m_aniRunning = false;
    setState(NoState);
}

void TrayGridView::clearDragModelIndex()
{
    TrayModel *listModel = static_cast<TrayModel *>(this->model());
    if (!listModel)
        return;

    listModel->clearDragDropIndex();
}

void TrayGridView::createAnimation(const int pos, const bool moveNext, const bool isLastAni)
{
    qDebug() << "create moveAnimation";
    const QModelIndex index(modelIndex(pos));
    if (!index.isValid())
        return;

    QLabel *floatLabel = new QLabel(this);
    QPropertyAnimation *ani = new QPropertyAnimation(floatLabel, "pos", floatLabel);
    qreal ratio = qApp->devicePixelRatio();

    BaseTrayWidget *widget = qobject_cast<BaseTrayWidget *>(indexWidget(index));
    if (!widget)
        return;

    QPixmap pixmap = widget->icon();

    QString text = index.data(Qt::DisplayRole).toString();

    pixmap.scaled(pixmap.size() * ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    pixmap.setDevicePixelRatio(ratio);

    floatLabel->setFixedSize(indexRect(index).size());
    floatLabel->setPixmap(pixmap);
    floatLabel->show();

    ani->setStartValue(indexRect(index).center() - QPoint(0, floatLabel->height() /2));
    ani->setEndValue(indexRect(modelIndex(moveNext ? pos - 1 : pos + 1)).center() - QPoint(0, floatLabel->height() /2));
    ani->setEasingCurve(m_aniCurveType);
    ani->setDuration(m_aniDuringTime);

    connect(ani, &QPropertyAnimation::finished, floatLabel, &QLabel::deleteLater);
    if (isLastAni) {
        m_aniRunning = true;

        TrayModel *model = qobject_cast<TrayModel *>(this->model());
        if (!model)
            return;

        connect(ani, &QPropertyAnimation::finished, this, &TrayGridView::dropSwap);
        connect(ani, &QPropertyAnimation::valueChanged, m_aniStartTime, &QTimer::stop);
    } else {
    }

    ani->start(QPropertyAnimation::DeleteWhenStopped);
}

void TrayGridView::mousePressEvent(QMouseEvent *e)
{
    if (e->buttons() == Qt::LeftButton && !m_aniRunning)
        m_dragPos = e->pos();

    m_pressed = true;
}

void TrayGridView::mouseMoveEvent(QMouseEvent *e)
{
    setState(QAbstractItemView::NoState);
    e->accept();

    if (e->buttons() == Qt::RightButton)
        return DListView::mouseMoveEvent(e);

    QModelIndex index = indexAt(e->pos());
    if (!index.isValid())
        return DListView::mouseMoveEvent(e);

    if (m_pressed && (qAbs(e->pos().x() - m_dragPos.x()) > m_dragDistance ||
                      qAbs(e->pos().y() - m_dragPos.y()) > m_dragDistance)) {
        qDebug() << "start drag";
        beginDrag(Qt::CopyAction | Qt::MoveAction);
    }
}

void TrayGridView::mouseReleaseEvent(QMouseEvent *e)
{
    m_pressed = false;
}

void TrayGridView::dragEnterEvent(QDragEnterEvent *e)
{
    const QModelIndex index = indexAt(e->pos());

    if (model()->canDropMimeData(e->mimeData(), e->dropAction(), index.row(),
                                 index.column(), QModelIndex()))
        e->accept();
    else
        e->ignore();

    Q_EMIT dragEntered();
}

void TrayGridView::dragLeaveEvent(QDragLeaveEvent *e)
{
    m_aniStartTime->stop();
    e->accept();
    dragLeaved();
}

void TrayGridView::dragMoveEvent(QDragMoveEvent *e)
{
    m_aniStartTime->stop();
    if (m_aniRunning)
        return;

    setState(QAbstractItemView::DraggingState);

    QModelIndex index = indexAt(e->pos());
    if (index.isValid()) {
        if (m_dropPos != indexRect(index).center()) {
            qDebug() << "update drop position: " << index.row();
            m_dropPos = indexRect(index).center();
        }
    }

    if (m_pressed)
        m_aniStartTime->start();
}

void TrayGridView::dropEvent(QDropEvent *e)
{
    setState(DListView::NoState);
    clearDragModelIndex();

    if (m_aniStartTime->isActive())
        m_aniStartTime->stop();

    if (e->mimeData()->formats().contains("type") && e->source() != this) {
        e->setDropAction(Qt::CopyAction);
        e->accept();

        TrayModel *m = qobject_cast<TrayModel *>(model());
        if (m) {
            WinInfo info;
            info.type = static_cast<TrayIconType>(e->mimeData()->data("type").toInt());
            info.key = static_cast<QString>(e->mimeData()->data("key"));
            info.winId = static_cast<quint32>(e->mimeData()->data("winId").toInt());
            info.servicePath = static_cast<QString>(e->mimeData()->data("servicePath"));
            m->addRow(info);
        }
    } else {
        e->ignore();
    }

    DListView::dropEvent(e);
}

void TrayGridView::beginDrag(Qt::DropActions supportedActions)
{
    QModelIndex modelIndex = indexAt(m_dragPos);
    TrayIconType trayType = modelIndex.data(TrayModel::Role::TypeRole).value<TrayIconType>();
    // 展开图标不能移动
    if (trayType == TrayIconType::EXPANDICON)
        return;

    m_dropPos = indexRect(modelIndex).center();

    TrayModel *listModel = qobject_cast<TrayModel *>(model());
    if (!listModel)
        return;

    BaseTrayWidget *widget = qobject_cast<BaseTrayWidget *>(indexWidget(modelIndex));
    if (!widget)
        return;

    auto pixmap = widget->icon();

    QString text = modelIndex.data(Qt::DisplayRole).toString();
    QString itemKey = modelIndex.data(TrayModel::Role::KeyRole).toString();
    qreal ratio = qApp->devicePixelRatio();
    // 创建拖拽释放时的应用图标
    QLabel *pixLabel = new QLabel(this);
    pixLabel->setPixmap(pixmap);
    pixLabel->setFixedSize(indexRect(modelIndex).size() / ratio);

    QRect rectIcon(pixLabel->rect().topLeft(), pixLabel->size());

    listModel->setDragingIndex(modelIndex);

    QDrag *drag = new QDrag(this);
    pixmap.scaled(pixmap.size() * ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    pixmap.setDevicePixelRatio(ratio);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center() / ratio);
    QMimeData *data = model()->mimeData(QModelIndexList() << modelIndex);
    data->setImageData(pixmap);
    drag->setMimeData(data);

    setState(DraggingState);

    listModel->setDragKey(itemKey);

    Qt::DropAction dropAct = drag->exec(supportedActions);

    // 拖拽完成结束动画
    m_aniStartTime->stop();
    m_pressed = false;

    Q_EMIT dragEntered();
    if (dropAct == Qt::IgnoreAction) {
        QPropertyAnimation *posAni = new QPropertyAnimation(pixLabel, "pos", pixLabel);
        connect(posAni, &QPropertyAnimation::finished, [&, listModel, pixLabel] () {
            pixLabel->hide();
            pixLabel->deleteLater();
            listModel->setDragKey(QString());
            clearDragModelIndex();

            m_dropPos = QPoint();
            m_dragPos = QPoint();
        });
        posAni->setEasingCurve(QEasingCurve::Linear);
        posAni->setDuration(m_aniDuringTime);
        posAni->setStartValue((QCursor::pos() - QPoint(0, pixLabel->height() / 2)));
        posAni->setEndValue(mapToGlobal(m_dropPos) - QPoint(0, pixLabel->height() / 2));
        pixLabel->show();
        posAni->start(QAbstractAnimation::DeleteWhenStopped);
    } else {
        listModel->setDragKey(QString());
        clearDragModelIndex();

        m_dropPos = QPoint();
        m_dragPos = QPoint();

        Q_EMIT requestRemove(itemKey);
    }
}

void TrayGridView::initUi()
{
    setAcceptDrops(true);
    setDragEnabled(true);
    setDragDropMode(QAbstractItemView::DragDrop);
    setDropIndicatorShown(false);

    setMouseTracking(false);
    setUniformItemSizes(true);
    setFocusPolicy(Qt::NoFocus);
    setMovement(DListView::Free);
    setOrientation(QListView::LeftToRight, true);
    setLayoutMode(DListView::Batched);
    setResizeMode(DListView::Adjust);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setFrameStyle(QFrame::NoFrame);
    setContentsMargins(0, 0, 0, 0);
    setSpacing(0);
    setItemSpacing(0);
    setBackgroundType(DStyledItemDelegate::RoundedBackground);
    setSelectionMode(QListView::SingleSelection);
    setVerticalScrollMode(QListView::ScrollPerPixel);

    viewport()->setAcceptDrops(true);
    viewport()->setAutoFillBackground(false);

    m_aniStartTime->setInterval(10);
    m_aniStartTime->setSingleShot(true);

    connect(m_aniStartTime, &QTimer::timeout, this, &TrayGridView::moveAnimation);
}
