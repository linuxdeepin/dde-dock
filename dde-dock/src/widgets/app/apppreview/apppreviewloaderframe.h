#ifndef APPPREVIEWLOADERFRAME_H
#define APPPREVIEWLOADERFRAME_H

#include <QWidget>
#include <QFrame>

class QLabel;
class QPushButton;
class QHBoxLayout;
class AppPreviewLoader;

class PopupFrame : public QWidget
{
    Q_OBJECT
public:
    explicit PopupFrame(QWidget *parent = 0);

signals:
    void mousePress();
    void mouseLeave();

protected:
    void mousePressEvent(QMouseEvent *);
    void leaveEvent(QEvent *);
};

class AppPreviewLoaderFrame : public QFrame
{
    Q_OBJECT
public:
    explicit AppPreviewLoaderFrame(const QString &title, int xid, QWidget *parent=0);
    ~AppPreviewLoaderFrame();
    void shrink(const QSize &size, bool miniStyle);

signals:
    void requestPreviewClose(int xid);
    void requestPreviewActive(int xid);

protected:
    void enterEvent(QEvent *);

private:
    void initPopupWidget();
    void initTitle(const QString &t);
    void initPreviewLoader(int xid);
    void initCloseButton();
    void updatePopWidgetGeometry();
    void updatePreviewLoaderGeometry();
    void updateCloseButtonGeometry();
    void updateTitleGeometry();
    void updateWidgetsGeometry();

private:
    AppPreviewLoader * m_previewLoader;
    QPushButton *m_closeButton;
    QHBoxLayout *m_layout;
    PopupFrame *m_popupWidget; //for show popup preview
    QWidget *m_parent;  //for reparent to make preview show like popup window
    QLabel *m_titleLabel;
    bool m_inMiniStyle;
    bool m_canShowTitle;
    int m_xid;
};

#endif // APPPREVIEWLOADERFRAME_H
