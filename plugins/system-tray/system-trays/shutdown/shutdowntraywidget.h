#ifndef SHUTDOWNTRAYWIDGET_H
#define SHUTDOWNTRAYWIDGET_H

#include "../abstractsystemtraywidget.h"
#include "../widgets/tipswidget.h"

#include <QWidget>
#include <QTimer>

class ShutdownTrayWidget : public AbstractSystemTrayWidget
{
    Q_OBJECT

public:
    explicit ShutdownTrayWidget(QWidget *parent = nullptr);

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *trayTipsWidget() Q_DECL_OVERRIDE;
    const QString trayClickCommand() Q_DECL_OVERRIDE;

    const QString contextMenu() const Q_DECL_OVERRIDE;
    void invokedMenuItem(const QString &menuId, const bool checked) Q_DECL_OVERRIDE;

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private:
    TipsWidget *m_tipsLabel;

    QPixmap m_pixmap;
};

#endif // SHUTDOWNTRAYWIDGET_H
