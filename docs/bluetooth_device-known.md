这里收集了一些在 `linux` 上蓝牙存在问题的设备及解决方法，由于没有环境来进行充分测试，可能里面的方法只是个例，如果不适用，请按照 [蓝牙常见问题处理](bluetooth_FAQ.md)。

有些条目由于不知到设备名，就按照 `dmesg` 里的 `firmware` 或者 `usb id` 作为名称。

## `0a5c:21e6 Broadcom Corp. BCM20702 Bluetooth 4.0`

* 问题

    此设备在连接低功耗蓝牙设备时，使用中会常常出现断开的问题，`sudo dmesg` 错误如下：
    
    ```shell
    [ 3084.215316] Bluetooth: hci0 advertising data length corrected
    [ 3084.220125] Bluetooth: hci0 advertising data length corrected
    [ 3086.117379] input: MiMouse as /devices/virtual/misc/uhid/0005:2717:0040.000C/input/input30
    [ 3086.120085] hid-generic 0005:2717:0040.000C: input,hidraw1: BLUETOOTH HID v0.23 Mouse [MiMouse] on 08:3E:8E:E5:83:F6
    [ 3176.471497] Bluetooth: Inquiry failed: status 0x1f
    [ 3176.475523] Bluetooth: hci0 command 0x0401 tx timeout
    [ 3184.119133] Bluetooth: hci0 command 0x2005 tx timeout
    [ 3184.523219] usb 1-1.4: USB disconnect, device number 14
    [ 3184.791085] usb 1-1.4: new full-speed USB device number 15 using ehci-pci
    [ 3184.904261] usb 1-1.4: New USB device found, idVendor=0a5c, idProduct=21e6
    [ 3184.904263] usb 1-1.4: New USB device strings: Mfr=1, Product=2, SerialNumber=3
    [ 3184.904264] usb 1-1.4: Product: BCM20702A0
    [ 3184.904265] usb 1-1.4: Manufacturer: Broadcom Corp
    [ 3184.904265] usb 1-1.4: SerialNumber: 083E8EE583F6
    [ 3184.913914] Bluetooth: hci0: BCM: chip id 63
    [ 3184.929900] Bluetooth: hci0: BCM20702A
    [ 3184.931077] Bluetooth: hci0: BCM20702A1 (001.002.014) build 0000
    [ 3184.931091] bluetooth hci0: firmware: failed to load brcm/BCM20702A1-0a5c-21e6.hcd (-2)
    [ 3184.931096] bluetooth hci0: Direct firmware load for brcm/BCM20702A1-0a5c-21e6.hcd failed with error -2
    [ 3184.931098] Bluetooth: hci0: BCM: Patch brcm/BCM20702A1-0a5c-21e6.hcd not found
    ```

* 解决方法

    去下载此电脑在 `window` 下最新版的蓝牙驱动，然后解压(`.exe` 文件使用 `innoextract` 解压)。查找 `0a53 21e6` 这个 `ID` 对应的 `.hex` 文件，最后使用 `hex2hcd` 将其转换成 `linux` 下的固件格式，最后将生成的固件按照上面错误里的名字放到 `/lib/firmware` 对应的目录中.


## ` intel/ibt-11-5.sfi`

* 问题

    此设备据反馈再连接 **罗技K370/375 蓝牙键盘** 时配对失败，或是连上就断开，`sudo dmesg` 错误如下：
    
    ```shell
    [    3.962407] Bluetooth: Core ver 2.22
    [    3.962417] Bluetooth: HCI device and connection manager initialized
    [    3.962535] Bluetooth: HCI socket layer initialized
    [    3.962537] Bluetooth: L2CAP socket layer initialized
    [    3.962751] Bluetooth: SCO socket layer initialized
    [    3.973863] Bluetooth: hci0: Bootloader revision 0.0 build 2 week 52 2014
    [    3.980870] Bluetooth: hci0: Device revision is 5
    [    3.980872] Bluetooth: hci0: Secure boot is enabled
    [    3.980873] Bluetooth: hci0: OTP lock is enabled
    [    3.980873] Bluetooth: hci0: API lock is enabled
    [    3.980874] Bluetooth: hci0: Debug lock is disabled
    [    3.980875] Bluetooth: hci0: Minimum firmware build 1 week 10 2014
    [    3.982063] Bluetooth: hci0: Found device firmware: intel/ibt-11-5.sfi
    [    4.325020] Bluetooth: BNEP (Ethernet Emulation) ver 1.3
    [    4.325021] Bluetooth: BNEP filters: protocol multicast
    [    4.325024] Bluetooth: BNEP socket layer initialized
    [    6.151587] Bluetooth: hci0 command 0xfc09 tx timeout
    [   14.312154] Bluetooth: hci0: Failed to send firmware data (-110)
    [   60.134877] Bluetooth: hci0: Bootloader revision 0.0 build 2 week 52 2014
    [   60.141870] Bluetooth: hci0: Device revision is 5
    [   60.141872] Bluetooth: hci0: Secure boot is enabled
    [   60.141872] Bluetooth: hci0: OTP lock is enabled
    [   60.141873] Bluetooth: hci0: API lock is enabled
    [   60.141874] Bluetooth: hci0: Debug lock is disabled
    [   60.141875] Bluetooth: hci0: Minimum firmware build 1 week 10 2014
    [   60.142080] Bluetooth: hci0: Found device firmware: intel/ibt-11-5.sfi
    [   61.604188] Bluetooth: hci0: Waiting for firmware download to complete
    [   61.604887] Bluetooth: hci0: Firmware loaded in 1437668 usecs
    [   61.604951] Bluetooth: hci0: Waiting for device to boot
    [   61.615992] Bluetooth: hci0: Device booted in 10804 usecs
    [   61.616569] Bluetooth: hci0: Found Intel DDC parameters: intel/ibt-11-5.ddc
    [   61.616951] Bluetooth: hci0: Failed to send Intel_Write_DDC (-22)
    [   61.667955] Bluetooth: RFCOMM TTY layer initialized
    [   61.667966] Bluetooth: RFCOMM socket layer initialized
    [   61.667977] Bluetooth: RFCOMM ver 1.11
    ```

* 解决方法

    去下载此电脑在 `window` 下最新版的蓝牙驱动，然后解压(`.exe` 文件使用 `innoextract` 解压)。查找相同 `.sfi` 后缀的文件，然后按照上面错误里的名字放到 `/lib/firmware` 对应的目录中。接着添加 `options iwlwifi bt_coex_active=0` 到 /etc/modprobe.d/iwlwifi.conf 这个文件中(关掉无线蓝牙共存模式，`bt_coex_active:enable wifi/bt co-exist`)。此方法来源与: [罗技K370/375 蓝牙键盘不可用](https://bbs.deepin.org/forum.php?mod=viewthread&tid=148631&page=1#pid403518)
