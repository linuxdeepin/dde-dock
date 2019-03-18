#include "holdcontainer.h"
#include "../fashiontrayconstants.h"

HoldContainer::HoldContainer(TrayPlugin *trayPlugin, QWidget *parent)
    : AbstractContainer(trayPlugin, parent),
      m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight)),
      m_holdSpliter(new SpliterAnimated(this))
{
    m_mainBoxLayout->setMargin(0);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(TraySpace);

    QBoxLayout *preLayout = wrapperLayout();
    QBoxLayout *newLayout = new QBoxLayout(QBoxLayout::Direction::LeftToRight);
    for (int i = 0; i < preLayout->count(); ++i) {
        newLayout->addItem(preLayout->takeAt(i));
    }
    setWrapperLayout(newLayout);

    m_mainBoxLayout->addWidget(m_holdSpliter);
    m_mainBoxLayout->addLayout(newLayout);

    m_mainBoxLayout->setAlignment(m_holdSpliter, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(newLayout, Qt::AlignCenter);

    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    setLayout(m_mainBoxLayout);
}

bool HoldContainer::acceptWrapper(FashionTrayWidgetWrapper *wrapper)
{
    const QString &key = wrapper->absTrayWidget()->itemKeyForConfig() + HoldKeySuffix;

    return trayPlugin()->getValue(wrapper->itemKey(), key, false).toBool();
}

void HoldContainer::addWrapper(FashionTrayWidgetWrapper *wrapper)
{
    AbstractContainer::addWrapper(wrapper);

    if (containsWrapper(wrapper)) {
        const QString &key = wrapper->absTrayWidget()->itemKeyForConfig() + HoldKeySuffix;
        trayPlugin()->saveValue(wrapper->itemKey(), key, true);
    }
}

void HoldContainer::refreshVisible()
{
    setVisible(expand() || !isEmpty());
}

void HoldContainer::setDockPosition(const Dock::Position pos)
{
    if (pos == Dock::Position::Top || pos == Dock::Position::Bottom) {
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    } else{
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::TopToBottom);
    }

    m_holdSpliter->setDockPosition(pos);

    AbstractContainer::setDockPosition(pos);
}

void HoldContainer::setExpand(const bool expand)
{
    m_holdSpliter->setVisible(expand);

    AbstractContainer::setExpand(expand);
}

QSize HoldContainer::totalSize() const
{
    QSize size = AbstractContainer::totalSize();

    if (expand()) {
        if (dockPosition() == Dock::Position::Top || dockPosition() == Dock::Position::Bottom) {
            size.setWidth(
                        size.width()
                        + SpliterSize
                        + TraySpace
                        );
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(
                        size.height()
                        + SpliterSize
                        + TraySpace
                        );
        }
    }

    return size;
}

void HoldContainer::setDragging(const bool dragging)
{
    if (dragging) {
        m_holdSpliter->startAnimation();
    } else {
        m_holdSpliter->stopAnimation();
    }
}

void HoldContainer::resizeEvent(QResizeEvent *event)
{
    const QSize &mSize = event->size();
    const Dock::Position dockPosition = trayPlugin()->dockPosition();

    if (dockPosition == Dock::Position::Top || dockPosition == Dock::Position::Bottom) {
        m_holdSpliterMiniSize = QSize(SpliterSize, mSize.height() * 0.3);
        m_holdSpliterMaxSize = QSize(SpliterSize, mSize.height() * 0.5);
        m_holdSpliter->setFixedSize(SpliterSize, mSize.height());
    } else{
        m_holdSpliterMiniSize = QSize(mSize.width() * 0.3, SpliterSize);
        m_holdSpliterMaxSize = QSize(mSize.width() * 0.5, SpliterSize);
        m_holdSpliter->setFixedSize(mSize.width(), SpliterSize);
    }

    m_holdSpliter->setStartValue(m_holdSpliterMiniSize);
    m_holdSpliter->setEndValue(m_holdSpliterMaxSize);

    AbstractContainer::resizeEvent(event);
}
