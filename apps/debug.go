// +build debug

package apps

import (
	"pkg.deepin.io/lib/dbus"
)

func (r *ALRecorder) DebugUserRemoved(dMsg dbus.DMessage) {
	uid := int(dMsg.GetSenderUID())
	r.handleUserRemoved(uid)
}
