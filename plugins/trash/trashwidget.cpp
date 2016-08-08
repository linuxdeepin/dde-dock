
#include "constants.h"
#include "trashwidget.h"

#include <QPainter>
#include <QIcon>
#include <QApplication>

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent),

      m_popupApplet(new PopupControlWidget(this))
{
    QIcon::setThemeName("deepin");

    m_popupApplet->setVisible(false);

    connect(m_popupApplet, &PopupControlWidget::emptyChanged, this, &TrashWidget::updateIcon);

    updateIcon();
}

QWidget *TrashWidget::popupApplet()
{
    return m_popupApplet;
}

QSize TrashWidget::sizeHint() const
{
    return QSize(20, 20);
}

void TrashWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void TrashWidget::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    updateIcon();
}

void TrashWidget::updateIcon()
{
    QString iconString = "user-trash";
    if (!m_popupApplet->empty())
        iconString.append("-full");
    if (qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>() == Dock::Efficient)
        iconString.append("-symbolic");

    const int size = std::min(width(), height()) * 0.8;
    QIcon icon = QIcon::fromTheme(iconString);
    m_icon = icon.pixmap(size, size);

    update();
}
