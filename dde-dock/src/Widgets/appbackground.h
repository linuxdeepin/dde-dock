#ifndef APPBACKGROUND_H
#define APPBACKGROUND_H

#include <QObject>
#include <QLabel>
#include <QStyle>
#include <QDebug>

class AppBackground : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(bool isActived READ getIsActived WRITE setIsActived)
    Q_PROPERTY(bool isCurrentOpened READ getIsCurrentOpened WRITE setIsCurrentOpened)
    Q_PROPERTY(bool isHovered READ getIsHovered WRITE setIsHovered)
public:
    explicit AppBackground(QWidget *parent = 0);

    bool getIsActived();
    void setIsActived(bool value);
    bool getIsCurrentOpened();
    void setIsCurrentOpened(bool value);
    bool getIsHovered();
    void setIsHovered(bool value);

signals:

public slots:

private:
    bool m_isActived = false;
    bool m_isCurrentOpened = false;
    bool m_isHovered = false;
};

#endif // APPBACKGROUND_H
