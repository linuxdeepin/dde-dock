package main

import (
	"testing"
)

func TestParseTitle(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{`menuentry 'LinuxDeepin GNU/Linux' --class linux $menuentry_id_option 'gnulinux-simple'`, `LinuxDeepin GNU/Linux`},
		{`  menuentry 'LinuxDeepin GNU/Linux' --class linux`, `LinuxDeepin GNU/Linux`},
		{`submenu 'Advanced options for LinuxDeepin GNU/Linux'`, `Advanced options for LinuxDeepin GNU/Linux`},
		{``, ``},
	}
	grub := &Grub2{}
	for _, c := range tests {
		got, _ := grub.parseTitle(c.s)
		if got != c.want {
			t.Errorf("parseTitle(%q) == %q, want %q", c.s, got, c.want)
		}
	}
}

func TestParseEntries(t *testing.T) {
	grub := &Grub2{}
	testMenuContent := `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
	load_video
}
submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
	recordfail
		echo	'载入 Linux 3.11.0-15-generic ...'
	}
    submenu 'Inner submenu for test' {
    	menuentry 'Menuentry in Level 3' {
        }
    }
}
menuentry 'Other OS' {
}
`
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test>Menuentry in Level 3`,
		`Other OS`,
	}

	grub.parseEntries(testMenuContent)
	if len(grub.entries) != len(wantEntyTitles) {
		t.Errorf("entriesCount == %v, want %v", len(grub.entries), len(wantEntyTitles))
	}
	for i, entry := range grub.entries {
		if entry.getFullTitle() != wantEntyTitles[i] {
			t.Errorf("entryTitle == %q, want %q", entry.getFullTitle(), wantEntyTitles[i])
		}
	}
}

func TestParseSettings(t *testing.T) {
	grub := &Grub2{}
	testConfigContent := `
# comment line
GRUB_DEFAULT="0"
GRUB_HIDDEN_TIMEOUT="0"
GRUB_HIDDEN_TIMEOUT_QUIET="true"
# comment line
GRUB_TIMEOUT="10"
GRUB_GFXMODE="1024x768"
  GRUB_BACKGROUND=/boot/grub/background.png
  GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`
	wantSettingCount := 7
	wantDefaultEntry := ""
	wantTimeout := "10"
	wantTheme := "/boot/grub/themes/demo/theme.txt"

	grub.parseSettings(testConfigContent)

	settingCount := len(grub.settings)
	if settingCount != wantSettingCount {
		t.Errorf("settingCount == %v, want %v", settingCount, wantSettingCount)
	}

	defaultEntry := grub.settings["GRUB_DEFAULT"]
	if defaultEntry != wantDefaultEntry {
		t.Errorf("defaultEntry == %q, want %q", defaultEntry, wantDefaultEntry)
	}

	timeout := grub.settings["GRUB_TIMEOUT"]
	if timeout != wantTimeout {
		t.Errorf("timeout == %q, want %q", timeout, wantTimeout)
	}

	theme := grub.settings["GRUB_THEME"]
	if theme != wantTheme {
		t.Errorf("theme == %q, want %q", theme, wantTheme)
	}
}

func TestGrub2(t *testing.T) {
	testMenuContent := `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
}
submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
	recordfail
		echo	'载入 Linux 3.11.0-15-generic ...'
	}
`
	testConfigContent := `
