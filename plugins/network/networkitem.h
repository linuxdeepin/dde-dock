#ifndef NETWORKITEM_H
#define NETWORKITEM_H

#include <DSwitchButton>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QTimer>

DWIDGET_USE_NAMESPACE

class PluginState;
class TipsWidget;
class WiredItem;
class WirelessItem;
class HorizontalSeperator;
class NetworkItem : public QWidget
{
    Q_OBJECT

    enum PluginState
    {
        Unknow              = 0,
        // A 无线 B 有线
        Disabled,
        Connected,
        Disconnected,
        Connecting,
//        Failed,
        ConnectNoInternet,
//        Aenabled,
//        Benabled,
        Adisabled,
        Bdisabled,
        Aconnected,
        Bconnected,
        Adisconnected,
        Bdisconnected,
        Aconnecting,
        Bconnecting,
        AconnectNoInternet,
        BconnectNoInternet,
//        Afailed,
        Bfailed,
        Nocable
    };
public:
    explicit NetworkItem(QWidget *parent = nullptr);

    QWidget *itemApplet();
    QWidget *itemTips();

    void updateDeviceItems(QMap<QString, WiredItem *> &wiredItems, QMap<QString, WirelessItem*> &wirelessItems);

    const QString contextMenu() const;
    void invokeMenuItem(const QString &menuId, const bool checked);
    void refreshTips();

public slots:
    void deviceDel();
    void updateSelf();
    void refreshIcon();

protected:
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);

private slots:
    void wiredsEnable(bool enable);
    void wirelessEnable(bool enable);

private:
    void getPluginState();
    void updateMasterControlSwitch();
    void updateView();
    int getStrongestAp();

private:
    TipsWidget *m_tipsWidget;
    QScrollArea *m_applet;

    QLabel *m_wiredTitle;
    DSwitchButton *m_switchWiredBtn;
    QVBoxLayout *m_wiredLayout;
    QWidget *m_wiredControlPanel;
    bool m_switchWiredBtnState;

    HorizontalSeperator *m_line;

    QLabel *m_wirelessTitle;
    DSwitchButton *m_switchWirelessBtn;
    QVBoxLayout *m_wirelessLayout;
    QWidget *m_wirelessControlPanel;
    bool m_switchWirelessBtnState;

    int m_sixteenDeviceHeight;
    bool m_switchWire;

    QMap<QString, WiredItem *> m_wiredItems;
    QMap<QString, WirelessItem *> m_wirelessItems;
    QMap<QString, WirelessItem *> m_connectedWirelessDevice;
    QMap<QString, WiredItem *> m_connectedWiredDevice;

    QPixmap m_iconPixmap;
    PluginState m_pluginState;
    QTimer *m_timer;
    QTimer *m_switchWireTimer;
};

#endif // NETWORKITEM_H
