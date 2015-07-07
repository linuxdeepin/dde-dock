#include "screenmask.h"

ScreenMask::ScreenMask(QWidget *parent) : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),rec.height());
    this->setWindowFlags(Qt::FramelessWindowHint);
    this->setWindowOpacity(0);

    setAcceptDrops(true);

    this->show();
}

void ScreenMask::dragEnterEvent(QDragEnterEvent *event)
{
    event->setDropAction(Qt::MoveAction);
    event->accept();

    emit itemEntered();
}

void ScreenMask::dragLeaveEvent(QDragLeaveEvent *event)
{
    emit itemExited();
}

void ScreenMask::dropEvent(QDropEvent *event)
{
    qWarning() << "Item drop here:" << event->pos() << event->mimeData()->hasImage();
    QImage image = qvariant_cast<QImage>(event->mimeData()->imageData());
    if (!image.isNull())
    {
        TransformLabel * imgLabel = new TransformLabel();
        imgLabel->setAttribute(Qt::WA_TranslucentBackground);
        imgLabel->setWindowFlags(Qt::ToolTip);
        imgLabel->setPixmap(QPixmap::fromImage(image));
        imgLabel->move(event->pos());
        imgLabel->show();

        //TODO add animation here
        QTimer::singleShot(1000,imgLabel,SLOT(deleteLater()));

        emit itemDropped(event->pos());
    }
    else
    {
        qWarning() << "Image is NULL!";
    }

    this->close();

}
