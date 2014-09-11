#!/usr/bin/python

import glob
import os


lang = ["zh_CN", "zh_TW","es","pt_BR"]

go_files = []
except_dirs = ["dbus_test", "memory_test", "dominant_color"]


def get_files():
    global go_files
    #for d in os.listdir("../app"):
        #if d not in except_dirs:
            #go_files += glob.glob("../app/%s/*.c" % d);
            #go_files += glob.glob("../app/%s/*.h" % d);

    go_files += [ "../../*/*.go" ]


POT_FILE = "dde-daemon.pot"
def scan():
    global go_files, lang
    #os.system("rm %s" % POT_FILE)
    cmd = "xgettext --from-code=utf-8 -C -kTr -o %s " % POT_FILE + " ".join(go_files)
    os.system(cmd)

if __name__ == '__main__':
    get_files()
    scan()
