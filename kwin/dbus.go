package kwin

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

var (
	dbusConnection *dbus.Conn
	dbusPath       string
	dbusInterface  string
	letters        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func Close() {
	for dbusConnection.Connected() {
		dbusConnection.Close()
		time.Sleep(200 * time.Millisecond)
	}
}

func GetDbusObject(dest string, path dbus.ObjectPath) dbus.BusObject {
	return dbusConnection.Object(dest, path)
}

func GetDbusObjectPathId(call *dbus.Call) dbus.ObjectPath {
	return dbus.ObjectPath(fmt.Sprintf("/%v", call.Body[0]))
}

// func ConnectSessionBus() (*dbus.Conn, error) {
// 	return dbus.ConnectSessionBus()
// }

func startEventListener() {
	rand.Seed(time.Now().UnixNano())
	dbusInterface = "no.hipsters.Kazinga." + randSeq(10)
	dbusPath = "/" + strings.ReplaceAll(dbusInterface, ".", "/")
	kazingaInterface := `
        <node>
            <interface name="` + dbusInterface + `">
                <method name="WindowList">
                    <arg direction="in" type="s"/>
                </method>
                <method name="DisplayInfo">
                    <arg direction="in" type="s"/>
                </method>
            </interface>` + introspect.IntrospectDataString + `
        </node>`
	var err error
	dbusConnection, err = dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	f := dbusResponse("I need to figure out why this is needed")
	dbusConnection.Export(f, dbus.ObjectPath(dbusPath), dbusInterface)
	dbusConnection.Export(introspect.Introspectable(kazingaInterface), dbus.ObjectPath(dbusPath), "org.freedesktop.DBus.Introspectable")
	reply, err := dbusConnection.RequestName(dbusInterface, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
}

func Wrap(method, data string) string {
	return `callDBus("` + dbusInterface + `","` + dbusPath + `","` + dbusInterface + `","` + method + `", ` + data + `);`
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
