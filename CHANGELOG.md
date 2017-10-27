## [Unreleased]

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
