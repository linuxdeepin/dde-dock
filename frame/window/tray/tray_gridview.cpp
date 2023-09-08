// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "tray_gridview.h"
#include "expandiconwidget.h"
#include "tray_model.h"
#include "basetraywidget.h"

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
#include <qwidget.h>

TrayGridView *TrayGridView::getDockTrayGridView(QWidget *parent)
{
    static TrayGridView *view = nullptr;
    if (!view)
        view = new TrayGridView(parent);
    return view;
}

TrayGridView *TrayGridView::getIconTrayGridView(QWidget *parent)
{
    static TrayGridView *view = nullptr;
    if (!view)
        view = new TrayGridView(parent);
    return view;
}

TrayGridView::TrayGridView(QWidget *parent)
    : DListView(parent)
    , m_aniCurveType(QEasingCurve::Linear)
    , m_aniDuringTime(250)
    , m_dragDistance(15)
    , m_aniStartTime(new QTimer(this))
    , m_pressed(false)
    , m_aniRunning(false)
    , m_positon(Dock::Position::Bottom)
{
    initUi();
}

void TrayGridView::setPosition(Dock::Position position)
{
    if (m_positon == position) {
        return;
    }
    m_positon = position;
    setOrientation(m_positon == Dock::Position::Top || m_positon == Dock::Position::Bottom ?
        QListView::Flow::LeftToRight : QListView::Flow::TopToBottom, false);
}

Dock::Position TrayGridView::position() const
{
    return m_positon;
}

QSize TrayGridView::suitableSize() const
{
    return suitableSize(m_positon);
}

QSize TrayGridView::suitableSize(const Dock::Position &position) const
{
    TrayModel *dataModel = qobject_cast<TrayModel *>(model());
    if (!dataModel)
        return QSize(-1, -1);

    if (dataModel->isIconTray()) {
        // 如果是托盘图标
        int width = 2;
        int height = 0;
        int count = dataModel->rowCount();
        if (count > 0) {
            int columnCount = qMin(count, 3);
            for (int i = 0; i < columnCount; i ++) {
                QModelIndex index = dataModel->index(i, 0);
                width += visualRect(index).width() + spacing() * 2;     // 左右边距加上单元格宽度
            }
            int rowCount = count / 3;
            if (count % 3 > 0)
                rowCount++;
            for (int i = 0; i < rowCount; i++) {
                QModelIndex index = dataModel->index(i * 3);
                height += visualRect(index).height() + spacing() * 2;
            }
        } else {
            width = spacing() * 2 + 30;
            height = spacing() * 2 + 30;
        }
        return QSize(width, height);
    }
    if (position == Dock::Position::Top || position == Dock::Position::Bottom) {
        int length = spacing() + 2;
        if (m_positon == Dock::Position::Top || m_positon == Dock::Position::Bottom) {
            for (int i = 0; i < dataModel->rowCount(); i++) {
                QModelIndex index = dataModel->index(i, 0);
                QRect indexRect = visualRect(index);
                length += indexRect.width() + spacing();
            }
        } else {
            // 如果是从左右切换过来的，此时还未进入上下位置，则将当前位置的高度作为计算左右位置的宽度
            for (int i = 0; i < dataModel->rowCount(); i++) {
                QModelIndex index = dataModel->index(i, 0);
                QRect indexRect = visualRect(index);
                length += indexRect.height() + spacing();
            }
        }

        return QSize(length, -1);
    }
    int height = spacing() + 2;
    if (m_positon == Dock::Position::Left || m_positon == Dock::Position::Right) {
        for (int i = 0; i < dataModel->rowCount(); i++) {
            QModelIndex index = dataModel->index(i, 0);
            QRect indexRect = visualRect(index);
            height += indexRect.height() + spacing();
        }
    } else {
        for (int i = 0; i < dataModel->rowCount(); i++) {
            QModelIndex index = dataModel->index(i, 0);
            QRect indexRect = visualRect(index);
            height += indexRect.width() + spacing();
        }
    }

    return QSize(-1, height);
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
    if (!m_pressed)
        return DListView::mouseMoveEvent(e);

    setState(QAbstractItemView::NoState);
    e->accept();

    if (e->buttons() == Qt::RightButton)
        return DListView::mouseMoveEvent(e);

    QModelIndex index = indexAt(e->pos());
    if (!index.isValid())
        return DListView::mouseMoveEvent(e);

    // 如果当前拖动的位置是托盘展开按钮，则不让其拖动
    TrayIconType iconType = index.data(TrayModel::Role::TypeRole).value<TrayIconType>();
    if (iconType == TrayIconType::ExpandIcon)
        return DListView::mouseMoveEvent(e);

    if ((qAbs(e->pos().x() - m_dragPos.x()) > m_dragDistance ||
                      qAbs(e->pos().y() - m_dragPos.y()) > m_dragDistance)) {
        qDebug() << "start drag";
        if (!beginDrag(Qt::CopyAction | Qt::MoveAction))
            DListView::mouseMoveEvent(e);
    }
}

