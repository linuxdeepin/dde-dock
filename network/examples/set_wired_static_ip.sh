#! /bin/bash

# 设置第一块有线网卡的IP地址为10.1.11.231、掩码为255.255.255.0、网关为
# 10.1.11.1、域名服务器为8.8.8.8

dbus_name="com.deepin.daemon.Network"
dbus_path="/com/deepin/daemon/Network"
dev_name="/org/freedesktop/NetworkManager/Devices/0"

ip_addr='"10.1.11.231"'
netmask='"255.255.255.0"'
gateway='"10.1.11.1"'
dns_server='"8.8.8.8"'

uuid=`qdbus ${dbus_name} ${dbus_path} ${dbus_name}.GetWiredConnectionUuid ${dev_name}`
echo $uuid
sess_path=`qdbus --literal ${dbus_name} ${dbus_path} ${dbus_name}.EditConnection ${uuid} ${dev_name} | awk '{print $NF}' | awk -F] '{print $1}'`
echo $sess_path
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.SetKey ipv4 method '"manual"'
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.SetKey ipv4 vk-addresses-address ${ip_addr}
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.SetKey ipv4 vk-addresses-mask ${netmask}
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.SetKey ipv4 vk-addresses-gateway ${gateway}
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.SetKey ipv4 vk-dns ${dns_server}
qdbus ${dbus_name} ${sess_path} com.deepin.daemon.ConnectionSession.Save
