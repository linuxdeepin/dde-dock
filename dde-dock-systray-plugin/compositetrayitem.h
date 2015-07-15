#ifndef COMPOSITETRAYITEM_H
#define COMPOSITETRAYITEM_H

#include <QFrame>

class CompositeTrayItem : public QFrame
{
    Q_OBJECT
public:
    explicit CompositeTrayItem(QWidget *parent = 0);
    virtual ~CompositeTrayItem();

    void addWidget(QWidget * widget);
    void removeWidget(QWidget * widget);

private:
    uint m_columnCount;

    void setBackground();
};

#endif // COMPOSITETRAYITEM_H
