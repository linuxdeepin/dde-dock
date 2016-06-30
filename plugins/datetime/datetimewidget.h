#ifndef DATETIMEWIDGET_H
#define DATETIMEWIDGET_H

#include <QWidget>

class DatetimeWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DatetimeWidget(QWidget *parent = 0);

private:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);

    const QPixmap loadSvg(const QString &fileName, const QSize size);
};

#endif // DATETIMEWIDGET_H