# comment line
GRUB_DEFAULT="0"
GRUB_HIDDEN_TIMEOUT="0"
GRUB_HIDDEN_TIMEOUT_QUIET="true"
# comment line
GRUB_TIMEOUT="10"
GRUB_GFXMODE="1024x768"
  GRUB_BACKGROUND=/boot/grub/background.png
  GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`

	grub := &Grub2{}
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	wantEntryCount := 2
	entriesCount := len(grub.GetEntryTitles())
	if entriesCount != wantEntryCount {
		t.Errorf("entriesCount == %v, want %v", entriesCount, wantEntryCount)
	}

	// default entry TODO
	wantDefaultEntry := `LinuxDeepin GNU/Linux`
	defaultEntry := grub.getDefaultEntry()
	if defaultEntry != wantDefaultEntry {
		t.Errorf("defaultEntry == %q, want %q", defaultEntry, wantDefaultEntry)
	}
	wantDefaultEntry = `Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`
	grub.setDefaultEntry(wantDefaultEntry)
	defaultEntry = grub.getDefaultEntry()
	if defaultEntry != wantDefaultEntry {
		t.Errorf("defaultEntry == %q, want %q", defaultEntry, wantDefaultEntry)
	}

	// timeout
	wantTimeout := int32(10)
	timeout := grub.getTimeout()
	if timeout != wantTimeout {
		t.Errorf("timeout == %v, want %v", timeout, wantTimeout)
	}
	wantTimeout = int32(15)
	grub.setTimeout(wantTimeout)
	timeout = grub.getTimeout()
	if timeout != wantTimeout {
		t.Errorf("timeout == %v, want %v", timeout, wantTimeout)
	}

	// gfxmode
	wantGfxmode := "1024x768"
	gfxmode := grub.getGfxmode()
	if gfxmode != wantGfxmode {
		t.Errorf("gfxmode == %q, want %q", gfxmode, wantGfxmode)
	}
	wantGfxmode = "saved"
	grub.setGfxmode(wantGfxmode)
	gfxmode = grub.getGfxmode()
	if gfxmode != wantGfxmode {
		t.Errorf("gfxmode == %q, want %q", gfxmode, wantGfxmode)
	}

	// background
	wantBackground := "/boot/grub/background.png"
	background := grub.getBackground()
	if background != wantBackground {
		t.Errorf("background == %q, want %q", background, wantBackground)
	}
	wantBackground = "another_background.png"
	grub.setBackground(wantBackground)
	background = grub.getBackground()
	if background != wantBackground {
		t.Errorf("background == %q, want %q", background, wantBackground)
	}

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	theme := grub.getTheme()
	if theme != wantTheme {
		t.Errorf("theme == %q, want %q", theme, wantTheme)
	}
	wantTheme = "another_theme.txt"
	grub.setTheme(wantTheme)
	theme = grub.getTheme()
	if theme != wantTheme {
		t.Errorf("theme == %q, want %q", theme, wantTheme)
	}
}

func TestSaveDefaultSettings(t *testing.T) {
	testMenuContent := `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
}
`
	testConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
`
	wantConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="LinuxDeepin GNU/Linux"
GRUB_TIMEOUT="5"
GRUB_GFXMODE="auto"
`
	grub := &Grub2{}
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)
	configContent := grub.getSettingContentToSave()
	if configContent != wantConfigContent {
		t.Errorf("configContent == %s, want %s", configContent, wantConfigContent)
	}
}

func TestSaveSettings(t *testing.T) {
	testMenuContent := `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
}
`
	testConfigContent := `GRUB_DEFAULT="0"
GRUB_TIMEOUT="10"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_GFXMODE="1024x768"
`
	wantConfigContent := `GRUB_DEFAULT="LinuxDeepin GNU/Linux"
GRUB_TIMEOUT="15"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_GFXMODE="auto"
GRUB_BACKGROUND="/boot/grub/background.png"
GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`

	grub := &Grub2{}
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	grub.setDefaultEntry(`LinuxDeepin GNU/Linux`)
	grub.setTimeout(15)
	grub.setGfxmode("auto")
	grub.setBackground("/boot/grub/background.png")
	grub.setTheme("/boot/grub/themes/demo/theme.txt")
	configContent := grub.getSettingContentToSave()
	if configContent != wantConfigContent {
		t.Errorf("configContent == %s, want %s", configContent, wantConfigContent)
	}
}
