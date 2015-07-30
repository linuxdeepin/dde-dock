#include "screenmask.h"

ScreenMask::ScreenMask(QWidget *parent) : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),rec.height());
    this->setWindowFlags(Qt::FramelessWindowHint | Qt::Dialog);
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
    AppItem *sourceItem = NULL;
    sourceItem = dynamic_cast<AppItem *>(event->source());
    if (sourceItem)
    {
        DBusDockedAppManager dda;
        if (dda.IsDocked(sourceItem->itemData().id))
            dda.RequestUndock(sourceItem->itemData().id);

        qWarning() << "Item drop here:" << event->pos() << event->mimeData()->hasImage();
        QImage image = qvariant_cast<QImage>(event->mimeData()->imageData());
        if (!image.isNull())
        {
            TransformLabel * imgLabel = new TransformLabel();
            imgLabel->setAttribute(Qt::WA_TranslucentBackground);
            imgLabel->setWindowFlags(Qt::ToolTip);
            imgLabel->setPixmap(QPixmap::fromImage(image).scaled(ICON_SIZE, ICON_SIZE));
            imgLabel->setFixedSize(ICON_SIZE, ICON_SIZE);
            imgLabel->move(event->pos());
            imgLabel->show();

            QPropertyAnimation *scaleAnimation = new QPropertyAnimation(imgLabel, "sValue");
            scaleAnimation->setDuration(1000);
            scaleAnimation->setStartValue(1);
            scaleAnimation->setEndValue(0.3);

            QPropertyAnimation *rotationAnimation = new QPropertyAnimation(imgLabel, "rValue");
            rotationAnimation->setDuration(1000);
            rotationAnimation->setStartValue(0);
            rotationAnimation->setEndValue(360);

            QParallelAnimationGroup * group = new QParallelAnimationGroup();
            group->addAnimation(scaleAnimation);
//            group->addAnimation(rotationAnimation);

            group->start();
            connect(group, &QPropertyAnimation::finished, [=]{
                imgLabel->deleteLater();
                scaleAnimation->deleteLater();
                rotationAnimation->deleteLater();
                group->deleteLater();
            });

            emit itemDropped(event->pos());
        }
        else
        {
            qWarning() << "Image is NULL!";
        }
    }

    this->close();

}

//Item drag to some place which not accept drop event will cause item missing
//Need layout restore state actively
void ScreenMask::enterEvent(QEvent *)
{
    emit itemMissing();

    this->close();
}
