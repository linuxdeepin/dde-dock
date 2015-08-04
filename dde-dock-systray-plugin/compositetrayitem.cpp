#include <QDebug>

#include <dock/dockconstants.h>

#include "compositetrayitem.h"

// these two variables are decided by the picture background.
static const int Margins = 4;
static const int ColumnWidth = 20;

CompositeTrayItem::CompositeTrayItem(QWidget *parent) :
    QFrame(parent)
{
    resize(1, 1);
}

CompositeTrayItem::~CompositeTrayItem()
{
    qDebug() << "CompositeTrayItem destroyed.";
}

void CompositeTrayItem::addItem(QString key, QWidget * widget)
{
    m_items[key] = widget;

    widget->setParent(this);

    this->relayout();
}

void CompositeTrayItem::removeItem(QString key)
{
    QWidget * widget = m_items.take(key);
    widget->setParent(NULL);
    widget->deleteLater();

    this->relayout();
}

Dock::DockMode CompositeTrayItem::mode() const
{
    return m_mode;
}

void CompositeTrayItem::setMode(const Dock::DockMode &mode)
{
    if (m_mode != mode) {
        m_mode = mode;

        this->relayout();
    }
}

void CompositeTrayItem::relayout()
{
    uint childrenCount = m_items.keys().length();
    uint columnCount = 2;

    if (childrenCount <= 4) {
        columnCount = 2;
    } else if (childrenCount <= 6) {
        columnCount = 3;
    } else if (childrenCount <= 8) {
        columnCount = 4;
    } else if (childrenCount <= 10) {
        columnCount = 5;
    } else if (childrenCount <= 12) {
        columnCount = 6;
    }

    if (m_mode == Dock::FashionMode) {
        setStyleSheet("QFrame { background-image: url(':/images/darea_container_4.svg') }");
        resize(Margins * 2 + ColumnWidth * columnCount, 48);

        QList<QWidget*> items = m_items.values();
        for (int i = 0; i < items.length(); i++) {
            QWidget * widget = items.at(i);

            int x = i % columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;
            int y = i / columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;

            widget->move(x, y);
            widget->show();
        }
    } else {
        setStyleSheet("");
        resize(childrenCount * Dock::APPLET_CLASSIC_ICON_SIZE + (childrenCount - 1) * Dock::APPLET_CLASSIC_ITEM_SPACING,
               Dock::APPLET_CLASSIC_ICON_SIZE);

        QList<QWidget*> items = m_items.values();
        for (int i = 0; i < items.length(); i++) {
            QWidget * widget = items.at(i);

            widget->move(i * (Dock::APPLET_CLASSIC_ICON_SIZE + Dock::APPLET_CLASSIC_ITEM_SPACING), 0);
            widget->show();
        }
    }
}
