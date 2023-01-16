// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "overflowitem.h"
#include "diconbutton.h"
#include "itemconsts.h"
#include "appitem.h"

#include <DGuiApplicationHelper>

#include <QBoxLayout>
#include <QLabel>
#include <QScrollArea>
#include <QScrollBar>
const QString ICON_DEFAULT = QStringLiteral(":/icons/resources/application-x-desktop");
const QString OVERFLOW_MORE = QStringLiteral(":/icons/resources/overflow-more");

const QString ARROW_UP = QStringLiteral(":/icons/resources/arrow-up");
const QString ARROW_DOWN = QStringLiteral(":/icons/resources/arrow-down");
const QString ARROW_LEFT = QStringLiteral(":/icons/resources/arrow-left");
const QString ARROW_RIGHT = QStringLiteral(":/icons/resources/arrow-right");

// INFO: check if is darktheme
inline bool isDarkTheme() {
    return DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::DarkType;
}

OverflowItem::OverflowItem(QWidget *parent)
    : DockItem(parent)
    , m_width(0)
    , m_clicked(false)
    , m_showpopup(false)
    , m_scrollarea(new QScrollArea)
    , m_centerScroll(new QWidget)
    , m_popuplayout(new QBoxLayout(QBoxLayout::LeftToRight))
    , m_popupwindow(new DockPopupWindow)
    , m_popupbtnslayout(new QBoxLayout(QBoxLayout::LeftToRight, m_scrollarea))
    , m_left(new DIconButton)
    , m_right(new DIconButton)
{
    initUI();
    initSlots();
    setbtnsVisible();

    m_centerScroll->installEventFilter(this);
    m_scrollarea->installEventFilter(this);
}

void OverflowItem::initUI() {
    m_popupwindow->setShadowBlurRadius(20);
    m_popupwindow->setRadius(6);
    m_popupwindow->setShadowYOffset(2);
    m_popupwindow->setShadowXOffset(0);
    m_popupwindow->setArrowWidth(18);
    m_popupwindow->setArrowHeight(10);
    m_popupwindow->setObjectName("overlaypopup");
    m_popupwindow->setLeftRightRadius(true);

    m_scrollarea->setFrameStyle(QFrame::NoFrame);
    m_scrollarea->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    m_scrollarea->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_scrollarea->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_scrollarea->setBackgroundRole(QPalette::Base);
    m_scrollarea->setWidgetResizable(true);
    m_scrollarea->setAutoFillBackground(true);

    m_centerScroll->setLayout(m_popuplayout);
    m_centerScroll->setAccessibleName(OVERFLOWWIDGET_ACCESS_NAME);
    m_centerScroll->setAttribute(Qt::WA_TranslucentBackground);
    m_centerScroll->setAutoFillBackground(true);
    m_scrollarea->setWidget(m_centerScroll);

    m_popupbtnslayout->addWidget(m_left, 0, Qt::AlignCenter);
    m_popupbtnslayout->addStretch(1);
    m_popupbtnslayout->addWidget(m_right, 0, Qt::AlignCenter);

    m_left->setIcon(QIcon(ARROW_LEFT));
    m_right->setIcon(QIcon(ARROW_RIGHT));
}

void OverflowItem::initSlots() {
    connect(m_left, &QPushButton::clicked, this, [this ]{
        switch (m_popuplayout->direction()) {
            case QBoxLayout::LeftToRight: {
                    int scroll_len = m_scrollarea->horizontalScrollBar()->value() - 50;
                    if (scroll_len <= 10) {
                        m_scrollarea->horizontalScrollBar()->setValue(0);
                    } else {
                        m_scrollarea->horizontalScrollBar()->setValue(scroll_len);
                    }
                }
                break;
            case QBoxLayout::TopToBottom: {
                    int scroll_len = m_scrollarea->verticalScrollBar()->value() - 50;
                    if (scroll_len <= 10) {
                        m_scrollarea->verticalScrollBar()->setValue(0);
                    } else {
                        m_scrollarea->verticalScrollBar()->setValue(scroll_len);
                    }
                }
                break;
            default:
                break;
        }
        setbtnsVisible();
    });
    connect(m_right, &QPushButton::clicked, this, [this ]{
        switch (m_popuplayout->direction()) {
            case QBoxLayout::LeftToRight: {
                    int maxlen = m_scrollarea->horizontalScrollBar()->maximum();
                    int scroll_len = m_scrollarea->horizontalScrollBar()->value() + 50;
                    if (scroll_len > maxlen - 10) {
                        m_scrollarea->horizontalScrollBar()->setValue(maxlen);
                    } else {
                        m_scrollarea->horizontalScrollBar()->setValue(scroll_len);
                    }
                }
                break;
            case QBoxLayout::TopToBottom: {
                    int maxlen = m_scrollarea->verticalScrollBar()->maximum();
                    int scroll_len = m_scrollarea->verticalScrollBar()->value() + 50;
                    if (scroll_len > maxlen - 10) {
                        m_scrollarea->verticalScrollBar()->setValue(maxlen);
                    } else {
                        m_scrollarea->verticalScrollBar()->setValue(scroll_len);
                    }
                }
                break;
            default:
                break;
        }
        setbtnsVisible();
    });
}

