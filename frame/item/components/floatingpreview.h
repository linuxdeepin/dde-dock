#ifndef FLOATINGPREVIEW_H
#define FLOATINGPREVIEW_H

#include <QWidget>
#include <QPointer>

#include <dimagebutton.h>

DWIDGET_USE_NAMESPACE

class AppSnapshot;
class FloatingPreview : public QWidget
{
    Q_OBJECT

public:
    explicit FloatingPreview(QWidget *parent = 0);

    WId trackedWid() const;

signals:
    void requestMove(const QPoint &p) const;

public slots:
    void trackWindow(AppSnapshot * const snap);

private:
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    bool eventFilter(QObject *watched, QEvent *event);

private slots:
    void onCloseBtnClicked();

private:
    QPointer<AppSnapshot> m_tracked;

    DImageButton *m_closeBtn;
};

#endif // FLOATINGPREVIEW_H
