#include "abstracttraywidget.h"

#include <xcb/xproto.h>
#include <QMouseEvent>

AbstractTrayWidget::AbstractTrayWidget(QWidget *parent, Qt::WindowFlags f):
    QWidget(parent, f)
{

}

AbstractTrayWidget::~AbstractTrayWidget()
{

}

void AbstractTrayWidget::mouseReleaseEvent(QMouseEvent *e)
{
    const QPoint point(e->pos() - rect().center());
    if (point.manhattanLength() > 24)
        return;

    e->accept();

    QPoint globalPos = QCursor::pos();
    uint8_t buttonIndex = XCB_BUTTON_INDEX_1;

    switch (e->button()) {
    case Qt:: MiddleButton:
        buttonIndex = XCB_BUTTON_INDEX_2;
        break;
    case Qt::RightButton:
        buttonIndex = XCB_BUTTON_INDEX_3;
        break;
    default:
        break;
    }

    sendClick(buttonIndex, globalPos.x(), globalPos.y());
}