void TrayGridView::mouseReleaseEvent(QMouseEvent *e)
{
    Q_UNUSED(e);
    m_pressed = false;
}

void TrayGridView::dragEnterEvent(QDragEnterEvent *e)
{
    const QModelIndex index = indexAt(e->pos());

    if (model()->canDropMimeData(e->mimeData(), e->dropAction(), index.row(),
                                 index.column(), index))
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

    QModelIndex index = indexAt(e->pos());
    if (!model()->canDropMimeData(e->mimeData(), e->dropAction(), index.row(),
                                 index.column(), index))
        return;

    setState(QAbstractItemView::DraggingState);

    if (index.isValid()) {
        if (m_dropPos != indexRect(index).center()) {
            qDebug() << "update drop position: " << index.row();
            m_dropPos = indexRect(index).center();
        }
    }

    if (m_pressed)
        m_aniStartTime->start();
}

const QModelIndex TrayGridView::getIndexFromPos(QPoint currentPoint) const
{
    QModelIndex index = indexAt(currentPoint);
    if (index.isValid())
        return index;

    if (model()->rowCount() == 0)
        return index;

    // 如果在第一个之前，则认为拖到了第一个的位置
    QRect indexRect0 = visualRect(model()->index(0, 0));
    if (currentPoint.x() < indexRect0.x() || currentPoint.y() < indexRect0.y())
        return model()->index(0, 0);

    // 如果从指定的位置没有找到索引，则依次从每个index中查找，先横向查找
    for (int i = 1; i < model()->rowCount(); i++) {
        QModelIndex lastIndex = model()->index(i - 1, 0);
        QModelIndex currentIndex = model()->index(i, 0);
        QRect lastIndexRect = visualRect(lastIndex);
        QRect indexRect = visualRect(currentIndex);
        if (lastIndexRect.x() + lastIndexRect.width() <= currentPoint.x()
                && indexRect.x() >= currentPoint.x())
            return currentIndex;
    }
    // 如果鼠标位置刚好在上下两个索引中间
    for (int i = 0; i < model()->rowCount(); i++) {
        QModelIndex currentIndex = model()->index(i, 0);
        QRect indexRect = visualRect(currentIndex);

        if (currentPoint.y() >= indexRect.y() - spacing() && currentPoint.y() < indexRect.y()
                && currentPoint.x() >= indexRect.x() - spacing() && currentPoint.x() < indexRect.x())
            return currentIndex;
    }

    return QModelIndex();
}

bool TrayGridView::mouseInDock()
{
    QPoint mousePosition = QCursor::pos();
    QRect dockRect(topLevelWidget()->pos(), topLevelWidget()->size());
    switch (m_positon) {
    case Dock::Position::Bottom: {
        return mousePosition.y() > dockRect.top();
    }
    case Dock::Position::Left: {
        return mousePosition.x() < dockRect.right();
    }
    case Dock::Position::Top: {
        return mousePosition.y() < dockRect.bottom();
    }
    case Dock::Position::Right: {
        return mousePosition.x() > dockRect.left();
    }
    }
    return false;
}

