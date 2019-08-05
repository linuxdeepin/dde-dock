#include "abstractcontainer.h"
#include "../fashiontrayconstants.h"

AbstractContainer::AbstractContainer(TrayPlugin *trayPlugin, QWidget *parent)
    : QWidget(parent),
      m_trayPlugin(trayPlugin),
      m_wrapperLayout(new QBoxLayout(QBoxLayout::LeftToRight)),
      m_currentDraggingWrapper(nullptr),
      m_expand(true),
      m_dockPosition(Dock::Position::Bottom),
      m_wrapperSize(QSize(TrayWidgetWidthMin, TrayWidgetHeightMin))
{
    setAcceptDrops(true);

    m_wrapperLayout->setMargin(0);
    m_wrapperLayout->setContentsMargins(0, 0, 0, 0);
    m_wrapperLayout->setSpacing(TraySpace);

    m_wrapperLayout->setAlignment(Qt::AlignCenter);

    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    setLayout(m_wrapperLayout);
}

void AbstractContainer::addWrapper(FashionTrayWidgetWrapper *wrapper)
{
    if (containsWrapper(wrapper)) {
        return;
    }

    const int index = whereToInsert(wrapper);
    m_wrapperLayout->insertWidget(index, wrapper);
    m_wrapperList.insert(index, wrapper);

    wrapper->setAttention(false);
    wrapper->setFixedSize(m_wrapperSize);

    connect(wrapper, &FashionTrayWidgetWrapper::attentionChanged, this, &AbstractContainer::onWrapperAttentionhChanged, static_cast<Qt::ConnectionType>(Qt::QueuedConnection | Qt::UniqueConnection));
    connect(wrapper, &FashionTrayWidgetWrapper::dragStart, this, &AbstractContainer::onWrapperDragStart, Qt::UniqueConnection);
    connect(wrapper, &FashionTrayWidgetWrapper::dragStop, this, &AbstractContainer::onWrapperDragStop, Qt::UniqueConnection);
    connect(wrapper, &FashionTrayWidgetWrapper::requestSwapWithDragging, this, &AbstractContainer::onWrapperRequestSwapWithDragging, Qt::UniqueConnection);

    refreshVisible();
}

bool AbstractContainer::removeWrapper(FashionTrayWidgetWrapper *wrapper)
{
    FashionTrayWidgetWrapper *w = takeWrapper(wrapper);
    if (!w) {
        return false;
    }

    // do not delete real tray object, just delete it's wrapper object
    // the real tray object should be deleted in TrayPlugin class
    w->absTrayWidget()->setParent(nullptr);
    w->deleteLater();

    refreshVisible();

    return true;
}

bool AbstractContainer::removeWrapperByTrayWidget(AbstractTrayWidget *trayWidget)
{
    FashionTrayWidgetWrapper *w = wrapperByTrayWidget(trayWidget);
    if (!w) {
        return false;
    }

    return removeWrapper(w);
}

FashionTrayWidgetWrapper *AbstractContainer::takeWrapper(FashionTrayWidgetWrapper *wrapper)
{
    if (!containsWrapper(wrapper)) {
        return nullptr;
    }

    if (m_currentDraggingWrapper == wrapper) {
        m_currentDraggingWrapper = nullptr;
    }

    wrapper->disconnect();
    m_wrapperLayout->removeWidget(wrapper);
    m_wrapperList.removeAll(wrapper);

    refreshVisible();

    return wrapper;
}

void AbstractContainer::setDockPosition(const Dock::Position pos)
{
    m_dockPosition = pos;

    if (pos == Dock::Position::Top || pos == Dock::Position::Bottom) {
        m_wrapperLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    } else{
        m_wrapperLayout->setDirection(QBoxLayout::Direction::TopToBottom);
    }

    refreshVisible();
}

void AbstractContainer::setExpand(const bool expand)
{
    m_expand = expand;

    refreshVisible();
}

QSize AbstractContainer::totalSize() const
{
    QSize size;

    const int wrapperWidth = m_wrapperSize.width();
    const int wrapperHeight = m_wrapperSize.height();

    if (m_dockPosition == Dock::Position::Top || m_dockPosition == Dock::Position::Bottom) {
        size.setWidth(
                    m_wrapperList.size() * wrapperWidth // 所有托盘图标
                    + m_wrapperList.size() * TraySpace // 所有托盘图标之间 + 一个尾部的 space
                    );
        size.setHeight(height());
    } else {
        size.setWidth(width());
        size.setHeight(
                    m_wrapperList.size() * wrapperHeight // 所有托盘图标
                    + m_wrapperList.size() * TraySpace // 所有托盘图标之间 + 一个尾部的 space
                    );
    }

    return size;
}

QSize AbstractContainer::sizeHint() const
{
    return totalSize();
}

void AbstractContainer::clearWrapper()
{
    QList<QPointer<FashionTrayWidgetWrapper>> mList = m_wrapperList;

    for (auto wrapper : mList) {
        removeWrapper(wrapper);
    }

    m_wrapperList.clear();

    refreshVisible();
}

void AbstractContainer::saveCurrentOrderToConfig()
{
    for (int i = 0; i < m_wrapperList.size(); ++i) {
        m_trayPlugin->setSortKey(m_wrapperList.at(i)->itemKey(), i + 1);
    }
}

void AbstractContainer::setWrapperSize(QSize size)
{
    m_wrapperSize = size;

    for (auto w : m_wrapperList) {
        w->setFixedSize(size);
    }
}

bool AbstractContainer::isEmpty()
{
    return m_wrapperList.isEmpty();
}

