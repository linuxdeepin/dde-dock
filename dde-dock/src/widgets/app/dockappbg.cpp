#include "dockappbg.h"

BGActiveIndicator::BGActiveIndicator(QWidget *parent)
    :QLabel(parent)
{
    setObjectName("AppBackgroundActiveLabel");
    setAlignment(Qt::AlignBottom | Qt::AlignHCenter);
}

void BGActiveIndicator::showActiveWithAnimation()
{
    if (m_loopCount != 0)
        return;
    m_loopCount = 0;
    setFixedSize(28, 13);
    emit sizeChange();
    setVisible(true);
    m_iconPath = m_openingIndicatorIcon;
    QPropertyAnimation *animation = new QPropertyAnimation(this, "opacity");
    animation->setDuration(500);
    animation->setStartValue(0);
    animation->setEndValue(1);
    animation->start();
    connect(animation, &QPropertyAnimation::finished, [=]{
        ++ m_loopCount;
        if (m_loopCount == 4){
            m_loopCount = 0;
            emit showAnimationFinish();
        }
        else{
            if (m_loopCount % 2 == 0){
                animation->setStartValue(0);
                animation->setEndValue(1);
                animation->start();
            }
            else{
                animation->setStartValue(1);
                animation->setEndValue(0);
                animation->start();
            }
        }
    });

}

void BGActiveIndicator::show()
{
    QLabel::show();
    setFixedSize(24, 5);
    setOpacity(1);
    m_iconPath = m_openIndicatorIcon;
    update();
    emit sizeChange();
}
double BGActiveIndicator::opacity() const
{
    return m_opacity;
}

void BGActiveIndicator::setOpacity(double opacity)
{
    m_opacity = opacity;
    update();
}

void BGActiveIndicator::paintEvent(QPaintEvent *event)
{
    if (m_iconPath.isEmpty()){
        QLabel::paintEvent(event);
        return;
    }
    QPainter painter;
    painter.begin(this);

    painter.setClipRect(rect());
    painter.setOpacity(m_opacity);
    painter.drawPixmap(0, 0, QPixmap(m_iconPath).scaled(size()));

    painter.end();
}
QString BGActiveIndicator::openingIndicatorIcon() const
{
    return m_openingIndicatorIcon;
}

void BGActiveIndicator::setOpeningIndicatorIcon(const QString &openingIndicatorIcon)
{
    m_openingIndicatorIcon = openingIndicatorIcon;
}

QString BGActiveIndicator::openIndicatorIcon() const
{
    return m_openIndicatorIcon;
}

void BGActiveIndicator::setOpenIndicatorIcon(const QString &openIndicatorIcon)
{
    m_openIndicatorIcon = openIndicatorIcon;
}


DockAppBG::DockAppBG(QWidget *parent) :
    QLabel(parent)
{
    this->setObjectName("AppBackground");
    initActiveLabel();
}

void DockAppBG::resize(int width, int height)
{
    QLabel::resize(width, height);
    updateActiveLabelPos();
}

bool DockAppBG::getIsActived()
{
    return m_isActived;
}

void DockAppBG::setIsActived(bool value)
{
    m_isActived = value;
    if (!m_isActived) {
        m_activeLabel->hide();
        m_bePress = false;
    }
    else if (!m_bePress && getIsFashionMode()){
        m_activeLabel->show();
    }

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}

bool DockAppBG::getIsCurrentOpened()
{
    return m_isCurrentOpened;
}

void DockAppBG::setIsCurrentOpened(bool value)
{
    m_isCurrentOpened = value;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}

bool DockAppBG::getIsHovered()
{
    return m_isHovered;
}

void DockAppBG::setIsHovered(bool value)
{
    m_isHovered = value;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}
bool DockAppBG::getIsFashionMode() const
{
    return DockModeData::instance()->getDockMode() == Dock::FashionMode;
}

void DockAppBG::slotMouseRelease(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    m_bePress = true;
    if (!m_isActived && getIsFashionMode())
        m_activeLabel->showActiveWithAnimation();
}

void DockAppBG::initActiveLabel()
{
    m_activeLabel = new BGActiveIndicator(this);
    connect(m_activeLabel, &BGActiveIndicator::sizeChange, this, &DockAppBG::updateActiveLabelPos);
    connect(DockModeData::instance(), &DockModeData::dockModeChanged, this, &DockAppBG::onDockModeChanged);
    connect(m_activeLabel, &BGActiveIndicator::showAnimationFinish, [=]{
        if (m_isActived)
            m_activeLabel->show();
        m_bePress = false;
    });
}

void DockAppBG::updateActiveLabelPos()
{
    if (m_activeLabel)
        m_activeLabel->move((width() - m_activeLabel->width()) / 2, height() - m_activeLabel->height());
}

void DockAppBG::onDockModeChanged()
{
    if (m_activeLabel && !getIsFashionMode())
        m_activeLabel->hide();
    else if (m_activeLabel && m_isActived)
        m_activeLabel->show();
}

