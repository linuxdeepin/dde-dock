#include <QLabel>
#include <QPixmap>
#include <QEvent>
#include <QFile>
#include <QDebug>

#include "interfaces/dockconstants.h"

#include "trayicon.h"
#include "compositetrayitem.h"

// these two variables are decided by the picture background.
static const int Margins = 4;
static const int ColumnWidth = 20;

CompositeTrayItem::CompositeTrayItem(QWidget *parent) :
    QFrame(parent),
    m_isCovered(true)
{
    resize(1, 1);

    m_cover = new QLabel(this);
    m_cover->setFixedSize(48, 48);
    m_cover->setPixmap(QPixmap(":/images/darea_cover.svg"));
    m_cover->move(QPoint(0, 0));
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
    if (icon) {

        icon->setParent(NULL);
        icon->deleteLater();

        this->relayout();
    }
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

void CompositeTrayItem::clear()
{
    foreach (TrayIcon * icon, m_icons.values()) {
        icon->deleteLater();
    }
    m_icons.clear();
}

bool CompositeTrayItem::exist(const QString &key)
{
    return m_icons.keys().indexOf(key) != -1;
}

QStringList CompositeTrayItem::trayIds() const
{
    return m_icons.keys();
}

void CompositeTrayItem::coverOn()
{
    m_cover->raise();
    m_cover->setVisible(true);
}

void CompositeTrayItem::coverOff()
{
    m_cover->lower();
    m_cover->setVisible(false);
}

void CompositeTrayItem::enterEvent(QEvent * event)
{
    coverOff();

    QFrame::enterEvent(event);
}

void CompositeTrayItem::leaveEvent(QEvent * event)
{
    QPoint globalPos = mapToGlobal(QPoint(0, 0));
    QRect globalGeometry(globalPos, size());

    if (!globalGeometry.contains(QCursor::pos())) {
        coverOn();
    }

    QFrame::leaveEvent(event);
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
        QString style = QString("QFrame { background-image: url(':/images/darea_container_%1.svg') }").arg(columnCount * 2);
        setStyleSheet(style);

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

    if (m_isCovered) {
        m_cover->raise();
        m_cover->show();
    } else {
        m_cover->hide();
    }
}
