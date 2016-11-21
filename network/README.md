**描述**: deepin 网络后端, 主要是对 NetworkManager DBus 接口进行包装,
提供全面的网络管理功能, 包括有线, 无线, PPPoE 拨号连接, 3G上网卡,
VPN(已支持 L2TP, PPTP, OpenConnect, OpenVPN, StrongSwan, VPNC等至少6
种 VPN 类型), 系统代理等功能.

**名词解释**:

- 字段(Setting): NetworkManager 网络连接原生支持的字段, 如
  `NM_SETTING_WIRED_SETTING_NAME`(802-3-ethernet)
- 键值(Key): NetworkManager 网络连接原生支持的键值, 如
  `NM_SETTING_WIRED_MTU`(mtu)
- 虚拟字段(Section 或 Vsection): 为方便用户配置网络, deepin 网络后
  端自定义的虚拟字段, 其值一般以 "vs-" 开头. 简单的理解, 前端控制中心网络编辑页面所显示的分类
  段落即分别对应一个虚拟字段, 如 `NM_SETTING_VS_GENERAL`("vs-general", 显
  示名称为 "General")
- 虚拟键值(Virtual Key): 为方便用户配置网络, deepin 网络后端自定义的虚
  拟键值, 其值一般以 "vk-" 开头. 如为了方便配置 3G/4G 网络, 默认在前端
  显示`NM_SETTING_VK_MOBILE_COUNTRY`("vk-mobile-country", 显示名称为
  "Country or region"),
  `NM_SETTING_VK_MOBILE_PROVIDER`("vk-mobile-provider", 显示名称为
  "Provider") 和`NM_SETTING_VK_MOBILE_PLAN`("vk-mobile-plan", 显示名称
  为 "Plan") 等几个虚拟键值, 而不是让用户手动配置 APN 等专业信息.

  另虚拟键值包括三类:
  1. wrapper, 对其他原生键值进行包装, 这种情况下前端只显示 wrapper 虚
     拟键值, 而隐藏对应的原生键值, 一般有两种用途:
     1. 改善键值设置时的交互方式, 例如"vk-no-permission" 对
        `NM_SETTING_CONNECTION_PERMISSIONS` 进行包装, 把原来需要用户手
        动输入用户名进行设置权限的交互方式改为开关(ktypeBoolean)
     1. 将某一个原生键分拆为多个虚拟子键(child key), 如
        `NM_SETTING_IP4_CONFIG_ADDRESSES` 被分拆成了
        "vk-addresses-address", "vk-addresses-mask" 和
        "vk-addresses-gateway"
  1. enable-wrapper, 对其他键值的可用性添加开关, 如 "vk-enable-mtu"
  1. controller, 纯控制型虚拟键, 一般没有相关的原生键, 如
     "vk-vpn-missing-plugin", 用来在前端显示所缺少的 VPN 插件包.

## 目录结构

- **examples**: 调用 deepin 网络后端 DBus 接口的 shell, python 脚本示
  例.

- **nm**: 定义 NetworkManager 常量的子库, 具体参考其
  [README](./nm/README.md).

- **nm_generator**: 辅助 NetworkManager Go 代码生成器, 后缀名
  为'*_gen.go' 的文件均由其生成, 包括 `nm/nm_consts_gen.go` 和
  `nm_setting_beans_gen.go`, 具体参考其
  [README](./nm_generator/README.md).

