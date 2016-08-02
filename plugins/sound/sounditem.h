#ifndef SOUNDITEM_H
#define SOUNDITEM_H

#include <QWidget>

class SoundItem : public QWidget
{
    Q_OBJECT

public:
    explicit SoundItem(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);
};

#endif // SOUNDITEM_H
