<a name=""></a>
##  4.6.4 (2018-05-24)


#### Features

*   use gsettings value delay load plugins ([3ac07aed](3ac07aed))
*   the maximum volume is 150 ([154df257](154df257))
*   add keybord_layout.json ([e4faf0b5](e4faf0b5))
*   add keyboard layout plugin. ([3d19f67c](3d19f67c))
*   support indicator traywidget ([58fd2df6](58fd2df6))
*   open dde-calendar with dbus call ([fb6ee52b](fb6ee52b))
* **keyboard_layout:**  add new dbus interface ([7ef2fea7](7ef2fea7))
* **network:**
  *  add refresh animation ([5bb7b22e](5bb7b22e))
  *  support cloned address ([e2983ded](e2983ded))
* **plugins:**  keep order ([e281f088](e281f088))
* **system-tray:**  using native popup ([354cfd9f](354cfd9f))

#### Bug Fixes

*   change display mode no signal. ([79922e89](79922e89))
*   identify launcher icon. ([a6b87571](a6b87571))
*   cookie type error. ([8d0bdbf6](8d0bdbf6))
*   geometry error. ([67250099](67250099))
*   plugin item not free. ([94fc261e](94fc261e))
*   plugin can not popup content window. ([a3c84e3c](a3c84e3c))
*   swing effect not allow mouse event ([61002cd3](61002cd3))
*   app volume can be adjusted to 150 ([acd2bd0d](acd2bd0d))
*   call dbus error ([d651fc9d](d651fc9d))
*   Popup Applet not hide when item removed ([0e9d9df0](0e9d9df0))
*   popup not shown ([0a0b60aa](0a0b60aa))
*   1px line on top and left ([65100785](65100785))
*   not installed keyboard_layout.json ([3776f361](3776f361))
*   dock hang when receive attention on fastion mode ([1444fa40](1444fa40))
*   sound input sink volume slider not init ([e06a0ac0](e06a0ac0))
*   update trayWidget on fashion mode ([ecb014c9](ecb014c9))
*   hidpi support for indicatortraywidget ([2704bf62](2704bf62))
*   indicatortraywidget connect error pertoires slots ([83cb34e8](83cb34e8))
*   Adapt lintian ([182ba8bd](182ba8bd))
*   dock hide behavior error when popup menu ([ad3f979e](ad3f979e))
*   menu auto hide ([51ee4140](51ee4140))
*   positioned at wrong place ([842e530e](842e530e))
*   show control-center not work ([2774afe4](2774afe4))
*   container item icon not clearly ([7b817246](7b817246))
*   dont show container tips if it's empty ([7ca10744](7ca10744))
*   aliasing border on highlighting app icon ([093d87e9](093d87e9))
*   dock frontend rect not support hiDPI ([9c6c41cb](9c6c41cb))
*   sound slider use svg ([a84062bc](a84062bc))
*   network refresh button support hidpi ([c784f3c5](c784f3c5))
*   network refresh button not visiable when device on ([3807c731](3807c731))
* **WiFi:**  connecting indicator icon is too large ([a005b64f](a005b64f))
* **attention:**  crash when app item destory. ([d40239a3](d40239a3))
* **container:**  item tips color wrong ([e9803aa2](e9803aa2))
* **dockitem:**  popup applet position error ([ec1ca53e](ec1ca53e))
* **item:**  panel will hide when show menu on keep always hide mode ([c29bca64](c29bca64))
* **network:**
  *  refresh button not set size ([05015516](05015516))
  *  add new wireless disconnect svg ([0421d796](0421d796))
  *  network plugin icon diable using theme icon ([e81ba0ba](e81ba0ba))
  *  plugins icon not support hidpi ([f5b69302](f5b69302))
  *  ap state error when hotspot open ([62eae043](62eae043))
* **plugin:**  wireless not refresh ([bcc5e866](bcc5e866))
* **power:**  power icon shown if switch display mode ([4ff28712](4ff28712))
* **sound:**
  *  resources not load ([7e94385a](7e94385a))
  *  app mute icon is not support HIDPI ([b5f8253a](b5f8253a))
  *  fallback to default icon ([3e7e65d7](3e7e65d7))
  *  plugins config file locate error ([34ce1bdd](34ce1bdd))
