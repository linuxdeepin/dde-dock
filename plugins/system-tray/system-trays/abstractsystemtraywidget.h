#ifndef ABSTRACTSYSTEMTRAYWIDGET_H
#define ABSTRACTSYSTEMTRAYWIDGET_H

#include "constants.h"
#include "abstracttraywidget.h"
#include "util/dockpopupwindow.h"
#include "dbus/dbusmenumanager.h"

class AbstractSystemTrayWidget : public AbstractTrayWidget
{
public:
    AbstractSystemTrayWidget(QWidget *parent = nullptr);
    virtual ~AbstractSystemTrayWidget();

    void sendClick(uint8_t mouseButton, int x, int y) Q_DECL_OVERRIDE;

    virtual inline QWidget *trayTipsWidget() { return nullptr; }
    virtual inline QWidget *trayPopupApplet() { return nullptr; }
    virtual inline const QString trayClickCommand() { return QString(); }
    virtual inline const QString contextMenu() const {return QString(); }
    virtual inline void invokedMenuItem(const QString &itemId, const bool checked) { Q_UNUSED(itemId); Q_UNUSED(checked); }

protected:
    bool event(QEvent *event) Q_DECL_OVERRIDE;
    void enterEvent(QEvent *event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;

protected:
    const QPoint popupMarkPoint() const;
    const QPoint topleftPoint() const;

    void hidePopup();
    void hideNonModel();
    void popupWindowAccept();
    void showPopupApplet(QWidget * const applet);

    virtual void showPopupWindow(QWidget * const content, const bool model = false);
    virtual void showHoverTips();

protected Q_SLOTS:
    void showContextMenu();

private:
    void updatePopupPosition();

private:
    bool m_popupShown;

    QPointer<QWidget> m_lastPopupWidget;

    QTimer *m_popupTipsDelayTimer;
    QTimer *m_popupAdjustDelayTimer;

    DBusMenuManager *m_menuManagerInter;

    static Dock::Position DockPosition;
    static QPointer<DockPopupWindow> PopupWindow;
};

#endif // ABSTRACTSYSTEMTRAYWIDGET_H
