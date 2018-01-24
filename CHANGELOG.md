## [Unreleased]

## [3.2.9] - 2018-01-24
*   inputdevices: use imwheel to speed up scrolling
*   langselector: use new lib `language_support`
*   dstore: fix waitJobDone for install job
*   swapsched: fix exec cgdelete error
*   keybinding: eliminate keystroke conflict during startup
*   fix: Adapt lintian
*   inputdevices: fix typo in write imwheel config file
*   network: fix nm code generate failure
*   network: add wifi security type 'wpa-eap'
*   inputdevices: fix property WheelSpeed is not writeable
*   network: optimize the method of updating active connections
*   accounts: add DesktopBackgrounds property for user

## [3.2.8] - 2017-12-13
*   add moudle swapsched
*   doc: update bluetooth faq
*   audio: fix update props after config applied
*   dock: fix method RequestDock ignore param index
*   launcher: add methods GetDisableScaling and SetDisableScaling
*   audio: filter out sound effect sink input
*   launcher: fix can not search for newly installed apps
*   appearance: support java scale
*   appearance: fix pam environment settings be override
*   support networkmanager 1.10
*   session/power: set dpms off before suspend
*   makefile GOLDFLAGS remove libcanberra, debian/control depends remove libcanberra-dev

## [3.2.7] - 2017-11-28
*   gesture: check keyboard grab status before do action
*   mime: add multi default app id
*   audio: select best port if config non-exist
*   plymouth: support ssd theme checker
*   dock: fix index in signal EntryAdded is wrong


## [3.2.6] - 2017-11-16
*   add flatpak to recommends

## [3.2.5] - 2017-11-16
*   audio: remove style in font config
*   network: fix wireless disconnect when delete inactive hotspot
*   logined: update 'UserList' when session removed
*   network: remove autoconnect from wireless hotspot
*   appearance: fix fonts memory used large when loading
*   audio: add switcher to decide whether auto switch port


## [3.2.4] - 2017-11-09
#### Features
*   add com.deepin.daemon.ImageBlur interface

#### Bug Fixes
*   not show newly installed wechat in launcher
*   failed to set some bmp image file as icon
*   the Accels field of two shortcuts is empty

#### Changed
*   make `install_to_hicolor.py` compatibility with older python3


## [3.2.3] - 2017-11-03
#### Features
*   automatic switch port when card changed
*   add shortcut for deepin-system-monitor and color-picker
*   support deepin qt theme settings
*   add touchpad tap gesture
*   add flatpak app window identify method


#### Bug Fixes
*   fix gccgo compile failed
*   fix syndaemon pid file not created
*   fix wireless not work after multiple toggle hotspot
*   fix active connections not updated when deleted the last connection
*   update font config xml version


#### Changed
*   refactor grub theme dbus interface
*   rename 'Logout' shortcut to 'Shutdown Interface'
*   add dependency 'dnsmasq'
*   update notifications for scale setting


##  [3.2.2] - 2017-10-27
#### Features
*   keybinding:  process grab pointer failed ([328aa07a](328aa07a))
*   add fprintd module ([1469e2d4](1469e2d4))

#### Bug Fixes
*   fix fprint dependencies missing ([22dc0735](22dc0735))
*   langselector:  write the configuration file wrong ([ee018ea2](ee018ea2))

#### Changed
*   network: remove band settings from hotspot
*   add proxychains-ng as suggested dependency


## [3.2.1] - 2017-10-25
#### Bug Fixes
*   launcher: RequestUninstall does not remove desktop file in autostart directory ([24d1b698](24d1b698))
*   grub2 policykit message not using user's locale ([aa461794](aa461794))
*   keybinding: failed to handle GSettings changed event correctly ([7583b35b](7583b35b))
*   network: delete dot at end ([800eb0c4](800eb0c4))
*   appearance: Fix scale set failed if file not found ([61b72897](61b72897))
*   keybinding can not use key Delete to delete keystroke ([deae5285](deae5285))

#### Features
*   support setting plymouth scale ([842a080e](842a080e))
*   add fprintd module ([1469e2d4](1469e2d4))
*   keybinding: AddCustomShortcut returns id and type of newly created shortcut ([d74f34f8](d74f34f8))
*   accounts: Add no password login ([b87c7448](b87c7448))
*   keybinding: update screenshot command ([64f62269](64f62269))
*   appearance: theme thumbnail support display scaling ([7cba49d6](7cba49d6))
*   dock: menu of entry add item "Force Quit" ([7b853187](7b853187))
*   appearance: Update greeter config when setting scale ([f1b37a80](f1b37a80))
*   network: Implement routes methods ([6889c2d3](6889c2d3))
*   Add 'dde-greeter-setter' ([4dd38e68](4dd38e68))

#### Changed
*   iw: replace 'iw' command with libnl


## [3.2.0] - 2017-10-12
#### Features
* Add scale factor setter
* Add touchpad palm setter
* Add 'Timedated' module to reduce authorization times
* Add the timer of detecting filesystem left space
* Add the methods of managing proxychains proxy
* Add the method of refreshing wireless list
* Add 'ClonedAddress' property to indicate current network device mac address

#### Changed
* Replace 'xfce/clipboard' with 'gnome/clipboard'
* Refactor 'keybinding' module, replace 'xgb' with 'go-x11-client'
* Update network event notify messages
* Update license
* Reset gesture event state when recieved the end event
* Support to hide apps by modify gsettings
* Support to uninstall 'deepin-fpapp-*' package
* Set the default font style when changing font
* Adjust network widgets layout

#### Bug Fixes
* Fix the bug of detecting network device properties error
* Fix activate network hotspot failed
* Fix 'SetProxy' failed if port is empty
