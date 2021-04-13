#!/bin/bash

# 需要先安装lcov，打开./tests/CMakeLists.txt 测试覆盖率的编译条件
# 将该脚本放置到dde-dock-unit_test二进制文件同级目录运行

workdir=.
executable=dde_dock_unit_test
build_dir=$workdir
result_coverage_dir=$build_dir/coverage
result_report_dir=$build_dir/report/report.xml
$build_dir/$executable --gtest_output=xml:$result_report_dir

# 剔除无效信息
lcov -d $build_dir -c -o $build_dir/coverage.info -o $build_dir/coverage.info
lcov --extract $build_dir/coverage.info '*/frame/*' '*/widgets/*' -o $build_dir/coverage.info
lcov --remove $build_dir/coverage.info '*/tests/*' '*/dbus/*' '*/xcb/*' -o $build_dir/coverage.info

lcov --list-full-path -e $build_dir/coverage.info –o $build_dir/coverage-stripped.info
genhtml -o $result_coverage_dir $build_dir/coverage.info
nohup x-www-browser $result_coverage_dir/index.html &
#nohup x-www-browser $result_report_dir &
lcov -d $build_dir –z
exit 0