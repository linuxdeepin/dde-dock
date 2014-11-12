#!/usr/bin/python

import os
import re


lang = ["zh_CN", "zh_TW","es","pt_BR"]

go_files = []
except_dirs = ["dbus_test", "memory_test", "dominant_color"]
except_files = []

source_dir = "../../"

def is_str_in_list(str, list):
    for d in list:
        if str == d:
            return True
    return False

def get_files():
    global go_files
    for root, dirs, files in os.walk(source_dir):
        for filepath in files:
            if filepath.endswith('.go'):
                if is_str_in_list(filepath, except_files):
                    continue
                go_files.append(os.path.join(root,filepath))

POT_FILE = "dde-daemon.pot"
def scan():
    global go_files, lang
    #os.system("rm %s" % POT_FILE)
    go_files.sort()
    cmd = "xgettext --from-code=utf-8 -C -kTr -o %s " % POT_FILE + " ".join(go_files)
    os.system(cmd)

if __name__ == '__main__':
    get_files()
    scan()
