#ifndef PREVIEWWIDGET_H
#define PREVIEWWIDGET_H

#include <QWidget>
#include <QDebug>
#include <QDragEnterEvent>
#include <QTimer>

#include <dimagebutton.h>

DWIDGET_USE_NAMESPACE

class PreviewWidget : public QWidget
{
    Q_OBJECT
public:
    explicit PreviewWidget(const WId wid, QWidget *parent = 0);

    void setTitle(const QString &title);

signals:
    void requestActivateWindow(const WId wid) const;
    void requestPreviewWindow(const WId wid) const;
    void requestCancelPreview() const;
    void requestHidePreview() const;

private slots:
    void refreshImage();
    void closeWindow();
    void showPreview();

private:
    void paintEvent(QPaintEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);
    void dropEvent(QDropEvent *e);

private:
    const WId m_wid;
    QImage m_image;
    QString m_title;

    DImageButton *m_closeButton;

    QTimer *m_droppedDelay;
    QTimer *m_mouseEnterTimer;

    bool m_hovered;
};

#endif // PREVIEWWIDGET_H
