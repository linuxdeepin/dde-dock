package item

import (
	"os"
	"path"

	C "launchpad.net/gocheck"
	"pkg.linuxdeepin.com/lib/utils"
)

type DesktopTestSuite struct {
	oldHome string
	testDataDir string
}

var _ = C.Suite(&DesktopTestSuite{})

// according to the sources of glib.
func (s *DesktopTestSuite) SetUpSuite(c *C.C) {
	s.oldHome = os.Getenv("HOME")
	s.testDataDir = "../testdata"
	os.Setenv("HOME", s.testDataDir)
}

func (s *DesktopTestSuite) TearDownSuite(c *C.C) {
	os.Setenv("HOME", s.oldHome)
}

func (s *DesktopTestSuite) TestgetDesktopPath(c *C.C) {
	c.Assert(getDesktopPath("firefox.desktop"), C.Equals, path.Join(s.testDataDir, "Desktop/firefox.desktop"))
}

func (s *DesktopTestSuite) TestisOnDesktop(c *C.C) {
	c.Assert(isOnDesktop("firefox.desktop"), C.Equals, true)
	c.Assert(isOnDesktop("google-chrome.desktop"), C.Equals, false)
}

func (s *DesktopTestSuite) TestSendRemoveDesktop(c *C.C) {
	srcFile := path.Join(s.testDataDir, "deepin-software-center.desktop")
	destFile := path.Join(s.testDataDir, "Desktop/deepin-software-center.desktop")
	sendToDesktop(srcFile)
	c.Assert(utils.IsFileExist(destFile), C.Equals, true)

	st, err := os.Lstat(destFile)
	if err != nil {
		c.Skip(err.Error())
	}

	var execPerm os.FileMode = 0100
	c.Assert(st.Mode().Perm()&execPerm, C.Equals, execPerm)

	removeFromDesktop(srcFile)
	c.Assert(utils.IsFileExist(destFile), C.Equals, false)
}