bool AbstractContainer::containsWrapper(FashionTrayWidgetWrapper *wrapper)
{
    for (auto w : m_wrapperList) {
        if (w == wrapper) {
            return true;
        }
    }

    return false;
}

bool AbstractContainer::containsWrapperByTrayWidget(AbstractTrayWidget *trayWidget)
{
    if (wrapperByTrayWidget(trayWidget)) {
        return true;
    }

    return false;
}

FashionTrayWidgetWrapper *AbstractContainer::wrapperByTrayWidget(AbstractTrayWidget *trayWidget)
{
    for (auto w : m_wrapperList) {
        if (w->absTrayWidget() == trayWidget) {
            return w;
        }
    }

    return nullptr;
}

void AbstractContainer::addDraggingWrapper(FashionTrayWidgetWrapper *wrapper)
{
    addWrapper(wrapper);

    if (containsWrapper(wrapper)) {
        m_currentDraggingWrapper = wrapper;
    }
}

FashionTrayWidgetWrapper *AbstractContainer::takeDraggingWrapper()
{
    if (!m_currentDraggingWrapper) {
        return nullptr;
    }

    return takeWrapper(m_currentDraggingWrapper);
}

int AbstractContainer::whereToInsert(FashionTrayWidgetWrapper *wrapper)
{
    if (m_wrapperList.isEmpty()) {
        return 0;
    }

    //根据配置文件记录的顺序排序
    const int destSortKey = m_trayPlugin->itemSortKey(wrapper->itemKey());

    if (destSortKey < -1) {
        return 0;
    }
    if (destSortKey == -1) {
        return m_wrapperList.size();
    }

    // 当目标插入位置为列表的大小时将从最后面追加到列表中
    int destIndex = m_wrapperList.size();
    for (int i = 0; i < m_wrapperList.size(); ++i) {
        if (destSortKey > m_trayPlugin->itemSortKey(m_wrapperList.at(i)->itemKey())) {
            continue;
        }
        destIndex = i;
        break;
    }

    return destIndex;
}

TrayPlugin *AbstractContainer::trayPlugin() const
{
    return m_trayPlugin;
}

QList<QPointer<FashionTrayWidgetWrapper> > AbstractContainer::wrapperList() const
{
    return m_wrapperList;
}

QBoxLayout *AbstractContainer::wrapperLayout() const
{
    return m_wrapperLayout;
}

// replace current WrapperLayout by "layout"
// but will not setLayout here, so the caller should handle the new WrapperLayout
void AbstractContainer::setWrapperLayout(QBoxLayout *layout)
{
    delete m_wrapperLayout;
    m_wrapperLayout = layout;
}

bool AbstractContainer::expand() const
{
    return m_expand;
}

Dock::Position AbstractContainer::dockPosition() const
{
    return m_dockPosition;
}

QSize AbstractContainer::wrapperSize() const
{
    return m_wrapperSize;
}

void AbstractContainer::dragEnterEvent(QDragEnterEvent *event)
{
    if (event->mimeData()->hasFormat(TRAY_ITEM_DRAG_MIMEDATA) && !m_currentDraggingWrapper) {
        event->accept();
        Q_EMIT requestDraggingWrapper();
        return;
    }

    QWidget::dragEnterEvent(event);
}

void AbstractContainer::onWrapperAttentionhChanged(const bool attention)
{
    FashionTrayWidgetWrapper *wrapper = dynamic_cast<FashionTrayWidgetWrapper *>(sender());
    if (!wrapper) {
        return;
    }

    Q_EMIT attentionChanged(wrapper, attention);
}

void AbstractContainer::onWrapperDragStart()
{
    FashionTrayWidgetWrapper *wrapper = static_cast<FashionTrayWidgetWrapper *>(sender());

    if (!wrapper) {
        return;
    }

    m_currentDraggingWrapper = wrapper;

    Q_EMIT draggingStateChanged(wrapper, true);
}

void AbstractContainer::onWrapperDragStop()
{
    FashionTrayWidgetWrapper *wrapper = static_cast<FashionTrayWidgetWrapper *>(sender());

    if (!wrapper) {
        return;
    }

    if (m_currentDraggingWrapper == wrapper) {
        m_currentDraggingWrapper = nullptr;
    } else {
        Q_UNREACHABLE();
    }

    saveCurrentOrderToConfig();

    Q_EMIT draggingStateChanged(wrapper, false);
}

void AbstractContainer::onWrapperRequestSwapWithDragging()
{
    FashionTrayWidgetWrapper *wrapper = static_cast<FashionTrayWidgetWrapper *>(sender());

    if (!wrapper || wrapper == m_currentDraggingWrapper) {
        return;
    }

    // the current dragging wrapper is null means that the dragging wrapper is contains by
    // another container, so this container need to emit requireDraggingWrapper() signal
    // to notify FashionTrayItem, the FashionTrayItem will move the dragging wrapper to this container
    if (!m_currentDraggingWrapper) {
        Q_EMIT requestDraggingWrapper();
        // here have to give up if dragging wrapper is still null
        if (!m_currentDraggingWrapper) {
            return;
        }
    }

    const int indexOfDest = m_wrapperLayout->indexOf(wrapper);
    const int indexOfDragging = m_wrapperLayout->indexOf(m_currentDraggingWrapper);

    m_wrapperLayout->removeWidget(m_currentDraggingWrapper);
    m_wrapperLayout->insertWidget(indexOfDest, m_currentDraggingWrapper);

    m_wrapperList.insert(indexOfDest, m_wrapperList.takeAt(indexOfDragging));
}
