#ifndef TIPSWIDGET_H
#define TIPSWIDGET_H

#include <QWidget>
#include <QBoxLayout>

class TrayWidget;
class TipsWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TipsWidget(QWidget *parent = 0);

    void clear();
    void addWidgets(QList<TrayWidget *> widgets);

private:
    QBoxLayout *m_mainLayout;
};

#endif // TIPSWIDGET_H