void OverflowItem::setbtnsVisible() {
    bool leftshow = true;
    bool rightshow = true;
    switch (m_popuplayout->direction()) {
        case QBoxLayout::LeftToRight:
            if (m_scrollarea->horizontalScrollBar()->value() <= 10) {
                leftshow = false;
            }
            if (m_scrollarea->horizontalScrollBar()->value() >=
                    m_scrollarea->horizontalScrollBar()->maximum() - 10) {
                rightshow = false;
            }
            break;
        case QBoxLayout::TopToBottom:
            if (m_scrollarea->verticalScrollBar()->value() <= 10) {
                leftshow = false;
            }
            if (m_scrollarea->verticalScrollBar()->value() >=
                    m_scrollarea->verticalScrollBar()->maximum() - 10) {
                rightshow = false;
            }
            break;
        default:
            break;
    }
    m_left->setVisible(leftshow);
    m_right->setVisible(rightshow);
}

void OverflowItem::setbtnsShape() {
    switch (m_popupbtnslayout->direction()) {
        case QBoxLayout::LeftToRight:
            m_left->setFixedSize(m_width / 3, m_width * 3 / 4);
            m_right->setFixedSize(m_width / 3, m_width * 3 / 4);
            break;
        case QBoxLayout::TopToBottom:
            m_left->setFixedSize(m_width * 3 / 4 , m_width / 3);
            m_right->setFixedSize(m_width * 3 / 4  , m_width / 3);
            break;
        default:
            break;
    }
}

void OverflowItem::hidePopUpWindow() {
    m_showpopup = false;
    m_popupwindow->hide();
}

void OverflowItem::setPopUpSize(int width, int height) {
    m_scrollarea->setFixedSize(width, height);
    m_width = qMin(width, height);
    setbtnsShape();
    setbtnsVisible();
}

void OverflowItem::addItem(QWidget *item) {
    m_popuplayout->addWidget(item,0, Qt::AlignCenter);
}

QPoint OverflowItem::OverflowIconPosition(const QPixmap &pixmap) const {
    const auto ratio = devicePixelRatioF();
    const QRectF itemRect = rect();
    const QRectF iconRect = pixmap.rect();
    const qreal iconX = itemRect.center().x() - iconRect.center().x() / ratio;
    const qreal iconY = itemRect.center().y() - iconRect.center().y() / ratio;
    return QPoint(iconX, iconY);
}

// FIXME: the size of app sometime cannot be controled
void OverflowItem::paintEvent(QPaintEvent *e) {

    DockItem::paintEvent(e);

    if (!isVisible()) {
        return;
    }
    QPainter painter(this);

    // Start paint image
    QPixmap image(ICON_DEFAULT);
    if (m_popuplayout->count() != 0) {
        image = static_cast<AppItem *>(m_popuplayout->itemAt(0)->widget())->appIcon();
    }
    QPoint realsize = OverflowIconPosition(image);
    painter.drawPixmap(realsize, image);
    // Paint End

    // Add Shadow
    if (isDarkTheme()) {
        painter.setOpacity(0.6);
    } else {
        painter.setOpacity(0.3);
    }
    if (m_hover) {
        painter.setOpacity(0.4);
    }
    if (m_clicked) {
        painter.setOpacity(0.7);
    }
    qreal min = qMin(rect().width(), rect().height());
    QRectF backgroundRect = QRectF(rect().x(), rect().y(), min, min);
    backgroundRect = backgroundRect.marginsRemoved(QMargins(2, 2, 2, 2));
    backgroundRect.moveCenter(rect().center());
    // Shadow end

    // Add More Icon
    QPainterPath path;
    path.addRoundedRect(backgroundRect, 8, 8);
    painter.fillPath(path, QColor(0, 0, 0, 255 * 0.8));

    painter.setOpacity(1);
    QPixmap moreicons(OVERFLOW_MORE);
    QPoint realsize_more = OverflowIconPosition(moreicons);
    moreicons.scaled(realsize.x(), realsize.y());
    painter.drawPixmap(realsize_more, moreicons);
    // Paint "More" End

}

