#ifndef APPBACKGROUND_H
#define APPBACKGROUND_H

#include <QObject>
#include <QLabel>
#include <QStyle>
#include <QPropertyAnimation>
#include <QPainter>
#include <QMouseEvent>
#include <QDebug>
#include "Controller/dockmodedata.h"

class ActiveLabel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(double opacity READ opacity WRITE setOpacity)
public:
    explicit ActiveLabel(QWidget *parent = 0);
    void showActiveWithAnimation();
    void show();

    double opacity() const;
    void setOpacity(double opacity);

protected:
    void paintEvent(QPaintEvent *event);

signals:
    void sizeChange();
    void showAnimationFinish();

private:
    int m_loopCount = 0;
    double m_opacity = 0;
    QString m_iconPath = "";
};

class AppBackground : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(bool isActived READ getIsActived WRITE setIsActived)
    Q_PROPERTY(bool isCurrentOpened READ getIsCurrentOpened WRITE setIsCurrentOpened)
    Q_PROPERTY(bool isHovered READ getIsHovered WRITE setIsHovered)
    Q_PROPERTY(bool isFashionMode READ getIsFashionMode)
public:
    explicit AppBackground(QWidget *parent = 0);

    void resize(int width, int height);

    bool getIsActived();
    void setIsActived(bool value);
    bool getIsCurrentOpened();
    void setIsCurrentOpened(bool value);
    bool getIsHovered();
    void setIsHovered(bool value);
    bool getIsFashionMode() const;

public slots:
    void slotMouseRelease(QMouseEvent *event);

private:
    void initActiveLabel();
    void updateActiveLabelPos();

private:
    bool m_isInit = true;
    bool m_bePress = false;
    bool m_isActived = false;
    bool m_isCurrentOpened = false;
    bool m_isHovered = false;
    bool m_isFashionMode = false;

    ActiveLabel *m_activeLabel = NULL;

    const int ACTIVE_LABEL_WIDTH = 30;
    const int ACTIVE_LABEL_HEIGHT = 10;
};

#endif // APPBACKGROUND_H