- **agent.go**: 实现 NetworkManager secret agent 的主要文件, 主要用于
  处理密码弹出框等问题, 相关文档可参考 `nm_generator/nm_docs`, 代码逻
  辑可参考
  [applet-agent.c](https://github.com/GNOME/network-manager-applet/blob/master/src/applet-agent.c).

- **connection_session.go**: ConnectionSession DBus 接口主要文件, 用于
  给前端提供接口来编辑网络连接, 配合 GetKey, SetKey, AvailableKeys,
  GetAvailableValues 等接口, 前端可以生成相关的控件界面, 并根据后端信
  号按需进行展示.

- **manager_accesspoint.go**: 处理 WiFi 热点及相关 DBus 接口.

- **manager_active.go**, **dbus_watcher.go**: 手动注册 DBus watcher监
  听 NetworkManager 所有激活连接的状态变更以避免调用 dbus-factory 接口
  一定概率导致信号不同步的问题, 同时提供了接口用来获取当前激活连接的相
  关信息.

- **manager_config.go**: 处理 deepin 网络后端的配置文件
  (~/.config/deepin/network.json), 主要保存一些网络开关状态和 VPN 自动
  连接的配置.

- **manager_connection.go**: 创建/激活/删除网络连接及相关 DBus 接口.

- **manager_device.go**: 处理网卡设备及相关 DBus 接口.

- **manager.go**: 主 Manager DBus 对象.

- **manager_proxy.go**: 处理系统代理及相关 DBus 接口.

- **manager_switch.go**: 处理设备开关及相关 DBus 接口. deepin 为每个网
  卡都单独提供一个虚拟开关, 同时会兼容 NetworkManager 本身的逻辑流程.

- **nm_custom_type.go**: 自定义的一些字符串格式的网络设备类型和网络连
  接类型, 主要方便与前端交互.

- **nm_key_xxx.go**: 通过 NetworkManager DBus 接口编辑网络连接时的一些辅
  助方法, 同时自定义了对应的键值类型, 包括添加了一些包装类型, 如
  `ktypeWrapperIpv4Addresses`, 用于方便处理相关逻辑. `nm_generator` 基
  本是围绕这部分内容生成 NetworkManager 键值的 getter, setter.

- **nm_setting_virtual_key.go**: 定义虚拟键值(virtual key) 及相关
  逻辑代码, 另外某些 "controller" 类型的虚拟键的 getter, setter, available
  key 即 available value 等逻辑代码也会定义到这里.

- **nm_setting_virtual_section.go**: 定义虚拟字段(virtual section) 及
  相关逻辑代码, 包括获取虚拟字段对应的原生字段列表, 判断虚拟字段对应的
  前端控件默认是否展开等.

- **nm_setting_beans_gen.go**: 通过 `nm_generator` 生成的辅助代码, 因
  为 NetworkManager 的字段键值比较多, 所以代码量比较大, 主要包括下面集
  中类型:

  - 虚拟字段列表及显示名称(翻译文本)
  - 虚拟键值列表, 包括前端展示时所用到的控件类型, 显示名称(翻译文本)
    和取值范围(可选)
  - 虚拟键值和原生键值的 getter 和 setter 代码
  - generalXXX 开头的通用辅助方法, 如 `generalGetSettingKeyType`,
    `generalGetSettingAvailableKeys`, `generalGetSettingDefaultValue`
    等

- **nm_setting_beans_extend.go**: 因为 `nm_generator` 会把
  NetworkManager 支持的所有字段(Setting)和键值(Key)都生成出来, 然后可
  能部分字段用不到, 为了保证编译通过, 需要补充
  getSettingXXXAvailableKeys, getSettingXXXAvailableValues,
  checkSettingXXXValues 等空方法.

- **nm_setting_xxx.go**: 放置 NetworkManager 字段(Setting)相关的逻辑代
  码, 主要包括 getSettingXXXAvailableKeys,
  getSettingXXXAvailableValues 和 checkSettingXXXValues.

- **nm_setting_vpn.go**: 放置 VPN 通用逻辑代码, 一部分用来处理 VPN 密
  码弹出框(一般调用 VPN 自己的密码弹出框, 如
  `/usr/lib/NetworkManager/nm-xxx-auth-dialog`), 另一部分特殊处理 VPN
  子键, 将其保存到 `NM_SETTING_VPN_DATA` 和 `NM_SETTING_VPN_SECRETS`这
  两个键(string dictionary)里面.

- **nm_setting_vpn_xxx.go**: 放置特定 VPN 相关的逻辑代码, 主要是
  getSettingVpnXXXAvailableKeys, getSettingVpnXXXAvailableValues,
  checkSettingVpnXXXValues.

- **state_handler.go**: 监听网络设备状态变更并按需弹出系统通知.

- **utils_xxx.go**: 一些辅助方法, 如处理 IPv6 地址, 读写 gnome-keyring,
  包装系统通知接口, 包装 NetworkManager 和 ModemManager DBus 接口等.

## NetworkManager DBus 接口简介

- `/org/freedesktop/NetworkManager`: NetworkManager DBus 主接口, 可用获
  取当前网络状态, 设备列表, 连接列表, 用户权限等
- `/org/freedesktop/NetworkManager/Devices/XXX`: 特定网络设备接口
- `/org/freedesktop/NetworkManager/Settings`: 配置文件相关接口
- `/org/freedesktop/NetworkManager/Settings/XXX`: 特定配置文件接口
- `/org/freedesktop/NetworkManager/AccessPoint/XXX`: WiFi 热点接口
- `/org/freedesktop/NetworkManager/ActiveConnection/XXX`: 已激活连接
  接口
- `/org/freedesktop/NetworkManager/IP4Config/XXX`: 当前分配的 IP4 地址接
  口
- `/org/freedesktop/NetworkManager/IP6Config/XXX`: 当前分配的 IP6 地址接
  口

具体请参考 `nm_generator/nm_docs`.

## deepin 网络后端 DBus 接口简介

**DBus 参数命名规范**:
- `uuid`: string 类型, 表示配置文件唯一的 UUID 值
- `apPath`: dbus.ObjectPath 类型, 表示 AccessPoint 对应的 DBus 路径, 如
  /org/freedesktop/NetworkManager/AccessPoint/XXX
- `devPath`: dbus.ObjectPath 类型, 表示网络设备对应的 DBus 路径, 如
  /org/freedesktop/NetworkManager/Devices/XXX
- `cpath`: dbus.ObjectPath 类型, 表示配置文件对应的 DBus 路径, 如
  /org/freedesktop/NetworkManager/Settings/XXX
- `connType`: string 类型, 表示配置文件类型, 如
  `connectionWired`("wired"), 全部定义位于 `nm_custom_type.go`.
- `xxxJSON`: JSON string 类型, 为了方便前后端交互, 同时减少 DBus 参数数
  量, deepin 网络后端用了很多 JSON string 类型的参数

### com.deepin.daemon.Network

- 当前网络状态
  - `GetActiveConnectionInfo() (acinfosJSON string)`
  - **prop** `State uint32`
  - **prop** `Devices string`
  - **prop** `Connections string`
  - **prop** `ActiveConnections string`

- 网络开关
  - `EnableDevice(devPath dbus.ObjectPath, enabled bool)`
  - `IsDeviceEnabled(devPath dbus.ObjectPath) (enabled bool)`
  - `SetDeviceManaged(devPathOrIfc string, managed bool)`
  - **prop-rw** `NetworkingEnabled bool`
  - **prop-rw** `VpnEnabled bool`
  - **signal** `DeviceEnabled func(devPath string, enabled bool)`

- 编辑网络连接
  - `CreateConnection(connType string, devPath dbus.ObjectPath) (session *ConnectionSession)`
  - `CreateConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (session *ConnectionSession)`
  - `DeleteConnection(uuid string)`
  - `EditConnection(uuid string, devPath dbus.ObjectPath) (session *ConnectionSession)`
  - `GetSupportedConnectionTypes() (types []string)`

- 激活网络连接
  - `ActivateConnection(uuid string, devPath dbus.ObjectPath) (cpath dbus.ObjectPath)`
  - `DeactivateConnection(uuid string)`
  - `DisconnectDevice(devPath dbus.ObjectPath)`
  - `GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string)`

- WiFi AccessPoint
  - `ActivateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (cpath dbus.ObjectPath)`
  - `GetAccessPoints(path dbus.ObjectPath) (apsJSON string)`
  - **signal** `AccessPointAdded func(devPath, apJSON string)`
  - **signal** `AccessPointRemoved func(devPath, apJSON string)`
  - **signal** `AccessPointPropertiesChanged func(devPath, apJSON string)`

- WiFi Hotspot 热点
  - `DisableWirelessHotspotMode(devPath dbus.ObjectPath)`
  - `EnableWirelessHotspotMode(devPath dbus.ObjectPath)`
  - `IsWirelessHotspotModeEnabled(devPath dbus.ObjectPath) (enabled bool)`

- 弹出密码输入框
  - `CancelSecret(path string, settingName string)`
  - `FeedSecret(path string, settingName, keyValue string, autoConnect bool)`
  - **signal** `NeedSecrets func(connPath, settingName, connectionId string, autoConnect bool)`
  - **signal** `NeedSecretsFinished func(connPath, settingName string)`

- 系统代理
  - `GetAutoProxy() (proxyAuto string)`
  - `GetProxy(proxyType string) (host, port string)`
  - `GetProxyIgnoreHosts() (ignoreHosts string)`
  - `GetProxyMethod() (proxyMode string)`
  - `SetAutoProxy(proxyAuto string)`
  - `SetProxy(proxyType, host, port string)`
  - `SetProxyIgnoreHosts(ignoreHosts string)`
  - `SetProxyMethod(proxyMode string)`

### com.deepin.daemon.Network.ConnectionSession

- DBus 属性
  - **prop** `ConnectionPath dbus.ObjectPath`
  - **prop** `Uuid string`
  - **prop** `Type string`
  - **prop** `AllowDelete bool`
  - **prop** `AllowEditConnectionId bool`
  - **prop** `AvailableVirtualSections []string`
  - **prop** `AvailableSections []string`
  - **prop** `AvailableKeys map[string][]string`
  - **prop** `Errors sessionErrors`

- DBus 信号
  - **signal** `ConnectionDataChanged func()`

- DBus 接口
  - `GetAllKeys() (infoJSON string)`
  - `GetAvailableValues(section, key string) (valuesJSON string)`
  - `GetKeyName(section, key string) (name string)`
  - `GetKey(section, key string) (valueJSON string)`
  - `SetKey(section, key, valueJSON string)`
  - `IsDefaultExpandedSection(vsection string) bool`
  - `Save() (ok bool)`
  - `Close()`

- Debug 辅助接口
  - `DebugGetConnectionData() connectionData`
  - `DebugGetErrors() sessionErrors`
  - `DebugListKeyDetail() (info string)`

详细 DBus 接口信息请参考
[godoc 文档](https://godoc.org/github.com/linuxdeepin/dde-daemon/network).

## HACKING 实践

### 添加一个虚拟键值

1. 编辑 `nm_generator/nm_virtual_sections.yml`, 在适当位置添加虚拟键值
   的定义, 如

   ```
   - KeyValue: vk-autoconnect
     Section: connection
     DisplayName: Automatically connect
     WidgetType: EditLineSwitchButton
     VKeyInfo:
       VirtualKeyName: NM_SETTING_VK_CONNECTION_AUTOCONNECT
       Type: ktypeBoolean
       VkType: vkTypeWrapper
       RelatedKeys:
       - NM_SETTING_CONNECTION_AUTOCONNECT
       ChildKey: false
       Optional: false
   ```

1. 运行 `make gen-nm-code` 重新生成代码

1. 编辑 `nm_setting_connection.go`, 实现
   `getSettingVkConnectionAutoconnect` 和
   `logicSetSettingVkConnectionAutoconnect`

1. 编辑 getSettingXXXAvailableKeys 控制何时显现该虚拟键

1. 如果需要检测用户输入信息的有效性, 还需要在
   checkSettingXXXValues 添加相关代码

1. 如果控件类型为 `EditLineComboBox` 或 `EditLineEditComboBox`, 还需要
   在 getSettingXXXAvailableValues 返回有效值列表

1. 处理国际化

### 支持新的 VPN 类型

1. 以 OpenVPN 为例, 查看 network-manager-openvpn 源码, 检索
   `src/nm-openvpn-service.c` 找到所有 VPN 键值名称, 并将它们定义到
   `nm_generator/nm_vpn_alias_settings.yml` 和
   `nm/nm_extends_consts.go`

1. 编辑 `nm_generator/nm_virtual_sections.yml`, 在
   `NM_SETTING_VS_VPN`虚拟字段下将需要展现到前端的键值按顺序补充完整,
   如果有必要还可用添加其他虚拟字段, 如
   `NM_SETTING_VS_VPN_OPENVPN_PROXIES`

1. 如果要为某些 VPN 键值添加 logic setter, 则编辑
   `nm_generator/nm_logicset_keys.yml`

1. 运行 `make gen-nm-code` 重新生成代码

1. 新建文件 `nm_setting_vpn_openvpn.go` 将相关代码补充完整

1. 编辑 `nm_custom_type.go`, 添加自定义 VPN 类型
   `connectionVpnOpenvpn`, 同时补充实现 `getCustomConnectionType`,
   `isVsectionExpandedDefault`, `doGetRelatedVsections` 和
   `doGetRelatedSectionsOfVsection`

1. 编辑 `nm_setting_virtual_key.go`, 补充实现
   `getLocalSupportedVpnTypes`, `getVpnNameFile`,
   `logicSetSettingVkVpnType` 和 `getSettingVkVpnMissingPlugin`

1. 处理国际化

## go test 单元测试

普通单元测试执行 `go test` 即可, 如果要执行 dev_test.go 里的内容则需要
添加 `dev` tag:

1. 获取当前系统支持的 VPN 类型

   ```sh
   go test -tags dev -gocheck.f TestLocalSupportedVpnTypes
   ```

1. 单独执行网络后端
   ```sh
   env DDE_DEBUG=t go test -tags dev -gocheck.f TestMain
   ```

## 自动化测试

目前配合 ansbile/docker/openwrt 等技术, 易经实现 deepin 网络功能自动化
测试, 具体请参考
[deepin-network-tests](https://github.com/x-deepin/deepin-network-tests)

如果要通过 iperf3 等工具测试网卡驱动的稳定性, 则参考 [network-testing-experiments](https://github.com/x-deepin/network-testing-experiments)

## DBus 脚本示例

1. python 脚本调用网络 DBus 接口示例, 连接 WiFi 网络

   [examples/python/main.py](./examples/python/main.py)

1. shell 脚本调用网络 DBus 接口示例, 设置有线网卡静态 IP 地址

   [examples/set_wired_static_ip.sh](./examples/set_wired_static_ip.sh)

1. 根据有线网卡地址获取其对应配置的 UUID

   ```sh
   macaddr_to_uuid() {
     local md5="$(printf "$1" | md5sum | awk '{print $1}' | tr 'A-Z' 'a-z')"
     echo "${md5:0:8}-${md5:8:4}-${md5:12:4}-${md5:16:4}-${md5:20:12}"
   }
   $ macaddr_to_uuid "00:12:34:56:ab:cd"
   > 086e214c-1f20-bca4-9816-c0a11c8c0e02
   ```

1. 监听 NetworkManager 服务的运行状态

   ```sh
   dbus-monitor --system sender=org.freedesktop.DBus,member=NameOwnerChanged
   ```

1. 监听 NetworkManager 配置文件的变更

   ```sh
   dbus-monitor --system sender=org.freedesktop.NetworkManager,interface=org.freedesktop.NetworkManager.Settings
   ```

1. 开关 NetworkManager 飞行模式

   ```sh
   dbus-send --system --type=method_call --print-reply \
       --dest=org.freedesktop.NetworkManager /org/freedesktop/NetworkManager  \
       org.freedesktop.DBus.Properties.Get \
       string:"org.freedesktop.NetworkManager" string:"NetworkingEnabled"
   dbus-send --system --type=method_call --print-reply \
       --dest=org.freedesktop.NetworkManager /org/freedesktop/NetworkManager  \
       org.freedesktop.NetworkManager.Enable boolean:"false"
   ```

1. 开关 NetworkManager WiFi 无线网络功能

   ```sh
   dbus-send --system --type=method_call --print-reply \
       --dest=org.freedesktop.NetworkManager /org/freedesktop/NetworkManager  \
       org.freedesktop.DBus.Properties.Set string:"org.freedesktop.NetworkManager" \
       string:"WirelessEnabled" variant:boolean:"false"
   ```

1. 设置 NetworkManager 设备托管状态

   ```sh
   dbus-send --system --type=method_call --print-reply \
       --dest=org.freedesktop.NetworkManager /org/freedesktop/NetworkManager/Devices/0  \
       org.freedesktop.DBus.Properties.Set string:"org.freedesktop.NetworkManager.Device" \
       string:"Managed" variant:boolean:"false"
   ```

1. 获取 deepin 网络后端提供的所有网络设备信息

   ```sh
   dbus-send --print-reply --dest=com.deepin.daemon.Network \
       /com/deepin/daemon/Network org.freedesktop.DBus.Properties.Get \
       string:"com.deepin.daemon.Network" string:"Devices"
   ```

1. 设置 deepin 网络后端提供的 WiFi 无线, VPN 等开关状态

   ```sh
   dbus-send --print-reply --type=method_call \
       --dest=com.deepin.daemon.Network /com/deepin/daemon/Network \
       org.freedesktop.DBus.Properties.Set \
       string:"com.deepin.daemon.Network" string:"NetworkingEnabled" variant:boolean:"true"
   dbus-send --print-reply --type=method_call \
       --dest=com.deepin.daemon.Network /com/deepin/daemon/Network \
       org.freedesktop.DBus.Properties.Set \
       string:"com.deepin.daemon.Network" string:"WirelessEnabled" variant:boolean:"true"
    dbus-send --print-reply --type=method_call \
       --dest=com.deepin.daemon.Network /com/deepin/daemon/Network \
       org.freedesktop.DBus.Properties.Set \
       string:"com.deepin.daemon.Network" string:"VpnEnabled" variant:boolean:"true"
   ```

## 相关网络配置文件

- **/etc/NetworkManager/system-connections/**: NetworkManager 的所有连
 接对应的配置文件所在目录, 可以手动创建配置文件放置到该目录下, 并确保
 文件权限为 0600, 然后重启 NetworkManager 就可以让其生效

  ```sh
  sudo systemctl restart NetworkManager
  ```

- **/etc/NetworkManager/NetworkManager.conf**: NetworkManager 自身的配
  置文件, 可以配置 DHCP, DNS 等选项, 如:

  1. 将网络 wlan0 设置为未托管
      ```
      [keyfile]
      unmanaged-devices=interface-name:wlan0
      ```

  1. 在 syslog 显示更详细的 NetworkManager 日志
     ```
     [logging]
     level=DEBUG
     ```

  具体请参考 `man NetworkManager.conf`

- **/etc/network/interfaces**: Linux 系统默认的网络接口配置文件, 一般
  不再直接编辑该文件, 否则可能导致 NetworkManager 异常, 例如如果在
  interfaces 文件里配置了某个网卡, NetworkManager 便无法正常托管该网卡

- **/etc/resolv.conf**: DNS 配置文件

## 相关网络工具

简单介绍一些网络工具以便调试时使用:

- nmcli, NetworkManager 自带的终端配置工具, 使用方便, 功能强大

  1. 查看 NetworkManager 当前的状态

     ```sh
     nmcli general
     ```

  1. 查看网卡状态

     ```sh
     nmcli device
     ```

  1. 监听 NetworkManager 变更

     ```sh
     nmcli monitor
     ```

  1. 设置网卡的托管状态

     ```sh
     nmcli device set wlan0 managed yes
     ```

  1. 创建一个 WiFi 连接

     ```sh
     nmcli connection add type wifi ifname '*' con-name 'test-ssid' ssid test-ssid
     nmcli connection modify test-ssid wifi-sec.key-mgmt wpa-psk
     nmcli connection modify test-ssid wifi-sec.psk password
     ```

  1. 连接指定 WiFi 网络

     ```sh
     nmcli device wifi connect "test-ssid" password "password"
     ```

  1. 显示特定连接详细信息

     ```sh
     nmcli conn list uuid 1ad8d2a5-d84f-4775-92e8-fae4e7273a76
     ```

  1. 手动激活某个连接, 如果失败则打印详细日志

     ```sh
     nmcli -p con up id "Wired Connection" iface eth0
     ```

- nmtui, NetworkManager 1.0 后新添加的终端配置工具, 提供 ncurses 界面,
  使用比较方便

- nm-connection-editor, network-manaager-gnome 包提供的配置工具, 简单
  易用, 需要注意安装 network-manaager-gnome 后会开机自动运行nm-applet,
  建议手动将 nm-applet 禁用, 否则可能和 deepin 网络后端有冲突

- mmcli, ModemManager 终端工具

  - 列出当前所有的 modem 设备

    ```sh
    mmcli -L
    ```

  - 监听 modem 设备列表变更

    ```sh
    mmcli -M
    ```

- usb-modeswitch, 通过更改 USB Modem 设备 ID 从而使其匹配到正确的驱动,
  很多 3G/4G USB 网卡通过配置 usb-modeswitch 来修复无法识别的问题

- dbus-send, dbus-monitor, 终端 DBus 工具

  1. 获取 deepin 网络后端 DBus 接口 XML 配置

     ```sh
     dbus-send --type=method_call --print-reply --dest=com.deepin.daemon.Network /com/deepin/daemon/Network org.freedesktop.DBus.Introspectable.Introspect | sed 1d | sed -e '1s/^   string "//' | sed '$s/"$//'
     ```

  1. 监听 WiFi 热点信号强度变化

     ```sh
     dbus-monitor --system sender=org.freedesktop.NetworkManager,interface=org.freedesktop.NetworkManager.AccessPoint,member=PropertiesChanged
     ```

- qdbus, Qt 提供的终端 DBus 工具, 和 dbus-send 类似, 但更加易用, 补全
  功能强大

  ```sh
  qdbus --literal --system org.freedesktop.NetworkManager /org/freedesktop/NetworkManager org.freedesktop.NetworkManager.state
  ```

- d-feet, 方便易用的 DBus GUI 查看工具

- dhclient, NetworkManager 默认的 DHCP 工具

  - 重新为 wlan0 分配 IP 地址, 可以解决某些情况 DHCP 分配 IP 地址成功
    但无法访问网络的问题
    ```sh
    sudo dhclient -x wlan0
    sudo dhclient wlan0
    ```

- iperf3, 网络性能测试工具, 可以用来测试网卡驱动是否稳定

  - 服务端

    ```sh
    iperf -s
    ```

  - 客户端

    ```sh
    iperf -c <server-address>
    ```

- rfkill, 无线网络控制开关工具

  ```sh
  rfkill list all
  rfkill unblock all
  ```

- iw, 无线网卡配置工具

  - 扫描 WiFi 网络

    ```sh
    sudo iw dev wlan0 scan
    ```

  - 查看无线网卡是否支持 AP 热点模式

    ```sh
    sudo iw phy0 info | grep '* AP'
    ```

  - 查看无线网卡是否支持 P2P/P2p-Client 模式

    ```sh
    sudo iw phy0 info | grep '* P2P'
    ```
- iwconfig, 无线网卡接口配置工具

   - 切换无线网卡工作模式
     ```sh
     sudo /sbin/ifconfig wlan0 down
     sudo /sbin/doiwconfig wlan0 mode monitor
     sudo /sbin/doifconfig wlan0 up
     ```

   - 获取无线网卡速率

     ```sh
     /sbin/iwconfig wlan0 | grep 'Bit Rate' | awk '{print $2}' | awk -F= '{print $2}'
     ```

- lshw, 系统硬件检索工具

  - 获取网卡设备信息

    ```sh
    sudo lshw -class network
    sudo lshw -class network -businfo
    ```

  - 获取网卡驱动信息

    ```sh
    sudo lshw -xml | xpath -q -e "//node[@id='network' and ./capabilities/capability[@id='wireless']]/configuration/setting[@id='driver']/@value" | cut -d"\"" -f2
    ```

- dmidecode, 另一款系统信息检索工具

- lspci, PCI 设备检索工具, 同时会显示对应加载的内核驱动. 建议使用时添
  加 `-vvnn` 选项, 从而可以定位设备对应的唯一 ID

- lsusb, USB 设备检索工具

- usb-devices, USB 设备检索工具, 会显示对应加载的内核驱动

## 常见网络问题

1. 无线网卡 rfkill 硬开关为关闭状态导致网络不可用

   一般是网卡驱动问题, 需要更换驱动或调整驱动参数.

1. 无线网卡 rfkill 软开关为关闭状态导致网络不可用

   一般是用户不小心手动关闭无线网卡开关导致, 一般键盘上有相应的快捷键,
   重新打开即可, 也可以调用 rkill 命令:

   ```sh
   sudo rfkill unblock all
   ```

1. 无线网卡异常, 频繁提示密码错误或过一段时间网卡会断开重连一次或网速
   很慢

   一般是网卡驱动问题, 需要更换驱动或调整驱动参数.

   对于频繁提示密码错误的问题, 也可能是 wpa_supplicant 的 Bug 导致, 需
   要跟踪上游最新的状态.

1. 有线网卡无法正常连接, 一直显示正在连接

   可能是 DHCP 分配 IP 出现了问题, 可能是路由器 DHCP 的问题, 也可能是
   局域网内有多个 DHCP 服务导致.

1. 网卡驱动正常, 用其他网络工具如 wicd 也能正常访问, 但 NetworkManager
   下无法显示该网卡

   可能是 NetworkManager 将该网卡标记了未托管状态, 需要确保
   /etc/network/interfaces 即 /etc/NetworkManager/NetworkManager.conf
   没有对该网卡做特殊配置.

1. 无法通过 WiFi 连接 PPPoE 拨号网络

   目前 NetworkManager 未实现该功能, 请使用原生的 pppoeconf 命令, 提供
   了 ncurses 终端界面.

1. 设置系统代理, 但并未生效

   目前 deepin 仅支持 GNOME gsettings 系统代理和环境变量系统代理, 如果
   目标网络应用不支持这两者系统代理便无法生效.

1. 插入 3G/4G modem 上网卡后无法识别

   一般 modem 上网卡插入后需要等待 10s 左右才能识别, 可以通过 `mmcli
   -M` 命令进行监听设备识别情况, 如果一直无法识别, 则可能 Linux 暂不支
   持该设备, 部分 modem 上网卡可以通过 usb-modeswitch 更改设备类型达到
   正常识别的目的.

1. 为 3G/4G 上网卡选择套餐后无法保存

   可能是缺乏 mobile-broadband-provider-info 包导致.

1. 苹果手机(iPhone) USB 热点可以显示, 但无法使用

   可能是缺乏相关包导致, 包括 libimobiledevice, usbmuxd, libusbmuxd,
   也可能是 usbmuxd, libusbmuxd 版本不兼容导致

## 参考资料

- NetworkManager
  - [NetworkManager Developer 主页](https://developer.gnome.org/NetworkManager/)
  - [NetworkManager 在线文档](https://developer.gnome.org/NetworkManager/unstable/index.html)
  - [NetworkManager 新版本变更日志](https://cgit.freedesktop.org/NetworkManager/NetworkManager/plain/NEWS)
  - [NetworkManager Debugging](https://wiki.gnome.org/Projects/NetworkManager/Debugging)
  - [NetworkManager Bug 列表](https://bugzilla.gnome.org/browse.cgi?product=NetworkManager)
- wpa_supplicant
  - [wpa_supplicant 项目主页](https://w1.fi/wpa_supplicant/)
  - [wpa_supplicant 邮件列表](http://lists.infradead.org/pipermail/hostap/)
  - [Ubuntu/wpa_supplicant Bug 列表](https://launchpad.net/ubuntu/+source/wpasupplicant/+bugs)
- Package
  - [debian NetworkManager 打包源码](https://anonscm.debian.org/cgit/pkg-utopia/)
  - [ubuntu NetworkManager 打包源码](https://code.launchpad.net/~network-manager/network-manager/ubuntu)
