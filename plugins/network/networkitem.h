#ifndef NETWORKITEM_H
#define NETWORKITEM_H

#include <DGuiApplicationHelper>
#include <DSwitchButton>
#include <dloadingindicator.h>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QTimer>
#include <QScrollBar>

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
        Unknown              = 0,
        // A 无线 B 有线
        Nocable             = 1,
        Bdisabled           = 2,
        Bdisconnected       = 8,
        Bconnected          = 16,
        Bconnecting         = 32,
        BconnectNoInternet  = 512,
        Bfailed             = 1024,

        //有线无线同时处于连接状态
        Connected,
        //无线网状态
        Adisabled,
        Aconnected,
        Adisconnected,
        Aconnecting,
        AconnectNoInternet
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
    /**
     * @def onConnected
     * @brief 连接中wifi动图
     */
    void onConnecting();

protected:
    void resizeEvent(QResizeEvent *e) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;
    bool eventFilter(QObject *obj,QEvent *event) Q_DECL_OVERRIDE;

private slots:
    void wiredsEnable(bool enable);
    void wirelessEnable(bool enable);
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    void getPluginState();
    void updateMasterControlSwitch();
    void updateView();
    int getStrongestAp();
    /**
     * @def wirelessItemsRequireScan
     * @brief 刷新wifi数据
     **/
    void wirelessItemsRequireScan();

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
    /**
     * @brief m_isWireless
     * @remark 这个参数是用来判断是否为wifi需要使用正在连接的状态
     */
    bool m_isWireless;

    QMap<QString, WiredItem *> m_wiredItems;
    QMap<QString, WirelessItem *> m_wirelessItems;
    QMap<QString, WirelessItem *> m_connectedWirelessDevice;
    QMap<QString, WiredItem *> m_connectedWiredDevice;

    QPixmap m_iconPixmap;
    PluginState m_pluginState;
    /**
     * @brief m_timer
     * @brief 做刷新操作的QTimer
     */
    QTimer *m_connectingTimer;
};

#endif // NETWORKITEM_H
