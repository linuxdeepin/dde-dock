#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include <QScrollArea>
#include <QVBoxLayout>

class SoundApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit SoundApplet(QWidget *parent = 0);

private:
    QWidget *m_centeralWidget;
    QVBoxLayout *m_centeralLayout;
};

#endif // SOUNDAPPLET_H
