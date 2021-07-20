#ifndef NETWORKITEM_H
#define NETWORKITEM_H

#include "com_deepin_daemon_network.h"

#include <DGuiApplicationHelper>
#include <DSwitchButton>
#include <dloadingindicator.h>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QTimer>

DGUI_USE_NAMESPACE
DWIDGET_USE_NAMESPACE

class PluginState;
namespace Dock {
class TipsWidget;
}
class WiredItem;
class WirelessItem;
class HorizontalSeperator;

using DbusNetwork = com::deepin::daemon::Network;

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
        //有线无线都失败
        Failed,
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
        Afailed,
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
    bool isShowControlCenter();

    const QStringList currentIpList();

public slots:
    void updateSelf();
    void refreshIcon();
    void wirelessScan();

protected:
    void resizeEvent(QResizeEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    bool eventFilter(QObject *obj,QEvent *event) override;
    QString getStrengthStateString(int strength = 0);

private slots:
    void wiredsEnable(bool enable);
    void wirelessEnable(bool enable);
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);
    void ipConfllict(const QString &in0, const QString &in1);

private:
    void getPluginState();
    void updateMasterControlSwitch();
    void updateView();
    int getStrongestAp();

private:
    Dock::TipsWidget *m_tipsWidget;
    QScrollArea *m_applet;

    DSwitchButton *m_switchWiredBtn;
    QVBoxLayout *m_wiredLayout;
    QWidget *m_wiredControlPanel;
    bool m_switchWiredBtnState;

    DLoadingIndicator *m_loadingIndicator;
    DSwitchButton *m_switchWirelessBtn;
    QVBoxLayout *m_wirelessLayout;
    QWidget *m_wirelessControlPanel;
    bool m_switchWirelessBtnState;

    bool m_switchWire;
    //判断定时的时间是否到,否则不重置计时器
    bool m_timeOut;

    QMap<QString, WiredItem *> m_wiredItems;
    QMap<QString, WirelessItem *> m_wirelessItems;
    QMap<QString, WirelessItem *> m_connectedWirelessDevice;
    QMap<QString, WiredItem *> m_connectedWiredDevice;

    QPixmap m_iconPixmap;
    PluginState m_pluginState;
    QTimer *refreshIconTimer;
    QTimer *m_switchWireTimer;
    QTimer *m_wirelessScanTimer;
    int m_wirelessScanInterval;

    HorizontalSeperator *m_firstSeparator;
    HorizontalSeperator *m_secondSeparator;
    HorizontalSeperator *m_thirdSeparator;

    DbusNetwork *m_networkInter;
    QTimer *m_detectTimer;
    QTime m_timeElapse;
    QString m_ipAddr;
    bool m_ipConflict;
};

#endif // NETWORKITEM_H
