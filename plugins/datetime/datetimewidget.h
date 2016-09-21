#ifndef DATETIMEWIDGET_H
#define DATETIMEWIDGET_H

#include <QWidget>

class DatetimeWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DatetimeWidget(QWidget *parent = 0);

signals:
    void requestContextMenu() const;

private:
    QSize sizeHint() const;
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

    const QPixmap loadSvg(const QString &fileName, const QSize size);

private:
    QPixmap m_cachedIcon;
    QString m_cachedTime;
};

#endif // DATETIMEWIDGET_H