void TrayGridView::handleDropEvent(QDropEvent *e)
{
    setState(DListView::NoState);
    clearDragModelIndex();

    if (m_aniStartTime->isActive())
        m_aniStartTime->stop();

    if (e->mimeData()->formats().contains("type") && e->source() != this) {
        e->setDropAction(Qt::CopyAction);
        e->accept();

        TrayModel *dataModel = qobject_cast<TrayModel *>(model());
        if (dataModel) {
            WinInfo info;
            info.type = static_cast<TrayIconType>(e->mimeData()->data("type").toInt());
            info.key = static_cast<QString>(e->mimeData()->data("key"));
            info.winId = static_cast<quint32>(e->mimeData()->data("winId").toInt());
            info.servicePath = static_cast<QString>(e->mimeData()->data("servicePath"));
            info.itemKey = static_cast<QString>(e->mimeData()->data("itemKey"));
            info.isTypeWriting = (static_cast<QString>(e->mimeData()->data("isTypeWritting")) == "1");
            info.expand = (static_cast<QString>(e->mimeData()->data("expand")) == "1");
            info.pluginInter = (PluginsItemInterface *)(e->mimeData()->imageData().value<qulonglong>());
            QModelIndex targetIndex = getIndexFromPos(e->pos());
            int index = -1;
            if (targetIndex.isValid() && targetIndex.row() < dataModel->rowCount()) {
                // 如果拖动的位置是合法的位置，则让其插入到当前的位置
                index = targetIndex.row();
                dataModel->insertRow(index, info);
            } else {
                // 在其他的情况下，让其插入到最后
                dataModel->addRow(info);
            }

            dataModel->saveConfig(index, info);
        }
    } else {
        e->ignore();
        DListView::dropEvent(e);
    }
}

void TrayGridView::onUpdateEditorView()
{
    for (int i = 0; i < model()->rowCount(); i++) {
        QModelIndex index = model()->index(i, 0);
        closePersistentEditor(index);
    }
    // 在关闭QWidget后不要立即调用openPersistentEditor来打开
    // 因为closePersistentEditor后，异步删除QWidget，在关闭后，如果立即调用openPersistentEditor，在删除的时候，会把
    // 通过openPersistentEditor新建的QWidget给删除，引起bug，因此，在所有的都closePersistentEditor后，异步来调用
    // openPersistentEditor就不会出现这种问题
    QMetaObject::invokeMethod(this, [ = ] {
        for (int i = 0; i < model()->rowCount(); i++) {
            QModelIndex index = model()->index(i, 0);
            openPersistentEditor(index);
        }
    }, Qt::QueuedConnection);
}

