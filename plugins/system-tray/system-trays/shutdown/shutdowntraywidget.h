#ifndef SHUTDOWNTRAYWIDGET_H
#define SHUTDOWNTRAYWIDGET_H

#include "abstracttraywidget.h"

#include <QWidget>
#include <QTimer>

class ShutdownTrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    explicit ShutdownTrayWidget(AbstractTrayWidget *parent = nullptr);

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    void sendClick(uint8_t mouseButton, int x, int y) Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private:
    QPixmap m_pixmap;
};

#endif // SHUTDOWNTRAYWIDGET_H
