## [Unreleased]

## [3.2.0] - 2017-10-12
### Added
- Add scale factor setter
- Add touchpad palm setter
- Add 'Timedated' module to reduce authorization times
- Add the timer of detecting filesystem left space
- Add the methods of managing proxychains proxy
- Add the method of refreshing wireless list
- Add 'ClonedAddress' property to indicate current network device mac address


### Changed
- Replace 'xfce/clipboard' with 'gnome/clipboard'
- Refactor 'keybinding' module, replace 'xgb' with 'go-x11-client'
- Update network event notify messages
- Update license
- Reset gesture event state when recieved the end event
- Support to hide apps by modify gsettings
- Support to uninstall 'deepin-fpapp-*' package
- Set the default font style when changing font
- Adjust network widgets layout


### Fixed
- Fix the bug of detecting network device properties error
- Fix activate network hotspot failed
- Fix 'SetProxy' failed if port is empty
