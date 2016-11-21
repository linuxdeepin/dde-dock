**dbus_gen**: deepin 网络后端 Python DBus 接口库, 由 gen_dbus_code.sh 生成

**gen_dbus_code.sh**: 用于生成 deepin 网络后端 Python DBus 接口库, 需要安
装 dbus2any 模块并下载模板文件
  ```sh
  sudo pip3 install dbus2any'
  dbus2any_tpl_dir=/usr/lib/python3.5/site-packages/dbus2any/templates
  # or dbus2any_tpl_dir=/usr/local/lib/python3.5/dist-packages/dbus2any/templates
  sudo mkdir ${dbus2any_tpl_dir}
  curl https://raw.githubusercontent.com/hugosenari/dbus2any/master/dbus2any/templates/pydbusclient.tpl | sudo tee ${dbus2any_tpl_dir}/pydbusclient.tpl
  ```

**main.py**: 一个 Python3 示例程序, 演示如何调用 deepin 网络后端 DBus
  接口来连接 WiFi 网络, 完整的自动化测试用例请参考 [deepin-network-tests](https://github.com/x-deepin/deepin-network-tests/)
