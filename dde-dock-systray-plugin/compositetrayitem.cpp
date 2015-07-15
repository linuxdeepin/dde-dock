#include <QStyle>
#include <QDebug>

#include <dock/dockconstants.h>

#include "compositetrayitem.h"

static const int Margins = 4;
static const int ColumnWidth = 20;

CompositeTrayItem::CompositeTrayItem(QWidget *parent) :
    QFrame(parent)
{
    m_columnCount = 2;

    setBackground();
}

CompositeTrayItem::~CompositeTrayItem()
{
    qDebug() << "CompositeTrayItem destroyed.";
}

void CompositeTrayItem::addWidget(QWidget * widget)
{
    widget->setParent(this);

    uint childrenCount = children().length();

    // update background and size
    if (childrenCount <= 4) {
        m_columnCount = 2;
    } else if (childrenCount <= 6) {
        m_columnCount = 3;
    } else if (childrenCount <= 8) {
        m_columnCount = 4;
    } else if (childrenCount <= 10) {
        m_columnCount = 5;
    } else if (childrenCount <= 12) {
        m_columnCount = 6;
    }

    setBackground();

    // move the widget to right position
    int x = (childrenCount - 1) % m_columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;
    int y = (childrenCount - 1) / m_columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;
    widget->move(x, y);
}

void CompositeTrayItem::removeWidget(QWidget *widget)
{
//    widget->setParent(NULL);

    uint childrenCount = children().length();

    // update background and size
    if (childrenCount <= 4) {
        m_columnCount = 2;
    } else if (childrenCount <= 6) {
        m_columnCount = 3;
    } else if (childrenCount <= 8) {
        m_columnCount = 4;
    } else if (childrenCount <= 10) {
        m_columnCount = 5;
    } else if (childrenCount <= 12) {
        m_columnCount = 6;
    }

    setBackground();
}


void CompositeTrayItem::setBackground()
{
    resize(Margins * 2 + ColumnWidth * m_columnCount, 48);
    setStyleSheet("QFrame { background-image: url(':/images/darea_container_4.svg') }");

    qDebug() << "CompositeTrayItem::setBackground()" << this->geometry();
}
