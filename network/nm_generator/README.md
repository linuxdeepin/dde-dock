**描述**: nm_generator 通过读取 NM.gir(隶属于 libnm-dev 包) 莱生成一些
辅助 Go 代码以方便 deepin network 模块通过 DBus 接口执行
NetworkManager 各操作.

## 目录结构

- **nm_docs**: NetworkManager DBus 接口 html 文档, 通过编译源码获取.
  之所以把这些接口文档复制过来, 一方面是这些文档没有归属任何包里, 不易
  获取, 另一方面是为了每次升级时对比文档的变更内容来了解上游详细的升级
  内容.

  ```
  $ cd $nm-src
  $ ./autogen.sh --enable-doc --enable-gtk-doc
  $ make
  $ ls docs/api/html/
  ```

- **nm_girs**: NM.gir 文件集合, 用于对比各版本的接口变更.

- **gen_nm_consts.py**: 通过 NM.gir 生成 nm_consts_gen.yml 文件.

- **nm_consts_gen.yml**: 通过 gen_nm_consts.py 读取 NM.gir 生成, 包含
  NetworkManager 的所有常量定义和编辑网络连接所用到的原始字段(setting)
  和键值(key, 包括键值的类型和默认值).

- **nm_consts_keys_override.yml**: 重载定义 nm_consts_gen.yml 里的部分
  键值, 主要是用来包装一下键值类型从而可以方便前端进行展现. 比如
  NM_SETTING_WIRELESS_SSID 默认为 ktypeArrayByte 类型, 目前包装成
  ktypeWrapperString 字符串类型, 这样在前端展现的时候就是普通的字符串,
  后端设置的时候会将其还原回原本的 byte 数组.

- **nm_logicset_keys.yml**: 定义默认需要开发者手动实现 setter/getter
  (即 logic setter/getter) 的键值列表. 这些键值变更时一般包含逻辑关系,
  比如 NM_SETTING_IP4_CONFIG_METHOD 设置获取 IP 地址类型时(静态/自动),
  还会添加或删除从属的键值, 故需要开发者手动编写相关逻辑代码.

- **nm_virtual_sections.yml**: 定义前端编辑网络连接是需要展示的虚拟字
  段(virtual section, 对应原始 setting)和虚拟键值(virtual key, 对应原
  始 key). 为了方便处理相关逻辑关系, 前端网络编辑页面基本是通过后端的
  接口动态生成的, 所以后端需要提供整套的解决方案包括对应的数据结构, 即
  重新定义虚拟字段和虚拟键值以涵盖 NetworkManager 的默认 setting 和
  key 不能完全满足的前端展示需求.

  简单的讲, 该文件定义的键值数据完全对应前端进入编辑页面所展现的 UI,
  包括所显示的国际化文本和控件的类型和次序, 所以要增删前端控件只需要在
  对应位置编辑该文件即可. 需要注意因为某些虚拟键可能会在多个虚拟字段出
  现, 其 VKeyInfo 只需定义一次即可, 如
  NM_SETTING_VK_VPN_VPNC_KEY_HYBRID_AUTHMODE.

- **nm_vpn_alias_settings.yml**: 定义原始字段 (setting) 的别名 (alias
  name), 主要用于定义各 VPN 类型对应的字段和键值, 因为所有 VPN 键值数
  据均保存在 NM_SETTING_VPN_SETTING_NAME 字段下的 NM_SETTING_VPN_DATA
  键下 (ktypeDictStringString 字典类型), 即它们共用一个键
  (NetworkManager 这样设计是为了扩展更多的 VPN 类型, 而每个 VPN 所需要
  的配置由它们自己处理), 所以需要为其添加别名进行分辨.

- **main.go**: nm_generator Go 主程序.

- **tpl.go**: 生成 Go 代码所用的模板.

- **utils.go**: 一些辅助函数.

## 运行 nm_generator

```sh
$ make gen-nm-code
```

## 特别注意

每次重新生成 nm_setting_beans_gen.go 文件后，需要手工给 generalSetSettingKeyJSON 打补丁，
否则就会出现想将 MAC 地址由非空设置为空，但不能成功的问题。

相关 Change-Id: Ie1ff2b980b601b04866d9b46f119da7d810b69da