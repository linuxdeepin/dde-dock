#ifndef REFRESHBUTTON_H
#define REFRESHBUTTON_H

#include <QLabel>

class RefreshButton : public QLabel
{
    Q_OBJECT
public:
    explicit RefreshButton(QWidget *parent = nullptr);

signals:
    void clicked();

protected:
    void enterEvent(QEvent *event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void mouseReleaseEvent(QMouseEvent *event) Q_DECL_OVERRIDE;

};

#endif // REFRESHBUTTON_H
