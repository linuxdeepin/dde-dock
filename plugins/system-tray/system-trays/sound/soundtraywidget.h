#ifndef SOUNDTRAYWIDGET_H
#define SOUNDTRAYWIDGET_H

#include "../abstractsystemtraywidget.h"
#include "soundapplet.h"
#include "dbus/dbussink.h"

#include <QWidget>

class TipsWidget;
class SoundTrayWidget : public AbstractSystemTrayWidget
{
public:
    SoundTrayWidget(QWidget *parent = nullptr);

    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *trayTipsWidget() Q_DECL_OVERRIDE;
    QWidget *trayPopupApplet() Q_DECL_OVERRIDE;

    const QString contextMenu() const Q_DECL_OVERRIDE;
    void invokedMenuItem(const QString &menuId, const bool checked) Q_DECL_OVERRIDE;

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
