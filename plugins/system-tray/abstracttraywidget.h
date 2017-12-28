#pragma once

#include <QWidget>

class QDBusMessage;
class AbstractTrayWidget: public QWidget
{
    Q_OBJECT
public:
    explicit AbstractTrayWidget(QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    virtual ~AbstractTrayWidget();

    virtual void setActive(const bool active) = 0;
    virtual void updateIcon() = 0;
    virtual void sendClick(uint8_t, int, int) = 0;
    virtual const QImage trayImage() = 0;

Q_SIGNALS:
    void iconChanged();

protected:
    void mouseReleaseEvent(QMouseEvent *e) Q_DECL_OVERRIDE;
};

