# com.deepin.daemon.Authority 服务

这是个在 system bus 上的 DBus 服务，提供给 UI 界面统一的认证服务。可以支持 登录界面、锁屏界面和 policykit 权限认证框的认证。

对于必须要进行 pam 认证的程序，如登录界面与 policykit 权限认证框 ，UI 程序需要先在本服务上进行认证，认证通过后可获得一个 cookie，把这个 cookie 作为 pam 认证的密码。

对于不需要进行 pam 认证的程序，如锁屏界面，UI 程序在本服务上进行认证，认证通过后即可放行。

## Authority 入口对象

对象路径：/com/deepin/daemon/Authority

### com.deepin.daemon.Authority 接口

#### 方法

Start(String authType, String user, Object Path agentObj) -> (Object Path transcationObj)

开始一次 pam 认证事务

- authType 只能是 fprint 或 keyboard, fprint 用于指纹识别，keyboard 用于键盘输入密码识别。

- user 要认证的用户，可以为空

- agentObj 实现了 agent 接口的对象路径

- transcationObj 返回的 pam 认证事务对象路径


---
CheckCookie(String user, String cookie) -> (Bool result)

检查 cookie 是否有效，仅供内部使用。

- user 用户名



## 认证事务

对象路径：/com/deepin/daemon/Authority/TranscationN

### com.deepin.daemon.Authority.Transcation 接口


#### 方法

Authenticate()

异步地执行认证，调用后不等认证成功就会返回。

---
End() 

结束此认证事务，会删除相关事务生成的cookie，因此必须在使用 cookie 进行 pam 认证后调用此方法，否则会导致 pam 认证失败。

---
SetUser(String user)

SetUser 修改用户名

注意以上方法都只能被注册了 agent 的连接调用，不能在 d-feet 中调试。


#### 属性

Authenticating Bool 表示是否正在进行认证


## Agent 对象

UI 程序在 system bus 上导出一个对象，这个对象实现 com.deepin.daemon.Authority.Agent 接口，把对象的路径作为 com.deepin.daemon.Authority.Start 方法的 agentObj 参数。

### com.deepin.daemon.Authority.Agent 接口

#### 方法

RequestEchoOff (String msg) -> (String result)

RequestEchoOff  从用户获取不回显的文本，比如密码

- msg 提示信息

- result 用户输入的文本

---
RequestEchoOn(String msg) -> (String result)

RequestEchoOn 从用户获取回显的文本，比如用户名

- msg 提示信息

- result 用户输入的文本

---
DisplayErrorMsg(String msg)

DisplayErrorMsg 向用户展示错误信息

- msg 错误信息

---
DisplayTextInfo(String msg)

DisplayTextInfo 向用户展示文本提示信息

- msg 文本提示信息

---
RespondResult(Stirng cookie)

RespondResult 响应认证的结果

- cookie 认证凭证，如果为空，表示认证失败，如果非空，表示认证成功。
