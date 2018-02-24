#ifndef INFORMATIONWIDGET_H
#define INFORMATIONWIDGET_H

#include <QWidget>
#include <QLabel>

class InformationWidget : public QWidget
{
    Q_OBJECT

public:
    explicit InformationWidget(QWidget *parent = nullptr);

private slots:
    void refreshInfo();

private:
    QLabel *m_infoLabel;
};

#endif // INFORMATIONWIDGET_H
