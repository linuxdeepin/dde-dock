# airplane mode 模块

在 system bus 上提供服务，服务名 com.deepin.daemon.AirplaneMode, 只有一个对象 /com/deepin/daemon/AirplaneMode，这个对象只有一个接口 com.deepin.daemon.AirplaneMode。

## 接口 com.deepin.daemon.AirplaneMode

### 属性
Enabled bool  飞行模式是否打开
WifiEnabled bool wifi无线是否打开
BluetoothEnabled bool 蓝牙是否打开

### 方法

Enable(enabled bool) -> ()

开启或关闭飞行模式

---

EnableWifi(enabled bool) -> ()

开启或关闭无线wifi

---

EnableBluetooth() ->()
开启或关闭蓝牙
