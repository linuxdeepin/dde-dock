<a name="4.8.4"></a>
## 4.8.4 (2018-12-14)


#### Bug Fixes

*   context menu invoke in touch screen ([77570b60](https://github.com/linuxdeepin/dde-dock/commit/77570b609a9666025b626cf5c98dc4024b9e9c17))
*   tips of network item not refresh ([635d4152](https://github.com/linuxdeepin/dde-dock/commit/635d41525f3390dfc09e2bdef28b54214bed5d5a))
*   sni tray context menu can not be initialize correctly after reboot ([8870cbc7](https://github.com/linuxdeepin/dde-dock/commit/8870cbc799b88f4350cd5ed5d882648c3d5c7496))
*   threshold to emit drag event ([06ca6df2](https://github.com/linuxdeepin/dde-dock/commit/06ca6df2057861fb70fca2ed212c8d41f458a601))



<a name="4.8.3.1"></a>
## 4.8.3.1 (2018-12-10)


#### Bug Fixes

*   compile failed in debian stable repo ([fc64302b](https://github.com/linuxdeepin/dde-dock/commit/fc64302b9e2b148faf3f66afa47e70a222633f23))
*   touch screen drag/drop tray icons to tray container by mistake ([0a3eb2b6](https://github.com/linuxdeepin/dde-dock/commit/0a3eb2b61922815b3cfe14fdcd8ee7faa4aa0534))
*   wired tray icon in hihdpi ([a9998b6c](https://github.com/linuxdeepin/dde-dock/commit/a9998b6cf2acf01fb39394ca3a23d2b278c6734e))
*   tray icon size in hidpi ([5b5d55d9](https://github.com/linuxdeepin/dde-dock/commit/5b5d55d9dbeb730e2636e62b62cb0b3b5dd6b6b9))
*   tray hover state not refresh after some mouse event ([4d8e0772](https://github.com/linuxdeepin/dde-dock/commit/4d8e077243bba234bd154b18b18dd874ad928cdd))



<a name="4.8.3"></a>
## 4.8.3 (2018-12-07)


#### Bug Fixes

*   tray icon do not change after system icon theme has changed ([16f10b66](https://github.com/linuxdeepin/dde-dock/commit/16f10b665c0136a984be882279494e3271a252be))
*   hover background while dragging ([3ee25e36](https://github.com/linuxdeepin/dde-dock/commit/3ee25e36b1d27a7b9ec69b78ef40738d2b8293f4))
*   sni tray icon size in hihdpi ([74cd9594](https://github.com/linuxdeepin/dde-dock/commit/74cd9594f9ae9a045ebe85d26cb2f0a26299c477))
*   sometimes got a invalid dbus menu path of some application sni tray ([717b30bc](https://github.com/linuxdeepin/dde-dock/commit/717b30bcec3883f3bd841c36f8e93a1654eab002))
*   build warning 2 ([dc1d415d](https://github.com/linuxdeepin/dde-dock/commit/dc1d415dc34ab1d3b4a02f0c3328392ad31820fc))
*   build warning ([61197b9d](https://github.com/linuxdeepin/dde-dock/commit/61197b9d3d9028a062edb72c1487b214600cceda))
*   typo2 ([63dc1029](https://github.com/linuxdeepin/dde-dock/commit/63dc102939072cfbf6a364498c1cbe55e5f13d0e))
*   typo2 ([9c372657](https://github.com/linuxdeepin/dde-dock/commit/9c3726571302a40f36e6db38e56d63402a8bc137))
*   typo ([a87911ce](https://github.com/linuxdeepin/dde-dock/commit/a87911ce8051dcd4ad74a5a8e379ba1e10c63f50))
*   dock hide problem and remove old imp ([091b52bc](https://github.com/linuxdeepin/dde-dock/commit/091b52bccddaa2ea9795b4542aca892d7a27c849))
* **TrayPlugin:**  send hover event to XEmbed trays ([dfa4dd9d](https://github.com/linuxdeepin/dde-dock/commit/dfa4dd9d24afa9577d92c78adeceb38bec529a8c))

#### Features

*   tray icon hover state ([caede05f](https://github.com/linuxdeepin/dde-dock/commit/caede05fe7bb395e080299e58ca4e5e2e4528205))



<a name="4.8.2.2"></a>
## 4.8.2.2 (2018-12-04)


#### Features

*   cherry-pick "561ed99 feat: more detail network status" ([c079511c](https://github.com/linuxdeepin/dde-dock/commit/c079511c6947ee9b1d00177b572720501cf43abf))



<a name="4.8.2.1"></a>
## 4.8.2.1 (2018-12-04)




<a name="4.8.2"></a>
## 4.8.2 (2018-12-03)


#### Features

*   more detail network status ([a84ef9d4](https://github.com/linuxdeepin/dde-dock/commit/a84ef9d4a8614938d3bc0f4bc96f39d674fa7b49))
*   save sort order to config and optimize sort algorithm ([805b5b56](https://github.com/linuxdeepin/dde-dock/commit/805b5b56e06f7431bfaefca4c55ff063502fadc7))
*   support drag and swap position ([b928e6fb](https://github.com/linuxdeepin/dde-dock/commit/b928e6fb4d919deba64086f384057ab932255ba1))
*   support drag fashion tray icon ([12434288](https://github.com/linuxdeepin/dde-dock/commit/124342881786a789f2e86bba5ea2393538958fb6))
*   fashion trays default sort order ([0f7a583d](https://github.com/linuxdeepin/dde-dock/commit/0f7a583d081f1f2c8a98ce9825e267ad9fce164f))
*   support item sort key interface of TrayPlugin inner plugin ([2ebef72a](https://github.com/linuxdeepin/dde-dock/commit/2ebef72a1845da199927a5cf84fbc1cbcde2faef))

#### Bug Fixes

*   remove unimplemented slots ([a244ea56](https://github.com/linuxdeepin/dde-dock/commit/a244ea567a7a9d60a75c0102c5541e3e7a15664d))
*   save fashion tray item  expand config in public config file ([a6e2546e](https://github.com/linuxdeepin/dde-dock/commit/a6e2546efcd2d55d67960e7413852e02cd5b29c3))
*   click item is ignored ([5e9886fa](https://github.com/linuxdeepin/dde-dock/commit/5e9886fa97ab413d32b22f47134029fbe492f9c8))
*   play swing effect when open app with mouse middle button ([fd3f5025](https://github.com/linuxdeepin/dde-dock/commit/fd3f502529a6f6bba14daa9bd9357fa04df44a6a))
*   connect to signal failed ([1dba9484](https://github.com/linuxdeepin/dde-dock/commit/1dba9484f2e64300874565d915f9242e30662bfa))
*   click on trays of wine app ([10da46b6](https://github.com/linuxdeepin/dde-dock/commit/10da46b64b07c7aacf22bab1daf93f0028a8df39))
*   handle tray mouseMoveEvent ([e025070a](https://github.com/linuxdeepin/dde-dock/commit/e025070ac557bacdf8c38471ead9b94c16787c0d))
*   forbidden lable when darg tray and move in fashion tray ([77f75382](https://github.com/linuxdeepin/dde-dock/commit/77f75382fdd0aa2958169c4fbb7139017749f7c8))
*   method of drag and swap item ([72a3d2d5](https://github.com/linuxdeepin/dde-dock/commit/72a3d2d52321f7987b81271cb5276ac0ef0e0ed5))
*   find the dest insert index ([83673322](https://github.com/linuxdeepin/dde-dock/commit/83673322caaa9f48b2e7cef43d1a89772e71313e))
*   tray strange fold animation when dock size is maxied ([7c76f8bf](https://github.com/linuxdeepin/dde-dock/commit/7c76f8bf4c580efab23d10eae79089d49eadc19e))
*   the end value of window size animation is error ([93601e9a](https://github.com/linuxdeepin/dde-dock/commit/93601e9aeb7e6a0f788acbc9058d74fc3b6267f2))
*   system tray icon still be shown while dragging ([7833bc93](https://github.com/linuxdeepin/dde-dock/commit/7833bc9344550a618872418d343927f8564c3b6c))
*   found removed fashion tray failed ([1272d9f9](https://github.com/linuxdeepin/dde-dock/commit/1272d9f9223663c0ddc80b5baf25adaacdd610e7))
*   crash after remove tray item from fashion tray ([6e883538](https://github.com/linuxdeepin/dde-dock/commit/6e883538091e39522a27e8ddc38f3c16205bd5c9))



<a name="4.8.1"></a>
## 4.8.1 (2018-11-23)


#### Bug Fixes

*   var name typo ([2304a7f5](https://github.com/linuxdeepin/dde-dock/commit/2304a7f506f619403d48d3f288314d034775145c))
*   tray strange fold animation when dock size is maxied ([57f5e5c7](https://github.com/linuxdeepin/dde-dock/commit/57f5e5c75cbc949b09380eed1bff387584e092e5))
*   the end value of window size animation is error ([7354cf51](https://github.com/linuxdeepin/dde-dock/commit/7354cf5166c658983d97615f02ddbc380b5aa9b0))
*   system tray icon still be shown while dragging ([9bfb8eee](https://github.com/linuxdeepin/dde-dock/commit/9bfb8eee372811355140dc49e8be7a634d4e452b))
*   drop to container item careless ([c61b620d](https://github.com/linuxdeepin/dde-dock/commit/c61b620debfb7fa1055eb3b58e3c9ebf803b733c))
*   date-time plugin enable status ([2d925fe7](https://github.com/linuxdeepin/dde-dock/commit/2d925fe74ed1fb2a62b5a857a1a5c6d472320bf8))
*   tray icon pixmap align ([d1fa5364](https://github.com/linuxdeepin/dde-dock/commit/d1fa5364065a4da924a4a2fe7346952706743162))
*   dock resize do not in time ([efd8e01e](https://github.com/linuxdeepin/dde-dock/commit/efd8e01e6ac508167d3110fcbb393db02d31109e))
*   invalid tray size when change display mode ([c7f953e1](https://github.com/linuxdeepin/dde-dock/commit/c7f953e121a600b24560b27718001628a34a7c69))
*   some types ([074e0ba4](https://github.com/linuxdeepin/dde-dock/commit/074e0ba4dbc2a5c181ba2ad4a952b7dac98131cb))
*   invalid icon size after dock size changed ([8dc212bc](https://github.com/linuxdeepin/dde-dock/commit/8dc212bc3a69913a1aa93d95b48dc8ca7e7336e4))
*   dock crash while loading plugins ([a06b7f9d](https://github.com/linuxdeepin/dde-dock/commit/a06b7f9dcd9ffcd03122b9007208300590afbc5a))

#### Features

*   touchscreen support ([ca085678](https://github.com/linuxdeepin/dde-dock/commit/ca085678617a5d4319924ecd0ba88b9fe6fe99a7))
*   optimiza size change animation of dock ([1a4652ca](https://github.com/linuxdeepin/dde-dock/commit/1a4652cab713d4824179d6d1bac0ca4c09675920))
* **system-tray:**  add animation for system tray expand and fold ([8090ef44](https://github.com/linuxdeepin/dde-dock/commit/8090ef445ebdd3653382540cf2858704626d3322))



<a name="4.8.0"></a>
## 4.8.0 (2018-11-13)




<a name="4.7.9"></a>
## 4.7.9 (2018-11-12)


#### Bug Fixes

*   error value of decrease fashion icon size ([92ac6dc3](https://github.com/linuxdeepin/dde-dock/commit/92ac6dc377fd49abc57db494516160a123b0e235))



<a name="4.7.8"></a>
## 4.7.8 (2018-11-12)


#### Bug Fixes

*   resize dock and fashion system tray recursively ([738f41aa](https://github.com/linuxdeepin/dde-dock/commit/738f41aa1728dacc2e1f94688c6f6a41c435e822))



<a name="4.7.7"></a>
### 4.7.7 (2018-11-09)


#### Bug Fixes

*   change min icon size ([8af71fae](https://github.com/linuxdeepin/dde-dock/commit/8af71faef3440ac18c1bb8c0c8a1141ca1b30378))



<a name="4.7.6"></a>
## 4.7.6 (2018-11-09)


#### Features

*   integrating plugins config files ([4837c9dd](https://github.com/linuxdeepin/dde-dock/commit/4837c9dd35e23ff8a08d063bc7a57fc03086ae20))
*   prevent unexpected trigger ([7ea3341a](https://github.com/linuxdeepin/dde-dock/commit/7ea3341aca124cf4f67cfced4acc3c04c6c7fca4))
*   auto adjust the icon size of fashion system tray ([598978d3](https://github.com/linuxdeepin/dde-dock/commit/598978d35043695e85e36e8ba70fb62ddf00061b))
*   support dde-dock.pc ([71b237bf](https://github.com/linuxdeepin/dde-dock/commit/71b237bfb9f4077e9158d0b0dd68dbbb42f62173))

#### Bug Fixes

*   not change dock visible in time ([93db674b](https://github.com/linuxdeepin/dde-dock/commit/93db674bb2ed5ead4527e273a4b85a10f0054c1d))
*   active connection info of device may not set ([bddcad5e](https://github.com/linuxdeepin/dde-dock/commit/bddcad5e6b92aae06b32921dd8114f93df850c93))
*   left click on sogou tray ([213db421](https://github.com/linuxdeepin/dde-dock/commit/213db42172176826f1f6d24238d51f79c3182f26))
*   backgroud of appitem ([949849af](https://github.com/linuxdeepin/dde-dock/commit/949849af62e2a0b6ba8b0d9bfe594329c1b757ec))
*   secret pixmap of wireless ap not be refresh ([084477f2](https://github.com/linuxdeepin/dde-dock/commit/084477f2082510bbc1a353c9199c4815d01bb5ad))
*   can not show context menu on system-tray control button ([bffcc3a1](https://github.com/linuxdeepin/dde-dock/commit/bffcc3a18573f252d8c0cc7711933518457dc517))
*   set the default theme to "dark" ([714c3ff1](https://github.com/linuxdeepin/dde-dock/commit/714c3ff1fa1e549cf85ac42fb45516494a23bb84))
*   duplicate include file ([0ff728b8](https://github.com/linuxdeepin/dde-dock/commit/0ff728b83f8ac5c7f053d17754c3f4cc8382c602))
*   non secret wireless align in hihdpi ([ce9dad99](https://github.com/linuxdeepin/dde-dock/commit/ce9dad99929cfb29028e1c389193b3f757330483))
*   place holder item still display in dock after drag leave ([725891a6](https://github.com/linuxdeepin/dde-dock/commit/725891a6b1cf9e141505d45e4280671f6debf2e4))
* **system-tray:**  new tray need attention policy ([cf5e0bc2](https://github.com/linuxdeepin/dde-dock/commit/cf5e0bc2d552303afbe9dc53e5af617c475984d6))



<a name="4.7.5.2"></a>
## 4.7.5.2 (2018-11-02)


#### Bug Fixes

*   dock crash by shutdown plugin config file ([e7677b6f](https://github.com/linuxdeepin/dde-dock/commit/e7677b6f9141b071037334fbf30ca93825259f98))



<a name="4.7.5.1"></a>
## 4.7.5.1 (2018-11-02)


#### Bug Fixes

*   translation of sound system tray load failed ([3ff8f35f](https://github.com/linuxdeepin/dde-dock/commit/3ff8f35ff49c81e7c838c65c13215b01fed88c43))



<a name="4.7.5"></a>
## 4.7.5 (2018-11-01)


#### Bug Fixes

*   the context menu can not be shown when container item is clicked ([22119b98](https://github.com/linuxdeepin/dde-dock/commit/22119b989a2466ef0223c832f679238b0d4c185f))
*   crash when drag desktop file leave ([ade41a05](https://github.com/linuxdeepin/dde-dock/commit/ade41a05056c746785f0f06d2bf326508e1d9418))
*   the dock is hidden when the tray area is expanded/closed ([ef6b78cc](https://github.com/linuxdeepin/dde-dock/commit/ef6b78ccbfac578476faf1b09b73672cc9062a34))
*   postion error of tip and applet of system trays ([26328147](https://github.com/linuxdeepin/dde-dock/commit/263281478f2a2ad89ed262e6e5745a6d3dd78ab3))
* **icon:**  limit min icon size to 24 ([297e0b57](https://github.com/linuxdeepin/dde-dock/commit/297e0b57f228afbadd3eb8b86de2fb44638d7a6c))
* **system-tray:**
  *  crash when refresh wired tray visible ([5c042701](https://github.com/linuxdeepin/dde-dock/commit/5c042701e1fbaafaef04013f5b1dcea88bd5342e))
  *  system tray icon tips and applet position ([d94da603](https://github.com/linuxdeepin/dde-dock/commit/d94da6033fa7b2c2d0ca99bf11b1140b702ccfcc))
  *  dock visible when context menu and tips of system tray is shown ([6ffd34ef](https://github.com/linuxdeepin/dde-dock/commit/6ffd34ef337b3ecdde827a495a8b7106f82e5a73))
  *  control widget icon hihdpi ([93f0ee27](https://github.com/linuxdeepin/dde-dock/commit/93f0ee27479ac64983bbe1111ded1856cfb42a63))
  *  system trays should not be drag into container ([d4641059](https://github.com/linuxdeepin/dde-dock/commit/d464105975b7761b24e68c11fe919fb4ab927f31))
  *  should resize system-tray after tray inactived ([cf72371f](https://github.com/linuxdeepin/dde-dock/commit/cf72371fd91d0976c6544318afcfbb33ed460f0a))
  *  icons position in hihdpi ([bd4bdfda](https://github.com/linuxdeepin/dde-dock/commit/bd4bdfdaebbd7df6518220d8f2248bf1d9ed9812))



<a name="4.7.4"></a>
## 4.7.4 (2018-10-26)


#### Bug Fixes

*   active connection identify ([62f381d6](https://github.com/linuxdeepin/dde-dock/commit/62f381d6ba08c6a1b5f72862bdac3ca4d549ed20))
*   frame opacity dbus valid ([e342a189](https://github.com/linuxdeepin/dde-dock/commit/e342a189535e3382144b1789f03d5f4ad984e860))
* **launcher:**  use show replace toggle for launcher dbus ([203a6eb9](https://github.com/linuxdeepin/dde-dock/commit/203a6eb96f2da63fe0370ed87aa74e02b7f4c0c5))
* **network:**
  *  loading indicator visible ([6ead8ff3](https://github.com/linuxdeepin/dde-dock/commit/6ead8ff3fa2bd3c8fe4df5f3a28ce7848f7d6a5c))
  *  wired item visiable when network plugin disabled ([b6a2a615](https://github.com/linuxdeepin/dde-dock/commit/b6a2a6152b0785a288f0025bad21bd72264935c5))
* **system-tray:**
  *  battery icon invisible ([fbff7b1f](https://github.com/linuxdeepin/dde-dock/commit/fbff7b1fff03a96a99e9b65733a548cd2acfa343))
  *  keyboard indicator not show ([99095cb2](https://github.com/linuxdeepin/dde-dock/commit/99095cb249b2a502e2775d4512a38eea912646fc))

#### Features

*   support change frame opacity ([0310a313](https://github.com/linuxdeepin/dde-dock/commit/0310a31352ed3fb2972ead029b42218460169cea))
* **system-tray:**  check have new tray protocol menu ([11d97f20](https://github.com/linuxdeepin/dde-dock/commit/11d97f206e4d50e566a3d66c3071eb7ec555ab92))



<a name="4.7.3.5"></a>
## 4.7.3.5 (2018-10-12)




<a name="4.7.3.4"></a>
## 4.7.3.4 (2018-10-12)


#### Bug Fixes

*   query active info when hotspot enabled ([4b557ae9](https://github.com/linuxdeepin/dde-dock/commit/4b557ae9a2d4e2b09dae4c655cb5b435eab3e191))



<a name="4.7.3.3"></a>
## 4.7.3.3 (2018-09-21)


#### Bug Fixes

*   pixmap size of wireless password dialog in hidpi mode ([f73428a4](https://github.com/linuxdeepin/dde-dock/commit/f73428a465a235ba04db62909fa3f5bd10265461))
* **network:**  wired item visiable when network plugin disabled ([fc42fd94](https://github.com/linuxdeepin/dde-dock/commit/fc42fd949c1be2d41731fd0904cd17329a053341))



## 4.7.3.2 (2018-09-13)
<a name="4.7.2.3"></a>
## 4.7.2.3 (2018-09-10)


#### Bug Fixes

*   app icon find mistake ([14894636](https://github.com/linuxdeepin/dde-dock/commit/14894636ab1e6d812cb3ab936e23b28af94f127e))
*   tips color error under 2D ([3ff02944](https://github.com/linuxdeepin/dde-dock/commit/3ff029440a4091cac0a9721e4091370bab1f44c0))
* **network:**  wired item visible logic when it's enable status changed ([f5da9e50](https://github.com/linuxdeepin/dde-dock/commit/f5da9e508cef1d1838cfa2ae151c660b1cace3f2))



<a name="4.7.2.2"></a>



*   apps preview position error ([14e1cd49](https://github.com/linuxdeepin/dde-dock/commit/14e1cd49edc990bce1e73a74d5011b19c40882fa))



<a name="4.7.3.1"></a>
## 4.7.3.1 (2018-09-10)


#### Features

*   remove spacing for launcher item ([f17cce91](https://github.com/linuxdeepin/dde-dock/commit/f17cce91d62bee357a657fd0c32bf993c4879ebb))

#### Bug Fixes

*   app icon find mistake ([7ae40f4b](https://github.com/linuxdeepin/dde-dock/commit/7ae40f4bb9ef23b9b6595e456a8428a81315bd8c))
*   tips color error under 2D ([009c03c9](https://github.com/linuxdeepin/dde-dock/commit/009c03c9fe70c80cdab77cde1091063951e92a3a))
* **network:**  wired item visible logic when it's enable status changed ([821e3937](https://github.com/linuxdeepin/dde-dock/commit/821e393753fdab98428092ef025bcd231a011e69))
* **plugins:**  use a fixed default sort for plugins ([e802e87e](https://github.com/linuxdeepin/dde-dock/commit/e802e87e4b743deba4c50213e0b302d12adda032))



<a name="4.7.3"></a>
## 4.7.3 (2018-08-31)
## 4.7.2.2 (2018-08-31)


#### Bug Fixes

*   tips size error ([48113767](https://github.com/linuxdeepin/dde-dock/commit/481137678b52c9f1e7b4aada9653863e813259c4))
*   strange structure when dock is start up ([752602d1](https://github.com/linuxdeepin/dde-dock/commit/752602d156d3c41596ff66c7c8663841b2d7ff97))
*   network pointer is not initialized ([da6d47f4](https://github.com/linuxdeepin/dde-dock/commit/da6d47f47b0af37af4db7a8322f386e65be94d6e))
*   icon is null ([da407f77](https://github.com/linuxdeepin/dde-dock/commit/da407f77dbf34b232ea3b50c51b094095206deef))
*   activated wireless info ([c287f7d6](https://github.com/linuxdeepin/dde-dock/commit/c287f7d6d72e1da8d713e19b8d91be07b7f2ac1f))
* **app:**  icon is null ([79f85f3d](https://github.com/linuxdeepin/dde-dock/commit/79f85f3d4af7ac71255a99f4df70db36cb1a3e28))
* **network:**  activated wireless info ([5f4df46a](https://github.com/linuxdeepin/dde-dock/commit/5f4df46a7626e041c9e81ad81fe9cab410929647))



<a name=""></a>
##  4.7.2.1 (2018-08-28)


#### Bug Fixes

*   strange structure when dock is start up ([d6532d2a](https://github.com/linuxdeepin/dde-dock/commit/d6532d2aa690c389a468c532fcb255ac84bb1ddd))
*   network pointer is not initialized ([5a63377f](https://github.com/linuxdeepin/dde-dock/commit/5a63377f870f2d91be04571b2b542968f131d501))



<a name=""></a>
##  4.7.2 (2018-08-12)


#### Bug Fixes

* **preview:**  invalid shm permission ([ac1fca3e](https://github.com/linuxdeepin/dde-dock/commit/ac1fca3e131b5590253120d18df4fdbb938283e1))



<a name=""></a>
##  4.7.1.1 (2018-08-09)




<a name="4.7.1"></a>
### 4.7.1 (2018-08-08)


#### Bug Fixes

*   drag widget follow the mouse all the time ([54593f55](https://github.com/linuxdeepin/dde-dock/commit/54593f553874a0026a4b68740e706eec28049dc2))



<a name="4.7.0"></a>
### 4.7.0 (2018-08-07)


#### Bug Fixes

*   triple dock size drag to remove distance ([87e6d18a](https://github.com/linuxdeepin/dde-dock/commit/87e6d18aaf99aa2c9b67643daae4881de663a97e))
* **network:**
  *  wired item visible ([23806991](https://github.com/linuxdeepin/dde-dock/commit/238069913130e492f00bea6edc67377e4cf6127b))
  *  refresh network plug theme icon ([c115fd15](https://github.com/linuxdeepin/dde-dock/commit/c115fd153fd110938eda2c541f0d7c96af77d033))

#### Features

*   add get preview image from shm ([081522f0](https://github.com/linuxdeepin/dde-dock/commit/081522f02ccbb628c75d8ae1e5c09da33d0ff5a2))
*   lazy loading of plugins which depends dbus daemon ([95b5c72f](https://github.com/linuxdeepin/dde-dock/commit/95b5c72f13e0ae35fae0c2b2a698bc567ae351ed))
 


<a name="4.6.9"></a>
### 4.6.9 (2018-07-31)


#### Bug Fixes

* **network:**  scan wireless ([5de234e6](https://github.com/linuxdeepin/dde-dock/commit/5de234e6cb5727c3b92bdaed05014f687a758955))



<a name="4.6.8"></a>
### 4.6.8 (2018-07-31)


#### Bug Fixes

*   close preview window ([97c9d45c](https://github.com/linuxdeepin/dde-dock/commit/97c9d45c164da93bc93025993584310cd0dd7367))
*   position and animation when show or hide ([96a70d76](https://github.com/linuxdeepin/dde-dock/commit/96a70d760eb9e84832c489e73a3c74ea2aaff80e))
* **network:**  refresh wireless list ([590652c7](https://github.com/linuxdeepin/dde-dock/commit/590652c709df3cc6673d946cba6cf35d0076ab1b))
* **plugin:**  unified popuptip style of container plugin ([ac4b76c3](https://github.com/linuxdeepin/dde-dock/commit/ac4b76c3e17f9d350e047261bf3ebbbabb98fa77))
* **preview:**  refresh preview snapshot ([ef2dc365](https://github.com/linuxdeepin/dde-dock/commit/ef2dc365dac4f5b1dea4b66fe9c6b41f50a579d3))

#### Features

*   add dock settings instance ([4dcb1d9c](https://github.com/linuxdeepin/dde-dock/commit/4dcb1d9c4f46b48c204e2c76720792d2b1c9957f))



<a name="4.6.7"></a>
### 4.6.7 (2018-07-19)


#### Bug Fixes

*   active connection info invalid ([e55d595d](https://github.com/linuxdeepin/dde-dock/commit/e55d595d1a4ee793c209fc57fbc929aad6ddf88c))
*   wireless crush ([af5b7bfa](https://github.com/linuxdeepin/dde-dock/commit/af5b7bfa98076f43b63f51877f03647d8cbaab41))
*   align of secret and open ap ([d09db529](https://github.com/linuxdeepin/dde-dock/commit/d09db5293c069b22820df9afcf49a1e1bffc7a3d))
*   load dde-network-utils translator ([3a10ce67](https://github.com/linuxdeepin/dde-dock/commit/3a10ce6773d35b5e4815825bb7122c0e8ba68245))
*   can not hide network plugin ([eb79f17b](https://github.com/linuxdeepin/dde-dock/commit/eb79f17b0b37a476817f2cba1ecf8d77cffa441f))
* **Icon:**  error intercepting the icon name ([6ec510ac](https://github.com/linuxdeepin/dde-dock/commit/6ec510ac12c321584605b2d305e49f043e45a6cd))
* **Network:**  Control bar layout disorder ([f73a4214](https://github.com/linuxdeepin/dde-dock/commit/f73a4214b99ce7eaeb8bc49616beb3d7009ad559))
* **network:**  crash with list is empty ([09bbd59a](https://github.com/linuxdeepin/dde-dock/commit/09bbd59adbd0ba7b4151dcf3fd8ddba9a50401eb))
* **network-plugin:**  wireless strength update ([60561a11](https://github.com/linuxdeepin/dde-dock/commit/60561a11394beaccf8e1fd881ff212f8516f7fc1))
* **sound:**
  *  not update when volume changed ([2e87a062](https://github.com/linuxdeepin/dde-dock/commit/2e87a062688ad1f6dfb877dd4ff415e38570ac1c))
  *  miss volume icon ([7547beb5](https://github.com/linuxdeepin/dde-dock/commit/7547beb578f3a8da93a5a43139b75dbde1368f9a))
  *  error icon find ([bd49207b](https://github.com/linuxdeepin/dde-dock/commit/bd49207b49a27dcfe6187017c8b321f5f391e0d7))
  *  refresh icon ([7023154b](https://github.com/linuxdeepin/dde-dock/commit/7023154becc0305010c085af21b2c381740e42de))
* **wireless:**
  *  activated ap ([eb9ea570](https://github.com/linuxdeepin/dde-dock/commit/eb9ea570ff2e3091a9e948d793cc43ca0a488638))
  *  ap property when received added,changed,removed signals ([b2ff74d2](https://github.com/linuxdeepin/dde-dock/commit/b2ff74d2f33f68a4e7e340e626c5586a0e76f816))

#### Performance

*   reduce memory usage by introducing cache ([2c0b14e4](https://github.com/linuxdeepin/dde-dock/commit/2c0b14e41befb42cecb39679ffdb94c0c5b47c6e))
* **mem:**  avoid some extra QImage copy ([f6c6a0e7](https://github.com/linuxdeepin/dde-dock/commit/f6c6a0e700cf4dc4827e991e24d753a07da51c5e))
* **network:**  update info instead of recreating ([c6e33dc3](https://github.com/linuxdeepin/dde-dock/commit/c6e33dc3e1d97dcd3f73beb3f7b575ad55be532d))

#### Features

* **TipsWidget:**  use one qss file ([052b6b29](https://github.com/linuxdeepin/dde-dock/commit/052b6b29d2a3d1c41ebee2f1505223fa04442f61))



<a name="4.6.6"></a>
### 4.6.6 (2018-06-07)


#### Features

* **tray:**  hide layout when zh_CN locale ([3486d0c1](https://github.com/linuxdeepin/dde-dock/commit/3486d0c108fe3ed8089dda2ef7629481c95ec104))



<a name="4.6.5"></a>
### 4.6.5 (2018-05-31)


#### Bug Fixes

* **tray:**  remove red background ([6514110b](https://github.com/linuxdeepin/dde-dock/commit/6514110b1dea173ac4989c385767b1e47f882cf9))



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
