package inputdevices

import (
	"encoding/json"
	"sync"
)

type appLayoutConfig struct {
	mu      sync.Mutex
	Layouts []string
	Map     map[string]int
}

func (kbd *Keyboard) loadAppLayoutConfig() error {
	mapStr := kbd.setting.GetString(kbdKeyAppLayoutMap)
	if mapStr == "" {
		return nil
	}
	err := json.Unmarshal([]byte(mapStr), &kbd.appLayoutCfg)
	if err != nil {
		return err
	}
	return nil
}

func (c *appLayoutConfig) set(app, layout string) (changed bool) {
	c.mu.Lock()
	changed = c.setNoLock(app, layout)
	c.mu.Unlock()
	return
}

func (c *appLayoutConfig) setNoLock(app, layout string) (changed bool) {
	idx := c.getLayoutIndex(layout)
	if idx == -1 {
		c.Layouts = append(c.Layouts, layout)
		idx = len(c.Layouts) - 1
		changed = true
	}

	if !changed {
		oldIdx, ok := c.Map[app]
		if ok {
			if oldIdx != idx {
				changed = true
			}
		} else {
			changed = true
		}
	}
	c.Map[app] = idx
	return
}

func (c *appLayoutConfig) toJson() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return toJSON(c)
}

func (c *appLayoutConfig) getLayoutIndex(layout string) int {
	idx := -1
	for i, l := range c.Layouts {
		if l == layout {
			idx = i
			break
		}
	}
	return idx
}

func (c *appLayoutConfig) deleteLayout(layout string) (changed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	idx := c.getLayoutIndex(layout)
	if idx == -1 {
		return false
	}

	for app, i := range c.Map {
		if i == idx {
			delete(c.Map, app)
		}
	}
	m := c.toMap()
	c.fromMap(m)
	return true
}

func (c *appLayoutConfig) toMap() map[string]string {
	result := make(map[string]string)
	for app, idx := range c.Map {
		if 0 <= idx && idx < len(c.Layouts) {
			result[app] = c.Layouts[idx]
		}
	}
	return result
}

func (c *appLayoutConfig) fromMap(m map[string]string) {
	c.Layouts = nil
	for key := range c.Map {
		delete(c.Map, key)
	}
	for key, value := range m {
		c.setNoLock(key, value)
	}
}

func (c *appLayoutConfig) get(app string) (layout string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var idx int
	idx, ok = c.Map[app]
	if !ok {
		return
	}

	if 0 <= idx && idx < len(c.Layouts) {
		ok = true
		layout = c.Layouts[idx]
	} else {
		return
	}
	return
}
