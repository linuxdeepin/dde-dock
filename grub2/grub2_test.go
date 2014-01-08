package main

import (
// "testing"
)

// func TestParseTitle(t *testing.T) {
// 	var tests = []struct {
// 		s, want string
// 	}{
// 		{`menuentry 'LinuxDeepin GNU/Linux' --class linux`, `LinuxDeepin GNU/Linux`},
// 		{`  menuentry 'LinuxDeepin GNU/Linux' --class linux`, `LinuxDeepin GNU/Linux`},
// 		{`submenu 'Advanced options for LinuxDeepin GNU/Linux'`, ``},
// 		{``, ``},
// 	}
// 	grub := &Grub2{}
// 	for _, c := range tests {
// 		got, _ := grub.parseTitle(c.s)
// 		if got != c.want {
// 			t.Errorf("parseTitle(%q) == %q, want %q", c.s, got, c.want)
// 		}
// 	}
// }

// func TestParseEntries(t *testing.T) {
// 	grub := &Grub2{}
// 	testContent := `
// menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
// recordfail
// }
// submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
// 	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
// 	recordfail
// 		echo	'载入 Linux 3.11.0-15-generic ...'
// 	}
// `
// 	wantEntryCount := 2
// 	wantEntryOne := `LinuxDeepin GNU/Linux`
// 	wantEntryTwo := `Advanced options for LinuxDeepin GNU/Linux`

// 	grub.parseEntries(testContent)
// 	entriesCount := len(grub.entries)
// 	if entriesCount != wantEntryCount {
// 		t.Errorf("entriesCount == %q, want %q", entriesCount, wantEntryCount)
// 	}
// 	if grub.entries[0] != wantEntryOne {
// 		t.Errorf("entryOne == %q, want %q", grub.entries[0], wantEntryOne)
// 	}
// 	if grub.entries[1] != wantEntryTwo {
// 		t.Errorf("entryTwo == %q, want %q", grub.entries[1], wantEntryTwo)
// 	}
// }

// func TestParseSettings(t *testing.T) {
// 	grub := &Grub2{}
// 	testContent := `
// # comment line
// GRUB_DEFAULT="0"
// GRUB_HIDDEN_TIMEOUT="0"
// GRUB_HIDDEN_TIMEOUT_QUIET="true"
// # comment line
// GRUB_TIMEOUT="10"
// GRUB_GFXMODE="1024x768"
// GRUB_BACKGROUND=/boot/grub/background.png
// GRUB_THEME="/boot/grub/themes/demo/theme.txt"
// `
// 	wantSettingCount := 7
// 	wantDefault := 0
// 	wantTimeout := 10
// 	wantGfxmode := "1024x768"
// 	wantBackground := "/boot/grub/background.png"
// 	wantTheme := "/boot/grub/themes/demo/theme.txt"

// 	grub.parseSettings(testContent)
// 	// TODO
// }
