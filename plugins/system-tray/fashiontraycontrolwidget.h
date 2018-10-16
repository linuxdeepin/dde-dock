#ifndef FASHIONTRAYCONTROLWIDGET_H
#define FASHIONTRAYCONTROLWIDGET_H

#include "constants.h"

#include <QLabel>
#include <QSettings>

class FashionTrayControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit FashionTrayControlWidget(Dock::Position position, QWidget *parent = nullptr);

    void setDockPostion(Dock::Position pos);

    bool expanded() const;
    void setExpanded(const bool &expanded);

Q_SIGNALS:
    void expandChanged(const bool expanded);

protected:
    void paintEvent(QPaintEvent *event) Q_DECL_OVERRIDE;
    void mouseReleaseEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void enterEvent(QEvent *event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *event) Q_DECL_OVERRIDE;

private:
    QSettings *m_settings;

    Dock::Position m_dockPosition;
    bool m_expanded;
    bool m_hover;
    bool m_pressed;
};

#endif // FASHIONTRAYCONTROLWIDGET_H
