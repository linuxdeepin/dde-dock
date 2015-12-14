#include <QLabel>
#include <QPixmap>
#include <QEvent>
#include <QTimer>
#include <QDebug>

#include "trayicon.h"
#include "compositetrayitem.h"

// these two variables are decided by the picture background.
static const int Margins = 4;
static const int ColumnWidth = 20;

CompositeTrayItem::CompositeTrayItem(QWidget *parent) :
    QFrame(parent),
    m_isCovered(true),
    m_isFolded(false)
{
    resize(1, 1);
    setObjectName("CompositeItem");

    m_cover = new QLabel(this);
    m_cover->setFixedSize(48, 48);
    m_cover->setPixmap(QPixmap(":/images/darea_cover.svg"));
    m_cover->move(QPoint(0, 0));

    m_coverTimer = new QTimer(this);
    m_coverTimer->setInterval(500);
    m_coverTimer->setSingleShot(true);

    m_updateTimer = new QTimer(this);
    m_updateTimer->setInterval(500);
    m_updateTimer->setSingleShot(false);
    m_updateTimer->start();


    m_foldButton = new DImageButton(":/images/fold-button-normal.svg",
                                    ":/images/fold-button-hover.svg",
                                    ":/images/fold-button-press.svg",
                                    this);
    m_foldButton->setFixedSize(18, 18);
    m_foldButton->hide();

    m_unfoldButton = new DImageButton(":/images/unfold-button-normal.svg",
                                      ":/images/unfold-button-hover.svg",
                                      ":/images/unfold-button-press.svg",
                                      this);
    m_unfoldButton->setFixedSize(18, 18);
    m_unfoldButton->hide();

    connect(m_coverTimer, &QTimer::timeout, this, &CompositeTrayItem::tryCoverOn);
    connect(m_updateTimer, &QTimer::timeout, this, &CompositeTrayItem::handleUpdateTimer);
    connect(m_foldButton, &DImageButton::clicked, this, &CompositeTrayItem::fold);
    connect(m_unfoldButton, &DImageButton::clicked, this, &CompositeTrayItem::unfold);
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
    m_coverTimer->stop();

    m_cover->raise();
    m_cover->setVisible(true);
    m_isCovered = true;
}

void CompositeTrayItem::coverOff()
{
    m_cover->lower();
    m_cover->setVisible(false);
    m_isCovered = false;
}

void CompositeTrayItem::tryCoverOn()
{
    QPoint globalPos = mapToGlobal(QPoint(0, 0));
    QRect globalGeometry(globalPos, size());

    if (!globalGeometry.contains(QCursor::pos()) &&
        (m_icons.keys().length() <= 4 || m_isFolded) &&
        m_mode == Dock::FashionMode)
    {
        coverOn();
    }
}

void CompositeTrayItem::handleTrayiconDamage()
{
    m_coverTimer->stop();

    unfold();

    QList<TrayIcon*> items = m_icons.values();
    for (int i = 0; i < items.length(); i++) {
        TrayIcon * icon = items.at(i);
        icon->updateIcon();
    }
}

void CompositeTrayItem::handleUpdateTimer()
{
    QList<TrayIcon*> items = m_icons.values();
    for (int i = 0; i < items.length(); i++) {
        TrayIcon * icon = items.at(i);
        icon->updateIcon();
    }
}

void CompositeTrayItem::resizeEvent(QResizeEvent * event)
{
    emit sizeChanged();

    QFrame::resizeEvent(event);
}

void CompositeTrayItem::enterEvent(QEvent * event)
{
    coverOff();

    QFrame::enterEvent(event);
}

void CompositeTrayItem::leaveEvent(QEvent * event)
{
    m_coverTimer->start();

    QFrame::leaveEvent(event);
}

void CompositeTrayItem::fold()
{
    m_isFolded = true;

    relayout();
}

void CompositeTrayItem::unfold()
{
    m_isFolded = false;

    coverOff();
    relayout();
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
        QList<TrayIcon*> items = m_icons.values();

        if (m_isFolded) {
            columnCount = 2;
        } else if (columnCount > 2 && childrenCount % 2 == 0) {
            columnCount += 1;
        }

        QString style = QString("QFrame#CompositeItem { background-image: url(':/images/darea_container_%1.svg') }").arg(columnCount * 2);
        setStyleSheet(style);

        resize(Margins * 2 + ColumnWidth * columnCount, 48);

        int placesCount = items.length();
        if (m_isFolded && placesCount > 3) { placesCount = 3;}

        for (int i = 0; i < items.length(); i++) {
            TrayIcon * icon = items.at(i);

            if (i < placesCount) {
                icon->maskOn();

                int x = i % columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;
                int y = i / columnCount * ColumnWidth + Margins + (ColumnWidth - 16) / 2;

                icon->move(x, y);
                icon->show();
                icon->updateIcon();
            } else {
                icon->hideIcon();
            }
        }

        if (columnCount > 2) {
            m_foldButton->move((columnCount - 1) * ColumnWidth + Margins + (ColumnWidth - 16) / 2,
                               ColumnWidth + Margins + (ColumnWidth - 16) / 2);
            m_foldButton->show();
            m_unfoldButton->hide();
        } else if (m_isFolded) {
            m_unfoldButton->move((columnCount - 1) * ColumnWidth + Margins + (ColumnWidth - 16) / 2,
                                 ColumnWidth + Margins + (ColumnWidth - 16) / 2);
            m_unfoldButton->show();
            m_foldButton->hide();
        } else {
            m_foldButton->hide();
            m_unfoldButton->hide();
        }

        if (m_isCovered) {
            m_cover->raise();
            m_cover->show();
        } else {
            m_cover->hide();
        }
    } else {
        m_cover->hide();

        setStyleSheet("");
        resize(childrenCount * Dock::APPLET_CLASSIC_ICON_SIZE + (childrenCount - 1) * Dock::APPLET_CLASSIC_ITEM_SPACING,
               Dock::APPLET_CLASSIC_ICON_SIZE);

        QList<TrayIcon*> items = m_icons.values();
        for (int i = 0; i < items.length(); i++) {
            TrayIcon * icon = items.at(i);
            icon->maskOff();

            icon->move(i * (Dock::APPLET_CLASSIC_ICON_SIZE + Dock::APPLET_CLASSIC_ITEM_SPACING), 0);
            icon->show();
            icon->updateIcon();
        }
    }
}