// INFO: public, be set in mainpanelcontrol.cpp, not use it in this cpp
void OverflowItem::setLayoutPosition(Dock::Position position) {
    switch (position) {
        case Top:
        case Bottom:
            m_popuplayout->setDirection(QBoxLayout::LeftToRight);
            m_popupbtnslayout->setDirection(QBoxLayout::LeftToRight);
            m_left->setIcon(QIcon(ARROW_LEFT));
            m_right->setIcon(QIcon(ARROW_RIGHT));
            break;
        case Left:
        case Right:
            m_popuplayout->setDirection(QBoxLayout::TopToBottom);
            m_popupbtnslayout->setDirection(QBoxLayout::TopToBottom);
            m_left->setIcon(QIcon(ARROW_UP));
            m_right->setIcon(QIcon(ARROW_DOWN));
            break;
    }
    setbtnsShape();
    setbtnsVisible();
}

void OverflowItem::mousePressEvent(QMouseEvent *e) {
    m_clicked = true;
    m_showpopup = !m_showpopup;
    if (m_showpopup) {
        m_popupwindow->setLeftRightRadius(DWindowManagerHelper::instance()->hasComposite());
        showPopupWindow(m_scrollarea, true, 12);
    } else {
        m_popupwindow->hide();
    }
    DockItem::mousePressEvent(e);
}

void OverflowItem::mouseMoveEvent(QMouseEvent *e) {
    DockItem::mouseMoveEvent(e);
}

void OverflowItem::mouseReleaseEvent(QMouseEvent *e) {
    m_clicked = false;
    DockItem::mouseReleaseEvent(e);
}

void OverflowItem::enterEvent(QEvent *e) {
    DockItem::enterEvent(e);
}

void OverflowItem::leaveEvent(QEvent *e) {
    m_clicked = false;
    DockItem::leaveEvent(e);
}

void OverflowItem::showPopupWindow(QWidget *const content, const bool model, const int radius) {

    m_popupShown = true;
    m_lastPopupWidget = content;

    if (model)
        emit requestWindowAutoHide(false);

    m_popupwindow->setRadius(radius);
    QWidget *lastContent = m_popupwindow->getContent();
    if (lastContent)
        lastContent->setVisible(false);

    switch (DockPosition) {
        case Top:   m_popupwindow->setArrowDirection(DockPopupWindow::ArrowTop);     break;
        case Bottom: m_popupwindow->setArrowDirection(DockPopupWindow::ArrowBottom);  break;
        case Left:  m_popupwindow->setArrowDirection(DockPopupWindow::ArrowLeft);    break;
        case Right: m_popupwindow->setArrowDirection(DockPopupWindow::ArrowRight);   break;
    }
    m_popupwindow->resize(content->sizeHint());
    m_popupwindow->setContent(content);

    const QPoint p = popupMarkPoint();
    if (!m_popupwindow->isVisible())
        QMetaObject::invokeMethod(m_popupwindow, "show", Qt::QueuedConnection, Q_ARG(QPoint, p), Q_ARG(bool, model));
    else
        m_popupwindow->show(p, model);
}

bool OverflowItem::eventFilter(QObject *watched, QEvent *event) {
    if (watched == m_centerScroll && event->type() == QEvent::Wheel) {
        QWheelEvent *wheelEvent = static_cast<QWheelEvent *>(event);
        const QPoint delta = wheelEvent->angleDelta();
        int scroll_len = qAbs(delta.x()) > qAbs(delta.y()) ? delta.x() : -1 * delta.y();
        if (m_popuplayout->direction() == QBoxLayout::LeftToRight) {
            if (m_scrollarea->horizontalScrollBar()->value() + scroll_len <= 0) {
                m_scrollarea->horizontalScrollBar()->setValue(0);
            } else if (m_scrollarea->horizontalScrollBar()->value() + scroll_len >=
                       m_scrollarea->horizontalScrollBar()->maximum()) {
                m_scrollarea->horizontalScrollBar()->setValue(m_scrollarea->horizontalScrollBar()->maximum());
            } else {
                m_scrollarea->horizontalScrollBar()->setValue(m_scrollarea->horizontalScrollBar()->value() + scroll_len);
            }
        } else {
            if (m_scrollarea->verticalScrollBar()->value() + scroll_len <= 0) {
                m_scrollarea->verticalScrollBar()->setValue(0);
            } else if (m_scrollarea->verticalScrollBar()->value() + scroll_len >=
                       m_scrollarea->verticalScrollBar()->maximum()) {
                m_scrollarea->verticalScrollBar()->setValue(m_scrollarea->verticalScrollBar()->maximum());
            } else {
                m_scrollarea->verticalScrollBar()->setValue(m_scrollarea->verticalScrollBar()->value() + scroll_len);
            }
        }
        setbtnsVisible();
        return true;
    }
    return DockItem::eventFilter(watched,event);
}
