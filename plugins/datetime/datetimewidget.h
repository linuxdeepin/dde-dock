#ifndef DATETIMEWIDGET_H
#define DATETIMEWIDGET_H

#include <QWidget>
#include <QSettings>

class DatetimeWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DatetimeWidget(QWidget *parent = 0);

    bool is24HourFormat() const { return m_24HourFormat; }

signals:
    void requestContextMenu() const;

public slots:
    void toggleHourFormat();

private:
    QSize sizeHint() const;
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

    const QPixmap loadSvg(const QString &fileName, const QSize size);

private:
    QPixmap m_cachedIcon;
    QString m_cachedTime;
    QSettings m_settings;
    bool m_24HourFormat;
};

#endif // DATETIMEWIDGET_H
