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

private slots:
    void fetchSnapshot();

private:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);

private:
    const WId m_wid;

    QImage m_snapshot;

    QTimer *m_fetchSnapshotTimer;
};

#endif // APPSNAPSHOT_H
