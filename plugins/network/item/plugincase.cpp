#include "plugincase.h"

PluginCase::PluginCase(QObject *parent)
    : QObject(parent)
{
    wireOpenState = new WireOpenState(this);
    wireCloseState = new WireCloseState(this);
    wireConnectingState = new WireConnetingState(this);
    wireNoCableState = new WireNoCableState(this);
    wireNoIpState = new WireNoIpState(this);
    wireNoInternetState = new WireNoInternetState(this);
    wirelessOpenState = new WirelessOpenState(this);
    wirelessCloseState = new WirelessCloseState(this);
    wirelessNoIp = new WirelessNoIp(this);
    wirelessNoInternet = new WirelessNoInternet(this);
}

void PluginCase::setState(PluginState *state)
{
    m_state = state;
    state->setPluginCase(this);
}

PluginState *PluginCase::getState()
{
    return m_state;
}

void PluginCase::openWire()
{
    m_state->openWire();
}

void PluginCase::closeWire()
{
    m_state->closeWire();
}

void PluginCase::wireConnect()
{
    m_state->wireConnect();
}

void PluginCase::openWireless()
{
    m_state->openWireless();
}

void PluginCase::closeWireless()
{
    m_state->closeWireless();
}

void PluginCase::wirelessConnect()
{
    m_state->wirelessConnect();
}

void WireOpenState::closeWire()
{
    m_pluginCase->setState(m_pluginCase->wireCloseState);
    m_pluginCase->closeWire();
}

void WireOpenState::wireConnect()
{
    m_pluginCase->setState(m_pluginCase->wireConnectingState);
    m_pluginCase->wireLinkCheck();
}

void WireConnetingState::wireLinkCheck()
{
    if (m_pluginCase->cablePluged) {
        m_pluginCase->setState(m_pluginCase->wireConnectingState);
        m_pluginCase->wireIpCheck();
    }
    else {
        m_pluginCase->setState(m_pluginCase->wireNoCableState);
    }
}
