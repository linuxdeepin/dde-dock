## 安装 `bluetooth firmware`

在 `Linux` 上遇到的蓝牙问题, 通常是 `firmware` 的问题. 所以当觉得自己的蓝牙出问题时, 首先应该考虑去更换 `firmware`, 然后测试能否修复问题.

本文就介绍下如何更换 `firmware`, 步骤大致如下:
1. 找出本机使用的 `frmware`

    使用蓝牙连接一个设备,然后查看 `sudo dmesg` 的输出, 这时可以看到本机需要的 `firmware`. 
    
    比如 `Thinkpad X230` 的既是: `[ 3184.931091] bluetooth hci0: firmware: failed to load brcm/BCM20702A1-0a5c-21e6.hcd (-2)`, 
    这就表示没有 `/lib/firmware/brcm/BCM20702A1-0a5c-21e6.hcd` 这个 `firmware`.

2. 找出本机蓝牙的 `ID`

    使用 `lsusb` 就可以看到, 如 `Thinkpad X230` 的既是: `ID 0a5c:21e6 Broadcom Corp. BCM20702 Bluetooth 4.0`, `ID` 就是 `0a5c 21e6`.
    
3. 去下载 `window` 下此蓝牙设备最新的驱动

    可以到本机的官网上去下载, 也可以到驱动精灵等网站上下载
    
4. 解压下载到的驱动
    
    如果驱动是 `.exe` 格式, 可以使用 `innoextract` 来解压
    
5. 查找蓝牙 `ID` 对应的 `firmware` 文件

    如果解压后的文件中只有一个文件与你需要的 `firmware` 扩展名相同, 那就是它.
    
    如果有很多, 如 `.hex` 扩展名的, 就需要查看 `.inf` 文件找到与 `ID` 对应的 `firmware` 文件.

    如果是 `.hex` 文件, 还需要使用 `hex2hcd` 来转换成 `.hcd` 格式.
    
6. 安装 `firmware`

    如果找到的 `window firmware` 与 `dmesg` 中需要的后缀一样, 就先备份一下系统中现在的 `firmware`, 然后将 `widnow firmware` 重命名复制到 `dmesg` 中得到的位置.
    
经过上面的操作, 就完成了 `firmware` 的安装, 然后执行以下命令, 来使用新的 `firmware`:

```shell
sudo modprobe -r btusb
sudo modprobe btusb
sudo systemctl restart bluetooth.service
```

注意观察有没有错误, 没有就开始连接一个设备, 测试是否成功. 如果失败就看下 `sudo dmesg` 的输出, 是否有明显的错误提示, 然后 `google` 或继续更换 `window firmware` 版本.


## 示例

下面以 `Thinkpad X230` 来详细描述下更换的过程.

首先连接一个设备, 然后获取 `sudo dmesg` 和 `lsusb` 的输出.

```shell
$ sudo dmesg
[ 3184.913914] Bluetooth: hci0: BCM: chip id 63
[ 3184.929900] Bluetooth: hci0: BCM20702A
[ 3184.931077] Bluetooth: hci0: BCM20702A1 (001.002.014) build 0000
[ 3184.931091] bluetooth hci0: firmware: failed to load brcm/BCM20702A1-0a5c-21e6.hcd (-2)
[ 3184.931096] bluetooth hci0: Direct firmware load for brcm/BCM20702A1-0a5c-21e6.hcd failed with error -2
[ 3184.931098] Bluetooth: hci0: BCM: Patch brcm/BCM20702A1-0a5c-21e6.hcd not found
$
$ lsusb
Bus 001 Device 003: ID 0a5c:21e6 Broadcom Corp. BCM20702 Bluetooth 4.0
...

```

这样可以知道 `firmware` 文件为: `/lib/firmware/brcm/BCM20702A1-0a5c-21e6.hcd`, 设备 `ID` 为: `0a5c 21e6`. 

然后去搜索 `window` 下的驱动, 选定驱动为 `g4wb12ww.exe`.然后下载, 完成后使用 `innoextract` 解压它. 发现里面有很多 `.hex` 文件, 然后搜索含有 `0a5c 21e6` 关键字的文件, 最后选定了文件 `app/Win64/bcbtums-win7x64-brcm.inf`.

接着打开这个文件, 搜索 `21e6` 来找到对应的 `.hex` 文件, 最终找到为 `BCM20702A1_001.002.014.0449.0462.hex`, 内容片段如下:

```shell
;;;;;;;;;;;;;RAMUSB21E6;;;;;;;;;;;;;;;;;

[RAMUSB21E6.CopyList]
bcbtums.sys
BCM20702A1_001.002.014.0449.0462.hex

[RAMUSB21E6.NTamd64]
Include=bth.inf
Needs=BthUsb.NT
...
```

接着执行 `hex2hcd BCM20702A1_001.002.014.0449.0462.hex` 转换成 `.hcd` 文件, 然后复制到 `/lib/firmware/brcm/BCM20702A1-0a5c-21e6.hcd`, 接着重启蓝牙相关服务来测试.

这样整个过程就完成了, 一般更换了正确的 `firmware` 都会让蓝牙问题消失或减轻, 如果问题没有消失, 可以尝试不同版本 `window driver`, 找到一个工作最好的.


最后贴一下遇到过有问题的蓝牙设备及解决方法: [已知蓝牙设备解决方法](bluetooth_device-known.md)


## 参考链接

* [[TUTORIAL] install Broadcom bluetooth on Lenovo ThinkPad X1 Carbon 1 gen ](https://forums.kali.org/showthread.php?37121-TUTORIAL-install-Broadcom-bluetooth-on-Lenovo-ThinkPad-X1-Carbon-1-gen)
