package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"pkg.deepin.io/lib/dbus1"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	for _, env := range os.Environ() {
		log.Println(env)
	}
	uid := os.Getuid()
	log.Println("uid:", uid)

	user := os.Getenv("PAM_USER")
	if user == "" {
		log.Fatal("user empty")
	}

	token, err := bufio.NewReader(os.Stdin).ReadString(0)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
	token = strings.TrimRight(token, "\x00\n")

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}

	authObj := conn.Object("com.deepin.daemon.Authority", "/com/deepin/daemon/Authority")
	var ok bool
	err = authObj.Call("com.deepin.daemon.Authority.CheckCookie", dbus.FlagNoAutoStart,
		user, token).Store(&ok)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("auth failed")
	}

	log.Println("auth success")
	return
}
