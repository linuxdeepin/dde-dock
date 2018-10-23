#ifndef SOUNDTRAYWIDGET_H
#define SOUNDTRAYWIDGET_H

#include "abstracttraywidget.h"
#include "soundapplet.h"
#include "dbus/dbussink.h"

#include <QWidget>

class TipsWidget;
class SoundTrayWidget : public AbstractTrayWidget
{
public:
    SoundTrayWidget(AbstractTrayWidget *parent = nullptr);

    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    void sendClick(uint8_t, int, int) Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *tipsWidget();
    QWidget *popupApplet();

    const QString contextMenu() const;
    void invokeMenuItem(const QString menuId, const bool checked);

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *e) Q_DECL_OVERRIDE;
    void wheelEvent(QWheelEvent *e) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private Q_SLOTS:
    void refreshTips(const bool force = false);
    void sinkChanged(DBusSink *sink);

private:
    TipsWidget *m_tipsLabel;
    SoundApplet *m_applet;
    DBusSink *m_sinkInter;
    QPixmap m_iconPixmap;
};

#endif // SOUNDTRAYWIDGET_H
