#ifndef CONTAINERWIDGET_H
#define CONTAINERWIDGET_H

#include <QWidget>
#include <QHBoxLayout>

class ContainerWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ContainerWidget(QWidget *parent = 0);

    void addWidget(QWidget * const w);
    void removeWidget(QWidget * const w);
    int itemCount() const;
    const QList<QWidget *> itemList() const;

    bool allowDragEnter(QDragEnterEvent *e);

protected:
    void dragEnterEvent(QDragEnterEvent *e);

private:
    QHBoxLayout *m_centralLayout;

    QList<QWidget *> m_itemList;
};

#endif // CONTAINERWIDGET_H