* **sys-tray:**  system tray alignment adjust ([9018711b](9018711b))
* **sys_tray:**  system tray not align vertial center ([a545ce08](a545ce08))
* **visible:**  visible error when display mode changed ([8ed399a5](8ed399a5))
* **wireless:**
  *  init device enabled state ([ca1c5acd](ca1c5acd))
  *  call dbus frequently ([4d217516](4d217516))

#### Performance

*   do not touch graphic scene when it's not necessary ([373fe443](373fe443))



<a name="4.6.3"></a>
### 4.6.3 (2018-03-28)


#### Bug Fixes

*   Popup Applet not hide when item removed ([0e9d9df0](0e9d9df0))
*   popup not shown ([0a0b60aa](0a0b60aa))



<a name=""></a>
##  4.6.2 (2018-03-22)


#### Bug Fixes

*   1px line on top and left ([65100785](65100785))
*   not installed keyboard_layout.json ([3776f361](3776f361))
*   dock hang when receive attention on fastion mode ([1444fa40](1444fa40))




<a name="4.6.1"></a>
## 4.6.1 (2018-03-15)


#### Bug Fixes

*   sound input sink volume slider not init ([e06a0ac0](e06a0ac0))
*   update trayWidget on fashion mode ([ecb014c9](ecb014c9))
*   hidpi support for indicatortraywidget ([2704bf62](2704bf62))
*   indicatortraywidget connect error pertoires slots ([83cb34e8](83cb34e8))
*   Adapt lintian ([182ba8bd](182ba8bd))
*   dock hide behavior error when popup menu ([ad3f979e](ad3f979e))
*   menu auto hide ([51ee4140](51ee4140))
*   positioned at wrong place ([842e530e](842e530e))
*   show control-center not work ([2774afe4](2774afe4))
*   container item icon not clearly ([7b817246](7b817246))
*   dont show container tips if it's empty ([7ca10744](7ca10744))
*   aliasing border on highlighting app icon ([093d87e9](093d87e9))
*   dock frontend rect not support hiDPI ([9c6c41cb](9c6c41cb))
*   sound slider use svg ([a84062bc](a84062bc))
*   network refresh button support hidpi ([c784f3c5](c784f3c5))
*   network refresh button not visiable when device on ([3807c731](3807c731))
* **WiFi:**  connecting indicator icon is too large ([a005b64f](a005b64f))
* **attention:**  crash when app item destory. ([d40239a3](d40239a3))
* **container:**  item tips color wrong ([e9803aa2](e9803aa2))
* **dockitem:**  popup applet position error ([ec1ca53e](ec1ca53e))
* **item:**  panel will hide when show menu on keep always hide mode ([c29bca64](c29bca64))
* **network:**
  *  refresh button not set size ([05015516](05015516))
  *  add new wireless disconnect svg ([0421d796](0421d796))
  *  network plugin icon diable using theme icon ([e81ba0ba](e81ba0ba))
  *  plugins icon not support hidpi ([f5b69302](f5b69302))
  *  ap state error when hotspot open ([62eae043](62eae043))
* **plugin:**  wireless not refresh ([bcc5e866](bcc5e866))
* **power:**  power icon shown if switch display mode ([4ff28712](4ff28712))
* **sound:**
  *  resources not load ([7e94385a](7e94385a))
  *  app mute icon is not support HIDPI ([b5f8253a](b5f8253a))
  *  fallback to default icon ([3e7e65d7](3e7e65d7))
  *  plugins config file locate error ([34ce1bdd](34ce1bdd))
* **sys-tray:**  system tray alignment adjust ([9018711b](9018711b))
* **sys_tray:**  system tray not align vertial center ([a545ce08](a545ce08))
* **visible:**  visible error when display mode changed ([8ed399a5](8ed399a5))
* **wireless:**
  *  init device enabled state ([ca1c5acd](ca1c5acd))
  *  call dbus frequently ([4d217516](4d217516))

#### Features

*   add keybord_layout.json ([e4faf0b5](e4faf0b5))
*   add keyboard layout plugin. ([3d19f67c](3d19f67c))
*   support indicator traywidget ([58fd2df6](58fd2df6))
*   open dde-calendar with dbus call ([fb6ee52b](fb6ee52b))
* **keyboard_layout:**  add new dbus interface ([7ef2fea7](7ef2fea7))
* **network:**
  *  add refresh animation ([5bb7b22e](5bb7b22e))
  *  support cloned address ([e2983ded](e2983ded))