bool TrayGridView::beginDrag(Qt::DropActions supportedActions)
{
    QModelIndex modelIndex = indexAt(m_dragPos);
    TrayIconType trayType = modelIndex.data(TrayModel::Role::TypeRole).value<TrayIconType>();
    // 展开图标不能移动
    if (trayType == TrayIconType::ExpandIcon)
        return false;

    m_dropPos = indexRect(modelIndex).center();

    TrayModel *listModel = qobject_cast<TrayModel *>(model());
    if (!listModel)
        return false;

    BaseTrayWidget *widget = qobject_cast<BaseTrayWidget *>(indexWidget(modelIndex));
    if (!widget)
        return false;

    QMimeData *data = model()->mimeData(QModelIndexList() << modelIndex);
    if (!data)
        return false;

    auto pixmap = widget->icon();

    qreal ratio = qApp->devicePixelRatio();
    // 创建拖拽释放时的应用图标
    QLabel *pixLabel = new QLabel(this);
    pixLabel->setPixmap(pixmap);
    pixLabel->setFixedSize(indexRect(modelIndex).size() / ratio);

    QDrag *drag = new QDrag(this);
    pixmap.scaled(pixmap.size() * ratio, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    pixmap.setDevicePixelRatio(ratio);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center() / ratio);

    data->setImageData(pixmap);
    drag->setMimeData(data);

    setState(DraggingState);

    listModel->setDragKey(modelIndex.data(TrayModel::Role::KeyRole).toString());
    listModel->setDragingIndex(modelIndex);
    // 删除当前的图标
    WinInfo winInfo = listModel->getWinInfo(modelIndex);

    Qt::DropAction dropAct = drag->exec(supportedActions);

    // 拖拽完成结束动画
    m_aniStartTime->stop();
    m_pressed = false;

    if (dropAct == Qt::IgnoreAction) {
        if (listModel->isIconTray()) {
            // 如果当前是从托盘区域释放，按照原来的流程走
            QPropertyAnimation *posAni = new QPropertyAnimation(pixLabel, "pos", pixLabel);
            connect(posAni, &QPropertyAnimation::finished, [ this, listModel, pixLabel, winInfo ] () {
                pixLabel->hide();
                pixLabel->deleteLater();
                listModel->setDragKey(QString());
                clearDragModelIndex();
                QModelIndex dropIndex = indexAt(m_dropPos);
                // 拖转完成后，将拖动的图标插入到新的位置
                //listModel->moveToIndex(winInfo, dropIndex.row());
                listModel->dropSwap(dropIndex.row());
                listModel->setExpandVisible(!TrayModel::getIconModel()->isEmpty());

                m_dropPos = QPoint();
                m_dragPos = QPoint();

                onUpdateEditorView();
                Q_EMIT dragFinished();
            });
            // 拖拽完成后，将当前拖拽的item从原来的列表中移除，后来会根据实际情况将item插入到特定的列表中
            posAni->setEasingCurve(QEasingCurve::Linear);
            posAni->setDuration(m_aniDuringTime);
            posAni->setStartValue((QCursor::pos() - QPoint(0, pixLabel->height() / 2)));
            posAni->setEndValue(mapToGlobal(m_dropPos) - QPoint(0, pixLabel->height() / 2));
            pixLabel->show();
            posAni->start(QAbstractAnimation::DeleteWhenStopped);

            Q_EMIT dragFinished();
        } else {
            listModel->setDragKey(QString());
            clearDragModelIndex();
            TrayModel *trayModel = TrayModel::getIconModel();
            if (!mouseInDock()) {
                listModel->removeWinInfo(winInfo);
                trayModel->addRow(winInfo);
                trayModel->saveConfig(-1, winInfo);
            }
            // 如果是任务栏的的托盘区，则更新是否显示展开入口
            listModel->setExpandVisible(trayModel->rowCount() > 0, false);

            m_dragPos = QPoint();
            m_dropPos = QPoint();
            Q_EMIT dragFinished();
        }
    } else {
        // 拖拽完成后，将当前拖拽的item从原来的列表中移除，后来会根据实际情况将item插入到特定的列表中
        listModel->removeWinInfo(winInfo);
        // 这里是将图标从一个区域移动到另外一个区域
        listModel->setDragKey(QString());
        clearDragModelIndex();
        if (listModel->isIconTray()) {
            // 如果当前是从托盘移动到任务栏，则根据托盘内部是否有应用来决定是否显示展开图标
            bool hasIcon = (listModel->rowCount() > 0);
            TrayModel::getDockModel()->setExpandVisible(hasIcon, hasIcon);
            // 如果没有图标，则隐藏托盘图标
            if (!hasIcon)
                ExpandIconWidget::popupTrayView()->hide();
        }

        m_dropPos = QPoint();
        m_dragPos = QPoint();
        Q_EMIT dragFinished();
    }

    return true;
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

void TrayGridView::dropEvent(QDropEvent *e)
{
    handleDropEvent(e);
}
