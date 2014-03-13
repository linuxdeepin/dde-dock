package main

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

type EntryProxyer struct {
	core *RemoteEntry
	ID   string `dmusic`
	Type string `applet/other`

	Tooltip string
	Icon    string

	Status int32 `Actived/Normal/`

	QuickWindowVieable bool
	Allocation         Rectangle
}

func NewEntryProxyer(entryId string) *EntryProxyer {
	e := &EntryProxyer{}
	core, err := NewRemoteEntry("dde.dock.entry."+entryId, "/dde/dock/entry/v1")
	if err != nil {
		return nil
	}
	e.core = core
	//TODO: init properties
	//TODO: monitor properties changed
	return e
}

func (e *EntryProxyer) QuickWindow(x, y int32)              { e.core.QuickWindow(x, y) }
func (e *EntryProxyer) ContextMenu(x, y int32)              { e.core.ContextMenu(x, y) }
func (e *EntryProxyer) Activate(x, y int32)                 { e.core.Activate(x, y) }
func (e *EntryProxyer) SecondaryActivate(x, y int32)        { e.core.SecondaryActivate(x, y) }
func (e *EntryProxyer) OnDragEnter(x, y int32, data string) { e.core.OnDragEnter(x, y, data) }
func (e *EntryProxyer) OnDragLeave(x, y int32, data string) { e.core.OnDragLeave(x, y, data) }
func (e *EntryProxyer) OnDragOver(x, y int32, data string)  { e.core.OnDragOver(x, y, data) }
func (e *EntryProxyer) OnDragDrop(x, y int32, data string)  { e.core.OnDragDrop(x, y, data) }
