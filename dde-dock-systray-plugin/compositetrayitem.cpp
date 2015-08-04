#include <QDebug>

#include <dock/dockconstants.h>

#include "trayicon.h"
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

void CompositeTrayItem::addTrayIcon(QString key, TrayIcon * icon)
{
    m_icons[key] = icon;

    icon->setParent(this);

    this->relayout();
}

void CompositeTrayItem::remove(QString key)
{
    TrayIcon * icon = m_icons.take(key);
    icon->setParent(NULL);
    icon->deleteLater();

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
    uint childrenCount = m_icons.keys().length();
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

        QList<TrayIcon*> items = m_icons.values();
        for (int i = 0; i < items.length(); i++) {
            TrayIcon * icon = items.at(i);
            icon->maskOn();

            int x = i % columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;
            int y = i / columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;

            icon->move(x, y);
            icon->show();
        }
    } else {
        setStyleSheet("");
        resize(childrenCount * Dock::APPLET_CLASSIC_ICON_SIZE + (childrenCount - 1) * Dock::APPLET_CLASSIC_ITEM_SPACING,
               Dock::APPLET_CLASSIC_ICON_SIZE);

        QList<TrayIcon*> items = m_icons.values();
        for (int i = 0; i < items.length(); i++) {
            TrayIcon * icon = items.at(i);
            icon->maskOff();

            icon->move(i * (Dock::APPLET_CLASSIC_ICON_SIZE + Dock::APPLET_CLASSIC_ITEM_SPACING), 0);
            icon->show();
        }
    }
}