* **plugins:**  keep order ([e281f088](e281f088))
* **system-tray:**  using native popup ([354cfd9f](354cfd9f))



<a name="4.6.0"></a>
## 4.6.0 (2018-03-12)


#### Features

*   add keyboard layout plugin. ([3d19f67c](3d19f67c))
*   support indicator traywidget ([58fd2df6](58fd2df6))
*   open dde-calendar with dbus call ([fb6ee52b](fb6ee52b))
* **network:**
  *  add refresh animation ([5bb7b22e](5bb7b22e))
  *  support cloned address ([e2983ded](e2983ded))
* **plugins:**  keep order ([e281f088](e281f088))
* **system-tray:**  using native popup ([354cfd9f](354cfd9f))

#### Bug Fixes

*   sound input sink volume slider not init ([e06a0ac0](e06a0ac0))
*   update trayWidget on fashion mode ([ecb014c9](ecb014c9))
*   hidpi support for indicatortraywidget ([2704bf62](2704bf62))
*   indicatortraywidget connect error pertoires slots ([83cb34e8](83cb34e8))
*   Adapt lintian ([182ba8bd](182ba8bd))
*   dock hide behavior error when popup menu ([ad3f979e](ad3f979e))
*   menu auto hide ([51ee4140](51ee4140))
*   positioned at wrong place ([842e530e](842e530e))
*   show control-center not work ([2774afe4](2774afe4))
*   container item icon not clearly ([7b817246](7b817246))
*   dont show container tips if it's empty ([7ca10744](7ca10744))
*   aliasing border on highlighting app icon ([093d87e9](093d87e9))
*   dock frontend rect not support hiDPI ([9c6c41cb](9c6c41cb))
*   sound slider use svg ([a84062bc](a84062bc))
*   network refresh button support hidpi ([c784f3c5](c784f3c5))
*   network refresh button not visiable when device on ([3807c731](3807c731))
* **WiFi:**  connecting indicator icon is too large ([a005b64f](a005b64f))
* **attention:**  crash when app item destory. ([d40239a3](d40239a3))
* **container:**  item tips color wrong ([e9803aa2](e9803aa2))
* **item:**  panel will hide when show menu on keep always hide mode ([c29bca64](c29bca64))
* **network:**
  *  refresh button not set size ([05015516](05015516))
  *  add new wireless disconnect svg ([0421d796](0421d796))
  *  network plugin icon diable using theme icon ([e81ba0ba](e81ba0ba))
  *  plugins icon not support hidpi ([f5b69302](f5b69302))
  *  ap state error when hotspot open ([62eae043](62eae043))
* **plugin:**  wireless not refresh ([bcc5e866](bcc5e866))
* **power:**  power icon shown if switch display mode ([4ff28712](4ff28712))
* **sound:**
  *  resources not load ([7e94385a](7e94385a))
  *  app mute icon is not support HIDPI ([b5f8253a](b5f8253a))
  *  fallback to default icon ([3e7e65d7](3e7e65d7))
  *  plugins config file locate error ([34ce1bdd](34ce1bdd))
* **sys_tray:**  system tray not align vertial center ([a545ce08](a545ce08))
* **visible:**  visible error when display mode changed ([8ed399a5](8ed399a5))



# 4.5.12
    Fix item not update when wm changed

# 4.5.11
    UI improve

# 4.5.10
    Minor bug fixes.

# 4.5.9
    Fix frontend rect error.

# 4.5.8
    Improve dock startup animation.

# 4.5.7
    Using system theme icon in network plugin

# 4.5.6
    Fix dock position error in multi screen

# 4.5.5
    Forbid using theme icon in network plugin.

# 4.5.4
    Fix dock crash if system tray' app is quit.

# 4.5.3
    Minor bug fixes.

# 4.5.2
    Improve popup auto hide algorithm

# 4.5.1
    Fix window size error when new item inserted.

# 4.5.0
    Improve HiDPI issues.
    Opmitize memory usage.

# 4.4.3
    Improve user experience.
    Memory issue fixes.

# 4.4.2
    Network add connecting animation
    Fix systemtray not align vertial center

# 4.4.1
    Update translations
    Fix power icon auto-hide not work perfectly

# 4.4.0
    New plugin API.
    HiDPI support.
