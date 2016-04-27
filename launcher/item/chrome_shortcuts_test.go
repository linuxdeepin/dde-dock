/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package item

import (
	. "github.com/smartystreets/goconvey/convey"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"testing"
)

func TestItemIsChromeShortcut(t *testing.T) {
	itemA := &Info{
		id: ItemID("chrome-nmmhkkegccagdldgiimedpiccmgmieda-Default"),
		xinfo: Xinfo{
			exec: "/opt/google/chrome/google-chrome --profile-directory=Default --app-id=nmmhkkegccagdldgiimedpiccmgmieda",
		},
	}

	itemB := &Info{
		id: ItemID("chrome-eejlmmonjkhhegpahogmbbpeoagmkcih-Default"),
		xinfo: Xinfo{
			exec: "/opt/google/chrome-unstable/google-chrome-unstable --user-data-dir=/home/tp/.config/google-chrome-unstable --profile-directory=Default --app-id=eejlmmonjkhhegpahogmbbpeoagmkcih",
		},
	}

	itemC := &Info{
		id: ItemID("firefox"),
		xinfo: Xinfo{
			exec: "firefox",
		},
	}

	Convey("Launcher/Chrome apps", t, func() {
		So(itemIsChromeShortcut(itemA), ShouldBeTrue)
		So(itemIsChromeShortcut(itemB), ShouldBeTrue)
		So(itemIsChromeShortcut(itemC), ShouldBeFalse)
	})
}
