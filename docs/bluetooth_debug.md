这里介绍一些蓝牙问题的调试工具及方法.

* **hciconfig**

    **注意:** 新版本的 `bluez` 使用 `btmgmt` 替换了 `hciconfig`.

    查看及更改本机蓝牙设备的状态, 例如:
    
    ```shell
    $ sudo hciconfig
    hci0:	Type: Primary  Bus: USB
	BD Address: 08:3E:8E:E5:83:F6  ACL MTU: 1021:8  SCO MTU: 64:1
	UP RUNNING 
	RX bytes:19528 acl:0 sco:0 events:2828 errors:0
	TX bytes:54242 acl:0 sco:0 commands:2378 errors:0
    

    ```
    
    可更改的状态: `UP`, `DOWN`, `PSCAN`(Page Scan状态表示设备可被连接), `ISCAN`(Inquiry Scan状态表示设备可被inquiry), `PISCAN`.
    

* **btmon**

    监听蓝牙事件, 如蓝牙意外断开时, 就可以用它来看是收到什么事件断开的, 为进一步调试提供基础.


* **rfkill**

    查看蓝牙及无线设备是否被 `blocked`, 如:
    
    ```shell
    $ sudo rfkill list
    0: ideapad_wlan: Wireless LAN
            Soft blocked: no
            Hard blocked: no
    1: ideapad_bluetooth: Bluetooth
            Soft blocked: no
            Hard blocked: no
    2: hci0: Bluetooth
            Soft blocked: no
            Hard blocked: no
    3: phy0: Wireless LAN
            Soft blocked: no
            Hard blocked: no

    ```

    如果设备被 `blocked`, 意味着不可用, 需要执行 `unblock` 操作, 如 `sudo rfkill unblock <id>`.


* **开启蓝牙的内核日志**

    蓝牙的驱动在 `drivers/bluetooth` 下, 对外通信是使用的 `socket`, 所以代码在 `net/bluetooth` 中.
    
    使用 `echo 'file net/bluetooth/* +p' | tee /sys/kernel/debug/dynamic_debug/control` 来查看 `net/bluetooth` 的日志.
    
    使用 `echo 'file drivers/bluetooth/* +p' | tee /sys/kernel/debug/dynamic_debug/control` 来查看 `net/bluetooth` 的日志.


* **查看日志**

    `kernel` 的日志通过 `sudo dmesg` 查看, 用户连接的日志通过 `sudo journalctl -u bluetooth` 来查看
