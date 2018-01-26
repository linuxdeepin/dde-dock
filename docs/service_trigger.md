# service_trigger 模块

通过编写配置文件，监听某种信号，触发 session 级别的命令执行。目前只实现了 DBus 信号的监听。



## 代码位置
二进制可执行文件: dde-session-daemon

代码: service_trigger 目录

## 配置文件
格式 json，后缀 .service.json。不会自动重新加载配置文件。

### 目录优先级
由高到低
- /etc/deepin-daemon/service-trigger/
- /usr/lib/deepin-daemon/service-trigger/


### 实例
文件名: apps-status-saved.service.json

```json
{
    "Monitor": {
        "Type": "DBus",
        "DBus": {
            "BusType": "System",
            "Sender": "com.deepin.daemon.Apps",
            "Interface": "com.deepin.daemon.Apps.LaunchedRecorder",
            "Path": "/com/deepin/daemon/Apps",
            "Signal": "StatusSaved"
        }
    },

    "Name": "test signal args",
    "Description": "...",
    "Exec": ["echo", "%arg1", "%arg2", "%arg3"]
}
```

这样会让 service_trigger 监听 system bus 的 com.deepin.daemon.Apps 服务的 /com/deepin/daemon/Apps 对象的 com.deepin.daemon.Apps.LaunchedRecorder 接口的 StatusSaved 信号，如果收到将执行 echo 命令输出信号参数。

### 详述
Name 名称，字符串，必填；

Description 描述，字符串，选填；

Exec 要执行的命令，字符串列表，必填，命令的参数可以使用 %argN ，表示信号的第N个参数， N 从1开始；

Monitor.Type 监听类型，字符串，目前只能为 "DBus"，"DBus" 用于监听 DBus 信号；

当 Monitor.Type 为 "DBus" 时，Monitor.DBus 不能为空；

Monitor.DBus.BusType  要监听的 bus 类型，字符串，必填，只能为 "Session" 或 "System"；

Monitor.DBus.Sender 发送者，字符串，必填，作为 dbus match rule 中的 sender；

Monitor.DBus.Interface 接口名，字符串，必填，作为 dbus match rule 中的 interface;

Monitor.DBus.Path 对象路径，字符串，选填，作为 dbus match rule 中的 path;

Monitor.DBus.Signal 信号名，字符串，选填，作为 dbus match rule 中的 member;