## Debug

`dde-session-daemon` 开启 `debug` 模式的方式有多种, 本文将举例介绍, 如下:

* `gsettings` 开启
    
    通过执行 `gsetiings set com.deepin.dde.daemon debug true` 开启 `debug` 模块, 设置为 `false` 就关闭此模块. 它会动态调节日志级别, 另外程序启动时也会读取这个值, 如果为 `true` 那程序启动时的日志级别就是 debug`.
    
    `com.deepin.dde.daemon` 这个 `schema` 下是所有支持的模块, 如果不想启动某个模块就将其设置为 `false`. 有些模块会被其他模块依赖, 要禁用就要一起禁用.
    
    另外 `debug` 模块还会打开 `pprof` 调试功能, 详情见下文.

* 命令行开启
    
    通过执行 `dde-session-daemon -v` 可以开启所有模块的 `debug` 日志. `dde-session-daemon` 也支持只启动某个模块, 命令如下 `dde-session-daemon enable <module name> -v`.

* 环境变量
    
    支持两个环境变量, 如下:
    - `DDE_DEBUG_LEVEL="debug"` 设置日志级别
    - `DDE_DEBUG_MATCH=<module name>` 只显示某个模块的 `debug` 日志

---------------------------------------------

## pprof

使用 `pprof` 之前需要先安装一些依赖, `deepin` 如命令: `sudo apt-get install golang golang-go golang-src graphviz`. 然后执行 `gsetiings set com.deepin.dde.daemon debug true` 开启 `pprof http server`, 然后就可以获取 `pprof` 信息了.

可以获取的信息如下:

* 堆内存信息 `go tool pprof http://localhost:6969/debug/pprof/heap`
* 30s 内 cpu 信息 `go tool pprof http://localhost:6969/debug/pprof/profile`
* 块信息 `go tool pprof http://localhost:6969/debug/pprof/block`

上面的命令执行后, 会进入一个 `shell` 输入 `pdf` 将结果保存为 `pdf` 文档来查看.

通过执行 `wget http://localhost:6969/debug/pprof/trace?seconds=5 > trace.out` 可以追踪 5s 内的调用操作, 然后使用 `go tool trace -http=':8080' trace trace.out` 来查看.

在浏览器里打开 `http://localhost:6969/debug/pprof/` 可以看到所有可用的 `profile` 的信息.

获取上述信息, 剩下的就是分析了. 另外最后也同时保存下 `dde-session-daemon` 的日志, 即在开启 `debug` 模块后执行 `journactl -f /usr/lib/deepin-daemon/dde-session-daemon > /tmp/daemon.log` 将日志保存到 文件 `/tmp/daemon.log` 里. 
