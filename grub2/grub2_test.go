package main

import (
	"testing"
)

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
	entries := grub.GetEntries()
	entriesCount := len(entries)
	if entriesCount != wantEntryCount {
		t.Errorf("entriesCount == %v, want %v", entriesCount, wantEntryCount)
	}

	// default entry
	wantDefaultEntry := uint32(0)
	defaultEntry := grub.GetDefaultEntry()
	if defaultEntry != wantDefaultEntry {
		t.Errorf("defaultEntry == %v, want %v", defaultEntry, wantDefaultEntry)
	}
	wantDefaultEntry = uint32(2)
	grub.SetDefaultEntry(wantDefaultEntry)
	defaultEntry = grub.GetDefaultEntry()
	if defaultEntry != wantDefaultEntry {
		t.Errorf("defaultEntry == %v, want %v", defaultEntry, wantDefaultEntry)
	}

	// timeout
	wantTimeout := int32(10)
	timeout := grub.GetTimeout()
	if timeout != wantTimeout {
		t.Errorf("timeout == %v, want %v", timeout, wantTimeout)
	}
	wantTimeout = int32(15)
	grub.SetTimeout(wantTimeout)
	timeout = grub.GetTimeout()
	if timeout != wantTimeout {
		t.Errorf("timeout == %v, want %v", timeout, wantTimeout)
	}

	// gfxmode
	wantGfxmode := "1024x768"
	gfxmode := grub.GetGfxmode()
	if gfxmode != wantGfxmode {
		t.Errorf("gfxmode == %q, want %q", gfxmode, wantGfxmode)
	}
	wantGfxmode = "saved"
	grub.SetGfxmode(wantGfxmode)
	gfxmode = grub.GetGfxmode()
	if gfxmode != wantGfxmode {
		t.Errorf("gfxmode == %q, want %q", gfxmode, wantGfxmode)
	}

	// background
	wantBackground := "/boot/grub/background.png"
	background := grub.GetBackground()
	if background != wantBackground {
		t.Errorf("background == %q, want %q", background, wantBackground)
	}
	wantBackground = "another_background.png"
	grub.SetBackground(wantBackground)
	background = grub.GetBackground()
	if background != wantBackground {
		t.Errorf("background == %q, want %q", background, wantBackground)
	}

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	theme := grub.GetTheme()
	if theme != wantTheme {
		t.Errorf("theme == %q, want %q", theme, wantTheme)
	}
	wantTheme = "another_theme.txt"
	grub.SetTheme(wantTheme)
	theme = grub.GetTheme()
	if theme != wantTheme {
		t.Errorf("theme == %q, want %q", theme, wantTheme)
	}
}

func TestSaveSettings(t *testing.T) {
	testConfigContent := `GRUB_DEFAULT="0"
GRUB_TIMEOUT="10"
GRUB_GFXMODE="1024x768"
`
	wantConfigContent := `GRUB_DEFAULT="1"
GRUB_TIMEOUT="15"
GRUB_GFXMODE="saved"
GRUB_BACKGROUND="/boot/grub/background.png"
GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`

	grub := &Grub2{}
	grub.parseSettings(testConfigContent)

	grub.SetDefaultEntry(1)
	grub.SetTimeout(15)
	grub.SetGfxmode("saved")
	grub.SetBackground("/boot/grub/background.png")
	grub.SetTheme("/boot/grub/themes/demo/theme.txt")
	configContent := grub.getSettingContentToSave()
	if configContent != wantConfigContent {
		t.Errorf("configContent == %s, want %s", configContent, wantConfigContent)
	}
}
