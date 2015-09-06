#include "appbackground.h"


ActiveLabel::ActiveLabel(QWidget *parent)
    :QLabel(parent)
{
    setObjectName("AppBackgroundActiveLabel");
    setAlignment(Qt::AlignBottom | Qt::AlignHCenter);
}

void ActiveLabel::showActiveWithAnimation()
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

void ActiveLabel::show()
{
    QLabel::show();
    setFixedSize(24, 5);
    setOpacity(1);
    m_iconPath = m_openIndicatorIcon;
    repaint();
    emit sizeChange();
}
double ActiveLabel::opacity() const
{
    return m_opacity;
}

void ActiveLabel::setOpacity(double opacity)
{
    m_opacity = opacity;
    repaint();
}

void ActiveLabel::paintEvent(QPaintEvent *event)
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
QString ActiveLabel::openingIndicatorIcon() const
{
    return m_openingIndicatorIcon;
}

void ActiveLabel::setOpeningIndicatorIcon(const QString &openingIndicatorIcon)
{
    m_openingIndicatorIcon = openingIndicatorIcon;
}

QString ActiveLabel::openIndicatorIcon() const
{
    return m_openIndicatorIcon;
}

void ActiveLabel::setOpenIndicatorIcon(const QString &openIndicatorIcon)
{
    m_openIndicatorIcon = openIndicatorIcon;
}


AppBackground::AppBackground(QWidget *parent) :
    QLabel(parent)
{
    this->setObjectName("AppBackground");
    initActiveLabel();
}

void AppBackground::resize(int width, int height)
{
    QLabel::resize(width, height);
    updateActiveLabelPos();
}

bool AppBackground::getIsActived()
{
    return m_isActived;
}

void AppBackground::setIsActived(bool value)
{
    m_isActived = value;
    if (!m_isActived)
        m_activeLabel->hide();
    else if (m_isInit && !m_bePress && getIsFashionMode()){
        m_activeLabel->show();
        m_isInit = false;
    }

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}

bool AppBackground::getIsCurrentOpened()
{
    return m_isCurrentOpened;
}

void AppBackground::setIsCurrentOpened(bool value)
{
    m_isCurrentOpened = value;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}

bool AppBackground::getIsHovered()
{
    return m_isHovered;
}

void AppBackground::setIsHovered(bool value)
{
    m_isHovered = value;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}
bool AppBackground::getIsFashionMode() const
{
    return DockModeData::instance()->getDockMode() == Dock::FashionMode;
}

void AppBackground::slotMouseRelease(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    m_bePress = true;
    if (!m_isActived && getIsFashionMode())
        m_activeLabel->showActiveWithAnimation();
}

void AppBackground::initActiveLabel()
{
    m_activeLabel = new ActiveLabel(this);
    connect(m_activeLabel, &ActiveLabel::sizeChange, this, &AppBackground::updateActiveLabelPos);
    connect(m_activeLabel, &ActiveLabel::showAnimationFinish, [=]{
        if (m_isActived)
            m_activeLabel->show();
        m_bePress = false;
    });
    connect(DockModeData::instance(), &DockModeData::dockModeChanged, [=]{
        if (!getIsFashionMode())
            m_activeLabel->hide();
        else if (m_isActived)
            m_activeLabel->show();
    });
}

void AppBackground::updateActiveLabelPos()
{
    if (m_activeLabel)
        m_activeLabel->move((width() - m_activeLabel->width()) / 2, height() - m_activeLabel->height());
}

