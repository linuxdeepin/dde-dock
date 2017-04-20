#ifndef PREVIEWWIDGET_H
#define PREVIEWWIDGET_H

#include <QWidget>
#include <QDebug>
#include <QPushButton>

class PreviewWidget : public QWidget
{
    Q_OBJECT
public:
    explicit PreviewWidget(const WId wid, QWidget *parent = 0);

    void setTitle(const QString &title);

signals:
    void requestActivateWindow(const WId wid) const;

private slots:
    void refershImage();
    void closeWindow();

private:
    void paintEvent(QPaintEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

private:
    const WId m_wid;
    QImage m_image;
    QString m_title;

    QPushButton *m_closeButton;

    bool m_hovered;
};

#endif // PREVIEWWIDGET_H
