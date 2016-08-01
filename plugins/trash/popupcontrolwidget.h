#ifndef POPUPCONTROLWIDGET_H
#define POPUPCONTROLWIDGET_H

#include <QWidget>
#include <QFileSystemWatcher>

#include <dlinkbutton.h>

class PopupControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PopupControlWidget(QWidget *parent = 0);

    bool empty() const;

signals:
    void emptyChanged(const bool empty) const;

private slots:
    void openTrashFloder();
    void clearTrashFloder();
    void trashStatusChanged();

private:
    bool m_empty;

    Dtk::Widget::DLinkButton *m_openBtn;
    Dtk::Widget::DLinkButton *m_clearBtn;

    QFileSystemWatcher *m_fsWatcher;
};

#endif // POPUPCONTROLWIDGET_H
