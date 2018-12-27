#include "fashiontrayholdcontainer.h"
#include "fashiontrayitem.h"

#define SpliterSize 2
#define TraySpace 10

FashionTrayHoldContainer::FashionTrayHoldContainer(TrayPlugin *trayPlugin, QWidget *parent)
    : QWidget(parent),
      m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight)),
      m_holdSpliter(new QLabel),
      m_trayPlugin(trayPlugin),
      m_dockPosistion(trayPlugin->dockPosition())
{
    setAcceptDrops(true);

    m_holdSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");

    m_mainBoxLayout->setMargin(0);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(TraySpace);

    m_mainBoxLayout->addWidget(m_holdSpliter);

    m_mainBoxLayout->setAlignment(Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_holdSpliter, Qt::AlignCenter);

    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    setLayout(m_mainBoxLayout);
}

void FashionTrayHoldContainer::setDockPostion(Dock::Position pos)
{
    m_dockPosistion = pos;

    if (pos == Dock::Position::Top || pos == Dock::Position::Bottom) {
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    } else{
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::TopToBottom);
    }
}

void FashionTrayHoldContainer::setTrayExpand(const bool expand)
{
    m_expand = expand;

    // 将显示与隐藏放在 timer 里做以避免收起动画的一些抖动发生
    QTimer::singleShot(200, this, [=] {
        // 这行代码的逻辑与下面被注释掉的部分相同
        setVisible(!(!m_expand && m_holdWrapperList.isEmpty()));

//        if (m_expand) {
//            setVisible(true);
//        } else {
//            if (m_holdWrapperList.isEmpty()) {
//                setVisible(false);
//            } else {
//                setVisible(true);
//            }
//        }
    });
}

bool FashionTrayHoldContainer::exists(FashionTrayWidgetWrapper *wrapper) const
{
    for (auto w : m_holdWrapperList) {

    }
    return false;
}

bool FashionTrayHoldContainer::isHoldTrayWrapper(FashionTrayWidgetWrapper *wrapper) const
{
    const QString &key = QString("%1_hold").arg(wrapper->itemKey());
    return m_trayPlugin->getValue(key, false).toBool();
}

void FashionTrayHoldContainer::addTrayWrapper(FashionTrayWidgetWrapper *wrapper)
{
    const QString &key = QString("%1_hold").arg(wrapper->itemKey());
    m_trayPlugin->saveValue(key, true);

    const int index = whereToInsert(wrapper);
    m_mainBoxLayout->insertWidget(index, wrapper);
    m_holdWrapperList.insert(index, wrapper);
}

bool FashionTrayHoldContainer::removeTrayWrapper(FashionTrayWidgetWrapper *wrapper)
{
    return false;

    const QString &key = QString("%1_hold").arg(wrapper->itemKey());
    m_trayPlugin->saveValue(key, false);
}

int FashionTrayHoldContainer::whereToInsert(FashionTrayWidgetWrapper *wrapper) const
{
    if (m_holdWrapperList.isEmpty()) {
        return 0;
    }

    const int destSortKey = m_trayPlugin->itemSortKey(wrapper->itemKey());

    if (destSortKey < -1) {
        return 0;
    }
    if (destSortKey == -1) {
        return m_holdWrapperList.size();
    }

    // 当目标插入位置为列表的大小时将从最后面追加到列表中
    int destIndex = m_holdWrapperList.size();
    for (int i = 0; i < m_holdWrapperList.size(); ++i) {
        if (destSortKey > m_trayPlugin->itemSortKey(m_holdWrapperList.at(i)->itemKey())) {
            continue;
        }
        destIndex = i;
        break;
    }

    return destIndex;
}

QSize FashionTrayHoldContainer::sizeHint() const
{
    QSize size;

    const int TrayWidgetWidth = FashionTrayItem::trayWidgetWidth();
    const int TrayWidgetHeight = FashionTrayItem::trayWidgetHeight();

    if (m_expand) {
        if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
            size.setWidth(
                        m_holdWrapperList.size() * TrayWidgetWidth // 所有保留显示的托盘图标
                        + SpliterSize // 一个分隔条
                        + (m_holdWrapperList.size() + 1) * TraySpace // 所有托盘图标之间的 space + 一个分隔条的 space
                        );
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(
                        m_holdWrapperList.size() * TrayWidgetHeight // 所有保留显示的托盘图标
                        + SpliterSize // 一个分隔条
                        + (m_holdWrapperList.size() + 1) * TraySpace // 所有托盘图标之间的 space + 一个分隔条的 space
                        );
        }
    } else {
        if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
            size.setWidth(
                        m_holdWrapperList.size() * TrayWidgetWidth // 所有保留显示的托盘图标
                        + m_holdWrapperList.size() * TraySpace // 所有托盘图标之间的 space
                        );
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(
                        m_holdWrapperList.size() * TrayWidgetHeight // 所有保留显示的托盘图标
                        + m_holdWrapperList.size() * TraySpace // 所有托盘图标之间的 space
                        );
        }
    }

    return size;
}

void FashionTrayHoldContainer::resizeEvent(QResizeEvent *event)
{
    const QSize &mSize = event->size();

    if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
        m_holdSpliter->setFixedSize(SpliterSize, mSize.height() * 0.3);
    } else{
        m_holdSpliter->setFixedSize(mSize.width() * 0.3, SpliterSize);
    }

    QWidget::resizeEvent(event);
}
