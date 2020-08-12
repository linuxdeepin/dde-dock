#ifndef NETWORKITEM_H
#define NETWORKITEM_H

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

public slots:
    void updateSelf();
    void refreshIcon();
    void wirelessScan();

protected:
    void resizeEvent(QResizeEvent *e) override;
    void paintEvent(QPaintEvent *e) override;
    bool eventFilter(QObject *obj,QEvent *event) override;

private slots:
    void wiredsEnable(bool enable);
    void wirelessEnable(bool enable);
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    void getPluginState();
    void updateMasterControlSwitch();
    void updateView();
    int getStrongestAp();

private:
    Dock::TipsWidget *m_tipsWidget;
    QScrollArea *m_applet;

    QLabel *m_wiredTitle;
    DSwitchButton *m_switchWiredBtn;
    QVBoxLayout *m_wiredLayout;
    QWidget *m_wiredControlPanel;
    bool m_switchWiredBtnState;

    QLabel *m_wirelessTitle;
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
    QTimer *m_timer;
    QTimer *m_switchWireTimer;
};

#endif // NETWORKITEM_H
