#ifndef APPSNAPSHOT_H
#define APPSNAPSHOT_H

#include <QWidget>
#include <QDebug>
#include <QTimer>
#include <QLabel>

#include <dimagebutton.h>
#include <DWindowManagerHelper>

DWIDGET_USE_NAMESPACE

class AppSnapshot : public QWidget
{
    Q_OBJECT

public:
    explicit AppSnapshot(const WId wid, QWidget *parent = 0);

    WId wid() const { return m_wid; }
    const QImage snapshot() const { return m_snapshot; }
    const QString title() const { return m_title->text(); }

signals:
    void entered(const WId wid) const;
    void clicked(const WId wid) const;
    void requestCheckWindow() const;

public slots:
    void fetchSnapshot();
    void closeWindow() const;
    void compositeChanged() const;
    void setWindowTitle(const QString &title);

private:
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    const WId m_wid;

    QImage m_snapshot;
    QLabel *m_title;
    DImageButton *m_closeBtn;

    DWindowManagerHelper *m_wmHelper;
};

#endif // APPSNAPSHOT_H
