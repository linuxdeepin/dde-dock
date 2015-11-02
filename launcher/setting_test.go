package launcher

import (
	C "launchpad.net/gocheck"
	"pkg.deepin.io/dde/daemon/launcher/mock"
	. "pkg.deepin.io/dde/daemon/launcher/setting"
)

type SettingTestSuite struct {
	s                               *Settings
	CategoryDisplayModeChangedCount int64
	SortMethodChangedCount          int64
}

// FIXME: gsetting cannot be mocked, because the signal callback must have type func(*gio.Settings, string)
// var _ = C.Suite(&SettingTestSuite{})

func (sts *SettingTestSuite) SetUpTest(c *C.C) {
	var err error
	core := mock.NewSettingCore()
	sts.s, err = NewSettings(core)

	sts.CategoryDisplayModeChangedCount = 0
	sts.s.CategoryDisplayModeChanged = func(int64) {
		sts.CategoryDisplayModeChangedCount++
	}

	sts.SortMethodChangedCount = 0
	sts.s.SortMethodChanged = func(int64) {
		sts.SortMethodChangedCount++
	}
	if err != nil {
		c.Fail()
	}
}

func (sts *SettingTestSuite) TestGetCategoryDisplayMode(c *C.C) {
	c.Assert(sts.s, C.NotNil)
	sts.s.GetCategoryDisplayMode()
	c.Assert(sts.s.GetCategoryDisplayMode(), C.Equals, int64(CategoryDisplayModeIcon))
}

func (sts *SettingTestSuite) TestSetCategoryDisplayMode(c *C.C) {
	c.Assert(sts.s, C.NotNil)

	c.Assert(sts.s.GetCategoryDisplayMode(), C.Equals, int64(CategoryDisplayModeIcon))

	sts.s.SetCategoryDisplayMode(int64(CategoryDisplayModeIcon))
	c.Assert(sts.s.GetCategoryDisplayMode(), C.Equals, int64(CategoryDisplayModeIcon))

	sts.s.SetCategoryDisplayMode(int64(CategoryDisplayModeText))
	c.Assert(sts.s.GetCategoryDisplayMode(), C.Equals, int64(CategoryDisplayModeText))
}

func (sts *SettingTestSuite) TestGetSortMethod(c *C.C) {
	c.Assert(sts.s, C.NotNil)
	c.Assert(sts.s.GetSortMethod(), C.Equals, int64(SortMethodByName))
}

func (sts *SettingTestSuite) TestSetSortMethod(c *C.C) {
	c.Assert(sts.s, C.NotNil)

	c.Assert(sts.s.GetSortMethod(), C.Equals, int64(SortMethodByName))

	sts.s.SetSortMethod(int64(SortMethodByName))
	c.Assert(sts.SortMethodChangedCount, C.Equals, int64(0))
	c.Assert(sts.s.GetSortMethod(), C.Equals, int64(SortMethodByName))

	sts.s.SetSortMethod(int64(SortMethodByCategory))
	c.Assert(sts.SortMethodChangedCount, C.Equals, int64(1))
	c.Assert(sts.s.GetSortMethod(), C.Equals, int64(SortMethodByCategory))
}
