#ifndef POPUPCONTROLWIDGET_H
#define POPUPCONTROLWIDGET_H

#include <QWidget>

#include <dlinkbutton.h>

class PopupControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PopupControlWidget(QWidget *parent = 0);

private slots:
    void openTrashFloder();

private:
    Dtk::Widget::DLinkButton *m_openBtn;
    Dtk::Widget::DLinkButton *m_clearBtn;
};

#endif // POPUPCONTROLWIDGET_H
