#ifndef APPSNAPSHOT_H
#define APPSNAPSHOT_H

#include <QWidget>
#include <QDebug>
#include <QTimer>

class AppSnapshot : public QWidget
{
    Q_OBJECT

public:
    explicit AppSnapshot(const WId wid, QWidget *parent = 0);

signals:
    void entered(const WId wid) const;
    void clicked(const WId wid) const;

private slots:
    void fetchSnapshot();

private:
    void enterEvent(QEvent *e);
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    const WId m_wid;

    QImage m_snapshot;

    QTimer *m_fetchSnapshotTimer;
};

#endif // APPSNAPSHOT_H
