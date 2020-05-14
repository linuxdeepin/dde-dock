#ifndef PLUGINCASE_H
#define PLUGINCASE_H

#include <QObject>

class PluginCase;
class PluginState : public QObject
{
    Q_OBJECT
public:
    void setPluginCase(PluginCase *pcase) {m_pluginCase = pcase;}
    virtual void openWire() {}
    virtual void closeWire() {}
    virtual void wireConnect() {}
    virtual void wireLinkCheck() {}
    virtual void wireIpCheck() {}
    virtual void wireInternetCheck() {}
    virtual void openWireless() {}
    virtual void closeWireless() {}
    virtual void wirelessConnect() {}
protected:
    explicit PluginState(QObject *parent = nullptr) : QObject(parent) {}

protected:
    PluginCase *m_pluginCase;
};

class WireOpenState : public PluginState
{
    Q_OBJECT
public:
    explicit WireOpenState(QObject *parent = nullptr);
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireCloseState : public PluginState
{
    Q_OBJECT
public:
    explicit WireCloseState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireOpenWirelessCloseState : public PluginState
{
    Q_OBJECT
public:
    explicit WireOpenWirelessCloseState(QObject *parent = nullptr);
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireCloseWirelessOpenState : public PluginState
{
    Q_OBJECT
public:
    explicit WireCloseWirelessOpenState(QObject *parent = nullptr);
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireConnetingState : public PluginState
{
    Q_OBJECT
public:
    explicit WireConnetingState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireLinkCheck() override;
    void wireIpCheck() override;
    void wireInternetCheck() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireNoCableState : public PluginState
{
    Q_OBJECT
public:
    explicit WireNoCableState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireNoIpState : public PluginState
{
    Q_OBJECT
public:
    explicit WireNoIpState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WireNoInternetState : public PluginState
{
    Q_OBJECT
public:
    explicit WireNoInternetState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WirelessOpenState : public PluginState
{
    Q_OBJECT
public:
    explicit WirelessOpenState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WirelessCloseState : public PluginState
{
    Q_OBJECT
public:
    explicit WirelessCloseState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WirelessNoIp : public PluginState
{
    Q_OBJECT
public:
    explicit WirelessNoIp(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class WirelessNoInternet : public PluginState
{
    Q_OBJECT
public:
    explicit WirelessNoInternet(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class AllOpenState : public PluginState
{
    Q_OBJECT
public:
    explicit AllOpenState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class AllCloseState : public PluginState
{
    Q_OBJECT
public:
    explicit AllCloseState(QObject *parent = nullptr);
    void openWire() override;
    void closeWire() override;
    void wireConnect() override;
    void openWireless() override;
    void closeWireless() override;
    void wirelessConnect() override;
};

class PluginCase : public QObject
{
    Q_OBJECT
public:
    explicit PluginCase(QObject *parent = nullptr);

    void setState(PluginState *state);
    PluginState *getState();

    void openWire();
    void closeWire();
    void wireConnect();
    void wireConnecting();
    void wireLinkCheck();
    void wireIpCheck();
    void wireInternetCheck();
    void openWireless();
    void closeWireless();
    void wirelessConnect();

signals:

public slots:

public:
    WireOpenState *wireOpenState;
    WireCloseState *wireCloseState;
    WireConnetingState *wireConnectingState;
    WireNoCableState *wireNoCableState;
    WireNoIpState *wireNoIpState;
    WireNoInternetState *wireNoInternetState;
    WirelessOpenState *wirelessOpenState;
    WirelessCloseState *wirelessCloseState;
    WirelessNoIp *wirelessNoIp;
    WirelessNoInternet *wirelessNoInternet;
    bool cablePluged;

private:
    PluginState *m_state = nullptr;
};

#endif // PLUGINCASE_H
