#include "appbackground.h"

AppBackground::AppBackground(QWidget *parent) :
    QLabel(parent)
{
    this->setObjectName("AppBackground");
}

bool AppBackground::getIsActived()
{
    return m_isActived;
}

void AppBackground::setIsActived(bool value)
{
    m_isActived = value;

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
