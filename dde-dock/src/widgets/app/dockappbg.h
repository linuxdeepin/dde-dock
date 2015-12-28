#ifndef DOCKAPPBG_H
#define DOCKAPPBG_H

#include <QDebug>
#include <QLabel>
#include <QStyle>
#include <QPainter>
#include <QMouseEvent>
#include <QPropertyAnimation>

#include "controller/dockmodedata.h"

class BGActiveIndicator : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(double opacity READ opacity WRITE setOpacity)
    Q_PROPERTY(QString openIndicatorIcon READ openIndicatorIcon WRITE setOpenIndicatorIcon)
    Q_PROPERTY(QString openingIndicatorIcon READ openingIndicatorIcon WRITE setOpeningIndicatorIcon)

public:
    explicit BGActiveIndicator(QWidget *parent = 0);
    void showActivatingAnimation();
    void show();

    double opacity() const;
    void setOpacity(double opacity);

    QString openIndicatorIcon() const;
    void setOpenIndicatorIcon(const QString &openIndicatorIcon);

    QString openingIndicatorIcon() const;
    void setOpeningIndicatorIcon(const QString &openingIndicatorIcon);

signals:
    void sizeChange();
    void showAnimationFinish();

protected:
    void paintEvent(QPaintEvent *event);

private:
    int m_loopCount = 0;
    double m_opacity = 0;
    QString m_openIndicatorIcon = "";
    QString m_openingIndicatorIcon = "";
    QString m_iconPath = "";
};

class DockAppBG : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(bool isActived READ isActived WRITE setIsActived)
    Q_PROPERTY(bool isCurrentOpened READ isCurrentOpened WRITE setIsCurrentOpened)
    Q_PROPERTY(bool isHovered READ isHovered WRITE setIsHovered)
    Q_PROPERTY(bool isFashionMode READ isFashionMode)
public:
    explicit DockAppBG(QWidget *parent = 0);
    void showActivatingAnimation();

    bool isActived();
    void setIsActived(bool value);
    bool isCurrentOpened();
    void setIsCurrentOpened(bool value);
    bool isHovered();
    void setIsHovered(bool value);
    bool isFashionMode() const;

protected:
    void resizeEvent(QResizeEvent *);

private:
    void initActiveLabel();
    void updateActiveLabelPos();
    void onDockModeChanged();

private:
    bool m_bePress = false;
    bool m_isActived = false;
    bool m_isCurrentOpened = false;
    bool m_isHovered = false;
    bool m_isFashionMode = false;

    BGActiveIndicator *m_activeLabel;

    const int ACTIVE_LABEL_WIDTH = 30;
    const int ACTIVE_LABEL_HEIGHT = 10;
};

#endif // DOCKAPPBG_H
