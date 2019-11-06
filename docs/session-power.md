# com.deepin.daemon.Power 服务

这是个在 Session Bus 上的服务

## Power 入口对象
 
对象路径：/com/deepin/daemon/Power

## 属性：

### PowerButtonAction
String
read/write
按下电源按钮后执行的命令
默认为 
```
dde-shutdown
```

### LidClosedAction
String
read/write
笔记本电脑关闭盖子后执行的命令
默认为 

```
dbus-send --print-reply --dest=com.deepin.SessionManager /com/deepin/SessionManager com.deepin.SessionManager.RequestSuspend
```


### LinePowerScreenBlackDelay
Int32
read/write
接通电源时，不做任何操作，到关闭屏幕需要的时间
单位：秒
值为 0 时表示从不

### LinePowerSleepDelay
Int32
read/write
接通电源时，不做任何操作，从黑屏到睡眠的时间
单位：秒
值为 0 时表示从不

### BatteryScreenBlackDelay
Int32
read/write
使用电池时，不做任何操作，到关闭屏幕需要的时间
单位：秒
值为 0 时表示从不

### BatterySleepDelay
Int32
read/write
使用电池时，不做任何操作，从黑屏到睡眠的时间
单位：秒
值为 0 时表示从不

### ScreenBlackLock
Boolean read/write
关闭显示器前是否锁定

### SleepLock
Boolean read/write
睡眠前是否锁定

### LidIsPresent
Int32
read
是否有盖子，一般笔记本电脑才有

### OnBattery
Boolean
read
是否使用电池
接通电源时为 false
使用电池时为 true

### BatteryIsPresent
Dict of {String,Boolean}
read
电池是否可用

例如：
{'BAT0':True}
表示 BAT0 可用

### BatteryPercentage
Dict of {String,Double}
电池电量百分比
例如：
{'BAT0': 50}
表示 电池 BAT0 的电量百分比是 50%

### BatteryState
Dict of {String,UInt32}
电池状态
例如：
{'BAT0': 1L}
表示 电池 BAT0 的状态为 1

状态数字代表的意义：
0 Unknown 未知
1 Charging 充电中
2 Discharging 不充电
3 Empty 空
4 FullyCharged 充满
5 PendingCharge
6 PendingDischarge

http://upower.freedesktop.org/docs/Device.html#Device:State



## 方法：

### Reset()
重置所有相关设置



