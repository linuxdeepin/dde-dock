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
    const QStringList getActiveWiredList();
    const QStringList getActiveWirelessList();

public slots:
    void updateSelf();
    void refreshIcon();
    void wirelessScan();

protected:
    void resizeEvent(QResizeEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    bool eventFilter(QObject *obj,QEvent *event) override;
    QString getStrengthStateString(int strength = 0);

signals:
    void sendIpConflictDect(int index);

private slots:
    void wiredsEnable(bool enable);
    void wirelessEnable(bool enable);
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);
    void ipConflict(const QString &ip, const QString &mac);
    void onSendIpConflictDect(int index = 0);
    void onDetectConflict();

private:
    void getPluginState();
    void updateMasterControlSwitch();
    void updateView();
    int getStrongestAp();
    bool isExistAvailableNetwork();

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
    QStringList m_disconflictList;         // 解除冲突数据列表
    QMap<QString, QString> m_conflictMap;  // 缓存有线和无线冲突的ip列表
    QTimer *m_detectConflictTimer;         // 定时器自检,当其他主机主动解除ip冲突，我方需要更新网络状态
    bool m_ipConflict;                     // ip冲突的标识
    bool m_ipConflictChecking;             // 标记是否正在检测中
};

#endif // NETWORKITEM_H
